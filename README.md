# uxswarm

Autonomous synthetic user persona swarm and pre-release UX friction profiler.

uxswarm launches a team of 5 to 10 simulated human personas with distinct psychological archetypes, patience thresholds, and technical literacy levels. The personas navigate web onboarding and checkout funnels via headless Chrome DevTools Protocol (CDP), recording hesitation intervals, dead clicks, and backtrack loops. The system calculates mathematical System Usability Scale (SUS) scores and maps friction points directly to Nielsen Norman usability heuristics.

## Key capabilities

- Diverse synthetic personas: Configurable via YAML with archetypes including rushed executives, skeptical security engineers, mobile thumb navigators, and accessibility focused users.
- Headless CDP automation: Direct browser interaction using Chrome DevTools Protocol for accessibility tree extraction, DOM coordinate mapping, and screenshot capture.
- Cognitive telemetry: Measures hesitation pauses exceeding 3.5 seconds, dead clicks on non-interactive elements, and navigation backtrack loops.
- Mathematical SUS evaluation: Computes standardized 10-item System Usability Scale scores (0 to 100) and friction index metrics.
- Heuristic mapping: Categorizes usability barriers under the 10 Nielsen Norman heuristics with concrete recommendations.
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

### 3. Generate an executive HTML report

Generate a standalone HTML report:

```bash
./bin/uxswarm report --input ux_report.json --format html --output usability_report.html
```

### 4. Start the interactive UX studio

Launch the web studio to review replays, heatmaps, and inner monologues:

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
  - `--headless`: Run browser in headless mode (default: true).
  - `--db`: SQLite database file path.
- `uxswarm audit`: Runs a fast single-page heuristic scan.
  - `--url`: Page URL to audit.
  - `--device`: Target device profile (desktop or mobile).
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
│   ├── config/               # Suite and persona YAML loader
│   ├── demo/                 # Built-in friction benchmark web app
│   ├── metrics/              # SUS scoring and Nielsen Norman mapping
│   ├── persona/              # Cognitive perception and decision brain
│   ├── report/               # HTML and Markdown report generator
│   └── store/                # SQLite database and query layer
├── web/                      # React 19 embedded UX studio
│   ├── src/
│   │   ├── components/       # Player, heatmap, and scorecards
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
