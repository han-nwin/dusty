# Dusty 🧹

A CleanMyMac-style TUI for macOS. Clean up caches, logs, and build artifacts safely from your terminal.

## Features

- Scan common cache directories (Xcode, npm, yarn, system caches)
- Interactive TUI with keyboard navigation
- Safe deletion with confirmation
- Move files to Trash or permanent delete
- Track cleanup history

## Installation

```bash
go install github.com/han/dusty@latest
```

Or build from source:

```bash
git clone https://github.com/han/dusty
cd dusty
go build -o dusty
```

## Usage

```bash
dusty
```

### Keyboard Shortcuts

- `↑↓` - Navigate
- `Space` - Toggle selection
- `Enter` - View details
- `C` - Clean selected
- `R` - Rescan
- `/` - Filter
- `?` - Help
- `Q` - Quit

## What Gets Cleaned

- `~/Library/Caches/**` - System and app caches
- `~/Library/Logs/**` - Log files
- `~/Library/Developer/Xcode/DerivedData` - Xcode build artifacts
- `~/.npm/_cacache` - npm cache
- `~/.cache/yarn` - Yarn cache
- `~/Library/Caches/pip` - Python pip cache

## Directory Structure

```
dusty/
├── main.go             # Application entry point
├── scanner/            # Directory scanning and size calculation
├── ui/                 # Bubble Tea TUI components
├── FEATURES.md         # Feature specifications
└── TODOs.md           # Development roadmap
```

## Safety

- Only scans allowlisted paths
- Never requires sudo
- Confirmation before deletion
- Saves cleanup manifest to `~/.dusty/undo/`

## License

MIT
