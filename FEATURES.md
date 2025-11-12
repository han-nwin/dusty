## 🧩 Core MVP Features (Needed)

These are essential for a usable, safe, minimal version 1:

### 1. **Scan & Display**

- Recursively scan known safe directories:

  - `~/Library/Caches/**`
  - `~/Library/Logs/**`
  - `~/Library/Developer/Xcode/{DerivedData,Archives}`
  - `~/.npm/_cacache`, `~/.cache/yarn`, `~/Library/Caches/pip`

- Compute:

  - Folder size (recursive)
  - File count
  - Last modified (newest + oldest)

- Display results in a TUI table (using Bubble Tea + Bubbles list).

### 2. **Selection & Filtering**

- Toggle entries (spacebar) to mark for cleaning.
- Sort by size, name, or last modified.
- Optional filter/search (`/chrome`, `/xcode`).

### 3. **Clean Action**

- On confirm:

  - Move to Trash via AppleScript (`osascript -e 'tell app "Finder" to delete POSIX file ...'`).
  - Or hard delete for temporary dirs (`os.RemoveAll()`).

- Dry-run mode by default (just print would-delete list).

### 4. **Progress Feedback**

- Progress bar or spinner (Bubble Tea progress component).
- Display total size being removed.
- Log summary: “Removed 4.2 GB across 6 targets”.

### 5. **Undo / Manifest**

- Save a JSON log in `~/.cleanlite/undo/`:

  ```json
  {
    "timestamp": "2025-11-11T18:45:00",
    "targets": ["~/Library/Caches", "~/Library/Logs"],
    "files": 51234,
    "bytes": 4512123134
  }
  ```

- Optional “Open Trash” or “Reveal manifest”.

### 6. **Keyboard UX**

| Key   | Action           |
| ----- | ---------------- |
| ↑↓    | Navigate         |
| Space | Toggle selection |
| Enter | Details          |
| C     | Clean            |
| R     | Rescan           |
| /     | Filter           |
| ?     | Help             |
| Q     | Quit             |

### 7. **Safety**

- Only allow known allow-listed paths.
- Ignore hidden/system dirs.
- Never escalate to `sudo`.
- Confirm before deletion with total GB display.

---

## 🌱 Nice-to-Have Features (v2+)

### Visual / UX

- Size color-coding (green < 500 MB, yellow < 2 GB, red > 2 GB).
- ASCII graph or bar showing relative sizes.
- Human-friendly “last cleaned X days ago”.

### Presets

- **Developer Mode:** includes Xcode, npm, yarn.
- **Browser Mode:** includes Safari, Chrome, Edge caches.
- **Full Scan:** everything in allowlist.

### Stats / History

- Track total reclaimed space over time.
- Show weekly cleanup trends.
- “You’ve saved 23 GB this month!” summary.

### Automation

- CLI flag `--auto` to clean everything silently.
- LaunchAgent integration for scheduled scans.
- Config file (`~/.cleanlite/config.yml`) for excludes, thresholds.

### Smart Scanning

- Skip files newer than 24 hours.
- Detect running apps and skip their cache (Safari, Chrome).
- Parallel scan with progress updates.

### Extensibility

- Plugin architecture for new “modules”:

  - `brew` (run `brew cleanup --dry-run`)
  - `docker` (show `docker system df`)
  - `npm`/`yarn` (show cache size, clean flag)

### Polish

- Help screen with markdown (rendered via Glamour).
- JSON/CSV export of scan results.
- Configurable trash vs delete.
- Light/dark theme toggle.

---

## 🧱 Suggested MVP Milestones

| Week | Focus                                      | Outcome                 |
| ---- | ------------------------------------------ | ----------------------- |
| 1    | Core scanner + TUI list + dry-run          | Show sizes safely       |
| 2    | Selectable items + Trash action + progress | Functional cleaner      |
| 3    | Undo manifest + presets + UX polish        | Safe, pleasant CLI tool |
