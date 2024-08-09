package display

import (
	"main/pkg/types"
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
}

func NewAllRoundsTableData(disableEmojis bool, transpose bool) *AllRoundsTableData {
	return &AllRoundsTableData{
		Validators:    types.ValidatorsWithInfoAndAllRoundVotes{},
		DisableEmojis: disableEmojis,
		Transpose:     transpose,
		cells:         [][]*tview.TableCell{},
	}
}

func (d *AllRoundsTableData) GetCell(row, column int) *tview.TableCell {
	if len(d.cells) <= row {
		return nil
	}

	if len(d.cells[row]) <= column {
		return nil
	}

	return d.cells[row][column]
}

func (d *AllRoundsTableData) GetRowCount() int {
	return len(d.cells)
}

func (d *AllRoundsTableData) GetColumnCount() int {
	if len(d.cells) == 0 {
		return 0
	}

	return len(d.cells[0])
}

func (d *AllRoundsTableData) SetValidators(validators types.ValidatorsWithInfoAndAllRoundVotes) {
	if d.Validators.Equals(validators) {
		return
	}

	d.Validators = validators
	d.redrawCells()
}

func (d *AllRoundsTableData) SetTranspose(transpose bool) {
	d.Transpose = transpose
	d.redrawCells()
}

func (d *AllRoundsTableData) redrawCells() {
	d.cells = make([][]*tview.TableCell, len(d.Validators.Validators)+1)

	for row := 0; row < len(d.Validators.Validators)+1; row++ {
		d.cells[row] = make([]*tview.TableCell, len(d.Validators.RoundsVotes)+1)

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

				d.cells[row][column] = tview.
					NewTableCell(text).
					SetAlign(tview.AlignCenter).
					SetStyle(tcell.StyleDefault.Bold(true))
				continue
			}

			// First column is always validators list.
			if column == 0 {
				text := d.Validators.Validators[row-1].Serialize()
				cell := tview.NewTableCell(text)
				d.cells[row][column] = cell
				continue
			}

			roundVotes := d.Validators.RoundsVotes[round]
			roundVote := roundVotes[row-1]
			text := roundVote.Serialize(d.DisableEmojis)

			cell := tview.NewTableCell(text)

			if roundVote.IsProposer {
				cell.SetBackgroundColor(tcell.ColorForestGreen)
			}

			d.cells[row][column] = cell
		}
	}
}
