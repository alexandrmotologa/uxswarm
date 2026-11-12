package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"strings"
	"time"

	"github.com/alexandrmotologa/uxswarm/internal/compare"
	"github.com/alexandrmotologa/uxswarm/internal/config"
	"github.com/alexandrmotologa/uxswarm/internal/gaze"
	"github.com/alexandrmotologa/uxswarm/internal/metrics"
	"github.com/alexandrmotologa/uxswarm/internal/patch"
	"github.com/alexandrmotologa/uxswarm/internal/report"
	"github.com/alexandrmotologa/uxswarm/internal/runner"
	"github.com/alexandrmotologa/uxswarm/internal/store"
)

// Server coordinates HTTP API routing, SSE streaming, and embedded studio serving.
type Server struct {
	store    *store.Store
	hub      *SSEHub
	staticFS fs.FS
	mux      *http.ServeMux
}

// NewServer constructs an API server with SQLite persistence and event streaming.
func NewServer(s *store.Store, staticFS fs.FS) *Server {
	srv := &Server{
		store:    s,
		hub:      NewSSEHub(),
		staticFS: staticFS,
		mux:      http.NewServeMux(),
	}
	srv.registerRoutes()
	return srv
}

// Hub returns the SSE event hub for broadcasting execution events.
func (s *Server) Hub() *SSEHub {
	return s.hub
}

// ServeHTTP delegates to the internal mux.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Enable CORS for development
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	s.mux.ServeHTTP(w, r)
}

func (s *Server) registerRoutes() {
	s.mux.HandleFunc("/api/health", s.handleHealth)
	s.mux.HandleFunc("/api/personas", s.handlePersonas)
	s.mux.HandleFunc("/api/runs", s.handleRuns)
	s.mux.HandleFunc("/api/runs/", s.handleRunRoutes)
	s.mux.HandleFunc("/api/sessions/", s.handleSessionRoutes)
	s.mux.HandleFunc("/api/compare", s.handleCompare)

	// Static web UI
	if s.staticFS != nil {
		fileServer := http.FileServer(http.FS(s.staticFS))
		s.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/api/") {
				http.NotFound(w, r)
				return
			}
			fileServer.ServeHTTP(w, r)
		})
	} else {
		s.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/api/") {
				http.NotFound(w, r)
				return
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			_, _ = w.Write([]byte(`<!DOCTYPE html><html><body><h1>uxswarm API Daemon</h1><p>Running API mode.</p></body></html>`))
		})
	}
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"status":    "healthy",
		"timestamp": time.Now(),
		"service":   "uxswarm",
		"version":   "1.0.0",
	})
}

func (s *Server) handlePersonas(w http.ResponseWriter, r *http.Request) {
	personas := config.GetDefaultPersonas()
	writeJSON(w, http.StatusOK, personas)
}

func (s *Server) handleRuns(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		runs, err := s.store.ListRuns(50)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, runs)

	case http.MethodPost:
		var req struct {
			Name      string `json:"name"`
			TargetURL string `json:"target_url"`
			Goal      string `json:"goal"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		if req.TargetURL == "" {
			http.Error(w, "target_url is required", http.StatusBadRequest)
			return
		}
		if req.Name == "" {
			req.Name = "Usability Study: " + req.TargetURL
		}
		if req.Goal == "" {
			req.Goal = "Navigate onboarding flow and complete goal"
		}

		runID := fmt.Sprintf("run_%d", time.Now().UnixNano()/1e6)
		run := &store.RunRecord{
			ID:                runID,
			Name:              req.Name,
			TargetURL:         req.TargetURL,
			Goal:              req.Goal,
			Status:            "running",
			CreatedAt:         time.Now(),
			MandatoryRedesign: false,
		}
		if err := s.store.CreateRun(run); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Asynchronously execute swarm in background and stream live events via SSE
		go func() {
			opts := runner.RunOptions{
				Name:      req.Name,
				TargetURL: req.TargetURL,
				Goal:      req.Goal,
				Store:     s.store,
				Headless:  true,
				OnStep: func(p config.Persona, evt metrics.StepEvent) {
					s.hub.Broadcast(runID, "step", map[string]any{
						"persona": p.Name,
						"event":   evt,
					})
				},
				OnSessionEnd: func(p config.Persona, summary *metrics.TelemetrySummary, sus *metrics.SUSResult) {
					s.hub.Broadcast(runID, "session_end", map[string]any{
						"persona": p.Name,
						"summary": summary,
						"sus":     sus,
					})
				},
			}
			rep, err := runner.ExecuteSwarm(context.Background(), opts)
			if err == nil && rep != nil {
				s.hub.Broadcast(runID, "complete", rep.Run)
			}
		}()

		writeJSON(w, http.StatusCreated, run)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleRunRoutes(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/runs/")
	parts := strings.Split(path, "/")
	if len(parts) == 0 || parts[0] == "" {
		http.NotFound(w, r)
		return
	}

	runID := parts[0]

	if len(parts) > 1 {
		switch parts[1] {
		case "stream":
			s.hub.ServeHTTP(w, r, runID)
			return
		case "sessions":
			sessions, err := s.store.ListSessionsByRun(runID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			writeJSON(w, http.StatusOK, sessions)
			return
		case "heuristics":
			findings, err := s.store.ListHeuristicsByRun(runID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			writeJSON(w, http.StatusOK, findings)
			return
		case "patches":
			findings, _ := s.store.ListHeuristicsByRun(runID)
			sessions, _ := s.store.ListSessionsByRun(runID)
			steps := make([]metrics.StepEvent, 0)
			for _, sess := range sessions {
				sessSteps, _ := s.store.ListStepsBySession(sess.ID)
				for _, st := range sessSteps {
					steps = append(steps, metrics.StepEvent{
						ActionType:        st.ActionType,
						TargetSelector:    st.TargetSelector,
						TargetDescription: st.TargetDescription,
						IsDeadClick:       st.IsDeadClick,
					})
				}
			}
			patches := patch.GeneratePatches(findings, steps)
			writeJSON(w, http.StatusOK, patches)
			return
		case "gaze":
			// Generate simulated gaze heatmap for active run
			viewport := config.Viewport{Width: 1440, Height: 900, IsMobile: false}
			gazeMap := gaze.GenerateGazeMap(nil, viewport, "Onboarding Flow")
			writeJSON(w, http.StatusOK, gazeMap)
			return
		default:
			http.NotFound(w, r)
			return
		}
	}

	// Single run GET
	if r.Method == http.MethodGet {
		run, err := s.store.GetRun(runID)
		if err != nil {
			http.Error(w, "run not found", http.StatusNotFound)
			return
		}
		writeJSON(w, http.StatusOK, run)
		return
	}

	http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
}

func (s *Server) handleSessionRoutes(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/sessions/")
	parts := strings.Split(path, "/")
	if len(parts) == 0 || parts[0] == "" {
		http.NotFound(w, r)
		return
	}

	sessionID := parts[0]

	if len(parts) > 1 && parts[1] == "steps" {
		steps, err := s.store.ListStepsBySession(sessionID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, steps)
		return
	}

	http.Error(w, "not found", http.StatusNotFound)
}

func (s *Server) handleCompare(w http.ResponseWriter, r *http.Request) {
	runIDA := r.URL.Query().Get("run_a")
	runIDB := r.URL.Query().Get("run_b")

	if runIDA == "" || runIDB == "" {
		http.Error(w, "parameters run_a and run_b are required", http.StatusBadRequest)
		return
	}

	runA, errA := s.store.GetRun(runIDA)
	runB, errB := s.store.GetRun(runIDB)
	if errA != nil || errB != nil {
		http.Error(w, "one or both runs not found", http.StatusNotFound)
		return
	}

	sessA, _ := s.store.ListSessionsByRun(runIDA)
	sessB, _ := s.store.ListSessionsByRun(runIDB)
	heurA, _ := s.store.ListHeuristicsByRun(runIDA)
	heurB, _ := s.store.ListHeuristicsByRun(runIDB)

	repA := &report.ReportData{Run: *runA, Sessions: sessA, Heuristics: heurA}
	repB := &report.ReportData{Run: *runB, Sessions: sessB, Heuristics: heurB}

	cmp := compare.CompareVariants(repA, repB)
	writeJSON(w, http.StatusOK, cmp)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
