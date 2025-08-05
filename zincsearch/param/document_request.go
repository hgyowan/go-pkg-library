package param

import (
	"github.com/hgyowan/go-pkg-library/zincsearch/model"
)

type DocumentCreateRequest struct {
	ID       string         `json:"id"`
	Document model.Document `json:"document"`
}

type DocumentUpdateRequest struct {
	ID       string         `json:"id"`
	Document model.Document `json:"document"`
}

type DocumentBulkV2Request struct {
	ID      string        `json:"id"`
	Records model.Records `json:"records"`
}

type DocumentSearchRequest struct {
	SearchType string             `json:"search_type,omitempty"`
	Query      model.QueryRequest `json:"query,omitempty"`
	SortFields []string           `json:"sort_fields,omitempty"`
	From       int                `json:"from,omitempty"`
	MaxResults int                `json:"max_results,omitempty"`
	Source     []string           `json:"_source,omitempty"`
}
