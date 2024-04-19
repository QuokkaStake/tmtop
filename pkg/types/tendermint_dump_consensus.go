package types

type DumpConsensusStateResponse struct {
	Result *DumpConsensusStateResult `json:"result"`
}

type DumpConsensusStateResult struct {
	RoundState *DumpConsensusStateRoundState `json:"round_state"`
}

type DumpConsensusStateRoundState struct {
	Validators DumpConsensusStateRoundStateValidators `json:"validators"`
}
type DumpConsensusStateRoundStateValidators struct {
	Validators []TendermintValidator `json:"validators"`
}
