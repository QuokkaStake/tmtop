package display

import (
	"fmt"
	"main/pkg/types"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog"
)

const (
	ColumnsAmount    = 3
	RowsAmount       = 10
	DebugBlockHeight = 2
)

type Wrapper struct {
	ConsensusInfoTextView *tview.TextView
	ChainInfoTextView     *tview.TextView
	ProgressTextView      *tview.TextView
	DebugTextView         *tview.TextView
	Table                 *tview.Table
	TableData             *TableData
	Grid                  *tview.Grid
	App                   *tview.Application
	InfoBlockWidth        int

	DebugEnabled bool

	Logger zerolog.Logger

	PauseChannel chan bool
	IsPaused     bool
}

func NewWrapper(logger zerolog.Logger, pauseChannel chan bool) *Wrapper {
	tableData := NewTableData(ColumnsAmount)

	table := tview.NewTable().
		SetBorders(false).
		SetSelectable(false, false).
		// SetSeparator(tview.Borders.Vertical).
		SetContent(tableData)

	consensusInfoTextView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true)

	chainInfoTextView := tview.NewTextView().
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
		ChainInfoTextView:     chainInfoTextView,
		ConsensusInfoTextView: consensusInfoTextView,
		ProgressTextView:      progressTextView,
		DebugTextView:         debugTextView,
		Table:                 table,
		TableData:             tableData,
		Grid:                  grid,
		App:                   app,
		Logger:                logger.With().Str("component", "display_wrapper").Logger(),
		DebugEnabled:          false,
		InfoBlockWidth:        2,
		PauseChannel:          pauseChannel,
		IsPaused:              false,
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

		if event.Rune() == 'p' {
			w.IsPaused = !w.IsPaused
			w.PauseChannel <- w.IsPaused
		}

		return event
	})

	w.Grid.SetBackgroundColor(tcell.ColorDefault)
	w.Table.SetBackgroundColor(tcell.ColorDefault)
	w.ChainInfoTextView.SetBackgroundColor(tcell.ColorDefault)
	w.ConsensusInfoTextView.SetBackgroundColor(tcell.ColorDefault)
	w.ProgressTextView.SetBackgroundColor(tcell.ColorDefault)
	w.DebugTextView.SetBackgroundColor(tcell.ColorDefault)

	w.Redraw()

	fmt.Fprint(w.ChainInfoTextView, "Loading...")
	fmt.Fprint(w.ConsensusInfoTextView, "Loading...")
	fmt.Fprint(w.ProgressTextView, "Loading...")

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

	w.ConsensusInfoTextView.Clear()
	w.ChainInfoTextView.Clear()
	fmt.Fprint(w.ConsensusInfoTextView, state.SerializeConsensus())
	fmt.Fprint(w.ChainInfoTextView, state.SerializeChainInfo())
	w.App.Draw()
}

func (w *Wrapper) DebugText(text string) {
	fmt.Fprint(w.DebugTextView, text)
	w.DebugTextView.ScrollToEnd()
}

func (w *Wrapper) ChangeInfoBlockHeight(increase bool) {
	if increase && w.InfoBlockWidth+1 <= RowsAmount-DebugBlockHeight-1 {
		w.InfoBlockWidth++
	} else if !increase && w.InfoBlockWidth-1 >= 1 {
		w.InfoBlockWidth--
	}

	w.Redraw()
}

func (w *Wrapper) Redraw() {
	w.Grid.RemoveItem(w.ConsensusInfoTextView)
	w.Grid.RemoveItem(w.ChainInfoTextView)
	w.Grid.RemoveItem(w.ProgressTextView)
	w.Grid.RemoveItem(w.Table)
	w.Grid.RemoveItem(w.DebugTextView)

	w.Grid.AddItem(w.ConsensusInfoTextView, 0, 0, w.InfoBlockWidth, 2, 1, 1, false)
	w.Grid.AddItem(w.ChainInfoTextView, 0, 2, w.InfoBlockWidth, 2, 1, 1, false)
	w.Grid.AddItem(w.ProgressTextView, 0, 4, w.InfoBlockWidth, 2, 1, 1, false)

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
