package types

import (
	"time"
)

type ConsensusStateResponse struct {
	Result *ConsensusStateResult `json:"result"`
}

type ConsensusStateResult struct {
	RoundState *ConsensusStateRoundState `json:"round_state"`
}

type ConsensusStateRoundState struct {
	HeightRoundStep string                   `json:"height/round/step"`
	StartTime       time.Time                `json:"start_time"`
	HeightVoteSet   []ConsensusHeightVoteSet `json:"height_vote_set"`
	Proposer        ConsensusStateProposer   `json:"proposer"`
}

type ConsensusHeightVoteSet struct {
	Round              int                   `json:"round"`
	Prevotes           []ConsensusVote       `json:"prevotes"`
	Precommits         []ConsensusVote       `json:"precommits"`
	PrevotesBitArray   ConsensusVoteBitArray `json:"prevotes_bit_array"`
	PrecommitsBitArray ConsensusVoteBitArray `json:"precommits_bit_array"`
}

type ConsensusStateProposer struct {
	Address string `json:"address"`
	Index   int    `json:"index"`
}

type ConsensusVote string
type ConsensusVoteBitArray string
