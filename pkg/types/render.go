package types

import (
	"fmt"
	"main/pkg/utils"
	"strings"
	"time"
)

type RenderInfo struct {
	Height     int64
	Round      int64
	Step       int64
	Validators Validators
	StartTime  time.Time
}

func (r RenderInfo) SerializeInfo() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("height=%d round=%d step=%d\n", r.Height, r.Round, r.Step))
	sb.WriteString(fmt.Sprintf("block time: %s\n", utils.ZeroOrPositiveDuration(time.Since(r.StartTime))))
	sb.WriteString(fmt.Sprintf(
		"prevote consensus (total/agreeing): %.2f / %.2f\n",
		r.Validators.GetTotalVotingPowerPrevotedPercent(true),
		r.Validators.GetTotalVotingPowerPrevotedPercent(false),
	))
	sb.WriteString(fmt.Sprintf(
		"precommit consensus (total/agreeing): %.2f / %.2f\n",
		r.Validators.GetTotalVotingPowerPrecommittedPercent(true),
		r.Validators.GetTotalVotingPowerPrecommittedPercent(false),
	))

	return sb.String()
}
