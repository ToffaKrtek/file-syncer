package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/ToffaKrtek/file-syncer/internal/hash"
	"github.com/joho/godotenv"
)

var (
	authEncryptionKey []byte
	authDataFile      = "./authdata.json"
	authLoaded        *AuthData
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic("Ошибка парсинга .env-файла")
	}

	authEncryptionKey = []byte(os.Getenv("AUTH_ENCRYPT_KEY"))
	if len(authEncryptionKey) != 32 {
		panic("AUTH_ENCRYPT_KEY должен быть длинной 32 байта для AES-256")
	}
}

type AuthData struct {
	Keys        map[string]KeyPair    `json:"keys"`
	Connections map[string]Connection `json:"connections"`
}

func Auth() *AuthData {
	if authLoaded == nil {
		loadAuthData()
	}
	return authLoaded
}

func loadAuthData() error {
	if _, err := os.Stat(authDataFile); os.IsNotExist(err) {
		authLoaded = &AuthData{
			Keys:        make(map[string]KeyPair),
			Connections: make(map[string]Connection),
		}
		// data, err := json.MarshalIndent(authDataDefault, "", "  ")
		// if err != nil {
		//   return err
		// }
		if _, err := authLoaded.NewKey("default"); err != nil {
			return err
		}
	}
	data, err := os.ReadFile(authDataFile)
	if err != nil {
		return err
	}
	decryptedData, err := hash.Decrypt(data, authEncryptionKey)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(decryptedData, &authLoaded); err != nil {
		return err
	}
	return nil
}

func saveAuthData() error {
	data, err := json.Marshal(authLoaded)
	if err != nil {
		return err
	}
	encryptedData, err := hash.Encrypt(data, authEncryptionKey)
	if err != nil {
		return err
	}
	if err := os.WriteFile(authDataFile, encryptedData, 0644); err != nil {
		return err
	}
	return nil
}

func (a *AuthData) NewKey(name string) (*KeyPair, error) {
	newKeyPair, err := GenerateKeyPair()
	if err != nil {
		return nil, err
	}
	return newKeyPair, addKey(name, newKeyPair)
}

func addKey(name string, key *KeyPair) error {
	authLoaded.Keys[name] = *key
	return saveAuthData()
}

func (a *AuthData) DeleteKey(name string) error {
	delete(a.Keys, name)
	return saveAuthData()
}

func (a *AuthData) AddConnection(host string, publicKey rsa.PublicKey) error {
	a.Connections[host] = Connection{PublicKey: publicKey, Host: host}
	return saveAuthData()
}

func (a *AuthData) DeleteConnection(host string) error {
	delete(a.Connections, host)
	return saveAuthData()
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
