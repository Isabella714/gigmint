package processor

import (
	"context"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/Isabella714/gigmint/pkg/contracts"
	"github.com/Isabella714/gigmint/service"
)

type GiGCampaignEnrollTuneProcessor struct {
	*AbstractEventProcessor[*contracts.GiGCampaign]
	campaignProcessor *service.CampaignService
}

func NewGiGCampaignEnrollTuneProcessor(contract *contracts.GiGCampaign) *GiGCampaignEnrollTuneProcessor {
	return &GiGCampaignEnrollTuneProcessor{
		AbstractEventProcessor: NewAbstractEventProcessor[*contracts.GiGCampaign]("EnrollTune(uint256,uint128,uint128,uint256,uint256)", contract),
		campaignProcessor:      service.NewCampaignService(),
	}
}

func (p *GiGCampaignEnrollTuneProcessor) Process(ctx context.Context, log *types.Log) error {
	event, err := p.contract.ParseEnrollTune(*log)
	if err != nil {
		return err
	}

	return p.campaignProcessor.SyncCampaign(ctx, event)
}
