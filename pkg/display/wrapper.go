package display

import (
	"fmt"
	"main/pkg/types"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog"
)

const ColumnsAmount = 3

type Wrapper struct {
	InfoTextView     *tview.TextView
	ProgressTextView *tview.TextView
	Table            *tview.Table
	TableData        *TableData
	Grid             *tview.Grid
	App              *tview.Application

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

	grid := tview.NewGrid().
		SetRows(0, 0, 0, 0, 0, 0, 0, 0, 0, 0).
		SetColumns(0, 0, 0, 0, 0, 0).
		SetBorders(true)

	app := tview.NewApplication().SetRoot(grid, true).SetFocus(table)

	return &Wrapper{
		InfoTextView:     infoTextView,
		ProgressTextView: progressTextView,
		Table:            table,
		TableData:        tableData,
		Grid:             grid,
		App:              app,
		Logger:           logger.With().Str("component", "display_wrapper").Logger(),
	}
}

func (w *Wrapper) Start() {
	w.Grid.SetBackgroundColor(tcell.ColorDefault)

	w.InfoTextView.SetBackgroundColor(tcell.ColorDefault)
	w.ProgressTextView.SetBackgroundColor(tcell.ColorDefault)

	w.Grid.AddItem(w.InfoTextView, 0, 0, 2, 3, 1, 1, false)
	w.Grid.AddItem(w.ProgressTextView, 0, 3, 2, 3, 1, 1, false)

	w.Table.SetBackgroundColor(tcell.ColorDefault)
	w.Grid.AddItem(w.Table, 2, 0, 8, 6, 0, 0, false)

	fmt.Fprint(w.InfoTextView, "testtesttest")
	fmt.Fprint(w.ProgressTextView, "testtesttest")

	w.App.SetBeforeDrawFunc(func(screen tcell.Screen) bool {
		screen.Clear()
		return false
	})

	if err := w.App.Run(); err != nil {
		w.Logger.Fatal().Err(err).Msg("Could not draw screen")
	}
}

func (w *Wrapper) SetState(state *types.State) {
	w.TableData.SetValidators(state.GetValidatorsWithInfo())

	w.InfoTextView.Clear()
	fmt.Fprint(w.InfoTextView, state.SerializeInfo())
	w.App.Draw()
}
