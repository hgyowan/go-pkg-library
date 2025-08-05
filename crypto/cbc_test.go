package aes

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_CBC(t *testing.T) {
	encrypt, err := CBCEncrypt("test@yopmail.com")
	require.NoError(t, err)

	t.Log(encrypt)

	decrypt, err := CBCDecrypt(encrypt)
	require.NoError(t, err)

	t.Log(decrypt)
}

func Test_CBCFixed(t *testing.T) {
	encrypt, err := CBCEncryptWithFixedKey("test@yopmail.com", "user:email", []byte{})
	require.NoError(t, err)

	t.Log(encrypt)

	decrypt, err := CBCDecryptWithIV(encrypt, "user:email", []byte{})
	require.NoError(t, err)

	t.Log(decrypt)
}

func Test_CBCRandom(t *testing.T) {
	encrypt, err := CBCEncryptWithRandomIV("test@yopmail.com", "user:email", []byte{})
	require.NoError(t, err)

	t.Log(encrypt)

	decrypt, err := CBCDecryptWithIV(encrypt, "user:email", []byte{})
	require.NoError(t, err)

	t.Log(decrypt)
}
