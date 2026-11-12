package metrics

import (
	"fmt"
	"math"

	"github.com/alexandrmotologa/uxswarm/internal/config"
)

// Standard SUS 10 Questions according to John Brooke (1996).
var SUSQuestions = [10]string{
	"I think that I would like to use this system frequently.",
	"I found the system unnecessarily complex.",
	"I thought the system was easy to use.",
	"I think that I would need the support of a technical person to be able to use this system.",
	"I found the various functions in this system were well integrated.",
	"I thought there was too much inconsistency in this system.",
	"I would imagine that most people would learn to use this system very quickly.",
	"I found the system very cumbersome to use.",
	"I felt very confident using the system.",
	"I needed to learn a lot of things before I could get going with this system.",
}

// SUSGrade represents qualitative SUS benchmark grading.
type SUSGrade string

const (
	GradeA SUSGrade = "A" // >= 80.3 (Excellent UX)
	GradeB SUSGrade = "B" // 74.0 - 80.2 (Good UX)
	GradeC SUSGrade = "C" // 68.0 - 73.9 (OK / Average Benchmark)
	GradeD SUSGrade = "D" // 51.0 - 67.9 (Poor UX / Significant Friction)
	GradeF SUSGrade = "F" // < 51.0 (Critical Usability Failures)
)

// SUSResult holds computed SUS scores and item responses for a persona session.
type SUSResult struct {
	Score      float64   `json:"score"`      // 0 to 100
	Grade      SUSGrade  `json:"grade"`      // A, B, C, D, F
	Adjective  string    `json:"adjective"`  // Excellent, Good, OK, Poor, Terrible
	Responses  [10]int   `json:"responses"`  // Raw answers (1 to 5)
	Normalized [10]int   `json:"normalized"` // Odd: R-1, Even: 5-R
	Questions  [10]string `json:"questions"`
}

// CalculateSUS computes the standard System Usability Scale score from 10 responses (1 to 5).
func CalculateSUS(responses [10]int) (*SUSResult, error) {
	for i, r := range responses {
		if r < 1 || r > 5 {
			return nil, fmt.Errorf("item %d response %d out of valid range (1-5)", i+1, r)
		}
	}

	var normalized [10]int
	sum := 0

	for i := 0; i < 10; i++ {
		if (i+1)%2 == 1 { // Odd items: 1, 3, 5, 7, 9
			normalized[i] = responses[i] - 1
		} else { // Even items: 2, 4, 6, 8, 10
			normalized[i] = 5 - responses[i]
		}
		sum += normalized[i]
	}

	score := float64(sum) * 2.5
	if score > 100.0 {
		score = 100.0
	}
	if score < 0.0 {
		score = 0.0
	}

	var grade SUSGrade
	var adjective string

	switch {
	case score >= 80.3:
		grade = GradeA
		adjective = "Excellent"
	case score >= 74.0:
		grade = GradeB
		adjective = "Good"
	case score >= 68.0:
		grade = GradeC
		adjective = "Acceptable"
	case score >= 51.0:
		grade = GradeD
		adjective = "Poor"
	default:
		grade = GradeF
		adjective = "Critical Failures"
	}

	return &SUSResult{
		Score:      math.Round(score*10) / 10,
		Grade:      grade,
		Adjective:  adjective,
		Responses:  responses,
		Normalized: normalized,
		Questions:  SUSQuestions,
	}, nil
}

// SynthesizeSUS generates authentic SUS responses for a persona based on their session telemetry.
func SynthesizeSUS(persona config.Persona, tel *TelemetrySummary) *SUSResult {
	// Base baseline rating starts neutral (3 on all items)
	// Base satisfaction shifts based on goal completion, frustration, hesitation, and dead clicks
	satisfaction := 0.5 // 0.0 (furious/blocked) to 1.0 (smooth/delighted)

	if tel.GoalCompleted {
		satisfaction += 0.25
	} else if tel.DroppedOff {
		satisfaction -= 0.35
	}

	// Hesitation impact: high hesitation reduces satisfaction
	hesitationRatio := 0.0
	if tel.TotalDurationSeconds > 0 {
		hesitationRatio = tel.HesitationDurationSeconds / tel.TotalDurationSeconds
	}
	satisfaction -= hesitationRatio * 0.25

	// Dead clicks impact
	satisfaction -= math.Min(float64(tel.DeadClicks)*0.08, 0.3)

	// Backtrack loop impact
	satisfaction -= math.Min(float64(tel.Backtracks)*0.06, 0.2)

	// Clamp satisfaction to [0.05, 0.95]
	if satisfaction < 0.05 {
		satisfaction = 0.05
	}
	if satisfaction > 0.95 {
		satisfaction = 0.95
	}

	// Low tech literacy personas get more frustrated when friction occurs
	if persona.TechnicalLiteracy <= 3 && satisfaction < 0.5 {
		satisfaction -= 0.1
		if satisfaction < 0.05 {
			satisfaction = 0.05
		}
	}

	var responses [10]int
	// Map satisfaction level to realistic response values (1-5)
	// Odd questions: positive sentiment (higher is better)
	// Even questions: negative sentiment (lower is better)
	for i := 0; i < 10; i++ {
		isPositiveItem := (i+1)%2 == 1
		if isPositiveItem {
			val := int(math.Round(1.0 + satisfaction*4.0))
			// Add archetype variation
			if persona.Archetype == config.ArchetypeHighImpatienceLowTech && !tel.GoalCompleted {
				val = int(math.Max(1, float64(val-1)))
			}
			responses[i] = clampResponse(val)
		} else {
			// Negative item: when satisfaction is high, response should be low (1 or 2)
			val := int(math.Round(5.0 - satisfaction*4.0))
			if tel.DeadClicks > 0 && (i == 1 || i == 7) { // Complexity / cumbersome
				val = int(math.Min(5, float64(val+1)))
			}
			responses[i] = clampResponse(val)
		}
	}

	res, _ := CalculateSUS(responses)
	return res
}

func clampResponse(v int) int {
	if v < 1 {
		return 1
	}
	if v > 5 {
		return 5
	}
	return v
}
