# Reference Analysis: PicoClaw + OpenClaw

Date: 2026-02-13
Cloned refs:
- d:/ClayBot/picoclaw (https://github.com/sipeed/picoclaw)
- d:/ClayBot/openclaw (https://github.com/openclaw/openclaw)

## Why these references matter for UpCraft
- PicoClaw gives a compact Go implementation of LLM + tool-call loops and provider abstraction.
- OpenClaw gives production-grade patterns for plugin loading, skills snapshots, session handling, and multi-surface orchestration.

## PicoClaw findings (Go)

### 1) Deterministic tool registry exists and is reusable
- `picoclaw/pkg/tools/registry.go` provides a mutex-protected registry with:
  - tool registration
  - execution with contextual metadata (channel/chat)
  - schema conversion to provider-compatible tool definitions
- Relevance to UpCraft:
  - Keep `core/skills` as interfaces and provide a registry in `core/engine` that maps deterministic JSON action -> concrete plugin call.

### 2) Core run loop is explicit LLM->tool->LLM iteration
- `picoclaw/pkg/tools/toolloop.go` implements the exact loop:
  1. Build tool definitions
  2. Call provider
  3. If tool calls returned, execute each tool
  4. Append tool results
  5. Iterate up to max
- Relevance to UpCraft:
  - This is the closest baseline for `picoclaw-core` behavior we want.
  - We should keep strict iteration caps and typed request/response structs in Go.

### 3) Agent loop separates runtime concerns from tool loop
- `picoclaw/pkg/agent/loop.go` handles:
  - message bus IO
  - session/state tracking
  - tool registry setup
  - defers iterative logic to tool loop
- Relevance to UpCraft:
  - Keep transport concerns (Android/Desktop events) out of the engine loop implementation.

### 4) Provider abstraction is practical and OpenRouter-ready
- `picoclaw/pkg/providers/http_provider.go` shows:
  - generic `/chat/completions` flow
  - tool definitions in request body
  - model/provider routing
  - OpenRouter fallback URL (`https://openrouter.ai/api/v1`)
- Relevance to UpCraft:
  - Matches your requirement to hardcode OpenRouter base and read key from env/config.

### 5) Gap vs UpCraft philosophy
- PicoClaw still allows general tools (including shell/exec).
- UpCraft requirement is stricter: no script-writing action path.
- Design implication:
  - UpCraft should only expose whitelisted skill actions with typed payload schemas.

## OpenClaw findings (TypeScript)

### 1) Mature plugin/registry model
- `openclaw/src/plugins/registry.ts` + `openclaw/src/plugins/tools.ts`:
  - plugin discovery and registration
  - conflict detection for tool names
  - optional tool gating via allowlist
  - diagnostics collection
- Relevance to UpCraft:
  - Add plugin conflict checks in Go registry early.
  - Add plugin metadata + diagnostics endpoint for observability.

### 2) Skills snapshot pattern is strong
- `openclaw/src/commands/agent.ts` builds and persists `skillsSnapshot` per session.
- Relevance to UpCraft:
  - For your cloud RAG sync, cache the skill-definition snapshot version locally.
  - Session should pin a snapshot to avoid mid-run schema drift.

### 3) Execution runner is production hardened
- `openclaw/src/agents/pi-embedded-runner/run.ts`
- `openclaw/src/agents/pi-embedded-runner/run/attempt.ts`
  - lane queues
  - auth profile failover
  - context-window guards
  - tool sanitization for provider quirks
  - session repair + transcript safety
- Relevance to UpCraft:
  - UpCraft can start lean, but should copy these hardening concepts:
    - bounded queueing
    - context budget checks
    - per-provider adaptation layer

### 4) Plugin runtime state is explicit
- `openclaw/src/plugins/runtime.ts` manages active registry state centrally.
- Relevance to UpCraft:
  - Use a single source of truth for active plugin instances in Go, especially with mobile lifecycle restarts.

## Direct mapping into UpCraft architecture

### Keep
- PicoClaw-style Go tool loop and provider abstraction.
- OpenClaw-style plugin registry diagnostics and skills snapshot/versioning.

### Reject
- Any freeform code-execution path (shell/script generation by LLM).
- Any tool contract that is not schema-validated before execution.

### Build next in UpCraft
1. `core/engine/registry.go`
- Register skill implementations by `skill` + `action`.
- Detect conflicts and expose deterministic dispatch.

2. `core/engine/loop.go`
- Request OpenRouter with strict JSON-mode instruction.
- Parse into typed command struct.
- Validate command against registry and schema.
- Execute plugin and return normalized result envelope.

3. `core/engine/openrouter.go`
- Hardcode `https://openrouter.ai/api/v1`.
- Read API key from env/local config.
- Support free-tier model IDs and timeout controls.

4. `backend/main.go`
- Real Qdrant client initialization.
- `/sync-skills` should return skill definitions + schema version.

## Risks discovered now
- OpenClaw is much broader than needed; copying it directly will overcomplicate mobile-first constraints.
- PicoClaw has useful loop patterns but looser tool safety than UpCraft requires.
- Therefore UpCraft should be a strict subset: PicoClaw loop simplicity + OpenClaw plugin discipline.

## Conclusion
Both repositories are now available locally for ongoing reference. The best architectural path is a constrained Go-first deterministic execution engine that borrows proven loop and registry patterns but enforces stricter typed-action boundaries than either upstream.