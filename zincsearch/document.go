package zincsearch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	pkgError "github.com/hgyowan/go-pkg-library/error"
	"github.com/hgyowan/go-pkg-library/zincsearch/model"
	"github.com/hgyowan/go-pkg-library/zincsearch/param"
	"resty.dev/v3"
)

type Document interface {
	Create(request *param.DocumentCreateRequest) error
	Update(request *param.DocumentUpdateRequest) error
	Delete(id string) error
	BulkV2(request []map[string]interface{}) error
	Search(request *param.DocumentSearchRequest) (*model.DocumentSearch, error)
	SearchES(request *param.DocumentSearchESRequest) (*model.DocumentSearch, error)
	DeleteBulkES(request []*param.DocumentDeleteBulkESRequest) error
}

type zinSearchDocument struct {
	zinSearchCli *resty.Request
	errHandler   func(body io.ReadCloser) error
	indexName    string
	options      options
}

func (z *zinSearchDocument) DeleteBulkES(request []*param.DocumentDeleteBulkESRequest) error {
	var buf bytes.Buffer

	for _, r := range request {
		line := map[string]interface{}{
			"delete": map[string]string{
				"_index": r.Index,
				"_id":    r.ID,
			},
		}

		b, err := json.Marshal(line)
		if err != nil {
			return pkgError.Wrap(err)
		}

		buf.Write(b)
		buf.WriteByte('\n')
	}

	res, err := z.zinSearchCli.
		SetBody(buf.Bytes()).
		SetHeader("Content-Type", "application/x-ndjson").
		Post("/es/_bulk")
	if err != nil {
		return pkgError.Wrap(err)
	}

	if res.StatusCode() != 200 {
		return pkgError.Wrap(z.errHandler(res.Body))
	}

	return nil
}

func (z *zinSearchDocument) Create(request *param.DocumentCreateRequest) error {
	z.indexName = fmt.Sprintf("%s%s", z.indexName, z.options.getSuffix())
	res, err := z.zinSearchCli.SetBody(map[string]interface{}{
		request.Document.Key: request.Document.Val,
	}).Put("/api/" + z.indexName + "/_doc/" + request.ID)
	if err != nil {
		return pkgError.Wrap(err)
	}

	if res.StatusCode() != 200 {
		return pkgError.Wrap(z.errHandler(res.Body))
	}

	return nil
}

func (z *zinSearchDocument) Update(request *param.DocumentUpdateRequest) error {
	z.indexName = fmt.Sprintf("%s%s", z.indexName, z.options.getSuffix())
	res, err := z.zinSearchCli.SetBody(map[string]interface{}{
		request.Document.Key: request.Document.Val,
	}).Post("/api/" + z.indexName + "/_update/" + request.ID)
	if err != nil {
		return pkgError.Wrap(err)
	}

	if res.StatusCode() != 200 {
		return pkgError.Wrap(z.errHandler(res.Body))
	}

	return nil
}

func (z *zinSearchDocument) Delete(id string) error {
	z.indexName = fmt.Sprintf("%s%s", z.indexName, z.options.getSuffix())
	res, err := z.zinSearchCli.Delete("/api/" + z.indexName + "/_doc/" + id)
	if err != nil {
		return pkgError.Wrap(err)
	}

	if res.StatusCode() != 200 {
		return pkgError.Wrap(z.errHandler(res.Body))
	}

	return nil
}

func (z *zinSearchDocument) BulkV2(request []map[string]interface{}) error {
	z.indexName = fmt.Sprintf("%s%s", z.indexName, z.options.getSuffix())
	res, err := z.zinSearchCli.SetBody(&param.DocumentBulkV2Request{
		Index:   z.indexName,
		Records: request,
	}).Post("/api/_bulkv2")
	if err != nil {
		return pkgError.Wrap(err)
	}

	if res.StatusCode() != 200 {
		return pkgError.Wrap(z.errHandler(res.Body))
	}

	return nil
}

func (z *zinSearchDocument) Search(request *param.DocumentSearchRequest) (*model.DocumentSearch, error) {
	z.indexName = fmt.Sprintf("%s%s", z.indexName, z.options.getSuffix())
	var resp *model.DocumentSearch
	res, err := z.zinSearchCli.
		SetBody(request).
		SetResult(&resp).
		Post("/api/" + z.indexName + "/_search")
	if err != nil {
		return nil, pkgError.Wrap(err)
	}

	if res.StatusCode() != 200 {
		return nil, pkgError.Wrap(z.errHandler(res.Body))
	}

	return resp, nil
}

func (z *zinSearchDocument) SearchES(request *param.DocumentSearchESRequest) (*model.DocumentSearch, error) {
	z.indexName = fmt.Sprintf("%s%s", z.indexName, z.options.getSuffix())
	var resp *model.DocumentSearch
	res, err := z.zinSearchCli.
		SetBody(request).
		SetResult(&resp).
		Post("/es/" + z.indexName + "/_search")
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
		options:      z.options,
	}
}
