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
	"github.com/alexandrmotologa/uxswarm/internal/compare"
	"github.com/alexandrmotologa/uxswarm/internal/config"
	"github.com/alexandrmotologa/uxswarm/internal/demo"
	"github.com/alexandrmotologa/uxswarm/internal/metrics"
	"github.com/alexandrmotologa/uxswarm/internal/patch"
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
	case "compare", "diff":
		cmdCompare(os.Args[2:])
	case "patch":
		cmdPatch(os.Args[2:])
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
  compare   Run an A/B usability study comparing two URLs with delta metrics
  patch     Generate CSS/Tailwind/HTML code patches for identified issues
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

	// Usability Budgets & CI/CD Quality Gates
	minSUS := fs.Float64("min-sus", 0.0, "Minimum required SUS score (fails CI if not met)")
	maxFriction := fs.Float64("max-friction", 0.0, "Maximum allowed Friction Index (fails CI if exceeded)")
	allowAbandonment := fs.Bool("allow-abandonment", true, "Allow persona abandonment (fails CI if false and drop-offs occur)")
	githubSummary := fs.Bool("github-summary", false, "Write Markdown report to GITHUB_STEP_SUMMARY environment variable")
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
		Name:      suiteName,
		TargetURL: *targetURL,
		Goal:      *goal,
		Personas:  personas,
		Headless:  *headless,
		Store:     db,
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

	// GitHub Actions Step Summary support
	if *githubSummary || os.Getenv("GITHUB_STEP_SUMMARY") != "" {
		summaryPath := os.Getenv("GITHUB_STEP_SUMMARY")
		if summaryPath != "" {
			mdReport, _ := report.GenerateMarkdown(data)
			_ = os.WriteFile(summaryPath, []byte(mdReport), 0644)
			fmt.Println("Appended report to GITHUB_STEP_SUMMARY.")
		}
	}

	// Usability Budget Validations
	failedBudgets := make([]string, 0)
	if *minSUS > 0.0 && data.Run.AverageSUS < *minSUS {
		failedBudgets = append(failedBudgets, fmt.Sprintf("Average SUS %.1f is below required budget %.1f", data.Run.AverageSUS, *minSUS))
	}
	if *maxFriction > 0.0 && data.Run.AverageFriction > *maxFriction {
		failedBudgets = append(failedBudgets, fmt.Sprintf("Friction Index %.2f exceeds maximum threshold %.2f", data.Run.AverageFriction, *maxFriction))
	}
	if !*allowAbandonment {
		for _, s := range data.Sessions {
			if s.DroppedOff {
				failedBudgets = append(failedBudgets, fmt.Sprintf("Persona %s abandoned session (zero abandonment allowed)", s.PersonaName))
				break
			}
		}
	}

	if len(failedBudgets) > 0 {
		fmt.Printf("\n[CI USABILITY BUDGET FAILURE]\n")
		for _, b := range failedBudgets {
			fmt.Printf("  ! %s\n", b)
		}
		os.Exit(1)
	}
}

func cmdCompare(args []string) {
	fs := flag.NewFlagSet("compare", flag.ExitOnError)
	urlA := fs.String("url-a", "", "Baseline variant URL (Variant A)")
	urlB := fs.String("url-b", "", "Candidate variant URL (Variant B)")
	cfgPath := fs.String("config", "", "Path to personas.yaml suite")
	goal := fs.String("goal", "Complete onboarding flow", "Task goal description")
	outputPath := fs.String("output", "ab_comparison_report.html", "Comparison report destination")
	format := fs.String("format", "html", "Output format (html, markdown, or json)")
	_ = fs.Parse(args)

	if *urlA == "" || *urlB == "" {
		fmt.Fprintln(os.Stderr, "Error: both --url-a and --url-b are required for comparison")
		os.Exit(1)
	}

	var personas []config.Persona
	if *cfgPath != "" {
		s, err := config.LoadSuite(*cfgPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading suite: %v\n", err)
			os.Exit(1)
		}
		personas = s.Personas
	} else {
		personas = config.GetDefaultPersonas()
	}

	fmt.Printf("\n=== uxswarm A/B Usability Comparison ===\n")
	fmt.Printf("Variant A : %s\n", *urlA)
	fmt.Printf("Variant B : %s\n", *urlB)
	fmt.Printf("Swarm     : %d personas\n\n", len(personas))

	fmt.Println("--> Running swarm across Variant A (Baseline)...")
	repA, err := runner.ExecuteSwarm(context.Background(), runner.RunOptions{
		Name:      "Variant A (Baseline)",
		TargetURL: *urlA,
		Goal:      *goal,
		Personas:  personas,
		Headless:  true,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Variant A run error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Variant A Finished: SUS %.1f (Grade %s) | FI %.2f\n\n", repA.Run.AverageSUS, repA.Run.Grade, repA.Run.AverageFriction)

	fmt.Println("--> Running swarm across Variant B (Candidate)...")
	repB, err := runner.ExecuteSwarm(context.Background(), runner.RunOptions{
		Name:      "Variant B (Candidate)",
		TargetURL: *urlB,
		Goal:      *goal,
		Personas:  personas,
		Headless:  true,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Variant B run error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Variant B Finished: SUS %.1f (Grade %s) | FI %.2f\n\n", repB.Run.AverageSUS, repB.Run.Grade, repB.Run.AverageFriction)

	cmp := compare.CompareVariants(repA, repB)

	fmt.Println("==================================================")
	fmt.Printf("A/B EVALUATION WINNER : %s\n", cmp.Winner)
	fmt.Printf("DELTA SUS SCORE       : %+.1f points\n", cmp.DeltaSUS)
	fmt.Printf("DELTA FRICTION (FI)   : %+.2f\n", cmp.DeltaFriction)
	fmt.Printf("DELTA DROP-OFFS       : %+d\n", cmp.DeltaDropOffs)
	fmt.Println("==================================================")
	fmt.Printf("Recommendation: %s\n", cmp.Recommendation)

	if err := compare.ExportDiffToFile(cmp, *format, *outputPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error exporting diff report: %v\n", err)
	} else {
		fmt.Printf("\nGenerated %s differential report: %s\n", strings.ToUpper(*format), *outputPath)
	}
}

func cmdPatch(args []string) {
	fs := flag.NewFlagSet("patch", flag.ExitOnError)
	inputFile := fs.String("input", "ux_report.json", "JSON report containing heuristic findings")
	framework := fs.String("format", "all", "Filter patch format: css, tailwind, html, react, all")
	_ = fs.Parse(args)

	raw, err := os.ReadFile(*inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read report %s: %v\n", *inputFile, err)
		os.Exit(1)
	}

	var data report.ReportData
	if err := json.Unmarshal(raw, &data); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse JSON report: %v\n", err)
		os.Exit(1)
	}

	steps := make([]metrics.StepEvent, 0)
	for _, sessSteps := range data.StepsMap {
		for _, s := range sessSteps {
			steps = append(steps, metrics.StepEvent{
				ActionType:        s.ActionType,
				TargetSelector:    s.TargetSelector,
				TargetDescription: s.TargetDescription,
				IsDeadClick:       s.IsDeadClick,
			})
		}
	}

	patches := patch.GeneratePatches(data.Heuristics, steps)
	if len(patches) == 0 {
		fmt.Println("No patches required. No actionable heuristic violations found.")
		return
	}

	fmt.Printf("Synthesized %d engineering code patches:\n\n", len(patches))
	targetFmt := strings.ToLower(*framework)
	for _, p := range patches {
		if targetFmt != "all" && !strings.Contains(strings.ToLower(p.Framework), targetFmt) {
			continue
		}
		fmt.Println(patch.FormatPatchText(p))
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
