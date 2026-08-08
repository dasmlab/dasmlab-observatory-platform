package campaign

import (
	"context"
	"fmt"
	"os"
	"strings"
)

// Brief is the input to channel renderers.
type Brief struct {
	CampaignID string
	Title      string
	Body       string
	URL        string
	UTM        UTM
}

// Adapter is the channel plugin contract (ADR-0402 / ADR-0008).
type Adapter interface {
	ID() string
	Kind() string
	Name() string
	CharLimit() int
	SetupHint() string
	CredsPresent() bool
	Render(ctx context.Context, b Brief) ([]Artifact, error)
	Send(ctx context.Context, art Artifact) (Receipt, error)
}

// Registry holds Wave-1 adapters.
type Registry struct {
	byID map[string]Adapter
}

func NewRegistry() *Registry {
	r := &Registry{byID: map[string]Adapter{}}
	for _, a := range DefaultAdapters() {
		r.byID[a.ID()] = a
	}
	return r
}

func (r *Registry) Get(id string) (Adapter, bool) {
	a, ok := r.byID[id]
	return a, ok
}

func (r *Registry) List() []Adapter {
	order := []string{
		"web_slash", "linkedin", "x_twitter", "bluesky", "mastodon",
		"meta_threads", "facebook_page", "reddit", "email", "sms",
	}
	out := make([]Adapter, 0, len(order))
	for _, id := range order {
		if a, ok := r.byID[id]; ok {
			out = append(out, a)
		}
	}
	return out
}

func (r *Registry) ChannelInfos() []ChannelInfo {
	var out []ChannelInfo
	for _, a := range r.List() {
		ready := a.CredsPresent() || a.ID() == "web_slash"
		live := "blocked"
		if ready {
			live = "ready"
		}
		out = append(out, ChannelInfo{
			ID: a.ID(), Kind: a.Kind(), Name: a.Name(), CharLimit: a.CharLimit(),
			Ready: ready, LiveStatus: live, SetupHint: a.SetupHint(),
			CanSend: ready && a.ID() != "web_slash", CanDryRun: true,
		})
	}
	return out
}

func TrackedURL(base string, utm UTM, channel string) string {
	base = strings.TrimRight(strings.TrimSpace(base), "/")
	if base == "" {
		base = "https://dasmlab.org/launch"
	}
	src := utm.Source
	if src == "" {
		src = channel
	}
	med := utm.Medium
	if med == "" {
		med = "campaign"
	}
	camp := utm.Campaign
	if camp == "" {
		camp = "dasmlab-2.0-launch"
	}
	sep := "?"
	if strings.Contains(base, "?") {
		sep = "&"
	}
	return fmt.Sprintf("%s%sutm_source=%s&utm_medium=%s&utm_campaign=%s",
		base, sep, urlQuery(src), urlQuery(med), urlQuery(camp))
}

func urlQuery(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(s, " ", "+"), "&", "")
}

func envAny(keys ...string) bool {
	for _, k := range keys {
		if strings.TrimSpace(os.Getenv(k)) != "" {
			return true
		}
	}
	return false
}

func truncate(s string, n int) string {
	if n <= 0 || len(s) <= n {
		return s
	}
	if n < 4 {
		return s[:n]
	}
	return s[:n-1] + "…"
}

func previewHTML(title, body string) string {
	esc := func(s string) string {
		s = strings.ReplaceAll(s, "&", "&amp;")
		s = strings.ReplaceAll(s, "<", "&lt;")
		s = strings.ReplaceAll(s, ">", "&gt;")
		return s
	}
	return "<div class=\"camp-preview\"><strong>" + esc(title) + "</strong><pre>" + esc(body) + "</pre></div>"
}

func artifact(format, body string, limit int, url string) Artifact {
	a := Artifact{
		Format: format, Body: body, CharLimit: limit, CharCount: len(body),
		OverLimit: limit > 0 && len(body) > limit, URL: url,
		PreviewHTML: previewHTML(format, body),
	}
	return a
}
