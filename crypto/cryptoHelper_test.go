package aes

import (
	pkgLogger "github.com/hgyowan/go-pkg-library/logger"
	"github.com/stretchr/testify/require"
	"testing"
)

type Model struct {
	Name  string  `crypto:"type:aes;"`
	Phone *string `crypto:"type:fixed_cbc;context:user:cell-no"`
	Email string  `crypto:"type:random_cbc;context:user:email"`
}

func (m *Model) BeforeCreate() error {
	return EncryptScheme(m)
}

func (m *Model) AfterFind() error {
	return DecryptScheme(m)
}

func Test_Crypto(t *testing.T) {
	pkgLogger.MustInitZapLogger()
	MustNewCryptoHelper([]byte("masterKey"))
	phone := "tell"
	model := &Model{
		Name:  "name",
		Phone: &phone,
		Email: "test@test.com",
	}

	err := model.BeforeCreate()
	require.NoError(t, err)

	t.Log(model)

	err = model.AfterFind()
	require.NoError(t, err)

	t.Log(model)
}
