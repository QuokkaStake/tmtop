package types

import "time"

type TendermintBlockResponse struct {
	Result TendermintBlockResult `json:"result"`
}

type TendermintBlockResult struct {
	Block TendermintBlock `json:"block"`
}

type TendermintBlock struct {
	Header TendermintBlockHeader `json:"header"`
}

type TendermintBlockHeader struct {
	Height string    `json:"height"`
	Time   time.Time `json:"time"`
}
