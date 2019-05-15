package cmd

import (
	"fmt"

	"github.com/zdnscloud/zke/cluster"
	"github.com/zdnscloud/zke/pki"

	"github.com/urfave/cli"
)

func VersionCommand() cli.Command {
	versionFlags := []cli.Flag{
		cli.StringFlag{
			Name:   "config",
			Usage:  "Specify an alternate cluster YAML file",
			Value:  pki.ClusterConfig,
			EnvVar: "ZKE_CONFIG",
		},
	}
	return cli.Command{
		Name:   "version",
		Usage:  "Show cluster Kubernetes version",
		Action: getClusterVersion,
		Flags:  versionFlags,
	}
}

func getClusterVersion(ctx *cli.Context) error {
	localKubeConfig := pki.GetLocalKubeConfig(ctx.String("config"), "")
	// not going to use a k8s dialer here.. this is a CLI command
	serverVersion, err := cluster.GetK8sVersion(localKubeConfig, nil)
	if err != nil {
		return err
	}
	fmt.Printf("Server Version: %s\n", serverVersion)
	return nil
}
