package browser

import (
	"strings"
	"testing"
)

func TestDOMNode_SummaryText(t *testing.T) {
	node := DOMNode{
		Tag:      "BUTTON",
		Role:     "button",
		Text:     "Complete Registration",
		CenterX:  250,
		CenterY:  400,
		Disabled: false,
	}

	summary := node.SummaryText()
	if !strings.Contains(summary, "[BUTTON]") {
		t.Errorf("expected [BUTTON] in summary, got %q", summary)
	}
	if !strings.Contains(summary, "Complete Registration") {
		t.Errorf("expected text in summary, got %q", summary)
	}
	if !strings.Contains(summary, "at (250, 400)") {
		t.Errorf("expected coordinates in summary, got %q", summary)
	}

	disabledNode := DOMNode{
		Tag:      "INPUT",
		Type:     "submit",
		Role:     "button",
		Value:    "Submit",
		Disabled: true,
	}
	if !strings.Contains(disabledNode.SummaryText(), "(disabled)") {
		t.Errorf("expected (disabled) in summary, got %q", disabledNode.SummaryText())
	}
}
