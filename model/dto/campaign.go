package dto

import (
	"github.com/Isabella714/gigmint/model/bo"
)

type Campaign struct {
	ID           string           `json:"id"`
	Tune         *Tune            `json:"tune"`
	MaximumScore uint64           `json:"maximum_score"`
	MinimumScore uint64           `json:"minimum_score"`
	Level        uint16           `json:"level"`
	Fee          string           `json:"fee"`
	RewardPool   string           `json:"reward_pool"`
	Stage        bo.CampaignStage `json:"stage"`
	Deadline     *int64           `json:"deadline"`
	CreatedAt    int64            `json:"created_at"`
}

type PagingCampaignRequest struct {
	*PagingRequest
	Level *uint16           `query:"level"`
	Stage *bo.CampaignStage `query:"stage"`
}
