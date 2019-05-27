package registry

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/zdnscloud/zke/core"
	"github.com/zdnscloud/zke/core/pki"
	"github.com/zdnscloud/zke/pkg/log"
	"github.com/zdnscloud/zke/pkg/templates"
	"github.com/zdnscloud/zke/registry/resources"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
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

	IngresscaCertBase64, IngresstlsCertBase64, IngresstlsKeyBase64, err := generateIngressCertsBase64(c, RegistryCertsCN)
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
		"IngresscaCertBase64":     IngresscaCertBase64,
		"IngresstlsCertBase64":    IngresstlsCertBase64,
		"IngresstlsKeyBase64":     IngresstlsKeyBase64,
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

func generateIngressCertsBase64(c *core.Cluster, commonName string) (string, string, string, error) {
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

	caCertPEMBase64 := base64.StdEncoding.EncodeToString([]byte(ca.CertificatePEM))
	tlsCertPEMBase64 := base64.StdEncoding.EncodeToString([]byte(tls.CertificatePEM))
	tlsKeyPEMBase64 := base64.StdEncoding.EncodeToString([]byte(tls.KeyPEM))
	return caCertPEMBase64, tlsCertPEMBase64, tlsKeyPEMBase64, nil
}

func connect(user, password, host string, port int) (*sftp.Client, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		sshClient    *ssh.Client
		sftpClient   *sftp.Client
		err          error
	)
	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(password))

	clientConfig = &ssh.ClientConfig{
		User:            user,
		Auth:            auth,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         30 * time.Second,
	}

	// connet to ssh
	addr = fmt.Sprintf("%s:%d", host, port)

	if sshClient, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}

	// create sftp client
	if sftpClient, err = sftp.NewClient(sshClient); err != nil {
		return nil, err
	}

	return sftpClient, nil
}
