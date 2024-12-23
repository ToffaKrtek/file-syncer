package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateKeyPair(t *testing.T) {
	_, err := GenerateKeyPair()
	assert := assert.New(t)
	assert.Equal(err, nil)
}

func TestEncrypt(t *testing.T) {
	texts := []string{
		"some text",
		"Какой-то текст",
		"123",
		"-!",
		" ",
	}
	assert := assert.New(t)
	kp, err := GenerateKeyPair()
	if err == nil {
		for _, text := range texts {
			res, err := kp.Encrypt([]byte(text))
			assert.Nil(err)
			assert.NotEqual(string(res), text)
		}
	} else {
		t.Errorf("Error generate key pair")
	}
}

func TestDecrypt(t *testing.T) {
	texts := []string{
		"some text",
		"Какой-то текст",
		"123",
		"-!",
		" ",
	}
	assert := assert.New(t)
	kp, err := GenerateKeyPair()

	if err == nil {
		for _, text := range texts {
			res, err := kp.Encrypt([]byte(text))
			assert.Nil(err)
			if err == nil {
				res, err = kp.Decrypt(res)
				assert.Nil(err)
				assert.Equal(string(res), text)
			}
		}
	} else {
		t.Errorf("Error generate key pair")
	}
}
