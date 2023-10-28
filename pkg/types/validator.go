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
	Prevote            Vote
	Precommit          Vote
	IsProposer         bool
}

type Validators []Validator

func (v Validators) GetTotalVotingPower() *big.Int {
	sum := big.NewInt(0)

	for _, validator := range v {
		sum = sum.Add(sum, validator.VotingPower)
	}

	return sum
}

func (v Validators) GetTotalVotingPowerPrevotedPercent(countDisagreeing bool) *big.Float {
	prevoted := big.NewInt(0)
	totalVP := big.NewInt(0)

	for _, validator := range v {
		totalVP = totalVP.Add(totalVP, validator.VotingPower)
		if validator.Prevote == Voted || (countDisagreeing && validator.Prevote == VotedZero) {
			prevoted = prevoted.Add(prevoted, validator.VotingPower)
		}
	}

	votingPowerPercent := big.NewFloat(0).SetInt(prevoted)
	votingPowerPercent = votingPowerPercent.Quo(votingPowerPercent, big.NewFloat(0).SetInt(totalVP))
	votingPowerPercent = votingPowerPercent.Mul(votingPowerPercent, big.NewFloat(100))

	return votingPowerPercent
}

func (v Validators) GetTotalVotingPowerPrecommittedPercent(countDisagreeing bool) *big.Float {
	precommitted := big.NewInt(0)
	totalVP := big.NewInt(0)

	for _, validator := range v {
		totalVP = totalVP.Add(totalVP, validator.VotingPower)
		if validator.Precommit == Voted || (countDisagreeing && validator.Precommit == VotedZero) {
			precommitted = precommitted.Add(precommitted, validator.VotingPower)
		}
	}

	votingPowerPercent := big.NewFloat(0).SetInt(precommitted)
	votingPowerPercent = votingPowerPercent.Quo(votingPowerPercent, big.NewFloat(0).SetInt(totalVP))
	votingPowerPercent = votingPowerPercent.Mul(votingPowerPercent, big.NewFloat(100))

	return votingPowerPercent
}

type ValidatorWithInfo struct {
	Validator      Validator
	ChainValidator *ChainValidator
}

func (v ValidatorWithInfo) Serialize() string {
	name := v.Validator.Address
	if v.ChainValidator != nil {
		name = v.ChainValidator.Moniker
	}

	return fmt.Sprintf(
		" %s %s %s %s%% %s ",
		v.Validator.Prevote.Serialize(),
		v.Validator.Precommit.Serialize(),
		utils.RightPadAndTrim(strconv.Itoa(v.Validator.Index+1), 3),
		utils.RightPadAndTrim(fmt.Sprintf("%.2f", v.Validator.VotingPowerPercent), 6),
		utils.LeftPadAndTrim(name, 25),
	)
}

type ValidatorsWithInfo []ValidatorWithInfo

func (v ValidatorsWithInfo) Serialise() []string {
	serialized := make([]string, len(v))

	for index, validator := range v {
		serialized[index] = validator.Serialize()
	}

	return serialized
}
