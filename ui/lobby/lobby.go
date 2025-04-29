package lobby

import (
	"log"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fprasx/secrets-and-spies/service"
	"github.com/fprasx/secrets-and-spies/ui/menu"
	"github.com/fprasx/secrets-and-spies/utils"
)

var (
	Service *service.Spies
)

var (
	appStyle     = lipgloss.NewStyle().Align(lipgloss.Center, lipgloss.Center)
	spinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))
	textStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#11111b")).
			Background(lipgloss.Color("#f5c2e7")).
			Padding(0, 1)
)

type player struct {
	name    string
	address string
}

func (i player) Title() string       { return i.name }
func (i player) Description() string { return i.address }
func (i player) FilterValue() string { return i.name }

var (
	startKey = key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "start game"),
	)
)

type model struct {
	loading bool
	width   int
	height  int
	service *service.Spies
	spinner spinner.Model
	players list.Model
}

func newModel() model {
	m := model{
		loading: true,
		service: service.New(menu.Address).
			WithName(menu.Name).
			WithHost(menu.Host),
		spinner: spinner.New(),
		players: list.New([]list.Item{player{name: "Joe", address: "USA"}}, list.NewDefaultDelegate(), 0, 0),
	}

	m.spinner.Spinner = spinner.Line
	m.players.Title = "Lobby"
	m.players.Styles.Title = titleStyle
	m.players.SetShowStatusBar(false)
	m.players.SetShowFilter(false)

	if menu.Host {
		m.players.AdditionalShortHelpKeys = func() []key.Binding {
			return []key.Binding{startKey}
		}

		m.players.AdditionalFullHelpKeys = func() []key.Binding {
			return []key.Binding{startKey}
		}
	}

	Service = m.service

	return m
}

type joinMsg struct{}
type startMsg struct{}

func connectToHost(srv *service.Spies) func() tea.Msg {
	return func() tea.Msg {
		if !srv.IsHost() {
			srv.Join(menu.HostAddress)
		}
		return joinMsg{}
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(connectToHost(m.service), m.spinner.Tick)
}

func (m *model) updatePeers() tea.Cmd {
	peers := m.service.Peers()

	players := []list.Item{}

	for _, peer := range peers {
		players = append(players, player{
			name:    peer.Name,
			address: utils.AddrString(peer.Addr),
		})
	}

	cmd := m.players.SetItems(players)
	return cmd
}

func startGame(srv *service.Spies) tea.Cmd {
	return func() tea.Msg {
		srv.HostStart()
		return startMsg{}
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if cmd := m.updatePeers(); cmd != nil {
		return m, cmd
	}

	switch msg := msg.(type) {
	case startMsg:
		return m, tea.Quit

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.players.SetSize(40, m.height/2)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			return m, tea.Interrupt
		}

		switch {
		case key.Matches(msg, startKey):
			if menu.Host && m.service.Players() > 1 {
				return m, startGame(m.service)
			}
		}
	case joinMsg:
		m.loading = false
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	if m.service.Started() {
		return m, func() tea.Msg { return startMsg{} }
	}


	var cmd tea.Cmd
	m.players, cmd = m.players.Update(msg)

	return m, cmd
}

func (m model) viewLobby() string {
	return lipgloss.NewStyle().
		PaddingTop(4).
		PaddingBottom(4).
		Render(m.players.View())
}

func (m model) viewLoading() string {
	return lipgloss.JoinHorizontal(
		lipgloss.Center,
		spinnerStyle.Render(m.spinner.View()),
		textStyle.MarginLeft(4).Render("Loading..."),
	)
}

func (m model) View() string {
	if m.loading {
		return appStyle.
			Width(m.width).
			Height(m.height).
			Render(m.viewLoading())
	} else {
		return appStyle.
			Width(m.width).
			Height(m.height).
			Render(m.viewLobby())
	}

}

func Show() {
	if _, err := tea.NewProgram(newModel()).Run(); err != nil {
		log.Fatal(err)
	}
}
