package utils

import (
	loggerPkg "main/pkg/logger"
	"strconv"
	"time"
)

func MustParseInt64(source string) int64 {
	result, err := strconv.ParseInt(source, 10, 64)
	if err != nil {
		loggerPkg.GetDefaultLogger().Fatal().Str("value", source).Msg("Could not parse int64")
	}

	return result
}

func ZeroOrPositiveDuration(duration time.Duration) time.Duration {
	if duration < 0 {
		return 0
	}

	return duration
}

func PadAndTrim(source string, desiredLength int, padLeft bool) string {
	if len(source) < desiredLength {
		result := source

		for len(result) < desiredLength {
			if padLeft {
				result += " "
			} else {
				result = " " + result
			}
		}

		return result
	}

	if len(source) > desiredLength {
		return source[:desiredLength-3] + "..."
	}

	return source
}

func RightPadAndTrim(source string, desiredLength int) string {
	return PadAndTrim(source, desiredLength, false)
}

func LeftPadAndTrim(source string, desiredLength int) string {
	return PadAndTrim(source, desiredLength, true)
}
