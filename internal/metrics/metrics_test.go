package metrics

import (
	"testing"
	"time"

	"github.com/alexandrmotologa/uxswarm/internal/config"
)

func TestCalculateSUS_PerfectScore(t *testing.T) {
	// Perfect score: 5 on odd, 1 on even
	// Odd: 5-1 = 4. Even: 5-1 = 4. Sum = 40. 40 * 2.5 = 100.
	responses := [10]int{5, 1, 5, 1, 5, 1, 5, 1, 5, 1}
	res, err := CalculateSUS(responses)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Score != 100.0 {
		t.Errorf("expected score 100, got %.1f", res.Score)
	}
	if res.Grade != GradeA {
		t.Errorf("expected GradeA, got %s", res.Grade)
	}
}

func TestCalculateSUS_WorstScore(t *testing.T) {
	// Worst score: 1 on odd, 5 on even
	// Odd: 1-1 = 0. Even: 5-5 = 0. Sum = 0. 0 * 2.5 = 0.
	responses := [10]int{1, 5, 1, 5, 1, 5, 1, 5, 1, 5}
	res, err := CalculateSUS(responses)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Score != 0.0 {
		t.Errorf("expected score 0, got %.1f", res.Score)
	}
	if res.Grade != GradeF {
		t.Errorf("expected GradeF, got %s", res.Grade)
	}
}

func TestCalculateSUS_BenchmarkScore(t *testing.T) {
	// Average neutral responses: all 3
	// Odd: 3-1 = 2. Even: 5-3 = 2. Sum = 20. 20 * 2.5 = 50.
	responses := [10]int{3, 3, 3, 3, 3, 3, 3, 3, 3, 3}
	res, err := CalculateSUS(responses)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Score != 50.0 {
		t.Errorf("expected score 50, got %.1f", res.Score)
	}
	if res.Grade != GradeF {
		t.Errorf("expected GradeF for 50, got %s", res.Grade)
	}
}

func TestTelemetryTracker_FrictionIndexAndHesitation(t *testing.T) {
	tracker := NewTelemetryTracker()

	now := time.Now()
	// Step 1: Navigated to homepage
	tracker.RecordStep(StepEvent{
		Timestamp:         now,
		ActionType:        "navigate",
		URL:               "http://localhost:3000/",
		TargetDescription: "Home page",
	})

	// Step 2: Hesitated 4.2 seconds then clicked button
	tracker.RecordStep(StepEvent{
		Timestamp:         now.Add(4200 * time.Millisecond),
		ActionType:        "click",
		URL:               "http://localhost:3000/step1",
		TargetDescription: "Get Started button",
		HesitationSeconds: 4.2,
		X:                 200,
		Y:                 150,
	})

	// Step 3: Dead click on non-clickable heading
	tracker.RecordStep(StepEvent{
		Timestamp:         now.Add(5000 * time.Millisecond),
		ActionType:        "click",
		URL:               "http://localhost:3000/step1",
		TargetDescription: "Features heading",
		IsDeadClick:       true,
		X:                 300,
		Y:                 200,
	})

	// Step 4: Backtrack to homepage
	tracker.RecordStep(StepEvent{
		Timestamp:         now.Add(6000 * time.Millisecond),
		ActionType:        "navigate",
		URL:               "http://localhost:3000/",
		TargetDescription: "Back to Home page",
	})

	tracker.MarkCompleted()
	summary := tracker.Summary()

	if summary.TotalSteps != 4 {
		t.Errorf("expected 4 steps, got %d", summary.TotalSteps)
	}
	if summary.HesitationCount != 1 {
		t.Errorf("expected 1 hesitation, got %d", summary.HesitationCount)
	}
	if summary.DeadClicks != 1 {
		t.Errorf("expected 1 dead click, got %d", summary.DeadClicks)
	}
	if summary.Backtracks != 1 {
		t.Errorf("expected 1 backtrack, got %d", summary.Backtracks)
	}
	if summary.FrictionIndex <= 0 {
		t.Errorf("expected positive friction index, got %.2f", summary.FrictionIndex)
	}
}

func TestEvaluateHeuristics(t *testing.T) {
	persona := config.Persona{
		ID:        "test-persona",
		Name:      "Test Persona",
		Archetype: config.ArchetypeHighImpatienceLowTech,
		Viewport:  config.Viewport{Width: 390, Height: 844, IsMobile: true},
	}

	summary := &TelemetrySummary{
		TotalSteps:                3,
		TotalDurationSeconds:      12.0,
		HesitationCount:           2,
		HesitationDurationSeconds: 8.5,
		DeadClicks:                2,
		Backtracks:                1,
		GoalCompleted:             false,
		DroppedOff:                true,
		DropOffReason:             "Patience threshold exceeded on complex checkout",
		FrictionIndex:             2.8,
		MandatoryRedesign:         true,
	}

	steps := []StepEvent{
		{
			IsDeadClick:       true,
			TargetDescription: "Unlinked promo banner",
			X:                 150,
			Y:                 80,
		},
	}

	findings := EvaluateHeuristics(
		[]config.Persona{persona},
		[]*TelemetrySummary{summary},
		map[string][]StepEvent{"test-persona": steps},
	)

	if len(findings) == 0 {
		t.Fatalf("expected heuristic findings, got 0")
	}

	foundVisibility := false
	foundDeadClick := false
	for _, f := range findings {
		if f.HeuristicID == 1 {
			foundVisibility = true
		}
		if f.HeuristicID == 4 {
			foundDeadClick = true
		}
	}
	if !foundVisibility {
		t.Errorf("expected finding for Heuristic 1 (Visibility of system status)")
	}
	if !foundDeadClick {
		t.Errorf("expected finding for Heuristic 4 (Consistency and standards)")
	}
}
