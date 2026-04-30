# Chinese (Simplified) Localization for Surge

## Overview

Add Simplified Chinese (简体中文 / zh-CN) language support to Surge — covering the Go backend (TUI, CLI, settings) and the browser extension — using a zero-dependency JSON dictionary approach. English text is used as the dictionary key so missing translations gracefully fall back to readable English, and upstream merges produce minimal conflicts.

## Design

### Core Principle: English Text as Key

```go
// Before — hardcoded English
"Loading..."

// After — T() looks up "Loading..." in zh-CN.json, returns "加载中..." or "Loading..."
T("Loading...")
```

- The dictionary key IS the English text — no abstract key names like `status_queued`
- No English baseline file needed — missing key → return key itself → readable English
- Minimal upstream merge conflict: source changes from `"Loading..."` to `T("Loading...")`, a single-line change per string

### Go i18n Package (`internal/i18n/`)

```
internal/i18n/
├── i18n.go       # T() function, JSON loading, language detection
├── zh-CN.json    # Chinese translations (only translated keys)
└── i18n_test.go
```

```go
func Init(lang string) error     // load locale file for given language
func T(key string) string        // translate key → zh-CN, fallback to key itself
func SetLanguage(lang string)    // hot-reload translations (for testing)
```

- `T()` is safe for concurrent use (read-only once loaded)
- Calls are **not** in hot paths (view render loops) — negligible overhead

### Language Selection

**Go TUI / CLI** — New field in `GeneralSettings`:

```go
// internal/config/settings.go
Language string `json:"language" ui_label:"Language" ui_desc:"Interface language (en, zh-CN). Restart required." ui_restart:"true"`
```

Wired in `initializeGlobalState()` (`cmd/root_startup.go`):

```go
i18n.Init(getSettings().General.Language)
```

Defaults to `"en"` when unset. Changing language in TUI settings requires restart (`ui_restart:"true"`).

**Browser Extension** — Language stored in `browser.storage.local` under key `language`:
- Default detected from `navigator.language` (startsWith "zh" → "zh-CN")
- Language can be changed in extension settings panel
- Extension language is independent of Surge daemon language

### String Migration By Area

#### 1. Settings Metadata (`internal/config/settings.go`)

- Do NOT modify struct tags (compile-time constant)
- Modify `GetSettingsMetadata()` to wrap returned `Label` and `Description` with `T()`
- ⚠️ Category map keys (e.g., `"General"`, `"Network"`) stay English — they're used for internal map lookups. Only individual setting `Label`/`Description` get translated.
- `CategoryOrder()` returns English labels; T() is applied at render time in `view_settings.go`
- `shortSettingsCategoryLabel()` gets a translation branch: `T("Gen")` → `"通用"`, etc.

```go
// GetSettingsMetadata reads struct tags then applies T()
meta.Label = T(meta.Label)         // "Default Download Dir" → 查找翻译
meta.Description = T(meta.Description)  // 同上
```

The zh-CN.json keys for settings are the English struct tag values themselves:
```json
{
  "Default Download Dir": "默认下载目录",
  "Default directory for new downloads. Leave empty to use current directory.": "新下载的默认目录。留空使用当前目录。",
  "Warn on Duplicate": "重复文件提醒",
  "Show warning when adding a download that already exists.": "添加已存在的下载时显示警告。"
}
```

#### 2. Status Labels (`internal/tui/components/status.go`)

Replace `statusMap` label values with `T()` calls at render time. The `Label()` method wraps:

```go
func (s DownloadStatus) Label() string {
    if info, ok := statusMap[s]; ok {
        return T(info.label)  // T("Queued") → "排队中"
    }
    return "Unknown"
}
```

#### 3. TUI View Strings (`internal/tui/view*.go`, `update*.go`)

Replace literal English strings with `T()`:
- Modal titles: `"Shutting Down"` → `T("Shutting Down")`
- Messages: `"Pausing downloads and saving resume state..."` → `T(...)`
- Labels: `"URL:"` → `T("URL:")`
- Empty states: `"No download selected"` → `T("No download selected")`
- Hints: `"esc: save/close  tab: next tab  enter: edit"` → `T(...)`

#### 4. CLI Command Descriptions (`cmd/*.go`)

```go
// Before
Short: "Add a new download to the running Surge instance"

// After
Short: T("Add a new download to the running Surge instance")
```

- `Use` field stays in English (commands are typed by user)

#### 5. CLI Output Messages (`cmd/*.go`)

```go
// Before
fmt.Printf("Successfully added %d downloads.\n", count)

// After — wrap format string, NOT the entire call
fmt.Printf(T("Successfully added %d downloads.\n"), count)
```

For `fmt.Println("text")` → `fmt.Println(T("text"))`.

#### 6. Headless Mode Output (`cmd/root_headless.go`)

Event messages like `"Started: %s [%s]"`, `"Completed: %s [%s] (in %s)"` — same pattern.

#### 7. Key Hint Labels (`internal/tui/keys.go`)

Help text in key binding definitions — wrap with `T()` at render time rather than at init time (since init runs before language is loaded).

#### 8. Browser Extension (`extension/lib/i18n/`)

```
extension/lib/i18n/
├── index.ts      # t() function + loading
└── zh-CN.json
```

```typescript
export function t(key: string): string;           // translate, fallback to key
export async function initLanguage(): Promise<void>;  // detect + load
```

Usage in JSX:
```tsx
// Before
<span>Intercept Downloads</span>

// After
<span>{t('Intercept Downloads')}</span>
```

Extension locale JSON:
```json
{
  "Intercept Downloads": "拦截下载",
  "Show Notifications": "显示通知",
  "Server URL": "服务器地址",
  "Auth Token": "认证令牌",
  "Active": "活动中",
  "History": "历史记录"
}
```

### JSON Dictionary Structure (`zh-CN.json`)

Only translated keys present. Missing keys fall back cleanly to the key itself (English):

```json
{
  "Queued": "排队中",
  "Downloading": "下载中",
  "Paused": "已暂停",
  "Completed": "已完成",
  "Error": "错误",
  "Unknown": "未知",
  "Loading...": "加载中...",
  "Add Download": "添加下载",
  "Shutting Down": "正在关闭",
  "Pausing downloads and saving resume state...": "正在暂停下载并保存断点续传状态...",
  "Please wait": "请稍候",
  ...
}
```

### Implementation Order

1. Create `internal/i18n/` package with `T()` + JSON loading + tests
2. Add `Language` field to `GeneralSettings` and wire `i18n.Init()` into startup
3. Migrate settings metadata labels/descriptions — grep `ui_label` and `ui_desc` tags
4. Migrate status labels (`components/status.go`)
5. Migrate TUI view strings (`view.go`, `view_*.go`, `update_*.go`)
6. Migrate CLI command descriptions (`cmd/*.go`)
7. Migrate CLI output messages (`cmd/*.go` — Println/Printf with user-facing text)
8. Migrate key hint labels
9. Create `extension/lib/i18n/` with `t()` function
10. Migrate extension strings
11. Add language selector to TUI settings UI (handled by `Language` field + `ui_restart`)
12. Add language selector to extension settings panel
13. Build and verify end-to-end

### Conflict Avoidance Strategy

| Scenario | Without i18n | With T() + JSON |
|----------|-------------|-----------------|
| Upstream changes "Loading..." to "Loading resources..." | Conflict if you translated it | No conflict — Go code changes from `T("Loading...")` to `T("Loading resources...")`, trivial auto-merge |
| Upstream adds a new UI string | Need to translate manually | Already in English, add to JSON when ready |
| Upstream rewords a help text | Conflict | One-line Go change, JSON key auto-resolves |
| Upstream deletes a UI string | Cleanup needed | `T()` call disappears, JSON key becomes unused (harmless) |

### Files Changed

**New:**
- `internal/i18n/i18n.go`
- `internal/i18n/zh-CN.json`
- `internal/i18n/i18n_test.go`
- `extension/lib/i18n/index.ts`
- `extension/lib/i18n/zh-CN.json`

**Modified (Go):**
- `internal/config/settings.go` — add `Language` field
- `cmd/root_startup.go` — wire `i18n.Init()`
- `internal/config/categories.go` — maybe
- `internal/config/settings_metadata.go` — maybe (GetSettingsMetadata call T())
- `internal/tui/components/status.go` — Label() wraps T()
- `internal/tui/view.go`
- `internal/tui/view_settings.go`
- `internal/tui/view_category.go`
- `internal/tui/view_dashboard_header.go`
- `internal/tui/view_dashboard_list.go`
- `internal/tui/view_dashboard_details.go`
- `internal/tui/view_dashboard_graph.go`
- `internal/tui/view_dashboard_log.go`
- `internal/tui/view_dashboard_chunkmap.go`
- `internal/tui/update_modals.go`
- `internal/tui/update_settings.go`
- `internal/tui/update_dashboard.go`
- `internal/tui/keys.go`
- `internal/tui/config.go`
- `internal/tui/constants.go`
- `internal/tui/helpers.go`
- `cmd/root.go`
- `cmd/add.go`
- `cmd/ls.go`
- `cmd/rm.go`
- `cmd/pause.go`
- `cmd/resume.go`
- `cmd/refresh.go`
- `cmd/server.go`
- `cmd/shutdown.go`
- `cmd/connect.go`
- `cmd/root_headless.go`
- `cmd/bugreport.go`

**Modified (Extension):**
- `extension/entrypoints/popup/App.tsx`
- `extension/entrypoints/popup/components/DownloadList.tsx`
- `extension/entrypoints/popup/components/DownloadItem.tsx`
- `extension/entrypoints/popup/components/SettingsView.tsx`
- `extension/entrypoints/popup/components/StatusBadge.tsx`
- `extension/entrypoints/popup/components/DuplicateModal.tsx`
- `extension/entrypoints/popup/components/ViewSwitch.tsx`
- `extension/lib/storage.ts`
