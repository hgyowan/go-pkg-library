package queue

import (
	"context"
)

type EventListener interface {
	Listen(ctx context.Context, events ...string) (<-chan []Event, <-chan error, error)
	ReceiveMessage(ctx context.Context, eventCh chan []Event, errorCh chan error, events []string, consumerID int)
	DeleteMessage(ctx context.Context, ReceiptHandle *string) error
	ChangeMessageVisibility(ctx context.Context, ReceiptHandle *string, VisibilityTimeout int32) error
}
