package types

import (
	"fmt"
	"main/pkg/utils"
	"math/big"
	"strconv"
)

type Validator struct {
	Index              int
	Address            string
	VotingPower        *big.Int
	VotingPowerPercent *big.Float
}

type Validators []Validator

type RoundVote struct {
	Address    string
	Prevote    Vote
	Precommit  Vote
	IsProposer bool
}
type RoundVotes []RoundVote

type ValidatorWithRoundVote struct {
	Validator Validator
	RoundVote RoundVote
}

type ValidatorsWithRoundVote []ValidatorWithRoundVote

type ValidatorsWithAllRoundsVotes struct {
	Validators  []Validator
	RoundsVotes []RoundVotes
}

func (v Validators) GetTotalVotingPower() *big.Int {
	sum := big.NewInt(0)

	for _, validator := range v {
		sum = sum.Add(sum, validator.VotingPower)
	}

	return sum
}
func (v ValidatorsWithRoundVote) GetTotalVotingPower() *big.Int {
	sum := big.NewInt(0)

	for _, validator := range v {
		sum = sum.Add(sum, validator.Validator.VotingPower)
	}

	return sum
}

func (v ValidatorsWithRoundVote) GetTotalVotingPowerPrevotedPercent(countDisagreeing bool) *big.Float {
	prevoted := big.NewInt(0)
	totalVP := big.NewInt(0)

	for _, validator := range v {
		totalVP = totalVP.Add(totalVP, validator.Validator.VotingPower)
		if validator.RoundVote.Prevote == Voted || (countDisagreeing && validator.RoundVote.Prevote == VotedZero) {
			prevoted = prevoted.Add(prevoted, validator.Validator.VotingPower)
		}
	}

	votingPowerPercent := big.NewFloat(0).SetInt(prevoted)
	votingPowerPercent = votingPowerPercent.Quo(votingPowerPercent, big.NewFloat(0).SetInt(totalVP))
	votingPowerPercent = votingPowerPercent.Mul(votingPowerPercent, big.NewFloat(100))

	return votingPowerPercent
}

func (v ValidatorsWithRoundVote) GetTotalVotingPowerPrecommittedPercent(countDisagreeing bool) *big.Float {
	precommitted := big.NewInt(0)
	totalVP := big.NewInt(0)

	for _, validator := range v {
		totalVP = totalVP.Add(totalVP, validator.Validator.VotingPower)
		if validator.RoundVote.Precommit == Voted || (countDisagreeing && validator.RoundVote.Precommit == VotedZero) {
			precommitted = precommitted.Add(precommitted, validator.Validator.VotingPower)
		}
	}

	votingPowerPercent := big.NewFloat(0).SetInt(precommitted)
	votingPowerPercent = votingPowerPercent.Quo(votingPowerPercent, big.NewFloat(0).SetInt(totalVP))
	votingPowerPercent = votingPowerPercent.Mul(votingPowerPercent, big.NewFloat(100))

	return votingPowerPercent
}

type ValidatorWithInfo struct {
	Validator      Validator
	RoundVote      RoundVote
	ChainValidator *ChainValidator
}

func (v ValidatorWithInfo) Serialize() string {
	name := v.Validator.Address
	if v.ChainValidator != nil {
		name = v.ChainValidator.Moniker
		if v.ChainValidator.AssignedAddress != "" {
			name = "🔑 " + name
		}
	}

	return fmt.Sprintf(
		" %s %s %s %s%% %s ",
		v.RoundVote.Prevote.Serialize(),
		v.RoundVote.Precommit.Serialize(),
		utils.RightPadAndTrim(strconv.Itoa(v.Validator.Index+1), 3),
		utils.RightPadAndTrim(fmt.Sprintf("%.2f", v.Validator.VotingPowerPercent), 6),
		utils.LeftPadAndTrim(name, 25),
	)
}

type ValidatorsWithInfo []ValidatorWithInfo

type ValidatorWithChainValidator struct {
	Validator      Validator
	ChainValidator *ChainValidator
}

type ValidatorsWithInfoAndAllRoundVotes struct {
	Validators  []ValidatorWithChainValidator
	RoundsVotes []RoundVotes
}
