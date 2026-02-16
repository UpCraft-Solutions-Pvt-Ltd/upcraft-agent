package engine

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/UpCraft-Solutions-Pvt-Ltd/upcraft-agent/core/memory"
)

// Agent is the lightweight runtime shell for syncing cloud skill definitions
// before entering the deterministic ReAct loop.
type Agent struct {
	RAG            *memory.RAGClient
	mu             sync.Mutex
	lastScreenJSON string
}

func NewAgent() *Agent {
	baseURL := os.Getenv("UPCRAFT_RAG_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}
	return &Agent{RAG: memory.NewRAGClient(baseURL)}
}

func (a *Agent) Start() {
	fmt.Println("UpCraft Agent Starting...")

	skills, err := a.RAG.FetchSkills()
	if err != nil {
		fmt.Printf("Warning: Could not sync skills (%v). Using cached defaults.\n", err)
		return
	}

	fmt.Printf("Synced %d skills from Cloud RAG.\n", len(skills))
}

// HandleScreenInput accepts simplified screen-state JSON from mobile/desktop UI layers
// and returns a deterministic JSON command envelope that the caller can execute.
func (a *Agent) HandleScreenInput(inputJSON string) string {
	a.mu.Lock()
	a.lastScreenJSON = inputJSON
	a.mu.Unlock()

	trimmed := strings.TrimSpace(inputJSON)
	if trimmed == "" {
		return `{"action":"NOOP","reason":"empty_input"}`
	}

	var payload map[string]interface{}
	if err := json.Unmarshal([]byte(trimmed), &payload); err != nil {
		return `{"action":"NOOP","reason":"invalid_json"}`
	}

	// Minimal deterministic behavior: if a visible label "Play" exists,
	// return a click command targeting that text.
	if containsText(payload, "Play") {
		return `{"action":"CLICK","text":"Play"}`
	}

	return `{"action":"NOOP","reason":"no_target"}`
}

func containsText(v interface{}, target string) bool {
	switch t := v.(type) {
	case map[string]interface{}:
		for k, value := range t {
			if strings.EqualFold(k, "text") {
				if s, ok := value.(string); ok && strings.EqualFold(strings.TrimSpace(s), target) {
					return true
				}
			}
			if containsText(value, target) {
				return true
			}
		}
	case []interface{}:
		for _, item := range t {
			if containsText(item, target) {
				return true
			}
		}
	}
	return false
}
