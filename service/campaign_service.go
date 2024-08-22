package service

import (
	"context"
	"time"

	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/gofiber/fiber/v3/log"
	"github.com/shopspring/decimal"

	"github.com/Isabella714/gigmint/component"
	"github.com/Isabella714/gigmint/component/blockchain"
	"github.com/Isabella714/gigmint/dao"
	"github.com/Isabella714/gigmint/model/bo"
	"github.com/Isabella714/gigmint/model/entity"
	"github.com/Isabella714/gigmint/pkg/contracts"
)

type CampaignService struct {
	campaignContract *contracts.GiGCampaign
	campaignDAO      *dao.CampaignDAO
	leaderboardDAO   *dao.LeaderboardDAO
}

func NewCampaignService() *CampaignService {
	gigcapaignAddress, _ := component.GetConfigString("contract.gigcampaign")
	campaignContract, err := contracts.NewGiGCampaign(gethCommon.HexToAddress(gigcapaignAddress), blockchain.Get(context.Background()))
	if err != nil {
		panic(err)
	}

	return &CampaignService{
		campaignContract: campaignContract,
		campaignDAO:      dao.NewCampaignDAO(),
		leaderboardDAO:   dao.NewLeaderboard(),
	}
}

func (s *CampaignService) SyncCampaign(ctx context.Context, event *contracts.GiGCampaignEnrollTune) (err error) {
	deadline := time.Unix(event.CompetitionDeadline.Int64(), 0)
	err = s.campaignDAO.CreateCampaign(ctx, &entity.CampaignEntity{
		ID:           event.TokenId.Uint64(),
		Tune:         event.TokenId.Uint64(),
		MaximumScore: event.MaximumScore.Uint64(),
		MinimumScore: event.MinimumScore.Uint64(),
		Deadline:     &deadline,
		Fee:          decimal.NewFromBigInt(event.CompetitionFee, -18),
		RewardPool:   decimal.Zero,
		Stage:        uint8(bo.CampaignStageAudition),
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *CampaignService) PagingCampaign(ctx context.Context, param *bo.PagingCampaignParam) (campaigns []*bo.Campaign, err error) {
	var stage *uint8
	if param.Stage != nil {
		var val = uint8(*param.Stage)
		stage = &val
	}

	entities, err := s.campaignDAO.PagingCampaign(ctx, int(param.Size), int((param.Page-1)*param.Size), param.Level, stage)
	if err != nil {
		return nil, err
	}

	for _, item := range entities {
		campaigns = append(campaigns, &bo.Campaign{
			ID:           item.ID,
			Tune:         item.Tune,
			MaximumScore: item.MaximumScore,
			MinimumScore: item.MinimumScore,
			Level:        item.Level,
			Fee:          item.Fee,
			RewardPool:   item.RewardPool,
			Stage:        bo.CampaignStage(item.Stage),
			Deadline:     item.Deadline,
			CreatedAt:    item.CreatedAt,
		})
	}

	return campaigns, nil
}

func (s *CampaignService) UploadCompetitionResult(ctx context.Context, tune uint64, player string, score uint64) error {
	transaction, err := s.campaignContract.SyncCompetitionResult(blockchain.TransactOpts(ctx),
		decimal.NewFromUint64(tune).BigInt(),
		gethCommon.HexToAddress(player),
		decimal.NewFromUint64(score).BigInt())
	if err != nil {
		return err
	}

	log.WithContext(ctx).Infow("SyncCompetitionResult", "transaction", transaction)

	return nil
}

func (s *CampaignService) SyncCompetitionResult(ctx context.Context, event *contracts.GiGCampaignSyncCompetitionResult) (err error) {
	err = s.leaderboardDAO.CreateLeaderboard(ctx, &entity.LeaderboardEntity{
		WalletAddress: event.Account.Hex(),
		Score:         event.Score.Uint64(),
		Tune:          event.TokenId.Uint64(),
	})
	if err != nil {
		return err
	}

	return nil
}
