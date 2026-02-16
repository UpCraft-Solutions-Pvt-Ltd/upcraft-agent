package skills

import "context"

// MusicPlayer defines deterministic controls for music playback providers.
type MusicPlayer interface {
	Play(ctx context.Context, query string) error
	Pause(ctx context.Context) error
	Resume(ctx context.Context) error
	Next(ctx context.Context) error
}
