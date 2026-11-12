# uxswarm

Autonomous synthetic user persona swarm and pre-release UX friction profiler.

uxswarm launches a team of 5 to 10 simulated human personas with distinct psychological archetypes, patience thresholds, and technical literacy levels. The personas navigate web onboarding and checkout funnels via headless Chrome DevTools Protocol (CDP), recording hesitation intervals, dead clicks, and backtrack loops. The system calculates mathematical System Usability Scale (SUS) scores and maps friction points directly to Nielsen Norman usability heuristics.

## Key capabilities

- Diverse synthetic personas: Configurable via YAML with archetypes including rushed executives, skeptical security engineers, mobile thumb navigators, and accessibility focused users.
- Headless CDP automation: Direct browser interaction using Chrome DevTools Protocol for accessibility tree extraction, DOM coordinate mapping, and screenshot capture.
- Cognitive telemetry: Measures hesitation pauses exceeding 3.5 seconds, dead clicks on non-interactive elements, and navigation backtrack loops.
- Mathematical SUS evaluation: Computes standardized 10-item System Usability Scale scores (0 to 100) and friction index metrics.
- A/B Usability Differential: Compares baseline and redesign variants (`uxswarm compare`) to measure delta SUS, delta friction index, and drop-off reductions.
- Automated code patch generator: Synthesizes concrete CSS, Tailwind, and React code fixes (`uxswarm patch`) to resolve identified heuristic violations.
- Eye-tracking gaze simulation: Models F-pattern visual hierarchy and saliency fixations across landing and onboarding screens.
- CI/CD Usability Budgets: Enforces minimum SUS scores and maximum friction thresholds directly in continuous integration workflows with GitHub Step Summary exports.
- Embedded UX studio: Single Go binary serving a React 19 interface for step-by-step playback, inner speech inspection, and click coordinate heatmaps.
- Built-in test target: Includes a local SaaS onboarding application (`uxswarm demo`) with intentional friction points for immediate offline evaluation.

## Installation

### Prerequisites

- Go 1.22 or higher
- Google Chrome or Microsoft Edge installed
- Node.js 20 or higher (for web studio development only)

### Build from source

```bash
git clone https://github.com/alexandrmotologa/uxswarm.git
cd uxswarm
go build -o bin/uxswarm.exe ./cmd/uxswarm
```

## Quick start

### 1. Run the built-in demo flow

Start the local SaaS onboarding app:

```bash
./bin/uxswarm demo --port 3000
```

### 2. Run an automated usability study

Run a 5-persona swarm against the demo flow:

```bash
./bin/uxswarm test \
  --config examples/personas_b2b_saas.yaml \
  --url http://localhost:3000 \
  --goal "Sign up, configure workspace, and complete onboarding" \
  --output ux_report.json
```

### 3. Compare two variants (A/B Usability Diff)

Run an A/B usability comparison across two design variants:

```bash
./bin/uxswarm compare \
  --url-a http://localhost:3000/v1 \
  --url-b http://localhost:3000/v2 \
  --output ab_comparison_report.html
```

### 4. Synthesize engineering code patches

Generate actionable CSS, Tailwind, and React code fixes:

```bash
./bin/uxswarm patch --input ux_report.json
```

### 5. Enforce usability budgets in CI/CD

Fail the pipeline if the flow does not meet usability standards:

```bash
./bin/uxswarm test \
  --url https://staging.myapp.com/signup \
  --min-sus 70.0 \
  --max-friction 1.5 \
  --allow-abandonment=false \
  --github-summary
```

### 6. Start the interactive UX studio

Launch the web studio to review replays, gaze heatmaps, and inner monologues:

```bash
./bin/uxswarm serve --port 8080 --db uxswarm.db
```

Open `http://localhost:8080` in your browser.

## Command reference

- `uxswarm serve`: Starts the HTTP API daemon and embedded React studio.
  - `--port`: HTTP listen port (default: 8080).
  - `--db`: SQLite database file path (default: uxswarm.db).
- `uxswarm test`: Executes a usability study against a target URL.
  - `--config`: Path to persona configuration YAML file.
  - `--url`: Target web application URL.
  - `--goal`: Natural language task description.
  - `--output`: File path for JSON result export.
  - `--report`: Optional path for standalone HTML report export.
  - `--headless`: Run browser in headless mode (default: true).
  - `--min-sus`: Minimum required average SUS score (fails CI if not met).
  - `--max-friction`: Maximum allowed Friction Index (fails CI if exceeded).
  - `--allow-abandonment`: Permitted persona drop-offs (fails CI if false and drop-offs occur).
  - `--github-summary`: Writes report to GITHUB_STEP_SUMMARY environment variable.
- `uxswarm compare`: Executes an A/B usability study comparing two variants.
  - `--url-a`: Baseline variant URL (Variant A).
  - `--url-b`: Candidate redesign URL (Variant B).
  - `--output`: Path for comparison report destination.
  - `--format`: Format (`html`, `markdown`, or `json`).
- `uxswarm patch`: Generates CSS, Tailwind, and React code fixes for findings.
  - `--input`: Path to JSON test report.
  - `--format`: Filter format (`css`, `tailwind`, `html`, `react`, or `all`).
- `uxswarm audit`: Runs a fast single-page heuristic scan.
  - `--url`: Page URL to audit.
  - `--device`: Target device profile (`desktop` or `mobile`).
  - `--output`: File path for audit findings.
- `uxswarm report`: Formats test results into human-readable documents.
  - `--run-id`: Run identifier from SQLite database.
  - `--input`: Alternatively, path to a JSON test report.
  - `--format`: Output format (`html`, `markdown`, or `json`).
  - `--output`: Output destination path.
- `uxswarm demo`: Runs the built-in mock SaaS onboarding application.
  - `--port`: HTTP listen port (default: 3000).
- `uxswarm version`: Prints version and build information.

## Repository layout

```
uxswarm/
├── cmd/
│   └── uxswarm/              # CLI command entrypoints
├── internal/
│   ├── api/                  # HTTP REST and SSE streaming handlers
│   ├── browser/              # CDP browser session and action actuator
│   ├── compare/              # A/B Usability Differential engine
│   ├── config/               # Suite and persona YAML loader
│   ├── demo/                 # Built-in friction benchmark web app
│   ├── gaze/                 # Simulated eye-tracking visual attention model
│   ├── metrics/              # SUS scoring and Nielsen Norman mapping
│   ├── patch/                # Automated engineering code patch generator
│   ├── persona/              # Cognitive perception and decision brain
│   ├── report/               # HTML and Markdown report generator
│   └── store/                # SQLite database and query layer
├── web/                      # React 19 embedded UX studio
│   ├── src/
│   │   ├── components/       # Player, heatmap, A/B compare, scorecards
│   │   └── pages/            # Studio, Dashboard, Persona Builder
│   ├── index.html
│   └── vite.config.ts
├── examples/                 # Persona suite definitions
├── docs/                     # Specifications and methodology
├── go.mod
└── README.md
```

## License

MIT License. See [LICENSE](LICENSE) for details.
