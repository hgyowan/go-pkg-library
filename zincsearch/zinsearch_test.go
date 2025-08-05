package zincsearch

import (
	"context"
	"github.com/hgyowan/go-pkg-library/envs"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_ZincSearch(t *testing.T) {
	zs := MustNewZincSearch(context.Background(), &ZinSearchConfig{
		Host:     envs.ZinSearchHost,
		Port:     envs.ZinSearchPort,
		Username: envs.ZinSearchUserName,
		Password: envs.ZinSearchPassword,
	})

	res, err := zs.Index().Exists("test")
	require.NoError(t, err)
	t.Log(res)
}
