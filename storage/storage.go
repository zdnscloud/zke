package storage

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/zdnscloud/zke/cluster"
	"github.com/zdnscloud/zke/cluster/services"
	"github.com/zdnscloud/zke/pkg/log"
	"github.com/zdnscloud/zke/pkg/templates"
	"github.com/zdnscloud/zke/storage/lvm"
	"github.com/zdnscloud/zke/storage/lvmd"
	"github.com/zdnscloud/zke/storage/nfs"
)

var (
	ErrLvmdNotReady = errors.New("some lvmd on node has not ready")
)

const (
	RBACConfig = "RBACConfig"

	LVMD              = "lvmd"
	LVMDPort          = "1736"
	LVMDProtocol      = "tcp"
	LVMDResourceName  = "zke-storage-agent-lvmd"
	LVMDCheckTimes    = 10
	LVMDCheckInterval = 6

	LVMStorageResourceName = "zke-storage-plugin-lvm"
	LVMResourceName        = "lvm-storageclass"
	LVMStorageClassName    = "lvm"
	LVMList                = "LVMList"
	Host                   = "Host"
	Devs                   = "Devs"
	NodeSelector           = "NodeSelector"

	NFSStorageResourceName = "zke-storage-plugin-nfs"
	NFSResourceName        = "nfs-storageclass"
	NFSStorageClassName    = "nfs"
	Size                   = "Size"

	StorageCSIAttacherImage     = "StorageCSIAttacherImage"
	StorageCSIProvisionerImage  = "StorageCSIProvisionerImage"
	StorageDriverRegistrarImage = "StorageDriverRegistrarImage"
	StorageCSILvmpluginImage    = "StorageCSILvmpluginImage"
	StorageLvmdImage            = "StorageLvmdImage"
	StorageNFSProvisionerImage  = "StorageNFSProvisionerImage"
)

func DeployStoragePlugin(ctx context.Context, c *cluster.Cluster) error {
	if len(c.Storage.Lvm) > 0 {
		if err := doLVMDDeploy(ctx, c); err != nil {
			return err
		}
		var ready bool
		for i := 0; i < LVMDCheckTimes; i++ {
			if checkLvmdReady(ctx, c) {
				ready = true
				break
			}
			time.Sleep(time.Duration(LVMDCheckInterval) * time.Second)
		}
		if !ready {
			return ErrLvmdNotReady
		}
		if err := doLVMStorageDeploy(ctx, c); err != nil {
			return err
		}
		if c.Storage.NFS.Size > 0 {
			if err := doNFSStorageDeploy(ctx, c); err != nil {
				return err
			}
		}
	}
	return nil
}

func doLVMDDeploy(ctx context.Context, c *cluster.Cluster) error {
	log.Infof(ctx, "[storage] Setting up StorageAgent: %s", LVMD)
	var arr = make([]map[string]string, 0)
	for _, v := range c.Storage.Lvm {
		var m = make(map[string]string)
		m[Host] = v.Host
		m[Devs] = strings.Replace(strings.Trim(fmt.Sprint(v.Devs), "[]"), " ", " ", -1)
		arr = append(arr, m)
	}
	lvmdConfig := map[string]interface{}{
		LVMList:          arr,
		StorageLvmdImage: c.SystemImages.StorageLvmd,
	}
	lvmdYaml, err := templates.CompileTemplateFromMap(lvmd.LVMDTemplate, lvmdConfig)
	if err != nil {
		return err
	}
	if err := c.DoAddonDeploy(ctx, lvmdYaml, LVMDResourceName, true); err != nil {
		return err
	}
	return nil
}

func doLVMStorageDeploy(ctx context.Context, c *cluster.Cluster) error {
	log.Infof(ctx, "[storage] Setting up StoragePlugin : %s", LVMStorageClassName)
	var arr = make([]map[string]string, 0)
	for _, v := range c.Storage.Lvm {
		var m = make(map[string]string)
		m[Host] = v.Host
		m[Devs] = strings.Replace(strings.Trim(fmt.Sprint(v.Devs), "[]"), " ", " ", -1)
		arr = append(arr, m)
	}
	lvmstorageConfig := map[string]interface{}{
		RBACConfig:                  c.Authorization.Mode,
		StorageCSIAttacherImage:     c.SystemImages.StorageCSIAttacher,
		StorageCSIProvisionerImage:  c.SystemImages.StorageCSIProvisioner,
		StorageDriverRegistrarImage: c.SystemImages.StorageDriverRegistrar,
		StorageCSILvmpluginImage:    c.SystemImages.StorageCSILvmplugin,
		NodeSelector:                cluster.StorageRoleLabel,
	}
	lvmstorageYaml, err := templates.CompileTemplateFromMap(lvm.LVMStorageTemplate, lvmstorageConfig)
	if err != nil {
		return err
	}
	if err := c.DoAddonDeploy(ctx, lvmstorageYaml, LVMStorageResourceName, true); err != nil {
		return err
	}
	return nil
}

func doNFSStorageDeploy(ctx context.Context, c *cluster.Cluster) error {
	log.Infof(ctx, "[storage] Setting up StoragePlugin : %s", NFSStorageClassName)
	nfsstorageConfig := map[string]interface{}{
		RBACConfig:                 c.Authorization.Mode,
		StorageNFSProvisionerImage: c.SystemImages.StorageNFSProvisioner,
		Size:                       c.Storage.NFS.Size,
	}
	nfsstorageYaml, err := templates.CompileTemplateFromMap(nfs.NFSStorageTemplate, nfsstorageConfig)
	if err != nil {
		return err
	}
	if err := c.DoAddonDeploy(ctx, nfsstorageYaml, NFSStorageResourceName, true); err != nil {
		return err
	}
	return nil
}

func checkLvmdReady(ctx context.Context, c *cluster.Cluster) bool {
	for _, n := range c.Nodes {
		for _, v := range n.Role {
			if v == services.StorageRole {
				addr := n.Address + ":" + LVMDPort
				_, err := net.Dial(LVMDProtocol, addr)
				if err != nil {
					return false
				}
			}
		}
	}
	return true
}
