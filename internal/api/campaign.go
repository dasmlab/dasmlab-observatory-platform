package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/dasmlab/dasmlab-observatory-platform/internal/campaign"
	collector "github.com/dasmlab/dasmlab-observatory-platform/platform/collector-sdk"
)

func (s *Server) channels(w http.ResponseWriter, _ *http.Request) {
	if s.d.Campaigns == nil {
		http.Error(w, "campaigns unavailable", 503)
		return
	}
	writeJSON(w, http.StatusOK, s.d.Campaigns.Registry().ChannelInfos())
}

func (s *Server) campaignsList(w http.ResponseWriter, _ *http.Request) {
	if s.d.Campaigns == nil {
		http.Error(w, "campaigns unavailable", 503)
		return
	}
	writeJSON(w, http.StatusOK, s.d.Campaigns.List())
}

func (s *Server) campaignCreate(w http.ResponseWriter, r *http.Request) {
	if s.d.Campaigns == nil {
		http.Error(w, "campaigns unavailable", 503)
		return
	}
	var body campaign.Campaign
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	c, err := s.d.Campaigns.Create(body)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	writeJSON(w, http.StatusCreated, c)
}

func (s *Server) campaignOne(w http.ResponseWriter, r *http.Request) {
	if s.d.Campaigns == nil {
		http.Error(w, "campaigns unavailable", 503)
		return
	}
	c, err := s.d.Campaigns.Get(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), 404)
		return
	}
	writeJSON(w, http.StatusOK, c)
}

func (s *Server) campaignRender(w http.ResponseWriter, r *http.Request) {
	if s.d.Campaigns == nil {
		http.Error(w, "campaigns unavailable", 503)
		return
	}
	id := chi.URLParam(r, "id")
	c, err := s.d.Campaigns.Render(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	s.emitCampaignEvents(c, "campaign_rendered", "dry_run")
	writeJSON(w, http.StatusOK, c)
}

func (s *Server) campaignArm(w http.ResponseWriter, r *http.Request) {
	if s.d.Campaigns == nil {
		http.Error(w, "campaigns unavailable", 503)
		return
	}
	c, err := s.d.Campaigns.Arm(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), 409)
		return
	}
	writeJSON(w, http.StatusOK, c)
}

func (s *Server) campaignSend(w http.ResponseWriter, r *http.Request) {
	if s.d.Campaigns == nil {
		http.Error(w, "campaigns unavailable", 503)
		return
	}
	c, err := s.d.Campaigns.Send(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), 409)
		return
	}
	s.emitCampaignEvents(c, "content_published", "")
	writeJSON(w, http.StatusOK, c)
}

func (s *Server) emitCampaignEvents(c *campaign.Campaign, metric, forceMode string) {
	if s.d.Store == nil || c == nil {
		return
	}
	var events []collector.Event
	for _, ch := range c.Channels {
		mode := forceMode
		if mode == "" && ch.Receipt != nil {
			mode = ch.Receipt.Mode
		}
		if mode == "" {
			mode = "dry_run"
		}
		events = append(events, collector.Event{
			SchemaVersion: "1",
			Collector:     "campaign",
			Type:          "campaign",
			Timestamp:     time.Now().UTC(),
			Tenant:        s.d.Tenant,
			Entity:        c.ID,
			Metric:        metric,
			Value:         1,
			Unit:          "count",
			Dims: map[string]string{
				"campaign": c.ID,
				"channel":  ch.ChannelID,
				"mode":     mode,
			},
		})
	}
	_ = s.d.Store.InsertEvents(events)
}
