package main

import (
	"fmt"
	"os"

	"github.com/fprasx/secrets-and-spies/ui/board"
)

func writeGraphviz(edges map[int][]int, filename string, cities []board.City) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Fprintln(file, "graph G {")
	for from, toList := range edges {
		fmt.Fprintf(file, "  %d [label = \"%v %v\"];\n", from, from, cities[from].Name)
		for _, to := range toList {
			fmt.Fprintf(file, "  %d -- %d;\n", from, to)
		}
	}
	fmt.Fprintln(file, "}")
	return nil
}

func main() {
	cities := []board.City{
		{Name: "Fugging, Austria", Color: "#ff0000"},
		{Name: "Poo, Spain", Color: "#00ff00"},
		{Name: "Bastardo, Italy", Color: "#0000ff"},
		{Name: "Condom, France", Color: "#ff0000"},
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
	edges := map[int][]int{
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
	activeLocation := board.City{Name: "Poo, Spain", Color: "#00ff00"}
	// b := board.NewBoard(cities, edges, activeLocation)
	// writeGraphviz(edges, "graph.dot", cities)
	board.Show(
		cities, edges, activeLocation,
	)
}
