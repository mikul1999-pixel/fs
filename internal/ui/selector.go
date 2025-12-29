package ui

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mikul1999-pixel/fs/internal/storage"
)

type model struct {
	shortcuts []storage.Shortcut
	cursor    int
	selected  *storage.Shortcut
	quitting  bool
}

func InitialModel(shortcuts []storage.Shortcut) model {
	return model{
		shortcuts: shortcuts,
		cursor:    0,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			m.quitting = true
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.shortcuts)-1 {
				m.cursor++
			}

		case "enter":
			m.selected = &m.shortcuts[m.cursor]
			m.quitting = true
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	if m.quitting {
		return ""
	}

	var s strings.Builder

	s.WriteString("Select a shortcut (↑/↓ or j/k to move, Enter to select, q to quit):\n\n")

	for i, shortcut := range m.shortcuts {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		tagStr := ""
		if len(shortcut.Tags) > 0 {
			tagStr = fmt.Sprintf(" [%s]", strings.Join(shortcut.Tags, ", "))
		}

		s.WriteString(fmt.Sprintf("%s %d. %s -> %s%s\n", cursor, i+1, shortcut.Name, shortcut.Path, tagStr))
	}

	return s.String()
}

// Run the interactive selector
func RunSelector(shortcuts []storage.Shortcut) (string, error) {
	if len(shortcuts) == 0 {
		return "", fmt.Errorf("no shortcuts to select from")
	}

	// stderr/Error Output for bash function
	p := tea.NewProgram(
		InitialModel(shortcuts),
		tea.WithInput(os.Stdin),
		tea.WithOutput(os.Stderr),
	)
	
	finalModel, err := p.Run()
	if err != nil {
		return "", err
	}

	m := finalModel.(model)
	if m.selected != nil {
		return m.selected.Path, nil
	}

	return "", fmt.Errorf("no selection made")
}