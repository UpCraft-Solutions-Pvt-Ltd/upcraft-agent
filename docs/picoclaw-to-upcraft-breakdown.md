# PicoClaw -> UpCraft Complete Breakdown

This document records a complete migration-style breakdown of the PicoClaw repository into the UpCraft architected folder buckets.

## Verification
- Source repository: `d:\ClayBot\picoclaw`
- Destination repository: `d:\ClayBot\upcraft-agent`
- Source file count: **111**
- Mapped file count: **111**
- Missing files after copy: **0**

## Mapping Rules Used
- `cmd/picoclaw/*` -> `core/cmd/picoclaw/*`
- `pkg/agent/*` -> `core/engine/picoclaw/agent/*`
- `pkg/auth/*` -> `core/engine/picoclaw/auth/*`
- `pkg/bus/*` -> `core/engine/picoclaw/bus/*`
- `pkg/config/*` -> `core/engine/picoclaw/config/*`
- `pkg/constants/*` -> `core/engine/picoclaw/constants/*`
- `pkg/cron/*` -> `core/engine/picoclaw/cron/*`
- `pkg/heartbeat/*` -> `core/engine/picoclaw/heartbeat/*`
- `pkg/logger/*` -> `core/engine/picoclaw/logger/*`
- `pkg/providers/*` -> `core/engine/picoclaw/providers/*`
- `pkg/tools/*` -> `core/engine/picoclaw/tools/*`
- `pkg/utils/*` -> `core/engine/picoclaw/utils/*`
- `pkg/migrate/*` -> `core/memory/picoclaw/migrate/*`
- `pkg/session/*` -> `core/memory/picoclaw/session/*`
- `pkg/state/*` -> `core/memory/picoclaw/state/*`
- `pkg/skills/*` -> `core/skills/picoclaw/runtime/*`
- `skills/*` -> `core/skills/picoclaw/catalog/*`
- `pkg/channels/*` -> `core/plugins/picoclaw/channels/*`
- `pkg/voice/*` -> `core/plugins/picoclaw/voice/*`
- `config/*` -> `core/engine/picoclaw/appconfig/*`
- `assets/*` -> `app/shared/picoclaw-assets/*`
- `.github/*` -> `docs/reference/picoclaw-github/*`
- root meta files (`README*`, `Dockerfile`, `go.mod`, etc.) -> `docs/reference/picoclaw-root/*`

## Full Per-File Manifest
- `docs/picoclaw-file-manifest.tsv`

## Notes
- This is a **complete source breakdown copy** with no mocked placeholders.
- The copied PicoClaw code is preserved under namespaced `picoclaw` subfolders to avoid clobbering UpCraft-native modules while retaining full reference fidelity.