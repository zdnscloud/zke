package ceph

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/zdnscloud/gok8s/client"
	"github.com/zdnscloud/gok8s/client/config"
	"github.com/zdnscloud/zke/core"
	"github.com/zdnscloud/zke/pkg/log"
	"github.com/zdnscloud/zke/pkg/templates"
	corev1 "k8s.io/api/core/v1"
	"strconv"
	"strings"
	"time"
)

const (
	StorageType        = "Ceph"
	StorageNamespace   = "zcloud"
	StorageClassName   = "cephfs"
	StorageHostLabels  = "storage.zcloud.cn/storagetype"
	CephMonSvcName     = "rook-ceph-mon-"
	CephOsdPodName     = "rook-ceph-osd-"
	CephMonSvcPort     = "6789"
	CephSecretName     = "rook-ceph-mon"
	CephSecretDataName = "admin-secret"
	CephAdminUser      = "admin"
	CephFilesystemName = "myfs"
	CheckInterval      = 6
	CephCheckTimes     = 50
)

func Deploy(ctx context.Context, c *core.Cluster) error {
	if len(c.Storage.Ceph) == 0 {
		return nil
	}
	if err := doCephCommonDeploy(ctx, c); err != nil {
		return err
	}
	if err := doCephClusterDeploy(ctx, c); err != nil {
		return err
	}
	if err := doWaitReady(ctx, c); err != nil {
		return err
	}
	if err := doCephFsDeploy(ctx, c); err != nil {
		return err
	}
	if err := doCephFsStorageDeploy(ctx, c); err != nil {
		return err
	}
	return nil
}

func doCephCommonDeploy(ctx context.Context, c *core.Cluster) error {
	log.Infof(ctx, "[storage] Setting up ceph common")
	cfg := map[string]interface{}{
		"RBACConfig":       c.Authorization.Mode,
		"StorageNamespace": StorageNamespace,
	}
	yaml, err := templates.CompileTemplateFromMap(commonTemplate, cfg)
	if err != nil {
		return err
	}
	if err := c.DoAddonDeploy(ctx, yaml, "zke-storage-ceph-common", true); err != nil {
		return err
	}
	return nil
}

func doCephClusterDeploy(ctx context.Context, c *core.Cluster) error {
	log.Infof(ctx, "[storage] Setting up ceph cluster")
	var arr = make([]map[string]interface{}, 0)
	for _, v := range c.Storage.Ceph {
		m := make(map[string]interface{})
		m["Host"] = v.Host
		devs := make([]map[string]string, 0)
		for _, dev := range v.Devs {
			n := make(map[string]string)
			n["Dev"] = dev[5:]
			devs = append(devs, n)
			m["Devs"] = devs
		}
		arr = append(arr, m)
	}
	cfg := map[string]interface{}{
		"CephList":                 arr,
		"RBACConfig":               c.Authorization.Mode,
		"StorageCephOperatorImage": c.SystemImages.StorageCephOperator,
		"StorageCephClusterImage":  c.SystemImages.StorageCephCluster,
		"StorageCephToolsImage":    c.SystemImages.StorageCephTools,
		"LabelKey":                 StorageHostLabels,
		"LabelValue":               StorageType,
		"StorageNamespace":         StorageNamespace,
	}
	yaml, err := templates.CompileTemplateFromMap(clusterTemplate, cfg)
	if err != nil {
		return err
	}
	if err := c.DoAddonDeploy(ctx, yaml, "zke-storage-ceph-cluster", true); err != nil {
		return err
	}
	return nil
}

func doCephFsDeploy(ctx context.Context, c *core.Cluster) error {
	log.Infof(ctx, "[storage] Setting up ceph filesystem")
	num := len(c.Storage.Ceph)
	cfg := map[string]interface{}{
		"CephFilesystem":   CephFilesystemName,
		"Replicas":         num,
		"StorageNamespace": StorageNamespace,
		"LabelKey":         StorageHostLabels,
		"LabelValue":       StorageType,
	}
	yaml, err := templates.CompileTemplateFromMap(filesystemTemplate, cfg)
	if err != nil {
		return err
	}
	if err := c.DoAddonDeploy(ctx, yaml, "zke-storage-cephfs", true); err != nil {
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
		"CephClusterMonitors":             monitors,
		"CephAdminUserEncode":             user,
		"CephAdminKeyEncode":              secret,
		"RBACConfig":                      c.Authorization.Mode,
		"StorageCephAttacherImage":        c.SystemImages.StorageCephAttacher,
		"StorageCephProvisionerImage":     c.SystemImages.StorageCephProvisioner,
		"StorageCephDriverRegistrarImage": c.SystemImages.StorageCephDriverRegistrar,
		"StorageCephFsCSIImage":           c.SystemImages.StorageCephFsCSI,
		"CephFilesystem":                  CephFilesystemName,
		"StorageNamespace":                StorageNamespace,
		"StorageClassName":                StorageClassName,
		"LabelKey":                        StorageHostLabels,
		"LabelValue":                      StorageType,
	}
	yaml, err := templates.CompileTemplateFromMap(fscsiTemplate, cfg)
	if err != nil {
		return err
	}
	if err := c.DoAddonDeploy(ctx, yaml, "zke-storage-cephfs-csi", true); err != nil {
		return err
	}
	return nil
}

func doWaitReady(ctx context.Context, c *core.Cluster) error {
	var ready bool
	var num int
	for _, v := range c.Storage.Ceph {
		num += len(v.Devs)
	}
	log.Infof(ctx, "[storage] Waiting for ceph cluster ready, it need %d osd proc to runing.", num)
	for i := 0; i < CephCheckTimes; i++ {
		if checkCephReady(ctx, c, num) {
			ready = true
			break
		}
		time.Sleep(time.Duration(CheckInterval) * time.Second)
	}
	if !ready {
		return errors.New("ceph cluster has not ready")
	}
	return nil
}

func getCephMonCfg(ctx context.Context, c *core.Cluster) (string, string, error) {
	config, err := config.GetConfigFromFile(c.LocalKubeConfigPath)
	cli, err := client.New(config, client.Options{})
	if err != nil {
		return "", "", err
	}
	services := corev1.ServiceList{}
	err = cli.List(context.TODO(), &client.ListOptions{Namespace: StorageNamespace}, &services)
	if err != nil {
		return "", "", err
	}
	secrets := corev1.SecretList{}
	err = cli.List(context.TODO(), &client.ListOptions{Namespace: StorageNamespace}, &secrets)
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

func checkCephReady(ctx context.Context, c *core.Cluster, num int) bool {
	config, err := config.GetConfigFromFile(c.LocalKubeConfigPath)
	cli, err := client.New(config, client.Options{})
	if err != nil {
		return false
	}
	pods := corev1.PodList{}
	err = cli.List(context.TODO(), &client.ListOptions{Namespace: StorageNamespace}, &pods)
	if err != nil {
		return false
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
		log.Infof(ctx, "[storage] %d osd proc has runing.", i+1)
	}
	return true
}
