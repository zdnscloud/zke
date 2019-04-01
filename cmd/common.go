package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/urfave/cli"
	"github.com/zdnscloud/zke/cluster"
	"github.com/zdnscloud/zke/hosts"
	"github.com/zdnscloud/zke/log"
	"github.com/zdnscloud/zke/types"
)

var commonFlags = []cli.Flag{
	cli.BoolFlag{
		Name:  "ssh-agent-auth",
		Usage: "Use SSH Agent Auth defined by SSH_AUTH_SOCK",
	},
	cli.BoolFlag{
		Name:  "ignore-docker-version",
		Usage: "Disable Docker version check",
	},
}

func resolveClusterFile(ctx *cli.Context) (string, string, error) {
	clusterFile := ctx.String("config")
	fp, err := filepath.Abs(clusterFile)
	if err != nil {
		return "", "", fmt.Errorf("failed to lookup current directory name: %v", err)
	}
	file, err := os.Open(fp)
	if err != nil {
		return "", "", fmt.Errorf("can not find cluster configuration file: %v", err)
	}
	defer file.Close()
	buf, err := ioutil.ReadAll(file)
	if err != nil {
		return "", "", fmt.Errorf("failed to read file: %v", err)
	}
	clusterFileBuff := string(buf)
	return clusterFileBuff, clusterFile, nil
}

func setOptionsFromCLI(c *cli.Context, rkeConfig *types.ZcloudKubernetesEngineConfig) (*types.ZcloudKubernetesEngineConfig, error) {
	// If true... override the file.. else let file value go through
	if c.Bool("ssh-agent-auth") {
		rkeConfig.SSHAgentAuth = c.Bool("ssh-agent-auth")
	}

	if c.Bool("ignore-docker-version") {
		rkeConfig.IgnoreDockerVersion = c.Bool("ignore-docker-version")
	}

	if c.Bool("s3") {
		if rkeConfig.Services.Etcd.BackupConfig == nil {
			rkeConfig.Services.Etcd.BackupConfig = &types.BackupConfig{}
		}
		rkeConfig.Services.Etcd.BackupConfig.S3BackupConfig = setS3OptionsFromCLI(c)
	}
	return rkeConfig, nil
}

func ClusterInit(ctx context.Context, rkeConfig *types.ZcloudKubernetesEngineConfig, dialersOptions hosts.DialersOptions, flags cluster.ExternalFlags) error {
	log.Infof(ctx, "Initiating Kubernetes cluster")
	var fullState *cluster.FullState
	stateFilePath := cluster.GetStateFilePath(flags.ClusterFilePath, flags.ConfigDir)
	if len(flags.CertificateDir) == 0 {
		flags.CertificateDir = cluster.GetCertificateDirPath(flags.ClusterFilePath, flags.ConfigDir)
	}
	rkeFullState, _ := cluster.ReadStateFile(ctx, stateFilePath)

	kubeCluster, err := cluster.InitClusterObject(ctx, rkeConfig, flags)
	if err != nil {
		return err
	}

	if err := kubeCluster.SetupDialers(ctx, dialersOptions); err != nil {
		return err
	}

	err = doUpgradeLegacyCluster(ctx, kubeCluster, rkeFullState)
	if err != nil {
		log.Warnf(ctx, "[state] can't fetch legacy cluster state from Kubernetes")
	}
	fullState, err = cluster.RebuildState(ctx, &kubeCluster.ZcloudKubernetesEngineConfig, rkeFullState, flags)
	if err != nil {
		return err
	}

	rkeState := cluster.FullState{
		DesiredState: fullState.DesiredState,
		CurrentState: fullState.CurrentState,
	}
	return rkeState.WriteStateFile(ctx, stateFilePath)
}

func setS3OptionsFromCLI(c *cli.Context) *types.S3BackupConfig {
	endpoint := c.String("s3-endpoint")
	bucketName := c.String("bucket-name")
	region := c.String("region")
	accessKey := c.String("access-key")
	secretKey := c.String("secret-key")
	var s3BackupBackend = &types.S3BackupConfig{}
	if len(endpoint) != 0 {
		s3BackupBackend.Endpoint = endpoint
	}
	if len(bucketName) != 0 {
		s3BackupBackend.BucketName = bucketName
	}
	if len(region) != 0 {
		s3BackupBackend.Region = region
	}
	if len(accessKey) != 0 {
		s3BackupBackend.AccessKey = accessKey
	}
	if len(secretKey) != 0 {
		s3BackupBackend.SecretKey = secretKey
	}
	return s3BackupBackend
}
