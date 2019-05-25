package nfs

import (
	"context"
	"errors"
	"github.com/zdnscloud/gok8s/client"
	"github.com/zdnscloud/gok8s/client/config"
	"github.com/zdnscloud/zke/core"
	"github.com/zdnscloud/zke/pkg/log"
	"github.com/zdnscloud/zke/pkg/templates"
	"github.com/zdnscloud/zke/storage/common"
)

const (
	StorageType       = "Nfs"
	StorageNamespace  = "zcloud"
	StorageClassName  = "nfs"
	StorageHostLabels = "storage.zcloud.cn/storagetype"
	CheckInterval     = 6
	NFS_DIR           = "/var/lib/singlecloud/nfs-export"
	NFSCheckTimes     = 10
)

func Deploy(ctx context.Context, c *core.Cluster) error {
	if len(c.Storage.Nfs) == 0 {
		return nil
	}
	if len(c.Storage.Nfs) > 1 {
		return errors.New("nfs only supports ont host!")
	}
	if err := doNFSInit(ctx, c); err != nil {
		return err
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

func doPartition(ctx context.Context, c *core.Cluster, name string) error {
	log.Infof(ctx, "[storage] Setting up nfs init")
	cfg := map[string]interface{}{
		"RBACConfig":          c.Authorization.Mode,
		"StorageNFSInitImage": c.SystemImages.StorageNFSInit,
		"LabelKey":            StorageHostLabels,
		"LabelValue":          StorageType,
		"StorageNamespace":    StorageNamespace,
		"NFS_DIR":             NFS_DIR,
	}
	yaml, err := templates.CompileTemplateFromMap(NFSInitTemplate, cfg)
	if err != nil {
		return err
	}
	if err := c.DoAddonDeploy(ctx, yaml, "zke-storage-nfs-init", true); err != nil {
		return err
	}
	return nil

}

func doNFSInit(ctx context.Context, c *core.Cluster) error {
	config, err := config.GetConfigFromFile(c.LocalKubeConfigPath)
	cli, err := client.New(config, client.Options{})
	if err != nil {
		return err
	}
	if !common.CheckStorageClassExist(cli, StorageClassName) {
		for _, h := range c.Storage.Nfs {
			err := doPartition(ctx, c, h.Host)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
