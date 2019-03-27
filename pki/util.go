package pki

import (
	cryptorand "crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"math/big"
	"net"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/zdnscloud/zke/hosts"
	"github.com/zdnscloud/zke/types"
	"k8s.io/client-go/util/cert"
)

func GenerateSignedCertAndKey(
	caCrt *x509.Certificate,
	caKey *rsa.PrivateKey,
	serverCrt bool,
	commonName string,
	altNames *cert.AltNames,
	reusedKey *rsa.PrivateKey,
	orgs []string) (*x509.Certificate, *rsa.PrivateKey, error) {
	// Generate a generic signed certificate
	var rootKey *rsa.PrivateKey
	var err error
	rootKey = reusedKey
	if reusedKey == nil {
		rootKey, err = cert.NewPrivateKey()
		if err != nil {
			return nil, nil, fmt.Errorf("Failed to generate private key for %s certificate: %v", commonName, err)
		}
	}
	usages := []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}
	if serverCrt {
		usages = append(usages, x509.ExtKeyUsageServerAuth)
	}
	if altNames == nil {
		altNames = &cert.AltNames{}
	}
	caConfig := cert.Config{
		CommonName:   commonName,
		Organization: orgs,
		Usages:       usages,
		AltNames:     *altNames,
	}
	clientCert, err := newSignedCert(caConfig, rootKey, caCrt, caKey)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to generate %s certificate: %v", commonName, err)
	}
	return clientCert, rootKey, nil
}

func GenerateCertSigningRequestAndKey(
	serverCrt bool,
	commonName string,
	altNames *cert.AltNames,
	reusedKey *rsa.PrivateKey,
	orgs []string) ([]byte, *rsa.PrivateKey, error) {
	// Generate a generic signed certificate
	var rootKey *rsa.PrivateKey
	var err error
	rootKey = reusedKey
	if reusedKey == nil {
		rootKey, err = cert.NewPrivateKey()
		if err != nil {
			return nil, nil, fmt.Errorf("Failed to generate private key for %s certificate: %v", commonName, err)
		}
	}
	usages := []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}
	if serverCrt {
		usages = append(usages, x509.ExtKeyUsageServerAuth)
	}
	if altNames == nil {
		altNames = &cert.AltNames{}
	}
	caConfig := cert.Config{
		CommonName:   commonName,
		Organization: orgs,
		Usages:       usages,
		AltNames:     *altNames,
	}
	clientCSR, err := newCertSigningRequest(caConfig, rootKey)

	if err != nil {
		return nil, nil, fmt.Errorf("Failed to generate %s certificate: %v", commonName, err)
	}
	return clientCSR, rootKey, nil
}

func GenerateCACertAndKey(commonName string, privateKey *rsa.PrivateKey) (*x509.Certificate, *rsa.PrivateKey, error) {
	var err error
	rootKey := privateKey
	if rootKey == nil {
		rootKey, err = cert.NewPrivateKey()
		if err != nil {
			return nil, nil, fmt.Errorf("Failed to generate private key for CA certificate: %v", err)
		}
	}
	caConfig := cert.Config{
		CommonName: commonName,
	}
	kubeCACert, err := cert.NewSelfSignedCACert(caConfig, rootKey)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to generate CA certificate: %v", err)
	}

	return kubeCACert, rootKey, nil
}

func GetAltNames(cpHosts []*hosts.Host, clusterDomain string, KubernetesServiceIP net.IP, SANs []string) *cert.AltNames {
	ips := []net.IP{}
	dnsNames := []string{}
	for _, host := range cpHosts {
		// Check if node address is a valid IP
		if nodeIP := net.ParseIP(host.Address); nodeIP != nil {
			ips = append(ips, nodeIP)
		} else {
			dnsNames = append(dnsNames, host.Address)
		}

		// Check if node internal address is a valid IP
		if len(host.InternalAddress) != 0 && host.InternalAddress != host.Address {
			if internalIP := net.ParseIP(host.InternalAddress); internalIP != nil {
				ips = append(ips, internalIP)
			} else {
				dnsNames = append(dnsNames, host.InternalAddress)
			}
		}
		// Add hostname to the ALT dns names
		if len(host.HostnameOverride) != 0 && host.HostnameOverride != host.Address {
			dnsNames = append(dnsNames, host.HostnameOverride)
		}
	}

	for _, host := range SANs {
		// Check if node address is a valid IP
		if nodeIP := net.ParseIP(host); nodeIP != nil {
			ips = append(ips, nodeIP)
		} else {
			dnsNames = append(dnsNames, host)
		}
	}

	ips = append(ips, net.ParseIP("127.0.0.1"))
	ips = append(ips, KubernetesServiceIP)
	dnsNames = append(dnsNames, []string{
		"localhost",
		"kubernetes",
		"kubernetes.default",
		"kubernetes.default.svc",
		"kubernetes.default.svc." + clusterDomain,
	}...)
	return &cert.AltNames{
		IPs:      ips,
		DNSNames: dnsNames,
	}
}

func (c *CertificatePKI) ToEnv() []string {
	env := []string{}
	if c.Key != nil {
		env = append(env, c.KeyToEnv())
	}
	if c.Certificate != nil {
		env = append(env, c.CertToEnv())
	}
	if c.Config != "" && c.ConfigEnvName != "" {
		env = append(env, c.ConfigToEnv())
	}
	return env
}

func (c *CertificatePKI) CertToEnv() string {
	encodedCrt := cert.EncodeCertPEM(c.Certificate)
	return fmt.Sprintf("%s=%s", c.EnvName, string(encodedCrt))
}

func (c *CertificatePKI) KeyToEnv() string {
	encodedKey := cert.EncodePrivateKeyPEM(c.Key)
	return fmt.Sprintf("%s=%s", c.KeyEnvName, string(encodedKey))
}

func (c *CertificatePKI) ConfigToEnv() string {
	return fmt.Sprintf("%s=%s", c.ConfigEnvName, c.Config)
}

func getEnvFromName(name string) string {
	return strings.Replace(strings.ToUpper(name), "-", "_", -1)
}

func getKeyEnvFromEnv(env string) string {
	return fmt.Sprintf("%s_KEY", env)
}

func getConfigEnvFromEnv(env string) string {
	return fmt.Sprintf("KUBECFG_%s", env)
}

func GetEtcdCrtName(address string) string {
	newAddress := strings.Replace(address, ".", "-", -1)
	return fmt.Sprintf("%s-%s", EtcdCertName, newAddress)
}

func GetCertPath(name string) string {
	return fmt.Sprintf("%s%s.pem", CertPathPrefix, name)
}

func GetKeyPath(name string) string {
	return fmt.Sprintf("%s%s-key.pem", CertPathPrefix, name)
}

func GetConfigPath(name string) string {
	return fmt.Sprintf("%skubecfg-%s.yaml", CertPathPrefix, name)
}

func GetCertTempPath(name string) string {
	return fmt.Sprintf("%s%s.pem", TempCertPath, name)
}

func GetKeyTempPath(name string) string {
	return fmt.Sprintf("%s%s-key.pem", TempCertPath, name)
}

func GetConfigTempPath(name string) string {
	return fmt.Sprintf("%skubecfg-%s.yaml", TempCertPath, name)
}

func ToCertObject(componentName, commonName, ouName string, certificate *x509.Certificate, key *rsa.PrivateKey, csrASN1 []byte) CertificatePKI {
	var config, configPath, configEnvName, certificatePEM, keyPEM string
	var csr *x509.CertificateRequest
	var csrPEM []byte
	if len(commonName) == 0 {
		commonName = getDefaultCN(componentName)
	}

	envName := getEnvFromName(componentName)
	keyEnvName := getKeyEnvFromEnv(envName)
	caCertPath := GetCertPath(CACertName)
	path := GetCertPath(componentName)
	keyPath := GetKeyPath(componentName)
	if certificate != nil {
		certificatePEM = string(cert.EncodeCertPEM(certificate))
	}
	if key != nil {
		keyPEM = string(cert.EncodePrivateKeyPEM(key))
	}
	if csrASN1 != nil {
		csr, _ = x509.ParseCertificateRequest(csrASN1)
		csrPEM = pem.EncodeToMemory(&pem.Block{
			Type: "CERTIFICATE REQUEST", Bytes: csrASN1,
		})
	}

	if componentName != CACertName && componentName != KubeAPICertName && !strings.Contains(componentName, EtcdCertName) && componentName != ServiceAccountTokenKeyName {
		config = getKubeConfigX509("https://127.0.0.1:6443", "local", componentName, caCertPath, path, keyPath)
		configPath = GetConfigPath(componentName)
		configEnvName = getConfigEnvFromEnv(envName)
	}

	return CertificatePKI{
		Certificate:    certificate,
		Key:            key,
		CSR:            csr,
		CertificatePEM: certificatePEM,
		KeyPEM:         keyPEM,
		CSRPEM:         string(csrPEM),
		Config:         config,
		Name:           componentName,
		CommonName:     commonName,
		OUName:         ouName,
		EnvName:        envName,
		KeyEnvName:     keyEnvName,
		ConfigEnvName:  configEnvName,
		Path:           path,
		KeyPath:        keyPath,
		ConfigPath:     configPath,
	}
}

func getDefaultCN(name string) string {
	return fmt.Sprintf("system:%s", name)
}

func getControlCertKeys() []string {
	return []string{
		CACertName,
		KubeAPICertName,
		ServiceAccountTokenKeyName,
		KubeControllerCertName,
		KubeSchedulerCertName,
		KubeProxyCertName,
		KubeNodeCertName,
		EtcdClientCertName,
		EtcdClientCACertName,
		RequestHeaderCACertName,
		APIProxyClientCertName,
	}
}

func getWorkerCertKeys() []string {
	return []string{
		CACertName,
		KubeProxyCertName,
		KubeNodeCertName,
	}
}

func getEtcdCertKeys(rkeNodes []types.RKEConfigNode, etcdRole string) []string {
	certList := []string{
		CACertName,
		KubeProxyCertName,
		KubeNodeCertName,
	}
	etcdHosts := hosts.NodesToHosts(rkeNodes, etcdRole)
	for _, host := range etcdHosts {
		certList = append(certList, GetEtcdCrtName(host.InternalAddress))
	}
	return certList

}

func GetKubernetesServiceIP(serviceClusterRange string) (net.IP, error) {
	ip, ipnet, err := net.ParseCIDR(serviceClusterRange)
	if err != nil {
		return nil, fmt.Errorf("Failed to get kubernetes service IP from Kube API option [service_cluster_ip_range]: %v", err)
	}
	ip = ip.Mask(ipnet.Mask)
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
	return ip, nil
}

func GetLocalKubeConfig(configPath, configDir string) string {
	baseDir := filepath.Dir(configPath)
	if len(configDir) > 0 {
		baseDir = filepath.Dir(configDir)
	}
	fileName := filepath.Base(configPath)
	baseDir += "/"
	return fmt.Sprintf("%s%s%s", baseDir, KubeAdminConfigPrefix, fileName)
}

func strCrtToEnv(crtName, crt string) string {
	return fmt.Sprintf("%s=%s", getEnvFromName(crtName), crt)
}

func strKeyToEnv(crtName, key string) string {
	envName := getEnvFromName(crtName)
	return fmt.Sprintf("%s=%s", getKeyEnvFromEnv(envName), key)
}

func getTempPath(s string) string {
	return TempCertPath + path.Base(s)
}

func populateCertMap(tmpCerts map[string]CertificatePKI, localConfigPath string, extraHosts []*hosts.Host) map[string]CertificatePKI {
	certs := make(map[string]CertificatePKI)
	// CACert
	certs[CACertName] = ToCertObject(CACertName, "", "", tmpCerts[CACertName].Certificate, tmpCerts[CACertName].Key, nil)
	// KubeAPI
	certs[KubeAPICertName] = ToCertObject(KubeAPICertName, "", "", tmpCerts[KubeAPICertName].Certificate, tmpCerts[KubeAPICertName].Key, nil)
	// kubeController
	certs[KubeControllerCertName] = ToCertObject(KubeControllerCertName, "", "", tmpCerts[KubeControllerCertName].Certificate, tmpCerts[KubeControllerCertName].Key, nil)
	// KubeScheduler
	certs[KubeSchedulerCertName] = ToCertObject(KubeSchedulerCertName, "", "", tmpCerts[KubeSchedulerCertName].Certificate, tmpCerts[KubeSchedulerCertName].Key, nil)
	// KubeProxy
	certs[KubeProxyCertName] = ToCertObject(KubeProxyCertName, "", "", tmpCerts[KubeProxyCertName].Certificate, tmpCerts[KubeProxyCertName].Key, nil)
	// KubeNode
	certs[KubeNodeCertName] = ToCertObject(KubeNodeCertName, KubeNodeCommonName, KubeNodeOrganizationName, tmpCerts[KubeNodeCertName].Certificate, tmpCerts[KubeNodeCertName].Key, nil)
	// KubeAdmin
	kubeAdminCertObj := ToCertObject(KubeAdminCertName, KubeAdminCertName, KubeAdminOrganizationName, tmpCerts[KubeAdminCertName].Certificate, tmpCerts[KubeAdminCertName].Key, nil)
	kubeAdminCertObj.Config = tmpCerts[KubeAdminCertName].Config
	kubeAdminCertObj.ConfigPath = localConfigPath
	certs[KubeAdminCertName] = kubeAdminCertObj
	// etcd
	for _, host := range extraHosts {
		etcdName := GetEtcdCrtName(host.InternalAddress)
		etcdCrt, etcdKey := tmpCerts[etcdName].Certificate, tmpCerts[etcdName].Key
		certs[etcdName] = ToCertObject(etcdName, "", "", etcdCrt, etcdKey, nil)
	}

	return certs
}

// Overriding k8s.io/client-go/util/cert.NewSignedCert function to extend the expiration date to 10 years instead of 1 year
func newSignedCert(cfg cert.Config, key *rsa.PrivateKey, caCert *x509.Certificate, caKey *rsa.PrivateKey) (*x509.Certificate, error) {
	serial, err := cryptorand.Int(cryptorand.Reader, new(big.Int).SetInt64(math.MaxInt64))
	if err != nil {
		return nil, err
	}
	if len(cfg.CommonName) == 0 {
		return nil, errors.New("must specify a CommonName")
	}
	if len(cfg.Usages) == 0 {
		return nil, errors.New("must specify at least one ExtKeyUsage")
	}

	certTmpl := x509.Certificate{
		Subject: pkix.Name{
			CommonName:   cfg.CommonName,
			Organization: cfg.Organization,
		},
		DNSNames:     cfg.AltNames.DNSNames,
		IPAddresses:  cfg.AltNames.IPs,
		SerialNumber: serial,
		NotBefore:    caCert.NotBefore,
		NotAfter:     time.Now().Add(duration365d * 10).UTC(),
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  cfg.Usages,
	}
	certDERBytes, err := x509.CreateCertificate(cryptorand.Reader, &certTmpl, caCert, key.Public(), caKey)
	if err != nil {
		return nil, err
	}
	return x509.ParseCertificate(certDERBytes)
}

func newCertSigningRequest(cfg cert.Config, key *rsa.PrivateKey) ([]byte, error) {
	if len(cfg.CommonName) == 0 {
		return nil, errors.New("must specify a CommonName")
	}
	if len(cfg.Usages) == 0 {
		return nil, errors.New("must specify at least one ExtKeyUsage")
	}

	certTmpl := x509.CertificateRequest{
		Subject: pkix.Name{
			CommonName:   cfg.CommonName,
			Organization: cfg.Organization,
		},
		DNSNames:    cfg.AltNames.DNSNames,
		IPAddresses: cfg.AltNames.IPs,
	}
	return x509.CreateCertificateRequest(cryptorand.Reader, &certTmpl, key)
}

func isFileNotFoundErr(e error) bool {
	if strings.Contains(e.Error(), "no such file or directory") ||
		strings.Contains(e.Error(), "Could not find the file") ||
		strings.Contains(e.Error(), "No such container:path:") {
		return true
	}
	return false
}

func deepEqualIPsAltNames(oldIPs, newIPs []net.IP) bool {
	if len(oldIPs) != len(newIPs) {
		return false
	}
	oldIPsStrings := make([]string, len(oldIPs))
	newIPsStrings := make([]string, len(newIPs))
	for i := range oldIPs {
		oldIPsStrings = append(oldIPsStrings, oldIPs[i].String())
		newIPsStrings = append(newIPsStrings, newIPs[i].String())
	}
	return reflect.DeepEqual(oldIPsStrings, newIPsStrings)
}

func TransformPEMToObject(in map[string]CertificatePKI) map[string]CertificatePKI {
	var certificate *x509.Certificate
	out := map[string]CertificatePKI{}
	for k, v := range in {
		certs, _ := cert.ParseCertsPEM([]byte(v.CertificatePEM))
		key, _ := cert.ParsePrivateKeyPEM([]byte(v.KeyPEM))
		if len(certs) > 0 {
			certificate = certs[0]
		}
		if key != nil {
			key = key.(*rsa.PrivateKey)
		}
		o := CertificatePKI{
			ConfigEnvName:  v.ConfigEnvName,
			Name:           v.Name,
			Config:         v.Config,
			CommonName:     v.CommonName,
			OUName:         v.OUName,
			EnvName:        v.EnvName,
			Path:           v.Path,
			KeyEnvName:     v.KeyEnvName,
			KeyPath:        v.KeyPath,
			ConfigPath:     v.ConfigPath,
			Certificate:    certificate,
			CertificatePEM: v.CertificatePEM,
			KeyPEM:         v.KeyPEM,
		}
		if key != nil {
			o.Key = key.(*rsa.PrivateKey)
		}

		out[k] = o
	}
	return out
}

func ReadCSRsAndKeysFromDir(certDir string) (map[string]CertificatePKI, error) {
	certMap := make(map[string]CertificatePKI)
	if _, err := os.Stat(certDir); os.IsNotExist(err) {
		return certMap, nil
	}

	files, err := ioutil.ReadDir(certDir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if strings.Contains(file.Name(), "-csr.pem") {
			certName := strings.TrimSuffix(file.Name(), "-csr.pem")
			logrus.Debugf("[certificates] Loading %s csr from directory [%s]", certName, certDir)
			// fetching csr
			csrASN1, err := getCSRFromFile(certDir, certName+"-csr.pem")
			if err != nil {
				return nil, err
			}
			// fetching key
			key, err := getKeyFromFile(certDir, certName+"-key.pem")
			if err != nil {
				return nil, err
			}
			certMap[certName] = ToCertObject(certName, getCommonName(certName), getOUName(certName), nil, key, csrASN1)
		}
	}

	return certMap, nil
}

func ReadCertsAndKeysFromDir(certDir string) (map[string]CertificatePKI, error) {
	certMap := make(map[string]CertificatePKI)
	if _, err := os.Stat(certDir); os.IsNotExist(err) {
		return certMap, nil
	}

	files, err := ioutil.ReadDir(certDir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		logrus.Debugf("[certificates] reading file %s from directory [%s]", file.Name(), certDir)
		// fetching cert
		cert, err := getCertFromFile(certDir, file.Name())
		if err != nil {
			continue
		}
		// fetching the cert's key
		certName := strings.TrimSuffix(file.Name(), ".pem")
		key, err := getKeyFromFile(certDir, certName+"-key.pem")
		if err != nil {
			continue
		}
		certMap[certName] = ToCertObject(certName, getCommonName(certName), getOUName(certName), cert, key, nil)
	}

	return certMap, nil
}

func getCommonName(certName string) string {
	switch certName {
	case KubeNodeCertName:
		return KubeNodeCommonName
	default:
		return certName
	}
}

func getOUName(certName string) string {
	switch certName {
	case KubeNodeCertName:
		return KubeNodeOrganizationName
	case KubeAdminCertName:
		return KubeAdminOrganizationName
	default:
		return ""
	}
}

func getCertFromFile(certDir string, fileName string) (*x509.Certificate, error) {
	var certificate *x509.Certificate
	certPEM, _ := ioutil.ReadFile(filepath.Join(certDir, fileName))
	if len(certPEM) > 0 {
		certificates, err := cert.ParseCertsPEM(certPEM)
		if err != nil {
			return nil, fmt.Errorf("failed to read certificate [%s]: %v", fileName, err)
		}
		certificate = certificates[0]
	}
	return certificate, nil
}

func getKeyFromFile(certDir string, fileName string) (*rsa.PrivateKey, error) {
	var key *rsa.PrivateKey
	keyPEM, _ := ioutil.ReadFile(filepath.Join(certDir, fileName))
	if len(keyPEM) > 0 {
		keyInterface, err := cert.ParsePrivateKeyPEM(keyPEM)
		if err != nil {
			return nil, fmt.Errorf("failed to read key [%s]: %v", fileName, err)
		}
		key = keyInterface.(*rsa.PrivateKey)
	}
	return key, nil
}

func getCSRFromFile(certDir string, fileName string) ([]byte, error) {
	csrPEM, err := ioutil.ReadFile(filepath.Join(certDir, fileName))
	if err != nil {
		return nil, fmt.Errorf("failed to read csr [%s]: %v", fileName, err)
	}
	csrASN1, _ := pem.Decode(csrPEM)
	return csrASN1.Bytes, nil
}

func WriteCertificates(certDirPath string, certBundle map[string]CertificatePKI) error {
	if _, err := os.Stat(certDirPath); os.IsNotExist(err) {
		err = os.MkdirAll(certDirPath, 0755)
		if err != nil {
			return err
		}
	}

	for certName, cert := range certBundle {
		if cert.CertificatePEM != "" {
			certificatePath := filepath.Join(certDirPath, certName+".pem")
			if err := ioutil.WriteFile(certificatePath, []byte(cert.CertificatePEM), 0640); err != nil {
				return fmt.Errorf("Failed to write certificate to path %v: %v", certificatePath, err)
			}
			logrus.Debugf("Successfully Deployed certificate file at [%s]", certificatePath)
		}

		if cert.KeyPEM != "" {
			keyPath := filepath.Join(certDirPath, certName+"-key.pem")
			if err := ioutil.WriteFile(keyPath, []byte(cert.KeyPEM), 0640); err != nil {
				return fmt.Errorf("Failed to write key to path %v: %v", keyPath, err)
			}
			logrus.Debugf("Successfully Deployed key file at [%s]", keyPath)
		}

		if cert.CSRPEM != "" {
			csrPath := filepath.Join(certDirPath, certName+"-csr.pem")
			if err := ioutil.WriteFile(csrPath, []byte(cert.CSRPEM), 0640); err != nil {
				return fmt.Errorf("Failed to write csr to path %v: %v", csrPath, err)
			}
			logrus.Debugf("Successfully Deployed csr file at [%s]", csrPath)
		}
	}
	logrus.Infof("Successfully Deployed certificates at [%s]", certDirPath)
	return nil
}

func ValidateBundleContent(rkeConfig *types.RancherKubernetesEngineConfig, certBundle map[string]CertificatePKI, configPath, configDir string) error {
	// ensure all needed certs exists
	// make sure all CA Certs exist
	if certBundle[CACertName].Certificate == nil {
		return fmt.Errorf("Failed to find master CA certificate")
	}
	if certBundle[RequestHeaderCACertName].Certificate == nil {
		logrus.Warnf("Failed to find RequestHeader CA certificate, using master CA certificate")
		certBundle[RequestHeaderCACertName] = ToCertObject(RequestHeaderCACertName, RequestHeaderCACertName, "", certBundle[CACertName].Certificate, nil, nil)
	}
	// make sure all components exists
	ComponentsCerts := []string{
		KubeAPICertName,
		KubeControllerCertName,
		KubeSchedulerCertName,
		KubeProxyCertName,
		KubeNodeCertName,
		KubeAdminCertName,
		APIProxyClientCertName,
	}
	for _, certName := range ComponentsCerts {
		if certBundle[certName].Certificate == nil || certBundle[certName].Key == nil {
			return fmt.Errorf("Failed to find [%s] Certificate or Key", certName)
		}
	}
	etcdHosts := hosts.NodesToHosts(rkeConfig.Nodes, etcdRole)
	for _, host := range etcdHosts {
		etcdName := GetEtcdCrtName(host.InternalAddress)
		if certBundle[etcdName].Certificate == nil || certBundle[etcdName].Key == nil {
			return fmt.Errorf("Failed to find etcd [%s] Certificate or Key", etcdName)
		}
	}
	// Configure kubeconfig
	cpHosts := hosts.NodesToHosts(rkeConfig.Nodes, controlRole)
	localKubeConfigPath := GetLocalKubeConfig(configPath, configDir)
	if len(cpHosts) > 0 {
		kubeAdminCertObj := certBundle[KubeAdminCertName]
		kubeAdminConfig := GetKubeConfigX509WithData(
			"https://"+cpHosts[0].Address+":6443",
			rkeConfig.ClusterName,
			KubeAdminCertName,
			string(cert.EncodeCertPEM(certBundle[CACertName].Certificate)),
			string(cert.EncodeCertPEM(certBundle[KubeAdminCertName].Certificate)),
			string(cert.EncodePrivateKeyPEM(certBundle[KubeAdminCertName].Key)))
		kubeAdminCertObj.Config = kubeAdminConfig
		kubeAdminCertObj.ConfigPath = localKubeConfigPath
		certBundle[KubeAdminCertName] = kubeAdminCertObj
	}
	return validateCAIssuer(rkeConfig, certBundle)
}

func validateCAIssuer(rkeConfig *types.RancherKubernetesEngineConfig, certBundle map[string]CertificatePKI) error {
	// make sure all certs are signed by CA cert
	caCert := certBundle[CACertName].Certificate
	ComponentsCerts := []string{
		KubeAPICertName,
		KubeControllerCertName,
		KubeSchedulerCertName,
		KubeProxyCertName,
		KubeNodeCertName,
		KubeAdminCertName,
	}
	etcdHosts := hosts.NodesToHosts(rkeConfig.Nodes, etcdRole)
	for _, host := range etcdHosts {
		etcdName := GetEtcdCrtName(host.InternalAddress)
		ComponentsCerts = append(ComponentsCerts, etcdName)
	}
	for _, componentCert := range ComponentsCerts {
		if certBundle[componentCert].Certificate.Issuer.CommonName != caCert.Subject.CommonName {
			return fmt.Errorf("Component [%s] is not signed by the custom CA certificate", componentCert)
		}
	}
	requestHeaderCACert := certBundle[RequestHeaderCACertName].Certificate
	if certBundle[APIProxyClientCertName].Certificate.Issuer.CommonName != requestHeaderCACert.Subject.CommonName {
		return fmt.Errorf("Component [%s] is not signed by the custom Request Header CA certificate", APIProxyClientCertName)
	}
	return nil
}
