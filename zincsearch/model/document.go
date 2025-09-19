package model

import "time"

type Document struct {
	Key string      `json:"key"`
	Val interface{} `json:"val"`
}

type RangeFilter struct {
	Gte    string `json:"gte,omitempty"`    // greater than or equal
	Lte    string `json:"lte,omitempty"`    // less than or equal
	Format string `json:"format,omitempty"` // 날짜 포맷
}

type QueryCondition struct {
	Match map[string]string      `json:"match,omitempty"` // 예: {"City": "paris"}
	Term  map[string]string      `json:"term,omitempty"`  // 예: {"Country": "ger"}
	Range map[string]RangeFilter `json:"range,omitempty"` // 예: {"@timestamp": {...}}
}

type BoolQuery struct {
	Must    []QueryCondition `json:"must,omitempty"`     // AND 조건
	Should  []QueryCondition `json:"should,omitempty"`   // OR 조건
	MustNot []QueryCondition `json:"must_not,omitempty"` // NOT 조건
	Filter  []QueryCondition `json:"filter,omitempty"`   // 필터 조건
}

type QueryRequest struct {
	Term      string     `json:"term,omitempty"`
	Field     string     `json:"field,omitempty"`
	StartTime *time.Time `json:"start_time,omitempty"`
	EndTime   *time.Time `json:"end_time,omitempty"`
	// ES 용 Bool 쿼리
	Bool *BoolQuery `json:"bool,omitempty"`
}

type DocumentSearch struct {
	Took     int     `json:"took,omitempty"`
	TimedOut bool    `json:"timed_out,omitempty"`
	MaxScore float64 `json:"max_score,omitempty"`
	Hits     Hits    `json:"hits,omitempty"`
	Buckets  any     `json:"buckets,omitempty"`
	Error    string  `json:"error,omitempty"`
}

type Hits struct {
	Total    Total        `json:"total,omitempty"`
	Hits     []*HitRecord `json:"hits,omitempty"`
	MaxScore float64      `json:"max_score,omitempty"`
}

type Total struct {
	Value int `json:"value,omitempty"`
}

type HitRecord struct {
	Index     string                 `json:"_index,omitempty"`
	Type      string                 `json:"_type,omitempty"`
	ID        string                 `json:"_id,omitempty"`
	Score     float64                `json:"_score,omitempty"`
	Timestamp *time.Time             `json:"@timestamp,omitempty"`
	Source    map[string]interface{} `json:"_source,omitempty"`
}
