package board

import (
	"fmt"
	"log"
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/tree"
)

type City struct {
	Name  string
	Color lipgloss.Color
}

type Board struct {
	// Game state info
	activeLocation int
	cities         []City
	edges          map[int][]int

	// Precomputed info for rendering
	maxNameLen int

	// Interactive board info
	currentSelection int
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
	activeLocation := indexOf(cities, initialLocation)
	return Board{
		activeLocation:   activeLocation,
		cities:           cities,
		edges:            edges,
		maxNameLen:       longestCityName(cities),
		currentSelection: edges[activeLocation][0],
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

func (board Board) CreateTree() *tree.Tree {
	visited := make(map[int]bool)
	queue := []int{board.activeLocation}
	parent := make(map[int]int)
	order := []int{}
	treebuilder := map[int]*tree.Tree{board.activeLocation: tree.Root(board.activeCity())}

	visited[board.activeLocation] = true

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]

		order = append(order, node)
		for _, neighbor := range board.edges[node] {
			if !visited[neighbor] {
				visited[neighbor] = true
				parent[neighbor] = node
				queue = append(queue, neighbor)
				newtree := tree.Root(board.cities[neighbor].Name)
				treebuilder[node].Child(newtree)
				treebuilder[neighbor] = newtree
			}
		}
	}
	return treebuilder[board.activeLocation]
}

func (board Board) View() string {
	fmt.Printf("")
	var rows []string
	for i := 0; i < 16; i += 4 {
		var cells []string
		for j := 0; j < 4; j += 1 {
			cellStyle := lipgloss.NewStyle().
				Width(board.maxNameLen).
				Align(lipgloss.Center).
				Border(lipgloss.RoundedBorder()).
				Margin(0, 1).MarginBottom(1)

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
	grid = lipgloss.JoinHorizontal(lipgloss.Left, grid, fmt.Sprintf("%v", board.CreateTree()))
	return grid
}

func Show(cities []City, edges map[int][]int, initialLocation City) {
	if _, err := tea.NewProgram(NewBoard(cities, edges, initialLocation)).Run(); err != nil {
		log.Fatal(err)
	}
}
