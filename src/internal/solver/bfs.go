// Package solver implements maze solving algorithms.
//
// The solver finds the shortest route between two cells, where the route length
// is measured as the number of visited cells. Movement is allowed only between
// orthogonally adjacent cells not separated by a wall.
package solver

import (
	"fmt"
	"slices"

	"maze/internal/maze"
)

type node struct {
	p    maze.Point
	prev int
}

// SolveShortest finds the shortest path (fewest cells) using BFS.
// It returns the path including start and end.
func SolveShortest(m *maze.Maze, start, end maze.Point) ([]maze.Point, error) {
	if err := m.Validate(); err != nil {
		return nil, err
	}
	if !m.InBounds(start) || !m.InBounds(end) {
		return nil, fmt.Errorf("start or end out of bounds")
	}

	if start == end {
		return []maze.Point{start}, nil
	}

	rows, cols := m.Rows, m.Cols
	idx := func(p maze.Point) int { return p.Row*cols + p.Col }

	visited := make([]bool, rows*cols)
	queue := make([]node, 0, rows*cols)
	queue = append(queue, node{p: start, prev: -1})
	visited[idx(start)] = true

	for head := 0; head < len(queue); head++ {
		cur := queue[head]
		for _, nb := range neighbors(m, cur.p) {
			i := idx(nb)
			if visited[i] {
				continue
			}
			visited[i] = true
			queue = append(queue, node{p: nb, prev: head})
			if nb == end {
				return reconstruct(queue, len(queue)-1), nil
			}
		}
	}

	return nil, fmt.Errorf("no path found")
}

func neighbors(m *maze.Maze, p maze.Point) []maze.Point {
	res := make([]maze.Point, 0, 4)

	// Left: check right wall of left cell.
	if p.Col > 0 && m.RightWalls[p.Row][p.Col-1] == 0 {
		res = append(res, maze.Point{Row: p.Row, Col: p.Col - 1})
	}
	// Right: check right wall of this cell.
	if p.Col < m.Cols-1 && m.RightWalls[p.Row][p.Col] == 0 {
		res = append(res, maze.Point{Row: p.Row, Col: p.Col + 1})
	}
	// Up: check bottom wall of upper cell.
	if p.Row > 0 && m.BottomWalls[p.Row-1][p.Col] == 0 {
		res = append(res, maze.Point{Row: p.Row - 1, Col: p.Col})
	}
	// Down: check bottom wall of this cell.
	if p.Row < m.Rows-1 && m.BottomWalls[p.Row][p.Col] == 0 {
		res = append(res, maze.Point{Row: p.Row + 1, Col: p.Col})
	}

	return res
}

func reconstruct(queue []node, endIdx int) []maze.Point {
	path := make([]maze.Point, 0)
	for i := endIdx; i >= 0; i = queue[i].prev {
		path = append(path, queue[i].p)
		if queue[i].prev == -1 {
			break
		}
	}
	slices.Reverse(path)
	return path
}
