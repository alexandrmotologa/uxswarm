package report

import (
	"encoding/json"
	"fmt"
	"html"
	"os"
	"strings"
	"time"

	"github.com/alexandrmotologa/uxswarm/internal/metrics"
	"github.com/alexandrmotologa/uxswarm/internal/store"
)

// ReportData consolidates all run data required to generate an executive usability report.
type ReportData struct {
	Run        store.RunRecord            `json:"run"`
	Sessions   []store.SessionRecord      `json:"sessions"`
	StepsMap   map[string][]store.StepRecord `json:"steps_map,omitempty"`
	Heuristics []metrics.HeuristicFinding `json:"heuristics"`
}

// GenerateHTML renders a standalone, zero-dependency executive HTML usability report.
func GenerateHTML(data *ReportData) (string, error) {
	var b strings.Builder

	// Calculate completion stats
	completedCount := 0
	for _, s := range data.Sessions {
		if s.GoalCompleted {
			completedCount++
		}
	}
	completionRate := 0.0
	if len(data.Sessions) > 0 {
		completionRate = float64(completedCount) / float64(len(data.Sessions)) * 100.0
	}

	gradeColor := "#16a34a" // Green
	if data.Run.Grade == "C" || data.Run.Grade == "B" {
		gradeColor = "#d97706" // Amber
	} else if data.Run.Grade == "D" || data.Run.Grade == "F" {
		gradeColor = "#dc2626" // Red
	}

	b.WriteString(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>uxswarm Usability Report - ` + html.EscapeString(data.Run.Name) + `</title>
  <style>
    :root {
      --bg: #0f172a;
      --card-bg: #1e293b;
      --card-border: #334155;
      --text-main: #f8fafc;
      --text-muted: #94a3b8;
      --primary: #3b82f6;
      --accent: #8b5cf6;
      --critical: #ef4444;
      --major: #f59e0b;
      --minor: #3b82f6;
    }
    * { box-sizing: border-box; margin: 0; padding: 0; font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; }
    body { background: var(--bg); color: var(--text-main); padding: 2rem 1.5rem; line-height: 1.5; }
    .container { max-width: 1100px; margin: 0 auto; }
    header { margin-bottom: 2.5rem; border-bottom: 1px solid var(--card-border); padding-bottom: 1.5rem; }
    h1 { font-size: 2rem; font-weight: 700; margin-bottom: 0.5rem; }
    .meta-line { color: var(--text-muted); font-size: 0.95rem; }
    
    .grid-stats { display: grid; grid-template-columns: repeat(auto-fit, minmax(220px, 1fr)); gap: 1.25rem; margin-bottom: 2.5rem; }
    .stat-card { background: var(--card-bg); border: 1px solid var(--card-border); border-radius: 12px; padding: 1.5rem; text-align: center; }
    .stat-value { font-size: 2.5rem; font-weight: 800; line-height: 1.1; margin-bottom: 0.35rem; }
    .stat-label { font-size: 0.85rem; color: var(--text-muted); text-transform: uppercase; letter-spacing: 0.05em; }

    section { margin-bottom: 3rem; }
    h2 { font-size: 1.4rem; font-weight: 600; margin-bottom: 1.25rem; display: flex; align-items: center; gap: 0.5rem; }

    table { width: 100%; border-collapse: collapse; background: var(--card-bg); border-radius: 12px; overflow: hidden; border: 1px solid var(--card-border); margin-bottom: 1.5rem; }
    th, td { padding: 1rem 1.25rem; text-align: left; font-size: 0.9rem; }
    th { background: #111827; color: var(--text-muted); font-weight: 600; text-transform: uppercase; font-size: 0.75rem; border-bottom: 1px solid var(--card-border); }
    tr:not(:last-child) td { border-bottom: 1px solid var(--card-border); }
    
    .badge { display: inline-block; padding: 2px 8px; border-radius: 999px; font-size: 0.75rem; font-weight: 600; }
    .badge-success { background: rgba(34, 197, 94, 0.2); color: #4ade80; }
    .badge-danger { background: rgba(239, 68, 68, 0.2); color: #f87171; }
    .badge-critical { background: #991b1b; color: #fecaca; }
    .badge-major { background: #92400e; color: #fde68a; }
    .badge-minor { background: #1e3a8a; color: #bfdbfe; }

    .heuristic-card { background: var(--card-bg); border: 1px solid var(--card-border); border-radius: 12px; padding: 1.5rem; margin-bottom: 1.25rem; }
    .heuristic-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 0.75rem; }
    .heuristic-title { font-size: 1.1rem; font-weight: 600; }
    .heuristic-desc { color: #cbd5e1; font-size: 0.95rem; margin-bottom: 0.75rem; }
    .recommendation-box { background: rgba(59, 130, 246, 0.1); border-left: 4px solid var(--primary); padding: 0.75rem 1rem; border-radius: 0 6px 6px 0; font-size: 0.9rem; margin-top: 0.5rem; }
    .recommendation-label { font-weight: 600; color: #93c5fd; margin-bottom: 0.2rem; }
    
    footer { text-align: center; color: var(--text-muted); font-size: 0.85rem; border-top: 1px solid var(--card-border); padding-top: 1.5rem; }
  </style>
</head>
<body>
  <div class="container">
    <header>
      <h1>Usability Friction & SUS Evaluation</h1>
      <p class="meta-line">Target: ` + html.EscapeString(data.Run.TargetURL) + ` | Study Goal: ` + html.EscapeString(data.Run.Goal) + `</p>
      <p class="meta-line">Evaluated on ` + data.Run.CreatedAt.Format(time.RFC1123) + ` with ` + fmt.Sprintf("%d", data.Run.TotalPersonas) + ` synthetic personas</p>
    </header>

    <div class="grid-stats">
      <div class="stat-card">
        <div class="stat-value" style="color: ` + gradeColor + `;">` + fmt.Sprintf("%.1f", data.Run.AverageSUS) + `</div>
        <div class="stat-label">SUS Benchmark (Grade ` + html.EscapeString(data.Run.Grade) + `)</div>
      </div>
      <div class="stat-card">
        <div class="stat-value" style="color: ` + (map[bool]string{true: "#ef4444", false: "#22c55e"})[data.Run.MandatoryRedesign] + `;">` + fmt.Sprintf("%.2f", data.Run.AverageFriction) + `</div>
        <div class="stat-label">Friction Index (FI ` + (map[bool]string{true: "> 2.5 Redesign", false: "Healthy"})[data.Run.MandatoryRedesign] + `)</div>
      </div>
      <div class="stat-card">
        <div class="stat-value" style="color: #3b82f6;">` + fmt.Sprintf("%.0f%%", completionRate) + `</div>
        <div class="stat-label">Funnel Completion Rate</div>
      </div>
      <div class="stat-card">
        <div class="stat-value" style="color: #a855f7;">` + fmt.Sprintf("%d", len(data.Heuristics)) + `</div>
        <div class="stat-label">Heuristic Findings</div>
      </div>
    </div>

    <section>
      <h2>Persona swarm results</h2>
      <table>
        <thead>
          <tr>
            <th>Persona</th>
            <th>Archetype</th>
            <th>Device</th>
            <th>SUS Score</th>
            <th>Friction</th>
            <th>Hesitations</th>
            <th>Dead Clicks</th>
            <th>Outcome</th>
          </tr>
        </thead>
        <tbody>`)

	for _, s := range data.Sessions {
		deviceStr := "Desktop (1440x900)"
		if s.IsMobile {
			deviceStr = "Mobile (390x844)"
		}

		statusBadge := `<span class="badge badge-success">Completed</span>`
		if s.DroppedOff {
			statusBadge = `<span class="badge badge-danger">Abandoned</span>`
		}

		b.WriteString(fmt.Sprintf(`
          <tr>
            <td><strong>%s</strong></td>
            <td>%s</td>
            <td>%s</td>
            <td><strong>%.1f (%s)</strong></td>
            <td>%.2f</td>
            <td>%d</td>
            <td>%d</td>
            <td>%s</td>
          </tr>`,
			html.EscapeString(s.PersonaName),
			html.EscapeString(s.Archetype),
			deviceStr,
			s.SUSScore, html.EscapeString(s.SUSGrade),
			s.FrictionIndex,
			s.HesitationCount,
			s.DeadClicks,
			statusBadge,
		))
	}

	b.WriteString(`
        </tbody>
      </table>
    </section>

    <section>
      <h2>Nielsen Norman heuristic violations</h2>`)

	if len(data.Heuristics) == 0 {
		b.WriteString(`<p style="color: var(--text-muted);">No critical or major heuristic violations detected in this test run.</p>`)
	} else {
		for _, h := range data.Heuristics {
			sevClass := "badge-minor"
			if h.Severity == metrics.SeverityCritical {
				sevClass = "badge-critical"
			} else if h.Severity == metrics.SeverityMajor {
				sevClass = "badge-major"
			}

			b.WriteString(fmt.Sprintf(`
      <div class="heuristic-card">
        <div class="heuristic-header">
          <div class="heuristic-title">Heuristic %d: %s</div>
          <span class="badge %s">%s</span>
        </div>
        <div style="font-weight: 600; margin-bottom: 0.4rem; color: #f1f5f9;">%s</div>
        <div class="heuristic-desc">%s</div>
        <div style="font-size: 0.85rem; color: var(--text-muted); margin-bottom: 0.5rem;">
          Affected personas: <em>%s</em>
        </div>
        <div class="recommendation-box">
          <div class="recommendation-label">Redesign recommendation</div>
          <div>%s</div>
        </div>
      </div>`,
				h.HeuristicID,
				html.EscapeString(h.HeuristicName),
				sevClass,
				html.EscapeString(string(h.Severity)),
				html.EscapeString(h.Title),
				html.EscapeString(h.Evidence),
				html.EscapeString(strings.Join(h.AffectedPersonas, ", ")),
				html.EscapeString(h.Recommendation),
			))
		}
	}

	b.WriteString(`
    </section>

    <footer>
      Generated autonomously by uxswarm. System Usability Scale methodology (Brooke, 1996).
    </footer>
  </div>
</body>
</html>`)

	return b.String(), nil
}

// GenerateMarkdown creates a human-readable markdown summary following the humanizer discipline.
func GenerateMarkdown(data *ReportData) (string, error) {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("# Usability evaluation: %s\n\n", data.Run.Name))
	b.WriteString(fmt.Sprintf("- Target URL: `%s`\n", data.Run.TargetURL))
	b.WriteString(fmt.Sprintf("- Goal: %s\n", data.Run.Goal))
	b.WriteString(fmt.Sprintf("- Evaluated: %s\n", data.Run.CreatedAt.Format(time.RFC1123)))
	b.WriteString(fmt.Sprintf("- Overall SUS Benchmark: **%.1f** (Grade %s)\n", data.Run.AverageSUS, data.Run.Grade))
	b.WriteString(fmt.Sprintf("- Friction Index: **%.2f** (Mandatory Redesign: %t)\n\n", data.Run.AverageFriction, data.Run.MandatoryRedesign))

	b.WriteString("## Persona swarm results\n\n")
	b.WriteString("| Persona | Archetype | SUS Score | Friction | Hesitations | Outcome |\n")
	b.WriteString("| :--- | :--- | :--- | :--- | :--- | :--- |\n")
	for _, s := range data.Sessions {
		outcome := "Completed"
		if s.DroppedOff {
			outcome = "Abandoned"
		}
		b.WriteString(fmt.Sprintf("| %s | %s | %.1f (%s) | %.2f | %d | %s |\n",
			s.PersonaName, s.Archetype, s.SUSScore, s.SUSGrade, s.FrictionIndex, s.HesitationCount, outcome))
	}

	b.WriteString("\n## Nielsen Norman heuristic findings\n\n")
	if len(data.Heuristics) == 0 {
		b.WriteString("No critical or major heuristic violations detected.\n")
	} else {
		for _, h := range data.Heuristics {
			b.WriteString(fmt.Sprintf("### Heuristic %d: %s [%s]\n\n", h.HeuristicID, h.HeuristicName, h.Severity))
			b.WriteString(fmt.Sprintf("- Problem: %s\n", h.Title))
			b.WriteString(fmt.Sprintf("- Evidence: %s\n", h.Evidence))
			b.WriteString(fmt.Sprintf("- Affected personas: %s\n", strings.Join(h.AffectedPersonas, ", ")))
			b.WriteString(fmt.Sprintf("- Redesign recommendation: %s\n\n", h.Recommendation))
		}
	}

	return b.String(), nil
}

// ExportToFile writes report in specified format (html, markdown, or json) to disk.
func ExportToFile(data *ReportData, format string, outputPath string) error {
	var content []byte
	var err error

	switch strings.ToLower(format) {
	case "html":
		htmlStr, htmlErr := GenerateHTML(data)
		if htmlErr != nil {
			return htmlErr
		}
		content = []byte(htmlStr)
	case "markdown", "md":
		mdStr, mdErr := GenerateMarkdown(data)
		if mdErr != nil {
			return mdErr
		}
		content = []byte(mdStr)
	case "json":
		content, err = json.MarshalIndent(data, "", "  ")
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported report format %q (use html, markdown, or json)", format)
	}

	return os.WriteFile(outputPath, content, 0644)
}
