package queue

import (
	"context"
	"encoding/json"
	"errors"
	pkgError "github.com/hgyowan/go-pkg-library/error"
	pkgLogger "github.com/hgyowan/go-pkg-library/logger"
	pkgVariable "github.com/hgyowan/go-pkg-library/variable"
	"github.com/segmentio/kafka-go"
	"strings"
	"sync"
)

type KafkaListenerConfig struct {
	ConsumerCount int `json:"consumerCount"`
}

type kafkaListener struct {
	reader *kafka.Reader
	config *KafkaListenerConfig
}

func (k *kafkaListener) Listen(ctx context.Context, events ...string) (<-chan []Event, <-chan error, error) {
	if k == nil {
		return nil, nil, pkgError.Wrap(errors.New("kafka listener is nil"))
	}

	eventCh := make(chan []Event)
	errorCh := make(chan error)

	wg := &sync.WaitGroup{}

	for i := 0; i < k.config.ConsumerCount; i++ {
		wg.Add(1)
		go func(consumerID int) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				default:
					k.ReceiveMessage(ctx, eventCh, errorCh, events, consumerID)
				}
			}
		}(i)
	}

	go func() {
		<-ctx.Done()
		wg.Wait()
		close(eventCh)
		close(errorCh)
		k.reader.Close()
	}()

	return eventCh, errorCh, nil
}

func (k *kafkaListener) ReceiveMessage(ctx context.Context, eventCh chan []Event, errorCh chan error, events []string, consumerID int) {
	msg, err := k.reader.FetchMessage(ctx)
	if err != nil {
		if ctx.Err() != nil {
			return
		}
		select {
		case errorCh <- pkgError.Wrap(err):
		case <-ctx.Done():
			return
		}
	}

	var event Event
	err = json.Unmarshal(msg.Value, &event)
	if err != nil {
		if deleteMessageErr := k.reader.CommitMessages(ctx, msg); deleteMessageErr != nil {
			pkgLogger.ZapLogger.Logger.Sugar().Error(pkgError.Wrap(deleteMessageErr))
		}
		errorCh <- pkgError.Wrap(err)
		return
	}

	bContinue := false
	for _, e := range events {
		if strings.EqualFold(event.EventName, e) {
			bContinue = true
			break
		}
	}

	b, _ := json.Marshal(msg)
	event.ReceiptHandle = pkgVariable.ConvertToPointer(string(b))

	if !bContinue {
		if deleteMessageErr := k.reader.CommitMessages(ctx, msg); deleteMessageErr != nil {
			pkgLogger.ZapLogger.Logger.Sugar().Error(pkgError.Wrap(deleteMessageErr))
		}
		return
	}

	eventCh <- []Event{event}
}

func (k *kafkaListener) DeleteMessage(ctx context.Context, ReceiptHandle *string) error {
	var msg kafka.Message
	if err := json.Unmarshal([]byte(pkgVariable.GetSafeValue(ReceiptHandle, "")), &msg); err != nil {
		return pkgError.WrapWithCode(err, pkgError.Delete)
	}

	if err := k.reader.CommitMessages(ctx, msg); err != nil {
		return pkgError.WrapWithCode(err, pkgError.Delete)
	}

	return nil
}

func (k *kafkaListener) ChangeMessageVisibility(ctx context.Context, ReceiptHandle *string, VisibilityTimeout int32) error {
	return nil
}

func MustNewKafkaEventListener(reader *kafka.Reader, config *KafkaListenerConfig) EventListener {
	if config == nil {
		pkgLogger.ZapLogger.Logger.Fatal("Failed to get kafka listener config")
		return nil
	}

	if reader == nil {
		pkgLogger.ZapLogger.Logger.Fatal("Failed to get kafka reader")
		return nil
	}

	return &kafkaListener{
		reader: reader,
		config: config,
	}
}
