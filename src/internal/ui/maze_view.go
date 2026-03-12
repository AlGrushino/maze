// Package ui contains the GUI widgets and rendering logic for the Maze project.
//
// It renders a thin-wall maze into a 500x500 pixel field and provides user
// interaction for selecting start/end cells and displaying the shortest path.
package ui

import (
	"image/color"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"

	"fyne.io/fyne/v2/container"

	"maze/internal/maze"
)

// FieldSizePx is the size (both width and height) of the maze drawing area in pixels.
const FieldSizePx float32 = 500

// WallWidthPx is the thickness of maze walls in pixels.
const WallWidthPx float32 = 2

// RouteWidthPx is the thickness of the solution route line in pixels.
const RouteWidthPx float32 = 2

var (
	wallColor                  = color.Black
	backgroundColor            = color.White
	routeColor                 = color.NRGBA{R: 220, G: 20, B: 60, A: 255}
	markerStartColor           = color.NRGBA{R: 34, G: 139, B: 34, A: 255}
	markerEndColor             = color.NRGBA{R: 30, G: 144, B: 255, A: 255}
	markerOutlineColor         = color.White
	markerRadius       float32 = 6
)

// MazeView renders a maze in a fixed 500x500 field and supports cell selection.
type MazeView struct {
	widget.BaseWidget

	m        *maze.Maze
	solution []maze.Point
	start    *maze.Point
	end      *maze.Point

	// OnCellTap is invoked when the user taps a cell inside the maze bounds.
	// If nil, taps are ignored.
	OnCellTap func(p maze.Point)
}

// NewMazeView creates a new MazeView with an empty state.
//
// The view starts with no maze loaded and no start/end selection.
func NewMazeView() *MazeView {
	v := &MazeView{}
	v.ExtendBaseWidget(v)
	return v
}

// SetMaze sets the maze to render and resets any cached render state.
//
// Passing nil clears the view.
func (v *MazeView) SetMaze(m *maze.Maze) {
	v.m = m
	v.solution = nil
	v.start = nil
	v.end = nil
	v.Refresh()
}

// SetSolution sets the solution path to render.
//
// The path is expected to be a sequence of maze cells including start and end.
// Passing an empty slice clears the solution.
func (v *MazeView) SetSolution(path []maze.Point) {
	v.solution = path
	v.Refresh()
}

// ClearSolution removes the currently displayed solution path.
func (v *MazeView) ClearSolution() {
	v.solution = nil
	v.Refresh()
}

// SetStart sets the start cell marker to render.
//
// Passing nil clears the start marker.
func (v *MazeView) SetStart(p *maze.Point) {
	v.start = p
	v.Refresh()
}

// SetEnd sets the end cell marker to render.
//
// Passing nil clears the end marker.
func (v *MazeView) SetEnd(p *maze.Point) {
	v.end = p
	v.Refresh()
}

// Tapped handles pointer taps and converts the tap position into a maze cell.
//
// If a maze is loaded and the tap is inside the rendered maze area, it calls
// OnCellTap with the corresponding cell coordinate.
func (v *MazeView) Tapped(ev *fyne.PointEvent) {
	if v.m == nil {
		return
	}
	row, col, ok := v.posToCell(ev.Position)
	if !ok {
		return
	}
	if v.OnCellTap != nil {
		v.OnCellTap(maze.Point{Row: row, Col: col})
	}
}

func (v *MazeView) posToCell(pos fyne.Position) (row, col int, ok bool) {
	if v.m == nil {
		return 0, 0, false
	}
	w := float64(FieldSizePx) / float64(v.m.Cols)
	h := float64(FieldSizePx) / float64(v.m.Rows)
	x, y := float64(pos.X), float64(pos.Y)
	if x < 0 || y < 0 || x >= float64(FieldSizePx) || y >= float64(FieldSizePx) {
		return 0, 0, false
	}
	col = int(math.Floor(x / w))
	row = int(math.Floor(y / h))
	if row < 0 || row >= v.m.Rows || col < 0 || col >= v.m.Cols {
		return 0, 0, false
	}
	return row, col, true
}

// CreateRenderer constructs the widget renderer for MazeView.
func (v *MazeView) CreateRenderer() fyne.WidgetRenderer {
	r := &mazeViewRenderer{view: v}
	r.Rebuild()
	return r
}

type mazeViewRenderer struct {
	view *MazeView
	objs []fyne.CanvasObject
}

// Layout lays out the renderer objects within the given size.
func (r *mazeViewRenderer) Layout(size fyne.Size) {
	if len(r.objs) == 0 {
		return
	}
	if bg, ok := r.objs[0].(*canvas.Rectangle); ok {
		bg.Resize(size)
	}
}

// MinSize returns the minimum size of the MazeView widget.
func (r *mazeViewRenderer) MinSize() fyne.Size {
	return fyne.NewSize(FieldSizePx, FieldSizePx)
}

// Refresh rebuilds the renderer objects from the current MazeView state and
// triggers a redraw.
func (r *mazeViewRenderer) Refresh() {
	r.Rebuild()
	canvas.Refresh(r.view)
}

// Destroy releases any resources held by the renderer (no-op because no external resources are allocated).
func (r *mazeViewRenderer) Destroy() {}

// Objects returns the canvas objects that make up the current MazeView rendering.
func (r *mazeViewRenderer) Objects() []fyne.CanvasObject {
	return r.objs
}

// Rebuild regenerates the full list of canvas objects for the MazeView.
func (r *mazeViewRenderer) Rebuild() {
	v := r.view
	objs := make([]fyne.CanvasObject, 0, 1024)

	bg := canvas.NewRectangle(backgroundColor)
	bg.Resize(fyne.NewSize(FieldSizePx, FieldSizePx))
	objs = append(objs, bg)

	if v.m == nil {
		r.objs = objs
		return
	}

	cellW := FieldSizePx / float32(v.m.Cols)
	cellH := FieldSizePx / float32(v.m.Rows)

	// Outer borders: top and left (right/bottom are stored as border walls).
	objs = append(objs,
		line(0, 0, FieldSizePx, 0, wallColor, WallWidthPx),
		line(0, 0, 0, FieldSizePx, wallColor, WallWidthPx),
	)

	// Internal + right/bottom borders.
	for row := 0; row < v.m.Rows; row++ {
		y0 := float32(row) * cellH
		y1 := y0 + cellH
		for col := 0; col < v.m.Cols; col++ {
			x0 := float32(col) * cellW
			x1 := x0 + cellW

			if v.m.RightWalls[row][col] == 1 {
				objs = append(objs, line(x1, y0, x1, y1, wallColor, WallWidthPx))
			}
			if v.m.BottomWalls[row][col] == 1 {
				objs = append(objs, line(x0, y1, x1, y1, wallColor, WallWidthPx))
			}
		}
	}

	// Solution (draw segments through centers).
	if len(v.solution) >= 2 {
		for i := 0; i < len(v.solution)-1; i++ {
			a := center(v.solution[i], cellW, cellH)
			b := center(v.solution[i+1], cellW, cellH)
			objs = append(objs, line(a.X, a.Y, b.X, b.Y, routeColor, RouteWidthPx))
		}
	}

	// Markers.
	if v.start != nil {
		c := center(*v.start, cellW, cellH)
		objs = append(objs, marker(c, markerStartColor))
	}
	if v.end != nil {
		c := center(*v.end, cellW, cellH)
		objs = append(objs, marker(c, markerEndColor))
	}

	r.objs = objs
}

func line(x1, y1, x2, y2 float32, c color.Color, w float32) *canvas.Line {
	l := canvas.NewLine(c)
	l.StrokeWidth = w
	l.Position1 = fyne.NewPos(x1, y1)
	l.Position2 = fyne.NewPos(x2, y2)
	return l
}

func center(p maze.Point, cellW, cellH float32) fyne.Position {
	return fyne.NewPos((float32(p.Col)+0.5)*cellW, (float32(p.Row)+0.5)*cellH)
}

func marker(pos fyne.Position, c color.Color) fyne.CanvasObject {
	circ := canvas.NewCircle(c)
	circ.Resize(fyne.NewSize(markerRadius*2, markerRadius*2))
	circ.Move(fyne.NewPos(pos.X-markerRadius, pos.Y-markerRadius))

	outline := canvas.NewCircle(color.NRGBA{A: 0})
	outline.StrokeWidth = 2
	outline.StrokeColor = markerOutlineColor
	outline.Resize(circ.Size())
	outline.Move(circ.Position())

	return container.NewWithoutLayout(outline, circ)
}
