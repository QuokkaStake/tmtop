package types

import (
	"fmt"
	"strings"
)

type ProgressBar struct {
	Width    int
	Height   int
	Progress int
	Prefix   string
}

func (p ProgressBar) Serialize() string {
	var sb strings.Builder

	percentText := fmt.Sprintf("%s %d%%", p.Prefix, p.Progress)
	percentTextStart := (p.Width - len(percentText)) / 2
	percentTextLine := (p.Height+1)/2 - 1

	for lineIndex := 0; lineIndex < p.Height; lineIndex++ {
		line := "[\"progress\"]"
		terminated := false

		for index := 0; index < p.Width; index++ {
			percent := int(float64(index) / float64(p.Width) * 100)

			if percent > p.Progress && !terminated {
				terminated = true
				line += "[\"\"]"
			}

			if lineIndex == percentTextLine && (index >= percentTextStart && index < percentTextStart+len(percentText)) {
				line += string(percentText[index-percentTextStart])
			} else {
				line += " "
			}
		}

		if !terminated {
			line += "[\"\"]"
		}

		sb.WriteString(line + "\n")
	}

	return sb.String()
}
