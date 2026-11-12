package persona

import (
	"context"
	"fmt"
	"time"

	"github.com/alexandrmotologa/uxswarm/internal/browser"
	"github.com/alexandrmotologa/uxswarm/internal/config"
	"github.com/alexandrmotologa/uxswarm/internal/metrics"
)

// Agent runs an autonomous persona navigation session.
type Agent struct {
	Persona     config.Persona
	Goal        string
	Session     *browser.Session
	Actuator    *browser.Actuator
	Recorder    *browser.Recorder
	Telemetry   *metrics.TelemetryTracker
	Frustration *FrustrationTracker
	Perceiver   *CognitivePerceiver
	OnStep      func(event metrics.StepEvent)
}

// NewAgent constructs a ready-to-run cognitive user agent.
func NewAgent(persona config.Persona, goal string, session *browser.Session, storageDir string) *Agent {
	return &Agent{
		Persona:     persona,
		Goal:        goal,
		Session:     session,
		Actuator:    browser.NewActuator(session),
		Recorder:    browser.NewRecorder(session, storageDir),
		Telemetry:   metrics.NewTelemetryTracker(),
		Frustration: NewFrustrationTracker(persona),
		Perceiver:   NewCognitivePerceiver(persona, goal),
	}
}

// Execute drives the browser through the goal flow, evaluating cognitive load at each step.
func (a *Agent) Execute(ctx context.Context, startURL string, maxSteps int) (*metrics.TelemetrySummary, []metrics.StepEvent, error) {
	if maxSteps <= 0 {
		maxSteps = 25
	}

	// 1. Initial navigation
	navStart := time.Now()
	if err := a.Session.Navigate(startURL); err != nil {
		return nil, nil, fmt.Errorf("initial navigation to %s: %w", startURL, err)
	}

	currentURL, _ := a.Session.CurrentURL()
	initialFrame, _ := a.Recorder.CaptureFrame(0)
	frameURI := ""
	if initialFrame != nil {
		frameURI = initialFrame.DataURI
	}

	initialEvt := metrics.StepEvent{
		Index:             0,
		Timestamp:         navStart,
		ActionType:        "navigate",
		URL:               currentURL,
		TargetDescription: "Initial Page Load",
		InnerSpeech:       SynthesizeInnerSpeech(a.Persona, "navigate", nil, 0.5, 0.0, false, false),
		Frustration:       0.0,
		FrameDataURI:      frameURI,
	}
	a.Telemetry.RecordStep(initialEvt)
	if a.OnStep != nil {
		a.OnStep(initialEvt)
	}

	// 2. Persona Action Loop
	for step := 1; step <= maxSteps; step++ {
		select {
		case <-ctx.Done():
			return a.Telemetry.Summary(), a.Telemetry.Steps(), ctx.Err()
		default:
		}

		currentURL, _ = a.Session.CurrentURL()
		nodes, err := browser.ExtractDOMNodes(a.Session.Context())
		if err != nil {
			// If extraction fails briefly, sleep and retry once
			time.Sleep(500 * time.Millisecond)
			nodes, err = browser.ExtractDOMNodes(a.Session.Context())
			if err != nil {
				return nil, nil, fmt.Errorf("extract dom nodes at step %d: %w", step, err)
			}
		}

		// Cognitive evaluation: what does the persona decide to do?
		decision := a.Perceiver.DecideNextAction(nodes, currentURL, step)

		// Check if goal completed
		if decision.ActionType == "complete" {
			a.Telemetry.MarkCompleted()
			frame, _ := a.Recorder.CaptureFrame(step)
			evt := metrics.StepEvent{
				Index:             step,
				Timestamp:         time.Now(),
				ActionType:        "complete",
				URL:               currentURL,
				TargetDescription: "Goal Accomplished",
				InnerSpeech:       decision.InnerSpeech,
				HesitationSeconds: decision.HesitationSeconds,
				Frustration:       a.Frustration.Frustration,
			}
			if frame != nil {
				evt.FrameDataURI = frame.DataURI
				evt.ScreenshotPath = frame.FilePath
			}
			a.Telemetry.RecordStep(evt)
			if a.OnStep != nil {
				a.OnStep(evt)
			}
			break
		}

		// Check frustration accumulation & possible abandonment
		isBacktrack := false
		urls := a.Telemetry.Summary().URLSequence
		if len(urls) >= 2 && currentURL == urls[len(urls)-2] {
			isBacktrack = true
		}

		frust, droppedOff := a.Frustration.Update(decision.HesitationSeconds, decision.IsDeadClick, isBacktrack)

		if droppedOff {
			a.Telemetry.MarkDroppedOff(a.Frustration.Reason)
			frame, _ := a.Recorder.CaptureFrame(step)
			evt := metrics.StepEvent{
				Index:             step,
				Timestamp:         time.Now(),
				ActionType:        "drop_off",
				URL:               currentURL,
				TargetDescription: "Session Abandoned",
				InnerSpeech:       SynthesizeInnerSpeech(a.Persona, "drop_off", decision.TargetNode, decision.HesitationSeconds, frust, decision.IsDeadClick, true),
				HesitationSeconds: decision.HesitationSeconds,
				Frustration:       frust,
			}
			if frame != nil {
				evt.FrameDataURI = frame.DataURI
				evt.ScreenshotPath = frame.FilePath
			}
			a.Telemetry.RecordStep(evt)
			if a.OnStep != nil {
				a.OnStep(evt)
			}
			break
		}

		// Execute physical action via actuator
		var actErr error
		switch decision.ActionType {
		case "type":
			if decision.TargetNode != nil && decision.TargetNode.Selector != "" {
				actErr = a.Actuator.Type(decision.TargetNode.Selector, decision.TypeValue)
			}
		case "click":
			if decision.TargetNode != nil {
				if decision.TargetNode.Selector != "" && !decision.IsDeadClick {
					actErr = a.Actuator.ClickSelector(decision.TargetNode.Selector)
				} else {
					actErr = a.Actuator.ClickCoordinate(decision.X, decision.Y)
				}
			}
		case "scroll":
			actErr = a.Actuator.Scroll(300)
		case "wait":
			actErr = a.Actuator.Wait(1 * time.Second)
		}

		if actErr != nil {
			// If action execution errored, record hesitation/error and proceed
			decision.HesitationSeconds += 1.5
		}

		// Capture frame for replay
		frame, _ := a.Recorder.CaptureFrame(step)
		targetDesc := "Page"
		targetRole := ""
		targetSel := ""
		if decision.TargetNode != nil {
			targetDesc = decision.TargetNode.SummaryText()
			targetRole = decision.TargetNode.Role
			targetSel = decision.TargetNode.Selector
		}

		evt := metrics.StepEvent{
			Index:             step,
			Timestamp:         time.Now(),
			ActionType:        decision.ActionType,
			TargetSelector:    targetSel,
			TargetDescription: targetDesc,
			TargetRole:        targetRole,
			X:                 decision.X,
			Y:                 decision.Y,
			IsDeadClick:       decision.IsDeadClick,
			HesitationSeconds: decision.HesitationSeconds,
			URL:               currentURL,
			InnerSpeech:       decision.InnerSpeech,
			Frustration:       frust,
		}
		if frame != nil {
			evt.FrameDataURI = frame.DataURI
			evt.ScreenshotPath = frame.FilePath
		}

		a.Telemetry.RecordStep(evt)
		if a.OnStep != nil {
			a.OnStep(evt)
		}

		// Small stabilization wait after action
		time.Sleep(300 * time.Millisecond)
	}

	return a.Telemetry.Summary(), a.Telemetry.Steps(), nil
}
