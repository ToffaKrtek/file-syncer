package config

import (
	"strconv"
	"time"

	"github.com/ToffaKrtek/file-syncer/hash"
)

type conf struct {
	Connections map[string]ServerConfig `json:"connections"`
	Server      ServerConfig            `json:"server"`
	TgToken     string                  `json:"tgToken"`
	SyncItems   []SyncItem              `json:"sync_items"`
}

var configData *conf

func Config() *conf {
	if configData == nil {
		// TODO:: init
	}
	return configData
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
