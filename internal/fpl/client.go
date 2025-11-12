package fpl

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const defaultBaseURL = "https://fantasy.premierleague.com/api"

// Client wraps HTTP access to the FPL API.
type Client struct {
	httpClient *http.Client
	baseURL    string
	cacheTTL   time.Duration

	bootstrapCache struct {
		data   *BootstrapStatic
		expiry time.Time
	}
}

// NewClient constructs a Client with sane defaults.
func NewClient(httpClient *http.Client, cacheTTL time.Duration) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 15 * time.Second}
	}
	return &Client{
		httpClient: httpClient,
		baseURL:    defaultBaseURL,
		cacheTTL:   cacheTTL,
	}
}

// Bootstrap fetches /bootstrap-static/, optionally serving from cache.
func (c *Client) Bootstrap(ctx context.Context) (*BootstrapStatic, error) {
	if c.cacheTTL > 0 && c.bootstrapCache.data != nil && time.Now().Before(c.bootstrapCache.expiry) {
		return c.bootstrapCache.data, nil
	}

	var payload BootstrapStatic
	if err := c.get(ctx, "/bootstrap-static/", &payload); err != nil {
		return nil, err
	}

	if c.cacheTTL > 0 {
		c.bootstrapCache.data = &payload
		c.bootstrapCache.expiry = time.Now().Add(c.cacheTTL)
	}

	return &payload, nil
}

// PlayerSummary fetches /element-summary/{id}/ for a player.
func (c *Client) PlayerSummary(ctx context.Context, id int) (*PlayerSummary, error) {
	var payload PlayerSummary
	if err := c.get(ctx, fmt.Sprintf("/element-summary/%d/", id), &payload); err != nil {
		return nil, err
	}
	return &payload, nil
}

func (c *Client) get(ctx context.Context, path string, target any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+path, nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return fmt.Errorf("fpl api %s: %s: %s", path, resp.Status, string(body))
	}

	return json.NewDecoder(resp.Body).Decode(target)
}
