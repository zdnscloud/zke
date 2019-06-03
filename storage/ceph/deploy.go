package ceph

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/zdnscloud/gok8s/client"
	"github.com/zdnscloud/gok8s/client/config"
	"github.com/zdnscloud/zke/core"
	"github.com/zdnscloud/zke/pkg/k8s"
	"github.com/zdnscloud/zke/pkg/log"
	"github.com/zdnscloud/zke/storage/common"
	corev1 "k8s.io/api/core/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"strconv"
	"strings"
	"time"
)

func doCephCommonDeploy(ctx context.Context, c *core.Cluster) error {
	log.Infof(ctx, "[storage] Setting up ceph common")
	cli, err := k8s.GetK8sClientFromConfig("./kube_config_cluster.yml")
	if err != nil {
		return err
	}
	cfg := map[string]interface{}{
		"RBACConfig":       c.Authorization.Mode,
		"StorageNamespace": common.StorageNamespace,
	}
	return k8s.DoDeployFromTemplate(cli, commonTemplate, cfg)
}

func doCephClusterDeploy(ctx context.Context, c *core.Cluster) error {
	log.Infof(ctx, "[storage] Setting up ceph cluster")
	cli, err := k8s.GetK8sClientFromConfig("./kube_config_cluster.yml")
	if err != nil {
		return err
	}
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
		"LabelKey":                 common.StorageHostLabels,
		"LabelValue":               StorageType,
		"StorageNamespace":         common.StorageNamespace,
	}
	return k8s.DoDeployFromTemplate(cli, clusterTemplate, cfg)
}

func doCephFsDeploy(ctx context.Context, c *core.Cluster) error {
	log.Infof(ctx, "[storage] Setting up ceph filesystem")
	cli, err := k8s.GetK8sClientFromConfig("./kube_config_cluster.yml")
	if err != nil {
		return err
	}
	num := len(c.Storage.Ceph)
	cfg := map[string]interface{}{
		"CephFilesystem":   CephFilesystemName,
		"Replicas":         num,
		"StorageNamespace": common.StorageNamespace,
		"LabelKey":         common.StorageHostLabels,
		"LabelValue":       StorageType,
	}
	return k8s.DoDeployFromTemplate(cli, filesystemTemplate, cfg)
}

func doCephFsStorageDeploy(ctx context.Context, c *core.Cluster) error {
	log.Infof(ctx, "[storage] Setting up storageclass cephfs")
	cli, err := k8s.GetK8sClientFromConfig("./kube_config_cluster.yml")
	if err != nil {
		return err
	}
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
		"StorageNamespace":                common.StorageNamespace,
		"StorageClassName":                StorageClassName,
		"LabelKey":                        common.StorageHostLabels,
		"LabelValue":                      StorageType,
	}
	return k8s.DoDeployFromTemplate(cli, fscsiTemplate, cfg)
}

func doWaitReady(ctx context.Context, c *core.Cluster) error {
	var num int
	for _, v := range c.Storage.Ceph {
		num += len(v.Devs)
	}
	log.Infof(ctx, "[storage] Waiting for ceph cluster ready, it need %d osd proc to runing.", num)
	for i := 0; i < CephCheckTimes; i++ {
		ready, err := checkCephReady(ctx, c, num)
		if err != nil {
			return err
		}
		if ready {
			return nil
		}
		time.Sleep(time.Duration(CheckInterval) * time.Second)
	}
	return errors.New("Timeout. Ceph cluster has not ready")
}

func getCephMonCfg(ctx context.Context, c *core.Cluster) (string, string, error) {
	config, err := config.GetConfigFromFile(c.LocalKubeConfigPath)
	cli, err := client.New(config, client.Options{})
	if err != nil {
		return "", "", err
	}
	sec := corev1.Secret{}
	err = cli.Get(context.TODO(), k8stypes.NamespacedName{common.StorageNamespace, CephSecretName}, &sec)
	if err != nil {
		return "", "", err
	}
	secret := base64.StdEncoding.EncodeToString(sec.Data[CephSecretDataName])

	services := corev1.ServiceList{}
	err = cli.List(context.TODO(), &client.ListOptions{Namespace: common.StorageNamespace}, &services)
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
	monitors := strings.Replace(strings.Trim(fmt.Sprint(addrs), "[]"), " ", ",", -1)
	return monitors, secret, nil
}

func checkCephReady(ctx context.Context, c *core.Cluster, num int) (bool, error) {
	config, err := config.GetConfigFromFile(c.LocalKubeConfigPath)
	cli, err := client.New(config, client.Options{})
	if err != nil {
		return false, err
	}
	pods := corev1.PodList{}
	err = cli.List(context.TODO(), &client.ListOptions{Namespace: common.StorageNamespace}, &pods)
	if err != nil {
		return false, err
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
			return false, nil
		}
		log.Infof(ctx, "[storage] %d osd proc has runing.", i+1)
	}
	return true, nil
}
