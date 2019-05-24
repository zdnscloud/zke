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
	"github.com/zdnscloud/zke/types"
	"time"
)

const (
	StorageType       = "Nfs"
	StorageNamespace  = "zcloud"
	StorageClassName  = "nfs"
	StorageHostLabels = "storage.zcloud.cn/storagetype"
	CheckInterval     = 6
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
			err = doNfsMount(ctx, c, h.Host)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func doNfsMount(ctx context.Context, c *core.Cluster, name string) error {
	var node types.ZKEConfigNode
	for _, n := range c.Nodes {
		if name == n.Address || name == n.HostnameOverride {
			node = n
		}
	}
	client, err := common.MakeSSHClient(node)
	if err != nil {
		return err
	}

	cmd := `ls /dev/mapper|grep -E nfs-data -q;if [ $? -eq 0 ];then echo true;else echo false;fi`
	var ready bool
	for i := 0; i < NFSCheckTimes; i++ {
		cmdout, _, err := common.GetSSHCmdOut(client, cmd)
		if err != nil {
			return err
		}
		if cmdout == "true" {
			ready = true
			break
		}
		time.Sleep(time.Duration(CheckInterval) * time.Second)
	}
	if ready {
		cmd := `sudo mkdir /var/lib/singlecloud/nfs-export -p;sleep 5;sudo mount /dev/mapper/nfs-data /var/lib/singlecloud/nfs-export;`
		cmdout, cmderr, err := common.GetSSHCmdOut(client, cmd)
		if err != nil || cmdout != "" || cmderr != "" {
			return errors.New("mount host path for nfs failed!")
		}
	}
	return nil
}
