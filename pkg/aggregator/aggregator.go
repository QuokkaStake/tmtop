package aggregator

import (
	configPkg "main/pkg/config"
	dataFetcher "main/pkg/fetcher"
	"main/pkg/tendermint"
	"main/pkg/types"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

type Aggregator struct {
	Config *configPkg.Config
	Logger zerolog.Logger

	TendermintClient *tendermint.RPC
	DataFetcher      dataFetcher.DataFetcher
}

func NewAggregator(config *configPkg.Config, state *types.State, logger zerolog.Logger) *Aggregator {
	return &Aggregator{
		Config:           config,
		Logger:           logger.With().Str("component", "aggregator").Logger(),
		TendermintClient: tendermint.NewRPC(config, state, logger),
		DataFetcher:      dataFetcher.GetDataFetcher(config, state, logger),
	}
}

func (a *Aggregator) GetData() (
	*types.ConsensusStateResponse,
	[]types.TendermintValidator,
	error,
) {
	var wg sync.WaitGroup

	var consensusError error
	var validatorsError error

	var validators []types.TendermintValidator
	var consensus *types.ConsensusStateResponse

	wg.Add(2)

	go func() {
		consensus, consensusError = a.TendermintClient.GetConsensusState()
		wg.Done()
	}()

	go func() {
		validators, validatorsError = a.TendermintClient.GetValidators()
		wg.Done()
	}()

	wg.Wait()

	if consensusError != nil {
		a.Logger.Error().Err(consensusError).Msg("Could not fetch consensus data")
		return nil, nil, consensusError
	}

	if validatorsError != nil {
		a.Logger.Error().Err(validatorsError).Msg("Could not fetch validators")
		return nil, nil, validatorsError
	}

	return consensus, validators, nil
}

func (a *Aggregator) GetChainValidators() (*types.ChainValidators, error) {
	return a.DataFetcher.GetValidators()
}

func (a *Aggregator) GetChainInfo(rpcURL string) (*types.TendermintStatusResponse, error) {
	return a.TendermintClient.GetStatus(rpcURL)
}

func (a *Aggregator) GetUpgrade() (*types.Upgrade, error) {
	return a.DataFetcher.GetUpgradePlan()
}

func (a *Aggregator) GetBlockTime() (time.Duration, error) {
	return a.TendermintClient.GetBlockTime()
}

func (a *Aggregator) GetNetInfo(host string) (*types.NetInfo, error) {
	return a.DataFetcher.GetNetInfo(host)
}
