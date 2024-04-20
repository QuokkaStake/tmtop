package display

import (
	"fmt"
	"main/pkg/types"
	"strconv"

	"github.com/gdamore/tcell/v2"

	"github.com/rivo/tview"
)

type AllRoundsTableHeader struct {
	Title string
	Value func(index int, votes types.RoundVotes, validators []types.ValidatorWithChainValidator) string
}

type AllRoundsTableData struct {
	tview.TableContentReadOnly

	Headers       []AllRoundsTableHeader
	Validators    types.ValidatorsWithInfoAndAllRoundVotes
	DisableEmojis bool
}

func NewAllRoundsTableData(disableEmojis bool) *AllRoundsTableData {
	headers := []AllRoundsTableHeader{
		{
			Title: "round #",
			Value: func(index int, votes types.RoundVotes, validators []types.ValidatorWithChainValidator) string {
				return strconv.Itoa(index)
			},
		},
		{
			Title: "prevoted count/total",
			Value: func(index int, votes types.RoundVotes, validators []types.ValidatorWithChainValidator) string {
				count := 0
				for _, vote := range votes {
					if vote.Prevote != types.VotedNil {
						count += 1
					}
				}
				return fmt.Sprintf("%d/%d", count, len(validators))
			},
		},
		{
			Title: "prevoted %",
			Value: func(index int, votes types.RoundVotes, validators []types.ValidatorWithChainValidator) string {
				prevotedPercent := types.
					ValidatorsWithRoundVoteFrom(validators, votes).
					GetTotalVotingPowerPrevotedPercent(true)
				return fmt.Sprintf("%.2f%%", prevotedPercent)
			},
		},
		{
			Title: "prevoted agreed %",
			Value: func(index int, votes types.RoundVotes, validators []types.ValidatorWithChainValidator) string {
				prevotedPercent := types.
					ValidatorsWithRoundVoteFrom(validators, votes).
					GetTotalVotingPowerPrevotedPercent(false)
				return fmt.Sprintf("%.2f%%", prevotedPercent)
			},
		},
		{
			Title: "precommitted %",
			Value: func(index int, votes types.RoundVotes, validators []types.ValidatorWithChainValidator) string {
				prevotedPercent := types.
					ValidatorsWithRoundVoteFrom(validators, votes).
					GetTotalVotingPowerPrecommittedPercent(true)
				return fmt.Sprintf("%.2f%%", prevotedPercent)
			},
		},
	}

	return &AllRoundsTableData{
		Validators:    types.ValidatorsWithInfoAndAllRoundVotes{},
		Headers:       headers,
		DisableEmojis: disableEmojis,
	}
}

func (d *AllRoundsTableData) GetCell(row, column int) *tview.TableCell {
	headersCount := len(d.Headers)

	// Table header.
	if row < headersCount {
		// First column is title.
		if column == 0 {
			return tview.
				NewTableCell(d.Headers[row].Title).
				SetAlign(tview.AlignCenter).
				SetStyle(tcell.StyleDefault.Bold(true))
		}

		roundVotes := d.Validators.RoundsVotes[column-1]

		return tview.
			NewTableCell(d.Headers[row].Value(column-1, roundVotes, d.Validators.Validators)).
			SetAlign(tview.AlignCenter).
			SetStyle(tcell.StyleDefault.Bold(true))
	}

	// First column is always validators list.
	if column == 0 {
		text := d.Validators.Validators[row-headersCount].Serialize()
		cell := tview.NewTableCell(text)
		return cell
	}

	roundVotes := d.Validators.RoundsVotes[column-1]
	roundVote := roundVotes[row-headersCount]
	text := roundVote.Serialize(d.DisableEmojis)

	cell := tview.NewTableCell(text)

	if roundVote.IsProposer {
		cell.SetBackgroundColor(tcell.ColorForestGreen)
	}

	return cell
}

func (d *AllRoundsTableData) GetRowCount() int {
	return len(d.Validators.Validators) + len(d.Headers)
}

func (d *AllRoundsTableData) GetColumnCount() int {
	return len(d.Validators.RoundsVotes) + 1 // first column is header
}

func (d *AllRoundsTableData) SetValidators(validators types.ValidatorsWithInfoAndAllRoundVotes) {
	d.Validators = validators
}
