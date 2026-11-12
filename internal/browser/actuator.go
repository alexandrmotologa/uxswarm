package browser

import (
	"fmt"
	"time"

	"github.com/chromedp/chromedp"
)

// Actuator simulates physical human user actions via Chrome DevTools Protocol.
type Actuator struct {
	session *Session
}

// NewActuator constructs an actuator bound to a browser session.
func NewActuator(session *Session) *Actuator {
	return &Actuator{session: session}
}

// ClickCoordinate moves the mouse cursor to (x, y) and fires a physical click event.
func (a *Actuator) ClickCoordinate(x, y int) error {
	ctx := a.session.Context()
	xf := float64(x)
	yf := float64(y)

	return chromedp.Do(ctx,
		chromedp.MouseClickXY(xf, yf),
		chromedp.Sleep(200*time.Millisecond),
	)
}

// ClickSelector locates the element by CSS selector and clicks it.
func (a *Actuator) ClickSelector(selector string) error {
	ctx := a.session.Context()
	return chromedp.Do(ctx,
		chromedp.ScrollIntoView(selector),
		chromedp.Click(selector),
		chromedp.Sleep(300*time.Millisecond),
	)
}

// Type focuses the input field, clears existing content, and sends keystrokes with human-like delays.
func (a *Actuator) Type(selector string, text string) error {
	ctx := a.session.Context()
	clearScript := fmt.Sprintf(`
		(function() {
			const el = document.querySelector(%q);
			if (el) {
				el.value = '';
				el.dispatchEvent(new Event('input', { bubbles: true }));
				el.dispatchEvent(new Event('change', { bubbles: true }));
			}
		})()
	`, selector)

	return chromedp.Do(ctx,
		chromedp.ScrollIntoView(selector),
		chromedp.Focus(selector),
		chromedp.Evaluate[chromedp.Void](clearScript),
		chromedp.SendKeys(selector, text),
		chromedp.Sleep(150*time.Millisecond),
	)
}

// Scroll scrolls the document by deltaY pixels.
func (a *Actuator) Scroll(deltaY int) error {
	ctx := a.session.Context()
	script := fmt.Sprintf("window.scrollBy({ top: %d, behavior: 'smooth' });", deltaY)
	return chromedp.Do(ctx,
		chromedp.Evaluate[chromedp.Void](script),
		chromedp.Sleep(300*time.Millisecond),
	)
}

// Wait pauses execution for the specified duration.
func (a *Actuator) Wait(duration time.Duration) error {
	ctx := a.session.Context()
	return chromedp.Do(ctx, chromedp.Sleep(duration))
}
