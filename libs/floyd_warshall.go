package libs

import (
	"math"
)

// Edge represents directed and weighted edge in graph
type Edge struct {
	From   int
	To     int
	Weight float64
}

// FloydWarshall finds shortest path between all pairs of vertices (all-pairs shortest path, APSP)
// using Floyd-Warshall algorithm
func FloydWarshall(
	n int,
	edges []Edge,
) ([][]float64, [][]int) {
	dist := make([][]float64, n)
	next := make([][]int, n)

	// Initialize dist and next
	for i := 0; i < n; i++ {
		dist[i] = make([]float64, n)
		next[i] = make([]int, n)

		for j := 0; j < n; j++ {
			dist[i][j] = math.Inf(1)
			next[i][j] = -1
		}
		dist[i][i] = 0
	}

	// Insert edges
	for _, e := range edges {
		dist[e.From][e.To] = e.Weight
		next[e.From][e.To] = e.To
	}

	// Floyd-Warshall algorithm
	for middle := 0; middle < n; middle++ {
		for start := 0; start < n; start++ {
			for end := 0; end < n; end++ {
				shortcutDist := dist[start][middle] + dist[middle][end]
				if shortcutDist < dist[start][end] {
					dist[start][end] = shortcutDist
					next[start][end] = next[start][middle]
				}
			}
		}
	}

	return dist, next
}
