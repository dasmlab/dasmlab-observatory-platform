package campaign

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Service stores campaigns and runs render/arm/send.
type Service struct {
	mu       sync.RWMutex
	reg      *Registry
	tenant   string
	byID     map[string]*Campaign
	briefDir string
}

func NewService(tenant, dataDir string) *Service {
	briefDir := filepath.Join("campaigns")
	if v := os.Getenv("DPO_CAMPAIGN_DIR"); v != "" {
		briefDir = v
	} else if dataDir != "" {
		// Prefer repo-mounted /app/campaigns in image; fallback dataDir.
		if _, err := os.Stat("/app/campaigns"); err == nil {
			briefDir = "/app/campaigns"
		}
	}
	s := &Service{
		reg: NewRegistry(), tenant: tenant, byID: map[string]*Campaign{}, briefDir: briefDir,
	}
	s.seedLaunch()
	return s
}

func (s *Service) Registry() *Registry { return s.reg }

func (s *Service) seedLaunch() {
	brief := loadBrief(filepath.Join(s.briefDir, "dasmlab-2.0-launch", "brief.md"))
	if brief == "" {
		brief = defaultLaunchBrief()
	}
	channels := make([]ChannelPlan, 0)
	for _, a := range s.reg.List() {
		channels = append(channels, ChannelPlan{ChannelID: a.ID(), Status: "pending"})
	}
	s.byID["dasmlab-2.0-launch"] = &Campaign{
		ID: "dasmlab-2.0-launch", Tenant: s.tenant, Type: "launch",
		Title: "DASMLAB 2.0 — Engineering Knowledge Network",
		Brief: brief, Status: "draft",
		UTM:      UTM{Medium: "campaign", Campaign: "dasmlab-2.0-launch"},
		Channels: channels, Updated: Now(),
	}
}

func defaultLaunchBrief() string {
	return `DASMLAB 2.0 is live: an Engineering Knowledge Network that answers for visitors and shows how we built it for engineers.

We also shipped the Digital Presence Observatory (DPO) — collectors, novel scores, and campaign orchestration — so we can observe the presence we create.

Start here: https://dasmlab.org/launch`
}

func loadBrief(path string) string {
	b, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(b))
}

func (s *Service) List() []Campaign {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Campaign, 0, len(s.byID))
	for _, c := range s.byID {
		out = append(out, *c)
	}
	return out
}

func (s *Service) Get(id string) (*Campaign, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	c, ok := s.byID[id]
	if !ok {
		return nil, fmt.Errorf("campaign not found")
	}
	cp := *c
	return &cp, nil
}

func (s *Service) Create(c Campaign) (*Campaign, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if c.ID == "" {
		return nil, fmt.Errorf("id required")
	}
	if _, exists := s.byID[c.ID]; exists {
		return nil, fmt.Errorf("campaign exists")
	}
	if c.Tenant == "" {
		c.Tenant = s.tenant
	}
	if c.Type == "" {
		c.Type = "launch"
	}
	if c.Status == "" {
		c.Status = "draft"
	}
	if len(c.Channels) == 0 {
		for _, a := range s.reg.List() {
			c.Channels = append(c.Channels, ChannelPlan{ChannelID: a.ID(), Status: "pending"})
		}
	}
	c.Updated = Now()
	s.byID[c.ID] = &c
	cp := c
	return &cp, nil
}

func (s *Service) Render(ctx context.Context, id string) (*Campaign, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	c, ok := s.byID[id]
	if !ok {
		return nil, fmt.Errorf("campaign not found")
	}
	baseURL := envOrCampaign("DPO_LAUNCH_URL", "https://dasmlab.org/launch")
	brief := Brief{
		CampaignID: c.ID, Title: c.Title, Body: c.Brief, URL: baseURL, UTM: c.UTM,
	}
	for i := range c.Channels {
		ch := &c.Channels[i]
		a, ok := s.reg.Get(ch.ChannelID)
		if !ok {
			ch.Status = "blocked"
			continue
		}
		arts, err := a.Render(ctx, brief)
		if err != nil {
			ch.Status = "blocked"
			continue
		}
		ch.Artifacts = arts
		ch.Status = "rendered"
		ch.Receipt = &Receipt{ChannelID: ch.ChannelID, Mode: "dry_run", At: Now()}
	}
	c.Status = "dry_run"
	c.Updated = Now()
	cp := *c
	return &cp, nil
}

func (s *Service) Arm(id string) (*Campaign, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	c, ok := s.byID[id]
	if !ok {
		return nil, fmt.Errorf("campaign not found")
	}
	if c.Status != "dry_run" && c.Status != "armed" && c.Status != "partial" {
		return nil, fmt.Errorf("render before arm (status=%s)", c.Status)
	}
	c.Status = "armed"
	c.Updated = Now()
	cp := *c
	return &cp, nil
}

func (s *Service) Send(ctx context.Context, id string) (*Campaign, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	c, ok := s.byID[id]
	if !ok {
		return nil, fmt.Errorf("campaign not found")
	}
	if c.Status != "armed" && c.Status != "partial" && os.Getenv("DPO_CAMPAIGN_ALLOW_SEND_WITHOUT_ARM") != "1" {
		return nil, fmt.Errorf("arm campaign before send (status=%s)", c.Status)
	}
	sent, manual := 0, 0
	for i := range c.Channels {
		ch := &c.Channels[i]
		a, ok := s.reg.Get(ch.ChannelID)
		if !ok {
			continue
		}
		if len(ch.Artifacts) == 0 {
			arts, err := a.Render(ctx, Brief{CampaignID: c.ID, Title: c.Title, Body: c.Brief,
				URL: envOrCampaign("DPO_LAUNCH_URL", "https://dasmlab.org/launch"), UTM: c.UTM})
			if err == nil {
				ch.Artifacts = arts
			}
		}
		if len(ch.Artifacts) == 0 {
			continue
		}
		rec, err := a.Send(ctx, ch.Artifacts[0])
		if err != nil {
			ch.Status = "blocked"
			ch.Receipt = &Receipt{ChannelID: ch.ChannelID, Mode: "manual", At: Now(), Detail: err.Error()}
			manual++
			continue
		}
		ch.Receipt = &rec
		if rec.Mode == "sent" {
			ch.Status = "sent"
			sent++
		} else {
			ch.Status = "manual"
			manual++
		}
	}
	if sent > 0 && manual == 0 {
		c.Status = "sent"
	} else if sent > 0 {
		c.Status = "partial"
	} else {
		c.Status = "partial"
	}
	c.Updated = Now()
	cp := *c
	return &cp, nil
}

func envOrCampaign(k, def string) string {
	if v := strings.TrimSpace(os.Getenv(k)); v != "" {
		return v
	}
	return def
}
