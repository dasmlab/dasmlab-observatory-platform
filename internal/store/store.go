package store

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "modernc.org/sqlite"

	collector "github.com/dasmlab/dasmlab-observatory-platform/platform/collector-sdk"
)

type Store struct {
	db *sql.DB
}

func Open(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path+"?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)")
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	s := &Store{db: db}
	if err := s.migrate(); err != nil {
		_ = db.Close()
		return nil, err
	}
	return s, nil
}

func (s *Store) Close() error { return s.db.Close() }

func (s *Store) migrate() error {
	_, err := s.db.Exec(`
CREATE TABLE IF NOT EXISTS events (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  collector TEXT NOT NULL,
  type TEXT NOT NULL,
  ts TEXT NOT NULL,
  tenant TEXT NOT NULL,
  entity TEXT NOT NULL,
  metric TEXT NOT NULL,
  value REAL NOT NULL,
  unit TEXT,
  dims TEXT
);
CREATE TABLE IF NOT EXISTS scores (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  tenant TEXT NOT NULL,
  name TEXT NOT NULL,
  value REAL NOT NULL,
  components TEXT NOT NULL,
  created_at TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS collector_status (
  name TEXT PRIMARY KEY,
  healthy INTEGER NOT NULL,
  message TEXT,
  last_run TEXT,
  last_error TEXT
);
CREATE TABLE IF NOT EXISTS entities (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  tenant TEXT NOT NULL,
  kind TEXT NOT NULL,
  key TEXT NOT NULL,
  UNIQUE(tenant, kind, key)
);
CREATE TABLE IF NOT EXISTS path_daily (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  tenant TEXT NOT NULL,
  path TEXT NOT NULL,
  day TEXT NOT NULL,
  page_views REAL DEFAULT 0,
  engaged REAL DEFAULT 0,
  impressions REAL DEFAULT 0,
  clicks REAL DEFAULT 0,
  bot_hits REAL DEFAULT 0,
  UNIQUE(tenant, path, day)
);
CREATE TABLE IF NOT EXISTS baselines (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  tenant TEXT NOT NULL,
  label TEXT NOT NULL,
  created_at TEXT NOT NULL,
  payload TEXT NOT NULL,
  UNIQUE(tenant, label)
);
`)
	return err
}

func (s *Store) InsertEvents(events []collector.Event) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()
	stmt, err := tx.Prepare(`INSERT INTO events (collector, type, ts, tenant, entity, metric, value, unit, dims)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, e := range events {
		dims, _ := json.Marshal(e.Dims)
		if _, err := stmt.Exec(e.Collector, e.Type, e.Timestamp.UTC().Format(time.RFC3339), e.Tenant, e.Entity, e.Metric, e.Value, e.Unit, string(dims)); err != nil {
			return err
		}
	}
	return tx.Commit()
}

type ScoreSnapshot struct {
	Tenant     string             `json:"tenant"`
	Name       string             `json:"name"`
	Value      float64            `json:"value"`
	Components map[string]float64 `json:"components"`
	CreatedAt  time.Time          `json:"created_at"`
}

func (s *Store) SaveScore(snap ScoreSnapshot) error {
	comp, _ := json.Marshal(snap.Components)
	_, err := s.db.Exec(`INSERT INTO scores (tenant, name, value, components, created_at) VALUES (?, ?, ?, ?, ?)`,
		snap.Tenant, snap.Name, snap.Value, string(comp), snap.CreatedAt.UTC().Format(time.RFC3339))
	return err
}

func (s *Store) LatestScore(tenant, name string) (*ScoreSnapshot, error) {
	row := s.db.QueryRow(`SELECT tenant, name, value, components, created_at FROM scores WHERE tenant=? AND name=? ORDER BY id DESC LIMIT 1`, tenant, name)
	var snap ScoreSnapshot
	var comp string
	var ts string
	if err := row.Scan(&snap.Tenant, &snap.Name, &snap.Value, &comp, &ts); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	_ = json.Unmarshal([]byte(comp), &snap.Components)
	snap.CreatedAt, _ = time.Parse(time.RFC3339, ts)
	return &snap, nil
}

func (s *Store) ScoreHistory(tenant, name string, limit int) ([]ScoreSnapshot, error) {
	if limit <= 0 {
		limit = 30
	}
	rows, err := s.db.Query(`SELECT tenant, name, value, components, created_at FROM scores WHERE tenant=? AND name=? ORDER BY id DESC LIMIT ?`, tenant, name, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []ScoreSnapshot
	for rows.Next() {
		var snap ScoreSnapshot
		var comp, ts string
		if err := rows.Scan(&snap.Tenant, &snap.Name, &snap.Value, &comp, &ts); err != nil {
			return nil, err
		}
		_ = json.Unmarshal([]byte(comp), &snap.Components)
		snap.CreatedAt, _ = time.Parse(time.RFC3339, ts)
		out = append(out, snap)
	}
	return out, rows.Err()
}

type CollectorStatus struct {
	Name      string `json:"name"`
	Healthy   bool   `json:"healthy"`
	Message   string `json:"message,omitempty"`
	LastRun   string `json:"last_run,omitempty"`
	LastError string `json:"last_error,omitempty"`
}

func (s *Store) UpsertCollectorStatus(st CollectorStatus) error {
	h := 0
	if st.Healthy {
		h = 1
	}
	_, err := s.db.Exec(`INSERT INTO collector_status (name, healthy, message, last_run, last_error) VALUES (?, ?, ?, ?, ?)
ON CONFLICT(name) DO UPDATE SET healthy=excluded.healthy, message=excluded.message, last_run=excluded.last_run, last_error=excluded.last_error`,
		st.Name, h, st.Message, st.LastRun, st.LastError)
	return err
}

func (s *Store) ListCollectorStatus() ([]CollectorStatus, error) {
	rows, err := s.db.Query(`SELECT name, healthy, COALESCE(message,''), COALESCE(last_run,''), COALESCE(last_error,'') FROM collector_status ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []CollectorStatus
	for rows.Next() {
		var st CollectorStatus
		var h int
		if err := rows.Scan(&st.Name, &h, &st.Message, &st.LastRun, &st.LastError); err != nil {
			return nil, err
		}
		st.Healthy = h == 1
		out = append(out, st)
	}
	return out, rows.Err()
}

func (s *Store) LatestMetrics(tenant string, since time.Time) (map[string]float64, error) {
	rows, err := s.db.Query(`
SELECT metric, value FROM events
WHERE tenant=? AND ts>=? AND id IN (
  SELECT MAX(id) FROM events WHERE tenant=? AND ts>=? GROUP BY metric
)`, tenant, since.UTC().Format(time.RFC3339), tenant, since.UTC().Format(time.RFC3339))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := map[string]float64{}
	for rows.Next() {
		var m string
		var v float64
		if err := rows.Scan(&m, &v); err != nil {
			return nil, err
		}
		out[m] = v
	}
	return out, rows.Err()
}

func (s *Store) MetricAvg(tenant, metric string, since time.Time) (float64, error) {
	row := s.db.QueryRow(`SELECT AVG(value) FROM events WHERE tenant=? AND metric=? AND ts>=?`, tenant, metric, since.UTC().Format(time.RFC3339))
	var v sql.NullFloat64
	if err := row.Scan(&v); err != nil {
		return 0, err
	}
	if !v.Valid {
		return 0, fmt.Errorf("no samples")
	}
	return v.Float64, nil
}

