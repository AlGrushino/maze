package solver

import (
	"math/rand"
	"testing"

	"maze/internal/generator"
	"maze/internal/maze"
)

func TestSolveShortest_Trivial(t *testing.T) {
	m, err := maze.New(1, 1)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	makeBorders(m)
	path, err := SolveShortest(m, maze.Point{Row: 0, Col: 0}, maze.Point{Row: 0, Col: 0})
	if err != nil {
		t.Fatalf("SolveShortest: %v", err)
	}
	if len(path) != 1 || path[0] != (maze.Point{Row: 0, Col: 0}) {
		t.Fatalf("unexpected path: %#v", path)
	}
}

func TestSolveShortest_NoPath(t *testing.T) {
	m, err := maze.New(1, 2)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	makeBorders(m)
	// Ensure wall between cells.
	m.RightWalls[0][0] = 1
	path, err := SolveShortest(m, maze.Point{Row: 0, Col: 0}, maze.Point{Row: 0, Col: 1})
	if err == nil {
		t.Fatalf("expected error, got path=%v", path)
	}
}

func TestSolveShortest_ValidOnPerfectMaze(t *testing.T) {
	rng := rand.New(rand.NewSource(7))
	g, err := generator.NewEller(rng)
	if err != nil {
		t.Fatalf("NewEller: %v", err)
	}
	m, err := g.Generate(10, 10)
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}

	start := maze.Point{Row: 0, Col: 0}
	end := maze.Point{Row: 9, Col: 9}

	path, err := SolveShortest(m, start, end)
	if err != nil {
		t.Fatalf("SolveShortest: %v", err)
	}
	if path[0] != start || path[len(path)-1] != end {
		t.Fatalf("path endpoints mismatch")
	}
	if !pathRespectsWalls(m, path) {
		t.Fatalf("path crosses walls")
	}
}

func makeBorders(m *maze.Maze) {
	for r := 0; r < m.Rows; r++ {
		for c := 0; c < m.Cols; c++ {
			m.RightWalls[r][c] = 1
			if r == m.Rows-1 {
				m.BottomWalls[r][c] = 1
			}
		}
	}
}

func pathRespectsWalls(m *maze.Maze, path []maze.Point) bool {
	if len(path) < 2 {
		return true
	}
	for i := 0; i < len(path)-1; i++ {
		a, b := path[i], path[i+1]
		dr := b.Row - a.Row
		dc := b.Col - a.Col

		switch {
		case dr == 0 && dc == 1:
			if m.RightWalls[a.Row][a.Col] != 0 {
				return false
			}
		case dr == 0 && dc == -1:
			if m.RightWalls[a.Row][a.Col-1] != 0 {
				return false
			}
		case dr == 1 && dc == 0:
			if m.BottomWalls[a.Row][a.Col] != 0 {
				return false
			}
		case dr == -1 && dc == 0:
			if m.BottomWalls[a.Row-1][a.Col] != 0 {
				return false
			}
		default:
			return false
		}
	}
	return true
}
