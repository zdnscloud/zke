package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/zdnscloud/gok8s/client"
	"github.com/zdnscloud/gok8s/client/config"
	"github.com/zdnscloud/zke/core"
	"github.com/zdnscloud/zke/pkg/log"
	"github.com/zdnscloud/zke/storage/ceph"
	"github.com/zdnscloud/zke/storage/common"
	"github.com/zdnscloud/zke/storage/lvm"
	"github.com/zdnscloud/zke/storage/nfs"
	"github.com/zdnscloud/zke/types"
	corev1 "k8s.io/api/core/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"strings"
)

const (
	StorageHostLabels        = "storage.zcloud.cn/storagetype"
	StorageBlocksAnnotations = "storage.zcloud.cn/blocks"
)

var storageTypes []string = []string{"Lvm", "Nfs", "Ceph"}

var storageClassMap = map[string]string{
	"Lvm":  lvm.StorageClassName,
	"Nfs":  nfs.StorageClassName,
	"Ceph": ceph.StorageClassName,
}

func DeployStoragePlugin(ctx context.Context, c *core.Cluster) error {
	if err := doPreparaJob(ctx, c); err != nil {
		return err
	}
	if err := nfs.Deploy(ctx, c); err != nil {
		return err
	}
	if err := lvm.Deploy(ctx, c); err != nil {
		return err
	}
	if err := ceph.Deploy(ctx, c); err != nil {
		return err
	}
	return nil
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

	for _, t := range storageTypes {
		cfg, ok := storageCfgMap[t]
		if !ok || len(cfg) == 0 {
			continue
		}
		for _, h := range cfg {
			if err = doUpdateNode(cli, h.Host, t, h.Devs); err != nil {
				return err
			}
			if common.CheckStorageClassExist(cli, storageClassMap[t]) {
				continue
			}
			if err = doCheckBlocks(ctx, c, h.Host, h.Devs); err != nil {
				return err
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

func doCheckBlocks(ctx context.Context, c *core.Cluster, name string, devs []string) error {
	var node types.ZKEConfigNode
	for _, n := range c.Nodes {
		if name == n.Address || name == n.HostnameOverride {
			node = n
		}
	}
	client, err := common.MakeSSHClient(node)
	if err != nil {
		return err
	}
	var errinfo string
	for _, d := range devs {
		cmd := "udevadm info --query=property " + d
		cmdout, cmderr, err := common.GetSSHCmdOut(client, cmd)
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
