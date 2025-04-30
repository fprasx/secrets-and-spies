package ui

import (
	"log"

	"github.com/fprasx/secrets-and-spies/service"
	"github.com/fprasx/secrets-and-spies/ui/board"
	"github.com/fprasx/secrets-and-spies/ui/button"
	"github.com/fprasx/secrets-and-spies/ui/palette"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	headerStyle    = lipgloss.NewStyle().Foreground(palette.Blue).Bold(true)
	containerStyle = lipgloss.NewStyle().Padding(1, 4, 0, 1).Align(lipgloss.Center, lipgloss.Center)
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

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit, k.Confirm}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right}, // first column
		{k.Help, k.Quit, k.Confirm},     // second column
	}
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "move left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "move right"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Confirm: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select action"),
	),
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
	help    help.Model
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

func (m *model) cursorLeft() {
	m.active = max(0, m.active-1)
}

func (m *model) cursorRight() {
	m.active = min(move, m.active+1)
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = min(msg.Width, 80)
		m.height = msg.Height
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Up):
		case key.Matches(msg, keys.Down):
		case key.Matches(msg, keys.Left):
			m.cursorLeft()
		case key.Matches(msg, keys.Right):
			m.cursorRight()
		case key.Matches(msg, keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, keys.Quit):
			cmds = append(cmds, tea.Quit)
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

	return lipgloss.JoinHorizontal(lipgloss.Center, buttons...)
}

func (m *model) viewAppBorder(text string) string {
	return lipgloss.PlaceHorizontal(
		m.width,
		lipgloss.Left,
		headerStyle.Render(text),
		lipgloss.WithWhitespaceChars("/"),
		lipgloss.WithWhitespaceForeground(palette.Blue),
	)
}

func (m *model) View() string {
	buttons := m.viewButtons()
	header := m.viewAppBorder("Secrets and Spies ")
	footer := m.viewAppBorder("")

	return containerStyle.Render(
		header + "\n" + buttons + "\n\n" + footer,
	)
}

func Show(service *service.Spies) {
	if _, err := tea.NewProgram(newModel(service), tea.WithAltScreen()).Run(); err != nil {
		log.Fatal(err)
	}
}
