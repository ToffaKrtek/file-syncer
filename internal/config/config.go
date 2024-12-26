package config

import (
	"encoding/json"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/ToffaKrtek/file-syncer/internal/hash"
)

type conf struct {
	Connections map[string]ServerConfig `json:"connections"`
	Server      ServerConfig            `json:"server"`
	TgToken     string                  `json:"tgToken"`
	SyncItems   []SyncItem              `json:"sync_items"`
}

var (
	configFile = "./config.json"
	configData *conf
	mu         sync.Mutex
)

func configRepository() *conf {
	mu.Lock()
	defer mu.Unlock()
	if configData == nil {
		var err error
		_, err = loadConfig()
		if err != nil {
			panic("Ошибка загрузки конфигурации")
		}
	}
	return configData
}

func loadConfig() (*conf, error) {
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		configData = &conf{
			Connections: make(map[string]ServerConfig),
			Server:      ServerConfig{},
			TgToken:     "", // TODO:: get from env package
			SyncItems:   []SyncItem{},
		}
		saveConfig(configData)
	}
	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, configData); err != nil {
		return nil, err
	}
	return configData, nil
}

func saveConfig(*conf) {
}

func (c conf) SetConnection(sc *ServerConfig) {
	c.Connections[sc.Name] = *sc
}

func (c conf) AddSyncItem(s *SyncItem) {
	c.SyncItems = append(
		c.SyncItems,
		*s,
	)
}

type ServerConfig struct {
	IpAddress string `json:"ip_address"`
	Name      string `json:"name"`
	Token     string `json:"token"`
}

func NewServer(ip string, name string, token string) *ServerConfig {
	return &ServerConfig{ip, name, token}
}

type SyncItem struct {
	Name   string `json:"name"`
	Source Item   `json:"source"`
	Target Item   `json:"target"`
}

type Item struct {
	Name  string `json:"name"`
	Path  string `json:"path"`
	Host  string `json:"host"`
	Hash  string `json:"hash"`
	IsDir bool   `json:"is_dir"`
}

func (i Item) CheckHash() (bool, error) {
	h, err := i.makeHash()
	check := false
	if err == nil {
		check = h == i.Hash
	}
	return check, err
}

func (i Item) SetHash() error {
	h, err := i.makeHash()
	if err == nil {
		i.Hash = h
	}
	return err
}

func (i Item) makeHash() (string, error) {
	return hash.Hash(i.Path, i.IsDir)
}

func NewItem(opts ...itemFunc) *Item {
	timestamp := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	item := &Item{
		timestamp,
		"./",
		"localhost",
		"",
		true,
	}
	for _, opt := range opts {
		opt(item)
	}
	return item
}

type itemFunc func(*Item)

func ItemName(name string) itemFunc {
	return func(i *Item) {
		i.Name = name
	}
}

func ItemPath(path string) itemFunc {
	return func(i *Item) {
		i.Path = path
	}
}

func ItemHost(host string) itemFunc {
	return func(i *Item) {
		i.Host = host
	}
}

func ItemIsDir(isdir bool) itemFunc {
	return func(i *Item) {
		i.IsDir = isdir
	}
}
