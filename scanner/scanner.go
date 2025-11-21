package scanner

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// CacheEntry represents a scanned directory or file
type CacheEntry struct {
	Name        string
	Path        string
	Size        int64
	FileCount   int
	LastMod     time.Time
	OldestMod   time.Time
	Selected    bool
	Description string
	IsParent    bool          // True if this is a parent category
	Children    []*CacheEntry // Sub-items within this category
	Expanded    bool          // Whether children are visible
	Depth       int           // Nesting level for display
}

// ScanResult holds all scan results
type ScanResult struct {
	Entries   []*CacheEntry
	TotalSize int64
	ScanTime  time.Duration
}

// Scanner handles directory scanning
type Scanner struct {
	HomeDir string
}

// NewScanner creates a new scanner
func NewScanner() (*Scanner, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	return &Scanner{HomeDir: home}, nil
}

// GetAllowedPaths returns the list of allowed paths to scan
func (s *Scanner) GetAllowedPaths() []struct {
	Path        string
	Description string
} {
	return []struct {
		Path        string
		Description string
	}{
		{filepath.Join(s.HomeDir, "Library", "Caches"), "System & App Caches"},
		{filepath.Join(s.HomeDir, "Library", "Logs"), "Log Files"},
		{filepath.Join(s.HomeDir, "Library", "Developer", "Xcode", "DerivedData"), "Xcode Build Data"},
		{filepath.Join(s.HomeDir, "Library", "Developer", "Xcode", "Archives"), "Xcode Archives"},
		{filepath.Join(s.HomeDir, ".npm", "_cacache"), "npm Cache"},
		{filepath.Join(s.HomeDir, ".cache", "yarn"), "Yarn Cache"},
		{filepath.Join(s.HomeDir, "Library", "Caches", "pip"), "Python pip Cache"},
		{filepath.Join(s.HomeDir, "Library", "Caches", "Homebrew"), "Homebrew Cache"},
		{filepath.Join(s.HomeDir, ".gradle", "caches"), "Gradle Cache"},
		{filepath.Join(s.HomeDir, ".cargo", "registry"), "Cargo Registry"},
		{filepath.Join(s.HomeDir, "Library", "Caches", "Google", "Chrome"), "Chrome Cache"},
		{filepath.Join(s.HomeDir, "Library", "Caches", "com.apple.Safari"), "Safari Cache"},
	}
}

// Scan performs the scan of all allowed paths
func (s *Scanner) Scan() (*ScanResult, error) {
	start := time.Now()
	var entries []*CacheEntry
	var totalSize int64

	for _, target := range s.GetAllowedPaths() {
		entry, err := s.scanPathWithChildren(target.Path, target.Description)
		if err != nil {
			continue // Skip paths that don't exist or can't be read
		}
		if entry.Size > 0 {
			entries = append(entries, entry)
			totalSize += entry.Size
		}
	}

	// Sort by size descending by default
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Size > entries[j].Size
	})

	return &ScanResult{
		Entries:   entries,
		TotalSize: totalSize,
		ScanTime:  time.Since(start),
	}, nil
}

// scanPathWithChildren scans a path and its immediate children
func (s *Scanner) scanPathWithChildren(path, description string) (*CacheEntry, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	entry := &CacheEntry{
		Name:        filepath.Base(path),
		Path:        path,
		Description: description,
		LastMod:     info.ModTime(),
		OldestMod:   info.ModTime(),
		IsParent:    true,
		Expanded:    false,
		Depth:       0,
	}

	// Read immediate children
	dirEntries, err := os.ReadDir(path)
	if err != nil {
		// If we can't read children, just scan the whole thing
		entry.IsParent = false
		s.scanSize(entry)
		return entry, nil
	}

	var children []*CacheEntry
	for _, de := range dirEntries {
		childPath := filepath.Join(path, de.Name())
		childInfo, err := os.Stat(childPath)
		if err != nil {
			continue
		}

		child := &CacheEntry{
			Name:     de.Name(),
			Path:     childPath,
			LastMod:  childInfo.ModTime(),
			OldestMod: childInfo.ModTime(),
			Depth:    1,
		}

		// Calculate size for each child
		if childInfo.IsDir() {
			s.scanSize(child)
		} else {
			child.Size = childInfo.Size()
			child.FileCount = 1
		}

		if child.Size > 0 {
			children = append(children, child)
			entry.Size += child.Size
			entry.FileCount += child.FileCount
			if child.LastMod.After(entry.LastMod) {
				entry.LastMod = child.LastMod
			}
		}
	}

	// Sort children by size
	sort.Slice(children, func(i, j int) bool {
		return children[i].Size > children[j].Size
	})

	entry.Children = children

	// If no children or only one, don't show as expandable
	if len(children) <= 1 {
		entry.IsParent = false
	}

	return entry, nil
}

// scanSize recursively calculates size of a directory
func (s *Scanner) scanSize(entry *CacheEntry) {
	filepath.Walk(entry.Path, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() {
			entry.Size += info.Size()
			entry.FileCount++
			if info.ModTime().After(entry.LastMod) {
				entry.LastMod = info.ModTime()
			}
			if info.ModTime().Before(entry.OldestMod) {
				entry.OldestMod = info.ModTime()
			}
		}
		return nil
	})
}

// FormatSize formats bytes into human-readable string
func FormatSize(bytes int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)

	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.1f GB", float64(bytes)/float64(GB))
	case bytes >= MB:
		return fmt.Sprintf("%.1f MB", float64(bytes)/float64(MB))
	case bytes >= KB:
		return fmt.Sprintf("%.1f KB", float64(bytes)/float64(KB))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

// ShortenPath shortens a path for display
func ShortenPath(path string) string {
	home, _ := os.UserHomeDir()
	if home != "" {
		return "~" + path[len(home):]
	}
	return path
}
