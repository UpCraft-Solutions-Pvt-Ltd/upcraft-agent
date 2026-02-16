package skills

import "context"

// Browser defines deterministic browser actions implemented by platform plugins.
type Browser interface {
	Visit(ctx context.Context, url string) error
	Search(ctx context.Context, query string) error
}
