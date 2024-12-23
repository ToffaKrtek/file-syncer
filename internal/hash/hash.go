package hash

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"

	"golang.org/x/mod/sumdb/dirhash"
)

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
