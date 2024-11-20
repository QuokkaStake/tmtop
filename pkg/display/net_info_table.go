package display

import (
	"fmt"
	"main/pkg/types"
	"main/pkg/utils"
	"slices"
	"strings"

	"github.com/rivo/tview"
)

type NetInfoTableData struct {
	tview.TableContentReadOnly

	NetInfo *types.NetInfo

	cells [][]*tview.TableCell
	mutex *utils.NoopLocker
}

func NewNetInfoTableData() *NetInfoTableData {
	return &NetInfoTableData{
		cells: [][]*tview.TableCell{},
		mutex: &utils.NoopLocker{},
	}
}

func (d *NetInfoTableData) GetCell(row, column int) *tview.TableCell {
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

func (d *NetInfoTableData) GetRowCount() int {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	return len(d.cells)
}

func (d *NetInfoTableData) GetColumnCount() int {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	if len(d.cells) == 0 {
		return 0
	}

	return len(d.cells[0])
}

func (d *NetInfoTableData) SetNetInfo(netInfo *types.NetInfo) {
	d.NetInfo = netInfo
	d.redrawData()
}

func (d *NetInfoTableData) redrawData() {
	cells := d.makeCells()

	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.cells = cells
}

func (d *NetInfoTableData) makeCells() [][]*tview.TableCell {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	if d.NetInfo == nil {
		return nil
	}

	cells := make([][]*tview.TableCell, len(d.NetInfo.Peers)+2)

	slices.SortFunc(d.NetInfo.Peers, func(a, b types.Peer) int {
		if a.ConnectionStatus.RecvMonitor.AvgRate > b.ConnectionStatus.RecvMonitor.AvgRate {
			return -1
		}
		return 1
	})

	cells[0] = []*tview.TableCell{
		tview.NewTableCell(""),
		tview.NewTableCell("IP"),
		tview.NewTableCell("Moniker"),
		tview.NewTableCell("Duration"),
		tview.NewTableCell("Send (cur)"),
		tview.NewTableCell("Recv (cur)"),
		tview.NewTableCell("Send (avg)"),
		tview.NewTableCell("Recv (avg)"),
		tview.NewTableCell("Node ID"),
		tview.NewTableCell("Version"),
		tview.NewTableCell("Proto"),
		tview.NewTableCell("RPC"),
	}

	cells[1] = []*tview.TableCell{
		tview.NewTableCell(""),
		tview.NewTableCell("=="),
		tview.NewTableCell("========"),
		tview.NewTableCell("========"),
		tview.NewTableCell("=========="),
		tview.NewTableCell("=========="),
		tview.NewTableCell("=========="),
		tview.NewTableCell("=========="),
		tview.NewTableCell("========================================"),
		tview.NewTableCell("======="),
		tview.NewTableCell("====="),
		tview.NewTableCell("==="),
	}

	for i, peer := range d.NetInfo.Peers {
		cells[i+2] = make([]*tview.TableCell, 12)

		duration := strings.Split(peer.ConnectionStatus.Duration.String(), ".")[0]

		direction := "in"
		if peer.IsOutbound {
			direction = "out"
		}

		cells[i+2][0] = tview.NewTableCell(direction)
		cells[i+2][1] = tview.NewTableCell(peer.RemoteIP)
		cells[i+2][2] = tview.NewTableCell(peer.NodeInfo.Moniker)
		cells[i+2][3] = tview.NewTableCell(duration)
		cells[i+2][4] = tview.NewTableCell(peer.ConnectionStatus.SendMonitor.CurRate.String() + "/s").SetAlign(tview.AlignRight)
		cells[i+2][5] = tview.NewTableCell(peer.ConnectionStatus.RecvMonitor.CurRate.String() + "/s").SetAlign(tview.AlignRight)
		cells[i+2][6] = tview.NewTableCell(peer.ConnectionStatus.SendMonitor.AvgRate.String() + "/s").SetAlign(tview.AlignRight)
		cells[i+2][7] = tview.NewTableCell(peer.ConnectionStatus.RecvMonitor.AvgRate.String() + "/s").SetAlign(tview.AlignRight)
		cells[i+2][8] = tview.NewTableCell(string(peer.NodeInfo.DefaultNodeID))
		cells[i+2][9] = tview.NewTableCell(peer.NodeInfo.Version)
		cells[i+2][10] = tview.NewTableCell(fmt.Sprintf("%v/%v/%v", peer.NodeInfo.ProtocolVersion.P2P, peer.NodeInfo.ProtocolVersion.Block, peer.NodeInfo.ProtocolVersion.App))
		cells[i+2][11] = tview.NewTableCell(peer.NodeInfo.Other.RPCAddress)
	}
	return cells
}
