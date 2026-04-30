# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Surge is a multi-connection TUI download manager written in Go with a browser extension. It has three operating modes: **TUI** (interactive Bubble Tea terminal UI), **Server** (headless HTTP API daemon), and **Connect** (remote TUI connecting to a running daemon).

## Build & Test Commands

```bash
# Build
go build -o surge .

# Run all tests
go test ./...

# Run tests for a specific package
go test ./internal/engine/concurrent
go test ./internal/tui

# Run a single test
go test ./internal/engine/concurrent -run TestConcurrentDownloader_SwitchOn429 -count=1
go test ./internal/download -run TestIntegration_PauseResume -count=1

# Run with race detection
go test -race ./...

# Extension (from extension/ directory)
cd extension && npm run dev           # dev mode (Chrome)
cd extension && npm run dev:firefox   # dev mode (Firefox)
cd extension && npm run build         # production build
cd extension && npm run check         # TypeScript type-check
cd extension && npm run lint          # ESLint
cd extension && npm run test          # Vitest
```

## Architecture

### Daemon / Single-Instance Model

Surge enforces a single running instance via a file lock (`internal/cmd/lock.go`). When you run `surge`, it starts an HTTP server on an available port (starting from 1700) and the TUI. Subsequent `surge add <url>` invocations detect the lock and forward the download to the running instance's HTTP API. The active port is written to a runtime file for discovery.

### Top-Level Package Layout

**`cmd/`** — Cobra CLI wiring, HTTP API routes, and global lifecycle orchestration. `root.go` is the central hub: it initializes the `WorkerPool`, `LocalDownloadService`, `LifecycleManager`, HTTP server, and TUI. Global variables (`GlobalPool`, `GlobalService`, `GlobalLifecycle`) are set in `PersistentPreRun`/`RunE`.

**`internal/core/`** — `DownloadService` interface with two implementations:
- `LocalDownloadService` — embedded engine, used when Surge runs locally
- `RemoteClient` — HTTP client that talks to a remote Surge daemon

**`internal/download/`** — `WorkerPool` manages concurrent download slots, `manager.go` (`TUIDownload`) orchestrates a single download's lifecycle (probe → engine → complete).

**`internal/engine/`** — Low-level HTTP downloaders:
- `concurrent/` — multi-connection chunked downloader with work stealing, slow-worker detection, and retry
- `single/` — single-connection fallback for servers that don't support Range requests
- `state/` — SQLite-backed persistence for pause/resume and download history
- `events/` — typed event messages (`ProgressMsg`, `DownloadCompleteMsg`, etc.) with SSE encode/decode for the HTTP API
- `types/` — shared type definitions (`DownloadConfig`, `DownloadStatus`, `DownloadState`)

**`internal/processing/`** — `LifecycleManager` orchestrates the flow from URL to download: probe (HEAD request for metadata, range support, redirects), duplicate detection, filename inference, then hands off to the engine via `AddFunc`.

**`internal/tui/`** — Bubble Tea TUI. `model.go` defines `RootModel` with ~20 UI states (Dashboard, Input, Settings, CategoryManager, etc.). The update/view pattern is split across files: `update.go` (main dispatch), `update_events.go`, `update_settings.go`, `update_modals.go`, etc. Views are similarly partitioned: `view.go`, `view_dashboard_*.go`, `view_settings.go`, `view_category.go`.

**`internal/config/`** — Typed `Settings` struct with categories (General, Network, Performance, Categories, Extension). Settings are loaded from JSON, validated on startup with rollback to defaults for invalid fields. `RuntimeConfig` converts user settings for the download engine.

**`extension/`** — Browser extension built with WXT + SolidJS + TypeScript. Background script intercepts browser downloads and forwards them to the local Surge HTTP API. Popup shows download list with status badges. Communicates via SSE for real-time updates.

### Event Flow

All components communicate via typed event messages sent on channels:
1. `WorkerPool` sends progress to `GlobalProgressCh`
2. `LocalDownloadService` fans out events to all registered listeners (TUI + SSE)
3. For remote connections, events are SSE-encoded over the `/events` endpoint
4. `LifecycleManager` listens to the event stream for pause/resume lifecycle management

### Key Patterns

- **File lock** at `config.GetRuntimeDir()/surge.lock` prevents multiple instances
- **Port discovery** at `config.GetRuntimeDir()/port` lets CLI and extension find the active server
- **Auth token** stored at `config.GetStateDir()/token`, Bearer-token protected HTTP API
- **Persistent state** via SQLite (`modernc.org/sqlite`, pure Go, no CGO) at `config.GetStateDir()/surge.db`
- **Startup integrity check** normalizes stale downloads (crashed mid-download → paused) and validates paused state files
- **Settings validation** runs on load; invalid fields are rolled back to defaults with warning messages shown in the TUI

### Build Flags

Version info is injected via `-ldflags`:
```
-X github.com/SurgeDM/Surge/cmd.Version={{.Version}}
-X github.com/SurgeDM/Surge/cmd.Commit={{.Commit}}
-X github.com/SurgeDM/Surge/cmd.BuildTime={{.Date}}
```

Builds use `CGO_ENABLED=0` for static binaries. Go 1.25+.
