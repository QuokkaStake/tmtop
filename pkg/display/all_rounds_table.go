package display

import (
    "fmt"
    "main/pkg/types"
    "main/pkg/utils"
    "sort"

    "github.com/gdamore/tcell/v2"
    "github.com/rivo/tview"
)

// BlockRoundHistory stores voting history for a block height
type BlockRoundHistory struct {
    Height int64
    Rounds map[int]map[string]types.RoundVote
}

type AllRoundsTableData struct {
    tview.TableContentReadOnly

    Validators    types.ValidatorsWithInfoAndAllRoundVotes
    DisableEmojis bool
    Transpose     bool
    CurrentHeight int64
    
    History        []BlockRoundHistory
    MaxHistorySize int
    
    cells [][]*tview.TableCell
    mutex *utils.NoopLocker
}

func NewAllRoundsTableData(disableEmojis bool, transpose bool) *AllRoundsTableData {
    return &AllRoundsTableData{
        DisableEmojis:  disableEmojis,
        Transpose:      transpose,
        History:        make([]BlockRoundHistory, 0),
        MaxHistorySize: 10,
        cells:          [][]*tview.TableCell{},
        mutex:          &utils.NoopLocker{},
    }
}

func (d *AllRoundsTableData) GetCell(row, column int) *tview.TableCell {
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

func (d *AllRoundsTableData) GetRowCount() int {
    d.mutex.RLock()
    defer d.mutex.RUnlock()

    return len(d.cells)
}

func (d *AllRoundsTableData) GetColumnCount() int {
    d.mutex.RLock()
    defer d.mutex.RUnlock()

    if len(d.cells) == 0 {
        return 0
    }

    return len(d.cells[0])
}

func (d *AllRoundsTableData) SetValidators(validators types.ValidatorsWithInfoAndAllRoundVotes, height int64) {
    d.mutex.Lock()
    
    // Update the history if height changed
    if height > 0 && d.CurrentHeight > 0 && height != d.CurrentHeight {
        // Store current data in history
        d.updateHistory(height)
    } else if d.CurrentHeight == 0 && height > 0 {
        d.CurrentHeight = height
    }
    
    d.Validators = validators
    d.mutex.Unlock()

    d.redrawData()
}

func (d *AllRoundsTableData) SetTranspose(transpose bool) {
    d.mutex.Lock()
    d.Transpose = transpose
    d.mutex.Unlock()
    
    d.redrawData()
}

func (d *AllRoundsTableData) updateHistory(newHeight int64) {
    // create history entry
    history := BlockRoundHistory{
        Height: d.CurrentHeight,
        Rounds: make(map[int]map[string]types.RoundVote),
    }
    
    // store votes by round
    for round, votes := range d.Validators.RoundsVotes {
        history.Rounds[round] = make(map[string]types.RoundVote)
        
        for i, vote := range votes {
            validatorID := fmt.Sprintf("validator-%d", i)
            
            if i < len(d.Validators.Validators) && 
               d.Validators.Validators[i].ChainValidator != nil &&
               d.Validators.Validators[i].ChainValidator.Address != "" {
                validatorID = d.Validators.Validators[i].ChainValidator.Address
            }
            
            history.Rounds[round][validatorID] = vote
        }
    }
    
    // add to history here
    d.History = append(d.History, history)
    if len(d.History) > d.MaxHistorySize {
        d.History = d.History[1:]
    }
    
    d.CurrentHeight = newHeight
}

func (d *AllRoundsTableData) redrawData() {
    cells := d.createCells()
    
    d.mutex.Lock()
    defer d.mutex.Unlock()
    d.cells = cells
}

// Helper function to get validator name
func getValidatorName(validator types.ValidatorWithChainValidator, index int) string {
    name := fmt.Sprintf("Validator %d", index)
    
    if validator.ChainValidator == nil {
        return name
    }
    
    if validator.ChainValidator.Moniker != "" {
        return validator.ChainValidator.Moniker
    }
    
    if validator.ChainValidator.Address != "" {
        addr := validator.ChainValidator.Address
        if len(addr) > 10 {
            addr = addr[:6] + "..." + addr[len(addr)-4:]
        }
        return addr
    }
    
    return name
}

// Create cells for the table
func (d *AllRoundsTableData) createCells() [][]*tview.TableCell {
    cells := [][]*tview.TableCell{}
    
    if d.Validators.Validators == nil || len(d.Validators.RoundsVotes) == 0 {
        return cells
    }
    
    // height+round structure definition
    type HeightRound struct {
        Height int64
        Round  int
    }
    
    // Collect all height+round combinations
    allHeightRounds := []HeightRound{}
    
    // add current rounds here
    for round := range d.Validators.RoundsVotes {
        allHeightRounds = append(allHeightRounds, HeightRound{
            Height: d.CurrentHeight,
            Round:  round,
        })
    }
    
    // add historical rounds here
    for _, history := range d.History {
        for round := range history.Rounds {
            allHeightRounds = append(allHeightRounds, HeightRound{
                Height: history.Height,
                Round:  round,
            })
        }
    }
    
    // sort by height (descending) then round (ascending)
    sort.Slice(allHeightRounds, func(i, j int) bool {
        if allHeightRounds[i].Height != allHeightRounds[j].Height {
            return allHeightRounds[i].Height > allHeightRounds[j].Height
        }
        return allHeightRounds[i].Round < allHeightRounds[j].Round
    })
    
    // Create header row with bold text
    headerRow := []*tview.TableCell{
        tview.NewTableCell("Validator").
            SetSelectable(false).
            SetStyle(tcell.StyleDefault.Bold(true)),
    }
    
	for _, hr := range allHeightRounds {
		// Format height to show only last 4 digits -> this can be adjusted by preference
		heightStr := fmt.Sprintf("%d", hr.Height)
		if len(heightStr) > 4 {
			heightStr = heightStr[len(heightStr)-4:]
		}
		
		headerCell := tview.NewTableCell(fmt.Sprintf("H%sR%d", heightStr, hr.Round)).
			SetSelectable(false).
			SetStyle(tcell.StyleDefault.Bold(true))
		headerRow = append(headerRow, headerCell)
	}
    cells = append(cells, headerRow)
    
    // Create validator rows here
    for i, validator := range d.Validators.Validators {
        row := []*tview.TableCell{}
        
        // enumerated validator name
        name := getValidatorName(validator, i)
        validatorCell := tview.NewTableCell(fmt.Sprintf("%d. %s", i+1, name))
        row = append(row, validatorCell)
        
        // Add vote cells here
        for _, hr := range allHeightRounds {
            cell := tview.NewTableCell("")
            found := false
            var vote types.RoundVote
            
            // Get vote from the current height or history
            if hr.Height == d.CurrentHeight {
                votes := d.Validators.RoundsVotes[hr.Round]
                if i < len(votes) {
                    vote = votes[i]
                    found = true
                }
            } else {
                // find in history
                for _, h := range d.History {
                    if h.Height == hr.Height {
                        roundVotes, ok := h.Rounds[hr.Round]
                        if !ok {
                            continue
                        }
                        
                        // match by the validator address here
                        validatorID := ""
                        if validator.ChainValidator != nil {
                            validatorID = validator.ChainValidator.Address
                        }
                        
                        if v, ok := roundVotes[validatorID]; ok {
                            vote = v
                            found = true
                            break
                        }
                    }
                }
            }
            
            if found {
                cell = tview.NewTableCell(vote.Serialize(d.DisableEmojis))
                if vote.IsProposer {
                    cell.SetBackgroundColor(tcell.ColorForestGreen)
                }
            }
            
            row = append(row, cell)
        }
        
        cells = append(cells, row)
    }
    
    return cells
}

func (d *AllRoundsTableData) HandleKey(event *tcell.EventKey) bool {
    return false
}