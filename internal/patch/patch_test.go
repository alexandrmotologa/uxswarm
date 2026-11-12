package patch

import (
	"strings"
	"testing"

	"github.com/alexandrmotologa/uxswarm/internal/metrics"
)

func TestGeneratePatches(t *testing.T) {
	findings := []metrics.HeuristicFinding{
		{
			HeuristicID:      1,
			HeuristicName:    "Visibility of system status",
			Severity:         metrics.SeverityCritical,
			Title:            "Hesitation latency high",
			Evidence:         "Users paused 4s",
			AffectedPersonas: []string{"Sarah Chen"},
		},
		{
			HeuristicID:      7,
			HeuristicName:    "Flexibility and efficiency of use",
			Severity:         metrics.SeverityMajor,
			Title:            "Mobile tap target below 44px",
			Evidence:         "Elena struggled with button",
			AffectedPersonas: []string{"Elena Rostova"},
		},
	}

	patches := GeneratePatches(findings, nil)
	if len(patches) != 2 {
		t.Fatalf("expected 2 patches, got %d", len(patches))
	}

	if patches[0].HeuristicID != 1 {
		t.Errorf("expected Heuristic 1, got %d", patches[0].HeuristicID)
	}
	if !strings.Contains(patches[0].AfterCode, "disabled={isSubmitting}") {
		t.Errorf("expected disabled state in Heuristic 1 patch")
	}

	if patches[1].HeuristicID != 7 {
		t.Errorf("expected Heuristic 7, got %d", patches[1].HeuristicID)
	}
	if !strings.Contains(patches[1].AfterCode, "min-height: 44px") {
		t.Errorf("expected 44px touch target in Heuristic 7 patch")
	}

	formatted := FormatPatchText(patches[0])
	if !strings.Contains(formatted, "BEFORE") || !strings.Contains(formatted, "AFTER") {
		t.Errorf("formatted text missing headers: %s", formatted)
	}
}
