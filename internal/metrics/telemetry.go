package metrics

import (
	"math"
	"time"
)

// HesitationThresholdSeconds is the cognitive delay threshold where pauses indicate confusion.
const HesitationThresholdSeconds = 3.5

// StepEvent represents a discrete interaction or cognitive event within a persona session.
type StepEvent struct {
	Index             int       `json:"index"`
	Timestamp         time.Time `json:"timestamp"`
	ActionType        string    `json:"action_type"` // navigate, click, type, wait, hesitate, drop_off, complete
	TargetSelector    string    `json:"target_selector,omitempty"`
	TargetDescription string    `json:"target_description,omitempty"`
	TargetRole        string    `json:"target_role,omitempty"`
	X                 int       `json:"x"`
	Y                 int       `json:"y"`
	IsDeadClick       bool      `json:"is_dead_click"`
	HesitationSeconds float64   `json:"hesitation_seconds"`
	URL               string    `json:"url"`
	InnerSpeech       string    `json:"inner_speech"`
	Frustration       float64   `json:"frustration"` // 0.0 (calm) to 1.0 (rage/abandonment)
	ScreenshotPath    string    `json:"screenshot_path,omitempty"`
	FrameDataURI      string    `json:"frame_data_uri,omitempty"` // Base64 data for studio replay
}

// TelemetrySummary aggregates session-level usability metrics.
type TelemetrySummary struct {
	TotalSteps               int       `json:"total_steps"`
	TotalDurationSeconds     float64   `json:"total_duration_seconds"`
	HesitationCount          int       `json:"hesitation_count"`
	HesitationDurationSeconds float64  `json:"hesitation_duration_seconds"`
	DeadClicks               int       `json:"dead_clicks"`
	Backtracks               int       `json:"backtracks"`
	GoalCompleted            bool      `json:"goal_completed"`
	DroppedOff               bool      `json:"dropped_off"`
	DropOffReason            string    `json:"drop_off_reason,omitempty"`
	FrictionIndex            float64   `json:"friction_index"`
	MandatoryRedesign        bool      `json:"mandatory_redesign"` // true if FI > 2.5
	MaxFrustration           float64   `json:"max_frustration"`
	URLSequence              []string  `json:"url_sequence"`
}

// TelemetryTracker tracks events during an individual persona's session.
type TelemetryTracker struct {
	startTime   time.Time
	lastStepTime time.Time
	steps       []StepEvent
	urls        []string
	deadClicks  int
	backtracks  int
	completed   bool
	droppedOff  bool
	dropReason  string
}

// NewTelemetryTracker initializes a fresh telemetry tracker.
func NewTelemetryTracker() *TelemetryTracker {
	now := time.Now()
	return &TelemetryTracker{
		startTime:    now,
		lastStepTime: now,
		steps:        make([]StepEvent, 0),
		urls:         make([]string, 0),
	}
}

// RecordStep appends a new interaction event and updates hesitation / loop detection.
func (t *TelemetryTracker) RecordStep(evt StepEvent) {
	now := time.Now()
	if evt.Timestamp.IsZero() {
		evt.Timestamp = now
	}
	evt.Index = len(t.steps) + 1

	// Check hesitation since last action if not already calculated
	if evt.HesitationSeconds == 0 && !t.lastStepTime.IsZero() {
		diff := evt.Timestamp.Sub(t.lastStepTime).Seconds()
		if diff >= HesitationThresholdSeconds {
			evt.HesitationSeconds = math.Round(diff*10) / 10
		}
	}

	if evt.IsDeadClick {
		t.deadClicks++
	}

	// Track URL transitions and detect backtrack loops (A -> B -> A)
	if evt.URL != "" {
		n := len(t.urls)
		if n == 0 || t.urls[n-1] != evt.URL {
			t.urls = append(t.urls, evt.URL)
			if len(t.urls) >= 3 {
				// Check for cycle: current URL equals URL two steps back
				if t.urls[len(t.urls)-1] == t.urls[len(t.urls)-3] {
					t.backtracks++
				}
			}
		}
	}

	if evt.ActionType == "complete" {
		t.completed = true
	}
	if evt.ActionType == "drop_off" {
		t.droppedOff = true
		t.dropReason = evt.InnerSpeech
	}

	t.steps = append(t.steps, evt)
	t.lastStepTime = evt.Timestamp
}

// MarkCompleted flags the session as successfully finishing the test goal.
func (t *TelemetryTracker) MarkCompleted() {
	t.completed = true
}

// MarkDroppedOff flags the session as abandoned due to frustration.
func (t *TelemetryTracker) MarkDroppedOff(reason string) {
	t.droppedOff = true
	t.dropReason = reason
}

// Steps returns the recorded step events.
func (t *TelemetryTracker) Steps() []StepEvent {
	return t.steps
}

// Summary calculates the aggregated metrics and mathematical Friction Index.
func (t *TelemetryTracker) Summary() *TelemetrySummary {
	totalDuration := time.Since(t.startTime).Seconds()
	if len(t.steps) > 0 {
		lastTime := t.steps[len(t.steps)-1].Timestamp
		totalDuration = lastTime.Sub(t.startTime).Seconds()
	}
	if totalDuration <= 0 {
		totalDuration = 0.1
	}

	hesitationCount := 0
	totalHesitation := 0.0
	maxFrustration := 0.0

	for _, s := range t.steps {
		if s.HesitationSeconds >= HesitationThresholdSeconds {
			hesitationCount++
			totalHesitation += s.HesitationSeconds
		}
		if s.Frustration > maxFrustration {
			maxFrustration = s.Frustration
		}
	}

	// Friction Index formulation:
	// FI = w1 * (T_hesitation / T_total) + w2 * N_misclicks + w3 * N_backtracks
	// w1 = 0.35, w2 = 0.40, w3 = 0.25
	hesitationRatio := totalHesitation / totalDuration
	if hesitationRatio > 1.0 {
		hesitationRatio = 1.0
	}

	fi := (0.35 * hesitationRatio) + (0.40 * float64(t.deadClicks)) + (0.25 * float64(t.backtracks))
	fi = math.Round(fi*100) / 100

	return &TelemetrySummary{
		TotalSteps:               len(t.steps),
		TotalDurationSeconds:     math.Round(totalDuration*10) / 10,
		HesitationCount:          hesitationCount,
		HesitationDurationSeconds: math.Round(totalHesitation*10) / 10,
		DeadClicks:               t.deadClicks,
		Backtracks:               t.backtracks,
		GoalCompleted:            t.completed,
		DroppedOff:               t.droppedOff,
		DropOffReason:            t.dropReason,
		FrictionIndex:            fi,
		MandatoryRedesign:        fi > 2.5,
		MaxFrustration:           math.Round(maxFrustration*100) / 100,
		URLSequence:              t.urls,
	}
}
