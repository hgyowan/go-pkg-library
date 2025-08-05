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
	UpdateMappings(request *param.IndexUpdateMappingRequest) error
	UpdateSettings(request *param.IndexUpdateSettingsRequest) error
	Delete(indexName string) error
	List(request *param.IndexListRequest) (*param.IndexListResponse, error)
	Exists(indexName string) (bool, error)
}

type zinSearchIndex struct {
	zinSearchCli *resty.Request
	errHandler   func(body io.ReadCloser) error
}

func (z *zinSearchIndex) Create(index *model.Index) error {
	res, err := z.zinSearchCli.SetBody(index).Post("/index")
	if err != nil {
		return pkgError.Wrap(err)
	}

	if res.StatusCode() != 200 {
		return pkgError.Wrap(z.errHandler(res.Body))
	}

	return nil
}

func (z *zinSearchIndex) UpdateMappings(request *param.IndexUpdateMappingRequest) error {
	res, err := z.zinSearchCli.SetBody(request.Mappings).Put("/" + request.IndexName + "/_mapping")
	if err != nil {
		return pkgError.Wrap(err)
	}

	if res.StatusCode() != 200 {
		return pkgError.Wrap(z.errHandler(res.Body))
	}

	return nil
}

func (z *zinSearchIndex) UpdateSettings(request *param.IndexUpdateSettingsRequest) error {
	res, err := z.zinSearchCli.SetBody(request.Settings).Put("/" + request.IndexName + "/_settings")
	if err != nil {
		return pkgError.Wrap(err)
	}

	if res.StatusCode() != 200 {
		return pkgError.Wrap(z.errHandler(res.Body))
	}

	return nil
}

func (z *zinSearchIndex) Delete(indexName string) error {
	res, err := z.zinSearchCli.Delete("/index/" + indexName)
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
	res, err := z.zinSearchCli.
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
	res, err := z.zinSearchCli.Get("/index/" + indexName)
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
		zinSearchCli: z.cli,
		errHandler:   z.errHandler,
	}
}
