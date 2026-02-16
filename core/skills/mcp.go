package skills

import (
	"context"
	"encoding/json"
)

// MCPClient defines a client that can invoke external MCP tool servers.
type MCPClient interface {
	Call(ctx context.Context, tool string, input json.RawMessage) (json.RawMessage, error)
}
