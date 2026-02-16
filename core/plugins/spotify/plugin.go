package spotify

import "context"

// MusicPlugin is a contributor-friendly wrapper around PlayerClient.
type MusicPlugin struct {
	client *PlayerClient
}

func NewMusicPluginFromEnv() (*MusicPlugin, error) {
	client, err := NewPlayerClientFromEnv()
	if err != nil {
		return nil, err
	}
	return &MusicPlugin{client: client}, nil
}

func (p *MusicPlugin) Play(query string) error {
	return p.client.Play(context.Background(), query)
}

func (p *MusicPlugin) Pause() error {
	return p.client.Pause(context.Background())
}

func (p *MusicPlugin) Resume() error {
	return p.client.Resume(context.Background())
}

func (p *MusicPlugin) Next() error {
	return p.client.Next(context.Background())
}
