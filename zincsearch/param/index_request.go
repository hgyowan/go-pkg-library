package param

import "github.com/hgyowan/go-pkg-library/zincsearch/model"

type IndexListRequest struct {
	PageNum  int    `json:"page_num"`
	PageSize int    `json:"page_size"`
	SortBy   string `json:"sort_by"`
	Desc     string `json:"desc"`
	Name     string `json:"name"`
}

type IndexListResponse struct {
	List []*model.Index `json:"list"`
}

type IndexUpdateMappingRequest struct {
	IndexName string          `json:"index_name"`
	Mappings  *model.Mappings `json:"mappings"`
}

type IndexUpdateSettingsRequest struct {
	IndexName string               `json:"index_name"`
	Settings  *model.IndexSettings `json:"settings"`
}
