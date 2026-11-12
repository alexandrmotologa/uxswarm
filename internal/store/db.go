package store

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/alexandrmotologa/uxswarm/internal/metrics"
	_ "modernc.org/sqlite"
)

// RunRecord represents an automated usability study run.
type RunRecord struct {
	ID                string    `json:"id"`
	Name              string    `json:"name"`
	TargetURL         string    `json:"target_url"`
	Goal              string    `json:"goal"`
	Status            string    `json:"status"` // running, completed, failed
	TotalPersonas     int       `json:"total_personas"`
	AverageSUS        float64   `json:"average_sus"`
	AverageFriction   float64   `json:"average_friction"`
	Grade             string    `json:"grade"`
	MandatoryRedesign bool      `json:"mandatory_redesign"`
	CreatedAt         time.Time `json:"created_at"`
	CompletedAt       time.Time `json:"completed_at,omitempty"`
}

// SessionRecord holds usability metrics and state for an individual persona's session.
type SessionRecord struct {
	ID                 string    `json:"id"`
	RunID              string    `json:"run_id"`
	PersonaID          string    `json:"persona_id"`
	PersonaName        string    `json:"persona_name"`
	Archetype          string    `json:"archetype"`
	ViewportWidth      int       `json:"viewport_width"`
	ViewportHeight     int       `json:"viewport_height"`
	IsMobile           bool      `json:"is_mobile"`
	SUSScore           float64   `json:"sus_score"`
	SUSGrade           string    `json:"sus_grade"`
	FrictionIndex      float64   `json:"friction_index"`
	TotalSteps         int       `json:"total_steps"`
	TotalDuration      float64   `json:"total_duration"`
	HesitationCount    int       `json:"hesitation_count"`
	DeadClicks         int       `json:"dead_clicks"`
	Backtracks         int       `json:"backtracks"`
	GoalCompleted      bool      `json:"goal_completed"`
	DroppedOff         bool      `json:"dropped_off"`
	DropOffReason      string    `json:"drop_off_reason,omitempty"`
	SUSResponsesJSON   string    `json:"sus_responses_json,omitempty"`
	CreatedAt          time.Time `json:"created_at"`
}

// StepRecord represents an individual interaction or cognitive frame in a session.
type StepRecord struct {
	ID                int64     `json:"id"`
	SessionID         string    `json:"session_id"`
	StepIndex         int       `json:"step_index"`
	Timestamp         time.Time `json:"timestamp"`
	ActionType        string    `json:"action_type"`
	TargetSelector    string    `json:"target_selector,omitempty"`
	TargetDescription string    `json:"target_description,omitempty"`
	TargetRole        string    `json:"target_role,omitempty"`
	X                 int       `json:"x"`
	Y                 int       `json:"y"`
	IsDeadClick       bool      `json:"is_dead_click"`
	HesitationSeconds float64   `json:"hesitation_seconds"`
	URL               string    `json:"url"`
	InnerSpeech       string    `json:"inner_speech"`
	Frustration       float64   `json:"frustration"`
	ScreenshotPath    string    `json:"screenshot_path,omitempty"`
	FrameDataURI      string    `json:"frame_data_uri,omitempty"`
}

// Store wraps SQLite operations for usability data.
type Store struct {
	db *sql.DB
}

// NewStore initializes a SQLite database and runs migrations.
func NewStore(dbPath string) (*Store, error) {
	if dbPath == "" {
		dbPath = "uxswarm.db"
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open sqlite db %s: %w", dbPath, err)
	}

	// Optimize SQLite for concurrent reads and writes
	pragmas := []string{
		"PRAGMA journal_mode = WAL;",
		"PRAGMA busy_timeout = 5000;",
		"PRAGMA synchronous = NORMAL;",
		"PRAGMA foreign_keys = ON;",
	}
	for _, p := range pragmas {
		if _, err := db.Exec(p); err != nil {
			_ = db.Close()
			return nil, fmt.Errorf("execute pragma %s: %w", p, err)
		}
	}

	s := &Store{db: db}
	if err := s.migrate(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("run migrations: %w", err)
	}

	return s, nil
}

func (s *Store) migrate() error {
	schema := `
	CREATE TABLE IF NOT EXISTS runs (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		target_url TEXT NOT NULL,
		goal TEXT NOT NULL,
		status TEXT NOT NULL,
		total_personas INTEGER NOT NULL DEFAULT 0,
		average_sus REAL NOT NULL DEFAULT 0.0,
		average_friction REAL NOT NULL DEFAULT 0.0,
		grade TEXT NOT NULL DEFAULT '',
		mandatory_redesign INTEGER NOT NULL DEFAULT 0,
		created_at TIMESTAMP NOT NULL,
		completed_at TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS sessions (
		id TEXT PRIMARY KEY,
		run_id TEXT NOT NULL,
		persona_id TEXT NOT NULL,
		persona_name TEXT NOT NULL,
		archetype TEXT NOT NULL,
		viewport_width INTEGER NOT NULL,
		viewport_height INTEGER NOT NULL,
		is_mobile INTEGER NOT NULL DEFAULT 0,
		sus_score REAL NOT NULL DEFAULT 0.0,
		sus_grade TEXT NOT NULL DEFAULT '',
		friction_index REAL NOT NULL DEFAULT 0.0,
		total_steps INTEGER NOT NULL DEFAULT 0,
		total_duration REAL NOT NULL DEFAULT 0.0,
		hesitation_count INTEGER NOT NULL DEFAULT 0,
		dead_clicks INTEGER NOT NULL DEFAULT 0,
		backtracks INTEGER NOT NULL DEFAULT 0,
		goal_completed INTEGER NOT NULL DEFAULT 0,
		dropped_off INTEGER NOT NULL DEFAULT 0,
		drop_off_reason TEXT,
		sus_responses_json TEXT,
		created_at TIMESTAMP NOT NULL,
		FOREIGN KEY (run_id) REFERENCES runs(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS steps (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		session_id TEXT NOT NULL,
		step_index INTEGER NOT NULL,
		timestamp TIMESTAMP NOT NULL,
		action_type TEXT NOT NULL,
		target_selector TEXT,
		target_description TEXT,
		target_role TEXT,
		x INTEGER NOT NULL DEFAULT 0,
		y INTEGER NOT NULL DEFAULT 0,
		is_dead_click INTEGER NOT NULL DEFAULT 0,
		hesitation_seconds REAL NOT NULL DEFAULT 0.0,
		url TEXT NOT NULL DEFAULT '',
		inner_speech TEXT NOT NULL DEFAULT '',
		frustration REAL NOT NULL DEFAULT 0.0,
		screenshot_path TEXT,
		frame_data_uri TEXT,
		FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS heuristics (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		run_id TEXT NOT NULL,
		heuristic_id INTEGER NOT NULL,
		heuristic_name TEXT NOT NULL,
		severity TEXT NOT NULL,
		title TEXT NOT NULL,
		evidence TEXT NOT NULL,
		affected_personas_json TEXT NOT NULL,
		recommendation TEXT NOT NULL,
		FOREIGN KEY (run_id) REFERENCES runs(id) ON DELETE CASCADE
	);

	CREATE INDEX IF NOT EXISTS idx_sessions_run_id ON sessions(run_id);
	CREATE INDEX IF NOT EXISTS idx_steps_session_id ON steps(session_id);
	CREATE INDEX IF NOT EXISTS idx_heuristics_run_id ON heuristics(run_id);
	`
	_, err := s.db.Exec(schema)
	return err
}

// CreateRun inserts a new test run record.
func (s *Store) CreateRun(r *RunRecord) error {
	query := `
	INSERT INTO runs (id, name, target_url, goal, status, total_personas, average_sus, average_friction, grade, mandatory_redesign, created_at)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
	`
	isRedesign := 0
	if r.MandatoryRedesign {
		isRedesign = 1
	}
	_, err := s.db.Exec(query,
		r.ID, r.Name, r.TargetURL, r.Goal, r.Status,
		r.TotalPersonas, r.AverageSUS, r.AverageFriction, r.Grade,
		isRedesign, r.CreatedAt,
	)
	return err
}

// UpdateRun updates an existing run with completed statistics.
func (s *Store) UpdateRun(r *RunRecord) error {
	query := `
	UPDATE runs SET
		status = ?,
		total_personas = ?,
		average_sus = ?,
		average_friction = ?,
		grade = ?,
		mandatory_redesign = ?,
		completed_at = ?
	WHERE id = ?;
	`
	isRedesign := 0
	if r.MandatoryRedesign {
		isRedesign = 1
	}
	_, err := s.db.Exec(query,
		r.Status, r.TotalPersonas, r.AverageSUS, r.AverageFriction,
		r.Grade, isRedesign, r.CompletedAt, r.ID,
	)
	return err
}

// GetRun fetches a single run by ID.
func (s *Store) GetRun(id string) (*RunRecord, error) {
	query := `
	SELECT id, name, target_url, goal, status, total_personas, average_sus, average_friction, grade, mandatory_redesign, created_at, completed_at
	FROM runs WHERE id = ?;
	`
	var r RunRecord
	var completedAt sql.NullTime
	var isRedesign int

	err := s.db.QueryRow(query, id).Scan(
		&r.ID, &r.Name, &r.TargetURL, &r.Goal, &r.Status,
		&r.TotalPersonas, &r.AverageSUS, &r.AverageFriction, &r.Grade,
		&isRedesign, &r.CreatedAt, &completedAt,
	)
	if err != nil {
		return nil, err
	}
	r.MandatoryRedesign = isRedesign == 1
	if completedAt.Valid {
		r.CompletedAt = completedAt.Time
	}
	return &r, nil
}

// ListRuns returns all runs ordered by creation time descending.
func (s *Store) ListRuns(limit int) ([]RunRecord, error) {
	if limit <= 0 {
		limit = 50
	}
	query := `
	SELECT id, name, target_url, goal, status, total_personas, average_sus, average_friction, grade, mandatory_redesign, created_at, completed_at
	FROM runs ORDER BY created_at DESC LIMIT ?;
	`
	rows, err := s.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	runs := make([]RunRecord, 0)
	for rows.Next() {
		var r RunRecord
		var completedAt sql.NullTime
		var isRedesign int

		if err := rows.Scan(
			&r.ID, &r.Name, &r.TargetURL, &r.Goal, &r.Status,
			&r.TotalPersonas, &r.AverageSUS, &r.AverageFriction, &r.Grade,
			&isRedesign, &r.CreatedAt, &completedAt,
		); err != nil {
			return nil, err
		}
		r.MandatoryRedesign = isRedesign == 1
		if completedAt.Valid {
			r.CompletedAt = completedAt.Time
		}
		runs = append(runs, r)
	}
	return runs, nil
}

// CreateSession inserts a persona session record.
func (s *Store) CreateSession(sess *SessionRecord) error {
	query := `
	INSERT INTO sessions (
		id, run_id, persona_id, persona_name, archetype,
		viewport_width, viewport_height, is_mobile,
		sus_score, sus_grade, friction_index,
		total_steps, total_duration, hesitation_count, dead_clicks, backtracks,
		goal_completed, dropped_off, drop_off_reason, sus_responses_json, created_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
	`
	isMobile := 0
	if sess.IsMobile {
		isMobile = 1
	}
	goalCompleted := 0
	if sess.GoalCompleted {
		goalCompleted = 1
	}
	droppedOff := 0
	if sess.DroppedOff {
		droppedOff = 1
	}

	_, err := s.db.Exec(query,
		sess.ID, sess.RunID, sess.PersonaID, sess.PersonaName, sess.Archetype,
		sess.ViewportWidth, sess.ViewportHeight, isMobile,
		sess.SUSScore, sess.SUSGrade, sess.FrictionIndex,
		sess.TotalSteps, sess.TotalDuration, sess.HesitationCount, sess.DeadClicks, sess.Backtracks,
		goalCompleted, droppedOff, sess.DropOffReason, sess.SUSResponsesJSON, sess.CreatedAt,
	)
	return err
}

// ListSessionsByRun retrieves all sessions associated with a run.
func (s *Store) ListSessionsByRun(runID string) ([]SessionRecord, error) {
	query := `
	SELECT id, run_id, persona_id, persona_name, archetype,
	       viewport_width, viewport_height, is_mobile,
	       sus_score, sus_grade, friction_index,
	       total_steps, total_duration, hesitation_count, dead_clicks, backtracks,
	       goal_completed, dropped_off, drop_off_reason, sus_responses_json, created_at
	FROM sessions WHERE run_id = ? ORDER BY created_at ASC;
	`
	rows, err := s.db.Query(query, runID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sessions := make([]SessionRecord, 0)
	for rows.Next() {
		var sess SessionRecord
		var isMobile, goalCompleted, droppedOff int
		var dropReason, responsesJSON sql.NullString

		if err := rows.Scan(
			&sess.ID, &sess.RunID, &sess.PersonaID, &sess.PersonaName, &sess.Archetype,
			&sess.ViewportWidth, &sess.ViewportHeight, &isMobile,
			&sess.SUSScore, &sess.SUSGrade, &sess.FrictionIndex,
			&sess.TotalSteps, &sess.TotalDuration, &sess.HesitationCount, &sess.DeadClicks, &sess.Backtracks,
			&goalCompleted, &droppedOff, &dropReason, &responsesJSON, &sess.CreatedAt,
		); err != nil {
			return nil, err
		}

		sess.IsMobile = isMobile == 1
		sess.GoalCompleted = goalCompleted == 1
		sess.DroppedOff = droppedOff == 1
		if dropReason.Valid {
			sess.DropOffReason = dropReason.String
		}
		if responsesJSON.Valid {
			sess.SUSResponsesJSON = responsesJSON.String
		}
		sessions = append(sessions, sess)
	}
	return sessions, nil
}

// CreateStep inserts an individual interaction step record.
func (s *Store) CreateStep(st *StepRecord) error {
	query := `
	INSERT INTO steps (
		session_id, step_index, timestamp, action_type,
		target_selector, target_description, target_role,
		x, y, is_dead_click, hesitation_seconds, url,
		inner_speech, frustration, screenshot_path, frame_data_uri
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
	`
	deadClick := 0
	if st.IsDeadClick {
		deadClick = 1
	}

	res, err := s.db.Exec(query,
		st.SessionID, st.StepIndex, st.Timestamp, st.ActionType,
		st.TargetSelector, st.TargetDescription, st.TargetRole,
		st.X, st.Y, deadClick, st.HesitationSeconds, st.URL,
		st.InnerSpeech, st.Frustration, st.ScreenshotPath, st.FrameDataURI,
	)
	if err != nil {
		return err
	}
	st.ID, _ = res.LastInsertId()
	return nil
}

// ListStepsBySession retrieves all interaction steps for session playback.
func (s *Store) ListStepsBySession(sessionID string) ([]StepRecord, error) {
	query := `
	SELECT id, session_id, step_index, timestamp, action_type,
	       target_selector, target_description, target_role,
	       x, y, is_dead_click, hesitation_seconds, url,
	       inner_speech, frustration, screenshot_path, frame_data_uri
	FROM steps WHERE session_id = ? ORDER BY step_index ASC;
	`
	rows, err := s.db.Query(query, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	steps := make([]StepRecord, 0)
	for rows.Next() {
		var st StepRecord
		var deadClick int
		var sel, desc, role, shotPath, frameURI sql.NullString

		if err := rows.Scan(
			&st.ID, &st.SessionID, &st.StepIndex, &st.Timestamp, &st.ActionType,
			&sel, &desc, &role,
			&st.X, &st.Y, &deadClick, &st.HesitationSeconds, &st.URL,
			&st.InnerSpeech, &st.Frustration, &shotPath, &frameURI,
		); err != nil {
			return nil, err
		}

		st.IsDeadClick = deadClick == 1
		if sel.Valid {
			st.TargetSelector = sel.String
		}
		if desc.Valid {
			st.TargetDescription = desc.String
		}
		if role.Valid {
			st.TargetRole = role.String
		}
		if shotPath.Valid {
			st.ScreenshotPath = shotPath.String
		}
		if frameURI.Valid {
			st.FrameDataURI = frameURI.String
		}
		steps = append(steps, st)
	}
	return steps, nil
}

// SaveHeuristics persists identified usability findings for a run.
func (s *Store) SaveHeuristics(runID string, findings []metrics.HeuristicFinding) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
	INSERT INTO heuristics (run_id, heuristic_id, heuristic_name, severity, title, evidence, affected_personas_json, recommendation)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?);
	`
	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, f := range findings {
		personasJSON, _ := json.Marshal(f.AffectedPersonas)
		if _, err := stmt.Exec(
			runID, f.HeuristicID, f.HeuristicName, string(f.Severity),
			f.Title, f.Evidence, string(personasJSON), f.Recommendation,
		); err != nil {
			return err
		}
	}

	return tx.Commit()
}

// ListHeuristicsByRun retrieves heuristic findings for a given test run.
func (s *Store) ListHeuristicsByRun(runID string) ([]metrics.HeuristicFinding, error) {
	query := `
	SELECT heuristic_id, heuristic_name, severity, title, evidence, affected_personas_json, recommendation
	FROM heuristics WHERE run_id = ? ORDER BY heuristic_id ASC;
	`
	rows, err := s.db.Query(query, runID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	findings := make([]metrics.HeuristicFinding, 0)
	for rows.Next() {
		var f metrics.HeuristicFinding
		var personasJSON string
		var sevStr string

		if err := rows.Scan(
			&f.HeuristicID, &f.HeuristicName, &sevStr,
			&f.Title, &f.Evidence, &personasJSON, &f.Recommendation,
		); err != nil {
			return nil, err
		}
		f.Severity = metrics.FindingSeverity(sevStr)
		_ = json.Unmarshal([]byte(personasJSON), &f.AffectedPersonas)
		findings = append(findings, f)
	}
	return findings, nil
}

// Close closes the underlying SQLite database connection.
func (s *Store) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}
