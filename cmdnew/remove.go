package cmdnew

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/zdnscloud/zke/core"
	"github.com/zdnscloud/zke/core/pki"
	"github.com/zdnscloud/zke/pkg/hosts"
	"github.com/zdnscloud/zke/pkg/log"
	"github.com/zdnscloud/zke/types"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func RemoveCommand() cli.Command {
	removeFlags := []cli.Flag{
		cli.StringFlag{
			Name:   "config",
			Usage:  "Specify an alternate cluster YAML file",
			Value:  pki.ClusterConfig,
			EnvVar: "ZKE_CONFIG",
		},
		cli.BoolFlag{
			Name:  "force",
			Usage: "Force removal of the cluster",
		},
	}
	removeFlags = append(removeFlags, commonFlags...)
	return cli.Command{
		Name:   "remove",
		Usage:  "Teardown the cluster and clean cluster nodes",
		Action: clusterRemoveFromCli,
		Flags:  removeFlags,
	}
}

func ClusterRemove(
	ctx context.Context,
	zkeConfig *types.ZKEConfig,
	dialersOptions hosts.DialersOptions,
	flags core.ExternalFlags) error {
	log.Infof(ctx, "Tearing down Kubernetes cluster")
	kubeCluster, err := core.InitClusterObject(ctx, zkeConfig, flags)
	if err != nil {
		return err
	}
	if err := kubeCluster.SetupDialers(ctx, dialersOptions); err != nil {
		return err
	}
	err = kubeCluster.TunnelHosts(ctx, flags)
	if err != nil {
		return err
	}
	logrus.Debugf("Starting Cluster removal")
	err = kubeCluster.ClusterRemove(ctx)
	if err != nil {
		return err
	}
	log.Infof(ctx, "Cluster removed successfully")
	return nil
}

func clusterRemoveFromCli(ctx *cli.Context) error {
	clusterFile, filePath, err := resolveClusterFile(ctx)
	if err != nil {
		return fmt.Errorf("Failed to resolve cluster file: %v", err)
	}
	force := ctx.Bool("force")
	if !force {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("Are you sure you want to remove Kubernetes cluster [y/n]: ")
		input, err := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if err != nil {
			return err
		}
		if input != "y" && input != "Y" {
			return nil
		}
	}
	zkeConfig, err := core.ParseConfig(clusterFile)
	if err != nil {
		return fmt.Errorf("Failed to parse cluster file: %v", err)
	}
	err = validateConfigVersion(zkeConfig)
	if err != nil {
		return err
	}
	zkeConfig, err = setOptionsFromCLI(ctx, zkeConfig)
	if err != nil {
		return err
	}
	// setting up the flags
	flags := core.GetExternalFlags(false, "", filePath)
	return ClusterRemove(context.Background(), zkeConfig, hosts.DialersOptions{}, flags)
}
