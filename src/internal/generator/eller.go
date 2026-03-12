// Package generator implements algorithms for generating perfect mazes.
//
// The project uses Eller's algorithm to generate "thin-wall" mazes such that
// every pair of cells is connected by exactly one simple path (i.e., the maze
// has no isolated areas and no loops).
package generator

import (
	"fmt"
	"math/rand"

	"maze/internal/maze"
)

const (
	wallAbsent  = 0
	wallPresent = 1
)

// Eller generates perfect mazes using Eller's algorithm.
type Eller struct {
	rng *rand.Rand
}

// NewEller creates an Eller generator.
func NewEller(rng *rand.Rand) (*Eller, error) {
	if rng == nil {
		return nil, fmt.Errorf("rng is nil")
	}
	return &Eller{rng: rng}, nil
}

// Generate returns a perfect maze of the given size.
func (g *Eller) Generate(rows, cols int) (*maze.Maze, error) {
	m, err := maze.New(rows, cols)
	if err != nil {
		return nil, err
	}

	setID := make([]int, cols)
	nextSet := 1

	for row := 0; row < rows; row++ {
		assignSets(setID, &nextSet)

		if row == rows-1 {
			g.lastRow(m, row, setID)
			break
		}

		g.makeRightWalls(m, row, setID)
		g.makeBottomWalls(m, row, setID)
		prepareNextRow(setID, m.BottomWalls[row])
	}

	if err := m.Validate(); err != nil {
		return nil, err
	}
	return m, nil
}

func assignSets(setID []int, nextSet *int) {
	for c := 0; c < len(setID); c++ {
		if setID[c] == 0 {
			setID[c] = *nextSet
			*nextSet++
		}
	}
}

func (g *Eller) makeRightWalls(m *maze.Maze, r int, setID []int) {
	for c := 0; c < m.Cols-1; c++ {
		if setID[c] == setID[c+1] {
			m.RightWalls[r][c] = wallPresent
			continue
		}

		if g.rng.Intn(2) == 0 {
			m.RightWalls[r][c] = wallPresent
			continue
		}

		m.RightWalls[r][c] = wallAbsent
		mergeSets(setID, setID[c+1], setID[c])
	}
	m.RightWalls[r][m.Cols-1] = wallPresent
}

func mergeSets(setID []int, from, to int) {
	if from == to {
		return
	}
	for i := 0; i < len(setID); i++ {
		if setID[i] == from {
			setID[i] = to
		}
	}
}

func (g *Eller) makeBottomWalls(m *maze.Maze, r int, setID []int) {
	indicesBySet := map[int][]int{}
	for c := 0; c < m.Cols; c++ {
		indicesBySet[setID[c]] = append(indicesBySet[setID[c]], c)
	}

	for _, indices := range indicesBySet {
		open := make([]bool, len(indices))
		openCount := 0
		for i := range indices {
			if g.rng.Intn(2) == 0 {
				open[i] = true
				openCount++
			}
		}
		if openCount == 0 {
			open[g.rng.Intn(len(indices))] = true
		}

		for i, c := range indices {
			if open[i] {
				m.BottomWalls[r][c] = wallAbsent
			} else {
				m.BottomWalls[r][c] = wallPresent
			}
		}
	}
}

func prepareNextRow(setID []int, bottomRow []int) {
	for c := 0; c < len(setID); c++ {
		if bottomRow[c] == wallPresent {
			setID[c] = 0
		}
	}
}

func (g *Eller) lastRow(m *maze.Maze, r int, setID []int) {
	for c := 0; c < m.Cols-1; c++ {
		if setID[c] == setID[c+1] {
			m.RightWalls[r][c] = wallPresent
			continue
		}
		m.RightWalls[r][c] = wallAbsent
		mergeSets(setID, setID[c+1], setID[c])
	}
	m.RightWalls[r][m.Cols-1] = wallPresent

	for c := 0; c < m.Cols; c++ {
		m.BottomWalls[r][c] = wallPresent
	}
}
