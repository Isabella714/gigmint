package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

const TableCampaign = "campaign"

type CampaignEntity struct {
	ID           uint64          `json:"id"`
	Tune         uint64          `json:"tune"`
	MaximumScore uint64          `gorm:"column:maximum_score"`
	MinimumScore uint64          `gorm:"column:minimum_score"`
	Level        uint16          `gorm:"column:level"`
	Fee          decimal.Decimal `gorm:"column:fee"`
	RewardPool   decimal.Decimal `gorm:"column:reward_pool"`
	Stage        uint8           `gorm:"column:stage"`
	Deadline     *time.Time      `gorm:"column:deadline"`
	CreatedAt    time.Time       `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time       `gorm:"column:updated_at;autoUpdateTime"`
}

func (*CampaignEntity) TableName() string {
	return TableCampaign
}
