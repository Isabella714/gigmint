package scanner

import (
	"context"
	"math/big"

	geth "github.com/ethereum/go-ethereum"
	gethCommon "github.com/ethereum/go-ethereum/common"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/gammazero/workerpool"
	"github.com/go-co-op/gocron"
	"github.com/gofiber/fiber/v3/log"
	"github.com/sirupsen/logrus"

	"github.com/Isabella714/gigmint/component"
	"github.com/Isabella714/gigmint/component/blockchain"
	"github.com/Isabella714/gigmint/pkg/contracts"
	"github.com/Isabella714/gigmint/scanner/processor"
)

type BlockEventScanner struct {
	scanners    []IEventScanner
	blockHeight uint64
}

func RegisterBlockEventScanner(scheduler *gocron.Scheduler) (*gocron.Job, error) {
	var p = &BlockEventScanner{
		scanners: []IEventScanner{
			NewGiGTuneEventScanner(),
			NewGiGCampaignEventScanner(),
		},
		blockHeight: 43206459,
	}

	return scheduler.Tag(p.Name()).CronWithSeconds(p.Cron()).Do(p.Run)
}

func (p *BlockEventScanner) Name() string {
	return "BlockEventScanner"
}

func (p *BlockEventScanner) Cron() string {
	return "*/5 * * * * *"
}

func (p *BlockEventScanner) Run() {
	var ctx = context.Background()

	log.Infow("BlockEventScanner start", "max_block", p.blockHeight)

	workers := workerpool.New(len(p.scanners))

	for _, item := range p.scanners {
		var scanner = item
		workers.Submit(func() {
			err := scanner.Scan(ctx, p.blockHeight)
			if err != nil {
				logrus.WithField("err", err).Error("scanner has error")
			}
		})
	}

	workers.StopWait()

	log.Infow("BlockEventScanner end", "max_block", p.blockHeight)

	p.blockHeight++
}

type IEventScanner interface {
	Scan(ctx context.Context, currentHeight uint64) error
}

type EventScanner struct {
	name           string
	blockHeightKey string
	controlKey     string
	address        gethCommon.Address
	topics         []gethCommon.Hash
	processors     []processor.IEventProcessor
}

func NewEventScanner(name string, address gethCommon.Address, topics []gethCommon.Hash, processors []processor.IEventProcessor) *EventScanner {
	return &EventScanner{
		name:       name,
		address:    address,
		topics:     topics,
		processors: processors,
	}
}

func (scanner *EventScanner) Scan(ctx context.Context, height uint64) error {
	log.Infow("scan block start", "block_height", height)
	defer log.Infow("scan block end", "block_height", height)

	events, err := scanner.scanEvents(ctx, height)
	if err != nil {
		log.Errorw("scan block has error", "block_height", height, "error", err)
		return err
	}

	err = scanner.processEvents(ctx, events)
	if err != nil {
		log.Errorw("process block events has error", "block_height", height, "error", err)
		return err
	}

	return nil
}

func (scanner *EventScanner) scanEvents(ctx context.Context, height uint64) ([]gethTypes.Log, error) {
	filter := geth.FilterQuery{
		FromBlock: new(big.Int).SetUint64(height),
		ToBlock:   new(big.Int).SetUint64(height),
		Addresses: []gethCommon.Address{scanner.address},
		Topics:    [][]gethCommon.Hash{scanner.topics},
	}

	return blockchain.Get(ctx).FilterLogs(ctx, filter)
}

func (scanner *EventScanner) processEvents(ctx context.Context, events []gethTypes.Log) error {
	for _, event := range events {
		for _, eventProcessor := range scanner.processors {
			if eventProcessor.Support(ctx, &event) {
				err := eventProcessor.Process(ctx, &event)
				if err != nil {
					return err
				}
				break
			}
		}
	}

	return nil
}

type GiGTuneEventScanner struct {
	*EventScanner
	contract *contracts.GiGTune
}

func NewGiGTuneEventScanner() *GiGTuneEventScanner {
	var rawAddress, _ = component.GetConfigString("contract.gigtune")
	var address = gethCommon.HexToAddress(rawAddress)
	var contract, _ = contracts.NewGiGTune(address, blockchain.Get(context.Background()))

	var processors = []processor.IEventProcessor{
		processor.NewGiGTuneMintTuneEventProcessor(contract),
	}
	var topics = make([]gethCommon.Hash, 0)
	for _, item := range processors {
		topics = append(topics, item.EventHash(context.Background()))
	}

	return &GiGTuneEventScanner{
		EventScanner: NewEventScanner("gig_tune", address, topics, processors),
		contract:     contract,
	}
}

type GiGCampaignEventScanner struct {
	*EventScanner
	contract *contracts.GiGCampaign
}

func NewGiGCampaignEventScanner() *GiGCampaignEventScanner {
	var rawAddress, _ = component.GetConfigString("contract.gigcampaign")
	var address = gethCommon.HexToAddress(rawAddress)
	var contract, _ = contracts.NewGiGCampaign(address, blockchain.Get(context.Background()))

	var processors = []processor.IEventProcessor{
		processor.NewGiGCampaignEnrollTuneProcessor(contract),
	}
	var topics = make([]gethCommon.Hash, 0)
	for _, item := range processors {
		topics = append(topics, item.EventHash(context.Background()))
	}

	return &GiGCampaignEventScanner{
		EventScanner: NewEventScanner("gig_campaign", address, topics, processors),
		contract:     contract,
	}
}
