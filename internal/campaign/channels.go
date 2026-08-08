package campaign

import (
	"context"
	"fmt"
	"os"
	"strings"
)

func DefaultAdapters() []Adapter {
	return []Adapter{
		&webSlashAdapter{},
		&socialAdapter{id: "linkedin", kind: "social", name: "LinkedIn", limit: 3000,
			hint: "LINKEDIN_ACCESS_TOKEN + LINKEDIN_ORG_URN", credKeys: []string{"LINKEDIN_ACCESS_TOKEN"}},
		&socialAdapter{id: "x_twitter", kind: "social", name: "X (Twitter)", limit: 280,
			hint: "X_BEARER_TOKEN or X_API_KEY+X_ACCESS_TOKEN", credKeys: []string{"X_BEARER_TOKEN", "X_ACCESS_TOKEN"}},
		&socialAdapter{id: "bluesky", kind: "social", name: "Bluesky", limit: 300,
			hint: "BLUESKY_HANDLE + BLUESKY_APP_PASSWORD", credKeys: []string{"BLUESKY_HANDLE", "BLUESKY_APP_PASSWORD"}},
		&socialAdapter{id: "mastodon", kind: "social", name: "Mastodon", limit: 500,
			hint: "MASTODON_BASE_URL + MASTODON_ACCESS_TOKEN", credKeys: []string{"MASTODON_ACCESS_TOKEN"}},
		&socialAdapter{id: "meta_threads", kind: "social", name: "Threads", limit: 500,
			hint: "THREADS_ACCESS_TOKEN + THREADS_USER_ID", credKeys: []string{"THREADS_ACCESS_TOKEN"}},
		&socialAdapter{id: "facebook_page", kind: "social", name: "Facebook Page", limit: 500,
			hint: "FACEBOOK_PAGE_ID + FACEBOOK_PAGE_TOKEN", credKeys: []string{"FACEBOOK_PAGE_TOKEN"}},
		&redditAdapter{},
		&emailAdapter{},
		&smsAdapter{},
	}
}

type webSlashAdapter struct{}

func (a *webSlashAdapter) ID() string          { return "web_slash" }
func (a *webSlashAdapter) Kind() string        { return "web" }
func (a *webSlashAdapter) Name() string        { return "Web slash (/launch)" }
func (a *webSlashAdapter) CharLimit() int      { return 155 }
func (a *webSlashAdapter) SetupHint() string   { return "Ship dasmlab_home /launch via GitOps" }
func (a *webSlashAdapter) CredsPresent() bool  { return true }
func (a *webSlashAdapter) Render(_ context.Context, b Brief) ([]Artifact, error) {
	url := TrackedURL(b.URL, b.UTM, "web")
	title := truncate(b.Title, 60)
	desc := truncate(strings.ReplaceAll(b.Body, "\n", " "), 155)
	body := fmt.Sprintf("SEO title: %s\nMeta description: %s\nCanonical: %s", title, desc, url)
	art := artifact("page", body, 155, url)
	art.Title = title
	art.Subject = desc
	return []Artifact{art}, nil
}
func (a *webSlashAdapter) Send(_ context.Context, art Artifact) (Receipt, error) {
	return Receipt{ChannelID: a.ID(), Mode: "manual", At: Now(),
		Detail: "Publish via home GitOps; page is the artifact: " + art.URL}, nil
}

type socialAdapter struct {
	id, kind, name, hint string
	limit                int
	credKeys             []string
}

func (a *socialAdapter) ID() string         { return a.id }
func (a *socialAdapter) Kind() string       { return a.kind }
func (a *socialAdapter) Name() string       { return a.name }
func (a *socialAdapter) CharLimit() int     { return a.limit }
func (a *socialAdapter) SetupHint() string  { return a.hint }
func (a *socialAdapter) CredsPresent() bool {
	if a.id == "bluesky" {
		return envAny("BLUESKY_HANDLE") && envAny("BLUESKY_APP_PASSWORD")
	}
	if a.id == "mastodon" {
		return envAny("MASTODON_BASE_URL") && envAny("MASTODON_ACCESS_TOKEN")
	}
	if a.id == "x_twitter" {
		return envAny("X_BEARER_TOKEN") || (envAny("X_API_KEY") && envAny("X_ACCESS_TOKEN"))
	}
	return envAny(a.credKeys...)
}

func (a *socialAdapter) Render(_ context.Context, b Brief) ([]Artifact, error) {
	url := TrackedURL(b.URL, b.UTM, a.id)
	core := strings.TrimSpace(b.Title) + "\n\n" + strings.TrimSpace(b.Body)
	tag := "\n\n" + url
	body := core + tag
	if a.limit > 0 && len(body) > a.limit {
		// Prefer keeping URL; trim core.
		room := a.limit - len(tag) - 1
		if room < 40 {
			body = truncate(core, a.limit)
		} else {
			body = truncate(core, room) + tag
		}
	}
	return []Artifact{artifact("post", body, a.limit, url)}, nil
}

func (a *socialAdapter) Send(ctx context.Context, art Artifact) (Receipt, error) {
	if !a.CredsPresent() {
		return Receipt{ChannelID: a.ID(), Mode: "manual", At: Now(),
			Detail: "No credentials — copy body manually. " + a.hint}, nil
	}
	// Live HTTP send stubs: attempt only when creds present; return honest detail.
	switch a.id {
	case "bluesky":
		return sendBluesky(ctx, art.Body)
	case "mastodon":
		return sendMastodon(ctx, art.Body)
	default:
		return Receipt{ChannelID: a.ID(), Mode: "manual", At: Now(),
			Detail: "Credentials present but live client not yet wired for " + a.id + "; use manual post this wave."}, nil
	}
}

type redditAdapter struct{}

func (a *redditAdapter) ID() string         { return "reddit" }
func (a *redditAdapter) Kind() string       { return "social" }
func (a *redditAdapter) Name() string       { return "Reddit" }
func (a *redditAdapter) CharLimit() int     { return 300 }
func (a *redditAdapter) SetupHint() string  { return "REDDIT_CLIENT_ID/SECRET + username/password + USER_AGENT" }
func (a *redditAdapter) CredsPresent() bool {
	return envAny("REDDIT_CLIENT_ID") && envAny("REDDIT_CLIENT_SECRET") && envAny("REDDIT_USERNAME")
}
func (a *redditAdapter) Render(_ context.Context, b Brief) ([]Artifact, error) {
	url := TrackedURL(b.URL, b.UTM, "reddit")
	title := truncate(b.Title, 300)
	body := strings.TrimSpace(b.Body) + "\n\n" + url
	art := artifact("post", body, 40000, url)
	art.Title = title
	return []Artifact{art}, nil
}
func (a *redditAdapter) Send(_ context.Context, art Artifact) (Receipt, error) {
	mode := "manual"
	detail := a.SetupHint()
	if a.CredsPresent() {
		detail = "Credentials present; live submit not wired this wave — post manually to chosen subreddit."
	}
	return Receipt{ChannelID: a.ID(), Mode: mode, At: Now(), Detail: detail}, nil
}

type emailAdapter struct{}

func (a *emailAdapter) ID() string         { return "email" }
func (a *emailAdapter) Kind() string       { return "email" }
func (a *emailAdapter) Name() string       { return "Email (Resend)" }
func (a *emailAdapter) CharLimit() int     { return 78 }
func (a *emailAdapter) SetupHint() string  { return "RESEND_API_KEY + RESEND_FROM (+ opt-in audience)" }
func (a *emailAdapter) CredsPresent() bool { return envAny("RESEND_API_KEY") && envAny("RESEND_FROM") }
func (a *emailAdapter) Render(_ context.Context, b Brief) ([]Artifact, error) {
	url := TrackedURL(b.URL, b.UTM, "email")
	subj := truncate(b.Title, 78)
	body := strings.TrimSpace(b.Body) + "\n\nRead more: " + url
	art := artifact("email", body, 0, url)
	art.Subject = subj
	return []Artifact{art}, nil
}
func (a *emailAdapter) Send(ctx context.Context, art Artifact) (Receipt, error) {
	if !a.CredsPresent() {
		return Receipt{ChannelID: a.ID(), Mode: "manual", At: Now(), Detail: a.SetupHint()}, nil
	}
	return sendResend(ctx, art)
}

type smsAdapter struct{}

func (a *smsAdapter) ID() string         { return "sms" }
func (a *smsAdapter) Kind() string       { return "sms" }
func (a *smsAdapter) Name() string       { return "SMS (Twilio)" }
func (a *smsAdapter) CharLimit() int     { return 160 }
func (a *smsAdapter) SetupHint() string  { return "TWILIO_ACCOUNT_SID + AUTH_TOKEN + FROM_NUMBER; opt-in list only" }
func (a *smsAdapter) CredsPresent() bool {
	return envAny("TWILIO_ACCOUNT_SID") && envAny("TWILIO_AUTH_TOKEN") && envAny("TWILIO_FROM_NUMBER")
}
func (a *smsAdapter) Render(_ context.Context, b Brief) ([]Artifact, error) {
	url := TrackedURL(b.URL, b.UTM, "sms")
	// Short SMS: title + URL
	body := truncate(b.Title, 100) + " " + url
	if len(body) > 160 {
		body = truncate(b.Title, 40) + " " + url
	}
	segs := 1
	if len(body) > 160 {
		segs = (len(body) + 152) / 153
	}
	art := artifact("sms", body, 160, url)
	art.PreviewHTML = previewHTML(fmt.Sprintf("sms (%d seg)", segs), body)
	return []Artifact{art}, nil
}
func (a *smsAdapter) Send(ctx context.Context, art Artifact) (Receipt, error) {
	if !a.CredsPresent() {
		return Receipt{ChannelID: a.ID(), Mode: "manual", At: Now(), Detail: a.SetupHint()}, nil
	}
	to := strings.TrimSpace(os.Getenv("TWILIO_TO_NUMBER"))
	if to == "" {
		return Receipt{ChannelID: a.ID(), Mode: "manual", At: Now(),
			Detail: "TWILIO_TO_NUMBER opt-in recipient required for send"}, nil
	}
	return sendTwilio(ctx, to, art.Body)
}
