package persona

import (
	"fmt"
	"strings"

	"github.com/alexandrmotologa/uxswarm/internal/browser"
	"github.com/alexandrmotologa/uxswarm/internal/config"
)

// SynthesizeInnerSpeech generates authentic first-person thoughts reflecting the persona's reaction.
func SynthesizeInnerSpeech(persona config.Persona, actionType string, target *browser.DOMNode, hesitation float64, frustration float64, isDeadClick bool, isDropOff bool) string {
	if isDropOff {
		switch persona.Archetype {
		case config.ArchetypeHighImpatienceLowTech:
			return "I've had enough of this confusing interface. This should take 30 seconds, not minutes. Abandoning."
		case config.ArchetypeScrutinizerHighTech:
			return "The workflow lacks standard security feedback and basic usability consistency. Abandoning evaluation."
		case config.ArchetypeMobileThumbDriven:
			return "Impossible to use comfortably on a phone. Tap targets are too small and the layout is frustrating."
		case config.ArchetypeAccessibilityReliant:
			return "Form feedback is unclear and focus states are uncommunicative. Cannot reliably proceed."
		default:
			return "Too much friction and confusion in this step. Leaving page."
		}
	}

	if isDeadClick {
		targetLabel := "this element"
		if target != nil && target.Text != "" {
			targetLabel = fmt.Sprintf("'%s'", target.Text)
		}
		switch persona.Archetype {
		case config.ArchetypeHighImpatienceLowTech:
			return fmt.Sprintf("I clicked %s expecting something to happen, but nothing did! Why is it styled like a button?", targetLabel)
		case config.ArchetypeMobileThumbDriven:
			return fmt.Sprintf("My thumb tapped %s, but it's completely unresponsive. Maybe the tap area is too small.", targetLabel)
		default:
			return fmt.Sprintf("Clicked on %s, but received zero visual feedback or navigation change.", targetLabel)
		}
	}

	if (persona.PatienceThresholdSeconds > 0 && hesitation >= persona.PatienceThresholdSeconds) || hesitation >= 3.5 {
		switch persona.Archetype {
		case config.ArchetypeHighImpatienceLowTech:
			return fmt.Sprintf("Stalled for %.1fs scanning the page. The primary button is buried under too much text.", hesitation)
		case config.ArchetypeScrutinizerHighTech:
			return fmt.Sprintf("Reviewing the page details and input labels carefully (%.1fs delay) before taking an action.", hesitation)
		case config.ArchetypeMobileThumbDriven:
			return fmt.Sprintf("Reorienting my grip to reach the top half of the screen (%.1fs pause).", hesitation)
		default:
			return fmt.Sprintf("Hesitated for %.1fs trying to identify what the interface expects me to do next.", hesitation)
		}
	}

	switch actionType {
	case "navigate":
		return fmt.Sprintf("Page loaded. As %s, scanning for the fastest way to achieve my goal.", persona.Name)
	case "type":
		if target != nil {
			label := target.Placeholder
			if label == "" {
				label = target.Name
			}
			return fmt.Sprintf("Filling in %q. Hopefully this doesn't have unannounced format restrictions.", label)
		}
		return "Entering required information."
	case "click":
		if target != nil {
			text := target.Text
			if text == "" {
				text = target.AriaLabel
			}
			if strings.EqualFold(text, "back") || strings.EqualFold(text, "cancel") {
				return "Taking a step back to re-evaluate the previous screen."
			}
			return fmt.Sprintf("Clicking %q to proceed forward.", text)
		}
		return "Clicking selected action control."
	case "complete":
		return "Goal successfully accomplished! The remaining flow was straightforward."
	default:
		return "Observing screen layout and next available action."
	}
}
