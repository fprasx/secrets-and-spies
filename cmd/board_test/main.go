package main

import (
	"github.com/fprasx/secrets-and-spies/ui/board"
)

func main() {
	board.Show(
		[]board.City{
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
		},
		map[int][]int{
			0:  {1},
			1:  {13},
			2:  {1, 2, 3, 13, 14},
			3:  {4, 9, 12},
			4:  {0, 4},
			5:  {2, 5},
			6:  {1},
			7:  {3, 6},
			8:  {1},
			9:  {1, 9},
			10: {10, 11},
			11: {2},
			12: {14},
			13: {0, 9, 13},
			14: {3, 6, 9, 12, 13},
			15: {1, 11, 14},
		},
		board.City{Name: "Poo, Spain", Color: "#00ff00"},
	)
}
