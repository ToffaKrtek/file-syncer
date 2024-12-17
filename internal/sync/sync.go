package sync

import "github.com/ToffaKrtek/file-syncer/internal/config"

type Upload struct{}

func (u Upload) Run(item config.Item) error {
	return nil
}

type Download struct{}

func (u Download) Run(item config.Item) error {
	return nil
}
