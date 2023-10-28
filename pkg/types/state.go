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
	Validators      *Validators
	ChainValidators *ChainValidators
	StartTime       time.Time
}

func NewState() *State {
	return &State{
		Height:          0,
		Round:           0,
		Step:            0,
		Validators:      nil,
		ChainValidators: nil,
		StartTime:       time.Now(),
	}
}

func (s *State) SetTendermintResponse(
	consensus *ConsensusStateResponse,
	tendermintValidators []TendermintValidator,
) error {
	validators, err := ValidatorsFromTendermintResponse(consensus, tendermintValidators)
	if err != nil {
		return err
	}

	hrsSplit := strings.Split(consensus.Result.RoundState.HeightRoundStep, "/")

	s.Validators = &validators
	s.Height = utils.MustParseInt64(hrsSplit[0])
	s.Round = utils.MustParseInt64(hrsSplit[1])
	s.Step = utils.MustParseInt64(hrsSplit[2])
	s.StartTime = consensus.Result.RoundState.StartTime

	return nil
}

func (s *State) SetChainValidators(validators *ChainValidators) {
	s.ChainValidators = validators
}

func (s *State) SerializeInfo() string {
	if s.Validators == nil {
		return ""
	}

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

func (s *State) GetValidatorsWithInfo() ValidatorsWithInfo {
	if s.Validators == nil {
		return ValidatorsWithInfo{}
	}

	validators := make(ValidatorsWithInfo, len(*s.Validators))

	for index, validator := range *s.Validators {
		validators[index] = ValidatorWithInfo{
			Validator: validator,
		}
	}

	if s.ChainValidators == nil {
		return validators
	}

	chainValidatorsMap := s.ChainValidators.ToMap()
	for index, validator := range *s.Validators {
		if chainValidator, ok := chainValidatorsMap[validator.Address]; ok {
			validators[index].ChainValidator = &chainValidator
		}
	}

	return validators
}
