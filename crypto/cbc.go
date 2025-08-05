package aes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"github.com/hgyowan/go-pkg-library/envs"
	pkgError "github.com/hgyowan/go-pkg-library/error"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"io"
	"strings"
)

// utf8 디코딩
func decodeUtf8(input []byte) (string, error) {
	decoder := unicode.UTF8.NewDecoder()

	reader := transform.NewReader(strings.NewReader(string(input)), decoder)

	decoded, err := io.ReadAll(reader)
	if err != nil {
		return "", pkgError.WrapWithCode(err, pkgError.CryptData)
	}

	return string(decoded), nil
}

// context 기반 파생키 생성
func deriveKey(context string, masterKey []byte) ([]byte, error) {
	decodeMasterKey, err := decodeUtf8(masterKey)
	if err != nil {
		return nil, pkgError.WrapWithCode(err, pkgError.CryptData)
	}

	key := sha256.Sum256([]byte(decodeMasterKey + context))
	return key[:], nil
}

// PKCS7 패딩 제거
func pkcs7UnPad(data []byte) []byte {
	length := len(data)
	unPadding := int(data[length-1])
	return data[:(length - unPadding)]
}

// PKCS7 패딩 적용
func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

// 비문 여부 판단 (Base64 디코딩 + aes.BlockSize 이상 여부)
func isEncrypted(cipherTextBase64 string) bool {
	data, err := base64.StdEncoding.DecodeString(cipherTextBase64)
	if err != nil {
		return false
	}
	return len(data) > aes.BlockSize
}

// 공통 암호화 로직
func encrypt(plainText string, context string, masterKey []byte, iv []byte) (string, error) {
	if plainText == "" {
		return "", nil
	}

	if isEncrypted(plainText) {
		return plainText, nil
	}
	key, err := deriveKey(context, masterKey)
	if err != nil {
		return "", pkgError.WrapWithCode(err, pkgError.CryptData)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", pkgError.WrapWithCode(err, pkgError.CryptData)
	}

	paddedText := pkcs7Pad([]byte(plainText), aes.BlockSize)

	mode := cipher.NewCBCEncrypter(block, iv)
	ciphertext := make([]byte, len(paddedText))
	mode.CryptBlocks(ciphertext, paddedText)

	// 암호문 = IV + ciphertext
	final := append(iv, ciphertext...)
	return base64.StdEncoding.EncodeToString(final), nil
}

func CBCEncrypt(plainText string) (string, error) {
	block, err := aes.NewCipher([]byte(envs.AesSecretKey))
	if err != nil {
		return "", pkgError.WrapWithCode(err, pkgError.CryptData)
	}

	paddedText := pkcs7Pad([]byte(plainText), block.BlockSize())

	iv, err := hex.DecodeString(envs.AesSecretIVKey)
	if err != nil {
		return "", pkgError.WrapWithCode(err, pkgError.CryptData)
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	ciphertext := make([]byte, len(paddedText))
	mode.CryptBlocks(ciphertext, paddedText)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func CBCDecrypt(cipherText string) (string, error) {
	block, err := aes.NewCipher([]byte(envs.AesSecretKey))
	if err != nil {
		return "", pkgError.WrapWithCode(err, pkgError.CryptData)
	}

	data, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", pkgError.WrapWithCode(err, pkgError.CryptData)
	}

	iv, err := hex.DecodeString(envs.AesSecretIVKey)
	if err != nil {
		return "", pkgError.WrapWithCode(err, pkgError.CryptData)
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	decrypted := make([]byte, len(data))
	mode.CryptBlocks(decrypted, data)

	unPaddedText := pkcs7UnPad(decrypted)
	return string(unPaddedText), nil
}

// CBCEncryptWithFixedKey
// 암호화 (IV + 암호문 → Base64 인코딩)
func CBCEncryptWithFixedKey(plainText string, context string, masterKey []byte) (string, error) {
	iv := make([]byte, aes.BlockSize)
	copy(iv, envs.AesSecretIVKey)

	return encrypt(plainText, context, masterKey, iv)
}

// CBCEncryptWithRandomIV
// 암호화 (랜덤 IV + 암호문 → Base64 인코딩)
func CBCEncryptWithRandomIV(plainText string, context string, masterKey []byte) (string, error) {
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}
	return encrypt(plainText, context, masterKey, iv)
}

// CBCDecryptWithIV
// 복호화 (Base64 → IV 분리 → 복호화)
func CBCDecryptWithIV(cipherTextBase64 string, context string, masterKey []byte) (string, error) {
	if cipherTextBase64 == "" {
		return "", nil
	}

	if !isEncrypted(cipherTextBase64) {
		return cipherTextBase64, nil
	}

	key, err := deriveKey(context, masterKey)
	if err != nil {
		return "", pkgError.WrapWithCode(err, pkgError.CryptData)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", pkgError.WrapWithCode(err, pkgError.CryptData)
	}

	data, err := base64.StdEncoding.DecodeString(cipherTextBase64)
	if err != nil {
		return "", pkgError.WrapWithCode(err, pkgError.CryptData)
	}

	// IV 추출
	iv := data[:aes.BlockSize]
	ciphertext := data[aes.BlockSize:]

	if len(ciphertext)%aes.BlockSize != 0 {
		return string(ciphertext), nil
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	decrypted := make([]byte, len(ciphertext))
	mode.CryptBlocks(decrypted, ciphertext)

	unPaddedText := pkcs7UnPad(decrypted)
	return string(unPaddedText), nil
}
