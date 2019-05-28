package ceph

import (
	"context"
	"github.com/zdnscloud/zke/core"
	"github.com/zdnscloud/zke/storage/common"
)

const (
	StorageType        = "Ceph"
	StorageClassName   = "cephfs"
	CephMonSvcName     = "rook-ceph-mon-"
	CephOsdPodName     = "rook-ceph-osd-"
	CephMonSvcPort     = "6789"
	CephSecretName     = "rook-ceph-mon"
	CephSecretDataName = "admin-secret"
	CephAdminUser      = "admin"
	CephFilesystemName = "myfs"
	CheckInterval      = 6
	CephCheckTimes     = 50
)

type Ceph struct{}

func (s *Ceph) Up(ctx context.Context, c *core.Cluster) error {
	if len(c.Storage.Ceph) == 0 {
		return nil
	}
	if err := common.Prepara(ctx, c, c.Storage.Ceph, StorageType, StorageClassName); err != nil {
		return err
	}
	if err := doCephCommonDeploy(ctx, c); err != nil {
		return err
	}
	if err := doCephClusterDeploy(ctx, c); err != nil {
		return err
	}
	if err := doWaitReady(ctx, c); err != nil {
		return err
	}
	if err := doCephFsDeploy(ctx, c); err != nil {
		return err
	}
	if err := doCephFsStorageDeploy(ctx, c); err != nil {
		return err
	}
	return nil
}

func (s *Ceph) Remove(ctx context.Context, c *core.Cluster) error {
	return nil
}
