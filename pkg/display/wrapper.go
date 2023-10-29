package display

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog"
	"main/pkg/types"
)

const (
	ColumnsAmount    = 3
	RowsAmount       = 10
	DebugBlockHeight = 2
)

type Wrapper struct {
	InfoTextView     *tview.TextView
	ProgressTextView *tview.TextView
	DebugTextView    *tview.TextView
	Table            *tview.Table
	TableData        *TableData
	Grid             *tview.Grid
	App              *tview.Application
	InfoBlockWidth   int

	DebugEnabled bool

	Logger zerolog.Logger
}

func NewWrapper(logger zerolog.Logger) *Wrapper {
	tableData := NewTableData(ColumnsAmount)

	table := tview.NewTable().
		SetBorders(false).
		SetSelectable(false, false).
		// SetSeparator(tview.Borders.Vertical).
		SetContent(tableData)

	infoTextView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true)

	progressTextView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true)

	debugTextView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true)

	grid := tview.NewGrid().
		SetRows(0, 0, 0, 0, 0, 0, 0, 0, 0, 0).
		SetColumns(0, 0, 0, 0, 0, 0).
		SetBorders(true)

	app := tview.NewApplication().SetRoot(grid, true).SetFocus(table)

	return &Wrapper{
		InfoTextView:     infoTextView,
		ProgressTextView: progressTextView,
		DebugTextView:    debugTextView,
		Table:            table,
		TableData:        tableData,
		Grid:             grid,
		App:              app,
		Logger:           logger.With().Str("component", "display_wrapper").Logger(),
		DebugEnabled:     false,
		InfoBlockWidth:   2,
	}
}

func (w *Wrapper) Start() {
	w.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'q' {
			w.App.Stop()
		}

		if event.Rune() == 'd' {
			w.ToggleDebug()
		}

		if event.Rune() == 'b' {
			w.ChangeInfoBlockHeight(true)
		}

		if event.Rune() == 's' {
			w.ChangeInfoBlockHeight(false)
		}

		return event
	})

	w.Grid.SetBackgroundColor(tcell.ColorDefault)
	w.Table.SetBackgroundColor(tcell.ColorDefault)
	w.InfoTextView.SetBackgroundColor(tcell.ColorDefault)
	w.ProgressTextView.SetBackgroundColor(tcell.ColorDefault)
	w.DebugTextView.SetBackgroundColor(tcell.ColorDefault)

	w.Redraw()

	fmt.Fprint(w.InfoTextView, "testtesttest")
	fmt.Fprint(w.ProgressTextView, "testtesttest")
	fmt.Fprint(w.DebugTextView, "testtesttest")

	w.App.SetBeforeDrawFunc(func(screen tcell.Screen) bool {
		screen.Clear()
		return false
	})

	if err := w.App.Run(); err != nil {
		w.Logger.Fatal().Err(err).Msg("Could not draw screen")
	}
}

func (w *Wrapper) ToggleDebug() {
	w.DebugEnabled = !w.DebugEnabled

	w.Redraw()
}

func (w *Wrapper) SetState(state *types.State) {
	w.TableData.SetValidators(state.GetValidatorsWithInfo())

	w.InfoTextView.Clear()
	fmt.Fprint(w.InfoTextView, state.SerializeInfo())
	w.App.Draw()
}

func (w *Wrapper) DebugText(text string) {
	fmt.Fprint(w.DebugTextView, text+"\n")
}

func (w *Wrapper) ChangeInfoBlockHeight(increase bool) {
	//w.InfoBlockWidth++

	if increase && w.InfoBlockWidth+1 <= RowsAmount-DebugBlockHeight-1 {
		w.InfoBlockWidth++
	} else if !increase && w.InfoBlockWidth-1 >= 1 {
		w.InfoBlockWidth--
	}

	w.Redraw()
}

func (w *Wrapper) Redraw() {
	w.Grid.RemoveItem(w.InfoTextView)
	w.Grid.RemoveItem(w.ProgressTextView)
	w.Grid.RemoveItem(w.Table)
	w.Grid.RemoveItem(w.DebugTextView)

	w.Grid.AddItem(w.InfoTextView, 0, 0, w.InfoBlockWidth, 3, 1, 1, false)
	w.Grid.AddItem(w.ProgressTextView, 0, 3, w.InfoBlockWidth, 3, 1, 1, false)

	if w.DebugEnabled {
		w.Grid.AddItem(
			w.Table,
			w.InfoBlockWidth,
			0,
			RowsAmount-w.InfoBlockWidth-DebugBlockHeight,
			6,
			0,
			0,
			false,
		)
		w.Grid.AddItem(
			w.DebugTextView,
			RowsAmount-DebugBlockHeight,
			0,
			DebugBlockHeight,
			6,
			0,
			0,
			false,
		)
	} else {
		w.Grid.AddItem(
			w.Table,
			w.InfoBlockWidth,
			0,
			RowsAmount-w.InfoBlockWidth,
			6,
			0,
			0,
			false,
		)
	}
}
