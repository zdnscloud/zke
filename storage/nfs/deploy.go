package nfs

import (
	"context"
	"github.com/zdnscloud/zke/core"
	"github.com/zdnscloud/zke/pkg/k8s"
	"github.com/zdnscloud/zke/pkg/log"
	"github.com/zdnscloud/zke/storage/common"
)

func doNFSStorageDeploy(ctx context.Context, c *core.Cluster) error {
	log.Infof(ctx, "[storage] Setting up storageclass nfs")
	cli, err := k8s.GetK8sClientFromConfig(c.LocalKubeConfigPath)
	if err != nil {
		return err
	}
	cfg := map[string]interface{}{
		"RBACConfig":                 c.Authorization.Mode,
		"StorageNFSProvisionerImage": c.SystemImages.StorageNFSProvisioner,
		"LabelKey":                   common.StorageHostLabels,
		"LabelValue":                 StorageType,
		"StorageClassName":           StorageClassName,
		"StorageNamespace":           common.StorageNamespace,
		"NFS_DIR":                    NFS_DIR,
		"StorageNFSInitImage":        c.SystemImages.StorageNFSInit,
	}
	return k8s.DoDeployFromTemplate(cli, NFSStorageTemplate, cfg)
}
