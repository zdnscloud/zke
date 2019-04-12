package cluster

import (
	"context"
	"fmt"
	"github.com/zdnscloud/zke/log"
	"github.com/zdnscloud/zke/templates"
	"strings"
)

const (
	LVMStorageResourceName = "zke-storage-plugin-lvm"
	LVMResourceName        = "lvm-storageclass"
	LVMStorageClassName    = "lvm"
	LVMList                = "LVMList"
	Host                   = "Host"
	Devs                   = "Devs"
	Storage                = "node-role.kubernetes.io/storage"
	NodeSelector           = "NodeSelector"

	StorageCSIAttacherImage     = "StorageCSIAttacherImage"
	StorageCSIProvisionerImage  = "StorageCSIProvisionerImage"
	StorageDriverRegistrarImage = "StorageDriverRegistrarImage"
	StorageCSILvmpluginImage    = "StorageCSILvmpluginImage"
	StorageLvmdImage            = "StorageLvmdImage"
)

func (c *Cluster) deployStoragePlugin(ctx context.Context) error {
	if len(c.Storage.Lvm) > 0 {
		if err := c.doLVMStorageDeploy(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (c *Cluster) doLVMStorageDeploy(ctx context.Context) error {
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
		NodeSelector:                Storage,
	}
	storageYaml, err := templates.GetManifest(lvmstorageConfig, LVMResourceName)
	if err != nil {
		return err
	}
	if err := c.doAddonDeploy(ctx, storageYaml, LVMStorageResourceName, true); err != nil {
		return err
	}
	return nil
}
