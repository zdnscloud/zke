package storage

import (
	"context"
	"fmt"
	"github.com/zdnscloud/zke/cluster"
	"github.com/zdnscloud/zke/pkg/log"
	"github.com/zdnscloud/zke/storage/lvm"
	"github.com/zdnscloud/zke/storage/nfs"
	"github.com/zdnscloud/zke/templates"
	"strings"
)

const (
	RBACConfig = "RBACConfig"

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
		if err := doLVMStorageDeploy(ctx, c); err != nil {
			return err
		}
	}
	if c.Storage.NFS.Size > 0 {
		if err := doNFSStorageDeploy(ctx, c); err != nil {
			return err
		}
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
		StorageLvmdImage:            c.SystemImages.StorageLvmd,
		LVMList:                     arr,
		NodeSelector:                cluster.StorageRoleLabel,
	}
	lvmstorageYaml, err := templates.CompileTemplateFromMap(lvm.LVMStorageTemplate, lvmstorageConfig)
	//lvmstorageYaml, err := templates.GetManifest(lvmstorageConfig, LVMResourceName)
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
	//nfsstorageYaml, err := templates.GetManifest(nfsstorageConfig, NFSResourceName)
	nfsstorageYaml, err := templates.CompileTemplateFromMap(nfs.NFSStorageTemplate, nfsstorageConfig)
	if err != nil {
		return err
	}
	if err := c.DoAddonDeploy(ctx, nfsstorageYaml, NFSStorageResourceName, true); err != nil {
		return err
	}
	return nil
}
