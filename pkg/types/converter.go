package types

import (
	"fmt"
	"main/pkg/utils"
	"math/big"
	"strings"
)

func RenderIntoFromTendermintResponse(
	consensus *ConsensusStateResponse,
	tendermintValidators []TendermintValidator,
	chainValidators ChainValidators,
) RenderInfo {
	hrsSplit := strings.Split(consensus.Result.RoundState.HeightRoundStep, "/")

	return RenderInfo{
		Validators: ValidatorsFromTendermintResponse(consensus, tendermintValidators, chainValidators),
		Height:     utils.MustParseInt64(hrsSplit[0]),
		Round:      utils.MustParseInt64(hrsSplit[1]),
		Step:       utils.MustParseInt64(hrsSplit[2]),
		StartTime:  consensus.Result.RoundState.StartTime,
	}
}

func ValidatorsFromTendermintResponse(
	consensus *ConsensusStateResponse,
	tendermintValidators []TendermintValidator,
	chainValidators ChainValidators,
) Validators {
	chainValidatorsMap := chainValidators.ToMap()

	validators := make(Validators, len(consensus.Result.RoundState.HeightVoteSet[0].Prevotes))

	for index, prevote := range consensus.Result.RoundState.HeightVoteSet[0].Prevotes {
		precommit := consensus.Result.RoundState.HeightVoteSet[0].Precommits[index]
		validator := tendermintValidators[index]

		vp := new(big.Int)
		vp, ok := vp.SetString(validator.VotingPower, 10)
		if !ok {
			fmt.Println("SetString: error")
			panic("no")
		}

		validators[index] = Validator{
			Address:     validator.Address,
			Moniker:     validator.Address,
			Precommit:   VoteFromString(precommit),
			Prevote:     VoteFromString(prevote),
			VotingPower: vp,
			IsProposer:  validator.Address == consensus.Result.RoundState.Proposer.Address,
		}

		if chainValidator, ok := chainValidatorsMap[validator.Address]; ok {
			validators[index].Moniker = chainValidator.Moniker
		}
	}

	totalVP := validators.GetTotalVotingPower()

	for index, validator := range validators {
		validators[index].Index = index

		votingPowerPercent := big.NewFloat(0).SetInt(validator.VotingPower)
		votingPowerPercent = votingPowerPercent.Quo(votingPowerPercent, big.NewFloat(0).SetInt(totalVP))
		votingPowerPercent = votingPowerPercent.Mul(votingPowerPercent, big.NewFloat(100))

		validators[index].VotingPowerPercent = votingPowerPercent
	}

	return validators
}

func VoteFromString(source ConsensusVote) Vote {
	if source == "nil-Vote" {
		return VotedNil
	}

	if strings.Contains(string(source), "SIGNED_MSG_TYPE_PREVOTE(Prevote) 000000000000") {
		return VotedZero
	}

	return Voted
}
