package metrics

import (
	"fmt"
	"strings"

	"github.com/alexandrmotologa/uxswarm/internal/config"
)

// HeuristicDefinition defines standard Nielsen Norman heuristic metadata.
type HeuristicDefinition struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// NielsenNormanHeuristics lists the 10 standard usability heuristics.
var NielsenNormanHeuristics = []HeuristicDefinition{
	{1, "Visibility of system status", "The design should always keep users informed about what is going on, through appropriate feedback within a reasonable time."},
	{2, "Match between system and the real world", "The design should speak the users' language, using words, phrases, and concepts familiar to the user, rather than internal jargon."},
	{3, "User control and freedom", "Users often perform actions by mistake. They need a clearly marked emergency exit to leave the unwanted action without an extended process."},
	{4, "Consistency and standards", "Users should not have to wonder whether different words, situations, or actions mean the same thing. Follow platform and industry conventions."},
	{5, "Error prevention", "Good error messages are important, but the best designs carefully prevent problems from occurring in the first place."},
	{6, "Recognition rather than recall", "Minimize the user's memory load by making elements, actions, and options visible. The user should not have to remember information across the interface."},
	{7, "Flexibility and efficiency of use", "Shortcuts hidden from novice users may speed up the interaction for the expert user such that the design caters to both inexperienced and experienced users."},
	{8, "Aesthetic and minimalist design", "Interfaces should not contain information that is irrelevant or rarely needed. Every extra unit of information competes with the relevant units."},
	{9, "Help users recognize, diagnose, and recover from errors", "Error messages should be expressed in plain language (no error codes), precisely indicate the problem, and constructively suggest a solution."},
	{10, "Help and documentation", "It's best if the system doesn't need any additional explanation. However, it may be necessary to provide documentation to help users understand their tasks."},
}

// FindingSeverity categorizes usability issue impact.
type FindingSeverity string

const (
	SeverityCritical FindingSeverity = "CRITICAL"
	SeverityMajor    FindingSeverity = "MAJOR"
	SeverityMinor    FindingSeverity = "MINOR"
)

// HeuristicFinding represents an actionable finding mapped to a Nielsen Norman rule.
type HeuristicFinding struct {
	HeuristicID      int             `json:"heuristic_id"`
	HeuristicName    string          `json:"heuristic_name"`
	Severity         FindingSeverity `json:"severity"`
	Title            string          `json:"title"`
	Evidence         string          `json:"evidence"`
	AffectedPersonas []string        `json:"affected_personas"`
	Recommendation   string          `json:"recommendation"`
}

// EvaluateHeuristics inspects persona telemetry to identify heuristic violations and synthesize recommendations.
func EvaluateHeuristics(personas []config.Persona, summaries []*TelemetrySummary, stepsMap map[string][]StepEvent) []HeuristicFinding {
	findings := make([]HeuristicFinding, 0)
	personaMap := make(map[string]config.Persona)
	for _, p := range personas {
		personaMap[p.ID] = p
	}

	// 1. Check for Dead Clicks -> Heuristic 4 (Consistency and standards)
	deadClickPersonas := make([]string, 0)
	totalDeadClicks := 0
	sampleDeadClickDesc := ""
	for i, summary := range summaries {
		if summary.DeadClicks > 0 {
			pID := personas[i].ID
			deadClickPersonas = append(deadClickPersonas, personas[i].Name)
			totalDeadClicks += summary.DeadClicks

			if sampleDeadClickDesc == "" {
				steps := stepsMap[pID]
				for _, s := range steps {
					if s.IsDeadClick {
						sampleDeadClickDesc = fmt.Sprintf("target '%s' at coordinates (%d, %d)", s.TargetDescription, s.X, s.Y)
						break
					}
				}
			}
		}
	}

	if totalDeadClicks > 0 {
		sev := SeverityMinor
		if totalDeadClicks >= 3 {
			sev = SeverityMajor
		}
		findings = append(findings, HeuristicFinding{
			HeuristicID:      4,
			HeuristicName:    NielsenNormanHeuristics[3].Name,
			Severity:         sev,
			Title:            fmt.Sprintf("Users misclicked non-interactive elements %d times", totalDeadClicks),
			Evidence:         fmt.Sprintf("Synthetic users clicked elements resembling buttons or interactive cards that did not trigger any action (%s).", sampleDeadClickDesc),
			AffectedPersonas: deadClickPersonas,
			Recommendation:   "Ensure clickable visual affordances (elevation, underlines, cursor pointer) are applied exclusively to interactive controls, and remove misleading card styling from static text.",
		})
	}

	// 2. Check for Hesitation after action / submission -> Heuristic 1 (Visibility of system status)
	hesitationPersonas := make([]string, 0)
	maxHesitation := 0.0
	for i, summary := range summaries {
		if summary.HesitationCount > 0 {
			hesitationPersonas = append(hesitationPersonas, personas[i].Name)
			if summary.HesitationDurationSeconds > maxHesitation {
				maxHesitation = summary.HesitationDurationSeconds
			}
		}
	}

	if len(hesitationPersonas) > 0 {
		sev := SeverityMajor
		if maxHesitation > 7.0 {
			sev = SeverityCritical
		}
		findings = append(findings, HeuristicFinding{
			HeuristicID:      1,
			HeuristicName:    NielsenNormanHeuristics[0].Name,
			Severity:         sev,
			Title:            fmt.Sprintf("Hesitation latency reached %.1fs during multi-step flow", maxHesitation),
			Evidence:         fmt.Sprintf("%d personas paused for more than 3.5 seconds before determining their next step or after submitting data.", len(hesitationPersonas)),
			AffectedPersonas: hesitationPersonas,
			Recommendation:   "Provide immediate visual feedback when buttons are clicked (e.g. inline spinner, progress bar) and clarify multi-step breadcrumbs so users know what to expect next.",
		})
	}

	// 3. Check for Backtrack Loops (A -> B -> A) -> Heuristic 3 (User control and freedom)
	backtrackPersonas := make([]string, 0)
	totalBacktracks := 0
	for i, summary := range summaries {
		if summary.Backtracks > 0 {
			backtrackPersonas = append(backtrackPersonas, personas[i].Name)
			totalBacktracks += summary.Backtracks
		}
	}

	if totalBacktracks > 0 {
		findings = append(findings, HeuristicFinding{
			HeuristicID:      3,
			HeuristicName:    NielsenNormanHeuristics[2].Name,
			Severity:         SeverityMajor,
			Title:            fmt.Sprintf("Navigation loops detected (%d backtrack transitions)", totalBacktracks),
			Evidence:         "Users navigated forward and immediately reversed their steps, indicating an ambiguous path forward or lack of a clear exit.",
			AffectedPersonas: backtrackPersonas,
			Recommendation:   "Add explicit 'Back' and 'Cancel' navigation controls with clear step preservation, avoiding forced wizard traps that cause users to hit the browser back button.",
		})
	}

	// 4. Check for Drop-offs / Abandonment -> Heuristic 8 (Aesthetic and minimalist design) & Heuristic 5 (Error prevention)
	dropOffPersonas := make([]string, 0)
	dropReasons := make([]string, 0)
	for i, summary := range summaries {
		if summary.DroppedOff {
			dropOffPersonas = append(dropOffPersonas, personas[i].Name)
			if summary.DropOffReason != "" {
				dropReasons = append(dropReasons, fmt.Sprintf("%s: %s", personas[i].Name, summary.DropOffReason))
			}
		}
	}

	if len(dropOffPersonas) > 0 {
		reasonText := "Users ran out of patience before reaching the goal."
		if len(dropReasons) > 0 {
			reasonText = strings.Join(dropReasons, "; ")
		}
		findings = append(findings, HeuristicFinding{
			HeuristicID:      8,
			HeuristicName:    NielsenNormanHeuristics[7].Name,
			Severity:         SeverityCritical,
			Title:            fmt.Sprintf("%d synthetic personas abandoned the session before completion", len(dropOffPersonas)),
			Evidence:         reasonText,
			AffectedPersonas: dropOffPersonas,
			Recommendation:   "Streamline form fields by deferring non-essential questions to post-onboarding, reducing visual noise and making the primary call to action immediately prominent.",
		})
	}

	// 5. Check Mobile-specific friction -> Heuristic 7 (Flexibility and efficiency of use)
	for i, p := range personas {
		if p.Viewport.IsMobile {
			summary := summaries[i]
			if summary.DeadClicks > 0 || summary.HesitationCount > 0 || summary.DroppedOff {
				findings = append(findings, HeuristicFinding{
					HeuristicID:      7,
					HeuristicName:    NielsenNormanHeuristics[6].Name,
					Severity:         SeverityMajor,
					Title:            "Mobile viewport tap friction in one-handed reach zone",
					Evidence:         fmt.Sprintf("Mobile persona %s experienced friction navigating 390x844 viewport.", p.Name),
					AffectedPersonas: []string{p.Name},
					Recommendation:   "Ensure all interactive tap targets satisfy Apple HIG and Material guidelines (minimum 44x44px) and place primary actions within the lower thumb reach zone.",
				})
				break
			}
		}
	}

	return findings
}
