package lvm

import (
	"context"
	"errors"
	"github.com/zdnscloud/zke/core"
	"github.com/zdnscloud/zke/pkg/log"
	"github.com/zdnscloud/zke/pkg/templates"
	"github.com/zdnscloud/zke/storage/common"
	"net"
	"time"
)

func doLVMDDeploy(ctx context.Context, c *core.Cluster) error {
	log.Infof(ctx, "[storage] Setting up storage agent lvmd")
	cfg := map[string]interface{}{
		"RBACConfig":       c.Authorization.Mode,
		"StorageLvmdImage": c.SystemImages.StorageLvmd,
		"LabelKey":         common.StorageHostLabels,
		"LabelValue":       StorageType,
		"StorageNamespace": common.StorageNamespace,
	}
	yaml, err := templates.CompileTemplateFromMap(LVMDTemplate, cfg)
	if err != nil {
		return err
	}
	return c.DoAddonDeploy(ctx, yaml, "zke-storage-agent-lvmd", true)
}

func doLVMStorageDeploy(ctx context.Context, c *core.Cluster) error {
	log.Infof(ctx, "[storage] Setting up storageclass lvm")
	cfg := map[string]interface{}{
		"RBACConfig":                     c.Authorization.Mode,
		"StorageLvmAttacherImage":        c.SystemImages.StorageLvmAttacher,
		"StorageLvmProvisionerImage":     c.SystemImages.StorageLvmProvisioner,
		"StorageLvmDriverRegistrarImage": c.SystemImages.StorageLvmDriverRegistrar,
		"StorageLvmCSIImage":             c.SystemImages.StorageLvmCSI,
		"LabelKey":                       common.StorageHostLabels,
		"LabelValue":                     StorageType,
		"StorageClassName":               StorageClassName,
		"StorageNamespace":               common.StorageNamespace,
	}
	yaml, err := templates.CompileTemplateFromMap(LVMStorageTemplate, cfg)
	if err != nil {
		return err
	}
	return c.DoAddonDeploy(ctx, yaml, "zke-storage-lvm", true)
}

func doWaitReady(ctx context.Context, c *core.Cluster) error {
	log.Infof(ctx, "[storage] Waiting for lvmd ready")
	for i := 0; i < LVMDCheckTimes; i++ {
		if checkLvmdReady(ctx, c) {
			return nil
		}
		time.Sleep(time.Duration(CheckInterval) * time.Second)
	}
	return errors.New("some lvmd on node has not ready")
}

func checkLvmdReady(ctx context.Context, c *core.Cluster) bool {
	for _, h := range c.Storage.Lvm {
		for _, n := range c.Nodes {
			if h.Host == n.Address || h.Host == n.HostnameOverride {
				addr := n.Address + ":" + LVMDPort
				if _, err := net.Dial(LVMDProtocol, addr); err != nil {
					log.Infof(ctx, "[storage] lvmd on %s not ready! Please check", n.Address)
					return false
				}
				log.Infof(ctx, "[storage] lvmd on %s has ready", n.Address)
				break
			}
		}
	}
	return true
}
