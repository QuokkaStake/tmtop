package display

import (
	"main/pkg/types"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type LastRoundTableData struct {
	tview.TableContentReadOnly

	Validators    types.ValidatorsWithInfo
	ColumnsCount  int
	DisableEmojis bool
	Transpose     bool
}

func NewLastRoundTableData(columnsCount int, disableEmojis bool, transpose bool) *LastRoundTableData {
	return &LastRoundTableData{
		ColumnsCount:  columnsCount,
		Validators:    make(types.ValidatorsWithInfo, 0),
		DisableEmojis: disableEmojis,
		Transpose:     transpose,
	}
}

func (d *LastRoundTableData) SetColumnsCount(count int) {
	d.ColumnsCount = count
}

func (d *LastRoundTableData) SetTranspose(transpose bool) {
	d.Transpose = transpose
}

func (d *LastRoundTableData) GetCell(row, column int) *tview.TableCell {
	index := row*d.ColumnsCount + column

	if d.Transpose {
		rows := d.GetRowCount()
		index = column*rows + row
	}

	text := ""

	if index < len(d.Validators) {
		text = d.Validators[index].Serialize(d.DisableEmojis)
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
