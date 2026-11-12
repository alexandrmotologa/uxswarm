# System architecture

uxswarm is an autonomous user research daemon that evaluates web applications before launch. It launches a swarm of synthetic user personas through headless Chrome DevTools Protocol (CDP) sessions, records user interactions, measures friction signals, and calculates mathematical System Usability Scale (SUS) scores.

## Architecture components

The system consists of five primary layers:

1. CLI and API daemon (`cmd/uxswarm`, `internal/api`):
   Provides command line workflows (`serve`, `test`, `report`, `audit`, `demo`) and an HTTP/SSE server. The server delivers JSON endpoints and real-time streaming updates during test runs.

2. Browser orchestration (`internal/browser`):
   Communicates directly with Chromium browsers (Google Chrome or Microsoft Edge) via Chrome DevTools Protocol over WebSockets. It manages headless viewport sizing, mobile device emulation, accessibility tree (AXTree) extraction, click/input dispatching, and frame-by-frame screenshot captures.

3. Cognitive persona engine (`internal/persona`):
   Coordinates simulated user agents defined in YAML. Each persona possesses a psychological profile, patience limit, technical literacy rating, and device preference. Decisions use either vision-language model inference or a deterministic heuristic perception engine that models human scanning patterns and reading velocity.

4. Usability telemetry and heuristic analyzer (`internal/metrics`):
   Monitors three core friction indicators: hesitation latency exceeding 3.5 seconds, dead clicks on non-interactive elements, and navigation backtrack loops. It calculates the Friction Index (FI) and computes 10-item System Usability Scale (SUS) scores. Issues map directly to the 10 Nielsen Norman usability heuristics with concrete copy and layout suggestions.

5. Embedded UX studio (`web`):
   A single-page React interface embedded into the Go binary. It provides a visual test dashboard, step-by-step session playback with persona inner thoughts, click coordinate heatmaps, and heuristic cards.

## Data flow

1. A test suite specification (`personas.yaml`) or CLI flag provides the target URL, user goal, and persona list.
2. The orchestrator initializes an isolated browser context for each persona in the swarm.
3. On every page state change, the browser layer extracts the accessibility tree, interactive bounding boxes, and viewport screenshot.
4. The persona brain evaluates visible UI elements against its goal and psychological biases, producing an action (click, type, navigate, wait) and internal reasoning monologue.
5. The actuator executes the action via CDP input dispatching while telemetry monitors hesitation, misclicks, and route loops.
6. When the persona reaches the goal, runs out of patience, or hits the session timeout, the session concludes.
7. The metrics engine calculates individual and aggregated SUS scores, identifies heuristic violations, and writes results to the SQLite store.
8. Reports are generated as JSON, standalone HTML, or viewed in the embedded web studio.
