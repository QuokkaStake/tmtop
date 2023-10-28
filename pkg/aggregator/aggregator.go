package aggregator

import (
	configPkg "main/pkg/config"
	dataFetcher "main/pkg/fetcher"
	"main/pkg/tendermint"
	"main/pkg/types"
	"sync"

	"github.com/rs/zerolog"
)

type Aggregator struct {
	Config configPkg.Config
	Logger zerolog.Logger

	TendermintClient *tendermint.RPC
	DataFetcher      dataFetcher.DataFetcher
}

func NewAggregator(config configPkg.Config, logger zerolog.Logger) *Aggregator {
	return &Aggregator{
		Config:           config,
		Logger:           logger.With().Str("component", "aggregator").Logger(),
		TendermintClient: tendermint.NewRPC(config, logger),
		DataFetcher:      dataFetcher.GetDataFetcher(config, logger),
	}
}

func (a *Aggregator) GetData() (
	*types.ConsensusStateResponse,
	[]types.TendermintValidator,
	types.ChainValidators,
	error,
) {
	var wg sync.WaitGroup

	var consensusError error
	var validatorsError error
	var chainValidatorsError error

	var validators []types.TendermintValidator
	var consensus *types.ConsensusStateResponse
	var chainValidators types.ChainValidators

	wg.Add(3)

	go func() {
		consensus, consensusError = a.TendermintClient.GetConsensusState()
		wg.Done()
	}()

	go func() {
		validators, validatorsError = a.TendermintClient.GetValidators()
		wg.Done()
	}()

	go func() {
		chainValidators, chainValidatorsError = a.DataFetcher.GetValidators()
		wg.Done()
	}()

	wg.Wait()

	if consensusError != nil {
		return nil, nil, nil, consensusError
	}

	if validatorsError != nil {
		return nil, nil, nil, validatorsError
	}

	if chainValidatorsError != nil {
		return nil, nil, nil, validatorsError
	}

	return consensus, validators, chainValidators, nil
}
