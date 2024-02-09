package types

import (
	"fmt"
	"math/big"
	"strings"
)

func ValidatorsFromTendermintResponse(
	consensus *ConsensusStateResponse,
	tendermintValidators []TendermintValidator,
	round int64,
) (ValidatorsWithRoundVote, error) {
	lastHeightVoteSet := consensus.Result.RoundState.HeightVoteSet[round]
	validators := make(ValidatorsWithRoundVote, len(lastHeightVoteSet.Prevotes))

	for index, prevote := range lastHeightVoteSet.Prevotes {
		precommit := lastHeightVoteSet.Precommits[index]
		validator := tendermintValidators[index]

		vp := new(big.Int)
		vp, ok := vp.SetString(validator.VotingPower, 10)
		if !ok {
			return nil, fmt.Errorf("error setting string")
		}

		validators[index] = ValidatorWithRoundVote{
			Validator: Validator{
				Address:     validator.Address,
				VotingPower: vp,
			},
			RoundVote: RoundVote{
				Address:    validator.Address,
				Precommit:  VoteFromString(precommit),
				Prevote:    VoteFromString(prevote),
				IsProposer: validator.Address == consensus.Result.RoundState.Proposer.Address,
			},
		}
	}

	totalVP := validators.GetTotalVotingPower()

	for index, validator := range validators {
		validators[index].Validator.Index = index

		votingPowerPercent := big.NewFloat(0).SetInt(validator.Validator.VotingPower)
		votingPowerPercent = votingPowerPercent.Quo(votingPowerPercent, big.NewFloat(0).SetInt(totalVP))
		votingPowerPercent = votingPowerPercent.Mul(votingPowerPercent, big.NewFloat(100))

		validators[index].Validator.VotingPowerPercent = votingPowerPercent
	}

	return validators, nil
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
