package ui

import (
	"log"

	"github.com/fprasx/secrets-and-spies/service"
	"github.com/fprasx/secrets-and-spies/ui/board"
	"github.com/fprasx/secrets-and-spies/ui/button"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type keyMap struct {
	Up      key.Binding
	Down    key.Binding
	Left    key.Binding
	Right   key.Binding
	Confirm key.Binding
	Help    key.Binding
	Quit    key.Binding
}

type buttonContent struct {
	title       string
	description string
}

const (
	capture = iota
	wait
	strike
	move
)

var (
	buttonContents = []buttonContent{
		{title: "Capture", description: "Capture this city"},
		{title: "Wait", description: "Stay at this city"},
		{title: "Strike", description: "Strike this city"},
		{title: "Move", description: "Move to city"},
	}
)

type model struct {
	service *service.Spies
	board   *board.Board
	buttons []*button.Model
	width   int
	height  int
	active  int
}

func newModel(service *service.Spies) *model {
	buttons := []*button.Model{}

	for i := range buttonContents {
		buttons = append(buttons, button.New(
			buttonContents[i].title,
			buttonContents[i].description,
		))
	}

	return &model{
		buttons: buttons,
		service: service,
		active:  capture,
	}
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

func (m *model) viewButtons() string {
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

func (m *model) View() string {
	return m.viewButtons()
}

func Show(service *service.Spies) {
	if _, err := tea.NewProgram(newModel(service), tea.WithAltScreen()).Run(); err != nil {
		log.Fatal(err)
	}
}
