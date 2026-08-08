package products

import (
	"time"

	"github.com/dasmlab/dasmlab-observatory-platform/internal/store"
)

// FromStore builds a Snapshot for a product code from latest metric events + optional overall score.
func FromStore(st *store.Store, tenant, code string) Snapshot {
	spec, ok := Spec(code)
	if !ok {
		return Snapshot{Code: code, Status: "scaffold", Mode: "scaffold"}
	}
	snap := Snapshot{
		Code:     code,
		Features: append([]string{}, spec.Features...),
		Proof:    append([]string{}, spec.Proof...),
		Scores:   nil,
	}

	if code == "dpo" {
		if sc, err := st.LatestScore(tenant, "overall"); err == nil && sc != nil {
			snap.Scores = append(snap.Scores, Score{Name: "overall", Value: sc.Value, Mode: "live"})
			snap.LastRun = sc.CreatedAt.UTC().Format(time.RFC3339)
		}
	} else if code == "duo" {
		// Filled by API from Compose; placeholder scaffold if called raw.
		snap.Scores = []Score{
			{Name: "business_impact", Value: 0, Mode: "scaffold"},
			{Name: "engineering_impact", Value: 0, Mode: "scaffold"},
			{Name: "operational_impact", Value: 0, Mode: "scaffold"},
		}
	} else {
		samples, _ := st.LatestMetricsByNames(tenant, spec.Metrics)
		var latest time.Time
		for _, name := range spec.Metrics {
			if sample, ok := samples[name]; ok {
				snap.Scores = append(snap.Scores, Score{Name: name, Value: sample.Value, Mode: sample.Mode})
				if sample.TS.After(latest) {
					latest = sample.TS
				}
			} else {
				snap.Scores = append(snap.Scores, Score{Name: name, Value: 0, Mode: "scaffold"})
			}
		}
		if !latest.IsZero() {
			snap.LastRun = latest.UTC().Format(time.RFC3339)
		}
	}

	snap.Mode = AggregateMode(snap.Scores)
	snap.Status = StatusFromMode(snap.Mode)
	return snap
}

// ListFromStore returns snapshots for all specs.
func ListFromStore(st *store.Store, tenant string) []Snapshot {
	specs := Specs()
	out := make([]Snapshot, 0, len(specs))
	for _, sp := range specs {
		out = append(out, FromStore(st, tenant, sp.Code))
	}
	return out
}

// SourceScoresForDUO flattens sibling product scores for Compose (excludes duo itself).
func SourceScoresForDUO(st *store.Store, tenant string) []Score {
	var out []Score
	for _, code := range []string{"dco", "dso", "dno", "dao", "daops", "dio"} {
		snap := FromStore(st, tenant, code)
		for _, sc := range snap.Scores {
			if sc.Mode == "scaffold" {
				continue
			}
			out = append(out, Score{Name: code + "." + sc.Name, Value: sc.Value, Mode: sc.Mode})
		}
	}
	return out
}
