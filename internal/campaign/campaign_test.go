package campaign_test

import (
	"context"
	"testing"

	"github.com/dasmlab/dasmlab-observatory-platform/internal/campaign"
)

func TestWave1DryRunAllChannels(t *testing.T) {
	svc := campaign.NewService("dasmlab.org", "")
	c, err := svc.Render(context.Background(), "dasmlab-2.0-launch")
	if err != nil {
		t.Fatal(err)
	}
	if c.Status != "dry_run" {
		t.Fatalf("status=%s", c.Status)
	}
	if len(c.Channels) < 10 {
		t.Fatalf("expected ≥10 channels, got %d", len(c.Channels))
	}
	for _, ch := range c.Channels {
		if len(ch.Artifacts) == 0 {
			t.Fatalf("channel %s missing artifacts", ch.ChannelID)
		}
		if ch.Artifacts[0].Body == "" {
			t.Fatalf("channel %s empty body", ch.ChannelID)
		}
	}
}

func TestChannelInfos(t *testing.T) {
	infos := campaign.NewRegistry().ChannelInfos()
	if len(infos) < 10 {
		t.Fatalf("got %d", len(infos))
	}
	var webOK bool
	for _, i := range infos {
		if i.ID == "web_slash" && i.LiveStatus == "ready" {
			webOK = true
		}
	}
	if !webOK {
		t.Fatal("web_slash should be ready")
	}
}
