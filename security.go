package roomapi

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"hash/crc32"
	"io"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
)

func encrypt(str, salt, aesKey string) (string, error) {
	str += salt

	key := []byte(aesKey)
	plaintext := []byte(str)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	padding := aes.BlockSize - (len(plaintext) % aes.BlockSize)
	paddedPlaintext := append(plaintext, bytes.Repeat([]byte{byte(padding)}, padding)...)

	ciphertext := make([]byte, aes.BlockSize+len(paddedPlaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], paddedPlaintext)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func decrypt(str, salt, aesKey string) (string, error) {
	key := []byte(aesKey)

	ciphertext, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)

	plaintext := string(ciphertext)
	plaintext = removeBackSpace(plaintext)
	plaintext = plaintext[:len(plaintext)-len(salt)]

	return plaintext, nil
}

func removeBackSpace(str string) string {
	for i, v := range str {
		if v < 32 {
			return str[:i]
		}
	}
	return str
}

func accessSecretVersion(name string) (secretData string) {
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		Error(errors.New("ERROR_CREATE_SECRETMANAGER_CLIENT"))
		return
	}
	defer client.Close()

	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	}

	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		Error(errors.New("ERROR_ACCESS_SECRET_VERSION"))
		return
	}

	crc32c := crc32.MakeTable(crc32.Castagnoli)
	checksum := int64(crc32.Checksum(result.Payload.Data, crc32c))
	if checksum != *result.Payload.DataCrc32C {
		Error(errors.New("ERROR_DATA_CORRUPTION"))
		return
	}

	secretData = string(result.Payload.Data)
	return
}
