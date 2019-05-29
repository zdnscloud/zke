package lvm

import (
	"context"
	"github.com/zdnscloud/zke/core"
	"github.com/zdnscloud/zke/storage/common"
)

const (
	StorageType      = "Lvm"
	StorageClassName = "lvm"
	CheckInterval    = 6
	LVMDCheckTimes   = 50
	LVMDPort         = "1736"
	LVMDProtocol     = "tcp"
)

type Lvm struct{}

func (s *Lvm) Up(ctx context.Context, c *core.Cluster) error {
	if len(c.Storage.Lvm) == 0 {
		return nil
	}
	if err := common.Prepara(ctx, c, c.Storage.Lvm, StorageType, StorageClassName); err != nil {
		return err
	}
	if err := doLVMDDeploy(ctx, c); err != nil {
		return err
	}
	if err := doLVMStorageDeploy(ctx, c); err != nil {
		return err
	}
	return nil
}

func (s *Lvm) Remove(ctx context.Context, c *core.Cluster) error {
	return nil
}
