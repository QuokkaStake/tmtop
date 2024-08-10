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

	cells [][]*tview.TableCell
}

func NewLastRoundTableData(columnsCount int, disableEmojis bool, transpose bool) *LastRoundTableData {
	return &LastRoundTableData{
		ColumnsCount:  columnsCount,
		Validators:    make(types.ValidatorsWithInfo, 0),
		DisableEmojis: disableEmojis,
		Transpose:     transpose,

		cells: [][]*tview.TableCell{},
	}
}

func (d *LastRoundTableData) SetColumnsCount(count int) {
	d.ColumnsCount = count
	d.redrawData()
}

func (d *LastRoundTableData) SetTranspose(transpose bool) {
	d.Transpose = transpose
	d.redrawData()
}

func (d *LastRoundTableData) GetCell(row, column int) *tview.TableCell {
	if len(d.cells) <= row {
		return nil
	}

	if len(d.cells[row]) <= column {
		return nil
	}

	return d.cells[row][column]
}

func (d *LastRoundTableData) GetRowCount() int {
	return len(d.cells)
}

func (d *LastRoundTableData) GetColumnCount() int {
	if len(d.cells) == 0 {
		return 0
	}

	return len(d.cells[0])
}

func (d *LastRoundTableData) SetValidators(validators types.ValidatorsWithInfo) {
	d.Validators = validators
	d.redrawData()
}

func (d *LastRoundTableData) redrawData() {
	rowsCount := len(d.Validators)/d.ColumnsCount + 1
	if len(d.Validators)%d.ColumnsCount == 0 {
		rowsCount = len(d.Validators) / d.ColumnsCount
	}

	d.cells = make([][]*tview.TableCell, rowsCount)

	for row := 0; row < rowsCount; row++ {
		d.cells[row] = make([]*tview.TableCell, d.ColumnsCount)

		for column := 0; column < d.ColumnsCount; column++ {
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

			d.cells[row][column] = cell
		}
	}
}
