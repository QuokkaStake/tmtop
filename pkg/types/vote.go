package types

import loggerPkg "main/pkg/logger"

type Vote int

const (
	Voted Vote = iota
	VotedNil
	VotedZero
)

func (v Vote) Serialize() string {
	switch v {
	case Voted:
		return "✅"
	case VotedZero:
		return "🤷"
	case VotedNil:
		return "❌"
	}

	loggerPkg.GetDefaultLogger().Fatal().Str("value", string(rune(v))).Msg("Error parsing vote")
	return ""
}
