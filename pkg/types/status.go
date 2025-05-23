package types

type TendermintStatusResponse struct {
	Result TendermintStatusResult `json:"result"`
}

type TendermintStatusResult struct {
	NodeInfo      TendermintNodeInfo      `json:"node_info"`
	ValidatorInfo TendermintValidatorInfo `json:"validator_info"`
}

type TendermintNodeInfo struct {
	Version string `json:"version"`
	Network string `json:"network"`
}

type TendermintValidatorInfo struct {
	Address string `json:"address"`
}
