package persona

import (
	"strings"
	"testing"

	"github.com/alexandrmotologa/uxswarm/internal/browser"
	"github.com/alexandrmotologa/uxswarm/internal/config"
)

func TestPerceiver_IdentifiesFormInputs(t *testing.T) {
	persona := config.Persona{
		ID:                "sarah-cfo",
		Name:              "Sarah Chen",
		Archetype:         config.ArchetypeHighImpatienceLowTech,
		TechnicalLiteracy: 3,
	}

	perceiver := NewCognitivePerceiver(persona, "Sign up for account")

	nodes := []browser.DOMNode{
		{
			Tag:         "INPUT",
			Type:        "email",
			Name:        "email",
			Placeholder: "Work Email",
			Selector:    "#email",
			Value:       "",
			CenterX:     150,
			CenterY:     120,
			Interactive: true,
			Visible:     true,
		},
		{
			Tag:         "BUTTON",
			Role:        "button",
			Text:        "Continue",
			Selector:    "#submit-btn",
			CenterX:     150,
			CenterY:     250,
			Interactive: true,
			Visible:     true,
		},
	}

	action := perceiver.DecideNextAction(nodes, "http://localhost:3000/signup", 1)
	if action.ActionType != "type" {
		t.Errorf("expected action type 'type', got %s", action.ActionType)
	}
	if !strings.Contains(action.TypeValue, "@example.com") {
		t.Errorf("expected email value, got %s", action.TypeValue)
	}
	if action.HesitationSeconds <= 0 {
		t.Errorf("expected positive hesitation seconds, got %.1f", action.HesitationSeconds)
	}
}

func TestPerceiver_DetectsGoalCompletion(t *testing.T) {
	persona := config.Persona{
		ID:        "test-user",
		Name:      "Test User",
		Archetype: config.ArchetypeScrutinizerHighTech,
	}

	perceiver := NewCognitivePerceiver(persona, "Complete onboarding")

	nodes := []browser.DOMNode{
		{
			Tag:     "H1",
			Text:    "Setup Successful! Welcome to Acme",
			Visible: true,
		},
	}

	action := perceiver.DecideNextAction(nodes, "http://localhost:3000/dashboard", 5)
	if action.ActionType != "complete" {
		t.Errorf("expected action type 'complete', got %s", action.ActionType)
	}
}

func TestFrustrationTracker_AccumulatesAndTriggersDropOff(t *testing.T) {
	persona := config.Persona{
		ID:                       "impatient-exec",
		Name:                     "Sarah Chen",
		Archetype:                config.ArchetypeHighImpatienceLowTech,
		PatienceThresholdSeconds: 3.0,
	}

	tracker := NewFrustrationTracker(persona)

	// Update with severe hesitation and dead click
	frust, dropped := tracker.Update(4.5, true, false)
	if frust <= 0 {
		t.Errorf("expected positive frustration, got %.2f", frust)
	}

	// Repeated dead click should push frustration >= 1.0
	frust, dropped = tracker.Update(4.0, true, false)
	if !dropped {
		t.Errorf("expected persona to drop off, got frustration %.2f", frust)
	}
	if tracker.Reason == "" {
		t.Errorf("expected non-empty drop-off reason")
	}
}
