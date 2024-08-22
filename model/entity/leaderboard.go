package entity

const TableLeaderboard = "leaderboard"

type LeaderboardEntity struct {
	ID            uint64 `gorm:"column:id"`
	WalletAddress string `gorm:"column:wallet_address"`
	Score         uint64 `gorm:"column:score"`
	Tune          uint64 `gorm:"column:tune"`
	CreatedAt     string `gorm:"column:created_at"`
	UpdatedAt     string `gorm:"column:updated_at"`
}

func (*LeaderboardEntity) TableName() string {
	return TableLeaderboard
}
