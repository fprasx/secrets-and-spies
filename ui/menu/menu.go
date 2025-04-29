package menu

import (
	_ "embed"
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/fprasx/secrets-and-spies/utils"
)

const minWidth = 50

//go:embed banner.txt
var banner string

var (
	Address string
	Name    string
	Host    bool
)

var (
	appStyle    = lipgloss.NewStyle().Align(lipgloss.Center, lipgloss.Center)
	headerStyle = lipgloss.NewStyle().Bold(true)
	formStyle   = lipgloss.NewStyle().MarginTop(1)
)

type model struct {
	form *huh.Form
}

func newModel() model {
	m := model{
		form: huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Name").
					Value(&Name).
					Description("Your display name").
					Prompt("").
					Validate(func(s string) error {
						if len(s) <= 0 {
							return fmt.Errorf("Name cannot be empty")
						}
						return nil
					}).
					Placeholder("Alice"),

				huh.NewInput().
					Value(&Address).
					Title("Address").
					Prompt("").
					Description("Your network address").
					Validate(utils.ValidateAddr).
					Placeholder("unix:///tmp/spies/socket"),

				huh.NewConfirm().
					Value(&Host).
					Affirmative("Host").
					Negative("Join"),
			),
		).
			WithWidth(minWidth).
			WithTheme(huh.ThemeCatppuccin()).
			WithShowHelp(false).
			WithShowErrors(true),
	}

	return m
}

func (m model) Init() tea.Cmd {
	return m.form.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		appStyle = appStyle.UnsetHeight().Height(msg.Height).Width(msg.Width)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Interrupt
		case "esc", "q":
			return m, tea.Quit
		}
	}

	var cmds []tea.Cmd

	// Process the form
	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
		cmds = append(cmds, cmd)
	}

	if m.form.State == huh.StateCompleted {
		// Quit when the form is done.
		cmds = append(cmds, tea.Quit)
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.form.State == huh.StateCompleted {
		return ""
	}

	return appStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			headerStyle.Render(banner),
			formStyle.Render(m.form.View()),
		),
	)
}

func Show() {
	if _, err := tea.NewProgram(newModel()).Run(); err != nil {
		log.Fatal(err)
	}
}
