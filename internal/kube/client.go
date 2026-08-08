// Package kube is a minimal in-cluster Kubernetes/OpenShift REST helper (no client-go).
package kube

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type Client struct {
	base   string
	token  string
	ns     string
	client *http.Client
}

func InCluster() (*Client, error) {
	host := os.Getenv("KUBERNETES_SERVICE_HOST")
	port := os.Getenv("KUBERNETES_SERVICE_PORT")
	if host == "" || port == "" {
		return nil, fmt.Errorf("not in cluster")
	}
	token, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/token")
	if err != nil {
		return nil, err
	}
	nsb, _ := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	ns := strings.TrimSpace(string(nsb))
	if ns == "" {
		ns = "dpo-system"
	}
	return &Client{
		base:  "https://" + host + ":" + port,
		token: strings.TrimSpace(string(token)),
		ns:    ns,
		client: &http.Client{
			Timeout: 15 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // cluster CA path varies; SA trust via skip for MVP
			},
		},
	}, nil
}

func (c *Client) Namespace() string { return c.ns }

func (c *Client) GetJSON(ctx context.Context, path string, out any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.base+path, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if resp.StatusCode >= 300 {
		return fmt.Errorf("kube %s: %s %s", path, resp.Status, truncate(string(body), 200))
	}
	return json.Unmarshal(body, out)
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}

// ListMeta is a thin items wrapper.
type ListMeta struct {
	Items []json.RawMessage `json:"items"`
}

func (c *Client) CountResource(ctx context.Context, apiPath string) (int, error) {
	var list ListMeta
	if err := c.GetJSON(ctx, apiPath, &list); err != nil {
		return 0, err
	}
	return len(list.Items), nil
}
