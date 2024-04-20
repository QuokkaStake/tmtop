package display

import (
	"fmt"
	configPkg "main/pkg/config"
	"main/pkg/types"
	"main/static"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog"
)

const (
	ModeLastRound = iota
	ModeAllRounds = iota
)

const (
	DefaultColumnsCount = 3
	RowsAmount          = 10
	DebugBlockHeight    = 2
	DefaultMode         = ModeLastRound
)

type Wrapper struct {
	ConsensusInfoTextView *tview.TextView
	ChainInfoTextView     *tview.TextView
	ProgressTextView      *tview.TextView
	DebugTextView         *tview.TextView
	LastRoundTable        *tview.Table
	LastRoundTableData    *LastRoundTableData
	AllRoundsTable        *tview.Table
	AllRoundsTableData    *AllRoundsTableData
	Grid                  *tview.Grid
	Pages                 *tview.Pages
	App                   *tview.Application
	HelpModal             *tview.Modal

	InfoBlockWidth int
	ColumnsCount   int
	Mode           int

	DebugEnabled bool

	Logger zerolog.Logger

	PauseChannel chan bool
	IsPaused     bool

	IsHelpDisplayed bool

	DisableEmojis bool
	Transpose     bool
}

func NewWrapper(
	config configPkg.Config,
	logger zerolog.Logger,
	pauseChannel chan bool,
	appVersion string,
) *Wrapper {
	lastRoundTableData := NewLastRoundTableData(DefaultColumnsCount, config.DisableEmojis, false)
	allRoundsTableData := NewAllRoundsTableData(config.DisableEmojis)

	helpTextBytes, _ := static.TemplatesFs.ReadFile("help.txt")
	helpText := strings.ReplaceAll(string(helpTextBytes), "{{ Version }}", appVersion)

	lastRoundTable := tview.NewTable().
		SetBorders(false).
		SetSelectable(false, false).
		SetContent(lastRoundTableData)

	allRoundsTable := tview.NewTable().
		SetBorders(false).
		SetSelectable(false, false).
		SetContent(allRoundsTableData).
		SetFixed(len(allRoundsTableData.Headers), 1)

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

	helpModal := tview.NewModal().SetText(helpText)

	grid := tview.NewGrid().
		SetRows(0, 0, 0, 0, 0, 0, 0, 0, 0, 0).
		SetColumns(0, 0, 0, 0, 0, 0).
		SetBorders(true)

	pages := tview.NewPages().AddPage("grid", grid, true, true)

	app := tview.NewApplication().SetRoot(pages, true).SetFocus(lastRoundTable)

	return &Wrapper{
		ChainInfoTextView:     chainInfoTextView,
		ConsensusInfoTextView: consensusInfoTextView,
		ProgressTextView:      progressTextView,
		DebugTextView:         debugTextView,
		LastRoundTable:        lastRoundTable,
		LastRoundTableData:    lastRoundTableData,
		AllRoundsTable:        allRoundsTable,
		AllRoundsTableData:    allRoundsTableData,
		HelpModal:             helpModal,
		Grid:                  grid,
		Pages:                 pages,
		App:                   app,
		Logger:                logger.With().Str("component", "display_wrapper").Logger(),
		DebugEnabled:          false,
		InfoBlockWidth:        2,
		ColumnsCount:          DefaultColumnsCount,
		Mode:                  DefaultMode,
		PauseChannel:          pauseChannel,
		IsPaused:              false,
		IsHelpDisplayed:       false,
		DisableEmojis:         config.DisableEmojis,
		Transpose:             false,
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

		if event.Rune() == 'h' {
			w.ToggleHelp()
		}

		if event.Rune() == 'm' {
			w.ChangeColumnsCount(true)
		}

		if event.Rune() == 'l' {
			w.ChangeColumnsCount(false)
		}

		if event.Rune() == 't' {
			w.Transpose = !w.Transpose
			w.LastRoundTableData.SetTranspose(w.Transpose)
		}

		if event.Rune() == 'p' {
			w.IsPaused = !w.IsPaused
			w.PauseChannel <- w.IsPaused
		}

		if event.Key() == tcell.KeyTAB {
			w.ChangeMode()
		}

		return event
	})

	w.Grid.SetBackgroundColor(tcell.ColorDefault)
	w.LastRoundTable.SetBackgroundColor(tcell.ColorDefault)
	w.AllRoundsTable.SetBackgroundColor(tcell.ColorDefault)
	w.ChainInfoTextView.SetBackgroundColor(tcell.ColorDefault)
	w.ConsensusInfoTextView.SetBackgroundColor(tcell.ColorDefault)
	w.ProgressTextView.SetBackgroundColor(tcell.ColorDefault)
	w.DebugTextView.SetBackgroundColor(tcell.ColorDefault)

	w.Redraw()

	_, _ = fmt.Fprint(w.ChainInfoTextView, "Loading...")
	_, _ = fmt.Fprint(w.ConsensusInfoTextView, "Loading...")
	_, _ = fmt.Fprint(w.ProgressTextView, "Loading...")

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

func (w *Wrapper) ToggleHelp() {
	w.IsHelpDisplayed = !w.IsHelpDisplayed

	w.Redraw()
}

func (w *Wrapper) SetState(state *types.State) {
	w.LastRoundTableData.SetValidators(state.GetValidatorsWithInfo())
	w.AllRoundsTableData.SetValidators(state.GetValidatorsWithInfoAndAllRoundVotes())

	w.ConsensusInfoTextView.Clear()
	w.ChainInfoTextView.Clear()
	w.ProgressTextView.Clear()
	_, _ = fmt.Fprint(w.ConsensusInfoTextView, state.SerializeConsensus())
	_, _ = fmt.Fprint(w.ChainInfoTextView, state.SerializeChainInfo())

	_, _, width, height := w.ConsensusInfoTextView.GetInnerRect()
	_, _ = fmt.Fprint(w.ProgressTextView, state.SerializePrevotesProgressbar(width, height/2))
	_, _ = fmt.Fprint(w.ProgressTextView, "\n")
	_, _ = fmt.Fprint(w.ProgressTextView, state.SerializePrecommitsProgressbar(width, height/2))

	w.App.Draw()
}

func (w *Wrapper) DebugText(text string) {
	_, _ = fmt.Fprint(w.DebugTextView, text)
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

func (w *Wrapper) ChangeColumnsCount(increase bool) {
	if increase {
		w.ColumnsCount++
	} else if !increase && w.ColumnsCount-1 >= 1 {
		w.ColumnsCount--
	}

	w.LastRoundTableData.SetColumnsCount(w.ColumnsCount)

	w.Redraw()
}

func (w *Wrapper) ChangeMode() {
	switch w.Mode {
	case ModeAllRounds:
		w.Mode = ModeLastRound
	case ModeLastRound:
		w.Mode = ModeAllRounds
	default:
		w.Mode = ModeLastRound
	}

	w.Redraw()
}

func (w *Wrapper) Redraw() {
	table := w.LastRoundTable
	if w.Mode == ModeAllRounds {
		table = w.AllRoundsTable
	}

	w.Grid.RemoveItem(w.ConsensusInfoTextView)
	w.Grid.RemoveItem(w.ChainInfoTextView)
	w.Grid.RemoveItem(w.ProgressTextView)
	w.Grid.RemoveItem(w.LastRoundTable)
	w.Grid.RemoveItem(w.AllRoundsTable)
	w.Grid.RemoveItem(w.DebugTextView)

	w.Grid.AddItem(w.ConsensusInfoTextView, 0, 0, w.InfoBlockWidth, 2, 1, 1, false)
	w.Grid.AddItem(w.ChainInfoTextView, 0, 2, w.InfoBlockWidth, 2, 1, 1, false)
	w.Grid.AddItem(w.ProgressTextView, 0, 4, w.InfoBlockWidth, 2, 1, 1, false)

	if w.DebugEnabled {
		w.Grid.AddItem(
			table,
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
			table,
			w.InfoBlockWidth,
			0,
			RowsAmount-w.InfoBlockWidth,
			6,
			0,
			0,
			false,
		)
	}

	if w.IsHelpDisplayed {
		w.Pages.AddPage("modal", w.HelpModal, true, true)
	} else {
		w.Pages.RemovePage("modal")
	}

	w.App.SetFocus(table)
}
