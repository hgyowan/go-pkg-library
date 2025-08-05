package zincsearch

import (
	pkgError "github.com/hgyowan/go-pkg-library/error"
	"github.com/hgyowan/go-pkg-library/zincsearch/model"
	"github.com/hgyowan/go-pkg-library/zincsearch/param"
	"io"
	"resty.dev/v3"
	"strconv"
)

type Index interface {
	Create(index *model.Index) error
	Update(index *model.Index) error
	Delete(indexName string) error
	List(request *param.IndexListRequest) (*param.IndexListResponse, error)
	Exists(indexName string) (bool, error)
}

type zinSearchIndex struct {
	cli        *resty.Request
	errHandler func(body io.ReadCloser) error
}

func (z *zinSearchIndex) Create(index *model.Index) error {
	res, err := z.cli.SetBody(index).Post("/index")
	if err != nil {
		return pkgError.Wrap(err)
	}

	if res.StatusCode() != 200 {
		return pkgError.Wrap(z.errHandler(res.Body))
	}

	return nil
}

func (z *zinSearchIndex) Update(index *model.Index) error {
	res, err := z.cli.SetBody(index).Put("/index")
	if err != nil {
		return pkgError.Wrap(err)
	}

	if res.StatusCode() != 200 {
		return pkgError.Wrap(z.errHandler(res.Body))
	}

	return nil
}

func (z *zinSearchIndex) Delete(indexName string) error {
	res, err := z.cli.Delete("/index/" + indexName)
	if err != nil {
		return pkgError.Wrap(err)
	}

	if res.StatusCode() != 200 {
		return pkgError.Wrap(z.errHandler(res.Body))
	}

	return nil
}

func (z *zinSearchIndex) List(request *param.IndexListRequest) (*param.IndexListResponse, error) {
	queryParams := make(map[string]string)
	if request.Name != "" {
		queryParams["name"] = request.Name
	}

	if request.Desc != "" {
		queryParams["desc"] = request.Desc
	}

	if request.SortBy != "" {
		queryParams["sort_by"] = request.SortBy
	}

	if request.PageNum != 0 {
		queryParams["page_num"] = strconv.Itoa(request.PageNum)
	}

	if request.PageSize != 0 {
		queryParams["page_size"] = strconv.Itoa(request.PageSize)
	}

	var resp *param.IndexListResponse
	res, err := z.cli.
		SetQueryParams(queryParams).
		SetResult(&resp).
		Get("/index")
	if err != nil {
		return nil, pkgError.Wrap(err)
	}

	if res.StatusCode() != 200 {
		return nil, pkgError.Wrap(z.errHandler(res.Body))
	}

	return resp, nil
}

func (z *zinSearchIndex) Exists(indexName string) (bool, error) {
	res, err := z.cli.Get("/index/" + indexName)
	if err != nil {
		return false, pkgError.Wrap(err)
	}

	if res.StatusCode() != 200 {
		return false, pkgError.Wrap(z.errHandler(res.Body))
	}

	return true, nil
}

func (z *zincSearch) Index() Index {
	return &zinSearchIndex{
		cli:        z.cli,
		errHandler: z.errHandler,
	}
}
