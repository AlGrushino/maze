# Maze — Go Maze Generator & Solver

A cross-platform desktop application for generating, rendering, solving, and saving perfect mazes. Built with Go and [Fyne](https://fyne.io/) GUI toolkit.

## Overview

Maze provides a graphical interface to:
- **Load** mazes from text files
- **Generate** perfect mazes using Eller's algorithm (up to 50x50)
- **Render** mazes in a 500×500 pixel canvas with interactive cell selection
- **Solve** mazes via BFS shortest-path search
- **Save** generated or loaded mazes to file

## Features

- Eller's algorithm guarantees a perfect maze (spanning tree — exactly one path between any two cells, no loops, no isolated areas)
- BFS solver finds the shortest path by cell count
- Interactive GUI: click cells to set start/end points
- File I/O: load and save mazes in a compact text format
- Responsive wall rendering with 2px-thick lines
- Unit tests for both generator and solver

## Project Structure

```
src/
├── cmd/maze_gui/main.go          # Application entry point
├── internal/
│   ├── maze/maze.go              # Core data model (Maze struct, validation)
│   ├── maze_io/maze_io.go        # File load/save (encode/decode)
│   ├── generator/eller.go        # Eller's maze generation algorithm
│   ├── solver/bfs.go             # BFS shortest-path solver
│   └── ui/
│       ├── buttons.go            # Control panel (load/save/generate/solve/clear)
│       └── maze_view.go          # Canvas widget for maze rendering
├── mazes/maze.txt                # Sample 4×4 maze
├── Makefile                      # Build targets
└── go.mod                        # Go module definition
```

## Prerequisites

- **Go** 1.22+
- **Linux:** system libraries for OpenGL and GLFW:
  ```bash
  sudo apt install libgl-dev libglfw3-dev libxcursor-dev libxrandr-dev libxinerama-dev libxi-dev
  ```

## Build & Run

```bash
cd src

# Build
make

# Run
./build/maze
```

## Usage

| Action | How |
|---|---|
| **Generate** | Enter rows & cols (1–50), click **Generate** |
| **Load** | Click **Load**, select a `.txt` maze file |
| **Solve** | Click **Solve**, then click a start cell (green), then an end cell (blue). The shortest path is drawn in crimson. |
| **Clear** | Click **Clear** to reset selection and solution |
| **Save** | Click **Save** to write the current maze to a file |

## Maze File Format

```
<rows> <cols>

<right_walls: rows × cols, space-separated 0/1>

<bottom_walls: rows × cols, space-separated 0/1>
```

**Example** (4×4 maze):

```
4 4
0 1 0 1
1 0 1 1
1 1 0 1
0 0 0 1

1 0 0 1
0 0 0 0
0 0 1 1
1 1 1 1
```

- `1` = wall present, `0` = passage
- `RightWalls[r][c]` — wall to the right of cell (r,c)
- `BottomWalls[r][c]` — wall below cell (r,c)
- Last column always has right walls (`1`), last row always has bottom walls (`1`)

## Algorithms

### Eller's Algorithm (Generation)

Generates a perfect maze row-by-row:
1. Each row starts with unassigned cells receiving unique set IDs
2. Adjacent cells in the same set get a wall (prevents loops); different sets are randomly merged
3. At least one cell per set keeps an open downward passage (ensures connectivity)
4. Last row: all different-set neighbors are forcibly merged
5. Result: a spanning tree with exactly `rows × cols − 1` edges

### BFS (Solving)

Standard breadth-first search:
- Explores all 4 neighbors per cell, respecting walls
- Tracks parent pointers for path reconstruction
- Guarantees shortest path by cell count

## Tests

```bash
make tests
# or
go test ./... -cover
```

Tests cover:
- **Generator** (`eller_test.go`): maze dimensions, connectivity, spanning tree properties
- **Solver** (`bfs_test.go`): path correctness, unreachable cells, edge cases

## Makefile Targets

| Target | Description |
|---|---|
| `make` / `make all` | Build to `build/maze` |
| `make install` | Install to `bin/maze` |
| `make uninstall` | Remove installed binary |
| `make clean` | Remove build artifacts and test cache |
| `make rebuild` | Clean + build |
| `make dist` | Create `dist/maze.tar.gz` archive |
| `make dvi` | Generate API docs to `docs/api.txt` |
| `make tests` | Run all tests with coverage |

## Tech Stack

- **Go 1.22** — core language
- **Fyne v2.5.4** — cross-platform GUI toolkit
- **Standard library only** for algorithms (no external deps for generation/solving)

## License

See [LICENSE](LICENSE).
