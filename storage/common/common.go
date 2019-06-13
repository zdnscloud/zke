package common

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/zdnscloud/gok8s/client"
	"github.com/zdnscloud/gok8s/client/config"
	"github.com/zdnscloud/zke/core"
	"github.com/zdnscloud/zke/pkg/hosts"
	"github.com/zdnscloud/zke/pkg/log"
	"github.com/zdnscloud/zke/types"

	corev1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
)

const (
	StorageHostLabels        = "storage.zcloud.cn/storagetype"
	StorageBlocksAnnotations = "storage.zcloud.cn/blocks"
	StorageNamespace         = "zcloud"
)

func Prepara(ctx context.Context, c *core.Cluster, cfg []types.Deviceconf, storagetype string, classname string) error {
	log.Infof(ctx, "[storage] Check storage blocks and update nodes Labels and Taints for %s", storagetype)
	config, err := config.GetConfigFromFile(c.LocalKubeConfigPath)
	if err != nil {
		return err
	}
	cli, err := client.New(config, client.Options{})
	if err != nil {
		return err
	}
	for _, h := range cfg {
		if err = updateNode(cli, h.Host, storagetype, h.Devs); err != nil {
			return err
		}
		if err = checkStorageClassExist(cli, classname); err == nil {
			continue
		}
		if err = checkBlocks(ctx, c, h.Host, h.Devs); err != nil {
			return err
		}
	}
	return nil
}

func updateNode(cli client.Client, hostname string, storagetype string, devs []string) error {
	node := corev1.Node{}
	err := cli.Get(context.TODO(), k8stypes.NamespacedName{"", hostname}, &node)
	if err != nil {
		return err
	}
	node.Labels[StorageHostLabels] = storagetype
	node.Annotations[StorageBlocksAnnotations] = strings.Replace(strings.Trim(fmt.Sprint(devs), "[]"), " ", ",", -1)
	return cli.Update(context.TODO(), &node)
}

func checkStorageClassExist(cli client.Client, classname string) error {
	sc := storagev1.StorageClass{}
	return cli.Get(context.TODO(), k8stypes.NamespacedName{"", classname}, &sc)
}

func checkBlocks(ctx context.Context, c *core.Cluster, name string, devs []string) error {
	var node *hosts.Host
	allHosts := hosts.GetUniqueHostList(c.EtcdHosts, c.ControlPlaneHosts, c.WorkerHosts, c.StorageHosts, c.EdgeHosts)
	for _, n := range allHosts {
		if name == n.Address || name == n.HostnameOverride {
			node = n
			break
		}
	}
	client, err := node.GetSSHClient()
	if err != nil {
		return err
	}
	var errinfo string
	for _, d := range devs {
		cmd := "udevadm info --query=property " + d
		cmdout, cmderr, err := node.GetSSHCmdOutput(client, cmd)
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
