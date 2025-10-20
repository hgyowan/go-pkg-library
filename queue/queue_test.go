package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	pkgError "github.com/hgyowan/go-pkg-library/error"
	pkgLogger "github.com/hgyowan/go-pkg-library/logger"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_Listener(t *testing.T) {
	pkgLogger.MustInitZapLogger()
	ctx := context.Background()
	listener := MustNewKafkaEventListener(kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"192.168.20.101:30796"},
		GroupID:  "email-queue-group",
		Topic:    "email-queue",
		MinBytes: 1,
		MaxBytes: 1024 * 1024 * 1,
		MaxWait:  0,
	}), &KafkaListenerConfig{
		ConsumerCount: 1,
	})
	eventCh, errorCh, _ := listener.Listen(ctx, "email-event")
	for {
		select {
		case events := <-eventCh:
			for _, event := range events {
				switch event.EventName {
				case "email-event":
					fmt.Println(event)
					if err := listener.DeleteMessage(ctx, event.ReceiptHandle); err != nil {
						pkgLogger.ZapLogger.Logger.Sugar().Error(pkgError.Wrap(err))
					}
				}
			}
		case err := <-errorCh:
			pkgLogger.ZapLogger.Logger.Sugar().Error(pkgError.Wrap(err))
		case <-ctx.Done():
			pkgLogger.ZapLogger.Logger.Info("Listener shutting down...")
			return
		}
	}
}

func Test_Emitter(t *testing.T) {
	pkgLogger.MustInitZapLogger()
	ctx := context.Background()
	emitter := MustNewKafkaEventEmitter(&kafka.Writer{
		Addr:     kafka.TCP("192.168.20.101:30796"),
		Topic:    "email-queue",
		Balancer: kafka.CRC32Balancer{},
	})

	b, _ := json.Marshal("test")
	err := emitter.Emit(ctx, Event{
		EventName:              "email-event",
		MessageDeduplicationId: uuid.NewString(),
		Data:                   b,
	})

	require.NoError(t, err)
	defer emitter.Close()
}
