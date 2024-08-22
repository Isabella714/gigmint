package dao

import (
	"context"

	"gorm.io/gorm/clause"

	"github.com/Isabella714/gigmint/component/mysql"
	"github.com/Isabella714/gigmint/model/entity"
)

type CampaignDAO struct {
}

func NewCampaignDAO() *CampaignDAO {
	return &CampaignDAO{}
}

func (*CampaignDAO) CreateCampaign(ctx context.Context, campaign *entity.CampaignEntity) error {
	return mysql.Get(ctx).WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(campaign).Error
}

func (*CampaignDAO) PagingCampaign(ctx context.Context, limit, offset int, level *uint16, stage *uint8) (campaigns []*entity.CampaignEntity, err error) {
	db := mysql.Get(ctx).WithContext(ctx)
	if level != nil {
		db = db.Where("level = ?", level)
	}
	if stage != nil {
		db = db.Where("stage = ?", stage)
	}
	err = db.Limit(limit).Offset(offset).
		Order("id DESC").
		Find(&campaigns).Error

	return
}
