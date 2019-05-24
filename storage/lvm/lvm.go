package lvm

import (
	"context"
	"errors"
	"github.com/zdnscloud/zke/core"
	"github.com/zdnscloud/zke/pkg/log"
	"github.com/zdnscloud/zke/pkg/templates"
	"net"
	"time"
)

const (
	StorageType       = "Lvm"
	StorageNamespace  = "zcloud"
	StorageClassName  = "lvm"
	StorageHostLabels = "storage.zcloud.cn/storagetype"
	CheckInterval     = 6
	LVMDCheckTimes    = 10
	LVMDPort          = "1736"
	LVMDProtocol      = "tcp"
)

func Deploy(ctx context.Context, c *core.Cluster) error {
	if len(c.Storage.Lvm) == 0 {
		return nil
	}
	if err := doLVMDDeploy(ctx, c); err != nil {
		return err
	}
	if err := doWaitReady(ctx, c); err != nil {
		return err
	}
	if err := doLVMStorageDeploy(ctx, c); err != nil {
		return err
	}
	return nil
}

func doLVMDDeploy(ctx context.Context, c *core.Cluster) error {
	log.Infof(ctx, "[storage] Setting up storage agent lvmd")
	cfg := map[string]interface{}{
		"RBACConfig":       c.Authorization.Mode,
		"StorageLvmdImage": c.SystemImages.StorageLvmd,
		"LabelKey":         StorageHostLabels,
		"LabelValue":       StorageType,
		"StorageNamespace": StorageNamespace,
	}
	yaml, err := templates.CompileTemplateFromMap(LVMDTemplate, cfg)
	if err != nil {
		return err
	}
	if err := c.DoAddonDeploy(ctx, yaml, "zke-storage-agent-lvmd", true); err != nil {
		return err
	}
	return nil
}

func doLVMStorageDeploy(ctx context.Context, c *core.Cluster) error {
	log.Infof(ctx, "[storage] Setting up storageclass lvm")
	cfg := map[string]interface{}{
		"RBACConfig":                     c.Authorization.Mode,
		"StorageLvmAttacherImage":        c.SystemImages.StorageLvmAttacher,
		"StorageLvmProvisionerImage":     c.SystemImages.StorageLvmProvisioner,
		"StorageLvmDriverRegistrarImage": c.SystemImages.StorageLvmDriverRegistrar,
		"StorageLvmCSIImage":             c.SystemImages.StorageLvmCSI,
		"LabelKey":                       StorageHostLabels,
		"LabelValue":                     StorageType,
		"StorageClassName":               StorageClassName,
		"StorageNamespace":               StorageNamespace,
	}
	yaml, err := templates.CompileTemplateFromMap(LVMStorageTemplate, cfg)
	if err != nil {
		return err
	}
	if err := c.DoAddonDeploy(ctx, yaml, "zke-storage-lvm", true); err != nil {
		return err
	}
	return nil
}

func doWaitReady(ctx context.Context, c *core.Cluster) error {
	log.Infof(ctx, "[storage] Waiting for lvmd ready")
	var ready bool
	for i := 0; i < LVMDCheckTimes; i++ {
		if checkLvmdReady(ctx, c) {
			ready = true
			break
		}
		time.Sleep(time.Duration(CheckInterval) * time.Second)
	}
	if !ready {
		return errors.New("some lvmd on node has not ready")
	}
	return nil
}

func checkLvmdReady(ctx context.Context, c *core.Cluster) bool {
	for _, h := range c.Storage.Lvm {
		for _, n := range c.Nodes {
			if h.Host == n.Address || h.Host == n.HostnameOverride {
				addr := n.Address + ":" + LVMDPort
				if _, err := net.Dial(LVMDProtocol, addr); err != nil {
					return false
				}
			}
		}
	}
	return true
}
