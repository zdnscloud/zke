package registry

import (
	"context"
	"crypto/rsa"
	"encoding/base64"

	"github.com/zdnscloud/zke/core"
	"github.com/zdnscloud/zke/core/pki"
	"github.com/zdnscloud/zke/pkg/docker"
	"github.com/zdnscloud/zke/pkg/hosts"
	"github.com/zdnscloud/zke/pkg/log"
	"github.com/zdnscloud/zke/pkg/templates"
	"github.com/zdnscloud/zke/registry/resources"
	"github.com/zdnscloud/zke/types"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/pkg/sftp"
	"github.com/zdnscloud/gok8s/client"
	"github.com/zdnscloud/gok8s/client/config"
	"github.com/zdnscloud/gok8s/helper"
	"k8s.io/client-go/util/cert"
)

const (
	AdminServerDeployJobName  = "zke-registry-adminserver-deploy-job"
	ChartMuseumDeployJobName  = "zke-registry-chartmuseum-deploy-job"
	ClairDeployJobName        = "zke-registry-clair-deploy-job"
	CoreDeployJobName         = "zke-registry-core-deploy-job"
	DatabaseDeployJobName     = "zke-registry-database-deploy-job"
	IngressDeployJobName      = "zke-registry-ingress-deploy-job"
	JobserviceDeployJobName   = "zke-registry-jobservice-deploy-job"
	NotaryServerDeployJobName = "zke-registry-notaryserver-deploy-job"
	NotarySignerDeployJobName = "zke-registry-notarysigner-deploy-job"
	PortalDeployJobName       = "zke-registry-portal-deploy-job"
	RedisDeployJobName        = "zke-registry-redis-deploy-job"
	RegistryDeployJobName     = "zke-registry-registry-deploy-job"
	RegistryCertsCN           = "harbor"
	DeployNamespace           = "zcloud"
)

type RegistryImage struct {
	HarborAdminserver  string `yaml:"harbor_adminserver" json:"harbor_adminserver"`
	HarborChartmuseum  string `yaml:"harbor_chartmuseum" json:"harbor_chartmuseum"`
	HarborClair        string `yaml:"harbor_clair" json:"harbor_clair"`
	HarborCore         string `yaml:"harbor_core" json:"harbor_core"`
	HarborDatabase     string `yaml:"harbor_database" json:"harbor_database"`
	HarborJobservice   string `yaml:"harbor_jobservice" json:"harbor_jobservice"`
	HarborNotaryServer string `yaml:"harbor_notaryserver" json:"harbor_notaryserver"`
	HarborNotarySigner string `yaml:"harbor_notarysigner" json:"harbor_notarysigner"`
	HarborPortal       string `yaml:"harbor_portal" json:"harbor_portal"`
	HarborRedis        string `yaml:"harbor_redis" json:"harbor_redis"`
	HarborRegistry     string `yaml:"harbor_registry" json:"harbor_registry"`
	HarborRegistryctl  string `yaml:"harbor_registryctl" json:"harbor_registryctl"`
}

var DefaultImage = RegistryImage{
	HarborAdminserver:  "goharbor/harbor-adminserver:v1.7.5",
	HarborChartmuseum:  "goharbor/chartmuseum-photon:v0.8.1-v1.7.5",
	HarborClair:        "goharbor/clair-photon:v2.0.8-v1.7.5",
	HarborCore:         "goharbor/harbor-core:v1.7.5",
	HarborDatabase:     "goharbor/harbor-db:v1.7.5",
	HarborJobservice:   "goharbor/harbor-jobservice:v1.7.5",
	HarborNotaryServer: "goharbor/notary-server-photon:v0.6.1-v1.7.5",
	HarborNotarySigner: "goharbor/notary-signer-photon:v0.6.1-v1.7.5",
	HarborPortal:       "goharbor/harbor-portal:v1.7.5",
	HarborRedis:        "goharbor/redis-photon:v1.7.5",
	HarborRegistry:     "goharbor/registry-photon:v2.6.2-v1.7.5",
	HarborRegistryctl:  "goharbor/harbor-registryctl:v1.7.5",
}

func DeployRegistry(ctx context.Context, c *core.Cluster) error {
	if c.Registry.Isenabled == false {
		log.Infof(ctx, "[Registry] Not enable registry plugin, skip it")
		return nil
	}
	log.Infof(ctx, "[Registry] Setting up Registry Plugin")

	IngresscaCert, IngresstlsCert, IngresstlsKey, err := generateRegistryCerts(c, RegistryCertsCN)
	if err != nil {
		return err
	}
	config := map[string]interface{}{
		"RedisImage":              DefaultImage.HarborRedis,
		"RedisDiskCapacity":       c.Registry.RedisDiskCapacity,
		"DatabaseImage":           DefaultImage.HarborDatabase,
		"DatabaseDiskCapacity":    c.Registry.DatabaseDiskCapacity,
		"CoreImage":               DefaultImage.HarborCore,
		"RegistryImage":           DefaultImage.HarborRegistry,
		"RegistryctlImage":        DefaultImage.HarborRegistryctl,
		"RegistryDiskCapacity":    c.Registry.RegistryDiskCapacity,
		"NotaryServerImage":       DefaultImage.HarborNotaryServer,
		"NotarySignerImage":       DefaultImage.HarborNotarySigner,
		"ChartmuseumImage":        DefaultImage.HarborChartmuseum,
		"ChartmuseumDiskCapacity": c.Registry.ChartmuseumDiskCapacity,
		"ClairImage":              DefaultImage.HarborClair,
		"JobserviceImage":         DefaultImage.HarborJobservice,
		"JobserviceDiskCapacity":  c.Registry.JobserviceDiskCapacity,
		"PortalImage":             DefaultImage.HarborPortal,
		"AdminserverImage":        DefaultImage.HarborAdminserver,
		"RegistryIngressURL":      c.Registry.RegistryIngressURL,
		"NotaryIngressURL":        c.Registry.NotaryIngressURL,
		"IngresscaCertBase64":     getB64Cert(IngresscaCert),
		"IngresstlsCertBase64":    getB64Cert(IngresstlsCert),
		"IngresstlsKeyBase64":     getB64Cert(IngresstlsKey),
		"DeployNamespace":         DeployNamespace,
	}
	// deploy redis
	if err := doOneDeploy(ctx, c, config, resources.RedisTemplate, RedisDeployJobName); err != nil {
		return err
	}
	// deploy database
	if err := doOneDeploy(ctx, c, config, resources.DatabaseTemplate, DatabaseDeployJobName); err != nil {
		return err
	}
	// deploy harbor-core
	if err := doOneDeploy(ctx, c, config, resources.CoreTemplate, CoreDeployJobName); err != nil {
		return err
	}
	// deploy harbor-registry
	if err := doOneDeploy(ctx, c, config, resources.RegistryTemplate, RegistryDeployJobName); err != nil {
		return err
	}
	// deploy notary-server
	if err := doOneDeploy(ctx, c, config, resources.NotaryServerTemplate, NotaryServerDeployJobName); err != nil {
		return err
	}
	// deploy notary-signer
	if err := doOneDeploy(ctx, c, config, resources.NotarySignerTemplate, NotarySignerDeployJobName); err != nil {
		return err
	}
	// deploy chartmuseum
	if err := doOneDeploy(ctx, c, config, resources.ChartMuseumTemplate, ChartMuseumDeployJobName); err != nil {
		return err
	}
	// deploy chair
	if err := doOneDeploy(ctx, c, config, resources.ClairTemplate, ClairDeployJobName); err != nil {
		return err
	}
	// deploy jobservice
	if err := doOneDeploy(ctx, c, config, resources.JobserviceTemplate, JobserviceDeployJobName); err != nil {
		return err
	}
	// deploy portal
	if err := doOneDeploy(ctx, c, config, resources.PortalTemplate, PortalDeployJobName); err != nil {
		return err
	}
	// deploy adminserver
	if err := doOneDeploy(ctx, c, config, resources.AdminServerTemplate, AdminServerDeployJobName); err != nil {
		return err
	}
	// deploy ingress
	if err := doOneDeploy(ctx, c, config, resources.IngressTemplate, IngressDeployJobName); err != nil {
		return err
	}
	if err := deployRegistryCert(ctx, c, IngresscaCert); err != nil {
		return err
	}
	return nil
}

func doOneDeploy(ctx context.Context, c *core.Cluster, config map[string]interface{}, resourcesTemplate string, deployJobName string) error {
	configYaml, err := templates.CompileTemplateFromMap(resourcesTemplate, config)
	if err != nil {
		return err
	}

	if err := c.DoAddonDeploy(ctx, configYaml, deployJobName, true); err != nil {
		return err
	}
	return nil
}

func generateRegistryCerts(c *core.Cluster, commonName string) (string, string, string, error) {
	caCert, caKey, err := pki.GenerateCACertAndKey(commonName, nil)
	if err != nil {
		return "", "", "", err
	}

	ca := pki.ToCertObject("", commonName, commonName, caCert, caKey, nil)

	var tlsTmpKey *rsa.PrivateKey
	tlsAltNames := cert.AltNames{}
	tlsAltNames.DNSNames = append(tlsAltNames.DNSNames, c.Registry.RegistryIngressURL)
	tlsAltNames.DNSNames = append(tlsAltNames.DNSNames, c.Registry.NotaryIngressURL)
	tlsCert, tlsKey, err := pki.GenerateSignedCertAndKey(ca.Certificate, ca.Key, true,
		c.Registry.RegistryIngressURL, &tlsAltNames, tlsTmpKey, nil)
	if err != nil {
		return "", "", "", err
	}
	tls := pki.ToCertObject("", "", "", tlsCert, tlsKey, nil)
	return ca.CertificatePEM, tls.CertificatePEM, tls.KeyPEM, nil
}

func getB64Cert(cert string) string {
	return base64.StdEncoding.EncodeToString([]byte(cert))
}

func deployRegistryCert(ctx context.Context, c *core.Cluster, registryCACert string) error {
	err := c.TunnelHosts(ctx, core.ExternalFlags{})
	if err != nil {
		return nil
	}
	hosts := hosts.GetUniqueHostList(c.EtcdHosts, c.ControlPlaneHosts, c.WorkerHosts, c.StorageHosts, c.EdgeHosts)
	for _, h := range hosts {
		certTmpBasePath := "/home/" + h.User + "/certs.d/"
		certTmpPath := certTmpBasePath + c.Registry.RegistryIngressURL
		sshClient, err := h.GetSSHClient()
		if err != nil {
			return err
		}
		sftpClient, err := h.GetSftpClient(sshClient)
		if err != nil {
			return err
		}
		err = transCertUseSftp(sftpClient, registryCACert, certTmpPath)
		if err != nil {
			return err
		}
		err = moveCerts(ctx, h, certTmpBasePath, c.SystemImages.Alpine, c.PrivateRegistriesMap)
		if err != nil {
			return err
		}
		err = sftpClient.RemoveDirectory(certTmpBasePath)
		if err != nil {
			return err
		}
	}
	return nil
}

func transCertUseSftp(cli *sftp.Client, fileContent string, dstPath string) error {
	if err := cli.MkdirAll(dstPath); err != nil {
		return err
	}
	dstFile, err := cli.Create(dstPath + "/ca.crt")
	defer dstFile.Close()
	if err != nil {
		return err
	}
	if _, err := dstFile.Write([]byte(fileContent)); err != nil {
		return err
	}
	return nil
}

func moveCerts(ctx context.Context, h *hosts.Host, tmpPath string, deployImage string, prsMap map[string]types.PrivateRegistry) error {
	imageCfg := &container.Config{
		Image: deployImage,
		Tty:   true,
		Cmd: []string{
			"mv",
			"/certs.d",
			"/etc/docker/",
		},
	}

	hostcfgMounts := []mount.Mount{
		{
			Type:        "bind",
			Source:      "/etc/docker",
			Target:      "/etc/docker",
			BindOptions: &mount.BindOptions{Propagation: "rshared"},
		},
		{
			Type:        "bind",
			Source:      tmpPath,
			Target:      "/certs.d",
			BindOptions: &mount.BindOptions{Propagation: "rshared"},
		},
	}
	hostCfg := &container.HostConfig{
		Mounts:     hostcfgMounts,
		Privileged: true,
	}

	if err := docker.DoRunContainer(ctx, h.DClient, imageCfg, hostCfg, "harbor-certs-deployer", h.Address, "cleanup", prsMap); err != nil {
		return err
	}
	if _, err := docker.WaitForContainer(ctx, h.DClient, h.Address, "harbor-certs-deployer"); err != nil {
		return err
	}
	if err := docker.DoRemoveContainer(ctx, h.DClient, "harbor-certs-deployer", h.Address); err != nil {
		return err
	}
	return nil
}

func doOneDeployFromYaml(yaml string) error {
	cfg, err := config.GetConfigFromFile("./kube_config_cluster.yml")
	if err != nil {
		return err
	}
	cli, err := client.New(cfg, client.Options{})
	if err != nil {
		return err
	}
	err = helper.CreateResourceFromYaml(cli, yaml)
	return err
}
