package main

import "fmt"

type empty struct{}

type GraphMtx struct {
	adjMatrix [][]int
}

func (g *GraphMtx) BFS(startVertex int) []int {
	visited := make(map[int]empty, len(g.adjMatrix))
	var queue []int
	path := make([]int, 0, len(g.adjMatrix))

	visited[startVertex] = empty{}
	queue = append(queue, startVertex)

	for len(queue) > 0 {
		currentVertex := queue[0]
		queue = queue[1:]
		path = append(path, currentVertex)

		for i := range g.adjMatrix {
			if _, ok := visited[i]; !ok {
				visited[i] = empty{}
				queue = append(queue, i)
			}
		}
	}

	return path
}

func main() {
	var adjMatrix = [][]int{
		{0, 1, 0, 1, 0},
		{1, 0, 1, 0, 1},
		{0, 1, 0, 1, 0},
		{1, 0, 1, 0, 1},
		{0, 1, 0, 1, 0},
	}
	graph := GraphMtx{adjMatrix}
	path := graph.BFS(0)
	fmt.Println(path)
}
