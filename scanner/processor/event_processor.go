package processor

import (
	"context"

	gethcmn "github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

type IEventProcessor interface {
	Process(ctx context.Context, log *gethtypes.Log) error
	Support(ctx context.Context, log *gethtypes.Log) bool
	EventSignature(ctx context.Context) string
	EventHash(ctx context.Context) gethcmn.Hash
}

type AbstractEventProcessor[T any] struct {
	eventSignature string
	eventHash      gethcmn.Hash
	contract       T
}

func NewAbstractEventProcessor[T any](eventSignature string, contract T) *AbstractEventProcessor[T] {
	hash := crypto.Keccak256Hash([]byte(eventSignature))

	return &AbstractEventProcessor[T]{
		eventSignature: eventSignature,
		eventHash:      hash,
		contract:       contract,
	}
}

func (p *AbstractEventProcessor[T]) Support(ctx context.Context, log *gethtypes.Log) bool {
	return log.Topics[0] == p.eventHash
}

func (p *AbstractEventProcessor[T]) EventSignature(ctx context.Context) string {
	return p.eventSignature
}

func (p *AbstractEventProcessor[T]) EventHash(ctx context.Context) gethcmn.Hash {
	return p.eventHash
}
