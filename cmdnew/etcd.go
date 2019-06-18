package cmdnew

import (
	"context"
	"fmt"
	"time"

	"github.com/zdnscloud/zke/core"
	"github.com/zdnscloud/zke/core/pki"
	"github.com/zdnscloud/zke/pkg/hosts"
	"github.com/zdnscloud/zke/pkg/log"
	"github.com/zdnscloud/zke/types"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func EtcdCommand() cli.Command {
	snapshotFlags := []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "Specify snapshot name",
		},
		cli.StringFlag{
			Name:   "config",
			Usage:  "Specify an alternate cluster YAML file",
			Value:  pki.ClusterConfig,
			EnvVar: "ZKE_CONFIG",
		},
	}
	snapshotFlags = append(snapshotFlags, commonFlags...)

	return cli.Command{
		Name:  "etcd",
		Usage: "etcd snapshot save/restore operations in k8s cluster",
		Subcommands: []cli.Command{
			{
				Name:   "snapshot-save",
				Usage:  "Take snapshot on all etcd hosts",
				Flags:  snapshotFlags,
				Action: SnapshotSaveEtcdHostsFromCli,
			},
			{
				Name:   "snapshot-restore",
				Usage:  "Restore existing snapshot",
				Flags:  snapshotFlags,
				Action: RestoreEtcdSnapshotFromCli,
			},
		},
	}
}

func SnapshotSaveEtcdHosts(
	ctx context.Context,
	zkeConfig *types.ZKEConfig,
	dialersOptions hosts.DialersOptions,
	flags core.ExternalFlags, snapshotName string) error {
	log.Infof(ctx, "Starting saving snapshot on etcd hosts")
	kubeCluster, err := core.InitClusterObject(ctx, zkeConfig, flags)
	if err != nil {
		return err
	}
	if err := kubeCluster.SetupDialers(ctx, dialersOptions); err != nil {
		return err
	}
	if err := kubeCluster.TunnelHosts(ctx, flags); err != nil {
		return err
	}
	if err := kubeCluster.SnapshotEtcd(ctx, snapshotName); err != nil {
		return err
	}
	log.Infof(ctx, "Finished saving snapshot [%s] on all etcd hosts", snapshotName)
	return nil
}

func RestoreEtcdSnapshot(
	ctx context.Context,
	zkeConfig *types.ZKEConfig,
	dialersOptions hosts.DialersOptions,
	flags core.ExternalFlags, snapshotName string) error {
	log.Infof(ctx, "Restoring etcd snapshot %s", snapshotName)
	stateFilePath := core.GetStateFilePath(flags.ClusterFilePath, flags.ConfigDir)
	zkeFullState, err := core.ReadStateFile(ctx, stateFilePath)
	if err != nil {
		return err
	}
	zkeFullState.CurrentState = core.State{}
	if err := zkeFullState.WriteStateFile(ctx, stateFilePath); err != nil {
		return err
	}
	kubeCluster, err := core.InitClusterObject(ctx, zkeConfig, flags)
	if err != nil {
		return err
	}
	if err := kubeCluster.SetupDialers(ctx, dialersOptions); err != nil {
		return err
	}
	if err := kubeCluster.TunnelHosts(ctx, flags); err != nil {
		return err
	}
	// first download and check
	if err := kubeCluster.PrepareBackup(ctx, snapshotName); err != nil {
		return err
	}
	log.Infof(ctx, "Cleaning old kubernetes cluster")
	if err := kubeCluster.CleanupNodes(ctx); err != nil {
		return err
	}
	if err := kubeCluster.RestoreEtcdSnapshot(ctx, snapshotName); err != nil {
		return err
	}
	if err := ClusterInit(ctx, zkeConfig, dialersOptions, flags); err != nil {
		return err
	}
	if _, _, _, _, _, err := ClusterUp(ctx, dialersOptions, flags); err != nil {
		return err
	}
	if err := core.RestartClusterPods(ctx, kubeCluster); err != nil {
		return nil
	}
	if err := kubeCluster.RemoveOldNodes(ctx); err != nil {
		return err
	}
	log.Infof(ctx, "Finished restoring snapshot [%s] on all etcd hosts", snapshotName)
	return nil
}

func SnapshotSaveEtcdHostsFromCli(ctx *cli.Context) error {
	clusterFile, filePath, err := resolveClusterFile(ctx)
	if err != nil {
		return fmt.Errorf("failed to resolve cluster file: %v", err)
	}
	zkeConfig, err := core.ParseConfig(clusterFile)
	if err != nil {
		return fmt.Errorf("failed to parse cluster file: %v", err)
	}
	zkeConfig, err = setOptionsFromCLI(ctx, zkeConfig)
	if err != nil {
		return err
	}
	// Check snapshot name
	etcdSnapshotName := ctx.String("name")
	if etcdSnapshotName == "" {
		etcdSnapshotName = fmt.Sprintf("zke_etcd_snapshot_%s", time.Now().Format(time.RFC3339))
		logrus.Warnf("Name of the snapshot is not specified using [%s]", etcdSnapshotName)
	}
	// setting up the flags
	flags := core.GetExternalFlags(false, "", filePath)
	return SnapshotSaveEtcdHosts(context.Background(), zkeConfig, hosts.DialersOptions{}, flags, etcdSnapshotName)
}

func RestoreEtcdSnapshotFromCli(ctx *cli.Context) error {
	clusterFile, filePath, err := resolveClusterFile(ctx)
	if err != nil {
		return fmt.Errorf("failed to resolve cluster file: %v", err)
	}
	zkeConfig, err := core.ParseConfig(clusterFile)
	if err != nil {
		return fmt.Errorf("failed to parse cluster file: %v", err)
	}
	zkeConfig, err = setOptionsFromCLI(ctx, zkeConfig)
	if err != nil {
		return err
	}
	etcdSnapshotName := ctx.String("name")
	if etcdSnapshotName == "" {
		return fmt.Errorf("you must specify the snapshot name to restore")
	}
	// setting up the flags
	flags := core.GetExternalFlags(false, "", filePath)
	return RestoreEtcdSnapshot(context.Background(), zkeConfig, hosts.DialersOptions{}, flags, etcdSnapshotName)
}
