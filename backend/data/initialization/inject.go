package initialization

import (
	"github.com/DMwangnima/easy-disk/data/api"
	"github.com/DMwangnima/easy-disk/data/storage/local"
)

func initInject() error {
	var err error
	api.Store, err = local.NewStorageChunk(LocalConfig.Storage.BasePath, DefaultBlockSize, LocalConfig.Storage.ChunkSize, LocalConfig.Storage.BlockNum)
	if err != nil {
		return err
	}
	return nil
}
