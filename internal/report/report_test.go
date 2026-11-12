package report

import (
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/alexandrmotologa/uxswarm/internal/metrics"
	"github.com/alexandrmotologa/uxswarm/internal/store"
)

func TestReportGeneration(t *testing.T) {
	data := &ReportData{
		Run: store.RunRecord{
			ID:                "run_1",
			Name:              "Onboarding Test",
			TargetURL:         "http://localhost:3000/signup",
			Goal:              "Sign up and create workspace",
			Status:            "completed",
			TotalPersonas:     1,
			AverageSUS:        82.5,
			AverageFriction:   1.2,
			Grade:             "A",
			MandatoryRedesign: false,
			CreatedAt:         time.Now(),
		},
		Sessions: []store.SessionRecord{
			{
				PersonaName:   "Sarah Chen",
				Archetype:     "HIGH_IMPATIENCE_LOW_TECH",
				SUSScore:      82.5,
				SUSGrade:      "A",
				FrictionIndex: 1.2,
				GoalCompleted: true,
			},
		},
		Heuristics: []metrics.HeuristicFinding{
			{
				HeuristicID:      1,
				HeuristicName:    "Visibility of system status",
				Severity:         metrics.SeverityMinor,
				Title:            "Slight hesitation on submit",
				Evidence:         "Sarah paused 3.6s",
				AffectedPersonas: []string{"Sarah Chen"},
				Recommendation:   "Add loading spinner",
			},
		},
	}

	htmlContent, err := GenerateHTML(data)
	if err != nil {
		t.Fatalf("GenerateHTML error: %v", err)
	}
	if !strings.Contains(htmlContent, "82.5") || !strings.Contains(htmlContent, "Sarah Chen") {
		t.Errorf("HTML missing expected data: %s", htmlContent)
	}

	mdContent, err := GenerateMarkdown(data)
	if err != nil {
		t.Fatalf("GenerateMarkdown error: %v", err)
	}
	if !strings.Contains(mdContent, "82.5") || !strings.Contains(mdContent, "Sarah Chen") {
		t.Errorf("Markdown missing expected data: %s", mdContent)
	}

	tmpDir := t.TempDir()
	htmlFile := filepath.Join(tmpDir, "report.html")
	if err := ExportToFile(data, "html", htmlFile); err != nil {
		t.Fatalf("ExportToFile html error: %v", err)
	}

	jsonFile := filepath.Join(tmpDir, "report.json")
	if err := ExportToFile(data, "json", jsonFile); err != nil {
		t.Fatalf("ExportToFile json error: %v", err)
	}
}
