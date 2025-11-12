# Dusty - TODOs

A CleanMyMac-style TUI for macOS.

## 🎯 MVP (v1.0)

### Core Scanner
- [ ] Scan allowlisted paths (~/Library/Caches, ~/Library/Logs, Xcode DerivedData, npm/yarn cache)
- [ ] Calculate size, file count, last modified date
- [ ] Return structured results

### TUI
- [ ] Display scan results in table (Bubble Tea + Bubbles)
- [ ] Navigate with ↑↓, toggle with Space
- [ ] Sort by size/name/date, filter with /
- [ ] Show details view (Enter)

### Clean
- [ ] Move to Trash via AppleScript
- [ ] Confirmation dialog with total size
- [ ] Progress bar during deletion
- [ ] Save JSON manifest to ~/.dusty/undo/

### Safety
- [ ] Enforce allowlist-only paths
- [ ] No sudo, confirm before delete
- [ ] Show summary after cleaning

### Keyboard Shortcuts
- [ ] ↑↓ navigate, Space toggle, Enter details
- [ ] C clean, R rescan, / filter, ? help, Q quit

---

## 🌱 v2+ (Future)

- [ ] Color-coded sizes, ASCII graphs
- [ ] Presets (Developer/Browser/Full modes)
- [ ] Stats & history tracking
- [ ] `--auto` flag and LaunchAgent scheduling
- [ ] Smart scanning (skip running apps, parallel scan)
- [ ] Plugins (brew cleanup, docker system df)
- [ ] Help screen with Glamour, JSON/CSV export

---

## 📦 Release

- [ ] Tests for scanner + allowlist validation
- [ ] README with usage examples
- [ ] GitHub release with binaries
