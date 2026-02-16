package spotify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	spotifyAPIBase = "https://api.spotify.com/v1"
)

// PlayerClient executes real Spotify Web API requests.
type PlayerClient struct {
	AccessToken string
	DeviceID    string
	HTTPClient  *http.Client
}

func NewPlayerClientFromEnv() (*PlayerClient, error) {
	token := strings.TrimSpace(os.Getenv("SPOTIFY_ACCESS_TOKEN"))
	if token == "" {
		return nil, fmt.Errorf("SPOTIFY_ACCESS_TOKEN is required")
	}
	return &PlayerClient{
		AccessToken: token,
		DeviceID:    strings.TrimSpace(os.Getenv("SPOTIFY_DEVICE_ID")),
		HTTPClient:  &http.Client{Timeout: 30 * time.Second},
	}, nil
}

func (c *PlayerClient) Play(ctx context.Context, query string) error {
	query = strings.TrimSpace(query)
	if query == "" {
		return fmt.Errorf("query is required")
	}

	trackURI, err := c.searchTopTrackURI(ctx, query)
	if err != nil {
		return err
	}

	payload := map[string]interface{}{"uris": []string{trackURI}}
	_, err = c.request(ctx, http.MethodPut, "/me/player/play", payload)
	return err
}

func (c *PlayerClient) Pause(ctx context.Context) error {
	_, err := c.request(ctx, http.MethodPut, "/me/player/pause", nil)
	return err
}

func (c *PlayerClient) Resume(ctx context.Context) error {
	_, err := c.request(ctx, http.MethodPut, "/me/player/play", nil)
	return err
}

func (c *PlayerClient) Next(ctx context.Context) error {
	_, err := c.request(ctx, http.MethodPost, "/me/player/next", nil)
	return err
}

func (c *PlayerClient) request(ctx context.Context, method, endpoint string, payload interface{}) ([]byte, error) {
	if c.HTTPClient == nil {
		c.HTTPClient = &http.Client{Timeout: 30 * time.Second}
	}

	u, err := url.Parse(spotifyAPIBase + endpoint)
	if err != nil {
		return nil, err
	}
	if c.DeviceID != "" {
		q := u.Query()
		q.Set("device_id", c.DeviceID)
		u.RawQuery = q.Encode()
	}

	var body io.Reader
	if payload != nil {
		encoded, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("marshal payload: %w", err)
		}
		body = bytes.NewReader(encoded)
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), body)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send spotify request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read spotify response: %w", err)
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return respBody, nil
	}

	if len(respBody) == 0 {
		return nil, fmt.Errorf("spotify request failed status=%d", resp.StatusCode)
	}
	return nil, fmt.Errorf("spotify request failed status=%d body=%s", resp.StatusCode, string(respBody))
}

func (c *PlayerClient) searchTopTrackURI(ctx context.Context, query string) (string, error) {
	u, _ := url.Parse(spotifyAPIBase + "/search")
	q := u.Query()
	q.Set("q", query)
	q.Set("type", "track")
	q.Set("limit", "1")
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return "", fmt.Errorf("create search request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.AccessToken)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("send search request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read search response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("spotify search failed status=%d body=%s", resp.StatusCode, string(body))
	}

	var decoded struct {
		Tracks struct {
			Items []struct {
				URI string `json:"uri"`
			} `json:"items"`
		} `json:"tracks"`
	}
	if err := json.Unmarshal(body, &decoded); err != nil {
		return "", fmt.Errorf("decode search response: %w", err)
	}
	if len(decoded.Tracks.Items) == 0 || strings.TrimSpace(decoded.Tracks.Items[0].URI) == "" {
		return "", fmt.Errorf("no track found for query %q", query)
	}

	return decoded.Tracks.Items[0].URI, nil
}
