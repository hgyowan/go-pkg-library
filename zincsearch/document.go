package zincsearch

import (
	pkgError "github.com/hgyowan/go-pkg-library/error"
	"github.com/hgyowan/go-pkg-library/zincsearch/model"
	"github.com/hgyowan/go-pkg-library/zincsearch/param"
	"io"
	"resty.dev/v3"
)

type Document interface {
	Create(request *param.DocumentCreateRequest) error
	Update(request *param.DocumentUpdateRequest) error
	Delete(id string) error
	BulkV2(request []map[string]interface{}) error
	Search(request *param.DocumentSearchRequest) (*model.DocumentSearch, error)
}

type zinSearchDocument struct {
	zinSearchCli *resty.Request
	errHandler   func(body io.ReadCloser) error
	indexName    string
}

func (z *zinSearchDocument) Create(request *param.DocumentCreateRequest) error {
	res, err := z.zinSearchCli.SetBody(map[string]interface{}{
		request.Document.Key: request.Document.Val,
	}).Put("/" + z.indexName + "/_doc/" + request.ID)
	if err != nil {
		return pkgError.Wrap(err)
	}

	if res.StatusCode() != 200 {
		return pkgError.Wrap(z.errHandler(res.Body))
	}

	return nil
}

func (z *zinSearchDocument) Update(request *param.DocumentUpdateRequest) error {
	res, err := z.zinSearchCli.SetBody(map[string]interface{}{
		request.Document.Key: request.Document.Val,
	}).Post("/" + z.indexName + "/_update/" + request.ID)
	if err != nil {
		return pkgError.Wrap(err)
	}

	if res.StatusCode() != 200 {
		return pkgError.Wrap(z.errHandler(res.Body))
	}

	return nil
}

func (z *zinSearchDocument) Delete(id string) error {
	res, err := z.zinSearchCli.Delete("/" + z.indexName + "/_doc/" + id)
	if err != nil {
		return pkgError.Wrap(err)
	}

	if res.StatusCode() != 200 {
		return pkgError.Wrap(z.errHandler(res.Body))
	}

	return nil
}

func (z *zinSearchDocument) BulkV2(request []map[string]interface{}) error {
	res, err := z.zinSearchCli.SetBody(&param.DocumentBulkV2Request{
		Index:   z.indexName,
		Records: request,
	}).Post("/_bulkv2")
	if err != nil {
		return pkgError.Wrap(err)
	}

	if res.StatusCode() != 200 {
		return pkgError.Wrap(z.errHandler(res.Body))
	}

	return nil
}

func (z *zinSearchDocument) Search(request *param.DocumentSearchRequest) (*model.DocumentSearch, error) {
	var resp *model.DocumentSearch
	res, err := z.zinSearchCli.
		SetBody(request).
		SetResult(&resp).
		Post("/" + z.indexName + "/_search")
	if err != nil {
		return nil, pkgError.Wrap(err)
	}

	if res.StatusCode() != 200 {
		return nil, pkgError.Wrap(z.errHandler(res.Body))
	}

	return resp, nil
}

func (z *zincSearch) Document(indexName string) Document {
	return &zinSearchDocument{
		zinSearchCli: z.cli,
		errHandler:   z.errHandler,
		indexName:    indexName,
	}
}
