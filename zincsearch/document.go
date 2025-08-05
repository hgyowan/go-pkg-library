package zincsearch

import (
	"github.com/hgyowan/go-pkg-library/zincsearch/model"
	"github.com/hgyowan/go-pkg-library/zincsearch/param"
	"io"
	"resty.dev/v3"
)

type Document interface {
	Create(request *param.DocumentCreateRequest) error
	Update(request *param.DocumentUpdateRequest) error
	Delete(id string) error
	BulkV2(request *param.DocumentBulkV2Request) error
	Search(request *param.DocumentSearchRequest) (*model.DocumentSearch, error)
}

type zinSearchDocument struct {
	zinSearchCli *resty.Request
	errHandler   func(body io.ReadCloser) error
	indexName    string
}

func (z *zinSearchDocument) Create(request *param.DocumentCreateRequest) error {
	//TODO implement me
	panic("implement me")
}

func (z *zinSearchDocument) Update(request *param.DocumentUpdateRequest) error {
	//TODO implement me
	panic("implement me")
}

func (z *zinSearchDocument) Delete(id string) error {
	//TODO implement me
	panic("implement me")
}

func (z *zinSearchDocument) BulkV2(request *param.DocumentBulkV2Request) error {
	//TODO implement me
	panic("implement me")
}

func (z *zinSearchDocument) Search(request *param.DocumentSearchRequest) (*model.DocumentSearch, error) {
	//TODO implement me
	panic("implement me")
}

func (z *zincSearch) Document(indexName string) Document {
	return &zinSearchDocument{
		zinSearchCli: z.cli,
		errHandler:   z.errHandler,
		indexName:    indexName,
	}
}
