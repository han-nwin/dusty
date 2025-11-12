package ui

import tea "github.com/charmbracelet/bubbletea"

// Model represents the TUI state
type Model struct {
	quitting bool
	message  string
}

func InitialModel() Model {
	return Model{
		quitting: false,
		message:  "Dusty - macOS cleanup TUI",
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			m.message = "Goodbye!"
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Model) View() string {
	if m.quitting {
		return m.message + "\n"
	}
	return m.message + "\nPress q to quit.\n"
}
