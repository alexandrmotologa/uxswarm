package store

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/alexandrmotologa/uxswarm/internal/metrics"
)

func TestStore_FullLifecycle(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test_uxswarm.db")

	s, err := NewStore(dbPath)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	defer s.Close()

	// 1. Create Run
	runID := "run_test_123"
	now := time.Now()
	run := &RunRecord{
		ID:                runID,
		Name:              "Onboarding Test",
		TargetURL:         "http://localhost:3000/signup",
		Goal:              "Complete signup",
		Status:            "running",
		TotalPersonas:     2,
		CreatedAt:         now,
		MandatoryRedesign: false,
	}
	if err := s.CreateRun(run); err != nil {
		t.Fatalf("create run error: %v", err)
	}

	// 2. Create Session
	sessID := "sess_1"
	session := &SessionRecord{
		ID:              sessID,
		RunID:           runID,
		PersonaID:       "sarah-cfo",
		PersonaName:     "Sarah Chen",
		Archetype:       "HIGH_IMPATIENCE_LOW_TECH",
		ViewportWidth:   1440,
		ViewportHeight:  900,
		IsMobile:        false,
		SUSScore:        72.5,
		SUSGrade:        "C",
		FrictionIndex:   1.8,
		TotalSteps:      3,
		TotalDuration:   15.4,
		HesitationCount: 1,
		DeadClicks:      1,
		GoalCompleted:   true,
		CreatedAt:       now,
	}
	if err := s.CreateSession(session); err != nil {
		t.Fatalf("create session error: %v", err)
	}

	// 3. Create Steps
	step := &StepRecord{
		SessionID:         sessID,
		StepIndex:         1,
		Timestamp:         now,
		ActionType:        "click",
		TargetDescription: "Submit Button",
		X:                 150,
		Y:                 200,
		IsDeadClick:       false,
		HesitationSeconds: 3.8,
		URL:               "http://localhost:3000/signup",
		InnerSpeech:       "Clicking submit button.",
		Frustration:       0.3,
	}
	if err := s.CreateStep(step); err != nil {
		t.Fatalf("create step error: %v", err)
	}

	// 4. Save Heuristics
	findings := []metrics.HeuristicFinding{
		{
			HeuristicID:      1,
			HeuristicName:    "Visibility of system status",
			Severity:         metrics.SeverityMajor,
			Title:            "Hesitation latency high",
			Evidence:         "Sarah paused 3.8s",
			AffectedPersonas: []string{"Sarah Chen"},
			Recommendation:   "Add progress spinner",
		},
	}
	if err := s.SaveHeuristics(runID, findings); err != nil {
		t.Fatalf("save heuristics error: %v", err)
	}

	// 5. Update Run to completed
	run.Status = "completed"
	run.AverageSUS = 72.5
	run.AverageFriction = 1.8
	run.Grade = "C"
	run.CompletedAt = now.Add(20 * time.Second)
	if err := s.UpdateRun(run); err != nil {
		t.Fatalf("update run error: %v", err)
	}

	// 6. Verify Queries
	fetchedRun, err := s.GetRun(runID)
	if err != nil {
		t.Fatalf("get run error: %v", err)
	}
	if fetchedRun.Status != "completed" || fetchedRun.AverageSUS != 72.5 {
		t.Errorf("unexpected run state: %+v", fetchedRun)
	}

	runs, err := s.ListRuns(10)
	if err != nil || len(runs) != 1 {
		t.Fatalf("expected 1 run in list, got %d (err: %v)", len(runs), err)
	}

	sessions, err := s.ListSessionsByRun(runID)
	if err != nil || len(sessions) != 1 {
		t.Fatalf("expected 1 session, got %d (err: %v)", len(sessions), err)
	}
	if sessions[0].PersonaName != "Sarah Chen" {
		t.Errorf("expected Sarah Chen, got %s", sessions[0].PersonaName)
	}

	steps, err := s.ListStepsBySession(sessID)
	if err != nil || len(steps) != 1 {
		t.Fatalf("expected 1 step, got %d (err: %v)", len(steps), err)
	}
	if steps[0].HesitationSeconds != 3.8 {
		t.Errorf("expected 3.8s hesitation, got %.1f", steps[0].HesitationSeconds)
	}

	fetchedFindings, err := s.ListHeuristicsByRun(runID)
	if err != nil || len(fetchedFindings) != 1 {
		t.Fatalf("expected 1 finding, got %d (err: %v)", len(fetchedFindings), err)
	}
	if fetchedFindings[0].HeuristicID != 1 {
		t.Errorf("expected heuristic id 1, got %d", fetchedFindings[0].HeuristicID)
	}
}
