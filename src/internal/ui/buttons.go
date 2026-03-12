package ui

import (
	"fmt"
	"math/rand"
	"maze/internal/generator"
	"maze/internal/maze"
	"maze/internal/maze_io"
	"maze/internal/solver"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

// DefaultIOPath is the initial directory used by file open/save dialogs.
const DefaultIOPath = "."

type uiState struct {
	current *maze.Maze
	start   *maze.Point
	end     *maze.Point
}

// CreateButtons builds the control panel for the GUI.
//
// The panel includes buttons for loading/saving a maze file, generating a new
// perfect maze, solving the current maze, and clearing selection/solution.
// It wires MazeView cell taps to start/end selection.
func CreateButtons(window fyne.Window, mazeView *MazeView) *fyne.Container {
	state, status := createStateAndStatus()

	wireCellTapHandler(mazeView, state, status)

	rowsEntry, colsEntry := createDimsEntries()

	loadBtn := createLoadButton(window, mazeView, state, status)
	saveBtn := createSaveButton(window, state, status)
	genBtn := createGenerateButton(window, mazeView, state, rowsEntry, colsEntry, status)
	solveBtn := createSolveButton(window, mazeView, state, status)
	clearBtn := createClearButton(mazeView, state, status)

	return container.NewVBox(
		widget.NewLabelWithStyle("Maze controls", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		container.NewHBox(loadBtn, saveBtn),
		widget.NewSeparator(),
		widget.NewLabel("Generate a maze (Eller):"),
		container.NewHBox(rowsEntry, colsEntry, genBtn),
		widget.NewSeparator(),
		widget.NewLabel("Solve: click start (green) and end (blue), then press Solve"),
		container.NewHBox(solveBtn, clearBtn),
		widget.NewSeparator(),
		status,
	)
}

func createStateAndStatus() (*uiState, *widget.Label) {
	state := &uiState{}
	status := widget.NewLabel("Load a maze or generate a new one.")
	return state, status
}

func createDimsEntries() (*widget.Entry, *widget.Entry) {
	rowsEntry := widget.NewEntry()
	rowsEntry.SetPlaceHolder("rows (1-50)")
	colsEntry := widget.NewEntry()
	colsEntry.SetPlaceHolder("cols (1-50)")
	return rowsEntry, colsEntry
}

func wireCellTapHandler(mazeView *MazeView, state *uiState, status *widget.Label) {
	mazeView.OnCellTap = func(p maze.Point) {
		if state.current == nil {
			return
		}

		// First click => start
		if state.start == nil {
			state.start = &maze.Point{Row: p.Row, Col: p.Col}
			mazeView.SetStart(state.start)
			mazeView.ClearSolution()
			status.SetText(fmt.Sprintf("Start set to (%d,%d). Click end cell.", p.Row, p.Col))
			return
		}

		// Second click => end
		if state.end == nil {
			state.end = &maze.Point{Row: p.Row, Col: p.Col}
			mazeView.SetEnd(state.end)
			mazeView.ClearSolution()
			status.SetText(fmt.Sprintf("End set to (%d,%d). Press Solve.", p.Row, p.Col))
			return
		}

		// Third click => restart selection
		state.start = &maze.Point{Row: p.Row, Col: p.Col}
		state.end = nil
		mazeView.SetStart(state.start)
		mazeView.SetEnd(nil)
		mazeView.ClearSolution()
		status.SetText(fmt.Sprintf("Start reset to (%d,%d). Click end cell.", p.Row, p.Col))
	}
}

func createLoadButton(window fyne.Window, mazeView *MazeView, state *uiState, status *widget.Label) *widget.Button {
	return widget.NewButton("Load", func() {
		fd := dialog.NewFileOpen(func(rc fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, window)
				return
			}
			if rc == nil {
				return
			}

			path := rc.URI().Path()
			_ = rc.Close()

			m, err := maze_io.Load(path)
			if err != nil {
				dialog.ShowError(err, window)
				return
			}

			clearStateSolution(mazeView, state)
			state.current = m
			mazeView.SetMaze(m)

			status.SetText(fmt.Sprintf("Loaded a %dx%d maze.", m.Rows, m.Cols))
		}, window)

		setDefaultDialogLocation(window, fd)
		fd.Show()
	})
}

func createSaveButton(window fyne.Window, state *uiState, status *widget.Label) *widget.Button {
	return widget.NewButton("Save", func() {
		if state.current == nil {
			status.SetText("Nothing to save.")
			return
		}

		fd := dialog.NewFileSave(func(wc fyne.URIWriteCloser, err error) {
			if err != nil {
				dialog.ShowError(err, window)
				return
			}
			if wc == nil {
				return
			}

			path := wc.URI().Path()
			_ = wc.Close()

			if err := maze_io.Save(path, state.current); err != nil {
				dialog.ShowError(err, window)
				return
			}
			status.SetText("Saved maze.")
		}, window)

		setDefaultDialogLocation(window, fd)
		fd.SetFileName("maze.txt")
		fd.Show()
	})
}

func createGenerateButton(
	window fyne.Window,
	mazeView *MazeView,
	state *uiState,
	rowsEntry *widget.Entry,
	colsEntry *widget.Entry,
	status *widget.Label,
) *widget.Button {
	return widget.NewButton("Generate", func() {
		rows, cols, err := parseDims(rowsEntry.Text, colsEntry.Text)
		if err != nil {
			dialog.ShowError(err, window)
			return
		}

		rng := rand.New(rand.NewSource(time.Now().UnixNano()))
		g, err := generator.NewEller(rng)
		if err != nil {
			dialog.ShowError(err, window)
			return
		}

		m, err := g.Generate(rows, cols)
		if err != nil {
			dialog.ShowError(err, window)
			return
		}

		clearStateSolution(mazeView, state)
		state.current = m
		mazeView.SetMaze(m)

		status.SetText(fmt.Sprintf("Generated a %dx%d maze. Save it and pick start/end.", rows, cols))
	})
}

func createSolveButton(window fyne.Window, mazeView *MazeView, state *uiState, status *widget.Label) *widget.Button {
	return widget.NewButton("Solve", func() {
		if state.current == nil {
			status.SetText("Load or generate a maze first.")
			return
		}
		if state.start == nil || state.end == nil {
			status.SetText("Click start and end cells on the maze.")
			return
		}

		path, err := solver.SolveShortest(state.current, *state.start, *state.end)
		if err != nil {
			dialog.ShowError(err, window)
			return
		}

		mazeView.SetSolution(path)
		status.SetText(fmt.Sprintf("Solved: %d cells in path.", len(path)))
	})
}

func createClearButton(mazeView *MazeView, state *uiState, status *widget.Label) *widget.Button {
	return widget.NewButton("Clear", func() {
		clearStateSolution(mazeView, state)
		status.SetText("Selection cleared.")
	})
}

func clearStateSolution(mazeView *MazeView, state *uiState) {
	state.start, state.end = nil, nil
	mazeView.SetStart(nil)
	mazeView.SetEnd(nil)
	mazeView.ClearSolution()
}

func setDefaultDialogLocation(window fyne.Window, fd *dialog.FileDialog) {
	uri := storage.NewFileURI(DefaultIOPath)
	location, err := storage.ListerForURI(uri)
	if err != nil {
		dialog.ShowError(err, window)
		return
	}
	fd.SetLocation(location)
}

func parseDims(rowsText, colsText string) (int, int, error) {
	rows, err := strconv.Atoi(rowsText)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid rows")
	}
	cols, err := strconv.Atoi(colsText)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid cols")
	}
	if rows <= 0 || cols <= 0 {
		return 0, 0, fmt.Errorf("rows and cols must be positive")
	}
	if rows > maze.MaxRows || cols > maze.MaxCols {
		return 0, 0, fmt.Errorf("max size is %dx%d", maze.MaxRows, maze.MaxCols)
	}
	return rows, cols, nil
}
