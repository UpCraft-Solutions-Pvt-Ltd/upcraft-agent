package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// QdrantClient uses Qdrant's HTTP API directly to avoid extra backend deps.
type QdrantClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewQdrantClient(baseURL string) *QdrantClient {
	return &QdrantClient{
		BaseURL: strings.TrimRight(baseURL, "/"),
		HTTPClient: &http.Client{
			Timeout: 8 * time.Second,
		},
	}
}

func (q *QdrantClient) EnsureCollection(name string) error {
	getResp, err := q.HTTPClient.Get(fmt.Sprintf("%s/collections/%s", q.BaseURL, name))
	if err == nil {
		defer getResp.Body.Close()
		if getResp.StatusCode == http.StatusOK {
			return nil
		}
	}

	payload := map[string]interface{}{
		"vectors": map[string]interface{}{
			"size":     1536,
			"distance": "Cosine",
		},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/collections/%s", q.BaseURL, name), bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := q.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		data, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to create collection: status=%s body=%s", resp.Status, string(data))
	}
	return nil
}

func (q *QdrantClient) ScrollSkills(collection string) ([]SkillDefinition, error) {
	scrollBody := map[string]interface{}{
		"limit":        200,
		"with_payload": true,
		"with_vector":  false,
	}
	body, err := json.Marshal(scrollBody)
	if err != nil {
		return nil, err
	}

	resp, err := q.HTTPClient.Post(
		fmt.Sprintf("%s/collections/%s/points/scroll", q.BaseURL, collection),
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		data, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("qdrant scroll error: status=%s body=%s", resp.Status, string(data))
	}

	var decoded struct {
		Result struct {
			Points []struct {
				Payload map[string]interface{} `json:"payload"`
			} `json:"points"`
		} `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return nil, err
	}

	out := make([]SkillDefinition, 0, len(decoded.Result.Points))
	for _, p := range decoded.Result.Points {
		raw, err := json.Marshal(p.Payload)
		if err != nil {
			return nil, err
		}
		var skill SkillDefinition
		if err := json.Unmarshal(raw, &skill); err != nil {
			return nil, err
		}
		if strings.TrimSpace(skill.Name) == "" {
			continue
		}
		out = append(out, skill)
	}
	return out, nil
}

func (q *QdrantClient) UpsertSkill(collection string, skill SkillDefinition) error {
	// Placeholder zero-vector until embedding service is added.
	zeroVector := make([]float32, 1536)

	payload := map[string]interface{}{
		"points": []map[string]interface{}{
			{
				"id":      stableSkillID(skill.Name),
				"vector":  zeroVector,
				"payload": skill,
			},
		},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/collections/%s/points?wait=true", q.BaseURL, collection), bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := q.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		data, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("qdrant upsert error: status=%s body=%s", resp.Status, string(data))
	}
	return nil
}

func stableSkillID(name string) uint64 {
	var h uint64 = 1469598103934665603
	const prime uint64 = 1099511628211
	for _, b := range []byte(strings.ToLower(strings.TrimSpace(name))) {
		h ^= uint64(b)
		h *= prime
	}
	if h == 0 {
		return 1
	}
	return h
}
