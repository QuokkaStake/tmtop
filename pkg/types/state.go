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
	ChainInfo       *TendermintStatusResult
	StartTime       time.Time
	Upgrade         *Upgrade
	BlockTime       time.Duration
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

	validators, err := ValidatorsFromTendermintResponse(consensus, tendermintValidators, s.Round)
	if err != nil {
		return err
	}

	s.Validators = &validators

	return nil
}

func (s *State) SetChainValidators(validators *ChainValidators) {
	s.ChainValidators = validators
}

func (s *State) SetChainInfo(info *TendermintStatusResult) {
	s.ChainInfo = info
}

func (s *State) SetUpgrade(upgrade *Upgrade) {
	s.Upgrade = upgrade
}

func (s *State) SetBlockTime(blockTime time.Duration) {
	s.BlockTime = blockTime
}

func (s *State) SerializeConsensus() string {
	if s.Validators == nil {
		return ""
	}

	var sb strings.Builder

	sb.WriteString(fmt.Sprintf(" height=%d round=%d step=%d\n", s.Height, s.Round, s.Step))
	sb.WriteString(fmt.Sprintf(
		" block time: %s\n",
		utils.ZeroOrPositiveDuration(utils.SerializeDuration(time.Since(s.StartTime))),
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

	for _, validator := range *s.Validators {
		if validator.Prevote != VotedNil {
			prevoted += 1
		}
		if validator.Precommit != VotedNil {
			precommitted += 1
		}
	}

	sb.WriteString(fmt.Sprintf(
		" prevoted/precommitted: %d/%d (out of %d)\n",
		prevoted,
		precommitted,
		len(*s.Validators),
	))

	sb.WriteString(fmt.Sprintf(" last updated at: %s\n", utils.SerializeTime(time.Now())))

	return sb.String()
}

func (s *State) SerializeChainInfo() string {
	if s.ChainInfo == nil {
		return ""
	}

	var sb strings.Builder

	sb.WriteString(fmt.Sprintf(" chain name: %s\n", s.ChainInfo.NodeInfo.Network))
	sb.WriteString(fmt.Sprintf(" tendermint version: v%s\n", s.ChainInfo.NodeInfo.Version))

	if s.BlockTime != 0 {
		sb.WriteString(fmt.Sprintf(" avg block time: %s\n", utils.SerializeDuration(s.BlockTime)))
	}

	if s.Upgrade == nil {
		sb.WriteString(" no chain upgrade scheduled\n")
	} else {
		sb.WriteString(fmt.Sprintf(
			" chain upgrade %s scheduled at block %d\n",
			s.Upgrade.Name,
			s.Upgrade.Height,
		))

		if s.BlockTime != 0 {
			upgradeTime := utils.CalculateTimeTillBlock(s.Height, s.Upgrade.Height, s.BlockTime)
			sb.WriteString(fmt.Sprintf(" upgrade estimated time: %s\n", utils.SerializeTime(upgradeTime)))
			sb.WriteString(fmt.Sprintf(
				" time till upgrade: %s\n",
				utils.SerializeDuration(time.Until(upgradeTime)),
			))
		}
	}

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
