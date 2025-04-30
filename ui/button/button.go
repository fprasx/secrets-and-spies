package button

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("#cba6f7"))
	descriptionStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#cdd6f4"))
	containerStyle   = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder(), true, false, false, false).
				Padding(0, 1).
				BorderForeground(lipgloss.Color("#7f849c"))
	activeContainerStyle = containerStyle.
				UnsetBorderForeground().
				BorderForeground(lipgloss.Color("#a6e3a1"))
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

func New(title, description string) *Model {
	return &Model{title: title, description: description, hovered: false}
}

func (m *Model) SetHovered(hovered bool) *Model {
	m.hovered = hovered
	return m
}
