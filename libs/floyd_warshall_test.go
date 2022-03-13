package libs

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func appendReversedEdges(edges []Edge) []Edge {
	n := len(edges)
	for i := 0; i < n; i++ {
		edges = append(edges, Edge{From: edges[i].To, To: edges[i].From, Weight: edges[i].Weight})
	}
	return edges
}

func TestFindAllPairsShortestPath(t *testing.T) {
	t.Run("undirected graph with 3 vertices and 2 edges", func(t *testing.T) {
		n := 3
		edges := []Edge{
			{From: 0, To: 1, Weight: 1},
			{From: 1, To: 2, Weight: 1},
		}

		// Make undirected
		edges = appendReversedEdges(edges)

		dist, next := FloydWarshall(n, edges)

		assert.Equal(t, [][]float64{
			{0, 1, 2},
			{1, 0, 1},
			{2, 1, 0},
		}, dist, "dist should match")

		assert.Equal(t, [][]int{
			{-1, 1, 1},
			{0, -1, 2},
			{1, 1, -1},
		}, next, "next should match")
	})

	t.Run("undirected graph with 8 vertices and 12 weighted edges", func(t *testing.T) {
		n := 8
		edges := []Edge{
			{From: 0, To: 1, Weight: 4},
			{From: 0, To: 2, Weight: 3},
			{From: 0, To: 4, Weight: 7},
			{From: 1, To: 2, Weight: 6},
			{From: 1, To: 3, Weight: 5},
			{From: 2, To: 3, Weight: 11},
			{From: 2, To: 4, Weight: 8},
			{From: 3, To: 4, Weight: 2},
			{From: 3, To: 5, Weight: 2},
			{From: 3, To: 6, Weight: 10},
			{From: 4, To: 6, Weight: 5},
			{From: 5, To: 6, Weight: 3},
		}

		edges = appendReversedEdges(edges)

		dist, next := FloydWarshall(n, edges)

		assert.Equal(t, 0., dist[0][0])
		assert.Equal(t, 4., dist[0][1])
		assert.Equal(t, 3., dist[0][2])
		assert.Equal(t, 9., dist[0][3])
		assert.Equal(t, 7., dist[0][4])
		assert.Equal(t, 11., dist[0][5])
		assert.Equal(t, 12., dist[0][6])
		assert.Equal(t, math.Inf(1), dist[0][7])

		assert.Equal(t, -1, next[0][0])
		assert.Equal(t, 1, next[0][1])
		assert.Equal(t, 2, next[0][2])
		assert.Equal(t, 1, next[0][3])
		assert.Equal(t, 3, next[1][3])
		assert.Equal(t, 1, next[0][5])
		assert.Equal(t, 3, next[1][5])
		assert.Equal(t, 5, next[3][5])
		assert.Equal(t, 4, next[0][6])
		assert.Equal(t, 6, next[4][6])
		assert.Equal(t, -1, next[0][7])
	})
}
