<p align="center">
  <img src="docs/images/logo.png?raw=true" alt="uxswarm Logo" width="130" style="border-radius: 24px;" />
</p>

<h1 align="center">uxswarm</h1>

<p align="center">
  <a href="https://github.com/alexandrmotologa/uxswarm/actions/workflows/test.yml"><img src="https://github.com/alexandrmotologa/uxswarm/actions/workflows/test.yml/badge.svg" alt="CI Status" /></a>
  <img src="https://img.shields.io/badge/go-1.22+-00ADD8.svg?style=flat&logo=go" alt="Go Version" />
  <img src="https://img.shields.io/badge/react-19-61DAFB.svg?style=flat&logo=react" alt="React 19" />
  <img src="https://img.shields.io/badge/engine-Chrome_CDP-4285F4.svg?style=flat&logo=googlechrome" alt="Chrome DevTools Protocol" />
  <img src="https://img.shields.io/badge/storage-SQLite_WAL-003B57.svg?style=flat&logo=sqlite" alt="SQLite WAL" />
  <img src="https://img.shields.io/badge/license-MIT-green.svg?style=flat" alt="MIT License" />
</p>

<p align="center">
  <strong>Autonomous synthetic user persona swarm and pre-release UX friction profiler.</strong><br />
  Simulates diverse behavioral archetypes through headless Chrome automation, calculates mathematical System Usability Scale (SUS) scores, and generates targeted engineering code patches.
</p>

<p align="center">
  <img src="docs/images/uxswarm_demo.gif?raw=true" alt="uxswarm Interactive Studio Demo" width="850" style="border-radius: 12px; box-shadow: 0 8px 30px rgba(0,0,0,0.3);" />
</p>

---

## Executive Overview

uxswarm deploys coordinated clusters of simulated human personas across web applications to audit onboarding funnels, checkout flows, and complex form interfaces prior to production release. 

Personas navigate pages via headless Chrome DevTools Protocol (CDP), extracting the accessibility tree and interactive coordinates. During navigation, the engine captures micro-behaviors including hesitation latencies exceeding 3.5 seconds, dead clicks on non-interactive visual decoys, and navigation backtracks. 

Collected telemetry is compiled into standardized System Usability Scale (SUS) benchmark grades, mapped against the Nielsen Norman 10 Usability Heuristics, and converted into concrete CSS, Tailwind, and React code patches.

---

## Visual Studio & Live Telemetry

### 1. Executive Usability Dashboard & SUS Benchmark
Aggregates swarm session results into an overarching System Usability Scale score, Friction Index (FI), funnel completion percentage, and per-persona outcome statuses.

<p align="center">
  <img src="docs/images/dashboard_overview.png?raw=true" alt="uxswarm Executive Dashboard" width="850" style="border-radius: 10px; border: 1px solid #1e293b;" />
</p>

### 2. Session Replay & Cognitive Thought Stream
Scrub through step-by-step persona execution, viewing exact mouse click coordinates, form inputs, frustration telemetry, and synthesized inner monologue subtitles.

<p align="center">
  <img src="docs/images/session_replay.png?raw=true" alt="uxswarm Session Replay" width="850" style="border-radius: 10px; border: 1px solid #1e293b;" />
</p>

### 3. Friction Coordinates & Simulated Gaze Heatmaps
Visualizes physical interaction clusters alongside simulated human eye-tracking visual attention models, identifying F-pattern reading drop-offs and dead-click traps.

<p align="center">
  <img src="docs/images/friction_heatmap.png?raw=true" alt="uxswarm Friction Heatmap" width="850" style="border-radius: 10px; border: 1px solid #1e293b;" />
</p>

### 4. Nielsen Norman Heuristic Violations & Code Patches
Every detected usability flaw maps to an established heuristic with severity ratings, affected personas, and actionable engineering code patches ready to copy into your codebase.

<p align="center">
  <img src="docs/images/heuristic_patches.png?raw=true" alt="uxswarm Heuristic Audit and Code Patches" width="850" style="border-radius: 10px; border: 1px solid #1e293b;" />
</p>

### 5. A/B Usability Differential Matrix
Quantifies usability improvements across baseline (Variant A) and redesigned (Variant B) candidate versions, displaying delta SUS score, delta friction, and abandonment reductions.

<p align="center">
  <img src="docs/images/ab_compare.png?raw=true" alt="uxswarm A/B Usability Differential Suite" width="850" style="border-radius: 10px; border: 1px solid #1e293b;" />
</p>

### 6. Responsive Viewport Profiling
Evaluates layouts across desktop viewports (1440x900) and mobile screens (390x844), checking one-handed thumb reach comfort and minimum 44px tap target requirements.

<p align="center">
  <img src="docs/images/mobile_responsive.png?raw=true" alt="uxswarm Mobile Responsive Profiling" width="360" style="border-radius: 10px; border: 1px solid #1e293b;" />
</p>

---

## Core Capabilities

- **Diverse Synthetic Personas**: Configurable via YAML with psychological archetypes including rushed CFOs, skeptical security leads, mobile thumb navigators, and accessibility-reliant users.
- **Pure Go CDP Engine**: Direct browser automation using Chrome DevTools Protocol for accessibility tree extraction, coordinate calculation, and screenshot capture.
- **Mathematical SUS Scoring**: Standardized 10-item System Usability Scale calculation generating 0 to 100 scores and letter grades (A to F).
- **A/B Differential Engine**: Direct comparison of baseline and candidate variants to measure delta SUS points, friction index changes, and drop-off reductions.
- **Automated Code Patching**: Generates production-ready CSS, Tailwind, and React code fixes to eliminate dead clicks and reduce submission hesitations.
- **Gaze Simulation**: Computational model mapping F-pattern scanning and visual saliency fixations across landing and onboarding screens.
- **CI/CD Usability Budgets**: Enforces minimum SUS scores and maximum friction thresholds directly in CI pipelines with automated GitHub Step Summary markdown output.
- **Embedded Web Studio**: Zero-dependency Go binary serving the React 19 interface via `//go:embed`.
- **Local Test Application**: Includes a built-in multi-step SaaS onboarding app with intentional friction points for testing without external services.

---

## Installation

### Prerequisites

- Go 1.22 or higher
- Google Chrome or Microsoft Edge installed
- Node.js 20 or higher (required only when modifying the web studio frontend)

### Build from Source

```bash
git clone https://github.com/alexandrmotologa/uxswarm.git
cd uxswarm
go build -o bin/uxswarm.exe ./cmd/uxswarm
```

---

## Quick Start

### 1. Launch the Built-in Demo Application

Start the local SaaS onboarding application:

```bash
./bin/uxswarm demo --port 3000
```

### 2. Execute a Usability Study

Run a 5-persona swarm against the onboarding flow:

```bash
./bin/uxswarm test \
  --config examples/personas_b2b_saas.yaml \
  --url http://localhost:3000 \
  --goal "Complete account registration and workspace onboarding" \
  --output ux_report.json
```

### 3. Compare Two Design Variants (A/B Usability Diff)

Evaluate candidate improvements against the baseline:

```bash
./bin/uxswarm compare \
  --url-a http://localhost:3000/v1 \
  --url-b http://localhost:3000/v2 \
  --output ab_comparison_report.html
```

### 4. Generate Usability Code Patches

Produce actionable CSS, Tailwind, and React fixes for detected issues:

```bash
./bin/uxswarm patch --input ux_report.json
```

### 5. Enforce Usability Budgets in CI/CD

Fail the build if the application does not meet usability thresholds:

```bash
./bin/uxswarm test \
  --url https://staging.myapp.com/signup \
  --min-sus 70.0 \
  --max-friction 1.5 \
  --allow-abandonment=false \
  --github-summary
```

### 6. Launch the Interactive UX Studio

Start the embedded studio to inspect session replays, gaze maps, and heuristics:

```bash
./bin/uxswarm serve --port 8080 --db uxswarm.db
```

Navigate to `http://localhost:8080` in your browser.

---

## Command Reference

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

---

## Repository Structure

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

---

## License

MIT License. See [LICENSE](LICENSE) for details.
