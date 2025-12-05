package queue

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	sqsType "github.com/aws/aws-sdk-go-v2/service/sqs/types"
	pkgError "github.com/hgyowan/go-pkg-library/error"
	pkgLogger "github.com/hgyowan/go-pkg-library/logger"
	"go.uber.org/zap"
	"strings"
	"sync"
	"time"
)

type SqsListenerConfig struct {
	MaxMsg            int32 `json:"maxMsg"`
	WaitTime          int32 `json:"waitTime"`
	VisibilityTimeout int32 `json:"visibilityTimeout"`
	SleepTime         int32 `json:"sleepTime"`
	ConsumerCount     int   `json:"consumerCount"`
}

type SqsListener struct {
	SqsClient           *sqs.Client
	QueueURL            *string
	maxNumberOfMessages int32
	waitTime            int32
	visibilityTimeOut   int32
	sleepTime           int32
	consumerCount       int
}

func (s *SqsListener) ChangeMessageVisibility(ctx context.Context, ReceiptHandle *string, VisibilityTimeout int32) error {
	_, err := s.SqsClient.ChangeMessageVisibility(ctx, &sqs.ChangeMessageVisibilityInput{
		QueueUrl:          s.QueueURL,
		ReceiptHandle:     ReceiptHandle,
		VisibilityTimeout: VisibilityTimeout, // (second)
	})
	if err != nil {
		return pkgError.WrapWithCode(err, pkgError.Update)
	}

	return nil
}

func (s *SqsListener) DeleteMessage(ctx context.Context, ReceiptHandle *string) error {
	_, err := s.SqsClient.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      s.QueueURL,
		ReceiptHandle: ReceiptHandle,
	})
	if err != nil {
		return pkgError.WrapWithCode(err, pkgError.Delete)
	}

	return nil
}

func (s *SqsListener) Listen(ctx context.Context, events ...string) (<-chan []Event, <-chan error, error) {
	if s == nil {
		return nil, nil, pkgError.Wrap(errors.New("queue listener is nil"))
	}

	sleepTime := 1 * time.Millisecond
	if s.sleepTime > 0 {
		sleepTime = time.Duration(s.sleepTime) * time.Second
	}
	eventCh := make(chan []Event, s.consumerCount*int(s.maxNumberOfMessages))
	errorCh := make(chan error, s.consumerCount*int(s.maxNumberOfMessages))

	wg := &sync.WaitGroup{}

	for i := 0; i < s.consumerCount; i++ {
		wg.Add(1)
		go func(consumerID int) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				default:
					s.ReceiveMessage(ctx, eventCh, errorCh, events, consumerID)
					time.Sleep(sleepTime)
				}
			}
		}(i)
	}

	go func() {
		<-ctx.Done()
		wg.Wait()
		close(eventCh)
		close(errorCh)
	}()

	return eventCh, errorCh, nil
}

func (s *SqsListener) ReceiveMessage(ctx context.Context, eventCh chan []Event, errorCh chan error, events []string, consumerID int) {
	receiveMsgResult, err := s.SqsClient.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		MessageAttributeNames: []string{
			string(sqsType.QueueAttributeNameAll),
		},
		QueueUrl:            s.QueueURL,
		MaxNumberOfMessages: s.maxNumberOfMessages,
		WaitTimeSeconds:     s.waitTime,
		VisibilityTimeout:   s.visibilityTimeOut,
	})
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

	receivedEvents := make([]Event, 0, len(receiveMsgResult.Messages))
	for _, msg := range receiveMsgResult.Messages {
		value, ok := msg.MessageAttributes["event_name"]
		if !ok {
			continue
		}

		eventName := aws.ToString(value.StringValue)

		bContinue := false
		for _, event := range events {
			if strings.EqualFold(eventName, event) {
				bContinue = true
				break
			}
		}

		if !bContinue {
			continue
		}

		message := aws.ToString(msg.Body)
		var event Event

		err = json.Unmarshal([]byte(message), &event)
		if err != nil {
			if deleteMessageErr := s.DeleteMessage(ctx, msg.ReceiptHandle); deleteMessageErr != nil {
				pkgLogger.ZapLogger.Logger.Sugar().Error(pkgError.Wrap(deleteMessageErr))
			}
			errorCh <- pkgError.Wrap(err)
			continue
		}

		event.ReceiptHandle = msg.ReceiptHandle

		receivedEvents = append(receivedEvents, event)
	}

	if len(receivedEvents) > 0 {
		eventCh <- receivedEvents
		pkgLogger.ZapLogger.Logger.Info(fmt.Sprintf("consumerID: %d, messageCount: %d", consumerID, len(receivedEvents)))
	}
}

func MustNewSqsEventListener(ctx context.Context, sqsClient *sqs.Client, queueName string, config SqsListenerConfig) EventListener {
	outInfo, err := sqsClient.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	})
	if err != nil {
		pkgLogger.ZapLogger.Logger.Fatal("Failed to get queue url", zap.Error(err))
	}

	// default setting
	if config.MaxMsg == 0 {
		config.MaxMsg = 1
	}

	if config.ConsumerCount == 0 {
		config.ConsumerCount = 1
	}

	return &SqsListener{
		SqsClient:           sqsClient,
		QueueURL:            outInfo.QueueUrl,
		maxNumberOfMessages: config.MaxMsg,
		waitTime:            config.WaitTime,
		visibilityTimeOut:   config.VisibilityTimeout,
		sleepTime:           config.SleepTime,
		consumerCount:       config.ConsumerCount,
	}
}
