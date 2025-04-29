package lobby

import (
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fprasx/secrets-and-spies/service"
	"github.com/fprasx/secrets-and-spies/ui/menu"
)

var (
	Service *service.Spies
)

var (
	appStyle     = lipgloss.NewStyle().Align(lipgloss.Center, lipgloss.Center)
	spinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))
	textStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
)

type model struct {
	loading bool
	width   int
	height  int
	service *service.Spies
	spinner spinner.Model
}

func newModel() model {
	m := model{
		loading: true,
		service: service.New(menu.Address).
			WithName(menu.Name).
			WithHost(menu.Host),
		spinner: spinner.New(),
	}

	m.spinner.Spinner = spinner.Line

	return m
}

type joinMsg struct{}

func connectToHost(srv *service.Spies) func() tea.Msg {
	return func() tea.Msg {
		srv.Join(menu.HostAddress)
		return joinMsg{}
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(connectToHost(m.service), m.spinner.Tick)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Interrupt
		case "esc", "q":
			return m, tea.Quit
		}
	case joinMsg:
		m.loading = false
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m model) View() string {
	var content strings.Builder

	if m.loading {
		content.WriteString(
			lipgloss.JoinHorizontal(
                lipgloss.Center,
				spinnerStyle.Render(m.spinner.View()),
				textStyle.MarginLeft(4).Render("Loading..."),
			),
		)
	}

	return appStyle.Width(m.width).
		Height(m.height).
		Render(content.String())
}

func Show() {
	if _, err := tea.NewProgram(newModel()).Run(); err != nil {
		log.Fatal(err)
	}
}
