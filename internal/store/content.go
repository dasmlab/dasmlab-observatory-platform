package store

import (
	"encoding/json"
	"time"
)

type PathStats struct {
	Path        string  `json:"path"`
	PageViews   float64 `json:"page_views"`
	Engaged     float64 `json:"engaged"`
	Impressions float64 `json:"impressions"`
	Clicks      float64 `json:"clicks"`
	BotHits     float64 `json:"bot_hits"`
}

func (s *Store) UpsertEntity(tenant, kind, key string) error {
	_, err := s.db.Exec(`INSERT INTO entities (tenant, kind, key) VALUES (?, ?, ?)
ON CONFLICT(tenant, kind, key) DO NOTHING`, tenant, kind, key)
	return err
}

func (s *Store) ListEntities(tenant, kind string) ([]string, error) {
	rows, err := s.db.Query(`SELECT key FROM entities WHERE tenant=? AND kind=? ORDER BY key`, tenant, kind)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []string
	for rows.Next() {
		var k string
		if err := rows.Scan(&k); err != nil {
			return nil, err
		}
		out = append(out, k)
	}
	return out, rows.Err()
}

func (s *Store) UpsertPathDaily(tenant, path, day string, pageViews, engaged, impressions, clicks, botHits float64) error {
	_, err := s.db.Exec(`INSERT INTO path_daily (tenant, path, day, page_views, engaged, impressions, clicks, bot_hits)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(tenant, path, day) DO UPDATE SET
  page_views=excluded.page_views,
  engaged=excluded.engaged,
  impressions=excluded.impressions,
  clicks=excluded.clicks,
  bot_hits=excluded.bot_hits`, tenant, path, day, pageViews, engaged, impressions, clicks, botHits)
	return err
}

func (s *Store) TopPaths(tenant string, limit int) ([]PathStats, error) {
	if limit <= 0 {
		limit = 25
	}
	rows, err := s.db.Query(`
SELECT path,
  SUM(page_views), SUM(engaged), SUM(impressions), SUM(clicks), SUM(bot_hits)
FROM path_daily WHERE tenant=?
GROUP BY path
ORDER BY SUM(page_views)+SUM(impressions)+SUM(bot_hits) DESC
LIMIT ?`, tenant, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []PathStats
	for rows.Next() {
		var p PathStats
		if err := rows.Scan(&p.Path, &p.PageViews, &p.Engaged, &p.Impressions, &p.Clicks, &p.BotHits); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

type Baseline struct {
	Tenant    string          `json:"tenant"`
	Label     string          `json:"label"`
	CreatedAt time.Time       `json:"created_at"`
	Payload   json.RawMessage `json:"payload"`
}

func (s *Store) SaveBaseline(tenant, label string, payload any) error {
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(`INSERT INTO baselines (tenant, label, created_at, payload) VALUES (?, ?, ?, ?)
ON CONFLICT(tenant, label) DO UPDATE SET created_at=excluded.created_at, payload=excluded.payload`,
		tenant, label, time.Now().UTC().Format(time.RFC3339), string(b))
	return err
}

func (s *Store) GetBaseline(tenant, label string) (*Baseline, error) {
	row := s.db.QueryRow(`SELECT tenant, label, created_at, payload FROM baselines WHERE tenant=? AND label=?`, tenant, label)
	var bl Baseline
	var ts string
	var payload string
	if err := row.Scan(&bl.Tenant, &bl.Label, &ts, &payload); err != nil {
		return nil, err
	}
	bl.CreatedAt, _ = time.Parse(time.RFC3339, ts)
	bl.Payload = json.RawMessage(payload)
	return &bl, nil
}

func (s *Store) ListBaselines(tenant string) ([]Baseline, error) {
	rows, err := s.db.Query(`SELECT tenant, label, created_at, payload FROM baselines WHERE tenant=? ORDER BY created_at DESC`, tenant)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Baseline
	for rows.Next() {
		var bl Baseline
		var ts, payload string
		if err := rows.Scan(&bl.Tenant, &bl.Label, &ts, &payload); err != nil {
			return nil, err
		}
		bl.CreatedAt, _ = time.Parse(time.RFC3339, ts)
		bl.Payload = json.RawMessage(payload)
		out = append(out, bl)
	}
	return out, rows.Err()
}

// EventsByMetricDim aggregates recent event values keyed by a dim (e.g. path).
func (s *Store) SumMetricByDim(tenant, metric, dimKey string, since time.Time) (map[string]float64, error) {
	rows, err := s.db.Query(`SELECT dims, value FROM events WHERE tenant=? AND metric=? AND ts>=?`,
		tenant, metric, since.UTC().Format(time.RFC3339))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := map[string]float64{}
	for rows.Next() {
		var dimsRaw string
		var v float64
		if err := rows.Scan(&dimsRaw, &v); err != nil {
			return nil, err
		}
		var dims map[string]string
		_ = json.Unmarshal([]byte(dimsRaw), &dims)
		k := dims[dimKey]
		if k == "" {
			continue
		}
		out[k] += v
	}
	return out, rows.Err()
}
