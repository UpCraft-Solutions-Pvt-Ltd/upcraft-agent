package desktop

import (
	"context"

	"github.com/UpCraft-Solutions-Pvt-Ltd/upcraft-agent/core/plugins/spotify"
)

// MusicPlugin delegates desktop playback to Spotify Web API.
type MusicPlugin struct {
	player *spotify.PlayerClient
}

func NewMusicPluginFromEnv() (*MusicPlugin, error) {
	player, err := spotify.NewPlayerClientFromEnv()
	if err != nil {
		return nil, err
	}
	return &MusicPlugin{player: player}, nil
}

func (p *MusicPlugin) Play(ctx context.Context, query string) error {
	return p.player.Play(ctx, query)
}

func (p *MusicPlugin) Pause(ctx context.Context) error {
	return p.player.Pause(ctx)
}

func (p *MusicPlugin) Resume(ctx context.Context) error {
	return p.player.Resume(ctx)
}

func (p *MusicPlugin) Next(ctx context.Context) error {
	return p.player.Next(ctx)
}
