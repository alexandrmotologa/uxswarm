package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/alexandrmotologa/uxswarm/internal/api"
	"github.com/alexandrmotologa/uxswarm/internal/config"
	"github.com/alexandrmotologa/uxswarm/internal/demo"
	"github.com/alexandrmotologa/uxswarm/internal/metrics"
	"github.com/alexandrmotologa/uxswarm/internal/report"
	"github.com/alexandrmotologa/uxswarm/internal/runner"
	"github.com/alexandrmotologa/uxswarm/internal/store"
	"github.com/alexandrmotologa/uxswarm/web"
)

const version = "1.0.0"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "serve":
		cmdServe(os.Args[2:])
	case "test":
		cmdTest(os.Args[2:])
	case "report":
		cmdReport(os.Args[2:])
	case "audit":
		cmdAudit(os.Args[2:])
	case "demo":
		cmdDemo(os.Args[2:])
	case "version", "--version", "-v":
		fmt.Printf("uxswarm v%s\n", version)
	case "help", "--help", "-h":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`uxswarm - Autonomous Synthetic User Persona Swarm & Pre-Release UX Friction Profiler

Usage:
  uxswarm <command> [arguments]

Commands:
  serve     Launch the HTTP API daemon and embedded React 19 UX studio
  test      Execute an automated usability study with a persona swarm
  report    Generate executive HTML, Markdown, or JSON usability reports
  audit     Perform a fast single-page heuristic accessibility audit
  demo      Run the built-in SaaS onboarding test target web application
  version   Print version and build information

Run 'uxswarm <command> --help' for details on each command.`)
}

func cmdServe(args []string) {
	fs := flag.NewFlagSet("serve", flag.ExitOnError)
	port := fs.Int("port", 8080, "HTTP listen port")
	dbPath := fs.String("db", "uxswarm.db", "SQLite database path")
	_ = fs.Parse(args)

	db, err := store.NewStore(*dbPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	staticFS, err := web.DistFS()
	if err != nil {
		fmt.Printf("Warning: embedded studio files unavailable: %v. Running API mode.\n", err)
	}

	srv := api.NewServer(db, staticFS)
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", *port),
		Handler: srv,
	}

	fmt.Printf("uxswarm Studio & API daemon active on http://localhost:%d\n", *port)
	fmt.Printf("Database: %s\n", *dbPath)
	fmt.Println("Press Ctrl+C to terminate.")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		}
	}()

	<-sigChan
	fmt.Println("\nShutting down uxswarm server...")
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_ = httpServer.Shutdown(ctx)
}

func cmdTest(args []string) {
	fs := flag.NewFlagSet("test", flag.ExitOnError)
	cfgPath := fs.String("config", "", "Path to personas.yaml suite")
	targetURL := fs.String("url", "", "Target web application URL")
	goal := fs.String("goal", "Complete account registration and workspace onboarding", "Task goal description")
	outputPath := fs.String("output", "ux_report.json", "JSON report output destination")
	htmlReport := fs.String("report", "", "Optional standalone HTML report output path")
	dbPath := fs.String("db", "uxswarm.db", "SQLite database path")
	headless := fs.Bool("headless", true, "Run Chromium in headless mode")
	isDemo := fs.Bool("demo", false, "Start built-in demo server automatically")
	_ = fs.Parse(args)

	var demoSrv *demo.Server
	if *isDemo || *targetURL == "" {
		srv, err := demo.Start(0)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to start demo test target: %v\n", err)
			os.Exit(1)
		}
		demoSrv = srv
		defer demoSrv.Close()
		*targetURL = srv.URL() + "/signup"
		fmt.Printf("Started built-in demo onboarding app on %s\n", *targetURL)
	}

	var personas []config.Persona
	var suiteName string

	if *cfgPath != "" {
		s, err := config.LoadSuite(*cfgPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading suite %s: %v\n", *cfgPath, err)
			os.Exit(1)
		}
		personas = s.Personas
		suiteName = s.Suite.Name
		if s.Suite.Goal != "" && *goal == "Complete account registration and workspace onboarding" {
			*goal = s.Suite.Goal
		}
		if s.Suite.TargetURL != "" && !*isDemo {
			*targetURL = s.Suite.TargetURL
		}
	} else {
		personas = config.GetDefaultPersonas()
		suiteName = "Automated Usability Study"
	}

	var db *store.Store
	if *dbPath != "" {
		s, err := store.NewStore(*dbPath)
		if err == nil {
			db = s
			defer db.Close()
		}
	}

	fmt.Printf("\n=== uxswarm Usability Study ===\n")
	fmt.Printf("Target URL : %s\n", *targetURL)
	fmt.Printf("User Goal  : %s\n", *goal)
	fmt.Printf("Swarm Size : %d synthetic personas\n\n", len(personas))

	opts := runner.RunOptions{
		Name:       suiteName,
		TargetURL:  *targetURL,
		Goal:       *goal,
		Personas:   personas,
		Headless:   *headless,
		Store:      db,
		OnStep: func(p config.Persona, evt metrics.StepEvent) {
			pauseNotice := ""
			if evt.HesitationSeconds >= 3.5 {
				pauseNotice = fmt.Sprintf(" [Paused %.1fs]", evt.HesitationSeconds)
			}
			deadNotice := ""
			if evt.IsDeadClick {
				deadNotice = " [DEAD CLICK]"
			}
			fmt.Printf("  [%s] Step #%d (%s): %s%s%s\n",
				p.Name, evt.Index, evt.ActionType, evt.TargetDescription, pauseNotice, deadNotice)
			if evt.InnerSpeech != "" {
				fmt.Printf("    * Inner thought: %q\n", evt.InnerSpeech)
			}
		},
		OnSessionEnd: func(p config.Persona, summary *metrics.TelemetrySummary, sus *metrics.SUSResult) {
			status := "Completed"
			if summary.DroppedOff {
				status = "Abandoned"
			}
			fmt.Printf("  --> %s: SUS %.1f (Grade %s) | FI %.2f | Status: %s\n\n",
				p.Name, sus.Score, sus.Grade, summary.FrictionIndex, status)
		},
	}

	data, err := runner.ExecuteSwarm(context.Background(), opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Study execution error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("==================================================")
	fmt.Printf("OVERALL SUS BENCHMARK : %.1f / 100 (Grade %s)\n", data.Run.AverageSUS, data.Run.Grade)
	fmt.Printf("FRICTION INDEX (FI)   : %.2f\n", data.Run.AverageFriction)
	if data.Run.MandatoryRedesign {
		fmt.Println("WARNING: High friction index (FI > 2.5). Mandatory redesign recommended.")
	}
	fmt.Printf("HEURISTIC FINDINGS    : %d violations identified\n", len(data.Heuristics))
	fmt.Println("==================================================")

	for _, h := range data.Heuristics {
		fmt.Printf("\n[%s] Heuristic %d (%s): %s\n", h.Severity, h.HeuristicID, h.HeuristicName, h.Title)
		fmt.Printf("  Evidence       : %s\n", h.Evidence)
		fmt.Printf("  Affected Users : %s\n", strings.Join(h.AffectedPersonas, ", "))
		fmt.Printf("  Recommendation : %s\n", h.Recommendation)
	}

	if *outputPath != "" {
		if err := report.ExportToFile(data, "json", *outputPath); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing json report: %v\n", err)
		} else {
			fmt.Printf("\nExported JSON report to %s\n", *outputPath)
		}
	}

	if *htmlReport != "" {
		if err := report.ExportToFile(data, "html", *htmlReport); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing html report: %v\n", err)
		} else {
			fmt.Printf("Exported standalone HTML report to %s\n", *htmlReport)
		}
	}
}

func cmdReport(args []string) {
	fs := flag.NewFlagSet("report", flag.ExitOnError)
	runID := fs.String("run-id", "", "Run identifier in database")
	inputFile := fs.String("input", "", "Alternative path to input JSON report")
	format := fs.String("format", "html", "Output format (html, markdown, json)")
	output := fs.String("output", "usability_report.html", "Output destination file")
	dbPath := fs.String("db", "uxswarm.db", "SQLite database path")
	_ = fs.Parse(args)

	var data report.ReportData

	if *inputFile != "" {
		raw, err := os.ReadFile(*inputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to read input report %s: %v\n", *inputFile, err)
			os.Exit(1)
		}
		if err := json.Unmarshal(raw, &data); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to parse json report: %v\n", err)
			os.Exit(1)
		}
	} else if *runID != "" {
		db, err := store.NewStore(*dbPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Database connection error: %v\n", err)
			os.Exit(1)
		}
		defer db.Close()

		run, err := db.GetRun(*runID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Run %s not found: %v\n", *runID, err)
			os.Exit(1)
		}
		sessions, _ := db.ListSessionsByRun(*runID)
		heuristics, _ := db.ListHeuristicsByRun(*runID)

		data.Run = *run
		data.Sessions = sessions
		data.Heuristics = heuristics
	} else {
		fmt.Fprintln(os.Stderr, "Please specify either --run-id <id> or --input <file.json>")
		os.Exit(1)
	}

	if err := report.ExportToFile(&data, *format, *output); err != nil {
		fmt.Fprintf(os.Stderr, "Error exporting report: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully generated %s report: %s\n", strings.ToUpper(*format), *output)
}

func cmdAudit(args []string) {
	fs := flag.NewFlagSet("audit", flag.ExitOnError)
	targetURL := fs.String("url", "", "Page URL to audit")
	device := fs.String("device", "desktop", "Device type (desktop or mobile)")
	outputPath := fs.String("output", "", "Optional JSON output path")
	_ = fs.Parse(args)

	if *targetURL == "" {
		fmt.Fprintln(os.Stderr, "Error: --url is required for audit")
		os.Exit(1)
	}

	isMobile := strings.ToLower(*device) == "mobile"
	personaConfig := config.Persona{
		ID:                       "audit-inspector",
		Name:                     "Usability Auditor",
		Archetype:                config.ArchetypeAccessibilityReliant,
		PatienceThresholdSeconds: 10.0,
		TechnicalLiteracy:        5,
		AttentionSpan:            "HIGH",
		Viewport: config.Viewport{
			Width:    map[bool]int{true: 390, false: 1440}[isMobile],
			Height:   map[bool]int{true: 844, false: 900}[isMobile],
			IsMobile: isMobile,
		},
	}

	fmt.Printf("Auditing %s on %s...\n", *targetURL, *device)

	data, err := runner.ExecuteSwarm(context.Background(), runner.RunOptions{
		Name:      fmt.Sprintf("Audit: %s", *targetURL),
		TargetURL: *targetURL,
		Goal:      "Inspect page layout, tap targets, and interactive controls",
		Personas:  []config.Persona{personaConfig},
		Headless:  true,
		MaxSteps:  6,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Audit error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\nAudit completed. Page SUS score: %.1f (Grade %s)\n", data.Run.AverageSUS, data.Run.Grade)
	fmt.Printf("Identified %d heuristic issues:\n", len(data.Heuristics))
	for _, h := range data.Heuristics {
		fmt.Printf("- [%s] Heuristic %d: %s -> %s\n", h.Severity, h.HeuristicID, h.Title, h.Recommendation)
	}

	if *outputPath != "" {
		_ = report.ExportToFile(data, "json", *outputPath)
		fmt.Printf("Saved audit findings to %s\n", *outputPath)
	}
}

func cmdDemo(args []string) {
	fs := flag.NewFlagSet("demo", flag.ExitOnError)
	port := fs.Int("port", 3000, "HTTP listen port")
	_ = fs.Parse(args)

	srv, err := demo.Start(*port)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start demo server: %v\n", err)
		os.Exit(1)
	}
	defer srv.Close()

	fmt.Printf("uxswarm Built-in SaaS Demo App running at: %s\n", srv.URL())
	fmt.Println("Flow routes: /signup -> /workspace -> /invite -> /complete")
	fmt.Println("Press Ctrl+C to stop.")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	fmt.Println("\nStopping demo server...")
}
