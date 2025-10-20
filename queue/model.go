package queue

import (
	"encoding/json"
)

type Event struct {
	EventName string `json:"eventName"`
	// 메시지 고유 ID
	MessageDeduplicationId string          `json:"messageDeduplicationId"`
	MessageGroupId         string          `json:"messageGroupId"`
	Data                   json.RawMessage `json:"data"`
	// 메시지 받은 후 메시지를 제거할 때 사용
	ReceiptHandle *string `json:"-"`
}
