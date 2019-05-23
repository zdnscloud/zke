package storage

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/zdnscloud/gok8s/client"
	"github.com/zdnscloud/gok8s/client/config"
	"github.com/zdnscloud/zke/core"
	"github.com/zdnscloud/zke/pkg/hosts"
	"github.com/zdnscloud/zke/pkg/log"
	"github.com/zdnscloud/zke/pkg/templates"
	"github.com/zdnscloud/zke/storage/ceph"
	"github.com/zdnscloud/zke/storage/lvm"
	"github.com/zdnscloud/zke/storage/lvmd"
	"github.com/zdnscloud/zke/storage/nfs"
	"github.com/zdnscloud/zke/types"
	"golang.org/x/crypto/ssh"
	corev1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"net"
	"strconv"
	"strings"
	"time"
)

var (
	ErrLvmdNotReady = errors.New("some lvmd on node has not ready")
)

const (
	RBACConfig = "RBACConfig"

	StorageNamespace = "kube-storage"

	LVMDResourceName          = "zke-storage-lvmd"
	LVMStorageResourceName    = "zke-storage-lvm"
	NFSStorageResourceName    = "zke-storage-nfs"
	NFSInitResourceName       = "zke-storage-nfs-init"
	CephClusterResourceName   = "zke-storage-ceph-cluster"
	CephFsResourceName        = "zke-storage-cephfs"
	CephFsStorageResourceName = "zke-storage-cephfs-csi"

	LVMStorageType  = "Lvm"
	NFSStorageType  = "Nfs"
	CephStorageType = "Ceph"

	CephNamespace       = "rook-ceph"
	CephMonSvcName      = "rook-ceph-mon-"
	CephOsdPodName      = "rook-ceph-osd-"
	CephMonSvcPort      = "6789"
	CephSecretName      = "rook-ceph-mon"
	CephSecretDataName  = "admin-secret"
	CephAdminUser       = "admin"
	CephFilesystemName  = "myfs"
	CephClusterMonitors = "CephClusterMonitors"
	CephClusterSecret   = "CephClusterSecret"
	CephAdminUserEncode = "CephAdminUserEncode"
	CephAdminKeyEncode  = "CephAdminKeyEncode"
	CephFilesystem      = "CephFilesystem"
	Replicas            = "Replicas"

	LVMD         = "lvmd"
	LVMDPort     = "1736"
	LVMDProtocol = "tcp"

	CheckInterval  = 6
	LVMDCheckTimes = 10
	CephCheckTimes = 50
	NFSCheckTimes  = 10

	LVMList  = "LVMList"
	NFSList  = "NFSList"
	CephList = "CephList"
	Host     = "Host"
	Devs     = "Devs"
	Dev      = "Dev"

	LabelKey   = "LabelKey"
	LabelValue = "LabelValue"

	StorageLvmAttacherImage        = "StorageLvmAttacherImage"
	StorageLvmProvisionerImage     = "StorageLvmProvisionerImage"
	StorageLvmDriverRegistrarImage = "StorageLvmDriverRegistrarImage"
	StorageLvmCSIImage             = "StorageLvmCSIImage"
	StorageLvmdImage               = "StorageLvmdImage"

	StorageNFSProvisionerImage = "StorageNFSProvisionerImage"
	StorageNFSInitImage        = "StorageNFSInitImage"

	StorageCephOperatorImage        = "StorageCephOperatorImage"
	StorageCephClusterImage         = "StorageCephClusterImage"
	StorageCephToolsImage           = "StorageCephToolsImage"
	StorageCephAttacherImage        = "StorageCephAttacherImage"
	StorageCephProvisionerImage     = "StorageCephProvisionerImage"
	StorageCephDriverRegistrarImage = "StorageCephDriverRegistrarImage"
	StorageCephFsCSIImage           = "StorageCephFsCSIImage"

	StorageHostLabels        = "storage.zcloud.cn/storagetype"
	StorageTypeAnnotations   = "storage.zcloud.cn/storagetype"
	StorageBlocksAnnotations = "storage.zcloud.cn/blocks"
)

func DeployStoragePlugin(ctx context.Context, c *core.Cluster) error {
	if err := doPreparaJob(ctx, c); err != nil {
		return err
	}
	if err := doNFSDeploy(ctx, c); err != nil {
		return err
	}
	if err := doLVMDeploy(ctx, c); err != nil {
		return err
	}
	if err := doCephDeploy(ctx, c); err != nil {
		return err
	}
	return nil
}

var storageClassMap = map[string]string{
	"Lvm":  "lvm",
	"Nfs":  "nfs",
	"Ceph": "cephfs",
}

func doPreparaJob(ctx context.Context, c *core.Cluster) error {
	log.Infof(ctx, "[storage] Check storage blocks and update nodes Labels and Taints ")
	config, err := config.GetConfigFromFile("kube_config_cluster.yml")
	cli, err := client.New(config, client.Options{})
	if err != nil {
		return err
	}
	var storageCfgMap = map[string][]types.Deviceconf{
		"Lvm":  c.Storage.Lvm,
		"Nfs":  c.Storage.Nfs,
		"Ceph": c.Storage.Ceph,
	}
	storagetypes := []string{LVMStorageType, NFSStorageType, CephStorageType}
	for _, t := range storagetypes {
		cfg, ok := storageCfgMap[t]
		if !ok || len(cfg) == 0 {
			continue
		}
		for _, s := range cfg {
			if err = doUpdateNode(cli, s.Host, t, s.Devs); err != nil {
				return err
			}
			if !checkStorageClassExist(cli, storageClassMap[t]) {
				if err = doCheckBlocks(ctx, c, s.Host, s.Devs); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func doUpdateNode(cli client.Client, name string, t string, devs []string) error {
	node := corev1.Node{}
	err := cli.Get(context.TODO(), k8stypes.NamespacedName{"", name}, &node)
	if err != nil {
		return err
	}
	annotations := strings.Replace(strings.Trim(fmt.Sprint(devs), "[]"), " ", ",", -1)
	node.Labels[StorageHostLabels] = t
	node.Annotations[StorageBlocksAnnotations] = annotations
	err = cli.Update(context.TODO(), &node)
	if err != nil {
		return err
	}
	return nil
}

func checkStorageClassExist(cli client.Client, name string) bool {
	var exist bool
	scs := storagev1.StorageClassList{}
	err := cli.List(context.TODO(), nil, &scs)
	if err != nil {
		return exist
	}
	for _, s := range scs.Items {
		if s.Name == name {
			exist = true
			break
		}
	}
	return exist
}

func doCheckBlocks(ctx context.Context, c *core.Cluster, name string, devs []string) error {
	var node types.ZKEConfigNode
	for _, n := range c.Nodes {
		if name == n.Address || name == n.HostnameOverride {
			node = n
		}
	}
	client, err := makeSSHClient(node)
	if err != nil {
		return err
	}
	var errinfo string
	for _, d := range devs {
		cmd := "udevadm info --query=property " + d
		cmdout, cmderr, err := getSSHCmdOut(client, cmd)
		if err != nil {
			return err
		}
		if cmderr != "" || strings.Contains(cmdout, "ID_PART_TABLE") || strings.Contains(cmdout, "ID_FS_TYPE") {
			info := name + ":" + d + "."
			errinfo += info
		}
	}
	if errinfo != "" {
		return errors.New("some blocks cat not be used!" + errinfo)
	}
	return nil
}

func doLVMDeploy(ctx context.Context, c *core.Cluster) error {
	if len(c.Storage.Lvm) == 0 {
		return nil
	}
	if err := doLVMDDeploy(ctx, c); err != nil {
		return err
	}
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
	if err := doLVMStorageDeploy(ctx, c); err != nil {
		return err
	}
	return nil
}

func doLVMDDeploy(ctx context.Context, c *core.Cluster) error {
	log.Infof(ctx, "[storage] Setting up lvmd")
	cfg := map[string]interface{}{
		RBACConfig:       c.Authorization.Mode,
		StorageLvmdImage: c.SystemImages.StorageLvmd,
		LabelKey:         StorageHostLabels,
		LabelValue:       LVMStorageType,
	}
	yaml, err := templates.CompileTemplateFromMap(lvmd.LVMDTemplate, cfg)
	if err != nil {
		return err
	}
	if err := c.DoAddonDeploy(ctx, yaml, LVMDResourceName, true); err != nil {
		return err
	}
	return nil
}

func doLVMStorageDeploy(ctx context.Context, c *core.Cluster) error {
	log.Infof(ctx, "[storage] Setting up storageclass lvm")
	cfg := map[string]interface{}{
		RBACConfig:                     c.Authorization.Mode,
		StorageLvmAttacherImage:        c.SystemImages.StorageLvmAttacher,
		StorageLvmProvisionerImage:     c.SystemImages.StorageLvmProvisioner,
		StorageLvmDriverRegistrarImage: c.SystemImages.StorageLvmDriverRegistrar,
		StorageLvmCSIImage:             c.SystemImages.StorageLvmCSI,
		LabelKey:                       StorageHostLabels,
		LabelValue:                     LVMStorageType,
	}
	yaml, err := templates.CompileTemplateFromMap(lvm.LVMStorageTemplate, cfg)
	if err != nil {
		return err
	}
	if err := c.DoAddonDeploy(ctx, yaml, LVMStorageResourceName, true); err != nil {
		return err
	}
	return nil
}

func doCephDeploy(ctx context.Context, c *core.Cluster) error {
	if len(c.Storage.Ceph) == 0 {
		return nil
	}
	if err := doCephClusterDeploy(ctx, c); err != nil {
		return err
	}
	log.Infof(ctx, "[storage] Waiting for ceph cluster ready")
	var ready bool
	for i := 0; i < CephCheckTimes; i++ {
		if checkCephReady(ctx, c) {
			ready = true
			break
		}
		time.Sleep(time.Duration(CheckInterval) * time.Second)
	}
	if !ready {
		return errors.New("ceph cluster has not ready")
	}
	if err := doCephFsDeploy(ctx, c); err != nil {
		return err
	}
	if err := doCephFsStorageDeploy(ctx, c); err != nil {
		return err
	}
	return nil
}

func doCephClusterDeploy(ctx context.Context, c *core.Cluster) error {
	log.Infof(ctx, "[storage] Setting up ceph cluster")
	var arr = make([]map[string]interface{}, 0)
	for _, v := range c.Storage.Ceph {
		m := make(map[string]interface{})
		m[Host] = v.Host
		devs := make([]map[string]string, 0)
		for _, dev := range v.Devs {
			n := make(map[string]string)
			n[Dev] = dev[5:]
			devs = append(devs, n)
			m[Devs] = devs
		}
		arr = append(arr, m)
	}
	cfg := map[string]interface{}{
		CephList:                 arr,
		RBACConfig:               c.Authorization.Mode,
		StorageCephOperatorImage: c.SystemImages.StorageCephOperator,
		StorageCephClusterImage:  c.SystemImages.StorageCephCluster,
		StorageCephToolsImage:    c.SystemImages.StorageCephTools,
		LabelKey:                 StorageHostLabels,
		LabelValue:               CephStorageType,
	}
	yaml, err := templates.CompileTemplateFromMap(ceph.ClusterTemplate, cfg)
	if err != nil {
		return err
	}
	if err := c.DoAddonDeploy(ctx, yaml, CephClusterResourceName, true); err != nil {
		return err
	}
	return nil
}

func doCephFsDeploy(ctx context.Context, c *core.Cluster) error {
	log.Infof(ctx, "[storage] Setting up ceph filesystem")
	num := len(c.Storage.Ceph)
	cfg := map[string]interface{}{
		CephFilesystem: CephFilesystemName,
		Replicas:       num,
	}
	yaml, err := templates.CompileTemplateFromMap(ceph.FilesystemTemplate, cfg)
	if err != nil {
		return err
	}
	if err := c.DoAddonDeploy(ctx, yaml, CephFsResourceName, true); err != nil {
		return err
	}
	return nil
}

func doCephFsStorageDeploy(ctx context.Context, c *core.Cluster) error {
	log.Infof(ctx, "[storage] Setting up storageclass cephfs")
	monitors, secret, err := getCephMonCfg(ctx, c)
	if err != nil {
		return err
	}
	user := base64.StdEncoding.EncodeToString([]byte(CephAdminUser))
	cfg := map[string]interface{}{
		CephClusterMonitors:             monitors,
		CephAdminUserEncode:             user,
		CephAdminKeyEncode:              secret,
		RBACConfig:                      c.Authorization.Mode,
		StorageCephAttacherImage:        c.SystemImages.StorageCephAttacher,
		StorageCephProvisionerImage:     c.SystemImages.StorageCephProvisioner,
		StorageCephDriverRegistrarImage: c.SystemImages.StorageCephDriverRegistrar,
		StorageCephFsCSIImage:           c.SystemImages.StorageCephFsCSI,
		CephFilesystem:                  CephFilesystemName,
	}
	yaml, err := templates.CompileTemplateFromMap(ceph.FscsiTemplate, cfg)
	if err != nil {
		return err
	}
	if err := c.DoAddonDeploy(ctx, yaml, CephFsStorageResourceName, true); err != nil {
		return err
	}
	return nil
}

func getCephMonCfg(ctx context.Context, c *core.Cluster) (string, string, error) {
	config, err := config.GetConfigFromFile("kube_config_cluster.yml")
	cli, err := client.New(config, client.Options{})
	if err != nil {
		return "", "", err
	}
	services := corev1.ServiceList{}
	err = cli.List(context.TODO(), &client.ListOptions{Namespace: CephNamespace}, &services)
	if err != nil {
		return "", "", err
	}
	secrets := corev1.SecretList{}
	err = cli.List(context.TODO(), &client.ListOptions{Namespace: CephNamespace}, &secrets)
	if err != nil {
		return "", "", err
	}
	var addrs []string
	for _, sv := range services.Items {
		if strings.Contains(sv.Name, CephMonSvcName) {
			addr := sv.Spec.ClusterIP + ":" + CephMonSvcPort
			addrs = append(addrs, addr)
		}
	}
	var monitors, secret string
	for _, sc := range secrets.Items {
		if sc.Name == CephSecretName {
			secret = base64.StdEncoding.EncodeToString(sc.Data[CephSecretDataName])
		}
	}
	monitors = strings.Replace(strings.Trim(fmt.Sprint(addrs), "[]"), " ", ",", -1)
	return monitors, secret, nil
}

func doNFSDeploy(ctx context.Context, c *core.Cluster) error {
	if len(c.Storage.Nfs) == 0 {
		return nil
	}
	if len(c.Storage.Nfs) > 1 {
		return errors.New("nfs only supports ont host!")
	}
	config, err := config.GetConfigFromFile("kube_config_cluster.yml")
	cli, err := client.New(config, client.Options{})
	if err != nil {
		return err
	}
	if !checkStorageClassExist(cli, storageClassMap["Nfs"]) {
		for _, h := range c.Storage.Nfs {
			name := h.Host
			err := doNfsInit(ctx, c, name)
			if err != nil {
				return err
			}
			ready, err := doNfsMount(ctx, c, name)
			if !ready || err != nil {
				return err
			}
		}
	}
	if err := doNFSStorageDeploy(ctx, c); err != nil {
		return err
	}
	return nil
}

func doNFSStorageDeploy(ctx context.Context, c *core.Cluster) error {
	log.Infof(ctx, "[storage] Setting up storageclass nfs")
	cfg := map[string]interface{}{
		RBACConfig:                 c.Authorization.Mode,
		StorageNFSProvisionerImage: c.SystemImages.StorageNFSProvisioner,
		LabelKey:                   StorageHostLabels,
		LabelValue:                 NFSStorageType,
	}
	yaml, err := templates.CompileTemplateFromMap(nfs.NFSStorageTemplate, cfg)
	if err != nil {
		return err
	}
	if err := c.DoAddonDeploy(ctx, yaml, NFSStorageResourceName, true); err != nil {
		return err
	}
	return nil
}

func checkLvmdReady(ctx context.Context, c *core.Cluster) bool {
	for _, h := range c.Storage.Lvm {
		for _, n := range c.Nodes {
			if h.Host == n.Address || h.Host == n.HostnameOverride {
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

func checkCephReady(ctx context.Context, c *core.Cluster) bool {
	config, err := config.GetConfigFromFile("kube_config_cluster.yml")
	cli, err := client.New(config, client.Options{})
	if err != nil {
		return false
	}
	pods := corev1.PodList{}
	err = cli.List(context.TODO(), &client.ListOptions{Namespace: CephNamespace}, &pods)
	if err != nil {
		return false
	}
	var num int
	for _, v := range c.Storage.Ceph {
		num += len(v.Devs)
	}
	for i := 0; i < num; i++ {
		name := CephOsdPodName + strconv.Itoa(i)
		pod := corev1.Pod{}
		for _, p := range pods.Items {
			if strings.Contains(p.Name, name) {
				pod = p
				break
			}
		}
		if pod.Status.Phase != "Running" {
			return false
		}
	}
	return true
}

func doNfsInit(ctx context.Context, c *core.Cluster, name string) error {
	cfg := map[string]interface{}{
		RBACConfig:          c.Authorization.Mode,
		StorageNFSInitImage: c.SystemImages.StorageNFSInit,
		LabelKey:            StorageHostLabels,
		LabelValue:          NFSStorageType,
	}
	yaml, err := templates.CompileTemplateFromMap(nfs.NFSInitTemplate, cfg)
	if err != nil {
		return err
	}
	if err := c.DoAddonDeploy(ctx, yaml, NFSInitResourceName, true); err != nil {
		return err
	}
	return nil

}
func doNfsMount(ctx context.Context, c *core.Cluster, name string) (bool, error) {
	var node types.ZKEConfigNode
	for _, n := range c.Nodes {
		if name == n.Address || name == n.HostnameOverride {
			node = n
		}
	}
	client, err := makeSSHClient(node)
	if err != nil {
		return false, err
	}

	cmd := `ls /dev/mapper|grep -E nfs-data -q;if [ $? -eq 0 ];then echo true;else echo false;fi`
	var ready bool
	for i := 0; i < NFSCheckTimes; i++ {
		cmdout, _, err := getSSHCmdOut(client, cmd)
		if err != nil {
			return false, err
		}
		if cmdout == "true" {
			ready = true
			break
		}
		time.Sleep(time.Duration(CheckInterval) * time.Second)
	}
	if ready {
		cmd := `sudo mkdir /var/lib/singlecloud/nfs-export -p;sleep 5;sudo mount /dev/mapper/nfs-data /var/lib/singlecloud/nfs-export;`
		cmdout, cmderr, err := getSSHCmdOut(client, cmd)
		if err != nil || cmdout != "" || cmderr != "" {
			return false, errors.New("mount host path for nfs failed!")
		}
	}
	return ready, nil
}

func makeSSHClient(node types.ZKEConfigNode) (*ssh.Client, error) {
	var sshKeyString, sshCertString string
	if !node.SSHAgentAuth {
		var err error
		sshKeyString, err = hosts.PrivateKeyPath(node.SSHKeyPath)
		if err != nil {
			return nil, err
		}

		if len(node.SSHCertPath) > 0 {
			sshCertString, err = hosts.CertificatePath(node.SSHCertPath)
			if err != nil {
				return nil, err
			}
		}
	}
	cfg, err := hosts.GetSSHConfig(node.User, sshKeyString, sshCertString, node.SSHAgentAuth)
	if err != nil {
		return nil, err
	}
	addr := node.Address + ":22"
	return ssh.Dial("tcp", addr, cfg)
}

func getSSHCmdOut(client *ssh.Client, cmd string) (string, string, error) {
	var cmdout, cmderr string
	session, err := client.NewSession()
	if err != nil {
		return cmdout, "error", err
	}
	defer session.Close()
	var stdOut, stdErr bytes.Buffer
	session.Stdout = &stdOut
	session.Stderr = &stdErr
	session.Run(cmd)
	cmdout = strings.Replace(stdOut.String(), "\n", "", -1)
	cmderr = strings.Replace(stdErr.String(), "\n", "", -1)
	return cmdout, cmderr, nil
}
