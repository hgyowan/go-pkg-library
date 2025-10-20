package queue

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	sqsType "github.com/aws/aws-sdk-go-v2/service/sqs/types"
	pkgError "github.com/hgyowan/go-pkg-library/error"
	pkgLogger "github.com/hgyowan/go-pkg-library/logger"
	"go.uber.org/zap"
)

type SqsEmitter struct {
	SqsClient *sqs.Client
	QueueURL  *string
}

func (s *SqsEmitter) Close() error {
	return nil
}

func (s *SqsEmitter) Emit(ctx context.Context, event Event) error {
	data, err := json.Marshal(event)
	if err != nil {
		return pkgError.Wrap(err)
	}

	attributes := map[string]sqsType.MessageAttributeValue{
		"event_name": {
			DataType:    aws.String("String"),
			StringValue: aws.String(event.EventName),
		},
	}

	var messageDeduplicationId *string
	if event.MessageDeduplicationId != "" {
		messageDeduplicationId = aws.String(event.MessageDeduplicationId)
	}

	_, err = s.SqsClient.SendMessage(ctx, &sqs.SendMessageInput{
		MessageAttributes:      attributes,
		MessageBody:            aws.String(string(data)),
		QueueUrl:               s.QueueURL,
		MessageDeduplicationId: messageDeduplicationId,
		MessageGroupId:         aws.String(event.MessageGroupId),
	})
	if err != nil {
		return pkgError.Wrap(err)
	}

	return nil
}

func MustNewSqsEventEmitter(ctx context.Context, sqsClient *sqs.Client, queueName string) EventEmitter {
	urlInfo, err := sqsClient.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	})
	if err != nil {
		pkgLogger.ZapLogger.Logger.Fatal("Failed to get queue url", zap.Error(err))
	}

	return &SqsEmitter{
		SqsClient: sqsClient,
		QueueURL:  urlInfo.QueueUrl,
	}
}
