package display

import (
	"main/pkg/types"
	"main/pkg/utils"
	"strconv"

	"github.com/gdamore/tcell/v2"

	"github.com/rivo/tview"
)

type AllRoundsTableData struct {
	tview.TableContentReadOnly

	Validators    types.ValidatorsWithInfoAndAllRoundVotes
	DisableEmojis bool
	Transpose     bool

	cells [][]*tview.TableCell
	mutex *utils.NoopLocker
}

func NewAllRoundsTableData(disableEmojis bool, transpose bool) *AllRoundsTableData {
	return &AllRoundsTableData{
		Validators:    types.ValidatorsWithInfoAndAllRoundVotes{},
		DisableEmojis: disableEmojis,
		Transpose:     transpose,
		cells:         [][]*tview.TableCell{},
		mutex:         &utils.NoopLocker{},
	}
}

func (d *AllRoundsTableData) GetCell(row, column int) *tview.TableCell {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	if len(d.cells) <= row {
		return nil
	}

	if len(d.cells[row]) <= column {
		return nil
	}

	return d.cells[row][column]
}

func (d *AllRoundsTableData) GetRowCount() int {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	return len(d.cells)
}

func (d *AllRoundsTableData) GetColumnCount() int {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	if len(d.cells) == 0 {
		return 0
	}

	return len(d.cells[0])
}

func (d *AllRoundsTableData) SetValidators(validators types.ValidatorsWithInfoAndAllRoundVotes) {
	d.mutex.RLock()
	if d.Validators.Equals(validators) {
		return
	}
	d.mutex.RUnlock()

	d.mutex.Lock()
	d.Validators = validators
	d.mutex.Unlock()

	d.redrawCells()
}

func (d *AllRoundsTableData) SetTranspose(transpose bool) {
	d.mutex.Lock()
	d.Transpose = transpose
	d.mutex.Unlock()
	d.redrawCells()
}

func (d *AllRoundsTableData) redrawCells() {
	cells := d.makeCells()

	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.cells = cells
}

func (d *AllRoundsTableData) makeCells() [][]*tview.TableCell {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	cells := make([][]*tview.TableCell, len(d.Validators.Validators)+1)

	for row := 0; row < len(d.Validators.Validators)+1; row++ {
		cells[row] = make([]*tview.TableCell, len(d.Validators.RoundsVotes)+1)

		for column := 0; column < len(d.Validators.RoundsVotes)+1; column++ {
			round := column - 1
			if d.Transpose {
				round = len(d.Validators.RoundsVotes) - column
			}

			// Table header.
			if row == 0 {
				text := "validator"
				if column != 0 {
					text = strconv.Itoa(round)
				}

				cells[row][column] = tview.
					NewTableCell(text).
					SetAlign(tview.AlignCenter).
					SetStyle(tcell.StyleDefault.Bold(true))
				continue
			}

			// First column is always validators list.
			if column == 0 {
				text := d.Validators.Validators[row-1].Serialize()
				cell := tview.NewTableCell(text)
				cells[row][column] = cell
				continue
			}

			roundVotes := d.Validators.RoundsVotes[round]
			roundVote := roundVotes[row-1]
			text := roundVote.Serialize(d.DisableEmojis)

			cell := tview.NewTableCell(text)

			if roundVote.IsProposer {
				cell.SetBackgroundColor(tcell.ColorForestGreen)
			}

			cells[row][column] = cell
		}
	}
	return cells
}
