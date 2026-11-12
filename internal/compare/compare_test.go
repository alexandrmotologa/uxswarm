package compare

import (
	"strings"
	"testing"
	"time"

	"github.com/alexandrmotologa/uxswarm/internal/report"
	"github.com/alexandrmotologa/uxswarm/internal/store"
)

func TestCompareVariants(t *testing.T) {
	repA := &report.ReportData{
		Run: store.RunRecord{
			Name:            "Variant A (Baseline)",
			TargetURL:       "http://localhost:3000/v1",
			AverageSUS:      52.0,
			Grade:           "D",
			AverageFriction: 2.10,
			CreatedAt:       time.Now(),
		},
		Sessions: []store.SessionRecord{
			{
				PersonaID:     "sarah",
				PersonaName:   "Sarah Chen",
				Archetype:     "HIGH_IMPATIENCE_LOW_TECH",
				SUSScore:      40.0,
				FrictionIndex: 2.50,
				DroppedOff:    true,
			},
		},
	}

	repB := &report.ReportData{
		Run: store.RunRecord{
			Name:            "Variant B (Redesign)",
			TargetURL:       "http://localhost:3000/v2",
			AverageSUS:      80.0,
			Grade:           "B",
			AverageFriction: 0.80,
			CreatedAt:       time.Now(),
		},
		Sessions: []store.SessionRecord{
			{
				PersonaID:     "sarah",
				PersonaName:   "Sarah Chen",
				Archetype:     "HIGH_IMPATIENCE_LOW_TECH",
				SUSScore:      80.0,
				FrictionIndex: 0.80,
				GoalCompleted: true,
			},
		},
	}

	cmp := CompareVariants(repA, repB)
	if cmp.DeltaSUS != 28.0 {
		t.Errorf("expected delta SUS 28.0, got %.1f", cmp.DeltaSUS)
	}
	if cmp.DeltaFriction != -1.30 {
		t.Errorf("expected delta friction -1.30, got %.2f", cmp.DeltaFriction)
	}
	if cmp.DeltaDropOffs != -1 {
		t.Errorf("expected drop-off reduction of -1, got %d", cmp.DeltaDropOffs)
	}
	if !strings.Contains(cmp.Winner, "Variant B") {
		t.Errorf("expected Variant B to win, got %s", cmp.Winner)
	}

	htmlContent, err := GenerateDiffHTML(cmp)
	if err != nil {
		t.Fatalf("GenerateDiffHTML error: %v", err)
	}
	if !strings.Contains(htmlContent, "+28.0") {
		t.Errorf("expected +28.0 in HTML diff")
	}

	mdContent, err := GenerateDiffMarkdown(cmp)
	if err != nil {
		t.Fatalf("GenerateDiffMarkdown error: %v", err)
	}
	if !strings.Contains(mdContent, "+28.0 pts") {
		t.Errorf("expected +28.0 pts in Markdown diff")
	}
}
