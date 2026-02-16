//go:build android

package android

import (
	"context"
	"encoding/json"
	"sync"
)

// MusicPlugin records deterministic playback intents for Android UI execution.
type MusicPlugin struct {
	mu          sync.Mutex
	lastCommand string
}

func NewMusicPlugin() *MusicPlugin {
	return &MusicPlugin{}
}

func (p *MusicPlugin) Play(_ context.Context, query string) error {
	p.setCommand("PLAY", query)
	return nil
}

func (p *MusicPlugin) Pause(_ context.Context) error {
	p.setCommand("PAUSE", "")
	return nil
}

func (p *MusicPlugin) Resume(_ context.Context) error {
	p.setCommand("RESUME", "")
	return nil
}

func (p *MusicPlugin) Next(_ context.Context) error {
	p.setCommand("NEXT", "")
	return nil
}

// ConsumeLastCommand returns and clears the latest Android action JSON.
func (p *MusicPlugin) ConsumeLastCommand() string {
	p.mu.Lock()
	defer p.mu.Unlock()
	out := p.lastCommand
	p.lastCommand = ""
	return out
}

func (p *MusicPlugin) setCommand(action, query string) {
	cmd := map[string]string{
		"tool":   "MusicPlayer",
		"action": action,
	}
	if query != "" {
		cmd["query"] = query
	}
	encoded, _ := json.Marshal(cmd)

	p.mu.Lock()
	p.lastCommand = string(encoded)
	p.mu.Unlock()
}
