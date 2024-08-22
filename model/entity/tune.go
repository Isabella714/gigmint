package entity

import "time"

const TableTune = "tune"

type TuneEntity struct {
	ID         uint64    `gorm:"column:id"`
	Name       string    `gorm:"column:name"`
	RhythmFile string    `gorm:"column:rhythm_file"`
	Owner      string    `gorm:"column:owner"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt  time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (*TuneEntity) TableName() string {
	return TableTune
}
