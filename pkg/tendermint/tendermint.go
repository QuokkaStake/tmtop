package tendermint

import (
	"encoding/json"
	"fmt"
	configPkg "main/pkg/config"
	"net/http"
	"strconv"
	"time"

	"main/pkg/types"

	"github.com/rs/zerolog"
)

type RPC struct {
	Config     configPkg.Config
	Logger     zerolog.Logger
	LogChannel chan string
}

func NewRPC(config configPkg.Config, logger zerolog.Logger) *RPC {
	return &RPC{
		Config: config,
		Logger: logger.With().Str("component", "rpc").Logger(),
	}
}

func (rpc *RPC) GetConsensusState() (*types.ConsensusStateResponse, error) {
	var response types.ConsensusStateResponse
	if err := rpc.Get("/consensus_state", &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (rpc *RPC) GetValidators() ([]types.TendermintValidator, error) {
	page := 1

	validators := make([]types.TendermintValidator, 0)

	for {
		response, err := rpc.GetValidatorsAtPage(page)
		if err != nil {
			return nil, err
		}

		total, err := strconv.ParseInt(response.Result.Total, 10, 64)
		if err != nil {
			return nil, err
		}

		validators = append(validators, response.Result.Validators...)
		if int64(len(validators)) >= total {
			break
		}

		page++
	}

	return validators, nil
}

func (rpc *RPC) GetStatus() (*types.TendermintStatusResponse, error) {
	var response types.TendermintStatusResponse
	if err := rpc.Get("/status", &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (rpc *RPC) GetValidatorsAtPage(page int) (*types.ValidatorsResponse, error) {
	var response types.ValidatorsResponse
	if err := rpc.Get(fmt.Sprintf("/validators?page=%d&per_page=100", page), &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (t *RPC) Block(height int64) (types.TendermintBlockResponse, error) {
	blockUrl := fmt.Sprintf("/block")
	if height != 0 {
		blockUrl = fmt.Sprintf("/block?height=%d", height)
	}

	res := types.TendermintBlockResponse{}
	err := t.Get(blockUrl, &res)
	return res, err
}

func (t *RPC) GetBlockTime() (time.Duration, error) {
	var blocksBehind int64 = 1000

	latestBlock, err := t.Block(0)
	if err != nil {
		t.Logger.Error().Err(err).Msg("Could not fetch current block")
		return 0, err
	}

	latestBlockHeight, err := strconv.ParseInt(latestBlock.Result.Block.Header.Height, 10, 64)
	if err != nil {
		t.Logger.Error().
			Err(err).
			Msg("Error converting latest block height to int64, which should never happen.")
		return 0, err
	}
	blockToCheck := latestBlockHeight - blocksBehind

	olderBlock, err := t.Block(blockToCheck)
	if err != nil {
		t.Logger.Error().Err(err).Msg("Could not fetch older block")
		return 0, err
	}

	blocksDiffTime := latestBlock.Result.Block.Header.Time.Sub(olderBlock.Result.Block.Header.Time)
	blockTime := blocksDiffTime.Seconds() / float64(blocksBehind)

	duration := time.Duration(int64(blockTime * float64(time.Second)))
	return duration, nil
}

func (rpc *RPC) Get(relativeURL string, target interface{}) error {
	client := &http.Client{Timeout: 300 * time.Second}
	start := time.Now()

	url := fmt.Sprintf("%s%s", rpc.Config.RPCHost, relativeURL)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", "tmtop")

	rpc.Logger.Debug().Str("url", url).Msg("Doing a query...")

	res, err := client.Do(req)
	if err != nil {
		rpc.Logger.Warn().Str("url", url).Err(err).Msg("Query failed")
		return err
	}
	defer res.Body.Close()

	rpc.Logger.Debug().Str("url", url).Dur("duration", time.Since(start)).Msg("Query is finished")

	return json.NewDecoder(res.Body).Decode(target)
}
