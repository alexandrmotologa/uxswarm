# `uxswarm` — Architecture & Engineering Master Plan
> **Autonomous Synthetic User Persona Swarm & Pre-Release UX Friction Profiler**  
> *Target: AI Product Engineering, Automated Usability Research, Cognitive Load Profiling & SUS Scoring*

---

## 1. Executive Summary & Market Problem

### 1.1 The Usability Testing Chokepoint in Modern Product Teams
Modern software product teams iterate rapidly, yet **user research and usability testing remain painfully slow and expensive**:
- Recruiting 5 to 10 real human participants, scheduling interviews, and moderating sessions takes **3 to 4 weeks** and costs **\$4,000 to \$10,000+** per sprint.
- As a result, product teams routinely skip user testing, releasing untested onboarding flows, redesigns, and checkout funnels directly to production.
- **The Cost of Blind Releases**: Over **40% of signups drop off** due to avoidable UX friction: confusing button copy, hidden form fields, lack of trust signals, or visual hierarchy failures.
- Traditional product analytics tools (Mixpanel, PostHog) only show *where* users dropped off after they left—they cannot explain *why* users felt confused or hesitated before leaving.

### 1.2 The `uxswarm` Solution
`uxswarm` is an autonomous, self-hosted Go and Playwright/CDP daemon with an embedded React 19 UX studio. It simulates realistic human user testing before a feature is released to real users:
- **Diverse Synthetic Persona Swarm (`personas.yaml`)**: Spins up a swarm of distinct cognitive personas (e.g., *Rushed Non-Technical Executive*, *Skeptical Security Architect*, *Impatient Mobile User*) with explicit psychological profiles, technical literacy levels, and domain biases.
- **Autonomous Headless Navigation**: Personas browse the web app (staging or localhost), interpret the UI visual hierarchy, and attempt to accomplish a business goal (e.g., *"Sign up, create a project, and invite a colleague"*).
- **Cognitive Friction Telemetry**: Tracks hesitation intervals, misclick rates (clicking non-interactive text), backtrack loops, and perceptual ambiguity.
- **Automated System Usability Scale (SUS) Score**: Generates mathematical usability scores (0–100) and actionable design recommendations categorized by Nielsen Norman heuristics.
- **Interactive Session Replay Studio**: Watch visual step-by-step replays of each persona's journey with overlaid "inner thoughts" and friction badges.

---

## 2. High-Level System Architecture

```
                       +---------------------------------------------+
                       |           Product Goal & Web URL            |
                       |  "Sign up for free trial & invite member"   |
                       +---------------------------------------------+
                                              |
                                              v
                       +---------------------------------------------+
                       |             uxswarm Orchestrator            |
                       |       (Single Go Binary, <30MB RAM)         |
                       +---------------------------------------------+
                                              |
                                              v
                       +---------------------------------------------+
                       |          Synthetic Persona Swarm            |
                       |  - Persona A: Non-tech CFO (low patience)   |
                       |  - Persona B: DevSecOps Lead (high scrutiny)|
                       |  - Persona C: Mobile User on 3G network     |
                       +---------------------------------------------+
                                              |
                         +--------------------+--------------------+
                         |                                         |
                         v                                         v
         +-------------------------------+         +-------------------------------+
         |     Headless CDP Browser      |         |    Vision-Language Perceiver  |
         |  - DOM Tree & Accessibility   |         |  - Screenshot Analysis        |
         |  - Real-Time Action Actuator  |         |  - Visual Hierarchy Scoring   |
         |  - Session Recording (Frames) |         |  - Cognitive Load Evaluator   |
         +-------------------------------+         +-------------------------------+
                         |                                         |
                         +--------------------+--------------------+
                                              |
                                              v
                       +---------------------------------------------+
                       |         Cognitive Friction Telemetry        |
                       |  - Hesitation Latency (> 3.5s pause)        |
                       |  - Dead Clicks / Misclick Coordinates       |
                       |  - Backtrack Loops (Route A -> B -> A)      |
                       |  - Drop-off Frustration Threshold Breaches  |
                       +---------------------------------------------+
                                              |
                                              v
                       +---------------------------------------------+
                       |        UX Scoring & Heuristic Engine        |
                       |  - Computes System Usability Scale (SUS)    |
                       |  - Maps Issues to Nielsen Norman 10 Rules   |
                       |  - Synthesizes Concrete Copy/Design Fixes   |
                       +---------------------------------------------+
                                              |
                         +--------------------+--------------------+
                         |                                         |
                         v                                         v
         +-------------------------------+         +-------------------------------+
         |       Executive UX Report     |         |     React 19 Replay Studio    |
         |  - SUS Benchmark Score (0-100)|         |  - Step-by-Step Video Replay  |
         |  - Funnel Conversion Dropouts |         |  - Persona "Inner Monologue"  |
         |  - Actionable Redesign Specs  |         |  - Interactive Friction Heat  |
         +-------------------------------+         +-------------------------------+
```

---

## 3. Mathematical Usability Metrics

### 3.1 System Usability Scale (SUS) Formulation
For each persona $p$ evaluating the 10 standard SUS questions with response values $R_{i} \in [1, 5]$:
- For odd items ($1, 3, 5, 7, 9$): $S_i = R_i - 1$
- For even items ($2, 4, 6, 8, 10$): $S_i = 5 - R_i$

$$\text{SUS Score} = \left( \sum_{i=1}^{10} S_i \right) \times 2.5$$
- $\text{Score} \ge 80.3$: **Grade A (Excellent UX)**
- $68.0 \le \text{Score} < 80.3$: **Grade B / C (Acceptable)**
- $\text{Score} < 68.0$: **Grade D / F (Critical Usability Failures)**

### 3.2 Friction Index ($FI$)
$$\text{FI} = w_1 \cdot \frac{T_{\text{hesitation}}}{T_{\text{total}}} + w_2 \cdot N_{\text{misclicks}} + w_3 \cdot N_{\text{backtracks}}$$
Where $w_1 = 0.35, w_2 = 0.40, w_3 = 0.25$. Any flow with $\text{FI} > 2.5$ triggers a mandatory redesign recommendation.

---

## 4. Declarative Persona Specification (`personas.yaml`)

```yaml
version: "1"
suite:
  target_url: "https://staging.example.com/signup"
  goal: "Complete account registration, set workspace name to 'Acme Marketing', and invite teammate 'colleague@example.com'."
  max_session_duration: "3m"

personas:
  - id: "impatient-executive"
    name: "Sarah Chen - CFO"
    age: 48
    archetype: "HIGH_IMPATIENCE_LOW_TECH"
    patience_threshold_seconds: 4.0
    technical_literacy: 3 # Scale 1-10
    attention_span: "LOW"
    traits:
      - "Hates reading long paragraphs of explanatory text"
      - "Scans solely for primary CTA buttons"
      - "Drops off immediately if credit card is required upfront"
      - "Expects social login (Google/Apple) over complex passwords"

  - id: "skeptical-security-lead"
    name: "Marcus Vance - DevSecOps Architect"
    age: 36
    archetype: "SCRUTINIZER_HIGH_TECH"
    patience_threshold_seconds: 15.0
    technical_literacy: 9
    attention_span: "HIGH"
    traits:
      - "Inspects terms of service and data residency disclaimers"
      - "Notices missing password strength requirements"
      - "Wants explicit SSO or 2FA setup options"

  - id: "distracted-mobile-user"
    name: "Elena Rostova - Field Sales Rep"
    age: 29
    archetype: "MOBILE_THUMB_DRIVEN"
    viewport:
      width: 390
      height: 844
      is_mobile: true
    traits:
      - "Uses one-handed thumb navigation"
      - "Struggles with small tap targets (< 44px)"
      - "Easily annoyed by modal popups blocking view"
```

---

## 5. Repository & Directory Structure

```
uxswarm/
├── .github/
│   └── workflows/
│       ├── test.yml
│       └── release.yml
├── cmd/
│   └── uxswarm/
│       └── main.go                 # CLI entrypoint (test, report, serve)
├── internal/
│   ├── config/                     # Persona & Suite YAML parser
│   ├── browser/                    # Headless Browser & CDP Engine
│   │   ├── session.go              # Playwright browser instance coordinator
│   │   ├── tree.go                 # Accessibility DOM tree extractor
│   │   ├── recorder.go             # Frame-by-frame screenshot recorder
│   │   └── actuator.go             # Mouse move, click, type, and scroll simulator
│   ├── persona/                    # Cognitive Agent Brain
│   │   ├── agent.go                # Vision-Language reasoning loop
│   │   ├── vision_perceiver.go     # Analyzes screenshot visual hierarchy
│   │   ├── inner_speech.go         # Synthesizes persona's thought stream
│   │   └── drop_off.go             # Evaluates frustration threshold breaches
│   ├── metrics/                    # Usability Analytics
│   │   ├── telemetry.go            # Tracks hesitation, dead clicks, and loops
│   │   ├── sus_calculator.go       # System Usability Scale math
│   │   └── nielsen_heuristics.go   # Maps errors to NN 10 usability heuristics
│   ├── store/                      # SQLite Persistence Layer
│   │   ├── db.go
│   │   └── sessions.go             # Historical test runs and replay frames
│   └── api/                        # HTTP REST & SSE Streaming
│       ├── handler.go
│       └── routes.go
├── web/                            # Embedded React 19 UX Studio
│   ├── index.html
│   ├── package.json
│   ├── vite.config.ts
│   └── src/
│       ├── App.tsx
│       ├── components/
│       │   ├── SessionReplay.tsx   # Video playback with inner-speech subtitles
│       │   ├── FrictionHeatmap.tsx # Click coordinates & hesitation hotspots
│       │   ├── SUSRadarChart.tsx   # 10-dimensional usability radar breakdown
│       │   └── HeuristicCard.tsx   # Actionable design recommendation cards
│       └── pages/
│           ├── Studio.tsx
│           ├── RunOverview.tsx
│           └── PersonaBuilder.tsx
├── examples/
│   ├── personas_b2b_saas.yaml
│   └── personas_ecommerce_checkout.yaml
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

---

## 6. CLI Command Surface

```bash
# Start the web UX studio and API server
uxswarm serve --port 8080 --db ./uxswarm.db

# Run an automated usability study against a staging URL
uxswarm test \
  --config ./personas.yaml \
  --url https://staging.myapp.com \
  --goal "Complete onboarding funnel" \
  --output ./ux_report.json

# Generate an executive HTML usability report
uxswarm report --run-id run_98a72b --format html --output ./usability_report.html

# Run quick heuristic check on a single web page
uxswarm audit --url https://staging.myapp.com/pricing --device mobile
```

---

## 7. Implementation Roadmap (5 Phases)

### Phase 1: Headless CDP Driver & Accessibility Extractor
- Implement Playwright/CDP browser session manager in Go.
- Extract accessibility DOM nodes and calculate interactive bounding boxes.
- Implement frame-by-frame screenshot recorder for video generation.

### Phase 2: Cognitive Persona Agent & Vision Perceiver
- Build persona YAML loader with patience thresholds and psychological archetypes.
- Implement Vision-Language decision loop (GPT-4o / Claude 3.5 Sonnet) choosing next DOM action.
- Capture persona "inner speech" reflecting confusion, hesitation, or satisfaction.

### Phase 3: Friction Telemetry & Heuristic Engine
- Implement mathematical detectors for hesitation intervals ($> 3.5$s), misclicks, and backtrack loops.
- Implement System Usability Scale (SUS) survey synthesizer based on session events.
- Map friction patterns to Nielsen Norman usability heuristics.

### Phase 4: Embedded React 19 Replay Studio
- Build interactive video player synchronizing screen frames with persona inner monologue.
- Display visual click heatmaps and cognitive load timeline graphs.
- Package frontend into single Go binary using `//go:embed`.

### Phase 5: Verification, Packaging & Documentation
- Test against complex real-world onboarding flows (Stripe, Linear, Supabase).
- Verify single run with 5 personas completes in under 2 minutes.
- Publish documentation, animated demo GIFs, and persona catalogs.
