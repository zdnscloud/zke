package cmd

import (
	"fmt"

	"github.com/zdnscloud/zke/core"
	"github.com/zdnscloud/zke/core/pki"
	"github.com/zdnscloud/zke/types"

	"github.com/urfave/cli"
)

const defaultConfigVersion = "v1.0.7"

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
	serverVersion, err := core.GetK8sVersion(localKubeConfig, nil)
	if err != nil {
		return err
	}
	fmt.Printf("Server Version: %s\n", serverVersion)
	return nil
}

func validateConfigVersion(zkeConfig *types.ZcloudKubernetesEngineConfig) error {
	if zkeConfig.ConfigVersion != defaultConfigVersion {
		return fmt.Errorf("config version not match[new version is %s, and current config file version is %s], please execut config command to generate new config", defaultConfigVersion, zkeConfig.ConfigVersion)
	}
	return nil
}
