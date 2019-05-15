package testenv

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"k8s.io/client-go/rest"
	"sigs.k8s.io/testing_frameworks/integration"
)

const (
	defaultStartStopTimeout = 20 * time.Second
)

var DefaultKubeAPIServerFlags = []string{
	"--etcd-servers={{ if .EtcdURL }}{{ .EtcdURL.String }}{{ end }}",
	"--cert-dir={{ .CertDir }}",
	"--insecure-port={{ if .URL }}{{ .URL.Port }}{{ end }}",
	"--insecure-bind-address={{ if .URL }}{{ .URL.Hostname }}{{ end }}",
	"--secure-port=0",
	"--admission-control=AlwaysAdmit",
}

type Environment struct {
	ControlPlane       integration.ControlPlane
	Config             *rest.Config
	KubeAPIServerFlags []string
	k8sBinPath         string
}

func NewEnv(k8sBinPath string, apiServerFlags []string) *Environment {
	return &Environment{
		KubeAPIServerFlags: apiServerFlags,
		k8sBinPath:         k8sBinPath,
	}
}

func (e *Environment) Stop() error {
	return e.ControlPlane.Stop()
}

func (e *Environment) Start() error {
	e.ControlPlane = integration.ControlPlane{}
	e.ControlPlane.APIServer = &integration.APIServer{Args: e.getAPIServerFlags()}
	e.ControlPlane.Etcd = &integration.Etcd{}

	e.ControlPlane.APIServer.Path = e.defaultAssetPath("kube-apiserver")
	e.ControlPlane.Etcd.Path = e.defaultAssetPath("etcd")
	if err := os.Setenv("TEST_ASSET_KUBECTL", e.defaultAssetPath("kubectl")); err != nil {
		return err
	}

	e.ControlPlane.Etcd.StartTimeout = defaultStartStopTimeout
	e.ControlPlane.Etcd.StopTimeout = defaultStartStopTimeout
	e.ControlPlane.APIServer.StartTimeout = defaultStartStopTimeout
	e.ControlPlane.APIServer.StopTimeout = defaultStartStopTimeout

	if err := e.startControlPlane(); err != nil {
		return err
	}

	e.Config = &rest.Config{
		Host: e.ControlPlane.APIURL().Host,
	}
	return nil
}

func (e Environment) getAPIServerFlags() []string {
	if len(e.KubeAPIServerFlags) == 0 {
		return DefaultKubeAPIServerFlags
	}
	return e.KubeAPIServerFlags
}

func (e *Environment) defaultAssetPath(binary string) string {
	return filepath.Join(e.k8sBinPath, binary)
}

func (e *Environment) startControlPlane() error {
	numTries, maxRetries := 0, 5
	for ; numTries < maxRetries; numTries++ {
		err := e.ControlPlane.Start()
		if err == nil {
			break
		}
	}

	if numTries == maxRetries {
		return fmt.Errorf("failed to start the controlplane. retried %d times", numTries)
	}
	return nil
}
