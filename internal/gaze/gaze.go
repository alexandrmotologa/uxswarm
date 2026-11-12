package gaze

import (
	"math"
	"strings"

	"github.com/alexandrmotologa/uxswarm/internal/browser"
	"github.com/alexandrmotologa/uxswarm/internal/config"
)

// GazePoint represents a simulated human eye-tracking visual fixation coordinate.
type GazePoint struct {
	X         int     `json:"x"`
	Y         int     `json:"y"`
	Intensity float64 `json:"intensity"` // 0.0 (peripheral) to 1.0 (peak focus)
	Duration  float64 `json:"duration"`  // Fixation duration in seconds
	Target    string  `json:"target,omitempty"`
}

// GenerateGazeMap models visual attention using F-shaped scanning and visual hierarchy saliency.
func GenerateGazeMap(nodes []browser.DOMNode, viewport config.Viewport, goal string) []GazePoint {
	points := make([]GazePoint, 0)
	lowerGoal := strings.ToLower(goal)

	// 1. F-Pattern anchor fixations across top reading band
	topScanX := []int{80, 200, 350, 520}
	for _, x := range topScanX {
		points = append(points, GazePoint{
			X:         x,
			Y:         90,
			Intensity: 0.65,
			Duration:  0.35,
			Target:    "Top Navigation / Headline Scan",
		})
	}

	// 2. F-Pattern vertical left-rail drop
	verticalY := []int{180, 280, 380, 480}
	for _, y := range verticalY {
		points = append(points, GazePoint{
			X:         90,
			Y:         y,
			Intensity: 0.55,
			Duration:  0.28,
			Target:    "Left-Rail Hierarchy Scan",
		})
	}

	// 3. Saliency fixations based on visible interactive nodes
	for _, n := range nodes {
		if !n.Visible {
			continue
		}

		lowerText := strings.ToLower(n.Text + " " + n.AriaLabel + " " + n.Placeholder)
		intensity := 0.4
		duration := 0.25

		// Primary CTA button has highest visual pull
		if n.Tag == "BUTTON" || n.Role == "button" || (n.Tag == "INPUT" && n.Type == "submit") {
			intensity = 0.90
			duration = 0.55
			if strings.Contains(lowerText, "continue") || strings.Contains(lowerText, "sign up") || strings.Contains(lowerText, "submit") {
				intensity = 0.98
				duration = 0.70
			}
		} else if n.Tag == "H1" || n.Tag == "H2" {
			intensity = 0.85
			duration = 0.45
		} else if n.Tag == "INPUT" || n.Tag == "TEXTAREA" {
			intensity = 0.75
			duration = 0.40
		} else if strings.Contains(lowerText, "free") || strings.Contains(lowerText, "trial") || strings.Contains(lowerText, "pro") {
			// Promotional badge trap draws visual fixation
			intensity = 0.80
			duration = 0.45
		}

		// Goal relevance boost
		if lowerGoal != "" && len(n.Text) > 2 && strings.Contains(lowerGoal, strings.ToLower(n.Text)) {
			intensity = math.Min(1.0, intensity+0.2)
		}

		// Mobile reach modifier
		if viewport.IsMobile {
			// Comfortable thumb zone on mobile is Y in [400, 750]
			if n.CenterY >= 400 && n.CenterY <= 750 {
				intensity = math.Min(1.0, intensity+0.15)
			} else if n.CenterY < 200 {
				// Far top requires stretch; visual attention arrives late
				duration += 0.2
			}
		}

		points = append(points, GazePoint{
			X:         n.CenterX,
			Y:         n.CenterY,
			Intensity: math.Round(intensity*100) / 100,
			Duration:  math.Round(duration*100) / 100,
			Target:    n.SummaryText(),
		})
	}

	return points
}
