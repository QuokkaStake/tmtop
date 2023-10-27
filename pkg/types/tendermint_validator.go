package types

type ValidatorsResponse struct {
	Result *ValidatorsResult `json:"result"`
}

type ValidatorsResult struct {
	Count      string                `json:"count"`
	Total      string                `json:"total"`
	Validators []TendermintValidator `json:"validators"`
}

type TendermintValidator struct {
	Address     string `json:"address"`
	VotingPower string `json:"voting_power"`
}
