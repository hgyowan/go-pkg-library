package model

import "time"

type Records struct {
	Docs map[string]interface{} `json:"docs"`
}

type Document struct {
	Key string      `json:"key"`
	Val interface{} `json:"val"`
}

type QueryRequest struct {
	Term      string     `json:"term,omitempty"`
	Field     string     `json:"field,omitempty"`
	StartTime *time.Time `json:"start_time,omitempty"`
	EndTime   *time.Time `json:"end_time,omitempty"`
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
	Total Total        `json:"total,omitempty"`
	Hits  []*HitRecord `json:"hits,omitempty"`
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
