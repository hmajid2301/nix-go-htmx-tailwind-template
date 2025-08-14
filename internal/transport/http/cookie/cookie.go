package cookie

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var (
	ErrInvalidKey       = errors.New("invalid encryption key")
	ErrDecryptionFailed = errors.New("decryption failed")
)

func Decode[T any](encryptionKey []byte, cookie *http.Cookie) (T, error) {
	var zero T

	decodedState, err := base64.URLEncoding.DecodeString(cookie.Value)
	if err != nil {
		return zero, fmt.Errorf("base64 decode failed: %w", err)
	}

	decryptedState, err := decryptState(encryptionKey, decodedState)
	if err != nil {
		return zero, fmt.Errorf("%w: %v", ErrDecryptionFailed, err)
	}

	var data T
	if err := json.Unmarshal(decryptedState, &data); err != nil {
		return zero, fmt.Errorf("json unmarshal failed: %w", err)
	}

	return data, nil
}

func decryptState(key []byte, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < gcm.NonceSize() {
		return nil, errors.New("malformed ciphertext")
	}

	nonce, ciphertext := ciphertext[:gcm.NonceSize()], ciphertext[gcm.NonceSize():]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

func Encode(encryptionKey []byte, name string, data any) (http.Cookie, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return http.Cookie{}, fmt.Errorf("json marshal failed: %w", err)
	}

	encryptedData, err := encryptState(encryptionKey, jsonData)
	if err != nil {
		return http.Cookie{}, fmt.Errorf("encryption failed: %w", err)
	}

	return http.Cookie{
		Name:     name,
		Value:    base64.URLEncoding.EncodeToString(encryptedData),
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
	}, nil
}

func encryptState(key []byte, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}