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
}

func NewAllRoundsTableData(disableEmojis bool, transpose bool) *AllRoundsTableData {
	return &AllRoundsTableData{
		Validators:    types.ValidatorsWithInfoAndAllRoundVotes{},
		DisableEmojis: disableEmojis,
		Transpose:     transpose,
	}
}

func (d *AllRoundsTableData) GetCell(row, column int) *tview.TableCell {
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

		return tview.
			NewTableCell(text).
			SetAlign(tview.AlignCenter).
			SetStyle(tcell.StyleDefault.Bold(true))
	}

	// First column is always validators list.
	if column == 0 {
		text := d.Validators.Validators[row-1].Serialize()
		cell := tview.NewTableCell(text)
		return cell
	}

	roundVotes := d.Validators.RoundsVotes[round]
	roundVote := roundVotes[row-1]
	text := roundVote.Serialize(d.DisableEmojis)

	cell := tview.NewTableCell(text)

	if roundVote.IsProposer {
		cell.SetBackgroundColor(tcell.ColorForestGreen)
	}

	return cell
}

func (d *AllRoundsTableData) GetRowCount() int {
	return len(d.Validators.Validators) + 1
}

func (d *AllRoundsTableData) GetColumnCount() int {
	return len(d.Validators.RoundsVotes) + 1
}

func (d *AllRoundsTableData) SetValidators(validators types.ValidatorsWithInfoAndAllRoundVotes) {
	d.Validators = validators
}

func (d *AllRoundsTableData) SetTranspose(transpose bool) {
	d.Transpose = transpose
}
