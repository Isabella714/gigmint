package dao

import (
	"context"

	"gorm.io/gorm/clause"

	"github.com/Isabella714/gigmint/component/mysql"
	"github.com/Isabella714/gigmint/model/entity"
)

type TuneDAO struct {
}

func NewTuneDAO() *TuneDAO {
	return &TuneDAO{}
}

func (*TuneDAO) CreateTune(ctx context.Context, tune *entity.TuneEntity) error {
	return mysql.Get(ctx).
		WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).
		Create(tune).Error
}

func (*TuneDAO) MGetTune(ctx context.Context, ids []uint64) (tunes []*entity.TuneEntity, err error) {
	err = mysql.Get(ctx).WithContext(ctx).
		Where("id IN ?", ids).
		Find(&tunes).Error

	return
}
