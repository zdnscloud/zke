package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/zdnscloud/zke/cluster"
	"github.com/zdnscloud/zke/pkg/log"
	"github.com/zdnscloud/zke/services"
	"github.com/zdnscloud/zke/storage/lvm"
	"github.com/zdnscloud/zke/storage/lvmd"
	"github.com/zdnscloud/zke/types"
	//"github.com/zdnscloud/zke/storage/nfs"
	"github.com/zdnscloud/gok8s/client"
	"github.com/zdnscloud/gok8s/client/config"
	"github.com/zdnscloud/zke/templates"
	corev1 "k8s.io/api/core/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"net"
	"strings"
	"time"
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

	CephStorageResourceName = "zke-storage-plugin-ceph"
	CephResourceName        = "ceph-storageclass"
	CephStorageClassName    = "ceph"

	StorageCSIAttacherImage     = "StorageCSIAttacherImage"
	StorageCSIProvisionerImage  = "StorageCSIProvisionerImage"
	StorageDriverRegistrarImage = "StorageDriverRegistrarImage"
	StorageCSILvmpluginImage    = "StorageCSILvmpluginImage"
	StorageLvmdImage            = "StorageLvmdImage"
	StorageNFSProvisionerImage  = "StorageNFSProvisionerImage"

	StorageTypeLabels = "storage.zcloud.cn/Storagetype"
)

func DeployStoragePlugin(ctx context.Context, c *cluster.Cluster) error {
	if err := doAddLabelsDeploy(ctx, c); err != nil {
		return err
	}

	if len(c.Storage.Lvm) > 0 {
		if err := doLVMDDeploy(ctx, c); err != nil {
			return err
		}
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
	/*
		if err := doCephDeploy(ctx, c); err != nil {
			return err
		}
			if err := doCephStorageDeploy(ctx, c); err != nil {
				return err
			}*/
	return nil
}

func doAddLabelsDeploy(ctx context.Context, c *cluster.Cluster) error {
	config, err := config.GetConfig()
	cli, err := client.New(config, client.Options{})
	if err != nil {
		fmt.Println(err)
	}
	var storageCfgMap = map[string][]types.Deviceconf{
		"Lvm":  c.Storage.Lvm,
		"Nfs":  c.Storage.Nfs,
		"Ceph": c.Storage.Ceph,
	}
	storagetypes := []string{"Lvm", "Nfs", "Ceph"}
	for _, t := range storagetypes {
		cfg, ok := storageCfgMap[t]
		if !ok || len(cfg) == 0 {
			return nil
		}
		for _, s := range cfg {
			fmt.Println("===========", s.Host)
			if err = doUpdateNode(cli, s.Host, t); err != nil {
				return err
			}
		}
	}
	return nil
}

func doUpdateNode(cli client.Client, name string, t string) error {
	node := corev1.Node{}
	err := cli.Get(context.TODO(), k8stypes.NamespacedName{"", name}, &node)
	if err != nil {
		return err
	}
	node.Labels[StorageTypeLabels] = t
	err = cli.Update(context.TODO(), &node)
	if err != nil {
		return err
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

func doCephDeploy(ctx context.Context, c *cluster.Cluster) error {
	log.Infof(ctx, "[storage] Setting up Ceph Cluster")
	var arr = make([]map[string]string, 0)
	for _, v := range c.Storage.Ceph {
		var m = make(map[string]string)
		m[Host] = v.Host
		m[Devs] = strings.Replace(strings.Trim(fmt.Sprint(v.Devs), "[]"), " ", " ", -1)
		arr = append(arr, m)
	}
	return nil
}

/*
func doCephStorageDeploy(ctx context.Context, c *cluster.Cluster) error {
	log.Infof(ctx, "[storage] Setting up StoragePlugin : %s", CephStorageClassName)
	var arr = make([]map[string]string, 0)
	for _, v := range c.Storage.Ceph {
		var m = make(map[string]string)
		m[Host] = v.Host
		m[Devs] = strings.Replace(strings.Trim(fmt.Sprint(v.Devs), "[]"), " ", " ", -1)
		arr = append(arr, m)
	}
	cephConfig := map[string]interface{}{
		CephList: arr,
	}
	lvmdYaml, err := templates.CompileTemplateFromMap(lvmd.LVMDTemplate, lvmdConfig)
	if err != nil {
		return err
	}
	if err := c.DoAddonDeploy(ctx, lvmdYaml, LVMDResourceName, true); err != nil {
		return err
	}
	return nil
}*/

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

/*
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
}*/

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
