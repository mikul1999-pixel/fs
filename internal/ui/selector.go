package ui

import (
	"fmt"
	"os"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-isatty"
	"github.com/mikul1999-pixel/fs/internal/storage"
)

type SelectorOptions struct {
	Query      string
	FilterTags []string
	NoColor    bool
}

type model struct {
	shortcuts []storage.Shortcut
	cursor    int
	selected  *storage.Shortcut
	quitting  bool
	query     string
	tagFilter map[string]struct{}
	useColor  bool
	styles    selectorStyles
}

type selectorStyles struct {
	cursor     lipgloss.Style
	highlight  lipgloss.Style
	tag        lipgloss.Style
	matchedTag lipgloss.Style
}

func InitialModel(shortcuts []storage.Shortcut, opts SelectorOptions) model {
	tagFilter := make(map[string]struct{}, len(opts.FilterTags))
	for _, t := range opts.FilterTags {
		tagFilter[strings.ToLower(t)] = struct{}{}
	}

	useColor := shouldUseColor(opts.NoColor)
	renderer := lipgloss.NewRenderer(os.Stderr)

	return model{
		shortcuts: shortcuts,
		cursor:    0,
		query:     opts.Query,
		tagFilter: tagFilter,
		useColor:  useColor,
		styles:    newSelectorStyles(renderer),
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

		case "enter", " ":
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

	s.WriteString("Select a shortcut (up/down or j/k to move, Enter to select, q to quit):\n\n")

	for i, shortcut := range m.shortcuts {
		cursor := " "
		if m.cursor == i {
			cursor = m.applyStyle(">", m.styles.cursor)
		}

		highlightedName := highlightByTokens(shortcut.Name, m.query, m.useColor, m.styles.highlight)
		highlightedPath := highlightByTokens(shortcut.Path, m.query, m.useColor, m.styles.highlight)
		tagStr := m.renderTags(shortcut.Tags)

		s.WriteString(fmt.Sprintf("%s %d. %s -> %s%s\n", cursor, i+1, highlightedName, highlightedPath, tagStr))
	}

	return s.String()
}

func (m model) renderTags(tags []string) string {
	if len(tags) == 0 {
		return ""
	}

	rendered := make([]string, len(tags))
	for i, tag := range tags {
		if _, ok := m.tagFilter[strings.ToLower(tag)]; ok {
			rendered[i] = m.applyStyle(tag, m.styles.matchedTag)
			continue
		}
		rendered[i] = m.applyStyle(tag, m.styles.tag)
	}

	return fmt.Sprintf(" [%s]", strings.Join(rendered, ", "))
}

func (m model) applyStyle(text string, style lipgloss.Style) string {
	if !m.useColor {
		return text
	}
	return style.Render(text)
}

func highlightByTokens(text, query string, useColor bool, style lipgloss.Style) string {
	if !useColor || strings.TrimSpace(query) == "" {
		return text
	}

	tokens := queryTokens(query)
	if len(tokens) == 0 {
		return text
	}

	lowerText := strings.ToLower(text)
	highlightMask := make([]bool, len(text))

	for _, token := range tokens {
		start := 0
		for {
			idx := strings.Index(lowerText[start:], token)
			if idx == -1 {
				break
			}

			idx += start
			for i := idx; i < idx+len(token) && i < len(highlightMask); i++ {
				highlightMask[i] = true
			}
			start = idx + len(token)
		}
	}

	var b strings.Builder
	inHighlight := false
	segmentStart := 0

	for i := 0; i < len(text); i++ {
		if highlightMask[i] == inHighlight {
			continue
		}

		segment := text[segmentStart:i]
		if inHighlight {
			b.WriteString(style.Render(segment))
		} else {
			b.WriteString(segment)
		}

		inHighlight = highlightMask[i]
		segmentStart = i
	}

	tail := text[segmentStart:]
	if inHighlight {
		b.WriteString(style.Render(tail))
	} else {
		b.WriteString(tail)
	}

	return b.String()
}

func queryTokens(query string) []string {
	rawTokens := strings.Fields(strings.ToLower(strings.TrimSpace(query)))
	if len(rawTokens) == 0 {
		return nil
	}

	unique := make(map[string]struct{}, len(rawTokens))
	for _, token := range rawTokens {
		unique[token] = struct{}{}
	}

	tokens := make([]string, 0, len(unique))
	for token := range unique {
		tokens = append(tokens, token)
	}

	sort.Slice(tokens, func(i, j int) bool {
		return len(tokens[i]) > len(tokens[j])
	})

	return tokens
}

func shouldUseColor(noColor bool) bool {
	if noColor || os.Getenv("NO_COLOR") != "" {
		return false
	}

	fd := os.Stderr.Fd()
	return isatty.IsTerminal(fd) || isatty.IsCygwinTerminal(fd)
}

func newSelectorStyles(renderer *lipgloss.Renderer) selectorStyles {
	return selectorStyles{
		cursor:     renderer.NewStyle().Foreground(lipgloss.Color("212")).Bold(true),
		highlight:  renderer.NewStyle().Foreground(lipgloss.Color("11")).Bold(true),
		tag:        renderer.NewStyle(), // keep default terminal color for non matching tags
		matchedTag: renderer.NewStyle().Foreground(lipgloss.Color("39")),
	}
}

// Run the interactive selector
func RunSelector(shortcuts []storage.Shortcut, opts SelectorOptions) (string, error) {
	if len(shortcuts) == 0 {
		return "", fmt.Errorf("no shortcuts to select from")
	}

	// stderr/Error Output for bash function
	p := tea.NewProgram(
		InitialModel(shortcuts, opts),
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
