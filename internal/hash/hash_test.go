package hash

import (
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHash(t *testing.T) {
	testFile, err := ioutil.TempFile("", "testfile")
	if err != nil {
		t.Errorf("Error generate temp file: %v", err)
	}
	defer os.Remove(testFile.Name())
	content := []byte("Test data")
	if _, err := testFile.Write(content); err != nil {
		t.Errorf("Error set data to temp file: %v", err)
	}
	if err := testFile.Close(); err != nil {
		t.Errorf("Error close temp file: %v", err)
	}
	expectedHash := sha256.Sum256(content)
	expectedHashString := hex.EncodeToString(expectedHash[:])
	assert := assert.New(t)

	hashFile, err := hashFile(testFile.Name())
	assert.Nil(err)
	assert.Equal(hashFile, expectedHashString)

	hashFile1, err := Hash(testFile.Name(), false)
	assert.Nil(err)
	assert.Equal(hashFile1, expectedHashString)
}
