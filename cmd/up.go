package cmd

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"github.com/zdnscloud/zke/cluster"
	"github.com/zdnscloud/zke/hosts"
	"github.com/zdnscloud/zke/monitoring"
	"github.com/zdnscloud/zke/network"
	"github.com/zdnscloud/zke/pkg/log"
	"github.com/zdnscloud/zke/pki"
	"github.com/zdnscloud/zke/storage"
	"github.com/zdnscloud/zke/types"
	"github.com/zdnscloud/zke/zcloud"
	"k8s.io/client-go/util/cert"
	"os"
	"strings"
)

const DINDWaitTime = 3

func UpCommand() cli.Command {
	upFlags := []cli.Flag{
		cli.StringFlag{
			Name:   "config",
			Usage:  "Specify an alternate cluster YAML file",
			Value:  pki.ClusterConfig,
			EnvVar: "ZKE_CONFIG",
		},
		cli.BoolFlag{
			Name:  "disable-port-check",
			Usage: "Disable port check validation between nodes",
		},
		cli.StringFlag{
			Name:  "cert-dir",
			Usage: "Specify a certificate dir path",
		},
		cli.BoolFlag{
			Name:  "custom-certs",
			Usage: "Use custom certificates from a cert dir",
		},
	}
	upFlags = append(upFlags, commonFlags...)
	return cli.Command{
		Name:   "up",
		Usage:  "Bring the cluster up",
		Action: clusterUpFromCli,
		Flags:  upFlags,
	}
}

func doUpgradeLegacyCluster(ctx context.Context, kubeCluster *cluster.Cluster, fullState *cluster.FullState) error {
	if _, err := os.Stat(kubeCluster.LocalKubeConfigPath); os.IsNotExist(err) {
		// there is no kubeconfig. This is a new cluster
		logrus.Debug("[state] local kubeconfig not found, this is a new cluster")
		return nil
	}
	if _, err := os.Stat(kubeCluster.StateFilePath); err == nil {
		// this cluster has a previous state, I don't need to upgrade!
		logrus.Debug("[state] previous state found, this is not a legacy cluster")
		return nil
	}
	// We have a kubeconfig and no current state. This is a legacy cluster or a new cluster with old kubeconfig
	// let's try to upgrade
	log.Infof(ctx, "[state] Possible legacy cluster detected, trying to upgrade")
	if err := cluster.RebuildKubeconfig(ctx, kubeCluster); err != nil {
		return err
	}
	recoveredCluster, err := cluster.GetStateFromKubernetes(ctx, kubeCluster)
	if err != nil {
		return err
	}
	// if we found a recovered cluster, we will need override the current state
	if recoveredCluster != nil {
		recoveredCerts, err := cluster.GetClusterCertsFromKubernetes(ctx, kubeCluster)
		if err != nil {
			return err
		}
		fullState.CurrentState.ZcloudKubernetesEngineConfig = recoveredCluster.ZcloudKubernetesEngineConfig.DeepCopy()
		fullState.CurrentState.CertificatesBundle = recoveredCerts
		// we don't want to regenerate certificates
		fullState.DesiredState.CertificatesBundle = recoveredCerts
		return fullState.WriteStateFile(ctx, kubeCluster.StateFilePath)
	}
	return nil
}

func ClusterUp(ctx context.Context, dialersOptions hosts.DialersOptions, flags cluster.ExternalFlags) (string, string, string, string, map[string]pki.CertificatePKI, error) {
	var APIURL, caCrt, clientCert, clientKey string
	clusterState, err := cluster.ReadStateFile(ctx, cluster.GetStateFilePath(flags.ClusterFilePath, flags.ConfigDir))
	if err != nil {
		return APIURL, caCrt, clientCert, clientKey, nil, err
	}

	kubeCluster, err := cluster.InitClusterObject(ctx, clusterState.DesiredState.ZcloudKubernetesEngineConfig.DeepCopy(), flags)
	if err != nil {
		return APIURL, caCrt, clientCert, clientKey, nil, err
	}
	log.Infof(ctx, "Building Kubernetes cluster")
	err = kubeCluster.SetupDialers(ctx, dialersOptions)
	if err != nil {
		return APIURL, caCrt, clientCert, clientKey, nil, err
	}
	err = kubeCluster.TunnelHosts(ctx, flags)
	if err != nil {
		return APIURL, caCrt, clientCert, clientKey, nil, err
	}
	currentCluster, err := kubeCluster.GetClusterState(ctx, clusterState)
	if err != nil {
		return APIURL, caCrt, clientCert, clientKey, nil, err
	}
	if !flags.DisablePortCheck {
		if err = kubeCluster.CheckClusterPorts(ctx, currentCluster); err != nil {
			return APIURL, caCrt, clientCert, clientKey, nil, err
		}
	}
	err = cluster.SetUpAuthentication(ctx, kubeCluster, currentCluster, clusterState)
	if err != nil {
		return APIURL, caCrt, clientCert, clientKey, nil, err
	}
	if len(kubeCluster.ControlPlaneHosts) > 0 {
		APIURL = fmt.Sprintf("https://" + kubeCluster.ControlPlaneHosts[0].Address + ":6443")
	}
	clientCert = string(cert.EncodeCertPEM(kubeCluster.Certificates[pki.KubeAdminCertName].Certificate))
	clientKey = string(cert.EncodePrivateKeyPEM(kubeCluster.Certificates[pki.KubeAdminCertName].Key))
	caCrt = string(cert.EncodeCertPEM(kubeCluster.Certificates[pki.CACertName].Certificate))
	// moved deploying certs before reconcile to remove all unneeded certs generation from reconcile
	err = kubeCluster.SetUpHosts(ctx, flags)
	if err != nil {
		return APIURL, caCrt, clientCert, clientKey, nil, err
	}
	err = cluster.ReconcileCluster(ctx, kubeCluster, currentCluster, flags)
	if err != nil {
		return APIURL, caCrt, clientCert, clientKey, nil, err
	}
	// update APIURL after reconcile
	if len(kubeCluster.ControlPlaneHosts) > 0 {
		APIURL = fmt.Sprintf("https://" + kubeCluster.ControlPlaneHosts[0].Address + ":6443")
	}
	if err := kubeCluster.PrePullK8sImages(ctx); err != nil {
		return APIURL, caCrt, clientCert, clientKey, nil, err
	}
	err = kubeCluster.DeployControlPlane(ctx)
	if err != nil {
		return APIURL, caCrt, clientCert, clientKey, nil, err
	}
	// Apply Authz configuration after deploying controlplane
	err = cluster.ApplyAuthzResources(ctx, kubeCluster.ZcloudKubernetesEngineConfig, flags, dialersOptions)
	if err != nil {
		return APIURL, caCrt, clientCert, clientKey, nil, err
	}
	err = kubeCluster.UpdateClusterCurrentState(ctx, clusterState)
	if err != nil {
		return APIURL, caCrt, clientCert, clientKey, nil, err
	}
	err = cluster.SaveFullStateToKubernetes(ctx, kubeCluster, clusterState)
	if err != nil {
		return APIURL, caCrt, clientCert, clientKey, nil, err
	}
	err = kubeCluster.DeployWorkerPlane(ctx)
	if err != nil {
		return APIURL, caCrt, clientCert, clientKey, nil, err
	}
	if err = kubeCluster.CleanDeadLogs(ctx); err != nil {
		return APIURL, caCrt, clientCert, clientKey, nil, err
	}
	err = kubeCluster.SyncLabelsAndTaints(ctx, currentCluster)
	if err != nil {
		return APIURL, caCrt, clientCert, clientKey, nil, err
	}
	err = ConfigureCluster(ctx, kubeCluster.ZcloudKubernetesEngineConfig, kubeCluster.Certificates, flags, dialersOptions, false)
	if err != nil {
		return APIURL, caCrt, clientCert, clientKey, nil, err
	}
	if err := checkAllIncluded(kubeCluster); err != nil {
		return APIURL, caCrt, clientCert, clientKey, nil, err
	}
	log.Infof(ctx, "Finished building Kubernetes cluster successfully")
	return APIURL, caCrt, clientCert, clientKey, kubeCluster.Certificates, nil
}

func checkAllIncluded(cluster *cluster.Cluster) error {
	if len(cluster.InactiveHosts) == 0 {
		return nil
	}
	var names []string
	for _, host := range cluster.InactiveHosts {
		names = append(names, host.Address)
	}
	return fmt.Errorf("Provisioning incomplete, host(s) [%s] skipped because they could not be contacted", strings.Join(names, ","))
}

func clusterUpFromCli(ctx *cli.Context) error {
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
	disablePortCheck := ctx.Bool("disable-port-check")
	// setting up the flags
	flags := cluster.GetExternalFlags(disablePortCheck, "", filePath)
	// Custom certificates and certificate dir flags
	flags.CertificateDir = ctx.String("cert-dir")
	flags.CustomCerts = ctx.Bool("custom-certs")
	if err := ClusterInit(context.Background(), zkeConfig, hosts.DialersOptions{}, flags); err != nil {
		return err
	}
	_, _, _, _, _, err = ClusterUp(context.Background(), hosts.DialersOptions{}, flags)
	return err
}

func ConfigureCluster(
	ctx context.Context,
	zkeConfig types.ZcloudKubernetesEngineConfig,
	crtBundle map[string]pki.CertificatePKI,
	flags cluster.ExternalFlags,
	dailersOptions hosts.DialersOptions,
	useKubectl bool) error {
	// dialer factories are not needed here since we are not uses docker only k8s jobs
	kubeCluster, err := cluster.InitClusterObject(ctx, &zkeConfig, flags)
	if err != nil {
		return err
	}
	if err := kubeCluster.SetupDialers(ctx, dailersOptions); err != nil {
		return err
	}
	kubeCluster.UseKubectlDeploy = useKubectl
	if len(kubeCluster.ControlPlaneHosts) > 0 {
		kubeCluster.Certificates = crtBundle
		if err := network.DeployNetwork(ctx, kubeCluster); err != nil {
			if err, ok := err.(*cluster.AddonError); ok && err.IsCritical {
				return err
			}
			log.Warnf(ctx, "Failed to deploy addon execute job [%s]: %v", network.NetworkPluginResourceName, err)
		}

		if err := storage.DeployStoragePlugin(ctx, kubeCluster); err != nil {
			return err
		}

		if err := kubeCluster.DeployAddons(ctx); err != nil {
			return err
		}

		if err := zcloud.DeployZcloudManager(ctx, kubeCluster); err != nil {
			return err
		}

		if err := monitoring.DeployMonitoring(ctx, kubeCluster); err != nil {
			return err
		}
	}
	return nil
}
