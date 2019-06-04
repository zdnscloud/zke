package registry

import (
	"context"
	"crypto/rsa"
	"encoding/base64"

	"github.com/zdnscloud/zke/core"
	"github.com/zdnscloud/zke/core/pki"
	"github.com/zdnscloud/zke/pkg/docker"
	"github.com/zdnscloud/zke/pkg/hosts"
	"github.com/zdnscloud/zke/types"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/pkg/sftp"
	"k8s.io/client-go/util/cert"
)

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
		certTmpBasePath := "/home/" + h.User + "/docker-certs-tmp/"
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
			"registry-cert",
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
			Target:      "/docker-certs-tmp",
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
