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
	options    options
}

type ZinSearchConfig struct {
	Host     string
	Port     string
	Username string
	Password string
}

func MustNewZincSearch(ctx context.Context, config *ZinSearchConfig, opts ...ZincSearchOption) ZincSearch {
	restyClient := resty.New().
		SetContext(ctx).
		SetBaseURL(fmt.Sprintf("http://%s:%s", config.Host, config.Port)).
		SetBasicAuth(config.Username, config.Password).R()

	// ping
	_, err := restyClient.Get("/version")
	if err != nil {
		pkgLogger.ZapLogger.Logger.Fatal(err.Error())
	}

	z := &zincSearch{
		cli: restyClient,
		errHandler: func(body io.ReadCloser) error {
			var resp model.ErrorResponse
			if err := json.NewDecoder(body).Decode(&resp); err != nil {
				return pkgError.Wrap(err)
			}

			switch resp.Error {
			case "id not found":
				return pkgError.WrapWithCode(pkgError.EmptyBusinessError(), pkgError.NotFound)
			}

			return pkgError.WrapWithMessage(pkgError.EmptyBusinessError(), resp.Error)
		},
	}

	for _, opt := range opts {
		opt(z)
	}

	if z.options.getSuffix == nil {
		z.options.getSuffix = func() string {
			return ""
		}
	}

	return z
}
