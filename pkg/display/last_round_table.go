package display

import (
	"fmt"
	"main/pkg/types"
	"main/pkg/utils"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type LastRoundTableData struct {
	tview.TableContentReadOnly

	Validators     types.ValidatorsWithInfo
	ConsensusError error
	ColumnsCount   int
	DisableEmojis  bool
	Transpose      bool

	cells [][]*tview.TableCell
	mutex *utils.NoopLocker
}

func NewLastRoundTableData(columnsCount int, disableEmojis bool, transpose bool) *LastRoundTableData {
	return &LastRoundTableData{
		ColumnsCount:  columnsCount,
		Validators:    make(types.ValidatorsWithInfo, 0),
		DisableEmojis: disableEmojis,
		Transpose:     transpose,

		cells: [][]*tview.TableCell{},
		mutex: &utils.NoopLocker{},
	}
}

func (d *LastRoundTableData) SetColumnsCount(count int) {
	d.mutex.Lock()
	d.ColumnsCount = count
	d.mutex.Unlock()

	d.redrawData()
}

func (d *LastRoundTableData) SetTranspose(transpose bool) {
	d.mutex.Lock()
	d.Transpose = transpose
	d.mutex.Unlock()

	d.redrawData()
}

func (d *LastRoundTableData) GetCell(row, column int) *tview.TableCell {
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

func (d *LastRoundTableData) GetRowCount() int {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	return len(d.cells)
}

func (d *LastRoundTableData) GetColumnCount() int {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	if len(d.cells) == 0 {
		return 0
	}

	return len(d.cells[0])
}

func (d *LastRoundTableData) SetValidators(validators types.ValidatorsWithInfo, consensusError error) {
	d.mutex.Lock()
	d.Validators = validators
	d.ConsensusError = consensusError
	d.mutex.Unlock()

	d.redrawData()
}

func (d *LastRoundTableData) redrawData() {
	cells := d.makeCells()

	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.cells = cells
}

func (d *LastRoundTableData) makeCells() [][]*tview.TableCell {
	if d.ConsensusError != nil {
		return [][]*tview.TableCell{
			{tview.NewTableCell(fmt.Sprintf(" Error fetching consensus: %s", d.ConsensusError))},
		}
	}

	rowsCount := len(d.Validators)/d.ColumnsCount + 1
	if len(d.Validators)%d.ColumnsCount == 0 {
		rowsCount = len(d.Validators) / d.ColumnsCount
	}

	cells := make([][]*tview.TableCell, rowsCount)

	for row := 0; row < rowsCount; row++ {
		cells[row] = make([]*tview.TableCell, d.ColumnsCount)

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

			cells[row][column] = cell
		}
	}
	return cells
}
