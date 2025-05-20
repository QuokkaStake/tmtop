package types

import (
	"fmt"
	"main/pkg/utils"
	"strings"
	"time"
)

type State struct {
	Height                       int64
	Round                        int64
	Step                         int64
	Validators                   *ValidatorsWithRoundVote
	ValidatorsWithAllRoundsVotes *ValidatorsWithAllRoundsVotes
	ChainValidators              *ChainValidators
	NodeStatus                   *TendermintStatusResult
	StartTime                    time.Time
	Upgrade                      *Upgrade
	BlockTime                    time.Duration

	ConsensusStateError  error
	ValidatorsError      error
	ChainValidatorsError error
	UpgradePlanError     error
	StatusError          error
}

func NewState() *State {
	return &State{
		Height:          0,
		Round:           0,
		Step:            0,
		Validators:      nil,
		ChainValidators: nil,
		StartTime:       time.Now(),
		BlockTime:       0,
	}
}

func (s *State) SetTendermintResponse(
	consensus *ConsensusStateResponse,
	tendermintValidators []TendermintValidator,
) error {
	hrsSplit := strings.Split(consensus.Result.RoundState.HeightRoundStep, "/")

	s.Height = utils.MustParseInt64(hrsSplit[0])
	s.Round = utils.MustParseInt64(hrsSplit[1])
	s.Step = utils.MustParseInt64(hrsSplit[2])
	s.StartTime = consensus.Result.RoundState.StartTime

	validators, err := ValidatorsWithLatestRoundFromTendermintResponse(consensus, tendermintValidators, s.Round)
	if err != nil {
		return err
	}

	s.Validators = &validators

	validatorsWithAllRounds, err := ValidatorsWithAllRoundsFromTendermintResponse(consensus, tendermintValidators)
	if err != nil {
		return err
	}

	s.ValidatorsWithAllRoundsVotes = &validatorsWithAllRounds

	return nil
}

func (s *State) SetChainValidators(validators *ChainValidators) {
	s.ChainValidators = validators
}

func (s *State) SetNodeStatus(status *TendermintStatusResult) {
	s.NodeStatus = status
}

func (s *State) SetUpgrade(upgrade *Upgrade) {
	s.Upgrade = upgrade
}

func (s *State) SetBlockTime(blockTime time.Duration) {
	s.BlockTime = blockTime
}

func (s *State) SetConsensusStateError(err error) {
	s.ConsensusStateError = err
}

func (s *State) SetValidatorsError(err error) {
	s.ValidatorsError = err
}

func (s *State) SetUpgradePlanError(err error) {
	s.UpgradePlanError = err
}

func (s *State) SetStatusError(err error) {
	s.StatusError = err
}

func (s *State) SerializeConsensus(timezone *time.Location) string {
	if s.ConsensusStateError != nil {
		return fmt.Sprintf(" consensus state error: %s", s.ConsensusStateError)
	}

	if s.Validators == nil {
		return ""
	}

	var sb strings.Builder

	sb.WriteString(fmt.Sprintf(" height=%d round=%d step=%d\n", s.Height, s.Round, s.Step))
	sb.WriteString(fmt.Sprintf(
		" block time: %s (%s)\n",
		utils.ZeroOrPositiveDuration(utils.SerializeDuration(time.Since(s.StartTime))),
		utils.SerializeTime(s.StartTime.In(timezone)),
	))
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

	prevoted := 0
	precommitted := 0
	prevotedAgreed := 0
	precommittedAgreed := 0

	for _, validator := range *s.Validators {
		if validator.RoundVote.Prevote != VotedNil {
			prevoted += 1
		}
		if validator.RoundVote.Precommit != VotedNil {
			precommitted += 1
		}

		if validator.RoundVote.Prevote == Voted {
			prevotedAgreed += 1
		}

		if validator.RoundVote.Precommit == Voted {
			precommittedAgreed += 1
		}
	}

	sb.WriteString(fmt.Sprintf(
		" prevoted/precommitted: %d/%d (out of %d)\n",
		prevoted,
		precommitted,
		len(*s.Validators),
	))
	sb.WriteString(fmt.Sprintf(
		" prevoted/precommitted agreed: %d/%d (out of %d)\n",
		prevotedAgreed,
		precommittedAgreed,
		len(*s.Validators),
	))

	sb.WriteString(fmt.Sprintf(" last updated at: %s\n", utils.SerializeTime(time.Now().In(timezone))))

	return sb.String()
}

func (s *State) SerializeChainInfo(timezone *time.Location) string {
	var sb strings.Builder

	if s.StatusError != nil {
		sb.WriteString(fmt.Sprintf(" chain info fetch error: %s\n", s.StatusError.Error()))
	} else if s.NodeStatus != nil {
		sb.WriteString(fmt.Sprintf(" chain name: %s\n", s.NodeStatus.NodeInfo.Network))
		sb.WriteString(fmt.Sprintf(" tendermint version: v%s\n", s.NodeStatus.NodeInfo.Version))

		if s.BlockTime != 0 {
			sb.WriteString(fmt.Sprintf(" avg block time: %s\n", utils.SerializeDuration(s.BlockTime)))
		}
	}

	if s.UpgradePlanError != nil {
		sb.WriteString(fmt.Sprintf(" upgrade plan fetch error: %s\n", s.UpgradePlanError))
	} else if s.Upgrade == nil {
		sb.WriteString(" no chain upgrade scheduled\n")
	} else {
		sb.WriteString(s.SerializeUpgradeInfo(timezone))
	}

	return sb.String()
}

func (s *State) SerializeUpgradeInfo(timezone *time.Location) string {
	var sb strings.Builder

	if s.Upgrade.Height+1 == s.Height {
		sb.WriteString(" upgrade in progress...\n")
		return sb.String()
	}

	if s.Upgrade.Height+1 < s.Height {
		sb.WriteString(fmt.Sprintf(
			" chain upgrade %s applied at block %d\n",
			s.Upgrade.Name,
			s.Upgrade.Height,
		))

		sb.WriteString(fmt.Sprintf(
			" blocks since upgrade: %d\n",
			s.Height-s.Upgrade.Height,
		))

		if s.BlockTime == 0 {
			return sb.String()
		}

		upgradeTime := utils.CalculateTimeTillBlock(s.Height, s.Upgrade.Height, s.BlockTime)
		sb.WriteString(fmt.Sprintf(
			" time since upgrade: %s\n",
			utils.SerializeDuration(time.Since(upgradeTime)),
		))

		sb.WriteString(fmt.Sprintf(" upgrade approximate time: %s\n", utils.SerializeTime(upgradeTime.In(timezone))))
		return sb.String()
	}

	sb.WriteString(fmt.Sprintf(
		" chain upgrade %s scheduled at block %d\n",
		s.Upgrade.Name,
		s.Upgrade.Height,
	))

	sb.WriteString(fmt.Sprintf(
		" blocks till upgrade: %d\n",
		s.Upgrade.Height-s.Height,
	))

	if s.BlockTime == 0 {
		return sb.String()
	}

	upgradeTime := utils.CalculateTimeTillBlock(s.Height, s.Upgrade.Height, s.BlockTime)

	sb.WriteString(fmt.Sprintf(
		" time till upgrade: %s\n",
		utils.SerializeDuration(time.Until(upgradeTime)),
	))

	sb.WriteString(fmt.Sprintf(" upgrade estimated time: %s\n", utils.SerializeTime(upgradeTime.In(timezone))))

	return sb.String()
}

func (s *State) SerializeProgressbar(width int, height int, prefix string, progress int) string {
	progressBar := ProgressBar{
		Width:    width,
		Height:   height,
		Progress: progress,
		Prefix:   prefix,
	}

	return progressBar.Serialize()
}

func (s *State) SerializePrevotesProgressbar(width int, height int) string {
	if s.Validators == nil {
		return ""
	}

	prevotePercent := s.Validators.GetTotalVotingPowerPrevotedPercent(true)
	prevotePercentFloat, _ := prevotePercent.Float64()
	prevotePercentInt := int(prevotePercentFloat)

	return s.SerializeProgressbar(width, height, "Prevotes: ", prevotePercentInt)
}

func (s *State) SerializePrecommitsProgressbar(width int, height int) string {
	if s.Validators == nil {
		return ""
	}

	precommitPercent := s.Validators.GetTotalVotingPowerPrecommittedPercent(true)
	precommitPercentFloat, _ := precommitPercent.Float64()
	precommitPercentInt := int(precommitPercentFloat)

	return s.SerializeProgressbar(width, height, "Precommits: ", precommitPercentInt)
}

func (s *State) GetValidatorsWithInfo() ValidatorsWithInfo {
	if s.Validators == nil {
		return ValidatorsWithInfo{}
	}

	validators := make(ValidatorsWithInfo, len(*s.Validators))

	for index, validator := range *s.Validators {
		validators[index] = ValidatorWithInfo{
			Validator: validator.Validator,
			RoundVote: validator.RoundVote,
		}
	}

	if s.ChainValidators == nil {
		return validators
	}

	chainValidatorsMap := s.ChainValidators.ToMap()
	for index, validator := range *s.Validators {
		if chainValidator, ok := chainValidatorsMap[validator.Validator.Address]; ok {
			validators[index].ChainValidator = &chainValidator
		}
	}

	return validators
}

func (s *State) GetValidatorsWithInfoAndAllRoundVotes() ValidatorsWithInfoAndAllRoundVotes {
	if s.ValidatorsWithAllRoundsVotes == nil {
		return ValidatorsWithInfoAndAllRoundVotes{}
	}

	validators := make([]ValidatorWithChainValidator, len(s.ValidatorsWithAllRoundsVotes.Validators))

	for index, validator := range s.ValidatorsWithAllRoundsVotes.Validators {
		validators[index] = ValidatorWithChainValidator{
			Validator: validator,
		}
	}

	if s.ChainValidators == nil {
		return ValidatorsWithInfoAndAllRoundVotes{
			Validators:  validators,
			RoundsVotes: s.ValidatorsWithAllRoundsVotes.RoundsVotes,
		}
	}

	chainValidatorsMap := s.ChainValidators.ToMap()
	for index, validator := range s.ValidatorsWithAllRoundsVotes.Validators {
		if chainValidator, ok := chainValidatorsMap[validator.Address]; ok {
			validators[index].ChainValidator = &chainValidator
		}
	}

	return ValidatorsWithInfoAndAllRoundVotes{
		Validators:  validators,
		RoundsVotes: s.ValidatorsWithAllRoundsVotes.RoundsVotes,
	}
}
