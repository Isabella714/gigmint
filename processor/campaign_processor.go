package processor

import (
	"context"
	"strconv"

	"github.com/Isabella714/gigmint/model/bo"
	"github.com/Isabella714/gigmint/model/dto"
	"github.com/Isabella714/gigmint/service"
)

type CampaignProcessor struct {
	campaignService *service.CampaignService
	tuneService     *service.TuneService
}

func NewCampaignProcessor() *CampaignProcessor {
	return &CampaignProcessor{
		campaignService: service.NewCampaignService(),
		tuneService:     service.NewTuneService(),
	}
}

func (p *CampaignProcessor) PagingCampaign(ctx context.Context, request *dto.PagingCampaignRequest) (response []*dto.Campaign, err error) {
	campaigns, err := p.campaignService.PagingCampaign(ctx, &bo.PagingCampaignParam{
		Page:  request.Page,
		Size:  request.Size,
		Stage: request.Stage,
		Level: request.Level,
	})
	if err != nil {
		return nil, err
	}

	if len(campaigns) == 0 {
		return make([]*dto.Campaign, 0), nil
	}

	var tuneIds = make([]uint64, 0)
	for _, item := range campaigns {
		tuneIds = append(tuneIds, item.Tune)
	}

	tunes, err := p.tuneService.MGetTune(ctx, tuneIds)
	if err != nil {
		return nil, err
	}

	for _, item := range campaigns {
		var deadline *int64
		if item.Deadline != nil {
			var val = item.Deadline.Unix()
			deadline = &val
		}
		var tune *dto.Tune
		if tuneItem, ok := tunes[item.Tune]; ok {
			tune = &dto.Tune{
				ID:         strconv.FormatUint(tuneItem.ID, 10),
				Name:       tuneItem.Name,
				RhythmFile: tuneItem.RhythmFile,
				Owner:      tuneItem.Owner,
				CreatedAt:  tuneItem.CreatedAt.Unix(),
			}
		}
		response = append(response, &dto.Campaign{
			ID:           strconv.FormatUint(item.ID, 10),
			Tune:         tune,
			MaximumScore: item.MaximumScore,
			MinimumScore: item.MinimumScore,
			Level:        item.Level,
			Fee:          item.Fee.String(),
			RewardPool:   item.RewardPool.String(),
			Stage:        item.Stage,
			Deadline:     deadline,
			CreatedAt:    item.CreatedAt.Unix(),
		})
	}

	return response, nil
}
