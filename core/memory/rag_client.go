package memory

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type RAGClient struct {
	BaseURL string
	Client  *http.Client
}

type RemoteSkill struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	JSONSchema  string `json:"json_schema"`
}

func NewRAGClient(baseURL string) *RAGClient {
	return &RAGClient{
		BaseURL: strings.TrimRight(baseURL, "/"),
		Client:  &http.Client{Timeout: 5 * time.Second},
	}
}

func (c *RAGClient) FetchSkills() ([]RemoteSkill, error) {
	resp, err := c.Client.Get(c.BaseURL + "/sync-skills")
	if err != nil {
		return nil, fmt.Errorf("failed to sync skills: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("backend returned status: %d", resp.StatusCode)
	}

	var skills []RemoteSkill
	if err := json.NewDecoder(resp.Body).Decode(&skills); err != nil {
		return nil, fmt.Errorf("failed to decode skills: %w", err)
	}

	return skills, nil
}
