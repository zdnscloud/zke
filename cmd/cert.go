package cmd

import (
	"context"
	"fmt"

	"github.com/urfave/cli"
	"github.com/zdnscloud/zke/cluster"
	"github.com/zdnscloud/zke/log"
	"github.com/zdnscloud/zke/pki"
	"github.com/zdnscloud/zke/types"
)

func CertificateCommand() cli.Command {
	return cli.Command{
		Name:  "cert",
		Usage: "Certificates management for RKE cluster",
		Subcommands: cli.Commands{
			cli.Command{
				Name:   "generate-csr",
				Usage:  "Generate certificate sign requests for k8s components",
				Action: generateCSRFromCli,
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:   "config",
						Usage:  "Specify an alternate cluster YAML file",
						Value:  pki.ClusterConfig,
						EnvVar: "ZKE_CONFIG",
					},
					cli.StringFlag{
						Name:  "cert-dir",
						Usage: "Specify a certificate dir path",
					},
				},
			},
		},
	}
}

func generateCSRFromCli(ctx *cli.Context) error {
	clusterFile, filePath, err := resolveClusterFile(ctx)
	if err != nil {
		return fmt.Errorf("Failed to resolve cluster file: %v", err)
	}
	zkeConfig, err := cluster.ParseConfig(clusterFile)
	if err != nil {
		return fmt.Errorf("Failed to parse cluster file: %v", err)
	}
	zkeConfig, err = setOptionsFromCLI(ctx, zkeConfig)
	if err != nil {
		return err
	}
	// setting up the flags
	externalFlags := cluster.GetExternalFlags(false, "", filePath)
	externalFlags.CertificateDir = ctx.String("cert-dir")
	externalFlags.CustomCerts = ctx.Bool("custom-certs")
	return GenerateRKECSRs(context.Background(), zkeConfig, externalFlags)
}

func GenerateRKECSRs(ctx context.Context, zkeConfig *types.ZcloudKubernetesEngineConfig, flags cluster.ExternalFlags) error {
	log.Infof(ctx, "Generating Kubernetes cluster CSR certificates")
	if len(flags.CertificateDir) == 0 {
		flags.CertificateDir = cluster.GetCertificateDirPath(flags.ClusterFilePath, flags.ConfigDir)
	}
	certBundle, err := pki.ReadCSRsAndKeysFromDir(flags.CertificateDir)
	if err != nil {
		return err
	}
	// initialze the cluster object from the config file
	kubeCluster, err := cluster.InitClusterObject(ctx, zkeConfig, flags)
	if err != nil {
		return err
	}
	// Generating csrs for kubernetes components
	if err := pki.GenerateRKEServicesCSRs(ctx, certBundle, kubeCluster.ZcloudKubernetesEngineConfig); err != nil {
		return err
	}
	return pki.WriteCertificates(kubeCluster.CertificateDir, certBundle)
}
