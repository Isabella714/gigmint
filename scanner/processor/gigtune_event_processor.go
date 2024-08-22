package processor

import (
	"context"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/Isabella714/gigmint/pkg/contracts"
	"github.com/Isabella714/gigmint/service"
)

type GiGTuneMintTuneEventProcessor struct {
	*AbstractEventProcessor[*contracts.GiGTune]
	tuneService *service.TuneService
}

func NewGiGTuneMintTuneEventProcessor(contract *contracts.GiGTune) *GiGTuneMintTuneEventProcessor {
	return &GiGTuneMintTuneEventProcessor{
		AbstractEventProcessor: NewAbstractEventProcessor[*contracts.GiGTune]("MintTune(address,uint256)", contract),
		tuneService:            service.NewTuneService(),
	}
}

func (p *GiGTuneMintTuneEventProcessor) Process(ctx context.Context, log *types.Log) error {
	event, err := p.contract.ParseMintTune(*log)
	if err != nil {
		return err
	}

	return p.tuneService.SyncTune(ctx, event)
}
