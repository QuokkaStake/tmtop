package display

import (
	"main/pkg/types"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type TableData struct {
	tview.TableContentReadOnly

	Validators   types.ValidatorsWithInfo
	ColumnsCount int
}

func NewTableData(columnsCount int) *TableData {
	return &TableData{
		ColumnsCount: columnsCount,
		Validators:   make(types.ValidatorsWithInfo, 0),
	}
}

func (d *TableData) GetCell(row, column int) *tview.TableCell {
	index := row*d.ColumnsCount + column
	text := ""

	if index < len(d.Validators) {
		text = d.Validators[index].Serialize()
	}

	cell := tview.NewTableCell(text)

	if index < len(d.Validators) && d.Validators[index].Validator.IsProposer {
		cell.SetBackgroundColor(tcell.ColorPeachPuff)
	}

	return cell
}

func (d *TableData) GetRowCount() int {
	return len(d.Validators) / d.ColumnsCount
}

func (d *TableData) GetColumnCount() int {
	return d.ColumnsCount
}

func (d *TableData) SetValidators(validators types.ValidatorsWithInfo) {
	d.Validators = validators
}