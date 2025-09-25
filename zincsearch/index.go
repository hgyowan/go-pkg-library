package zincsearch

import (
	"fmt"
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
	options      options
}

func (z *zinSearchIndex) Create(index *model.Index) error {
	index.Name = fmt.Sprintf("%s%s", index.Name, z.options.getSuffix())
	res, err := z.zinSearchCli.SetBody(index).Post("/api/index")
	if err != nil {
		return pkgError.Wrap(err)
	}

	if res.StatusCode() != 200 {
		return pkgError.Wrap(z.errHandler(res.Body))
	}

	return nil
}

func (z *zinSearchIndex) UpdateMappings(request *param.IndexUpdateMappingRequest) error {
	request.IndexName = fmt.Sprintf("%s%s", request.IndexName, z.options.getSuffix())
	res, err := z.zinSearchCli.SetBody(request.Mappings).Put("/api/" + request.IndexName + "/_mapping")
	if err != nil {
		return pkgError.Wrap(err)
	}

	if res.StatusCode() != 200 {
		return pkgError.Wrap(z.errHandler(res.Body))
	}

	return nil
}

func (z *zinSearchIndex) UpdateSettings(request *param.IndexUpdateSettingsRequest) error {
	request.IndexName = fmt.Sprintf("%s%s", request.IndexName, z.options.getSuffix())
	res, err := z.zinSearchCli.SetBody(request.Settings).Put("/api/" + request.IndexName + "/_settings")
	if err != nil {
		return pkgError.Wrap(err)
	}

	if res.StatusCode() != 200 {
		return pkgError.Wrap(z.errHandler(res.Body))
	}

	return nil
}

func (z *zinSearchIndex) Delete(indexName string) error {
	indexName = fmt.Sprintf("%s%s", indexName, z.options.getSuffix())
	res, err := z.zinSearchCli.Delete("/api/index/" + indexName)
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
		request.Name = fmt.Sprintf("%s%s", request.Name, z.options.getSuffix())
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
		Get("/api/index")
	if err != nil {
		return nil, pkgError.Wrap(err)
	}

	if res.StatusCode() != 200 {
		return nil, pkgError.Wrap(z.errHandler(res.Body))
	}

	return resp, nil
}

func (z *zinSearchIndex) Exists(indexName string) (bool, error) {
	indexName = fmt.Sprintf("%s%s", indexName, z.options.getSuffix())
	res, err := z.zinSearchCli.Get("/api/index/" + indexName)
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
		options:      z.options,
	}
}
