package menu

import (
	_ "embed"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

var (
	HostAddress string
)

type joinModel struct {
	form *huh.Form
}

func newJoinModel() joinModel {
	m := joinModel{
		form: huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Value(&HostAddress).
					Title("Host Address").
					Prompt("").
					Description("Address of host to connect to").
					Placeholder("unix:///tmp/spies/socket"),
			),
		).
			WithWidth(minWidth).
			WithTheme(huh.ThemeCatppuccin()).
			WithShowHelp(false).
			WithShowErrors(true),
	}

	return m
}

func (m joinModel) Init() tea.Cmd {
	return m.form.Init()
}

func (m joinModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m joinModel) View() string {
	if m.form.State == huh.StateCompleted {
		return ""
	}

	return appStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Center,
			headerStyle.Render(banner),
			formStyle.Render(m.form.View()),
		),
	)
}

func ShowJoin() {
	if _, err := tea.NewProgram(newJoinModel()).Run(); err != nil {
		log.Fatal(err)
	}
}
