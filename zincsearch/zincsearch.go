package zincsearch

import (
	"context"
	"encoding/json"
	"fmt"
	pkgError "github.com/hgyowan/go-pkg-library/error"
	pkgLogger "github.com/hgyowan/go-pkg-library/logger"
	"github.com/hgyowan/go-pkg-library/zincsearch/model"
	"io"
	"resty.dev/v3"
)

type ZincSearch interface {
	Index() Index
	Document(indexName string) Document
}

type zincSearch struct {
	cli        *resty.Request
	errHandler func(body io.ReadCloser) error
}

type ZinSearchConfig struct {
	Host     string
	Port     string
	Username string
	Password string
}

func MustNewZincSearch(ctx context.Context, config *ZinSearchConfig) ZincSearch {
	restyClient := resty.New().
		SetContext(ctx).
		SetBaseURL(fmt.Sprintf("http://%s:%s/api", config.Host, config.Port)).
		SetBasicAuth(config.Username, config.Password).R()

	// ping
	_, err := restyClient.Get("/version")
	if err != nil {
		pkgLogger.ZapLogger.Logger.Fatal(err.Error())
	}

	return &zincSearch{
		cli: restyClient,
		errHandler: func(body io.ReadCloser) error {
			var resp model.ErrorResponse
			if err := json.NewDecoder(body).Decode(&resp); err != nil {
				return pkgError.Wrap(err)
			}

			return pkgError.WrapWithMessage(pkgError.EmptyBusinessError(), resp.Error)
		},
	}
}
