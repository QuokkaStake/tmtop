package types

import "strings"

type ProgressBar struct {
	Width    int
	Height   int
	Progress int
}

func (p ProgressBar) Serialize() string {
	line := ""

	for index := 0; index < p.Width; index++ {
		percent := int(float64(index) / float64(p.Width) * 100)
		if percent > p.Progress {
			line += " "
		} else {
			line += "█"
		}
	}

	var sb strings.Builder
	for index := 0; index < p.Height; index++ {
		sb.WriteString(line + "\n")
	}

	return sb.String()
}
