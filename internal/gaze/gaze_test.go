package gaze

import (
	"testing"

	"github.com/alexandrmotologa/uxswarm/internal/browser"
	"github.com/alexandrmotologa/uxswarm/internal/config"
)

func TestGenerateGazeMap(t *testing.T) {
	viewport := config.Viewport{Width: 1440, Height: 900, IsMobile: false}
	nodes := []browser.DOMNode{
		{
			Tag:         "BUTTON",
			Role:        "button",
			Text:        "Continue to Workspace",
			CenterX:     500,
			CenterY:     600,
			Interactive: true,
			Visible:     true,
		},
		{
			Tag:     "H1",
			Text:    "Create Account",
			CenterX: 300,
			CenterY: 150,
			Visible: true,
		},
	}

	points := GenerateGazeMap(nodes, viewport, "Create account and continue")
	if len(points) == 0 {
		t.Fatalf("expected gaze points, got 0")
	}

	foundButtonFixation := false
	for _, p := range points {
		if p.X == 500 && p.Y == 600 {
			foundButtonFixation = true
			if p.Intensity < 0.9 {
				t.Errorf("expected high intensity for primary button, got %.2f", p.Intensity)
			}
		}
	}
	if !foundButtonFixation {
		t.Errorf("expected gaze fixation point on button at (500, 600)")
	}
}
