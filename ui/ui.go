package ui

import (
	"log"

	"github.com/fprasx/secrets-and-spies/game"
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
	wait = iota
	strike
	move
)

var (
	buttonContents = []buttonContent{
		{title: "Wait", description: "Stay at this city"},
		{title: "Strike", description: "Strike this city"},
		{title: "Move", description: "Move to city"},
	}
)

var (
	cities = []board.City{
		{Name: "Graz, Austria", Color: "#ff0000"},
		{Name: "Barcelona, Spain", Color: "#00ff00"},
		{Name: "Sicily, Italy", Color: "#0000ff"},
		{Name: "Lyon, France", Color: "#ff0000"},
		{Name: "Batman, Turkey", Color: "#ff0000"},
		{Name: "Vienna, Austria", Color: "#ff0000"},
		{Name: "Berlin, Germany", Color: "#00ff00"},
		{Name: "Paris, France", Color: "#0000ff"},
		{Name: "Madrid, Spain", Color: "#ffff00"},
		{Name: "Rome, Italy", Color: "#ff00ff"},
		{Name: "Warsaw, Poland", Color: "#00ffff"},
		{Name: "Prague, Czech Republic", Color: "#800000"},
		{Name: "Amsterdam, Netherlands", Color: "#008000"},
		{Name: "Copenhagen, Denmark", Color: "#000080"},
		{Name: "Lisbon, Portugal", Color: "#808000"},
		{Name: "Oslo, Norway", Color: "#800080"},
	}
	edges = map[int][]int{
		0:  {3, 7},
		1:  {10, 6},
		2:  {8, 15, 7},
		3:  {0, 12, 5, 7},
		4:  {12},
		5:  {10, 3, 12},
		6:  {1, 7},
		7:  {0, 2, 3, 6, 9, 14, 15},
		8:  {2},
		9:  {7},
		10: {1, 5},
		11: {13, 15},
		12: {3, 4, 5, 13, 15},
		13: {11, 12, 15},
		14: {7},
		15: {2, 7, 11, 12, 13},
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
	dead    bool
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
		active:  wait,
		board:   board.NewBoard(cities, edges),
		dead:    false,
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
	var cmd tea.Cmd

	cmds := []tea.Cmd{}

	player := m.service.MyPlayer()
	m.board.ActiveLocation = player.City

	if len(m.board.Edges[m.board.ActiveLocation]) <= m.board.CurrentSelection {
		m.board.CurrentSelection = 0
	}

	revealed := []int{}
	states := m.service.PlayerStates()

	for i := range states {
		if states[i].Revealed {
			revealed = append(revealed, states[i].City)
		}
	}

	if player.Dead {
		m.dead = true
	}

	m.board.Revealed = revealed

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
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
		case key.Matches(msg, keys.Confirm):
			if m.service.IsMyTurn() {
				currentLocation := m.board.ActiveLocation
				currentSelection := m.board.Edges[m.board.ActiveLocation][m.board.CurrentSelection]
				active := m.active

				m.service.DoTurn(func(s *service.Spies) game.Action {
					log.Printf("currentSelection: %v", currentSelection)

					switch active {
					case move:
						return game.Action{
							Type:       game.Move,
							TargetCity: currentSelection,
						}
					case strike:
						return game.Action{
							Type: game.Strike,
						}
					default:
						return game.Action{
							Type:       game.Move,
							TargetCity: currentLocation,
						}
					}

				})

				cmds = append(cmds, func() tea.Msg { return struct{}{} })
			}
		}
	}

	player = m.service.MyPlayer()
	m.board.ActiveLocation = player.City

	if len(m.board.Edges[m.board.ActiveLocation]) <= m.board.CurrentSelection {
		m.board.CurrentSelection = 0
	}

	if !m.service.IsMyTurn() {
		m.service.DoTurn(func(s *service.Spies) game.Action { return game.Action{} })
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

	_, cmd = m.board.Update(msg)
	cmds = append(cmds, cmd)

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

func (m *model) View() string {
	if len(m.board.Edges[m.board.ActiveLocation]) <= m.board.CurrentSelection {
		m.board.CurrentSelection = 0
	}

	if !m.dead {
		buttons := m.viewButtons()

		return containerStyle.
			Width(m.width).
			Height(m.height).
			Render(
				m.board.View() + "\n\n" + buttons + "\n\n\n",
			)
	} else {
		return containerStyle.
			Width(m.width).
			Height(m.height).
			Foreground(palette.Red).
			Bold(true).
			Render("You have died")
	}
}

func Show(service *service.Spies) {
	if _, err := tea.NewProgram(newModel(service), tea.WithAltScreen()).Run(); err != nil {
		log.Fatal(err)
	}
}
