package browser

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/chromedp/chromedp"
)

// DOMNode represents a visible or interactive element extracted from the browser DOM.
type DOMNode struct {
	Index       int    `json:"index"`
	Tag         string `json:"tag"`
	Role        string `json:"role"`
	Type        string `json:"type"`
	Text        string `json:"text"`
	Placeholder string `json:"placeholder"`
	Name        string `json:"name"`
	Value       string `json:"value"`
	AriaLabel   string `json:"aria_label"`
	Selector    string `json:"selector"`
	X           int    `json:"x"`
	Y           int    `json:"y"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	CenterX     int    `json:"center_x"`
	CenterY     int    `json:"center_y"`
	Interactive bool   `json:"interactive"`
	Visible     bool   `json:"visible"`
	Disabled    bool   `json:"disabled"`
}

// SummaryText returns a concise human-readable description for persona perception.
func (n DOMNode) SummaryText() string {
	parts := make([]string, 0)
	if n.Role != "" {
		parts = append(parts, fmt.Sprintf("[%s]", strings.ToUpper(n.Role)))
	} else {
		parts = append(parts, fmt.Sprintf("<%s>", strings.ToLower(n.Tag)))
	}

	label := n.Text
	if label == "" {
		label = n.AriaLabel
	}
	if label == "" {
		label = n.Placeholder
	}
	if label == "" {
		label = n.Name
	}
	if label != "" {
		if len(label) > 60 {
			label = label[:57] + "..."
		}
		parts = append(parts, fmt.Sprintf("%q", label))
	}

	if n.Disabled {
		parts = append(parts, "(disabled)")
	}

	parts = append(parts, fmt.Sprintf("at (%d, %d)", n.CenterX, n.CenterY))
	return strings.Join(parts, " ")
}

// ExtractionScript is the JavaScript payload executed in the target page context.
const ExtractionScript = `
(function() {
  const elements = [];
  const candidates = document.querySelectorAll(
    'button, a, input, select, textarea, [role="button"], [role="link"], [role="checkbox"], [tabindex], h1, h2, h3, [data-interactive="true"], [data-ux-target]'
  );

  function getUniqueSelector(el) {
    if (el.id) return '#' + CSS.escape(el.id);
    if (el.getAttribute('data-testid')) return '[data-testid="' + el.getAttribute('data-testid') + '"]';
    if (el.getAttribute('name')) return el.tagName.toLowerCase() + '[name="' + el.getAttribute('name') + '"]';
    
    let path = [];
    let current = el;
    while (current && current.nodeType === Node.ELEMENT_NODE && current !== document.body) {
      let tag = current.tagName.toLowerCase();
      let siblings = Array.from(current.parentNode ? current.parentNode.children : []).filter(e => e.tagName === current.tagName);
      if (siblings.length > 1) {
        let index = siblings.indexOf(current) + 1;
        path.unshift(tag + ':nth-of-type(' + index + ')');
      } else {
        path.unshift(tag);
      }
      current = current.parentNode;
    }
    return path.join(' > ');
  }

  let idx = 0;
  for (const el of candidates) {
    const rect = el.getBoundingClientRect();
    const style = window.getComputedStyle(el);
    
    const visible = rect.width > 0 && rect.height > 0 && 
                    style.visibility !== 'hidden' && 
                    style.display !== 'none' && 
                    style.opacity !== '0';

    if (!visible) continue;

    const tag = el.tagName.toUpperCase();
    const isInteractive = ['BUTTON', 'A', 'INPUT', 'SELECT', 'TEXTAREA'].includes(tag) ||
                          el.getAttribute('role') === 'button' ||
                          el.getAttribute('role') === 'link' ||
                          el.getAttribute('role') === 'checkbox' ||
                          el.hasAttribute('onclick') ||
                          el.getAttribute('data-interactive') === 'true';

    let text = (el.innerText || el.textContent || '').trim();
    if (text.length > 150) text = text.substring(0, 150);

    idx++;
    elements.push({
      index: idx,
      tag: tag,
      role: el.getAttribute('role') || (tag === 'BUTTON' ? 'button' : tag === 'A' ? 'link' : tag === 'INPUT' ? 'textbox' : ''),
      type: el.getAttribute('type') || '',
      text: text,
      placeholder: el.getAttribute('placeholder') || '',
      name: el.getAttribute('name') || '',
      value: el.value || '',
      aria_label: el.getAttribute('aria-label') || el.getAttribute('aria-labelledby') || '',
      selector: getUniqueSelector(el),
      x: Math.round(rect.left),
      y: Math.round(rect.top),
      width: Math.round(rect.width),
      height: Math.round(rect.height),
      center_x: Math.round(rect.left + rect.width / 2),
      center_y: Math.round(rect.top + rect.height / 2),
      interactive: isInteractive,
      visible: true,
      disabled: el.disabled || el.getAttribute('aria-disabled') === 'true'
    });
  }

  return JSON.stringify(elements);
})();
`

// ExtractDOMNodes evaluates the ExtractionScript in the current page and parses the node tree.
func ExtractDOMNodes(ctx context.Context) ([]DOMNode, error) {
	jsonStr, err := chromedp.Run(ctx, chromedp.Evaluate[string](ExtractionScript))
	if err != nil {
		return nil, fmt.Errorf("evaluate extraction script: %w", err)
	}

	var nodes []DOMNode
	if err := json.Unmarshal([]byte(jsonStr), &nodes); err != nil {
		return nil, fmt.Errorf("unmarshal extracted nodes: %w", err)
	}

	return nodes, nil
}
