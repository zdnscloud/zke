package cmd

import (
	"context"
	"fmt"

	"github.com/rancher/rke/cluster"
	"github.com/rancher/rke/hosts"
	"github.com/rancher/rke/log"
	"github.com/rancher/rke/pki"
	"github.com/rancher/rke/services"
	"github.com/rancher/types/apis/management.cattle.io/v3"
	"github.com/urfave/cli"
	"k8s.io/client-go/util/cert"
)

func CertificateCommand() cli.Command {
	return cli.Command{
		Name:  "cert",
		Usage: "Certificates management for RKE cluster",
		Subcommands: cli.Commands{
			cli.Command{
				Name:   "rotate",
				Usage:  "Rotate RKE cluster certificates",
				Action: rotateRKECertificatesFromCli,
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:   "config",
						Usage:  "Specify an alternate cluster YAML file",
						Value:  pki.ClusterConfig,
						EnvVar: "RKE_CONFIG",
					},
					cli.StringSliceFlag{
						Name: "service",
						Usage: fmt.Sprintf("Specify a k8s service to rotate certs, (allowed values: %s, %s, %s, %s, %s, %s)",
							services.KubeAPIContainerName,
							services.KubeControllerContainerName,
							services.SchedulerContainerName,
							services.KubeletContainerName,
							services.KubeproxyContainerName,
							services.EtcdContainerName,
						),
					},
					cli.BoolFlag{
						Name:  "rotate-ca",
						Usage: "Rotate all certificates including CA certs",
					},
				},
			},
			cli.Command{
				Name:   "generate-csr",
				Usage:  "Generate certificate sign requests for k8s components",
				Action: generateCSRFromCli,
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:   "config",
						Usage:  "Specify an alternate cluster YAML file",
						Value:  pki.ClusterConfig,
						EnvVar: "RKE_CONFIG",
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

func rotateRKECertificatesFromCli(ctx *cli.Context) error {
	k8sComponents := ctx.StringSlice("service")
	rotateCACerts := ctx.Bool("rotate-ca")
	clusterFile, filePath, err := resolveClusterFile(ctx)
	if err != nil {
		return fmt.Errorf("Failed to resolve cluster file: %v", err)
	}

	rkeConfig, err := cluster.ParseConfig(clusterFile)
	if err != nil {
		return fmt.Errorf("Failed to parse cluster file: %v", err)
	}
	rkeConfig, err = setOptionsFromCLI(ctx, rkeConfig)
	if err != nil {
		return err
	}
	// setting up the flags
	externalFlags := cluster.GetExternalFlags(false, false, false, "", filePath)
	// setting up rotate flags
	rkeConfig.RotateCertificates = &v3.RotateCertificates{
		CACertificates: rotateCACerts,
		Services:       k8sComponents,
	}
	if err := ClusterInit(context.Background(), rkeConfig, hosts.DialersOptions{}, externalFlags); err != nil {
		return err
	}
	_, _, _, _, _, err = ClusterUp(context.Background(), hosts.DialersOptions{}, externalFlags)
	return err
}

func generateCSRFromCli(ctx *cli.Context) error {
	clusterFile, filePath, err := resolveClusterFile(ctx)
	if err != nil {
		return fmt.Errorf("Failed to resolve cluster file: %v", err)
	}

	rkeConfig, err := cluster.ParseConfig(clusterFile)
	if err != nil {
		return fmt.Errorf("Failed to parse cluster file: %v", err)
	}
	rkeConfig, err = setOptionsFromCLI(ctx, rkeConfig)
	if err != nil {
		return err
	}
	// setting up the flags
	externalFlags := cluster.GetExternalFlags(false, false, false, "", filePath)
	externalFlags.CertificateDir = ctx.String("cert-dir")
	externalFlags.CustomCerts = ctx.Bool("custom-certs")

	return GenerateRKECSRs(context.Background(), rkeConfig, externalFlags)
}

func showRKECertificatesFromCli(ctx *cli.Context) error {
	return nil
}

func rebuildClusterWithRotatedCertificates(ctx context.Context,
	dialersOptions hosts.DialersOptions,
	flags cluster.ExternalFlags) (string, string, string, string, map[string]pki.CertificatePKI, error) {
	var APIURL, caCrt, clientCert, clientKey string
	log.Infof(ctx, "Rebuilding Kubernetes cluster with rotated certificates")
	clusterState, err := cluster.ReadStateFile(ctx, cluster.GetStateFilePath(flags.ClusterFilePath, flags.ConfigDir))
	if err != nil {
		return APIURL, caCrt, clientCert, clientKey, nil, err
	}

	kubeCluster, err := cluster.InitClusterObject(ctx, clusterState.DesiredState.RancherKubernetesEngineConfig.DeepCopy(), flags)
	if err != nil {
		return APIURL, caCrt, clientCert, clientKey, nil, err
	}
	if err := kubeCluster.SetupDialers(ctx, dialersOptions); err != nil {
		return APIURL, caCrt, clientCert, clientKey, nil, err
	}

	if err := kubeCluster.TunnelHosts(ctx, flags); err != nil {
		return APIURL, caCrt, clientCert, clientKey, nil, err
	}

	if err := cluster.SetUpAuthentication(ctx, kubeCluster, nil, clusterState); err != nil {
		return APIURL, caCrt, clientCert, clientKey, nil, err
	}
	APIURL = fmt.Sprintf("https://" + kubeCluster.ControlPlaneHosts[0].Address + ":6443")
	clientCert = string(cert.EncodeCertPEM(kubeCluster.Certificates[pki.KubeAdminCertName].Certificate))
	clientKey = string(cert.EncodePrivateKeyPEM(kubeCluster.Certificates[pki.KubeAdminCertName].Key))
	caCrt = string(cert.EncodeCertPEM(kubeCluster.Certificates[pki.CACertName].Certificate))

	if err := kubeCluster.SetUpHosts(ctx, flags); err != nil {
		return APIURL, caCrt, clientCert, clientKey, nil, err
	}
	// Save new State
	if err := kubeCluster.UpdateClusterCurrentState(ctx, clusterState); err != nil {
		return APIURL, caCrt, clientCert, clientKey, nil, err
	}

	// Restarting Kubernetes components
	servicesMap := make(map[string]bool)
	for _, component := range kubeCluster.RotateCertificates.Services {
		servicesMap[component] = true
	}

	if len(kubeCluster.RotateCertificates.Services) == 0 || kubeCluster.RotateCertificates.CACertificates || servicesMap[services.EtcdContainerName] {
		if err := services.RestartEtcdPlane(ctx, kubeCluster.EtcdHosts); err != nil {
			return APIURL, caCrt, clientCert, clientKey, nil, err
		}
	}

	if err := services.RestartControlPlane(ctx, kubeCluster.ControlPlaneHosts); err != nil {
		return APIURL, caCrt, clientCert, clientKey, nil, err
	}

	allHosts := hosts.GetUniqueHostList(kubeCluster.EtcdHosts, kubeCluster.ControlPlaneHosts, kubeCluster.WorkerHosts)
	if err := services.RestartWorkerPlane(ctx, allHosts); err != nil {
		return APIURL, caCrt, clientCert, clientKey, nil, err
	}

	if kubeCluster.RotateCertificates.CACertificates {
		if err := cluster.RestartClusterPods(ctx, kubeCluster); err != nil {
			return APIURL, caCrt, clientCert, clientKey, nil, err
		}
	}
	return APIURL, caCrt, clientCert, clientKey, kubeCluster.Certificates, nil
}

func rotateRKECertificates(ctx context.Context, kubeCluster *cluster.Cluster, flags cluster.ExternalFlags, rkeFullState *cluster.FullState) (*cluster.FullState, error) {
	log.Infof(ctx, "Rotating Kubernetes cluster certificates")
	currentCluster, err := kubeCluster.GetClusterState(ctx, rkeFullState)
	if err != nil {
		return nil, err
	}
	if currentCluster == nil {
		return nil, fmt.Errorf("Failed to rotate certificates: can't find old certificates")
	}
	currentCluster.RotateCertificates = kubeCluster.RotateCertificates
	if err := cluster.RotateRKECertificates(ctx, currentCluster, flags, rkeFullState); err != nil {
		return nil, err
	}
	return rkeFullState, nil
}

func GenerateRKECSRs(ctx context.Context, rkeConfig *v3.RancherKubernetesEngineConfig, flags cluster.ExternalFlags) error {
	log.Infof(ctx, "Generating Kubernetes cluster CSR certificates")
	if len(flags.CertificateDir) == 0 {
		flags.CertificateDir = cluster.GetCertificateDirPath(flags.ClusterFilePath, flags.ConfigDir)
	}

	certBundle, err := pki.ReadCSRsAndKeysFromDir(flags.CertificateDir)
	if err != nil {
		return err
	}

	// initialze the cluster object from the config file
	kubeCluster, err := cluster.InitClusterObject(ctx, rkeConfig, flags)
	if err != nil {
		return err
	}

	// Generating csrs for kubernetes components
	if err := pki.GenerateRKEServicesCSRs(ctx, certBundle, kubeCluster.RancherKubernetesEngineConfig); err != nil {
		return err
	}
	return pki.WriteCertificates(kubeCluster.CertificateDir, certBundle)
}
