package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"github.com/hgyowan/go-pkg-library/envs"
	pkgError "github.com/hgyowan/go-pkg-library/error"
	"io"
)

func Encrypt(plainText string) (string, error) {
	hash := sha256.Sum256([]byte(envs.AesSecretKey))
	aesBlock, err := aes.NewCipher(hash[:])
	if err != nil {
		return "", pkgError.WrapWithCode(err, pkgError.CryptData)
	}

	gcm, err := cipher.NewGCM(aesBlock)
	if err != nil {
		return "", pkgError.WrapWithCode(err, pkgError.CryptData)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", pkgError.WrapWithCode(err, pkgError.CryptData)
	}

	cipherText := gcm.Seal(nonce, nonce, []byte(plainText), nil)
	return base64.URLEncoding.EncodeToString(cipherText), nil
}

func Decrypt(cipherText []byte) (string, error) {
	hash := sha256.Sum256([]byte(envs.AesSecretKey))
	aesBlock, err := aes.NewCipher(hash[:])
	if err != nil {
		return "", pkgError.WrapWithCode(err, pkgError.CryptData)
	}

	gcm, err := cipher.NewGCM(aesBlock)
	if err != nil {
		return "", pkgError.WrapWithCode(err, pkgError.CryptData)
	}

	data, err := base64.URLEncoding.DecodeString(string(cipherText))
	if err != nil {
		return "", pkgError.WrapWithCode(err, pkgError.CryptData)
	}

	nonceSize := gcm.NonceSize()
	nonce, cipherText := data[:nonceSize], data[nonceSize:]

	plainText, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", pkgError.WrapWithCode(err, pkgError.CryptData)
	}

	return string(plainText), nil
}
