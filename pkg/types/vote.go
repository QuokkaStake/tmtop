package types

type Vote int

const (
	Voted Vote = iota
	VotedNil
	VotedZero
)

func (v Vote) Serialize(disableEmojis bool) string {
	if disableEmojis {
		switch v {
		case Voted:
			return "[X[]"
		case VotedZero:
			return "[0[]"
		case VotedNil:
			return "[ []"
		default:
			return ""
		}
	}

	switch v {
	case Voted:
		return "✅"
	case VotedZero:
		return "🤷"
	case VotedNil:
		return "❌"
	default:
		return ""
	}
}
