package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

type AuthData struct {
	Keys        map[string]KeyPair    `json:"keys"`
	Connections map[string]Connection `json:"connections"`
}

type Connection struct {
	PublicKey rsa.PublicKey `json:"public_key"`
	Host      string        `json:"host"`
}

type KeyPair struct {
	Name       string `json:"name"`
	PublicKey  string `json:"public"`
	PrivateKey string `json:"private"`
}

func GenerateKeyPair() (*KeyPair, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPem := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	publicKey := &privateKey.PublicKey
	publicKeyBytes := x509.MarshalPKCS1PublicKey(publicKey)
	publicKeyPem := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	return &KeyPair{
		PrivateKey: string(privateKeyPem),
		PublicKey:  string(publicKeyPem),
	}, nil
}

func (kp *KeyPair) Decrypt(ciphertext []byte) ([]byte, error) {
	block, _ := pem.Decode([]byte(kp.PrivateKey))
	if block == nil || block.Type != "PRIVATE KEY" {
		return nil, fmt.Errorf("wrong format for private key")
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	hash := sha256.New()
	return rsa.DecryptOAEP(hash, rand.Reader, privateKey, ciphertext, nil)
}

func (kp *KeyPair) Encrypt(plaintext []byte) ([]byte, error) {
	block, rest := pem.Decode([]byte(kp.PublicKey))
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, fmt.Errorf("wrong format for public key")
	}
	if len(rest) > 0 {
		return nil, fmt.Errorf("wrong format for public key")
	}

	publicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	hash := sha256.New()
	return rsa.EncryptOAEP(hash, rand.Reader, publicKey, plaintext, nil)
}
