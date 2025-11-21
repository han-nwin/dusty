package ui

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/han/dusty/scanner"
)

// Catppuccin Mocha colors
var (
	colorRosewater = lipgloss.Color("#f5e0dc")
	colorFlamingo  = lipgloss.Color("#f2cdcd")
	colorPink      = lipgloss.Color("#f5c2e7")
	colorMauve     = lipgloss.Color("#cba6f7")
	colorRed       = lipgloss.Color("#f38ba8")
	colorMaroon    = lipgloss.Color("#eba0ac")
	colorPeach     = lipgloss.Color("#fab387")
	colorYellow    = lipgloss.Color("#f9e2af")
	colorGreen     = lipgloss.Color("#a6e3a1")
	colorTeal      = lipgloss.Color("#94e2d5")
	colorSky       = lipgloss.Color("#89dceb")
	colorSapphire  = lipgloss.Color("#74c7ec")
	colorBlue      = lipgloss.Color("#89b4fa")
	colorLavender  = lipgloss.Color("#b4befe")
	colorText      = lipgloss.Color("#cdd6f4")
	colorSubtext1  = lipgloss.Color("#bac2de")
	colorSubtext0  = lipgloss.Color("#a6adc8")
	colorOverlay2  = lipgloss.Color("#9399b2")
	colorOverlay1  = lipgloss.Color("#7f849c")
	colorOverlay0  = lipgloss.Color("#6c7086")
	colorSurface2  = lipgloss.Color("#585b70")
	colorSurface1  = lipgloss.Color("#45475a")
	colorSurface0  = lipgloss.Color("#313244")
	colorBase      = lipgloss.Color("#1e1e2e")
	colorMantle    = lipgloss.Color("#181825")
	colorCrust     = lipgloss.Color("#11111b")
)

// Styles
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorMauve).
			MarginBottom(1)

	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorBlue).
			Background(colorSurface0).
			Padding(0, 1)

	selectedStyle = lipgloss.NewStyle().
			Foreground(colorBase).
			Background(colorBlue).
			Bold(true)

	normalStyle = lipgloss.NewStyle().
			Foreground(colorText)

	dimStyle = lipgloss.NewStyle().
			Foreground(colorOverlay0)

	pathStyle = lipgloss.NewStyle().
			Foreground(colorOverlay1).
			Italic(true)

	checkboxStyle = lipgloss.NewStyle().
			Foreground(colorPink)

	// Size colors - pastel Catppuccin
	sizeStyleSmall = lipgloss.NewStyle().
			Foreground(colorGreen)

	sizeStyleMedium = lipgloss.NewStyle().
			Foreground(colorYellow)

	sizeStyleLarge = lipgloss.NewStyle().
			Foreground(colorMaroon)

	statusStyle = lipgloss.NewStyle().
			Foreground(colorSubtext0).
			MarginTop(1)

	helpStyle = lipgloss.NewStyle().
			Foreground(colorOverlay1)

	confirmStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorRed)

	successStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorGreen)

	borderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorSurface2).
			Padding(1, 2)

	expandIcon   = lipgloss.NewStyle().Foreground(colorPeach).Render("▶")
	collapseIcon = lipgloss.NewStyle().Foreground(colorPeach).Render("▼")
)

type viewState int

const (
	viewList viewState = iota
	viewScanning
	viewConfirm
	viewCleaning
	viewHelp
	viewFilter
)

// Messages
type scanCompleteMsg struct {
	result *scanner.ScanResult
	err    error
}

type cleanCompleteMsg struct {
	cleaned int64
	err     error
}

// displayEntry is a flattened entry for display
type displayEntry struct {
	entry    *scanner.CacheEntry
	isChild  bool
	parentIdx int
}

// Model represents the TUI state
type Model struct {
	entries       []*scanner.CacheEntry
	displayList   []displayEntry
	cursor        int
	totalSize     int64
	selectedSize  int64
	scanTime      time.Duration
	state         viewState
	spinner       spinner.Model
	filterInput   textinput.Model
	filter        string
	width         int
	height        int
	message       string
	err           error
	confirmAction string // "delete" or "trash"
}

func InitialModel() Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(colorMauve)

	ti := textinput.New()
	ti.Placeholder = "Filter..."
	ti.CharLimit = 50
	ti.Width = 30

	return Model{
		state:       viewScanning,
		spinner:     s,
		filterInput: ti,
		width:       80,
		height:      24,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, scanCmd())
}

func scanCmd() tea.Cmd {
	return func() tea.Msg {
		s, err := scanner.NewScanner()
		if err != nil {
			return scanCompleteMsg{err: err}
		}
		result, err := s.Scan()
		return scanCompleteMsg{result: result, err: err}
	}
}

func (m *Model) rebuildDisplayList() {
	m.displayList = nil
	filterLower := strings.ToLower(m.filter)

	for i, entry := range m.entries {
		// Check if entry matches filter
		if m.filter != "" {
			if !strings.Contains(strings.ToLower(entry.Name), filterLower) &&
				!strings.Contains(strings.ToLower(entry.Description), filterLower) {
				continue
			}
		}

		m.displayList = append(m.displayList, displayEntry{
			entry:    entry,
			isChild:  false,
			parentIdx: i,
		})

		// Add children if expanded
		if entry.Expanded && entry.IsParent {
			for _, child := range entry.Children {
				if m.filter != "" {
					if !strings.Contains(strings.ToLower(child.Name), filterLower) {
						continue
					}
				}
				m.displayList = append(m.displayList, displayEntry{
					entry:    child,
					isChild:  true,
					parentIdx: i,
				})
			}
		}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case scanCompleteMsg:
		if msg.err != nil {
			m.err = msg.err
			m.state = viewList
			return m, nil
		}
		m.entries = msg.result.Entries
		m.totalSize = msg.result.TotalSize
		m.scanTime = msg.result.ScanTime
		m.state = viewList
		m.rebuildDisplayList()
		m.updateSelectedSize()
		return m, nil

	case cleanCompleteMsg:
		if msg.err != nil {
			m.message = fmt.Sprintf("Error: %v", msg.err)
		} else {
			m.message = fmt.Sprintf("Cleaned %s!", scanner.FormatSize(msg.cleaned))
		}
		m.state = viewList
		return m, scanCmd()

	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	}

	return m, nil
}

func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Handle filter input mode
	if m.state == viewFilter {
		switch msg.String() {
		case "enter":
			m.filter = m.filterInput.Value()
			m.state = viewList
			m.rebuildDisplayList()
			m.cursor = 0
			return m, nil
		case "esc":
			m.filterInput.SetValue(m.filter)
			m.state = viewList
			return m, nil
		default:
			var cmd tea.Cmd
			m.filterInput, cmd = m.filterInput.Update(msg)
			return m, cmd
		}
	}

	// Handle confirmation mode
	if m.state == viewConfirm {
		switch msg.String() {
		case "y", "Y":
			m.state = viewCleaning
			return m, m.cleanCmd()
		case "n", "N", "esc":
			m.state = viewList
			return m, nil
		}
		return m, nil
	}

	// Handle help mode
	if m.state == viewHelp {
		m.state = viewList
		return m, nil
	}

	// Normal list navigation
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit

	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}

	case "down", "j":
		if m.cursor < len(m.displayList)-1 {
			m.cursor++
		}

	case "enter", "l":
		// Toggle expand/collapse for parent items
		if len(m.displayList) > 0 && m.cursor < len(m.displayList) {
			de := m.displayList[m.cursor]
			if !de.isChild && de.entry.IsParent {
				de.entry.Expanded = !de.entry.Expanded
				m.rebuildDisplayList()
			}
		}

	case " ":
		if len(m.displayList) > 0 && m.cursor < len(m.displayList) {
			de := m.displayList[m.cursor]
			de.entry.Selected = !de.entry.Selected

			// If selecting a parent, select/deselect all children too
			if !de.isChild && de.entry.IsParent {
				for _, child := range de.entry.Children {
					child.Selected = de.entry.Selected
				}
			}
			m.updateSelectedSize()
		}

	case "a":
		// Select all
		for _, entry := range m.entries {
			entry.Selected = true
			for _, child := range entry.Children {
				child.Selected = true
			}
		}
		m.updateSelectedSize()

	case "A":
		// Deselect all
		for _, entry := range m.entries {
			entry.Selected = false
			for _, child := range entry.Children {
				child.Selected = false
			}
		}
		m.updateSelectedSize()

	case "c":
		// Clean (permanent delete)
		if m.selectedSize > 0 {
			m.confirmAction = "delete"
			m.state = viewConfirm
		}

	case "t":
		// Move to trash
		if m.selectedSize > 0 {
			m.confirmAction = "trash"
			m.state = viewConfirm
		}

	case "r", "R":
		m.state = viewScanning
		m.message = ""
		m.cursor = 0
		return m, tea.Batch(m.spinner.Tick, scanCmd())

	case "/":
		m.state = viewFilter
		m.filterInput.Focus()
		return m, textinput.Blink

	case "?":
		m.state = viewHelp

	case "esc":
		if m.filter != "" {
			m.filter = ""
			m.filterInput.SetValue("")
			m.rebuildDisplayList()
			m.cursor = 0
		}
	}

	return m, nil
}

func (m *Model) updateSelectedSize() {
	m.selectedSize = 0
	for _, entry := range m.entries {
		if entry.Selected {
			m.selectedSize += entry.Size
		} else {
			// Check children
			for _, child := range entry.Children {
				if child.Selected {
					m.selectedSize += child.Size
				}
			}
		}
	}
}

func (m Model) cleanCmd() tea.Cmd {
	action := m.confirmAction
	return func() tea.Msg {
		var cleaned int64
		var toClean []string

		// Collect paths to clean
		for _, entry := range m.entries {
			if entry.Selected {
				toClean = append(toClean, entry.Path)
				cleaned += entry.Size
			} else {
				for _, child := range entry.Children {
					if child.Selected {
						toClean = append(toClean, child.Path)
						cleaned += child.Size
					}
				}
			}
		}

		// Clean each path
		for _, path := range toClean {
			var err error
			if action == "trash" {
				// Move to trash using AppleScript
				script := fmt.Sprintf(`tell app "Finder" to delete POSIX file "%s"`, path)
				cmd := exec.Command("osascript", "-e", script)
				err = cmd.Run()
			} else {
				// Permanent delete
				err = os.RemoveAll(path)
			}
			if err != nil {
				return cleanCompleteMsg{err: fmt.Errorf("failed to clean %s: %v", path, err)}
			}
		}

		return cleanCompleteMsg{cleaned: cleaned}
	}
}

func (m Model) View() string {
	switch m.state {
	case viewScanning:
		return m.viewScanning()
	case viewConfirm:
		return m.viewConfirm()
	case viewHelp:
		return m.viewHelp()
	case viewFilter:
		return m.viewFilter()
	default:
		return m.viewList()
	}
}

func (m Model) viewScanning() string {
	content := fmt.Sprintf("\n\n   %s Scanning directories...\n\n", m.spinner.View())
	return lipgloss.NewStyle().Foreground(colorText).Render(content)
}

func (m Model) viewList() string {
	var b strings.Builder

	// Title
	title := titleStyle.Render("  Dusty")
	subtitle := dimStyle.Render("  Clean up your Mac")
	b.WriteString(title + "\n" + subtitle + "\n\n")

	// Error display
	if m.err != nil {
		b.WriteString(confirmStyle.Render(fmt.Sprintf("  Error: %v", m.err)) + "\n\n")
	}

	// Success message
	if m.message != "" {
		b.WriteString(successStyle.Render("  " + m.message) + "\n\n")
	}

	// Entries
	if len(m.displayList) == 0 {
		b.WriteString(dimStyle.Render("\n  No items found.\n"))
	}

	// Calculate visible range
	visibleHeight := m.height - 14
	if visibleHeight < 5 {
		visibleHeight = 5
	}
	start := 0
	if m.cursor >= visibleHeight {
		start = m.cursor - visibleHeight + 1
	}
	end := start + visibleHeight
	if end > len(m.displayList) {
		end = len(m.displayList)
	}

	for i := start; i < end; i++ {
		de := m.displayList[i]
		e := de.entry

		var line string
		if de.isChild {
			// Child item - indented
			line = m.renderChildItem(e, i == m.cursor)
		} else {
			// Parent item
			line = m.renderParentItem(e, i == m.cursor)
		}

		b.WriteString(line + "\n")
	}

	// Status bar
	b.WriteString("\n")

	// Stats line
	statsLine := fmt.Sprintf("  Total: %s  │  Selected: %s  │  Items: %d",
		lipgloss.NewStyle().Foreground(colorBlue).Render(scanner.FormatSize(m.totalSize)),
		lipgloss.NewStyle().Foreground(colorPink).Render(scanner.FormatSize(m.selectedSize)),
		len(m.displayList))

	if m.filter != "" {
		statsLine += fmt.Sprintf("  │  Filter: %s", lipgloss.NewStyle().Foreground(colorYellow).Render(m.filter))
	}

	b.WriteString(statusStyle.Render(statsLine) + "\n\n")

	// Help
	help := "  ↑↓ navigate • space select • enter expand • a/A all/none • 🗑️ t trash • 💀 c clean • 🔄 r rescan • 🔍 / filter • ❓ ? help • 👋 q quit"
	b.WriteString(helpStyle.Render(help) + "\n")

	return b.String()
}

func (m Model) renderParentItem(e *scanner.CacheEntry, isCursor bool) string {
	cursor := "  "
	if isCursor {
		cursor = "👉"
	}

	checkbox := dimStyle.Render("[ ]")
	if e.Selected {
		checkbox = lipgloss.NewStyle().Foreground(colorRed).Render("[✓]")
	}

	// Expand/collapse icon
	icon := "  "
	if e.IsParent {
		if e.Expanded {
			icon = collapseIcon + " "
		} else {
			icon = expandIcon + " "
		}
	}

	// Size with color
	sizeStr := m.colorSize(e.Size)

	// Name and description
	name := e.Name
	if len(name) > 20 {
		name = name[:17] + "..."
	}

	// Path
	path := scanner.ShortenPath(e.Path)

	// File count and date with colors
	files := lipgloss.NewStyle().Foreground(colorSapphire).Render(fmt.Sprintf("%d files", e.FileCount))
	date := lipgloss.NewStyle().Foreground(colorLavender).Render(e.LastMod.Format("Jan 02"))

	// Build the line
	line := fmt.Sprintf("%s%s %s%-20s  %10s  %12s  %s",
		cursor, checkbox, icon, name, sizeStr, files, date)

	// Second line with path
	pathLine := fmt.Sprintf("       %s", pathStyle.Render(path))

	if isCursor {
		return selectedStyle.Render(line) + "\n" + pathLine
	}
	return normalStyle.Render(line) + "\n" + pathLine
}

func (m Model) renderChildItem(e *scanner.CacheEntry, isCursor bool) string {
	cursor := "    "
	if isCursor {
		cursor = "  👉"
	}

	checkbox := dimStyle.Render("[ ]")
	if e.Selected {
		checkbox = lipgloss.NewStyle().Foreground(colorRed).Render("[✓]")
	}

	// Size with color
	sizeStr := m.colorSize(e.Size)

	// Name
	name := e.Name
	if len(name) > 25 {
		name = name[:22] + "..."
	}

	// File count and date with colors
	files := lipgloss.NewStyle().Foreground(colorSapphire).Render(fmt.Sprintf("%d files", e.FileCount))
	date := lipgloss.NewStyle().Foreground(colorLavender).Render(e.LastMod.Format("Jan 02"))

	line := fmt.Sprintf("%s%s   %-25s  %10s  %12s  %s",
		cursor, checkbox, name, sizeStr, files, date)

	if isCursor {
		return selectedStyle.Render(line)
	}
	return dimStyle.Render(line)
}

func (m Model) colorSize(size int64) string {
	sizeStr := scanner.FormatSize(size)
	switch {
	case size >= 2*1024*1024*1024: // > 2GB
		return sizeStyleLarge.Render(sizeStr)
	case size >= 500*1024*1024: // > 500MB
		return sizeStyleMedium.Render(sizeStr)
	default:
		return sizeStyleSmall.Render(sizeStr)
	}
}

func (m Model) viewConfirm() string {
	var b strings.Builder

	actionEmoji := "🗑️"
	actionText := "Move to Trash"
	warning := ""
	if m.confirmAction == "delete" {
		actionEmoji = "💀"
		actionText = "Clean"
		warning = confirmStyle.Render("  ⚠️  WARNING: This will PERMANENTLY remove these files! They cannot be recovered!\n\n")
	}

	b.WriteString(titleStyle.Render(fmt.Sprintf("  %s Confirm %s", actionEmoji, actionText)) + "\n\n")
	if warning != "" {
		b.WriteString(warning)
	}

	var count int
	var items []string
	for _, entry := range m.entries {
		if entry.Selected {
			count++
			items = append(items, fmt.Sprintf("  • %s (%s)\n    %s",
				entry.Name,
				scanner.FormatSize(entry.Size),
				pathStyle.Render(scanner.ShortenPath(entry.Path))))
		} else {
			for _, child := range entry.Children {
				if child.Selected {
					count++
					items = append(items, fmt.Sprintf("  • %s (%s)\n    %s",
						child.Name,
						scanner.FormatSize(child.Size),
						pathStyle.Render(scanner.ShortenPath(child.Path))))
				}
			}
		}
	}

	for _, item := range items {
		b.WriteString(normalStyle.Render(item) + "\n\n")
	}

	b.WriteString(confirmStyle.Render(fmt.Sprintf("  %s %d items (%s)?", actionText, count, scanner.FormatSize(m.selectedSize))))
	b.WriteString("\n\n")
	b.WriteString(helpStyle.Render("  Press y to confirm, n to cancel"))
	b.WriteString("\n")

	return b.String()
}

func (m Model) viewHelp() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("  Help") + "\n\n")

	help := []struct {
		key  string
		desc string
	}{
		{"↑/k, ↓/j", "Navigate up/down"},
		{"Enter/l", "Expand/collapse"},
		{"Space", "Toggle selection"},
		{"a", "Select all"},
		{"A", "Deselect all"},
		{"t", "🗑️  Move to Trash"},
		{"c", "💀 Clean (permanent)"},
		{"r", "🔄 Rescan directories"},
		{"/", "🔍 Filter items"},
		{"Esc", "Clear filter"},
		{"?", "❓ Show this help"},
		{"q", "👋 Quit"},
	}

	for _, h := range help {
		key := lipgloss.NewStyle().Foreground(colorPeach).Width(12).Render(h.key)
		b.WriteString(fmt.Sprintf("  %s %s\n", key, h.desc))
	}

	b.WriteString("\n")
	b.WriteString(helpStyle.Render("  Press any key to return"))
	b.WriteString("\n")

	return b.String()
}

func (m Model) viewFilter() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("  Filter") + "\n\n")
	b.WriteString("  " + m.filterInput.View() + "\n\n")
	b.WriteString(helpStyle.Render("  Press Enter to apply, Esc to cancel"))
	b.WriteString("\n")

	return b.String()
}
