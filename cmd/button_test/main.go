package main

import (
	"log"

	"github.com/fprasx/secrets-and-spies/ui/button"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	buttons       []*button.Model
	width, height int
	active        int
}

func newModel() *model {
	m := &model{
		buttons: []*button.Model{
			button.New().WithTitle("Capture").WithDesciption("Capture this city"),
			button.New().WithTitle("Wait").WithDesciption("Stay at this city"),
			button.New().WithTitle("Strike").WithDesciption("Strike this city"),
			button.New().WithTitle("Move").WithDesciption("Move to city"),
		},
		active: 0,
	}
	return m
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "right", "l":
			m.active = min(len(m.buttons)-1, m.active+1)
		case "left", "h":
			m.active = max(0, m.active-1)
		}
	}

	for i := range m.buttons {
		if i == m.active {
			m.buttons[i].SetHovered(true)
		} else {
            m.buttons[i].SetHovered(false)
        }
		_, cmd := m.buttons[i].Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *model) View() string {
	buttons := []string{}

	for i := range m.buttons {
        view := m.buttons[i].View()

        if i != 0 {
            view = lipgloss.NewStyle().MarginLeft(4).Render(view)
        }

		buttons = append(buttons, view)
	}

	return lipgloss.
		NewStyle().
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center, lipgloss.Center).
		Render(
			lipgloss.JoinHorizontal(lipgloss.Center, buttons...),
		)
}

func main() {
	if _, err := tea.NewProgram(newModel(), tea.WithAltScreen()).Run(); err != nil {
		log.Fatal(err)
	}
}
