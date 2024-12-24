package hash

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"os"

	"golang.org/x/mod/sumdb/dirhash"
)

func Encrypt(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	cipherText := gcm.Seal(nonce, nonce, data, nil)
	return cipherText, nil
}

func Decrypt(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	if len(data) < gcm.NonceSize() {
		return nil, errors.New("неправильный размер шифро-текста")
	}
	nonce, cipherText := data[:gcm.NonceSize()], data[gcm.NonceSize():]
	res, err := gcm.Open(nil, nonce, cipherText, nil)
	return res, err
}

func Hash(path string, isDir bool) (string, error) {
	if isDir {
		return hashDir(path)
	}
	return hashFile(path)
}

func hashFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func hashDir(path string) (string, error) {
	return dirhash.HashDir(path, "", dirhash.DefaultHash)
}
