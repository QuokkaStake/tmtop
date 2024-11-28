package types

import (
	"encoding/base64"
	"errors"
	"math/big"
	"strings"

	"github.com/cometbft/cometbft/crypto/ed25519"
	"github.com/cometbft/cometbft/p2p"
)

func ValidatorsWithLatestRoundFromTendermintResponse(
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
			return nil, errors.New("error setting string")
		}

		pubkey := make([]byte, 0, p2p.IDByteLength)
		defer func() {
			if perr := recover(); perr != nil {
				panic("BAD BASE64: " + validator.PubKey.PubKeyBase64)
			}
		}()
		pubkey, err := base64.StdEncoding.DecodeString(validator.PubKey.PubKeyBase64)
		if err != nil {
			return nil, err
		}

		validators[index] = ValidatorWithRoundVote{
			Validator: Validator{
				Address:     validator.Address,
				VotingPower: vp,
				PubKey:      pubkey,
				PeerID:      p2p.PubKeyToID(ed25519.PubKey(pubkey)),
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

func ValidatorsWithAllRoundsFromTendermintResponse(
	consensus *ConsensusStateResponse,
	tendermintValidators []TendermintValidator,
) (ValidatorsWithAllRoundsVotes, error) {
	validators := make(Validators, len(tendermintValidators))
	for index, validator := range tendermintValidators {
		vp := new(big.Int)
		vp, ok := vp.SetString(validator.VotingPower, 10)
		if !ok {
			return ValidatorsWithAllRoundsVotes{}, errors.New("error setting string")
		}

		pubkey, err := base64.StdEncoding.DecodeString(validator.PubKey.PubKeyBase64)
		if err != nil {
			return ValidatorsWithAllRoundsVotes{}, err
		}

		validators[index] = Validator{
			Address:     validator.Address,
			VotingPower: vp,
			PubKey:      pubkey,
			PeerID:      p2p.PubKeyToID(ed25519.PubKey(pubkey)),
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

	roundsVotes := make([]RoundVotes, len(consensus.Result.RoundState.HeightVoteSet))

	for round, roundHeightVoteSet := range consensus.Result.RoundState.HeightVoteSet {
		currentRoundVotes := make(RoundVotes, len(roundHeightVoteSet.Prevotes))

		for index, prevote := range roundHeightVoteSet.Prevotes {
			precommit := roundHeightVoteSet.Precommits[index]
			validator := tendermintValidators[index]
			currentRoundVotes[index] = RoundVote{
				Address:    validator.Address,
				Precommit:  VoteFromString(precommit),
				Prevote:    VoteFromString(prevote),
				IsProposer: validator.Address == consensus.Result.RoundState.Proposer.Address,
			}
		}

		roundsVotes[round] = currentRoundVotes
	}

	return ValidatorsWithAllRoundsVotes{
		Validators:  validators,
		RoundsVotes: roundsVotes,
	}, nil
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
