package linkerd

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/zdnscloud/zke/zcloud/linkerd/pkg/charts"
	"github.com/zdnscloud/zke/zcloud/linkerd/pkg/k8s"
	"github.com/zdnscloud/zke/zcloud/linkerd/pkg/tls"
)

type installIdentityOptions struct {
	replicas           uint
	trustDomain        string
	trustPEMFile       string
	crtPEMFilei        string
	keyPEMFile         string
	issuanceLifetime   time.Duration
	clockSkewAllowance time.Duration
}

func (idopts *installIdentityOptions) issuerName() string {
	return fmt.Sprintf("identity.%s.%s", defaultNamespace, idopts.trustDomain)
}

func (idopts *installIdentityOptions) genValues() (*charts.Identity, error) {
	root, err := tls.GenerateRootCAWithDefaults(idopts.issuerName())
	if err != nil {
		return nil, fmt.Errorf("failed to generate root certificate for identity: %s", err)
	}

	return &charts.Identity{
		TrustDomain:     idopts.trustDomain,
		TrustAnchorsPEM: root.Cred.Crt.EncodeCertificatePEM(),
		Issuer: &charts.Issuer{
			ClockSkewAllowance:  idopts.clockSkewAllowance.String(),
			IssuanceLifetime:    idopts.issuanceLifetime.String(),
			CrtExpiry:           root.Cred.Crt.Certificate.NotAfter,
			CrtExpiryAnnotation: k8s.IdentityIssuerExpiryAnnotation,
			TLS: &charts.TLS{
				KeyPEM: root.Cred.EncodePrivateKeyPEM(),
				CrtPEM: root.Cred.Crt.EncodeCertificatePEM(),
			},
		},
	}, nil
}

func (idopts *installIdentityOptions) readValues() (*charts.Identity, error) {
	creds, err := tls.ReadPEMCreds(idopts.keyPEMFile, idopts.crtPEMFile)
	if err != nil {
		return nil, err
	}

	trustb, err := ioutil.ReadFile(idopts.trustPEMFile)
	if err != nil {
		return nil, err
	}
	trustAnchorsPEM := string(trustb)
	roots, err := tls.DecodePEMCertPool(trustAnchorsPEM)
	if err != nil {
		return nil, err
	}

	if err := creds.Verify(roots, idopts.issuerName()); err != nil {
		return nil, fmt.Errorf("invalid credentials: %s", err)
	}

	return &charts.Identity{
		TrustDomain:     idopts.trustDomain,
		TrustAnchorsPEM: trustAnchorsPEM,
		Issuer: &charts.Issuer{
			ClockSkewAllowance:  idopts.clockSkewAllowance.String(),
			IssuanceLifetime:    idopts.issuanceLifetime.String(),
			CrtExpiry:           creds.Crt.Certificate.NotAfter,
			CrtExpiryAnnotation: k8s.IdentityIssuerExpiryAnnotation,
			TLS: &charts.TLS{
				KeyPEM: creds.EncodePrivateKeyPEM(),
				CrtPEM: creds.EncodeCertificatePEM(),
			},
		},
	}, nil
}
