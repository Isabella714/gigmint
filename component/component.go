package component

import "context"

type IComponent[T any] interface {
	Instantiate() error
	Close() error
	Get(ctx context.Context) T
}
