package generator

import (
	"math/rand"
	"strconv"
	"testing"

	"maze/internal/maze"
)

func TestEllerGenerate(t *testing.T) {
	rng := rand.New(rand.NewSource(42))

	eller, ellerErr := NewEller(rng)
	if ellerErr != nil {
		t.Fatalf("NewEller: %v", ellerErr)
	}

	_, genErr := eller.Generate(7, 7)
	if genErr != nil {
		t.Fatalf("Generate: %v", genErr)
	}
}

func TestEllerGenerate_PerfectProperties(t *testing.T) {
	cases := []struct {
		rows int
		cols int
		seed int64
	}{
		{1, 1, 1},
		{10, 10, 2},
		{10, 10, 3},
		{2, 2, 4},
		{7, 7, 5},
		{50, 50, 6},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(fmtName(tc.rows, tc.cols), func(t *testing.T) {
			rng := rand.New(rand.NewSource(tc.seed))
			g, err := NewEller(rng)
			if err != nil {
				t.Fatalf("NewEller: %v", err)
			}
			m, err := g.Generate(tc.rows, tc.cols)
			if err != nil {
				t.Fatalf("Generate: %v", err)
			}
			if err := m.Validate(); err != nil {
				t.Fatalf("Validate: %v", err)
			}
			if !isConnected(m) {
				t.Fatalf("maze is not connected")
			}
			edges := edgeCount(m)
			cells := tc.rows * tc.cols
			if edges != cells-1 && cells > 0 {
				t.Fatalf("expected edges=%d, got %d (not a tree)", cells-1, edges)
			}
		})
	}
}

func fmtName(r, c int) string { return "r" + strconv.Itoa(r) + "c" + strconv.Itoa(c) }

func edgeCount(m *maze.Maze) int {
	edges := 0
	for r := 0; r < m.Rows; r++ {
		for c := 0; c < m.Cols; c++ {
			if c < m.Cols-1 && m.RightWalls[r][c] == 0 {
				edges++
			}
			if r < m.Rows-1 && m.BottomWalls[r][c] == 0 {
				edges++
			}
		}
	}
	return edges
}

func isConnected(m *maze.Maze) bool {
	cells := m.Rows * m.Cols
	if cells == 0 {
		return true
	}

	visited := make([]bool, cells)
	cols := m.Cols
	idx := func(p maze.Point) int { return p.Row*cols + p.Col }

	q := make([]maze.Point, 0, cells)
	q = append(q, maze.Point{Row: 0, Col: 0})
	visited[idx(q[0])] = true

	for head := 0; head < len(q); head++ {
		p := q[head]
		for _, nb := range neighbors(m, p) {
			i := idx(nb)
			if visited[i] {
				continue
			}
			visited[i] = true
			q = append(q, nb)
		}
	}

	for _, v := range visited {
		if !v {
			return false
		}
	}
	return true
}

func neighbors(m *maze.Maze, p maze.Point) []maze.Point {
	res := make([]maze.Point, 0, 4)
	if p.Col > 0 && m.RightWalls[p.Row][p.Col-1] == 0 {
		res = append(res, maze.Point{Row: p.Row, Col: p.Col - 1})
	}
	if p.Col < m.Cols-1 && m.RightWalls[p.Row][p.Col] == 0 {
		res = append(res, maze.Point{Row: p.Row, Col: p.Col + 1})
	}
	if p.Row > 0 && m.BottomWalls[p.Row-1][p.Col] == 0 {
		res = append(res, maze.Point{Row: p.Row - 1, Col: p.Col})
	}
	if p.Row < m.Rows-1 && m.BottomWalls[p.Row][p.Col] == 0 {
		res = append(res, maze.Point{Row: p.Row + 1, Col: p.Col})
	}
	return res
}
