package types

type ValidatorsResponse struct {
	Result *ValidatorsResult `json:"result"`
	Error  *ValidatorsError  `json:"error"`
}

type ValidatorsError struct {
	Message string `json:"message"`
	Data    string `json:"data"`
}

type ValidatorsResult struct {
	Count      string                `json:"count"`
	Total      string                `json:"total"`
	Validators []TendermintValidator `json:"validators"`
}

type TendermintValidator struct {
	Address     string                    `json:"address"`
	VotingPower string                    `json:"voting_power"`
	PubKey      TendermintValidatorPubKey `json:"pub_key"`
}

type TendermintValidatorPubKey struct {
	Type         string `json:"type"`
	PubKeyBase64 string `json:"value"`
}
