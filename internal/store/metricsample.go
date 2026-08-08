package store

import (
	"database/sql"
	"encoding/json"
	"time"
)

// MetricSample is the latest event for a metric including mode dim.
type MetricSample struct {
	Metric    string
	Value     float64
	Mode      string
	Collector string
	TS        time.Time
}

// LatestMetricSample returns the newest event for tenant+metric (any collector).
func (s *Store) LatestMetricSample(tenant, metric string) (*MetricSample, error) {
	row := s.db.QueryRow(`
SELECT collector, metric, value, COALESCE(dims,''), ts FROM events
WHERE tenant=? AND metric=? ORDER BY id DESC LIMIT 1`, tenant, metric)
	var sample MetricSample
	var dims, ts string
	if err := row.Scan(&sample.Collector, &sample.Metric, &sample.Value, &dims, &ts); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	sample.TS, _ = time.Parse(time.RFC3339, ts)
	var m map[string]string
	_ = json.Unmarshal([]byte(dims), &m)
	if m != nil {
		sample.Mode = m["mode"]
	}
	if sample.Mode == "" {
		sample.Mode = "demo"
	}
	return &sample, nil
}

// LatestMetricsByNames returns samples for each metric name (missing → absent from map).
func (s *Store) LatestMetricsByNames(tenant string, names []string) (map[string]MetricSample, error) {
	out := map[string]MetricSample{}
	for _, n := range names {
		sample, err := s.LatestMetricSample(tenant, n)
		if err != nil {
			return nil, err
		}
		if sample != nil {
			out[n] = *sample
		}
	}
	return out, nil
}
