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
	Config configPkg.Config
	Logger zerolog.Logger
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

func (rpc *RPC) GetValidatorsAtPage(page int) (*types.ValidatorsResponse, error) {
	var response types.ValidatorsResponse
	if err := rpc.Get(fmt.Sprintf("/validators?page=%d&per_page=100", page), &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (rpc *RPC) Get(relativeUrl string, target interface{}) error {
	client := &http.Client{Timeout: 300 * time.Second}
	start := time.Now()

	url := fmt.Sprintf("%s%s", rpc.Config.RPCHost, relativeUrl)

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
