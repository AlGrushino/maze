// Package maze defines the core maze data structures.
//
// A maze is represented as a grid of cells with "thin walls" between cells.
// The maze is stored using two wall matrices:
//   - RightWalls[r][c] indicates a wall between (r,c) and (r,c+1)
//   - BottomWalls[r][c] indicates a wall between (r,c) and (r+1,c)
//
// Outer borders are represented by setting the appropriate right/bottom walls
// on the last column/row.
package maze

import "fmt"

// MaxRows is the maximum supported maze height (number of rows).
const MaxRows = 50

// MaxCols is the maximum supported maze width (number of columns).
const MaxCols = 50

// Point is a maze cell coordinate.
// Row and Col are zero-based.
type Point struct {
	Row int
	Col int
}

// Maze is a thin-wall maze.
//
// RightWalls[r][c] indicates the wall to the right of cell (r,c).
// BottomWalls[r][c] indicates the wall at the bottom of cell (r,c).
// Values are 0 (no wall) or 1 (wall).
type Maze struct {
	Rows int
	Cols int

	RightWalls  [][]int
	BottomWalls [][]int
}

// New creates an empty maze with correct outer borders.
// Internal walls are initialized to 0.
func New(rows, cols int) (*Maze, error) {
	if rows <= 0 || cols <= 0 {
		return nil, fmt.Errorf("rows and cols must be positive")
	}
	if rows > MaxRows || cols > MaxCols {
		return nil, fmt.Errorf("maze size exceeds %dx%d", MaxRows, MaxCols)
	}

	right := make([][]int, rows)
	bottom := make([][]int, rows)
	for r := 0; r < rows; r++ {
		right[r] = make([]int, cols)
		bottom[r] = make([]int, cols)
	}

	return &Maze{Rows: rows, Cols: cols, RightWalls: right, BottomWalls: bottom}, nil
}

// InBounds reports whether p is inside the maze.
func (m *Maze) InBounds(p Point) bool {
	return p.Row >= 0 && p.Row < m.Rows && p.Col >= 0 && p.Col < m.Cols
}

// Validate checks dimensions, matrix sizes, and values.
func (m *Maze) Validate() error {
	if m == nil {
		return fmt.Errorf("maze is nil")
	}
	if m.Rows <= 0 || m.Cols <= 0 {
		return fmt.Errorf("invalid dimensions")
	}
	if m.Rows != m.Cols {
		return fmt.Errorf("maze must have same rows and cols")
	}
	if m.Rows > MaxRows || m.Cols > MaxCols {
		return fmt.Errorf("maze size exceeds %dx%d", MaxRows, MaxCols)
	}
	if len(m.RightWalls) != m.Rows || len(m.BottomWalls) != m.Rows {
		return fmt.Errorf("invalid wall matrix row count")
	}
	for r := 0; r < m.Rows; r++ {
		if len(m.RightWalls[r]) != m.Cols || len(m.BottomWalls[r]) != m.Cols {
			return fmt.Errorf("invalid wall matrix col count at row %d", r)
		}
		for c := 0; c < m.Cols; c++ {
			if m.RightWalls[r][c] != 0 && m.RightWalls[r][c] != 1 {
				return fmt.Errorf("invalid right wall value at (%d,%d)", r, c)
			}
			if m.BottomWalls[r][c] != 0 && m.BottomWalls[r][c] != 1 {
				return fmt.Errorf("invalid bottom wall value at (%d,%d)", r, c)
			}
		}
		if m.RightWalls[r][m.Cols-1] != 1 {
			return fmt.Errorf("right border must be a wall")
		}
	}
	for c := 0; c < m.Cols; c++ {
		if m.BottomWalls[m.Rows-1][c] != 1 {
			return fmt.Errorf("bottom border must be a wall")
		}
	}
	return nil
}
