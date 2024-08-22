package dao

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/Isabella714/gigmint/component/mysql"
	"github.com/Isabella714/gigmint/model/entity"
)

type LeaderboardDAO struct {
}

func NewLeaderboard() *LeaderboardDAO {
	return &LeaderboardDAO{}
}

func (*LeaderboardDAO) CreateLeaderboard(ctx context.Context, record *entity.LeaderboardEntity) error {
	return mysql.Get(ctx).WithContext(ctx).
		Clauses(clause.OnConflict{
			DoUpdates: []clause.Assignment{
				{
					Column: clause.Column{
						Name: "score",
					},
					Value: gorm.Expr("GREATEST(score, ?)", record.Score),
				},
			},
		}).
		Create(record).Error
}
