package types

import "encoding/json"

type Genesis struct {
	AppState AppState `json:"app_state"`
}

type AppState struct {
	Staking json.RawMessage `json:"staking"`
}
