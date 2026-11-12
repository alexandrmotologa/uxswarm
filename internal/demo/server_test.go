package demo

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestDemoServer_Lifecycle(t *testing.T) {
	srv, err := Start(0) // Port 0 uses ephemeral port
	if err != nil {
		t.Fatalf("failed to start demo server: %v", err)
	}
	defer srv.Close()

	if srv.Port <= 0 {
		t.Fatalf("expected positive port number, got %d", srv.Port)
	}

	url := srv.URL() + "/signup"
	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("get /signup error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(body), "Create your account") {
		t.Errorf("expected signup content in response")
	}
	if !strings.Contains(string(body), "promo-badge") {
		t.Errorf("expected promo-badge in signup response")
	}
}
