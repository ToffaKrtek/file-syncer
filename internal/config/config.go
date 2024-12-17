package config

type conf struct {
	Connections map[string]ServerConfig `json:"connections"`
	Server      ServerConfig            `json:"server"`
	TgToken     string                  `json:"tgToken"`
	SyncItems   []SyncItem              `json:"sync_items"`
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
