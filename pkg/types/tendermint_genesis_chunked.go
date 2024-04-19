package types

type TendermintGenesisChunkResponse struct {
	Result *TendermintGenesisChunkResult `json:"result"`
}

type TendermintGenesisChunkResult struct {
	Chunk string `json:"string"`
	Total string `json:"total"`
	Data  []byte `json:"data"`
}
