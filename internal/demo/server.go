package demo

import (
	"fmt"
	"net"
	"net/http"
	"time"
)

// Server provides an embedded, self-contained SaaS onboarding web app with intentional UX friction points.
type Server struct {
	server   *http.Server
	listener net.Listener
	Port     int
}

// Start launches the demo server on the specified port. If port is 0, an ephemeral port is selected.
func Start(port int) (*Server, error) {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handleIndex)
	mux.HandleFunc("/signup", handleSignup)
	mux.HandleFunc("/workspace", handleWorkspace)
	mux.HandleFunc("/invite", handleInvite)
	mux.HandleFunc("/complete", handleComplete)
	mux.HandleFunc("/shop", handleShop)

	addr := fmt.Sprintf("127.0.0.1:%d", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("listen on %s: %w", addr, err)
	}

	actualPort := listener.Addr().(*net.TCPAddr).Port
	srv := &http.Server{
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		_ = srv.Serve(listener)
	}()

	return &Server{
		server:   srv,
		listener: listener,
		Port:     actualPort,
	}, nil
}

// URL returns the base URL of the running demo server.
func (s *Server) URL() string {
	return fmt.Sprintf("http://127.0.0.1:%d", s.Port)
}

// Close gracefully stops the demo HTTP server.
func (s *Server) Close() error {
	if s.server != nil {
		return s.server.Close()
	}
	return nil
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, "/signup", http.StatusFound)
}

const baseStyles = `
<style>
  :root {
    --primary: #2563eb;
    --primary-hover: #1d4ed8;
    --bg: #f8fafc;
    --card: #ffffff;
    --border: #e2e8f0;
    --text: #0f172a;
    --muted: #64748b;
  }
  * { box-sizing: border-box; margin: 0; padding: 0; font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif; }
  body { background: var(--bg); color: var(--text); display: flex; align-items: center; justify-content: center; min-height: 100vh; padding: 1.5rem; }
  .card { background: var(--card); border: 1px solid var(--border); border-radius: 12px; max-width: 480px; width: 100%; padding: 2.5rem; box-shadow: 0 4px 6px -1px rgba(0,0,0,0.05); }
  .step-pill { display: inline-block; background: #eff6ff; color: var(--primary); padding: 4px 10px; border-radius: 999px; font-size: 0.75rem; font-weight: 600; text-transform: uppercase; margin-bottom: 1rem; }
  h1 { font-size: 1.5rem; font-weight: 700; margin-bottom: 0.5rem; }
  p.subtitle { color: var(--muted); font-size: 0.9rem; margin-bottom: 1.5rem; line-height: 1.4; }
  .form-group { margin-bottom: 1.25rem; }
  label { display: block; font-size: 0.85rem; font-weight: 500; margin-bottom: 0.4rem; color: #334155; }
  input[type="text"], input[type="email"], input[type="password"], select {
    width: 100%; padding: 10px 12px; border: 1px solid var(--border); border-radius: 8px; font-size: 0.95rem; outline: none; transition: border-color 0.2s;
  }
  input:focus { border-color: var(--primary); ring: 2px solid #bfdbfe; }
  
  /* Intentional UX friction element: Deceptive promo card that looks like a button */
  .promo-banner {
    background: #f1f5f9; border: 1px dashed #94a3b8; border-radius: 8px; padding: 12px; margin-bottom: 1.25rem; cursor: pointer; text-align: center; font-size: 0.85rem; font-weight: 600; color: #1e293b;
  }
  .promo-banner:hover { background: #e2e8f0; }

  /* Intentional UX friction: Small tap target padding (below 44px on mobile) */
  .btn-submit {
    display: block; width: 100%; background: var(--primary); color: white; padding: 8px 16px; border: none; border-radius: 8px; font-size: 0.95rem; font-weight: 600; cursor: pointer; margin-top: 1rem; height: 38px;
  }
  .btn-submit:hover { background: var(--primary-hover); }
  .btn-secondary { background: transparent; color: var(--muted); border: 1px solid var(--border); padding: 8px 16px; border-radius: 8px; text-decoration: none; display: inline-block; font-size: 0.85rem; margin-top: 0.75rem; text-align: center; width: 100%; }

  .success-box { text-align: center; padding: 2rem 1rem; }
  .success-icon { font-size: 3rem; margin-bottom: 1rem; }
</style>
`

func handleSignup(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		http.Redirect(w, r, "/workspace", http.StatusSeeOther)
		return
	}

	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <title>Create your Acme Cloud Account</title>
  %s
</head>
<body>
  <div class="card">
    <span class="step-pill">Step 1 of 3</span>
    <h1>Create your account</h1>
    <p class="subtitle">Start your 14-day evaluation and connect your team.</p>
    
    <!-- Intentional Friction: Misleading banner card inviting dead clicks -->
    <div class="promo-banner" data-interactive="false" data-testid="promo-badge">
      Free 14-Day Enterprise Trial (Click to activate)
    </div>

    <form method="POST" action="/signup">
      <div class="form-group">
        <label for="name">Full Name</label>
        <input type="text" id="name" name="name" placeholder="Alex Morgan" required>
      </div>

      <div class="form-group">
        <label for="email">Work Email</label>
        <input type="email" id="email" name="email" placeholder="alex@company.com" required>
      </div>

      <div class="form-group">
        <label for="password">Password</label>
        <input type="password" id="password" name="password" placeholder="At least 8 characters" required>
      </div>

      <button type="submit" class="btn-submit" id="btn-signup">Continue to Workspace</button>
    </form>
  </div>
</body>
</html>`, baseStyles)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(html))
}

func handleWorkspace(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Simulate server delay causing user hesitation
		time.Sleep(300 * time.Millisecond)
		http.Redirect(w, r, "/invite", http.StatusSeeOther)
		return
	}

	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <title>Configure Workspace</title>
  %s
</head>
<body>
  <div class="card">
    <span class="step-pill">Step 2 of 3</span>
    <h1>Configure workspace</h1>
    <p class="subtitle">Set up your shared environment for projects and telemetry.</p>

    <form method="POST" action="/workspace">
      <div class="form-group">
        <!-- Ambiguous technical jargon label -->
        <label for="workspace">Enterprise Tenant Namespace (Workspace Name)</label>
        <input type="text" id="workspace" name="workspace" placeholder="Acme Labs" required>
      </div>

      <div class="form-group">
        <label for="tier">Workspace Tier</label>
        <select id="tier" name="tier">
          <option value="cloud">Cloud Native (Recommended)</option>
          <option value="hybrid">Hybrid Dedicated</option>
          <option value="enterprise">Air-gapped VPC</option>
        </select>
      </div>

      <button type="submit" class="btn-submit" id="btn-workspace">Save Workspace & Continue</button>
      <a href="/signup" class="btn-secondary">Back to Account</a>
    </form>
  </div>
</body>
</html>`, baseStyles)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(html))
}

func handleInvite(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		http.Redirect(w, r, "/complete", http.StatusSeeOther)
		return
	}

	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <title>Invite Teammates</title>
  %s
</head>
<body>
  <div class="card">
    <span class="step-pill">Step 3 of 3</span>
    <h1>Invite your team</h1>
    <p class="subtitle">Collaborate with colleagues inside your new workspace.</p>

    <form method="POST" action="/invite">
      <div class="form-group">
        <label for="colleague_email">Teammate Email</label>
        <input type="email" id="colleague_email" name="colleague_email" placeholder="colleague@example.com" required>
      </div>

      <!-- Ambiguous action label: Commit Sync instead of Complete Setup -->
      <button type="submit" class="btn-submit" id="btn-invite">Complete Setup</button>
    </form>
  </div>
</body>
</html>`, baseStyles)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(html))
}

func handleComplete(w http.ResponseWriter, r *http.Request) {
	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <title>Onboarding Complete</title>
  %s
</head>
<body>
  <div class="card">
    <div class="success-box">
      <div class="success-icon">&#x2714;</div>
      <h1>Setup Successful! Welcome to Acme</h1>
      <p class="subtitle">Your enterprise workspace is ready. You have successfully finished onboarding.</p>
      <a href="/signup" class="btn-submit" style="display:inline-block; text-decoration:none; line-height:22px;">Launch Dashboard</a>
    </div>
  </div>
</body>
</html>`, baseStyles)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(html))
}

func handleShop(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		http.Redirect(w, r, "/complete", http.StatusSeeOther)
		return
	}

	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <title>E-Commerce Store Checkout</title>
  %s
</head>
<body>
  <div class="card">
    <span class="step-pill">Express Checkout</span>
    <h1>Order Summary</h1>
    <p class="subtitle">UXSwarm Pro Subscription (12 months).</p>

    <form method="POST" action="/shop">
      <div class="form-group">
        <label for="name">Cardholder Name</label>
        <input type="text" id="name" name="name" placeholder="Arthur Pendelton" required>
      </div>

      <div class="form-group">
        <label for="card">Card Number</label>
        <input type="text" id="card" name="card" placeholder="4242 &bull;&bull;&bull;&bull; &bull;&bull;&bull;&bull; 4242" required>
      </div>

      <button type="submit" class="btn-submit" id="btn-order">Place Order</button>
    </form>
  </div>
</body>
</html>`, baseStyles)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(html))
}
