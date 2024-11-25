package types

type TendermintStatusResponse struct {
	Result TendermintStatusResult `json:"result"`
}

type TendermintStatusResult struct {
	NodeInfo TendermintNodeInfo `json:"node_info"`
}

type TendermintNodeInfo struct {
	ID      string `json:"id"`
	Version string `json:"version"`
	Network string `json:"network"`
	Moniker string `json:"moniker"`
	Other   struct {
		RPCAddress string `json:"rpc_address"`
	} `json:"other"`
}

type TendermintValidatorInfo struct {
	Address     string `json:"address"`
	VotingPower string `json:"voting_power"`
}
