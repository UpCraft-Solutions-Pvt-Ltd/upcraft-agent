package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

// SkillDefinition is served to the agent as the cloud skill contract.
type SkillDefinition struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	JSONSchema  string `json:"json_schema"`
}

var qdrant *QdrantClient

func main() {
	qdrantURL := strings.TrimSpace(os.Getenv("QDRANT_URL"))
	if qdrantURL == "" {
		qdrantURL = "http://localhost:6333"
	}
	qdrant = NewQdrantClient(qdrantURL)

	if err := qdrant.EnsureCollection("skills"); err != nil {
		log.Printf("Warning: failed to ensure Qdrant collection: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/sync-skills", handleSyncSkills)
	mux.HandleFunc("/admin/ingest", handleIngestSkill)
	mux.HandleFunc("/health", handleHealth)

	log.Printf("UpCraft Backend running on :%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}

func handleHealth(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}

func handleSyncSkills(w http.ResponseWriter, _ *http.Request) {
	skills, err := qdrant.ScrollSkills("skills")
	if err != nil {
		http.Error(w, "DB Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(skills); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleIngestSkill(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}

	var skill SkillDefinition
	if err := json.Unmarshal(body, &skill); err != nil {
		http.Error(w, "Bad JSON", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(skill.Name) == "" {
		http.Error(w, "Field 'name' is required", http.StatusBadRequest)
		return
	}

	if err := qdrant.UpsertSkill("skills", skill); err != nil {
		http.Error(w, fmt.Sprintf("Failed to upsert: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Skill Ingested"))
}
