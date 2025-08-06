package zincsearch

import (
	"context"
	"github.com/hgyowan/go-pkg-library/envs"
	"github.com/hgyowan/go-pkg-library/zincsearch/model"
	"github.com/hgyowan/go-pkg-library/zincsearch/param"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_ZincSearchIndex(t *testing.T) {
	zs := MustNewZincSearch(context.Background(), &ZinSearchConfig{
		Host:     envs.ZincSearchHost,
		Port:     envs.ZincSearchPort,
		Username: envs.ZincSearchUserName,
		Password: envs.ZincSearchPassword,
	})

	err := zs.Index().Create(&model.Index{
		Name: "test",
	})
	require.NoError(t, err)

	res, err := zs.Index().List(&param.IndexListRequest{
		Name: "test",
	})
	require.NoError(t, err)
	t.Log(res)

	err = zs.Index().UpdateMappings(&param.IndexUpdateMappingRequest{
		IndexName: "test",
		Mappings: &model.Mappings{
			Properties: map[string]model.Property{
				"name": {
					Type: "text",
				},
			},
		},
	})
	require.NoError(t, err)

	err = zs.Index().UpdateSettings(&param.IndexUpdateSettingsRequest{
		IndexName: "test",
		Settings: &model.IndexSettings{Analysis: &model.IndexAnalysis{
			Filter: map[string]interface{}{
				"lowercase": map[string]interface{}{
					"type": "lowercase",
				},
			},
		}},
	})
	require.NoError(t, err)

	res, err = zs.Index().List(&param.IndexListRequest{
		Name: "test",
	})
	require.NoError(t, err)
	t.Log(res)

	exists, err := zs.Index().Exists("test")
	require.NoError(t, err)
	t.Log(exists)

	err = zs.Index().Delete("test")
	require.NoError(t, err)

	exists, err = zs.Index().Exists("test")
	require.Equal(t, false, exists)
}

func Test_ZincSearchDocument(t *testing.T) {
	zs := MustNewZincSearch(context.Background(), &ZinSearchConfig{
		Host:     envs.ZincSearchHost,
		Port:     envs.ZincSearchPort,
		Username: envs.ZincSearchUserName,
		Password: envs.ZincSearchPassword,
	})

	err := zs.Document("test").Create(&param.DocumentCreateRequest{
		ID: "1",
		Document: &model.Document{
			Key: "name",
			Val: "test",
		},
	})
	require.NoError(t, err)

	err = zs.Document("test").Update(&param.DocumentUpdateRequest{
		ID: "1",
		Document: &model.Document{
			Key: "name",
			Val: "test2",
		},
	})
	require.NoError(t, err)

	err = zs.Document("test").BulkV2([]map[string]interface{}{
		{
			"name": "test3",
		},
		{
			"name": "test4",
		},
	})
	require.NoError(t, err)

	res, err := zs.Document("test").Search(&param.DocumentSearchRequest{
		Query: &model.QueryRequest{
			Field: "_all",
		},
		MaxResults: 10,
	})
	require.NoError(t, err)
	t.Log(res)

	for _, hit := range res.Hits.Hits {
		err = zs.Document("test").Delete(hit.ID)
		require.NoError(t, err)
	}
}
