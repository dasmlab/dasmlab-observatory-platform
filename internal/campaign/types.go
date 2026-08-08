package campaign

import "time"

// Campaign is a presence-engineering campaign (ADR-0402).
type Campaign struct {
	ID       string        `json:"id"`
	Tenant   string        `json:"tenant"`
	Type     string        `json:"type"` // launch | presence | research
	Title    string        `json:"title"`
	Brief    string        `json:"brief"`
	UTM      UTM           `json:"utm"`
	Status   string        `json:"status"` // draft | dry_run | armed | sent | partial
	Channels []ChannelPlan `json:"channels"`
	Updated  string        `json:"updated_at,omitempty"`
}

type UTM struct {
	Source   string `json:"source,omitempty"`
	Medium   string `json:"medium"`
	Campaign string `json:"campaign"`
}

type ChannelPlan struct {
	ChannelID string     `json:"channel_id"`
	Status    string     `json:"status"` // pending | rendered | blocked | sent | manual
	Artifacts []Artifact `json:"artifacts,omitempty"`
	Receipt   *Receipt   `json:"receipt,omitempty"`
}

type Artifact struct {
	Format      string   `json:"format"` // post | sms | email | page
	Body        string   `json:"body"`
	Subject     string   `json:"subject,omitempty"`
	Title       string   `json:"title,omitempty"`
	CharLimit   int      `json:"char_limit"`
	CharCount   int      `json:"char_count"`
	OverLimit   bool     `json:"over_limit"`
	PreviewHTML string   `json:"preview_html,omitempty"`
	MediaRefs   []string `json:"media_refs,omitempty"`
	URL         string   `json:"url,omitempty"`
}

type Receipt struct {
	ChannelID  string `json:"channel_id"`
	Mode       string `json:"mode"` // dry_run | manual | sent
	ExternalID string `json:"external_id,omitempty"`
	At         string `json:"at"`
	Detail     string `json:"detail,omitempty"`
}

// ChannelInfo is catalog + readiness.
type ChannelInfo struct {
	ID          string `json:"id"`
	Kind        string `json:"kind"` // social | sms | email | web
	Name        string `json:"name"`
	CharLimit   int    `json:"char_limit"`
	Ready       bool   `json:"ready"`
	LiveStatus  string `json:"live_status"` // ready | blocked
	SetupHint   string `json:"setup_hint,omitempty"`
	CanSend     bool   `json:"can_send"`
	CanDryRun   bool   `json:"can_dry_run"`
}

func Now() string { return time.Now().UTC().Format(time.RFC3339) }
