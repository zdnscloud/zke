package pki

import (
	"context"
	"crypto/rsa"
	"fmt"
	"reflect"

	"github.com/rancher/rke/hosts"
	"github.com/rancher/rke/log"
	"github.com/rancher/types/apis/management.cattle.io/v3"
	"k8s.io/client-go/util/cert"
)

func GenerateKubeAPICertificate(ctx context.Context, certs map[string]CertificatePKI, rkeConfig v3.RancherKubernetesEngineConfig, configPath, configDir string, rotate bool) error {
	// generate API certificate and key
	caCrt := certs[CACertName].Certificate
	caKey := certs[CACertName].Key
	kubernetesServiceIP, err := GetKubernetesServiceIP(rkeConfig.Services.KubeAPI.ServiceClusterIPRange)
	if err != nil {
		return fmt.Errorf("Failed to get Kubernetes Service IP: %v", err)
	}
	clusterDomain := rkeConfig.Services.Kubelet.ClusterDomain
	cpHosts := hosts.NodesToHosts(rkeConfig.Nodes, controlRole)
	kubeAPIAltNames := GetAltNames(cpHosts, clusterDomain, kubernetesServiceIP, rkeConfig.Authentication.SANs)
	kubeAPICert := certs[KubeAPICertName].Certificate
	if kubeAPICert != nil &&
		reflect.DeepEqual(kubeAPIAltNames.DNSNames, kubeAPICert.DNSNames) &&
		deepEqualIPsAltNames(kubeAPIAltNames.IPs, kubeAPICert.IPAddresses) && !rotate {
		return nil
	}
	log.Infof(ctx, "[certificates] Generating Kubernetes API server certificates")
	var serviceKey *rsa.PrivateKey
	if !rotate {
		serviceKey = certs[KubeAPICertName].Key
	}
	kubeAPICrt, kubeAPIKey, err := GenerateSignedCertAndKey(caCrt, caKey, true, KubeAPICertName, kubeAPIAltNames, serviceKey, nil)
	if err != nil {
		return err
	}
	certs[KubeAPICertName] = ToCertObject(KubeAPICertName, "", "", kubeAPICrt, kubeAPIKey, nil)
	// handle service account tokens in old clusters
	apiCert := certs[KubeAPICertName]
	if certs[ServiceAccountTokenKeyName].Key == nil {
		log.Infof(ctx, "[certificates] Generating Service account token key")
		certs[ServiceAccountTokenKeyName] = ToCertObject(ServiceAccountTokenKeyName, ServiceAccountTokenKeyName, "", apiCert.Certificate, apiCert.Key, nil)
	}
	return nil
}

func GenerateKubeAPICSR(ctx context.Context, certs map[string]CertificatePKI, rkeConfig v3.RancherKubernetesEngineConfig) error {
	// generate API csr and key
	kubernetesServiceIP, err := GetKubernetesServiceIP(rkeConfig.Services.KubeAPI.ServiceClusterIPRange)
	if err != nil {
		return fmt.Errorf("Failed to get Kubernetes Service IP: %v", err)
	}
	clusterDomain := rkeConfig.Services.Kubelet.ClusterDomain
	cpHosts := hosts.NodesToHosts(rkeConfig.Nodes, controlRole)
	kubeAPIAltNames := GetAltNames(cpHosts, clusterDomain, kubernetesServiceIP, rkeConfig.Authentication.SANs)
	kubeAPICert := certs[KubeAPICertName].Certificate
	oldKubeAPICSR := certs[KubeAPICertName].CSR
	if oldKubeAPICSR != nil &&
		reflect.DeepEqual(kubeAPIAltNames.DNSNames, oldKubeAPICSR.DNSNames) &&
		deepEqualIPsAltNames(kubeAPIAltNames.IPs, oldKubeAPICSR.IPAddresses) {
		return nil
	}
	log.Infof(ctx, "[certificates] Generating Kubernetes API server csr")
	kubeAPICSR, kubeAPIKey, err := GenerateCertSigningRequestAndKey(true, KubeAPICertName, kubeAPIAltNames, certs[KubeAPICertName].Key, nil)
	if err != nil {
		return err
	}
	certs[KubeAPICertName] = ToCertObject(KubeAPICertName, "", "", kubeAPICert, kubeAPIKey, kubeAPICSR)
	return nil
}

func GenerateKubeControllerCertificate(ctx context.Context, certs map[string]CertificatePKI, rkeConfig v3.RancherKubernetesEngineConfig, configPath, configDir string, rotate bool) error {
	// generate Kube controller-manager certificate and key
	caCrt := certs[CACertName].Certificate
	caKey := certs[CACertName].Key
	if certs[KubeControllerCertName].Certificate != nil && !rotate {
		return nil
	}
	log.Infof(ctx, "[certificates] Generating Kube Controller certificates")
	var serviceKey *rsa.PrivateKey
	if !rotate {
		serviceKey = certs[KubeControllerCertName].Key
	}
	kubeControllerCrt, kubeControllerKey, err := GenerateSignedCertAndKey(caCrt, caKey, false, getDefaultCN(KubeControllerCertName), nil, serviceKey, nil)
	if err != nil {
		return err
	}
	certs[KubeControllerCertName] = ToCertObject(KubeControllerCertName, "", "", kubeControllerCrt, kubeControllerKey, nil)
	return nil
}

func GenerateKubeControllerCSR(ctx context.Context, certs map[string]CertificatePKI, rkeConfig v3.RancherKubernetesEngineConfig) error {
	// generate Kube controller-manager csr and key
	kubeControllerCrt := certs[KubeControllerCertName].Certificate
	kubeControllerCSRPEM := certs[KubeControllerCertName].CSRPEM
	if kubeControllerCSRPEM != "" {
		return nil
	}
	log.Infof(ctx, "[certificates] Generating Kube Controller csr")
	kubeControllerCSR, kubeControllerKey, err := GenerateCertSigningRequestAndKey(false, getDefaultCN(KubeControllerCertName), nil, certs[KubeControllerCertName].Key, nil)
	if err != nil {
		return err
	}
	certs[KubeControllerCertName] = ToCertObject(KubeControllerCertName, "", "", kubeControllerCrt, kubeControllerKey, kubeControllerCSR)
	return nil
}

func GenerateKubeSchedulerCertificate(ctx context.Context, certs map[string]CertificatePKI, rkeConfig v3.RancherKubernetesEngineConfig, configPath, configDir string, rotate bool) error {
	// generate Kube scheduler certificate and key
	caCrt := certs[CACertName].Certificate
	caKey := certs[CACertName].Key
	if certs[KubeSchedulerCertName].Certificate != nil && !rotate {
		return nil
	}
	log.Infof(ctx, "[certificates] Generating Kube Scheduler certificates")
	var serviceKey *rsa.PrivateKey
	if !rotate {
		serviceKey = certs[KubeSchedulerCertName].Key
	}
	kubeSchedulerCrt, kubeSchedulerKey, err := GenerateSignedCertAndKey(caCrt, caKey, false, getDefaultCN(KubeSchedulerCertName), nil, serviceKey, nil)
	if err != nil {
		return err
	}
	certs[KubeSchedulerCertName] = ToCertObject(KubeSchedulerCertName, "", "", kubeSchedulerCrt, kubeSchedulerKey, nil)
	return nil
}

func GenerateKubeSchedulerCSR(ctx context.Context, certs map[string]CertificatePKI, rkeConfig v3.RancherKubernetesEngineConfig) error {
	// generate Kube scheduler csr and key
	kubeSchedulerCrt := certs[KubeSchedulerCertName].Certificate
	kubeSchedulerCSRPEM := certs[KubeSchedulerCertName].CSRPEM
	if kubeSchedulerCSRPEM != "" {
		return nil
	}
	log.Infof(ctx, "[certificates] Generating Kube Scheduler csr")
	kubeSchedulerCSR, kubeSchedulerKey, err := GenerateCertSigningRequestAndKey(false, getDefaultCN(KubeSchedulerCertName), nil, certs[KubeSchedulerCertName].Key, nil)
	if err != nil {
		return err
	}
	certs[KubeSchedulerCertName] = ToCertObject(KubeSchedulerCertName, "", "", kubeSchedulerCrt, kubeSchedulerKey, kubeSchedulerCSR)
	return nil
}

func GenerateKubeProxyCertificate(ctx context.Context, certs map[string]CertificatePKI, rkeConfig v3.RancherKubernetesEngineConfig, configPath, configDir string, rotate bool) error {
	// generate Kube Proxy certificate and key
	caCrt := certs[CACertName].Certificate
	caKey := certs[CACertName].Key
	if certs[KubeProxyCertName].Certificate != nil && !rotate {
		return nil
	}
	log.Infof(ctx, "[certificates] Generating Kube Proxy certificates")
	var serviceKey *rsa.PrivateKey
	if !rotate {
		serviceKey = certs[KubeProxyCertName].Key
	}
	kubeProxyCrt, kubeProxyKey, err := GenerateSignedCertAndKey(caCrt, caKey, false, getDefaultCN(KubeProxyCertName), nil, serviceKey, nil)
	if err != nil {
		return err
	}
	certs[KubeProxyCertName] = ToCertObject(KubeProxyCertName, "", "", kubeProxyCrt, kubeProxyKey, nil)
	return nil
}

func GenerateKubeProxyCSR(ctx context.Context, certs map[string]CertificatePKI, rkeConfig v3.RancherKubernetesEngineConfig) error {
	// generate Kube Proxy csr and key
	kubeProxyCrt := certs[KubeProxyCertName].Certificate
	kubeProxyCSRPEM := certs[KubeProxyCertName].CSRPEM
	if kubeProxyCSRPEM != "" {
		return nil
	}
	log.Infof(ctx, "[certificates] Generating Kube Proxy csr")
	kubeProxyCSR, kubeProxyKey, err := GenerateCertSigningRequestAndKey(false, getDefaultCN(KubeProxyCertName), nil, certs[KubeProxyCertName].Key, nil)
	if err != nil {
		return err
	}
	certs[KubeProxyCertName] = ToCertObject(KubeProxyCertName, "", "", kubeProxyCrt, kubeProxyKey, kubeProxyCSR)
	return nil
}

func GenerateKubeNodeCertificate(ctx context.Context, certs map[string]CertificatePKI, rkeConfig v3.RancherKubernetesEngineConfig, configPath, configDir string, rotate bool) error {
	// generate kubelet certificate
	caCrt := certs[CACertName].Certificate
	caKey := certs[CACertName].Key
	if certs[KubeNodeCertName].Certificate != nil && !rotate {
		return nil
	}
	log.Infof(ctx, "[certificates] Generating Node certificate")
	var serviceKey *rsa.PrivateKey
	if !rotate {
		serviceKey = certs[KubeProxyCertName].Key
	}
	nodeCrt, nodeKey, err := GenerateSignedCertAndKey(caCrt, caKey, false, KubeNodeCommonName, nil, serviceKey, []string{KubeNodeOrganizationName})
	if err != nil {
		return err
	}
	certs[KubeNodeCertName] = ToCertObject(KubeNodeCertName, KubeNodeCommonName, KubeNodeOrganizationName, nodeCrt, nodeKey, nil)
	return nil
}

func GenerateKubeNodeCSR(ctx context.Context, certs map[string]CertificatePKI, rkeConfig v3.RancherKubernetesEngineConfig) error {
	// generate kubelet csr and key
	nodeCrt := certs[KubeNodeCertName].Certificate
	nodeCSRPEM := certs[KubeNodeCertName].CSRPEM
	if nodeCSRPEM != "" {
		return nil
	}
	log.Infof(ctx, "[certificates] Generating Node csr and key")
	nodeCSR, nodeKey, err := GenerateCertSigningRequestAndKey(false, KubeNodeCommonName, nil, certs[KubeNodeCertName].Key, []string{KubeNodeOrganizationName})
	if err != nil {
		return err
	}
	certs[KubeNodeCertName] = ToCertObject(KubeNodeCertName, KubeNodeCommonName, KubeNodeOrganizationName, nodeCrt, nodeKey, nodeCSR)
	return nil
}

func GenerateKubeAdminCertificate(ctx context.Context, certs map[string]CertificatePKI, rkeConfig v3.RancherKubernetesEngineConfig, configPath, configDir string, rotate bool) error {
	// generate Admin certificate and key
	log.Infof(ctx, "[certificates] Generating admin certificates and kubeconfig")
	caCrt := certs[CACertName].Certificate
	caKey := certs[CACertName].Key
	cpHosts := hosts.NodesToHosts(rkeConfig.Nodes, controlRole)
	if len(configPath) == 0 {
		configPath = ClusterConfig
	}
	localKubeConfigPath := GetLocalKubeConfig(configPath, configDir)
	var serviceKey *rsa.PrivateKey
	if !rotate {
		serviceKey = certs[KubeAdminCertName].Key
	}
	kubeAdminCrt, kubeAdminKey, err := GenerateSignedCertAndKey(caCrt, caKey, false, KubeAdminCertName, nil, serviceKey, []string{KubeAdminOrganizationName})
	if err != nil {
		return err
	}
	kubeAdminCertObj := ToCertObject(KubeAdminCertName, KubeAdminCertName, KubeAdminOrganizationName, kubeAdminCrt, kubeAdminKey, nil)
	if len(cpHosts) > 0 {
		kubeAdminConfig := GetKubeConfigX509WithData(
			"https://"+cpHosts[0].Address+":6443",
			rkeConfig.ClusterName,
			KubeAdminCertName,
			string(cert.EncodeCertPEM(caCrt)),
			string(cert.EncodeCertPEM(kubeAdminCrt)),
			string(cert.EncodePrivateKeyPEM(kubeAdminKey)))
		kubeAdminCertObj.Config = kubeAdminConfig
		kubeAdminCertObj.ConfigPath = localKubeConfigPath
	} else {
		kubeAdminCertObj.Config = ""
	}
	certs[KubeAdminCertName] = kubeAdminCertObj
	return nil
}

func GenerateKubeAdminCSR(ctx context.Context, certs map[string]CertificatePKI, rkeConfig v3.RancherKubernetesEngineConfig) error {
	// generate Admin certificate and key
	kubeAdminCrt := certs[KubeAdminCertName].Certificate
	kubeAdminCSRPEM := certs[KubeAdminCertName].CSRPEM
	if kubeAdminCSRPEM != "" {
		return nil
	}
	kubeAdminCSR, kubeAdminKey, err := GenerateCertSigningRequestAndKey(false, KubeAdminCertName, nil, certs[KubeAdminCertName].Key, []string{KubeAdminOrganizationName})
	if err != nil {
		return err
	}
	log.Infof(ctx, "[certificates] Generating admin csr and kubeconfig")
	kubeAdminCertObj := ToCertObject(KubeAdminCertName, KubeAdminCertName, KubeAdminOrganizationName, kubeAdminCrt, kubeAdminKey, kubeAdminCSR)
	certs[KubeAdminCertName] = kubeAdminCertObj
	return nil
}

func GenerateAPIProxyClientCertificate(ctx context.Context, certs map[string]CertificatePKI, rkeConfig v3.RancherKubernetesEngineConfig, configPath, configDir string, rotate bool) error {
	//generate API server proxy client key and certs
	caCrt := certs[RequestHeaderCACertName].Certificate
	caKey := certs[RequestHeaderCACertName].Key
	if certs[APIProxyClientCertName].Certificate != nil && !rotate {
		return nil
	}
	log.Infof(ctx, "[certificates] Generating Kubernetes API server proxy client certificates")
	var serviceKey *rsa.PrivateKey
	if !rotate {
		serviceKey = certs[APIProxyClientCertName].Key
	}
	apiserverProxyClientCrt, apiserverProxyClientKey, err := GenerateSignedCertAndKey(caCrt, caKey, true, APIProxyClientCertName, nil, serviceKey, nil)
	if err != nil {
		return err
	}
	certs[APIProxyClientCertName] = ToCertObject(APIProxyClientCertName, "", "", apiserverProxyClientCrt, apiserverProxyClientKey, nil)
	return nil
}

func GenerateAPIProxyClientCSR(ctx context.Context, certs map[string]CertificatePKI, rkeConfig v3.RancherKubernetesEngineConfig) error {
	//generate API server proxy client key and certs
	apiserverProxyClientCrt := certs[APIProxyClientCertName].Certificate
	apiserverProxyClientCSRPEM := certs[APIProxyClientCertName].CSRPEM
	if apiserverProxyClientCSRPEM != "" {
		return nil
	}
	log.Infof(ctx, "[certificates] Generating Kubernetes API server proxy client csr")
	apiserverProxyClientCSR, apiserverProxyClientKey, err := GenerateCertSigningRequestAndKey(true, APIProxyClientCertName, nil, certs[APIProxyClientCertName].Key, nil)
	if err != nil {
		return err
	}
	certs[APIProxyClientCertName] = ToCertObject(APIProxyClientCertName, "", "", apiserverProxyClientCrt, apiserverProxyClientKey, apiserverProxyClientCSR)
	return nil
}

func GenerateExternalEtcdCertificates(ctx context.Context, certs map[string]CertificatePKI, rkeConfig v3.RancherKubernetesEngineConfig, configPath, configDir string, rotate bool) error {
	clientCert, err := cert.ParseCertsPEM([]byte(rkeConfig.Services.Etcd.Cert))
	if err != nil {
		return err
	}
	clientKey, err := cert.ParsePrivateKeyPEM([]byte(rkeConfig.Services.Etcd.Key))
	if err != nil {
		return err
	}
	certs[EtcdClientCertName] = ToCertObject(EtcdClientCertName, "", "", clientCert[0], clientKey.(*rsa.PrivateKey), nil)

	caCert, err := cert.ParseCertsPEM([]byte(rkeConfig.Services.Etcd.CACert))
	if err != nil {
		return err
	}
	certs[EtcdClientCACertName] = ToCertObject(EtcdClientCACertName, "", "", caCert[0], nil, nil)
	return nil
}

func GenerateEtcdCertificates(ctx context.Context, certs map[string]CertificatePKI, rkeConfig v3.RancherKubernetesEngineConfig, configPath, configDir string, rotate bool) error {
	caCrt := certs[CACertName].Certificate
	caKey := certs[CACertName].Key
	kubernetesServiceIP, err := GetKubernetesServiceIP(rkeConfig.Services.KubeAPI.ServiceClusterIPRange)
	if err != nil {
		return fmt.Errorf("Failed to get Kubernetes Service IP: %v", err)
	}
	clusterDomain := rkeConfig.Services.Kubelet.ClusterDomain
	etcdHosts := hosts.NodesToHosts(rkeConfig.Nodes, etcdRole)
	etcdAltNames := GetAltNames(etcdHosts, clusterDomain, kubernetesServiceIP, []string{})
	for _, host := range etcdHosts {
		etcdName := GetEtcdCrtName(host.InternalAddress)
		if _, ok := certs[etcdName]; ok && !rotate {
			continue
		}
		var serviceKey *rsa.PrivateKey
		if !rotate {
			serviceKey = certs[etcdName].Key
		}
		log.Infof(ctx, "[certificates] Generating etcd-%s certificate and key", host.InternalAddress)
		etcdCrt, etcdKey, err := GenerateSignedCertAndKey(caCrt, caKey, true, EtcdCertName, etcdAltNames, serviceKey, nil)
		if err != nil {
			return err
		}
		certs[etcdName] = ToCertObject(etcdName, "", "", etcdCrt, etcdKey, nil)
	}
	return nil
}

func GenerateEtcdCSRs(ctx context.Context, certs map[string]CertificatePKI, rkeConfig v3.RancherKubernetesEngineConfig) error {
	kubernetesServiceIP, err := GetKubernetesServiceIP(rkeConfig.Services.KubeAPI.ServiceClusterIPRange)
	if err != nil {
		return fmt.Errorf("Failed to get Kubernetes Service IP: %v", err)
	}
	clusterDomain := rkeConfig.Services.Kubelet.ClusterDomain
	etcdHosts := hosts.NodesToHosts(rkeConfig.Nodes, etcdRole)
	etcdAltNames := GetAltNames(etcdHosts, clusterDomain, kubernetesServiceIP, []string{})
	for _, host := range etcdHosts {
		etcdName := GetEtcdCrtName(host.InternalAddress)
		etcdCrt := certs[etcdName].Certificate
		etcdCSRPEM := certs[etcdName].CSRPEM
		if etcdCSRPEM != "" {
			return nil
		}
		log.Infof(ctx, "[certificates] Generating etcd-%s csr and key", host.InternalAddress)
		etcdCSR, etcdKey, err := GenerateCertSigningRequestAndKey(true, EtcdCertName, etcdAltNames, certs[etcdName].Key, nil)
		if err != nil {
			return err
		}
		certs[etcdName] = ToCertObject(etcdName, "", "", etcdCrt, etcdKey, etcdCSR)
	}
	return nil
}

func GenerateServiceTokenKey(ctx context.Context, certs map[string]CertificatePKI, rkeConfig v3.RancherKubernetesEngineConfig, configPath, configDir string, rotate bool) error {
	// generate service account token key
	privateAPIKey := certs[ServiceAccountTokenKeyName].Key
	caCrt := certs[CACertName].Certificate
	caKey := certs[CACertName].Key
	if certs[ServiceAccountTokenKeyName].Certificate != nil {
		return nil
	}
	// handle rotation on old clusters
	if certs[ServiceAccountTokenKeyName].Key == nil {
		privateAPIKey = certs[KubeAPICertName].Key
	}
	tokenCrt, tokenKey, err := GenerateSignedCertAndKey(caCrt, caKey, false, ServiceAccountTokenKeyName, nil, privateAPIKey, nil)
	if err != nil {
		return fmt.Errorf("Failed to generate private key for service account token: %v", err)
	}
	certs[ServiceAccountTokenKeyName] = ToCertObject(ServiceAccountTokenKeyName, ServiceAccountTokenKeyName, "", tokenCrt, tokenKey, nil)
	return nil
}

func GenerateRKECACerts(ctx context.Context, certs map[string]CertificatePKI, configPath, configDir string) error {
	// generate kubernetes CA certificate and key
	log.Infof(ctx, "[certificates] Generating CA kubernetes certificates")

	caCrt, caKey, err := GenerateCACertAndKey(CACertName, nil)
	if err != nil {
		return err
	}
	certs[CACertName] = ToCertObject(CACertName, "", "", caCrt, caKey, nil)

	// generate request header client CA certificate and key
	log.Infof(ctx, "[certificates] Generating Kubernetes API server aggregation layer requestheader client CA certificates")
	requestHeaderCACrt, requestHeaderCAKey, err := GenerateCACertAndKey(RequestHeaderCACertName, nil)
	if err != nil {
		return err
	}
	certs[RequestHeaderCACertName] = ToCertObject(RequestHeaderCACertName, "", "", requestHeaderCACrt, requestHeaderCAKey, nil)
	return nil
}

func GenerateRKEServicesCerts(ctx context.Context, certs map[string]CertificatePKI, rkeConfig v3.RancherKubernetesEngineConfig, configPath, configDir string, rotate bool) error {
	RKECerts := []GenFunc{
		GenerateKubeAPICertificate,
		GenerateServiceTokenKey,
		GenerateKubeControllerCertificate,
		GenerateKubeSchedulerCertificate,
		GenerateKubeProxyCertificate,
		GenerateKubeNodeCertificate,
		GenerateKubeAdminCertificate,
		GenerateAPIProxyClientCertificate,
		GenerateEtcdCertificates,
	}
	for _, gen := range RKECerts {
		if err := gen(ctx, certs, rkeConfig, configPath, configDir, rotate); err != nil {
			return err
		}
	}
	if len(rkeConfig.Services.Etcd.ExternalURLs) > 0 {
		return GenerateExternalEtcdCertificates(ctx, certs, rkeConfig, configPath, configDir, false)
	}
	return nil
}

func GenerateRKEServicesCSRs(ctx context.Context, certs map[string]CertificatePKI, rkeConfig v3.RancherKubernetesEngineConfig) error {
	RKECerts := []CSRFunc{
		GenerateKubeAPICSR,
		GenerateKubeControllerCSR,
		GenerateKubeSchedulerCSR,
		GenerateKubeProxyCSR,
		GenerateKubeNodeCSR,
		GenerateKubeAdminCSR,
		GenerateAPIProxyClientCSR,
		GenerateEtcdCSRs,
	}
	for _, csr := range RKECerts {
		if err := csr(ctx, certs, rkeConfig); err != nil {
			return err
		}
	}
	return nil
}
