package storage

import (
	"context"
	"github.com/zdnscloud/zke/core"
	"github.com/zdnscloud/zke/storage/ceph"
	"github.com/zdnscloud/zke/storage/lvm"
	"github.com/zdnscloud/zke/storage/nfs"
)

type Storage interface {
	Up(ctx context.Context, c *core.Cluster) error
	Remove(ctx context.Context, c *core.Cluster) error
}

func DeployStoragePlugin(ctx context.Context, c *core.Cluster) error {
	storages := []Storage{
		&lvm.Lvm{},
		&nfs.Nfs{},
		&ceph.Ceph{},
	}
	for _, s := range storages {
		if err := s.Up(ctx, c); err != nil {
			return err
		}
	}
	return nil
}
