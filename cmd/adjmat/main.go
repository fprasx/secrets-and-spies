package main

import (
	"fmt"
)

var (
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

func createAdjacencyMatrix(edges map[int][]int) [][]int {
	n := len(edges)
	adj := make([][]int, n)
	for i := range n {
		adj[i] = make([]int, n)
	}
	for i := range n {
		adj[i][i] = 1
	}
	for i := range n {
		for _, neighbor := range edges[i] {
			adj[i][neighbor] = 1
			adj[neighbor][i] = 1
		}
	}
	return adj
}

func main() {
	adj := createAdjacencyMatrix(edges)
	fmt.Printf("%v\n", adj)
}
