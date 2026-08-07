package collector

import (
	"context"
	"time"
)

// Event is the normalized observability event (ADR-0002 / ADR-0003).
type Event struct {
	SchemaVersion string            `json:"schema_version"`
	Collector     string            `json:"collector"`
	Type          string            `json:"type"`
	Timestamp     time.Time         `json:"timestamp"`
	Tenant        string            `json:"tenant"`
	Entity        string            `json:"entity"`
	Metric        string            `json:"metric"`
	Value         float64           `json:"value"`
	Unit          string            `json:"unit,omitempty"`
	Dims          map[string]string `json:"dims,omitempty"`
	TraceID       string            `json:"trace_id,omitempty"`
}

// Collector is the Prometheus-exporter-like plugin contract (ADR-0003).
type Collector interface {
	Name() string
	Discover(ctx context.Context) error
	Collect(ctx context.Context) error
	Normalize(ctx context.Context) ([]Event, error)
	Health(ctx context.Context) error
}

// Registry holds enabled collectors for the Collector Manager.
type Registry struct {
	byName map[string]Collector
}

func NewRegistry() *Registry {
	return &Registry{byName: make(map[string]Collector)}
}

func (r *Registry) Register(c Collector) {
	r.byName[c.Name()] = c
}

func (r *Registry) Get(name string) (Collector, bool) {
	c, ok := r.byName[name]
	return c, ok
}

func (r *Registry) List() []Collector {
	out := make([]Collector, 0, len(r.byName))
	for _, c := range r.byName {
		out = append(out, c)
	}
	return out
}
