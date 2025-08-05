package aes

import (
	pkgError "github.com/hgyowan/go-pkg-library/error"
	"reflect"
	"strings"
	"sync"
)

type Field struct {
	Name       string
	CryptoType string /*aes, cbc, fixed_cbc, random_cbc*/
	Context    string
}

type Schema struct {
	Fields map[string]*Field
}

type CryptoStorage struct {
	MasterKey []byte
	Storage   map[reflect.Type]*Schema
	Mu        sync.RWMutex
}

var storage *CryptoStorage

func MustNewCryptoHelper(masterKey []byte) {
	storage = &CryptoStorage{
		Storage:   make(map[reflect.Type]*Schema),
		MasterKey: masterKey,
	}
}

func (cs *CryptoStorage) loadSchema(model interface{}) (*Schema, error) {
	v := reflect.ValueOf(model)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return nil, pkgError.WrapWithCode(pkgError.EmptyBusinessError(), pkgError.CryptData, "loadSchema: must pass a non-nil pointer")
	}
	v = v.Elem()
	t := v.Type()

	storage.Mu.RLock()
	schema, ok := storage.Storage[t]
	storage.Mu.RUnlock()

	if !ok {
		schema = &Schema{Fields: make(map[string]*Field)}
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			tagVal, tagOk := field.Tag.Lookup("crypto")
			isString := field.Type.Kind() == reflect.String
			isStringPtr := field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.String

			if tagOk && (isString || isStringPtr) {
				tagMap, err := parseCryptoTag(tagVal)
				if err != nil {
					return nil, pkgError.Wrap(err)
				}
				switch tagMap[CryptoType] {
				case CryptoTypeFixedCBC, CryptoTypeRandomCBC:
					contextVal, contextOk := tagMap[CryptoContext]
					if contextVal == "" || !contextOk {
						continue
					}
				}
				schema.Fields[field.Name] = &Field{
					Name:       field.Name,
					CryptoType: tagMap[CryptoType],
					Context:    tagMap[CryptoContext],
				}

			}
		}
		storage.Mu.Lock()
		storage.Storage[t] = schema
		storage.Mu.Unlock()
	}

	return schema, nil
}

func EncryptScheme(model interface{}) error {
	v := reflect.ValueOf(model)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return pkgError.WrapWithCode(pkgError.EmptyBusinessError(), pkgError.CryptData, "EncryptScheme: model is nil")
	}
	v = v.Elem()

	schema, err := storage.loadSchema(model)
	if err != nil {
		return pkgError.Wrap(err)
	}

	for fieldName, fieldMeta := range schema.Fields {
		fieldVal := v.FieldByName(fieldName)
		isString := fieldVal.Kind() == reflect.String
		isStringPtr := fieldVal.Kind() == reflect.Ptr && fieldVal.Type().Elem().Kind() == reflect.String

		if fieldVal.IsValid() && fieldVal.CanSet() && (isString || isStringPtr) {
			var err error

			var original string

			if isString {
				original = fieldVal.String()
			}

			if isStringPtr {
				if fieldVal.IsNil() {
					continue
				}
				original = fieldVal.Elem().String()
			}

			encrypted := original
			switch fieldMeta.CryptoType {
			case CryptoTypeAES:
				encrypted, err = Encrypt(original)
				if err != nil {
					return pkgError.Wrap(err)
				}
			case CryptoTypeCBC:
				encrypted, err = CBCEncrypt(original)
				if err != nil {
					return pkgError.Wrap(err)
				}
			case CryptoTypeFixedCBC:
				encrypted, err = CBCEncryptWithFixedKey(original, fieldMeta.Context, storage.MasterKey)
				if err != nil {
					return pkgError.Wrap(err)
				}
			case CryptoTypeRandomCBC:
				encrypted, err = CBCEncryptWithRandomIV(original, fieldMeta.Context, storage.MasterKey)
				if err != nil {
					return pkgError.Wrap(err)
				}
			}

			if isString {
				fieldVal.SetString(encrypted)
			}

			if isStringPtr {
				if fieldVal.IsNil() {
					continue
				}
				fieldVal.Elem().SetString(encrypted)
			}

		}
	}

	return nil
}

func DecryptScheme(model interface{}) error {
	v := reflect.ValueOf(model)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return pkgError.WrapWithCode(pkgError.EmptyBusinessError(), pkgError.CryptData, "DecryptScheme: model is nil")
	}
	v = v.Elem()

	schema, err := storage.loadSchema(model)
	if err != nil {
		return pkgError.Wrap(err)
	}

	for fieldName, fieldMeta := range schema.Fields {
		fieldVal := v.FieldByName(fieldName)
		isString := fieldVal.Kind() == reflect.String
		isStringPtr := fieldVal.Kind() == reflect.Ptr && fieldVal.Type().Elem().Kind() == reflect.String

		if fieldVal.IsValid() && fieldVal.CanSet() && (isString || isStringPtr) {
			var err error

			var original string

			if isString {
				original = fieldVal.String()
			}

			if isStringPtr {
				if fieldVal.IsNil() {
					continue
				}
				original = fieldVal.Elem().String()
			}

			decrypted := original
			switch fieldMeta.CryptoType {
			case CryptoTypeAES:
				decrypted, err = Decrypt([]byte(original))
				if err != nil {
					return pkgError.Wrap(err)
				}
			case CryptoTypeCBC:
				decrypted, err = CBCDecrypt(original)
				if err != nil {
					return pkgError.Wrap(err)
				}
			case CryptoTypeFixedCBC, CryptoTypeRandomCBC:
				decrypted, err = CBCDecryptWithIV(original, fieldMeta.Context, storage.MasterKey)
				if err != nil {
					return pkgError.Wrap(err)
				}
			}

			if isString {
				fieldVal.SetString(decrypted)
			}

			if isStringPtr {
				if fieldVal.IsNil() {
					continue
				}
				fieldVal.Elem().SetString(decrypted)
			}

		}
	}

	return nil
}

func parseCryptoTag(tag string) (map[string]string, error) {
	result := make(map[string]string)

	if strings.TrimSpace(tag) == "" {
		return result, nil
	}

	pairs := strings.Split(tag, ";")
	for _, pair := range pairs {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}

		kv := strings.SplitN(pair, ":", 2)
		key := strings.TrimSpace(kv[0])
		val := ""
		if len(kv) == 2 {
			val = strings.TrimSpace(kv[1])
		}

		if key != "" {
			result[key] = val
		}
	}

	return result, nil
}
