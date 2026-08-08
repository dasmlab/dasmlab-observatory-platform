package collectors

import (
	"os"
	"path/filepath"

	collector "github.com/dasmlab/dasmlab-observatory-platform/platform/collector-sdk"
)

// Default returns Phase-1 collectors. Live paths activate when env credentials / logs exist.
func Default(tenant, dataDir string) []collector.Collector {
	edgePath := filepath.Join(dataDir, "edge", "access.log")
	if v := os.Getenv("DPO_EDGE_LOG_PATH"); v != "" {
		edgePath = v
	}
	return []collector.Collector{
		&gscCollector{tenant: tenant},
		&githubCollector{tenant: tenant},
		&edgeCollector{tenant: tenant, logPath: edgePath},
		&activityCollector{tenant: tenant},
		&siteCollector{tenant: tenant},
		// Family product runtimes (F1/F2 novel scores → DUO).
		&dcoCollector{tenant: tenant},
		&dsoCollector{tenant: tenant},
		&dnoCollector{tenant: tenant},
		&daoCollector{tenant: tenant, dataDir: dataDir},
		&daopsCollector{tenant: tenant},
		&dioCollector{tenant: tenant},
	}
}

type missingErr string

func (e missingErr) Error() string { return string(e) }

func errMissing(env string) error { return missingErr("missing " + env) }
