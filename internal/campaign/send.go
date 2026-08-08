package campaign

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func httpClient() *http.Client {
	return &http.Client{Timeout: 20 * time.Second}
}

func sendBluesky(ctx context.Context, text string) (Receipt, error) {
	handle := strings.TrimSpace(os.Getenv("BLUESKY_HANDLE"))
	pass := strings.TrimSpace(os.Getenv("BLUESKY_APP_PASSWORD"))
	sessReq, _ := http.NewRequestWithContext(ctx, http.MethodPost,
		"https://bsky.social/xrpc/com.atproto.server.createSession",
		strings.NewReader(fmt.Sprintf(`{"identifier":%q,"password":%q}`, handle, pass)))
	sessReq.Header.Set("Content-Type", "application/json")
	resp, err := httpClient().Do(sessReq)
	if err != nil {
		return Receipt{}, err
	}
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	_ = resp.Body.Close()
	if resp.StatusCode >= 300 {
		return Receipt{ChannelID: "bluesky", Mode: "manual", At: Now(),
			Detail: fmt.Sprintf("session failed: %s", truncate(string(body), 120))}, nil
	}
	var sess struct {
		AccessJwt string `json:"accessJwt"`
		Did       string `json:"did"`
	}
	_ = json.Unmarshal(body, &sess)
	payload := fmt.Sprintf(`{"repo":%q,"collection":"app.bsky.feed.post","record":{"$type":"app.bsky.feed.post","text":%q,"createdAt":%q}}`,
		sess.Did, text, time.Now().UTC().Format(time.RFC3339))
	postReq, _ := http.NewRequestWithContext(ctx, http.MethodPost,
		"https://bsky.social/xrpc/com.atproto.repo.createRecord", strings.NewReader(payload))
	postReq.Header.Set("Content-Type", "application/json")
	postReq.Header.Set("Authorization", "Bearer "+sess.AccessJwt)
	resp2, err := httpClient().Do(postReq)
	if err != nil {
		return Receipt{}, err
	}
	b2, _ := io.ReadAll(io.LimitReader(resp2.Body, 1<<20))
	_ = resp2.Body.Close()
	if resp2.StatusCode >= 300 {
		return Receipt{ChannelID: "bluesky", Mode: "manual", At: Now(),
			Detail: fmt.Sprintf("createRecord failed: %s", truncate(string(b2), 120))}, nil
	}
	var out struct {
		URI string `json:"uri"`
	}
	_ = json.Unmarshal(b2, &out)
	return Receipt{ChannelID: "bluesky", Mode: "sent", ExternalID: out.URI, At: Now()}, nil
}

func sendMastodon(ctx context.Context, text string) (Receipt, error) {
	base := strings.TrimRight(strings.TrimSpace(os.Getenv("MASTODON_BASE_URL")), "/")
	token := strings.TrimSpace(os.Getenv("MASTODON_ACCESS_TOKEN"))
	form := url.Values{"status": {text}}
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, base+"/api/v1/statuses", strings.NewReader(form.Encode()))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := httpClient().Do(req)
	if err != nil {
		return Receipt{}, err
	}
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	_ = resp.Body.Close()
	if resp.StatusCode >= 300 {
		return Receipt{ChannelID: "mastodon", Mode: "manual", At: Now(),
			Detail: fmt.Sprintf("status failed: %s", truncate(string(body), 120))}, nil
	}
	var out struct {
		ID  string `json:"id"`
		URL string `json:"url"`
	}
	_ = json.Unmarshal(body, &out)
	return Receipt{ChannelID: "mastodon", Mode: "sent", ExternalID: out.ID, At: Now(), Detail: out.URL}, nil
}

func sendResend(ctx context.Context, art Artifact) (Receipt, error) {
	key := strings.TrimSpace(os.Getenv("RESEND_API_KEY"))
	from := strings.TrimSpace(os.Getenv("RESEND_FROM"))
	to := strings.TrimSpace(os.Getenv("RESEND_TO"))
	if to == "" {
		return Receipt{ChannelID: "email", Mode: "manual", At: Now(),
			Detail: "RESEND_TO required for live send (opt-in)"}, nil
	}
	payload := fmt.Sprintf(`{"from":%q,"to":[%q],"subject":%q,"text":%q}`,
		from, to, art.Subject, art.Body)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.resend.com/emails", strings.NewReader(payload))
	req.Header.Set("Authorization", "Bearer "+key)
	req.Header.Set("Content-Type", "application/json")
	resp, err := httpClient().Do(req)
	if err != nil {
		return Receipt{}, err
	}
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	_ = resp.Body.Close()
	if resp.StatusCode >= 300 {
		return Receipt{ChannelID: "email", Mode: "manual", At: Now(),
			Detail: truncate(string(body), 160)}, nil
	}
	var out struct {
		ID string `json:"id"`
	}
	_ = json.Unmarshal(body, &out)
	return Receipt{ChannelID: "email", Mode: "sent", ExternalID: out.ID, At: Now()}, nil
}

func sendTwilio(ctx context.Context, to, text string) (Receipt, error) {
	sid := strings.TrimSpace(os.Getenv("TWILIO_ACCOUNT_SID"))
	token := strings.TrimSpace(os.Getenv("TWILIO_AUTH_TOKEN"))
	from := strings.TrimSpace(os.Getenv("TWILIO_FROM_NUMBER"))
	form := url.Values{"To": {to}, "From": {from}, "Body": {text}}
	endpoint := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", sid)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(sid+":"+token)))
	resp, err := httpClient().Do(req)
	if err != nil {
		return Receipt{}, err
	}
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	_ = resp.Body.Close()
	if resp.StatusCode >= 300 {
		return Receipt{ChannelID: "sms", Mode: "manual", At: Now(),
			Detail: truncate(string(body), 160)}, nil
	}
	var out struct {
		SID string `json:"sid"`
	}
	_ = json.Unmarshal(body, &out)
	return Receipt{ChannelID: "sms", Mode: "sent", ExternalID: out.SID, At: Now()}, nil
}
