package board

import (
	"log"
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type City struct {
	Name  string
	Color lipgloss.Color
}

type Board struct {
	activeLocation int
	cities         []City
	edges          map[int][]int

	maxNameLen int
}

func (b *Board) nextCity() int {
	return (b.activeLocation + 1) % len(b.cities)
}

func (b *Board) prevCity() int {
	return (b.activeLocation - 1 + len(b.cities)) % len(b.cities)
}

// item must be in slice
func indexOf[T comparable](slice []T, item T) int {
	for i, v := range slice {
		if v == item {
			return i
		}
	}
	panic("item was not in slice")
}

func longestCityName(cities []City) int {
	maxLen := 0
	for _, c := range cities {
		if len(c.Name) > maxLen {
			maxLen = len(c.Name)
		}
	}
	return maxLen
}

// initialLocation must be one of the specified cities
func NewBoard(cities []City, edges map[int][]int, initialLocation City) Board {
	return Board{
		activeLocation: indexOf(cities, initialLocation),
		cities:         cities,
		edges:          edges,
		maxNameLen:     longestCityName(cities),
	}
}

func (board Board) Init() tea.Cmd {
	return nil
}

func (board Board) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
	case tea.KeyMsg:
		switch msg.String() {
		case "right":
			board.activeLocation = board.nextCity()
		case "left":
			board.activeLocation = board.prevCity()
		case "ctrl+c", "esc", "q":
			return board, tea.Interrupt
		}
	}
	return board, nil
}

func PadRight(s string, length int) string {
	if len(s) >= length {
		return s
	}
	return s + strings.Repeat(" ", length-len(s))
}

func (b Board) activeCity() string {
	return b.cities[b.activeLocation].Name
}

func (board Board) View() string {
	var rows []string
	for i := 0; i < 16; i += 4 {
		var cells []string
		for j := 0; j < 4; j += 1 {
			cellStyle := lipgloss.NewStyle().
				Width(board.maxNameLen+2).
				Align(lipgloss.Center).
				Border(lipgloss.RoundedBorder()).
				Padding(1, 0).
				Margin(2, 2)

			if i+j == board.activeLocation {
				cellStyle = cellStyle.
					BorderForeground(lipgloss.Color("#00ff00")).
					Foreground(lipgloss.Color("#00ff00"))
			} else if slices.Contains(board.edges[board.activeLocation], i+j) {
				cellStyle = cellStyle.
					BorderForeground(lipgloss.Color("#0000ff")).
					Foreground(lipgloss.Color("#0000ff"))

			}
			cells = append(cells, cellStyle.Render(board.cities[i+j].Name))
		}
		rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, cells...))
	}
	grid := lipgloss.JoinVertical(lipgloss.Left, rows...)
	return grid
}

func Show(cities []City, edges map[int][]int, initialLocation City) {
	if _, err := tea.NewProgram(NewBoard(cities, edges, initialLocation)).Run(); err != nil {
		log.Fatal(err)
	}
}
