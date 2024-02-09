package display

import (
	"main/pkg/types"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type LastRoundTableData struct {
	tview.TableContentReadOnly

	Validators   types.ValidatorsWithInfo
	ColumnsCount int
}

func NewLastRoundTableData(columnsCount int) *LastRoundTableData {
	return &LastRoundTableData{
		ColumnsCount: columnsCount,
		Validators:   make(types.ValidatorsWithInfo, 0),
	}
}

func (d *LastRoundTableData) SetColumnsCount(count int) {
	d.ColumnsCount = count
}

func (d *LastRoundTableData) GetCell(row, column int) *tview.TableCell {
	index := row*d.ColumnsCount + column
	text := ""

	if index < len(d.Validators) {
		text = d.Validators[index].Serialize()
	}

	cell := tview.NewTableCell(text)

	if index < len(d.Validators) && d.Validators[index].RoundVote.IsProposer {
		cell.SetBackgroundColor(tcell.ColorForestGreen)
	}

	return cell
}

func (d *LastRoundTableData) GetRowCount() int {
	if len(d.Validators)%d.ColumnsCount == 0 {
		return len(d.Validators) / d.ColumnsCount
	}

	return len(d.Validators)/d.ColumnsCount + 1
}

func (d *LastRoundTableData) GetColumnCount() int {
	return d.ColumnsCount
}

func (d *LastRoundTableData) SetValidators(validators types.ValidatorsWithInfo) {
	d.Validators = validators
}
