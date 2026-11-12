package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestGetDefaultPersonas(t *testing.T) {
	personas := GetDefaultPersonas()
	if len(personas) != 5 {
		t.Fatalf("expected 5 personas, got %d", len(personas))
	}

	mobileFound := false
	for _, p := range personas {
		if p.ID == "" {
			t.Errorf("persona missing ID: %+v", p)
		}
		if p.TechnicalLiteracy < 1 || p.TechnicalLiteracy > 10 {
			t.Errorf("persona %s invalid technical literacy: %d", p.ID, p.TechnicalLiteracy)
		}
		if p.Viewport.IsMobile {
			mobileFound = true
		}
	}
	if !mobileFound {
		t.Errorf("expected at least one mobile persona")
	}
}

func TestLoadSuite(t *testing.T) {
	tmpDir := t.TempDir()
	suiteFile := filepath.Join(tmpDir, "test_suite.yaml")

	content := `
version: "1"
suite:
  name: "Sample Test Suite"
  target_url: "https://example.com/signup"
  goal: "Sign up and create an account"
  max_session_duration: "2m"
personas:
  - id: "tester-1"
    name: "Test User"
    age: 30
    archetype: "HIGH_IMPATIENCE_LOW_TECH"
    patience_threshold_seconds: 5.5
    technical_literacy: 4
    attention_span: "LOW"
`
	if err := os.WriteFile(suiteFile, []byte(content), 0644); err != nil {
		t.Fatalf("write temp suite: %v", err)
	}

	suite, err := LoadSuite(suiteFile)
	if err != nil {
		t.Fatalf("failed to load suite: %v", err)
	}

	if suite.Suite.Name != "Sample Test Suite" {
		t.Errorf("expected suite name 'Sample Test Suite', got '%s'", suite.Suite.Name)
	}
	if len(suite.Personas) != 1 {
		t.Fatalf("expected 1 persona, got %d", len(suite.Personas))
	}
	p := suite.Personas[0]
	if p.ID != "tester-1" {
		t.Errorf("expected persona id tester-1, got %s", p.ID)
	}
	if p.Viewport.Width != 1440 || p.Viewport.Height != 900 {
		t.Errorf("expected default desktop viewport 1440x900, got %dx%d", p.Viewport.Width, p.Viewport.Height)
	}
	if dur := suite.ParsedMaxDuration(); dur != 2*time.Minute {
		t.Errorf("expected 2m duration, got %v", dur)
	}
}
