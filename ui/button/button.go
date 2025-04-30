package button

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("#cba6f7"))
	descriptionStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#7f849c"))
	containerStyle   = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder(), true, false, false, false).
				Padding(0, 1)
	activeContainerStyle = containerStyle.BorderForeground(lipgloss.Color("#a6e3a1"))
)

type HoverMsg struct {
	Hover bool
}

type Model struct {
	title       string
	description string
	hovered     bool
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case HoverMsg:
		m.hovered = msg.Hover
	}

	return m, nil
}

func (m *Model) View() string {
	containerStyle := containerStyle

	if m.hovered {
		containerStyle = activeContainerStyle
	}

	return containerStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			titleStyle.Render(m.title),
			descriptionStyle.Render(m.description),
		),
	)
}

func New() *Model {
	return new(Model)
}

func (m *Model) WithTitle(title string) *Model {
	m.title = title
	return m
}

func (m *Model) WithDesciption(description string) *Model {
	m.description = description
	return m
}

func (m *Model) SetHovered(hovered bool) *Model {
	m.hovered = hovered
	return m
}
