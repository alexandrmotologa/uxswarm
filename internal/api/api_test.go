package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"github.com/alexandrmotologa/uxswarm/internal/store"
)

func TestAPIServer_Endpoints(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "api_test.db")

	s, err := store.NewStore(dbPath)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	defer s.Close()

	srv := NewServer(s, nil)

	// 1. Test /api/health
	reqHealth := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	wHealth := httptest.NewRecorder()
	srv.ServeHTTP(wHealth, reqHealth)

	if wHealth.Code != http.StatusOK {
		t.Errorf("health endpoint returned %d", wHealth.Code)
	}

	// 2. Test /api/personas
	reqPersonas := httptest.NewRequest(http.MethodGet, "/api/personas", nil)
	wPersonas := httptest.NewRecorder()
	srv.ServeHTTP(wPersonas, reqPersonas)

	if wPersonas.Code != http.StatusOK {
		t.Errorf("personas endpoint returned %d", wPersonas.Code)
	}

	// 3. Test /api/runs (empty list)
	reqRuns := httptest.NewRequest(http.MethodGet, "/api/runs", nil)
	wRuns := httptest.NewRecorder()
	srv.ServeHTTP(wRuns, reqRuns)

	if wRuns.Code != http.StatusOK {
		t.Errorf("runs endpoint returned %d", wRuns.Code)
	}
	var runs []store.RunRecord
	if err := json.Unmarshal(wRuns.Body.Bytes(), &runs); err != nil {
		t.Fatalf("unmarshal runs: %v", err)
	}
	if len(runs) != 0 {
		t.Errorf("expected 0 runs initially, got %d", len(runs))
	}

	// Insert test run in store
	testRun := &store.RunRecord{
		ID:        "run_test_api",
		Name:      "API Test Run",
		TargetURL: "http://localhost:3000",
		Goal:      "Test goal",
		Status:    "completed",
		CreatedAt: time.Now(),
	}
	_ = s.CreateRun(testRun)

	// 4. Test /api/runs/run_test_api
	reqRun := httptest.NewRequest(http.MethodGet, "/api/runs/run_test_api", nil)
	wRun := httptest.NewRecorder()
	srv.ServeHTTP(wRun, reqRun)

	if wRun.Code != http.StatusOK {
		t.Errorf("get run returned %d", wRun.Code)
	}
}
