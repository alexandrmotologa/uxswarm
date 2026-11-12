package browser

import (
	"context"
	"fmt"
	"time"

	"github.com/alexandrmotologa/uxswarm/internal/config"
	"github.com/chromedp/chromedp"
)

// Session represents an isolated browser session running a synthetic persona.
type Session struct {
	ID          string
	Persona     config.Persona
	allocCtx    context.Context
	allocCancel context.CancelFunc
	ctx         context.Context
	cancel      context.CancelFunc
	isClosed    bool
}

// SessionOptions configures browser startup parameters.
type SessionOptions struct {
	Headless bool
	ExecPath string
	Timeout  time.Duration
}

const (
	MobileUserAgent  = "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Mobile/15E148 Safari/604.1"
	DesktopUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36"
)

// NewSession creates an isolated CDP browser context tailored to the persona's viewport and archetype.
func NewSession(persona config.Persona, opts SessionOptions) (*Session, error) {
	width := persona.Viewport.Width
	height := persona.Viewport.Height
	if width <= 0 || height <= 0 {
		if persona.Viewport.IsMobile {
			width, height = 390, 844
		} else {
			width, height = 1440, 900
		}
	}

	userAgent := DesktopUserAgent
	if persona.Viewport.IsMobile {
		userAgent = MobileUserAgent
	}

	execOpts := append(
		chromedp.DefaultExecAllocatorOptions[:],
		chromedp.WindowSize(width, height),
		chromedp.UserAgent(userAgent),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("hide-scrollbars", false),
	)

	if opts.ExecPath != "" {
		execOpts = append(execOpts, chromedp.ExecPath(opts.ExecPath))
	}

	if opts.Headless {
		execOpts = append(execOpts, chromedp.Headless)
	} else {
		execOpts = append(execOpts, chromedp.Flag("headless", false))
	}

	allocCtx, allocCancel := chromedp.NewExecAllocator(context.Background(), execOpts...)
	ctx, cancel := chromedp.NewContext(allocCtx)

	// Set viewport and mobile emulation if needed
	if persona.Viewport.IsMobile {
		err := chromedp.Do(ctx,
			chromedp.EmulateViewport(
				int64(width),
				int64(height),
				chromedp.EmulateMobile,
				chromedp.EmulateTouch,
				chromedp.EmulatePortrait,
			),
		)
		if err != nil {
			cancel()
			allocCancel()
			return nil, fmt.Errorf("configure mobile emulation: %w", err)
		}
	}

	return &Session{
		ID:          persona.ID,
		Persona:     persona,
		allocCtx:    allocCtx,
		allocCancel: allocCancel,
		ctx:         ctx,
		cancel:      cancel,
	}, nil
}

// Context returns the active chromedp context for actions.
func (s *Session) Context() context.Context {
	return s.ctx
}

// Navigate opens the target URL and waits for initial page stabilization.
func (s *Session) Navigate(url string) error {
	return chromedp.Do(s.ctx,
		chromedp.Navigate(url),
		chromedp.Sleep(600*time.Millisecond),
	)
}

// CurrentURL returns the active document URL.
func (s *Session) CurrentURL() (string, error) {
	return chromedp.Run(s.ctx, chromedp.Location())
}

// Title returns the current page title.
func (s *Session) Title() (string, error) {
	return chromedp.Run(s.ctx, chromedp.Title())
}

// Close terminates the browser context and cleans up resources.
func (s *Session) Close() {
	if s.isClosed {
		return
	}
	s.isClosed = true
	if s.cancel != nil {
		s.cancel()
	}
	if s.allocCancel != nil {
		s.allocCancel()
	}
}
