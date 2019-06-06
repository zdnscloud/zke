package nfs

import (
	"context"
	"errors"
	"github.com/zdnscloud/zke/core"
	"github.com/zdnscloud/zke/storage/common"
)

const (
	StorageType      = "Nfs"
	StorageClassName = "nfs"
	NFS_DIR          = "/var/lib/singlecloud/nfs-export"
)

type Nfs struct{}

func (s *Nfs) Up(ctx context.Context, c *core.Cluster) error {
	if len(c.Storage.Nfs) == 0 {
		return nil
	}
	if len(c.Storage.Nfs) > 1 {
		return errors.New("nfs only supports one host!")
	}
	if err := common.Prepara(ctx, c, c.Storage.Nfs, StorageType, StorageClassName); err != nil {
		return err
	}
	if err := doNFSStorageDeploy(ctx, c); err != nil {
		return err
	}
	return nil
}

func (s *Nfs) Remove(ctx context.Context, c *core.Cluster) error {
	return nil
}
