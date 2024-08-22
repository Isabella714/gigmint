package service

import (
	"context"

	"github.com/Isabella714/gigmint/dao"
	"github.com/Isabella714/gigmint/model/bo"
	"github.com/Isabella714/gigmint/model/entity"
	"github.com/Isabella714/gigmint/pkg/contracts"
)

type TuneService struct {
	tuneDAO *dao.TuneDAO
}

func NewTuneService() *TuneService {
	return &TuneService{
		tuneDAO: dao.NewTuneDAO(),
	}
}

func (s *TuneService) SyncTune(ctx context.Context, event *contracts.GiGTuneMintTune) (err error) {
	err = s.tuneDAO.CreateTune(ctx, &entity.TuneEntity{
		ID:    event.TokenId.Uint64(),
		Owner: event.Owner.Hex(),
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *TuneService) MGetTune(ctx context.Context, tuneIds []uint64) (tunes map[uint64]*bo.Tune, err error) {
	entities, err := s.tuneDAO.MGetTune(ctx, tuneIds)
	if err != nil {
		return nil, err
	}

	tunes = make(map[uint64]*bo.Tune)
	for _, item := range entities {
		tunes[item.ID] = &bo.Tune{
			ID:         item.ID,
			Name:       item.Name,
			RhythmFile: item.RhythmFile,
			Owner:      item.Owner,
			CreatedAt:  item.CreatedAt,
		}
	}

	return tunes, nil
}
