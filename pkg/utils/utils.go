package utils

import (
	"bytes"
	loggerPkg "main/pkg/logger"
	"math"
	"strconv"
	"time"

	"github.com/btcsuite/btcutil/bech32"
	"github.com/cometbft/cometbft/crypto"
	"github.com/cometbft/cometbft/crypto/ed25519"
	comet_secp256k1 "github.com/cometbft/cometbft/crypto/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func PubKeyToPeerID(pubKeyBytes string) (string, error) {
	// Convert to ed25519 public key
	pubKey := make(ed25519.PubKey, ed25519.PubKeySize)
	copy(pubKey[:], pubKeyBytes)

	// Get the peer ID from the public key
	peerID := crypto.Address(pubKey).String()
	return peerID, nil
}

func ValidatorAddr(pubkeyBytes []byte) string {
	pubkey := comet_secp256k1.PubKey(pubkeyBytes)
	pubKeyConvertedToAddress := sdk.ValAddress(pubkey.Address().Bytes())
	return pubKeyConvertedToAddress.String()
}

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

func CalculateTimeTillBlock(currentHeight, requiredHeight int64, blockTime time.Duration) time.Time {
	blocksTillEstimatedBlock := requiredHeight - currentHeight
	secondsTillEstimatedBlock := int64(float64(blocksTillEstimatedBlock) * blockTime.Seconds())
	durationTillEstimatedBlock := time.Duration(secondsTillEstimatedBlock * int64(time.Second))

	return time.Now().Add(durationTillEstimatedBlock)
}

func SerializeTime(date time.Time) string {
	return date.Format(time.RFC850)
}

func SerializeDuration(duration time.Duration) time.Duration {
	digits := 3
	denom := time.Duration(math.Pow(float64(10), float64(digits)))

	if duration > time.Second {
		return duration.Round(time.Second / denom)
	}

	return duration.Round(time.Millisecond / denom)
}

func Find[T any](slice []T, f func(T) bool) (T, bool) {
	for _, elt := range slice {
		if f(elt) {
			return elt, true
		}
	}

	return *new(T), false
}

func CompareTwoBech32(first, second string) (bool, error) {
	_, firstBytes, err := bech32.Decode(first)
	if err != nil {
		return false, err
	}

	_, secondBytes, err := bech32.Decode(second)
	if err != nil {
		return false, err
	}

	return bytes.Equal(firstBytes, secondBytes), nil
}
