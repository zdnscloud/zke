package nfs

import (
	"context"
	"errors"
	"github.com/zdnscloud/zke/core"
	"github.com/zdnscloud/zke/pkg/log"
	"github.com/zdnscloud/zke/pkg/templates"
)

const (
	StorageType       = "Nfs"
	StorageNamespace  = "zcloud"
	StorageClassName  = "nfs"
	StorageHostLabels = "storage.zcloud.cn/storagetype"
	NFS_DIR           = "/var/lib/singlecloud/nfs-export"
)

func Deploy(ctx context.Context, c *core.Cluster) error {
	if len(c.Storage.Nfs) == 0 {
		return nil
	}
	if len(c.Storage.Nfs) > 1 {
		return errors.New("nfs only supports ont host!")
	}
	if err := doNFSStorageDeploy(ctx, c); err != nil {
		return err
	}
	return nil
}

func doNFSStorageDeploy(ctx context.Context, c *core.Cluster) error {
	log.Infof(ctx, "[storage] Setting up storageclass nfs")
	cfg := map[string]interface{}{
		"RBACConfig":                 c.Authorization.Mode,
		"StorageNFSProvisionerImage": c.SystemImages.StorageNFSProvisioner,
		"LabelKey":                   StorageHostLabels,
		"LabelValue":                 StorageType,
		"StorageClassName":           StorageClassName,
		"StorageNamespace":           StorageNamespace,
		"NFS_DIR":                    NFS_DIR,
		"StorageNFSInitImage":        c.SystemImages.StorageNFSInit,
	}
	yaml, err := templates.CompileTemplateFromMap(NFSStorageTemplate, cfg)
	if err != nil {
		return err
	}
	if err := c.DoAddonDeploy(ctx, yaml, "zke-storage-nfs", true); err != nil {
		return err
	}
	return nil
}
