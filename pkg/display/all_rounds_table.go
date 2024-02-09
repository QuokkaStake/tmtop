package display

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"main/pkg/types"

	"github.com/rivo/tview"
)

type AllRoundsTableData struct {
	tview.TableContentReadOnly

	Validators types.ValidatorsWithInfoAndAllRoundVotes
}

func NewAllRoundsTableData() *AllRoundsTableData {
	return &AllRoundsTableData{
		Validators: types.ValidatorsWithInfoAndAllRoundVotes{},
	}
}

func (d *AllRoundsTableData) GetCell(row, column int) *tview.TableCell {
	// Table header.
	if row == 0 {
		text := "validator"
		if column != 0 {
			text = fmt.Sprintf("%d", column-1)
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

	roundVotes := d.Validators.RoundsVotes[column-1]
	roundVote := roundVotes[row-1]
	text := roundVote.Serialize()

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
