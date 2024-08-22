package bo

import (
	"time"

	"github.com/shopspring/decimal"
)

type CampaignStage uint8

const (
	CampaignStageAudition CampaignStage = iota
	CampaignStageCompete
	CampaignStageHallOfFame
)

type Campaign struct {
	ID           uint64          `json:"id"`
	Tune         uint64          `json:"tune"`
	MaximumScore uint64          `json:"maximum_score"`
	MinimumScore uint64          `json:"minimum_score"`
	Level        uint16          `json:"level"`
	Fee          decimal.Decimal `json:"fee"`
	RewardPool   decimal.Decimal `json:"reward_pool"`
	Stage        CampaignStage   `json:"stage"`
	Deadline     *time.Time      `json:"deadline"`
	CreatedAt    time.Time       `json:"created_at"`
}

type PagingCampaignParam struct {
	Page  uint32         `json:"page"`
	Size  uint32         `json:"size"`
	Level *uint16        `json:"level"`
	Stage *CampaignStage `json:"stage"`
}
