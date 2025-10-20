package queue

import (
	"context"
)

type EventEmitter interface {
	Emit(ctx context.Context, event Event) error
	Close() error
}
