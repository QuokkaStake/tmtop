package display

import (
	"fmt"
	"main/pkg/types"
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type LastRoundTableData struct {
	tview.TableContentReadOnly

	Validators              types.ValidatorsWithInfo
	CurrentValidatorAddress string
	ConsensusError          error
	ColumnsCount            int
	DisableEmojis           bool
	Transpose               bool

	cells [][]*tview.TableCell
	mutex sync.Mutex
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
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if len(d.cells) <= row {
		return nil
	}

	if len(d.cells[row]) <= column {
		return nil
	}

	return d.cells[row][column]
}

func (d *LastRoundTableData) GetRowCount() int {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	return len(d.cells)
}

func (d *LastRoundTableData) GetColumnCount() int {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if len(d.cells) == 0 {
		return 0
	}

	return len(d.cells[0])
}

func (d *LastRoundTableData) SetValidators(
	validators types.ValidatorsWithInfo,
	consensusError error,
	validatorInfo types.TendermintValidatorInfo,
) {
	d.Validators = validators
	d.ConsensusError = consensusError
	d.CurrentValidatorAddress = validatorInfo.Address
	d.redrawData()
}

func (d *LastRoundTableData) redrawData() {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if d.ConsensusError != nil {
		d.cells = [][]*tview.TableCell{
			{
				tview.NewTableCell(fmt.Sprintf(" Error fetching consensus: %s", d.ConsensusError)),
			},
		}
		return
	}

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
				rows := len(d.cells)
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

			if index < len(d.Validators) && d.Validators[index].Validator.Address == d.CurrentValidatorAddress {
				cell.SetBackgroundColor(tcell.ColorMediumTurquoise)
			}

			d.cells[row][column] = cell
		}
	}
}
