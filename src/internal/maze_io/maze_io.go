// Package maze_io implements loading and saving mazes in the project file format.
//
// The format starts with: "<rows> <cols>", followed by two integer matrices:
//  1. the right-wall matrix (Rows x Cols)
//  2. the bottom-wall matrix (Rows x Cols)
//
// Matrix values are 0 (no wall) or 1 (wall). Whitespace is used as a separator;
// empty lines are ignored.
package maze_io

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"maze/internal/maze"
)

const (
	wallAbsent  = 0
	wallPresent = 1
)

// Load reads a maze in the given format.
func Load(path string) (*maze.Maze, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Decode(f)
}

// Save writes a maze in the given format.
func Save(path string, m *maze.Maze) error {
	if err := m.Validate(); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	if err := Encode(w, m); err != nil {
		return err
	}
	return w.Flush()
}

// Decode reads from r using whitespace-separated ints.
func Decode(r io.Reader) (*maze.Maze, error) {
	in := bufio.NewReader(r)
	var rows, cols int
	if _, err := fmt.Fscan(in, &rows, &cols); err != nil {
		return nil, fmt.Errorf("read dimensions: %w", err)
	}

	m, err := maze.New(rows, cols)
	if err != nil {
		return nil, err
	}

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if _, err := fmt.Fscan(in, &m.RightWalls[i][j]); err != nil {
				return nil, fmt.Errorf("read right walls: %w", err)
			}
		}
	}
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if _, err := fmt.Fscan(in, &m.BottomWalls[i][j]); err != nil {
				return nil, fmt.Errorf("read bottom walls: %w", err)
			}
		}
	}

	if err := validateWalls(m); err != nil {
		return nil, err
	}
	return m, nil
}

func validateWalls(m *maze.Maze) error {
	if err := m.Validate(); err != nil {
		return err
	}
	for r := 0; r < m.Rows; r++ {
		for c := 0; c < m.Cols; c++ {
			if m.RightWalls[r][c] != wallAbsent && m.RightWalls[r][c] != wallPresent {
				return fmt.Errorf("invalid right wall value at (%d,%d)", r, c)
			}
			if m.BottomWalls[r][c] != wallAbsent && m.BottomWalls[r][c] != wallPresent {
				return fmt.Errorf("invalid bottom wall value at (%d,%d)", r, c)
			}
		}
	}
	return nil
}

// Encode writes to w in the given format.
func Encode(w io.Writer, m *maze.Maze) error {
	if err := m.Validate(); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "%d %d\n", m.Rows, m.Cols); err != nil {
		return err
	}
	writeWalls := func(walls [][]int) error {
		for r := 0; r < m.Rows; r++ {
			for c := 0; c < m.Cols; c++ {
				sep := " "
				if c == m.Cols-1 {
					sep = "\n"
				}
				if _, err := fmt.Fprintf(w, "%d%s", walls[r][c], sep); err != nil {
					return err
				}
			}
		}
		return nil
	}
	if err := writeWalls(m.RightWalls); err != nil {
		return err
	}
	if err := writeWalls(m.BottomWalls); err != nil {
		return err
	}
	return nil
}
