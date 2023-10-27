package state

import "main/pkg/types"

type State struct {
	Proposer   types.Validator
	Validators []types.Validator
	Height     int64
	Round      int64
	Step       int64
}

func NewState() {

}
