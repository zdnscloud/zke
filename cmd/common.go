package cmd

import (
	"context"
	"fmt"
	"github.com/urfave/cli"
	"github.com/zdnscloud/zke/cluster"
	"github.com/zdnscloud/zke/hosts"
	"github.com/zdnscloud/zke/pkg/log"
	"github.com/zdnscloud/zke/types"
	"io/ioutil"
	"os"
	"path/filepath"
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

func setOptionsFromCLI(c *cli.Context, zkeConfig *types.ZcloudKubernetesEngineConfig) (*types.ZcloudKubernetesEngineConfig, error) {
	// If true... override the file.. else let file value go through
	if c.Bool("ssh-agent-auth") {
		zkeConfig.SSHAgentAuth = c.Bool("ssh-agent-auth")
	}
	if c.Bool("ignore-docker-version") {
		zkeConfig.IgnoreDockerVersion = c.Bool("ignore-docker-version")
	}
	return zkeConfig, nil
}

func ClusterInit(ctx context.Context, zkeConfig *types.ZcloudKubernetesEngineConfig, dialersOptions hosts.DialersOptions, flags cluster.ExternalFlags) error {
	log.Infof(ctx, "Initiating Kubernetes cluster")
	var fullState *cluster.FullState
	stateFilePath := cluster.GetStateFilePath(flags.ClusterFilePath, flags.ConfigDir)
	if len(flags.CertificateDir) == 0 {
		flags.CertificateDir = cluster.GetCertificateDirPath(flags.ClusterFilePath, flags.ConfigDir)
	}
	zkeFullState, _ := cluster.ReadStateFile(ctx, stateFilePath)
	kubeCluster, err := cluster.InitClusterObject(ctx, zkeConfig, flags)
	if err != nil {
		return err
	}
	if err := kubeCluster.SetupDialers(ctx, dialersOptions); err != nil {
		return err
	}
	err = doUpgradeLegacyCluster(ctx, kubeCluster, zkeFullState)
	if err != nil {
		log.Warnf(ctx, "[state] can't fetch legacy cluster state from Kubernetes")
	}
	fullState, err = cluster.RebuildState(ctx, &kubeCluster.ZcloudKubernetesEngineConfig, zkeFullState, flags)
	if err != nil {
		return err
	}
	zkeState := cluster.FullState{
		DesiredState: fullState.DesiredState,
		CurrentState: fullState.CurrentState,
	}
	return zkeState.WriteStateFile(ctx, stateFilePath)
}
