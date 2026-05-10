package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
	"strings"
)

func EncryptString(keyText, plainText string) (string, error) {
	if strings.TrimSpace(plainText) == "" {
		return "", errors.New("plain text is required")
	}

	gcm, err := newGCM(keyText)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	sealed := gcm.Seal(nonce, nonce, []byte(plainText), nil)
	return base64.StdEncoding.EncodeToString(sealed), nil
}

func DecryptString(keyText, cipherText string) (string, error) {
	if strings.TrimSpace(cipherText) == "" {
		return "", errors.New("cipher text is required")
	}

	gcm, err := newGCM(keyText)
	if err != nil {
		return "", err
	}

	raw, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(raw) < nonceSize {
		return "", errors.New("cipher text is invalid")
	}

	nonce, data := raw[:nonceSize], raw[nonceSize:]
	plain, err := gcm.Open(nil, nonce, data, nil)
	if err != nil {
		return "", err
	}

	return string(plain), nil
}

func newGCM(keyText string) (cipher.AEAD, error) {
	if strings.TrimSpace(keyText) == "" {
		return nil, errors.New("secret cipher key is required")
	}

	sum := sha256.Sum256([]byte(keyText))
	block, err := aes.NewCipher(sum[:])
	if err != nil {
		return nil, err
	}

	return cipher.NewGCM(block)
}
