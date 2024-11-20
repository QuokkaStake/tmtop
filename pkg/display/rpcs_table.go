package display

import (
	"main/pkg/types"
	"main/pkg/utils"

	"github.com/rivo/tview"
)

type RPCsTableData struct {
	tview.TableContentReadOnly

	knownRPCs []types.RPC

	cells [][]*tview.TableCell
	mutex *utils.NoopLocker
}

func NewRPCsTableData() *RPCsTableData {
	return &RPCsTableData{
		cells: [][]*tview.TableCell{},
		mutex: &utils.NoopLocker{},
	}
}

func (d *RPCsTableData) SetKnownRPCs(rpcs []types.RPC) {
	d.mutex.Lock()
	d.knownRPCs = rpcs
	d.mutex.Unlock()

	d.redrawData()
}

func (d *RPCsTableData) GetCell(row, column int) *tview.TableCell {
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

func (d *RPCsTableData) GetRowCount() int {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	return len(d.cells)
}

func (d *RPCsTableData) GetColumnCount() int {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	if len(d.cells) == 0 {
		return 0
	}

	return len(d.cells[0])
}

func (d *RPCsTableData) redrawData() {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.cells = make([][]*tview.TableCell, len(d.knownRPCs))
	for i, rpc := range d.knownRPCs {
		d.cells[i] = []*tview.TableCell{tview.NewTableCell(rpc.URL), tview.NewTableCell(rpc.Moniker)}
	}
}
