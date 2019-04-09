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
	Default                = "Default"
	LVMList                = "LVMList"
	Host                   = "Host"
	Devs                   = "Devs"
)

func (c *Cluster) deployStorageClass(ctx context.Context) error {
	if err := c.doLVMStorageclassDeploy(ctx); err != nil {
		return err
	}
	return nil
}

func (c *Cluster) doLVMStorageclassDeploy(ctx context.Context) error {
	log.Infof(ctx, "[storage] Setting up StorageClass: %s", LVMStorageClassName)
	var arr = make([]map[string]string, 0)
	for _, v := range c.Storage.Lvm {
		var m = make(map[string]string)
		m[Host] = v.Host
		m[Devs] = strings.Replace(strings.Trim(fmt.Sprint(v.Devs), "[]"), " ", ",", -1)
		arr = append(arr, m)
	}
	lvmstorageClassConfig := map[string]interface{}{
		RBACConfig: c.Authorization.Mode,
		LVMList:    arr,
	}
	storageYaml, err := templates.GetManifest(lvmstorageClassConfig, LVMResourceName)
	if err != nil {
		return err
	}
	if err := c.doAddonDeploy(ctx, storageYaml, LVMStorageResourceName, true); err != nil {
		return err
	}
	return nil
}
