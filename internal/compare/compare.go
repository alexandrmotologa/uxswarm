package compare

import (
	"encoding/json"
	"fmt"
	"html"
	"math"
	"os"
	"strings"
	"time"

	"github.com/alexandrmotologa/uxswarm/internal/report"
)

// PersonaDelta compares individual persona outcomes between variants.
type PersonaDelta struct {
	PersonaID     string  `json:"persona_id"`
	PersonaName   string  `json:"persona_name"`
	Archetype     string  `json:"archetype"`
	SUSA          float64 `json:"sus_a"`
	SUSB          float64 `json:"sus_b"`
	DeltaSUS      float64 `json:"delta_sus"`
	FrictionA     float64 `json:"friction_a"`
	FrictionB     float64 `json:"friction_b"`
	DeltaFriction float64 `json:"delta_friction"`
	OutcomeA      string  `json:"outcome_a"`
	OutcomeB      string  `json:"outcome_b"`
}

// ComparisonReport holds statistical differential analysis between two UX variants.
type ComparisonReport struct {
	VariantAName      string         `json:"variant_a_name"`
	VariantBName      string         `json:"variant_b_name"`
	URLA              string         `json:"url_a"`
	URLB              string         `json:"url_b"`
	Goal              string         `json:"goal"`
	CreatedAt         time.Time      `json:"created_at"`
	SUSA              float64        `json:"sus_a"`
	SUSB              float64        `json:"sus_b"`
	DeltaSUS          float64        `json:"delta_sus"`
	GradeA            string         `json:"grade_a"`
	GradeB            string         `json:"grade_b"`
	FrictionA         float64        `json:"friction_a"`
	FrictionB         float64        `json:"friction_b"`
	DeltaFriction     float64        `json:"delta_friction"`
	DropOffsA         int            `json:"drop_offs_a"`
	DropOffsB         int            `json:"drop_offs_b"`
	DeltaDropOffs     int            `json:"delta_drop_offs"`
	HesitationsA      int            `json:"hesitations_a"`
	HesitationsB      int            `json:"hesitations_b"`
	DeltaHesitations  int            `json:"delta_hesitations"`
	DeadClicksA       int            `json:"dead_clicks_a"`
	DeadClicksB       int            `json:"dead_clicks_b"`
	DeltaDeadClicks   int            `json:"delta_dead_clicks"`
	Winner            string         `json:"winner"`
	Recommendation    string         `json:"recommendation"`
	PersonaDeltas     []PersonaDelta `json:"persona_deltas"`
	ResolvedIssues    []string       `json:"resolved_issues"`
	RemainingIssues   []string       `json:"remaining_issues"`
}

// CompareVariants computes the usability differential between Variant A and Variant B.
func CompareVariants(repA *report.ReportData, repB *report.ReportData) *ComparisonReport {
	deltaSUS := math.Round((repB.Run.AverageSUS-repA.Run.AverageSUS)*10) / 10
	deltaFriction := math.Round((repB.Run.AverageFriction-repA.Run.AverageFriction)*100) / 100

	dropOffsA, dropOffsB := 0, 0
	hesitationsA, hesitationsB := 0, 0
	deadClicksA, deadClicksB := 0, 0

	for _, s := range repA.Sessions {
		if s.DroppedOff {
			dropOffsA++
		}
		hesitationsA += s.HesitationCount
		deadClicksA += s.DeadClicks
	}

	for _, s := range repB.Sessions {
		if s.DroppedOff {
			dropOffsB++
		}
		hesitationsB += s.HesitationCount
		deadClicksB += s.DeadClicks
	}

	// Match personas
	personaMapB := make(map[string]report.ReportData)
	_ = personaMapB
	sessBMap := make(map[string]report.ReportData)
	_ = sessBMap

	personaDeltas := make([]PersonaDelta, 0)
	for _, sa := range repA.Sessions {
		var matchedB *report.ReportData
		_ = matchedB
		var foundB bool
		var sbSUS, sbFriction float64
		var sbOutcome = "Not Run"

		for _, sb := range repB.Sessions {
			if sb.PersonaID == sa.PersonaID || sb.PersonaName == sa.PersonaName {
				foundB = true
				sbSUS = sb.SUSScore
				sbFriction = sb.FrictionIndex
				if sb.GoalCompleted {
					sbOutcome = "Completed"
				} else if sb.DroppedOff {
					sbOutcome = "Abandoned"
				}
				break
			}
		}

		outcomeA := "Completed"
		if sa.DroppedOff {
			outcomeA = "Abandoned"
		}

		dSUS := 0.0
		dFriction := 0.0
		if foundB {
			dSUS = math.Round((sbSUS-sa.SUSScore)*10) / 10
			dFriction = math.Round((sbFriction-sa.FrictionIndex)*100) / 100
		}

		personaDeltas = append(personaDeltas, PersonaDelta{
			PersonaID:     sa.PersonaID,
			PersonaName:   sa.PersonaName,
			Archetype:     sa.Archetype,
			SUSA:          sa.SUSScore,
			SUSB:          sbSUS,
			DeltaSUS:      dSUS,
			FrictionA:     sa.FrictionIndex,
			FrictionB:     sbFriction,
			DeltaFriction: dFriction,
			OutcomeA:      outcomeA,
			OutcomeB:      sbOutcome,
		})
	}

	// Determine winner
	winner := "Variant B (Recommended)"
	recommendation := fmt.Sprintf("Variant B increased SUS score by +%.1f points and reduced friction by %.2f. Recommended for production deployment.", deltaSUS, math.Abs(deltaFriction))
	if deltaSUS < 0 || deltaFriction > 0 {
		if deltaSUS < -5.0 {
			winner = "Variant A (Baseline Preferred)"
			recommendation = "Variant B introduced regression in usability. Maintain Variant A and address identified friction points."
		} else {
			winner = "Inconclusive / Neutral"
			recommendation = "No statistically material difference in usability metrics between Variant A and Variant B."
		}
	}

	// Compare heuristics: what was resolved in B?
	resolved := make([]string, 0)
	remaining := make([]string, 0)

	for _, ha := range repA.Heuristics {
		foundInB := false
		for _, hb := range repB.Heuristics {
			if hb.HeuristicID == ha.HeuristicID {
				foundInB = true
				break
			}
		}
		if !foundInB {
			resolved = append(resolved, fmt.Sprintf("Heuristic %d (%s): %s", ha.HeuristicID, ha.HeuristicName, ha.Title))
		}
	}

	for _, hb := range repB.Heuristics {
		remaining = append(remaining, fmt.Sprintf("Heuristic %d (%s): %s [%s]", hb.HeuristicID, hb.HeuristicName, hb.Title, hb.Severity))
	}

	return &ComparisonReport{
		VariantAName:     repA.Run.Name,
		VariantBName:     repB.Run.Name,
		URLA:             repA.Run.TargetURL,
		URLB:             repB.Run.TargetURL,
		Goal:             repA.Run.Goal,
		CreatedAt:        time.Now(),
		SUSA:             repA.Run.AverageSUS,
		SUSB:             repB.Run.AverageSUS,
		DeltaSUS:         deltaSUS,
		GradeA:           repA.Run.Grade,
		GradeB:           repB.Run.Grade,
		FrictionA:        repA.Run.AverageFriction,
		FrictionB:        repB.Run.AverageFriction,
		DeltaFriction:    deltaFriction,
		DropOffsA:        dropOffsA,
		DropOffsB:        dropOffsB,
		DeltaDropOffs:    dropOffsB - dropOffsA,
		HesitationsA:     hesitationsA,
		HesitationsB:     hesitationsB,
		DeltaHesitations: hesitationsB - hesitationsA,
		DeadClicksA:      deadClicksA,
		DeadClicksB:      deadClicksB,
		DeltaDeadClicks:  deadClicksB - deadClicksA,
		Winner:           winner,
		Recommendation:   recommendation,
		PersonaDeltas:    personaDeltas,
		ResolvedIssues:   resolved,
		RemainingIssues:  remaining,
	}
}

// GenerateDiffHTML renders an interactive side-by-side executive A/B report.
func GenerateDiffHTML(cmp *ComparisonReport) (string, error) {
	var b strings.Builder

	deltaSign := "+"
	if cmp.DeltaSUS < 0 {
		deltaSign = ""
	}

	susColor := "#16a34a"
	if cmp.DeltaSUS < 0 {
		susColor = "#dc2626"
	}

	b.WriteString(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>uxswarm A/B Usability Differential Report</title>
  <style>
    * { box-sizing: border-box; margin: 0; padding: 0; font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; }
    body { background: #090d16; color: #f8fafc; padding: 2.5rem 1.5rem; line-height: 1.5; }
    .container { max-width: 1100px; margin: 0 auto; }
    header { border-bottom: 1px solid #1e293b; padding-bottom: 1.5rem; margin-bottom: 2rem; }
    h1 { font-size: 1.85rem; font-weight: 800; margin-bottom: 0.5rem; }
    .winner-banner { background: rgba(37,99,235,0.15); border: 1px solid #3b82f6; border-radius: 12px; padding: 1.25rem 1.5rem; margin-bottom: 2rem; }
    .grid-cmp { display: grid; grid-template-columns: 1fr 1fr 1fr; gap: 1.25rem; margin-bottom: 2.5rem; }
    .card { background: #141e33; border: 1px solid #1e293b; border-radius: 12px; padding: 1.5rem; text-align: center; }
    .metric-val { font-size: 2.5rem; font-weight: 800; margin-bottom: 0.25rem; }
    .metric-lbl { font-size: 0.8rem; color: #94a3b8; text-transform: uppercase; letter-spacing: 0.05em; }
    table { width: 100%; border-collapse: collapse; background: #141e33; border: 1px solid #1e293b; border-radius: 12px; overflow: hidden; margin-bottom: 2rem; }
    th, td { padding: 1rem 1.25rem; text-align: left; font-size: 0.9rem; }
    th { background: #020617; color: #94a3b8; text-transform: uppercase; font-size: 0.75rem; border-bottom: 1px solid #1e293b; }
    tr:not(:last-child) td { border-bottom: 1px solid #1e293b; }
    .resolved-box { background: rgba(22,163,74,0.1); border-left: 4px solid #16a34a; padding: 1rem; border-radius: 0 8px 8px 0; margin-bottom: 1.5rem; font-size: 0.9rem; }
  </style>
</head>
<body>
  <div class="container">
    <header>
      <h1>A/B Usability Differential Evaluation</h1>
      <p style="color: #94a3b8; font-size: 0.9rem;">
        Variant A: ` + html.EscapeString(cmp.URLA) + ` &bull; Variant B: ` + html.EscapeString(cmp.URLB) + `
      </p>
    </header>

    <div class="winner-banner">
      <div style="font-size: 0.8rem; text-transform: uppercase; color: #93c5fd; font-weight: 700;">Evaluation Outcome</div>
      <div style="font-size: 1.4rem; font-weight: 800; color: #f8fafc; margin: 0.25rem 0;">` + html.EscapeString(cmp.Winner) + `</div>
      <div style="font-size: 0.9rem; color: #cbd5e1;">` + html.EscapeString(cmp.Recommendation) + `</div>
    </div>

    <div class="grid-cmp">
      <div class="card">
        <div class="metric-val" style="color: ` + susColor + `;">` + fmt.Sprintf("%s%.1f", deltaSign, cmp.DeltaSUS) + `</div>
        <div class="metric-lbl">&Delta; SUS Benchmark (A: ` + fmt.Sprintf("%.1f", cmp.SUSA) + ` &rarr; B: ` + fmt.Sprintf("%.1f", cmp.SUSB) + `)</div>
      </div>
      <div class="card">
        <div class="metric-val" style="color: ` + (map[bool]string{true: "#16a34a", false: "#dc2626"})[cmp.DeltaFriction <= 0] + `;">` + fmt.Sprintf("%.2f", cmp.DeltaFriction) + `</div>
        <div class="metric-lbl">&Delta; Friction Index (A: ` + fmt.Sprintf("%.2f", cmp.FrictionA) + ` &rarr; B: ` + fmt.Sprintf("%.2f", cmp.FrictionB) + `)</div>
      </div>
      <div class="card">
        <div class="metric-val" style="color: ` + (map[bool]string{true: "#16a34a", false: "#dc2626"})[cmp.DeltaDropOffs <= 0] + `;">` + fmt.Sprintf("%d", cmp.DeltaDropOffs) + `</div>
        <div class="metric-lbl">&Delta; Session Drop-offs (A: ` + fmt.Sprintf("%d", cmp.DropOffsA) + ` &rarr; B: ` + fmt.Sprintf("%d", cmp.DropOffsB) + `)</div>
      </div>
    </div>

    <h2 style="font-size: 1.25rem; font-weight: 700; margin-bottom: 1rem;">Persona Swarm Breakdown</h2>
    <table>
      <thead>
        <tr>
          <th>Persona</th>
          <th>Variant A (SUS)</th>
          <th>Variant B (SUS)</th>
          <th>&Delta; SUS</th>
          <th>Variant A (FI)</th>
          <th>Variant B (FI)</th>
          <th>Outcome (A &rarr; B)</th>
        </tr>
      </thead>
      <tbody>`)

	for _, p := range cmp.PersonaDeltas {
		sign := "+"
		if p.DeltaSUS < 0 {
			sign = ""
		}
		pColor := "#16a34a"
		if p.DeltaSUS < 0 {
			pColor = "#dc2626"
		}

		b.WriteString(fmt.Sprintf(`
        <tr>
          <td><strong>%s</strong><div style="font-size: 0.75rem; color: #94a3b8;">%s</div></td>
          <td>%.1f</td>
          <td><strong>%.1f</strong></td>
          <td><strong style="color: %s;">%s%.1f</strong></td>
          <td>%.2f</td>
          <td><strong>%.2f</strong></td>
          <td>%s &rarr; %s</td>
        </tr>`,
			html.EscapeString(p.PersonaName),
			html.EscapeString(p.Archetype),
			p.SUSA, p.SUSB,
			pColor, sign, p.DeltaSUS,
			p.FrictionA, p.FrictionB,
			p.OutcomeA, p.OutcomeB,
		))
	}

	b.WriteString(`
      </tbody>
    </table>`)

	if len(cmp.ResolvedIssues) > 0 {
		b.WriteString(`<h2 style="font-size: 1.25rem; font-weight: 700; margin-bottom: 0.75rem;">Resolved Usability Friction</h2>
    <div class="resolved-box">
      <ul style="padding-left: 1.25rem;">`)
		for _, res := range cmp.ResolvedIssues {
			b.WriteString("<li>" + html.EscapeString(res) + "</li>")
		}
		b.WriteString(`</ul>
    </div>`)
	}

	b.WriteString(`
  </div>
</body>
</html>`)

	return b.String(), nil
}

// GenerateDiffMarkdown renders an executive markdown differential summary.
func GenerateDiffMarkdown(cmp *ComparisonReport) (string, error) {
	var b strings.Builder
	b.WriteString("# A/B Usability Differential Report\n\n")
	b.WriteString(fmt.Sprintf("- Variant A: `%s`\n", cmp.URLA))
	b.WriteString(fmt.Sprintf("- Variant B: `%s`\n", cmp.URLB))
	b.WriteString(fmt.Sprintf("- Winner: **%s**\n", cmp.Winner))
	b.WriteString(fmt.Sprintf("- Recommendation: %s\n\n", cmp.Recommendation))

	b.WriteString(fmt.Sprintf("## Overall Metrics Differential\n\n"))
	b.WriteString(fmt.Sprintf("| Metric | Variant A | Variant B | Differential (&Delta;) |\n"))
	b.WriteString(fmt.Sprintf("| :--- | :--- | :--- | :--- |\n"))
	b.WriteString(fmt.Sprintf("| SUS Benchmark | %.1f (Grade %s) | %.1f (Grade %s) | **%+.1f pts** |\n", cmp.SUSA, cmp.GradeA, cmp.SUSB, cmp.GradeB, cmp.DeltaSUS))
	b.WriteString(fmt.Sprintf("| Friction Index | %.2f | %.2f | **%+.2f** |\n", cmp.FrictionA, cmp.FrictionB, cmp.DeltaFriction))
	b.WriteString(fmt.Sprintf("| Session Drop-offs | %d | %d | **%+d** |\n", cmp.DropOffsA, cmp.DropOffsB, cmp.DeltaDropOffs))
	b.WriteString(fmt.Sprintf("| Dead Clicks | %d | %d | **%+d** |\n\n", cmp.DeadClicksA, cmp.DeadClicksB, cmp.DeltaDeadClicks))

	b.WriteString("## Persona Swarm Breakdown\n\n")
	b.WriteString("| Persona | Variant A (SUS) | Variant B (SUS) | &Delta; SUS | Outcome A &rarr; B |\n")
	b.WriteString("| :--- | :--- | :--- | :--- | :--- |\n")
	for _, p := range cmp.PersonaDeltas {
		b.WriteString(fmt.Sprintf("| %s | %.1f | %.1f | %+.1f | %s &rarr; %s |\n",
			p.PersonaName, p.SUSA, p.SUSB, p.DeltaSUS, p.OutcomeA, p.OutcomeB))
	}

	return b.String(), nil
}

// ExportDiffToFile writes the A/B comparison report in specified format.
func ExportDiffToFile(cmp *ComparisonReport, format string, outputPath string) error {
	var content []byte
	var err error

	switch strings.ToLower(format) {
	case "html":
		s, errH := GenerateDiffHTML(cmp)
		if errH != nil {
			return errH
		}
		content = []byte(s)
	case "markdown", "md":
		s, errM := GenerateDiffMarkdown(cmp)
		if errM != nil {
			return errM
		}
		content = []byte(s)
	case "json":
		content, err = json.MarshalIndent(cmp, "", "  ")
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported diff format %q (use html, markdown, or json)", format)
	}

	return os.WriteFile(outputPath, content, 0644)
}
