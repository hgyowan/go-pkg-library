package queue

import (
	"context"
	"encoding/json"
	pkgError "github.com/hgyowan/go-pkg-library/error"
	pkgLogger "github.com/hgyowan/go-pkg-library/logger"
	"github.com/segmentio/kafka-go"
)

type kafkaEmitter struct {
	writer *kafka.Writer
}

func (k *kafkaEmitter) Close() error {
	return pkgError.Wrap(k.writer.Close())
}

func (k *kafkaEmitter) Emit(ctx context.Context, event Event) error {
	data, err := json.Marshal(event)
	if err != nil {
		return pkgError.Wrap(err)
	}

	err = k.writer.WriteMessages(ctx,
		kafka.Message{
			Key:   []byte(event.MessageDeduplicationId),
			Value: data,
		},
	)
	if err != nil {
		return pkgError.Wrap(err)
	}

	return nil
}

func MustNewKafkaEventEmitter(writer *kafka.Writer) EventEmitter {
	if writer == nil {
		pkgLogger.ZapLogger.Logger.Fatal("Failed to get kafka writer")
		return nil
	}

	return &kafkaEmitter{
		writer: writer,
	}
}
