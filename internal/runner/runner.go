package runner

import (
	"context"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"time"

	"github.com/alexandrmotologa/uxswarm/internal/browser"
	"github.com/alexandrmotologa/uxswarm/internal/config"
	"github.com/alexandrmotologa/uxswarm/internal/metrics"
	"github.com/alexandrmotologa/uxswarm/internal/persona"
	"github.com/alexandrmotologa/uxswarm/internal/report"
	"github.com/alexandrmotologa/uxswarm/internal/store"
)

// RunOptions configures the execution parameters of a usability study.
type RunOptions struct {
	Name        string
	TargetURL   string
	Goal        string
	Personas    []config.Persona
	Headless    bool
	MaxSteps    int
	StorageDir  string
	Store       *store.Store
	OnStep      func(persona config.Persona, event metrics.StepEvent)
	OnSessionEnd func(persona config.Persona, summary *metrics.TelemetrySummary, sus *metrics.SUSResult)
}

// ExecuteSwarm drives the complete persona swarm through the target application flow.
func ExecuteSwarm(ctx context.Context, opts RunOptions) (*report.ReportData, error) {
	if len(opts.Personas) == 0 {
		opts.Personas = config.GetDefaultPersonas()
	}
	if opts.MaxSteps <= 0 {
		opts.MaxSteps = 20
	}
	if opts.Name == "" {
		opts.Name = fmt.Sprintf("Study: %s", opts.TargetURL)
	}

	runID := fmt.Sprintf("run_%d", time.Now().UnixNano()/1e6)
	startTime := time.Now()

	runRecord := store.RunRecord{
		ID:            runID,
		Name:          opts.Name,
		TargetURL:     opts.TargetURL,
		Goal:          opts.Goal,
		Status:        "running",
		TotalPersonas: len(opts.Personas),
		CreatedAt:     startTime,
	}

	if opts.Store != nil {
		if err := opts.Store.CreateRun(&runRecord); err != nil {
			return nil, fmt.Errorf("create run in store: %w", err)
		}
	}

	summaries := make([]*metrics.TelemetrySummary, 0, len(opts.Personas))
	susResults := make([]*metrics.SUSResult, 0, len(opts.Personas))
	sessionRecords := make([]store.SessionRecord, 0, len(opts.Personas))
	stepsMap := make(map[string][]metrics.StepEvent)
	storeStepsMap := make(map[string][]store.StepRecord)

	totalSUS := 0.0
	totalFriction := 0.0

	for _, p := range opts.Personas {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		sessDir := ""
		if opts.StorageDir != "" {
			sessDir = filepath.Join(opts.StorageDir, runID, p.ID)
			_ = os.MkdirAll(sessDir, 0755)
		}

		// Launch isolated headless browser session
		sess, err := browser.NewSession(p, browser.SessionOptions{
			Headless: opts.Headless,
		})
		if err != nil {
			return nil, fmt.Errorf("launch browser session for %s: %w", p.Name, err)
		}

		agent := persona.NewAgent(p, opts.Goal, sess, sessDir)
		agent.OnStep = func(evt metrics.StepEvent) {
			if opts.OnStep != nil {
				opts.OnStep(p, evt)
			}
		}

		// Run navigation loop
		summary, steps, runErr := agent.Execute(ctx, opts.TargetURL, opts.MaxSteps)
		sess.Close()

		if runErr != nil && ctx.Err() != nil {
			return nil, ctx.Err()
		}

		// Synthesize SUS responses based on recorded session telemetry
		susRes := metrics.SynthesizeSUS(p, summary)

		summaries = append(summaries, summary)
		susResults = append(susResults, susRes)
		stepsMap[p.ID] = steps

		totalSUS += susRes.Score
		totalFriction += summary.FrictionIndex

		sessID := fmt.Sprintf("%s_%s", runID, p.ID)
		sessRec := store.SessionRecord{
			ID:              sessID,
			RunID:           runID,
			PersonaID:       p.ID,
			PersonaName:     p.Name,
			Archetype:       p.Archetype,
			ViewportWidth:   p.Viewport.Width,
			ViewportHeight:  p.Viewport.Height,
			IsMobile:        p.Viewport.IsMobile,
			SUSScore:        susRes.Score,
			SUSGrade:        string(susRes.Grade),
			FrictionIndex:   summary.FrictionIndex,
			TotalSteps:      summary.TotalSteps,
			TotalDuration:   summary.TotalDurationSeconds,
			HesitationCount: summary.HesitationCount,
			DeadClicks:      summary.DeadClicks,
			Backtracks:      summary.Backtracks,
			GoalCompleted:   summary.GoalCompleted,
			DroppedOff:      summary.DroppedOff,
			DropOffReason:   summary.DropOffReason,
			CreatedAt:       time.Now(),
		}
		sessionRecords = append(sessionRecords, sessRec)

		if opts.Store != nil {
			if err := opts.Store.CreateSession(&sessRec); err != nil {
				fmt.Printf("Warning: failed to create session %s: %v\n", sessID, err)
			}
		}

		// Convert steps to store records
		storeSteps := make([]store.StepRecord, 0, len(steps))
		for _, st := range steps {
			sr := store.StepRecord{
				SessionID:         sessID,
				StepIndex:         st.Index,
				Timestamp:         st.Timestamp,
				ActionType:        st.ActionType,
				TargetSelector:    st.TargetSelector,
				TargetDescription: st.TargetDescription,
				TargetRole:        st.TargetRole,
				X:                 st.X,
				Y:                 st.Y,
				IsDeadClick:       st.IsDeadClick,
				HesitationSeconds: st.HesitationSeconds,
				URL:               st.URL,
				InnerSpeech:       st.InnerSpeech,
				Frustration:       st.Frustration,
				ScreenshotPath:    st.ScreenshotPath,
				FrameDataURI:      st.FrameDataURI,
			}
			storeSteps = append(storeSteps, sr)
			if opts.Store != nil {
				if err := opts.Store.CreateStep(&sr); err != nil {
					fmt.Printf("Warning: failed to create step %d for session %s: %v\n", st.Index, sessID, err)
				}
			}
		}
		storeStepsMap[sessID] = storeSteps

		if opts.OnSessionEnd != nil {
			opts.OnSessionEnd(p, summary, susRes)
		}
	}

	// Compute aggregate scores
	avgSUS := 0.0
	avgFriction := 0.0
	if len(opts.Personas) > 0 {
		avgSUS = math.Round((totalSUS/float64(len(opts.Personas)))*10) / 10
		avgFriction = math.Round((totalFriction/float64(len(opts.Personas)))*100) / 100
	}

	overallGrade := "F"
	switch {
	case avgSUS >= 80.3:
		overallGrade = "A"
	case avgSUS >= 74.0:
		overallGrade = "B"
	case avgSUS >= 68.0:
		overallGrade = "C"
	case avgSUS >= 51.0:
		overallGrade = "D"
	}

	mandatoryRedesign := avgFriction > 2.5

	// Evaluate against Nielsen Norman 10 Usability Heuristics
	heuristicFindings := metrics.EvaluateHeuristics(opts.Personas, summaries, stepsMap)

	runRecord.Status = "completed"
	runRecord.AverageSUS = avgSUS
	runRecord.AverageFriction = avgFriction
	runRecord.Grade = overallGrade
	runRecord.MandatoryRedesign = mandatoryRedesign
	runRecord.CompletedAt = time.Now()

	if opts.Store != nil {
		_ = opts.Store.UpdateRun(&runRecord)
		_ = opts.Store.SaveHeuristics(runID, heuristicFindings)
	}

	return &report.ReportData{
		Run:        runRecord,
		Sessions:   sessionRecords,
		StepsMap:   storeStepsMap,
		Heuristics: heuristicFindings,
	}, nil
}
