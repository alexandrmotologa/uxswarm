package persona

import (
	"fmt"
	"math"
	"strings"

	"github.com/alexandrmotologa/uxswarm/internal/browser"
	"github.com/alexandrmotologa/uxswarm/internal/config"
)

// PlannedAction represents the persona's next cognitive decision.
type PlannedAction struct {
	ActionType        string            `json:"action_type"` // click, type, scroll, wait, complete
	TargetNode        *browser.DOMNode  `json:"target_node,omitempty"`
	TypeValue         string            `json:"type_value,omitempty"`
	X                 int               `json:"x"`
	Y                 int               `json:"y"`
	HesitationSeconds float64           `json:"hesitation_seconds"`
	IsDeadClick       bool              `json:"is_dead_click"`
	InnerSpeech       string            `json:"inner_speech"`
	Reasoning         string            `json:"reasoning"`
}

// CognitivePerceiver evaluates visible DOM nodes through the lens of a persona's psychological archetype.
type CognitivePerceiver struct {
	persona config.Persona
	goal    string
}

// NewCognitivePerceiver constructs a new perceiver for a persona and goal.
func NewCognitivePerceiver(persona config.Persona, goal string) *CognitivePerceiver {
	return &CognitivePerceiver{
		persona: persona,
		goal:    goal,
	}
}

// DecideNextAction evaluates the screen state and determines the next interaction.
func (cp *CognitivePerceiver) DecideNextAction(nodes []browser.DOMNode, currentURL string, stepIndex int) PlannedAction {
	lowerURL := strings.ToLower(currentURL)
	lowerGoal := strings.ToLower(cp.goal)

	// 1. Check if goal is already completed
	for _, n := range nodes {
		txt := strings.ToLower(n.Text)
		if strings.Contains(txt, "onboarding complete") ||
			strings.Contains(txt, "welcome to") ||
			strings.Contains(txt, "setup successful") ||
			strings.Contains(txt, "order confirmed") ||
			strings.Contains(txt, "all set!") ||
			(lowerGoal != "" && strings.Contains(txt, "success") && strings.Contains(txt, "complete")) {
			return PlannedAction{
				ActionType:        "complete",
				HesitationSeconds: 0.5,
				InnerSpeech:       SynthesizeInnerSpeech(cp.persona, "complete", &n, 0.5, 0.0, false, false),
				Reasoning:         "Detected confirmation landmark indicating goal completion.",
			}
		}
	}
	if strings.Contains(lowerURL, "/complete") || strings.Contains(lowerURL, "/success") || strings.Contains(lowerURL, "/dashboard") {
		return PlannedAction{
			ActionType:        "complete",
			HesitationSeconds: 0.5,
			InnerSpeech:       SynthesizeInnerSpeech(cp.persona, "complete", nil, 0.5, 0.0, false, false),
			Reasoning:         "URL matches destination success path.",
		}
	}

	// 2. Identify candidate form inputs that need data
	for i := range nodes {
		n := &nodes[i]
		if (n.Tag == "INPUT" || n.Tag == "TEXTAREA") && n.Type != "submit" && n.Type != "button" && n.Type != "hidden" {
			if n.Value == "" && !n.Disabled {
				val := cp.generateInputValue(n)
				hesitation := cp.calculateHesitation(n, false)

				return PlannedAction{
					ActionType:        "type",
					TargetNode:        n,
					TypeValue:         val,
					X:                 n.CenterX,
					Y:                 n.CenterY,
					HesitationSeconds: hesitation,
					IsDeadClick:       false,
					InnerSpeech:       SynthesizeInnerSpeech(cp.persona, "type", n, hesitation, 0.0, false, false),
					Reasoning:         fmt.Sprintf("Filling required field %s with %q.", n.SummaryText(), val),
				}
			}
		}
	}

	// 3. Check for checkboxes (e.g. Terms of Service, consent)
	for i := range nodes {
		n := &nodes[i]
		if n.Type == "checkbox" && !n.Disabled && n.Value != "on" && n.Value != "true" {
			hesitation := cp.calculateHesitation(n, false)
			// Scrutinizer persona reads terms for longer
			if cp.persona.Archetype == config.ArchetypeScrutinizerHighTech {
				hesitation = math.Max(hesitation, 3.8)
			}

			return PlannedAction{
				ActionType:        "click",
				TargetNode:        n,
				X:                 n.CenterX,
				Y:                 n.CenterY,
				HesitationSeconds: hesitation,
				IsDeadClick:       false,
				InnerSpeech:       SynthesizeInnerSpeech(cp.persona, "click", n, hesitation, 0.0, false, false),
				Reasoning:         "Checking required checkbox.",
			}
		}
	}

	// 4. Check for Dead Click traps (e.g. non-interactive banners or text styled like buttons)
	// Impatient or mobile personas occasionally click deceptive elements
	if (cp.persona.Archetype == config.ArchetypeHighImpatienceLowTech || cp.persona.Archetype == config.ArchetypeMobileThumbDriven) && stepIndex == 2 {
		for i := range nodes {
			n := &nodes[i]
			lowerTxt := strings.ToLower(n.Text)
			// Element that looks like a button but is not an interactive button
			if !n.Interactive && (strings.Contains(lowerTxt, "pro plan") || strings.Contains(lowerTxt, "fast setup") || strings.Contains(lowerTxt, "free trial badge")) {
				return PlannedAction{
					ActionType:        "click",
					TargetNode:        n,
					X:                 n.CenterX,
					Y:                 n.CenterY,
					HesitationSeconds: 1.2,
					IsDeadClick:       true,
					InnerSpeech:       SynthesizeInnerSpeech(cp.persona, "click", n, 1.2, 0.4, true, false),
					Reasoning:         "Misclicked non-interactive element that mimicked a button.",
				}
			}
		}
	}

	// 5. Look for primary navigation and submission buttons
	var primaryButton *browser.DOMNode
	var fallbackButton *browser.DOMNode

	for i := range nodes {
		n := &nodes[i]
		if n.Interactive && !n.Disabled {
			lowerText := strings.ToLower(n.Text + " " + n.AriaLabel + " " + n.Value)
			if strings.Contains(lowerText, "continue") ||
				strings.Contains(lowerText, "next") ||
				strings.Contains(lowerText, "sign up") ||
				strings.Contains(lowerText, "create account") ||
				strings.Contains(lowerText, "get started") ||
				strings.Contains(lowerText, "invite") ||
				strings.Contains(lowerText, "submit") ||
				strings.Contains(lowerText, "proceed") ||
				strings.Contains(lowerText, "finish") {
				primaryButton = n
				break
			}
			if n.Tag == "BUTTON" || (n.Tag == "INPUT" && n.Type == "submit") || n.Role == "button" {
				if fallbackButton == nil {
					fallbackButton = n
				}
			}
		}
	}

	targetBtn := primaryButton
	if targetBtn == nil {
		targetBtn = fallbackButton
	}

	if targetBtn != nil {
		hesitation := cp.calculateHesitation(targetBtn, false)
		// Check for small mobile tap target penalty
		if cp.persona.Viewport.IsMobile && (targetBtn.Width < 44 || targetBtn.Height < 44) {
			hesitation += 1.5
		}

		return PlannedAction{
			ActionType:        "click",
			TargetNode:        targetBtn,
			X:                 targetBtn.CenterX,
			Y:                 targetBtn.CenterY,
			HesitationSeconds: hesitation,
			IsDeadClick:       false,
			InnerSpeech:       SynthesizeInnerSpeech(cp.persona, "click", targetBtn, hesitation, 0.0, false, false),
			Reasoning:         fmt.Sprintf("Clicking primary action %s.", targetBtn.SummaryText()),
		}
	}

	// 6. If no button found, scroll down to reveal more content
	return PlannedAction{
		ActionType:        "scroll",
		HesitationSeconds: 1.5,
		InnerSpeech:       "No obvious path forward visible above the fold. Scrolling down to find next step.",
		Reasoning:         "Scrolling viewport.",
	}
}

func (cp *CognitivePerceiver) calculateHesitation(node *browser.DOMNode, isDeadClick bool) float64 {
	base := 1.2

	// Low literacy personas take longer to read and comprehend
	if cp.persona.TechnicalLiteracy <= 3 {
		base += 2.0
	} else if cp.persona.TechnicalLiteracy >= 8 {
		base -= 0.4
	}

	// Mobile thumb reach zone consideration (viewport height ~ 844)
	if cp.persona.Viewport.IsMobile && node != nil {
		// Reach zone: Y between 350 and 700 is comfortable. Y < 200 requires top reach.
		if node.CenterY < 200 {
			base += 1.6 // Reach stretch delay
		}
		// Small tap targets induce hesitation
		if node.Height < 40 || node.Width < 40 {
			base += 1.2
		}
	}

	// Scrutinizer archetype pauses to verify security and privacy text
	if cp.persona.Archetype == config.ArchetypeScrutinizerHighTech {
		base += 1.8
	}

	// High impatience archetype acts faster on simple primary targets
	if cp.persona.Archetype == config.ArchetypeHighImpatienceLowTech {
		base *= 0.8
	}

	return math.Round(base*10) / 10
}

func (cp *CognitivePerceiver) generateInputValue(node *browser.DOMNode) string {
	descriptor := strings.ToLower(node.Name + " " + node.Placeholder + " " + node.AriaLabel + " " + node.Type)

	switch {
	case strings.Contains(descriptor, "password"):
		if cp.persona.Archetype == config.ArchetypeScrutinizerHighTech {
			return "Secure#P@ssw0rd2026!"
		}
		return "SecretPassword123!"
	case strings.Contains(descriptor, "email") || strings.Contains(descriptor, "invite"):
		if strings.Contains(descriptor, "colleague") || strings.Contains(descriptor, "team") || strings.Contains(descriptor, "invite") {
			return "colleague@example.com"
		}
		cleanID := strings.ReplaceAll(cp.persona.ID, "-", "")
		return fmt.Sprintf("%s@example.com", cleanID)
	case strings.Contains(descriptor, "workspace") || strings.Contains(descriptor, "company") || strings.Contains(descriptor, "tenant") || strings.Contains(descriptor, "org"):
		return "Acme Labs"
	case strings.Contains(descriptor, "name"):
		return cp.persona.Name
	default:
		return "Test Value"
	}
}
