package types

import (
	"fmt"
	"main/pkg/utils"
	"strings"
	"time"
)

type State struct {
	Height          int64
	Round           int64
	Step            int64
	Validators      Validators
	ChainValidators ChainValidators
	StartTime       time.Time
}

func (s State) SerializeInfo() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf(" height=%d round=%d step=%d\n", s.Height, s.Round, s.Step))
	sb.WriteString(fmt.Sprintf(" block time: %s\n", utils.ZeroOrPositiveDuration(time.Since(s.StartTime))))
	sb.WriteString(fmt.Sprintf(
		" prevote consensus (total/agreeing): %.2f / %.2f\n",
		s.Validators.GetTotalVotingPowerPrevotedPercent(true),
		s.Validators.GetTotalVotingPowerPrevotedPercent(false),
	))
	sb.WriteString(fmt.Sprintf(
		" precommit consensus (total/agreeing): %.2f / %.2f\n",
		s.Validators.GetTotalVotingPowerPrecommittedPercent(true),
		s.Validators.GetTotalVotingPowerPrecommittedPercent(false),
	))

	return sb.String()
}

func (s State) GetValidatorsWithInfo() ValidatorsWithInfo {
	validators := make(ValidatorsWithInfo, len(s.Validators))
	chainValidatorsMap := s.ChainValidators.ToMap()

	for index, validator := range s.Validators {
		validators[index] = ValidatorWithInfo{
			Validator: validator,
		}

		if chainValidator, ok := chainValidatorsMap[validator.Address]; ok {
			validators[index].ChainValidator = &chainValidator
		}
	}

	return validators
}
