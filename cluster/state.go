package cluster

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rancher/rke/hosts"
	"github.com/rancher/rke/k8s"
	"github.com/rancher/rke/log"
	"github.com/rancher/rke/pki"
	"github.com/rancher/types/apis/management.cattle.io/v3"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"k8s.io/api/core/v1"
)

const (
	stateFileExt = ".rkestate"
	certDirExt   = "_certs"
)

type FullState struct {
	DesiredState State `json:"desiredState,omitempty"`
	CurrentState State `json:"currentState,omitempty"`
}

type State struct {
	RancherKubernetesEngineConfig *v3.RancherKubernetesEngineConfig `json:"rkeConfig,omitempty"`
	CertificatesBundle            map[string]pki.CertificatePKI     `json:"certificatesBundle,omitempty"`
}

func (c *Cluster) UpdateClusterCurrentState(ctx context.Context, fullState *FullState) error {
	fullState.CurrentState.RancherKubernetesEngineConfig = c.RancherKubernetesEngineConfig.DeepCopy()
	fullState.CurrentState.CertificatesBundle = c.Certificates
	return fullState.WriteStateFile(ctx, c.StateFilePath)
}

func (c *Cluster) GetClusterState(ctx context.Context, fullState *FullState) (*Cluster, error) {
	var err error
	if fullState.CurrentState.RancherKubernetesEngineConfig == nil {
		return nil, nil
	}

	// resetup external flags
	flags := GetExternalFlags(false, false, false, c.ConfigDir, c.ConfigPath)
	currentCluster, err := InitClusterObject(ctx, fullState.CurrentState.RancherKubernetesEngineConfig, flags)
	if err != nil {
		return nil, err
	}
	currentCluster.Certificates = fullState.CurrentState.CertificatesBundle

	// resetup dialers
	dialerOptions := hosts.GetDialerOptions(c.DockerDialerFactory, c.LocalConnDialerFactory, c.K8sWrapTransport)
	if err := currentCluster.SetupDialers(ctx, dialerOptions); err != nil {
		return nil, err
	}
	return currentCluster, nil
}

func SaveFullStateToKubernetes(ctx context.Context, kubeCluster *Cluster, fullState *FullState) error {
	k8sClient, err := k8s.NewClient(kubeCluster.LocalKubeConfigPath, kubeCluster.K8sWrapTransport)
	if err != nil {
		return fmt.Errorf("Failed to create Kubernetes Client: %v", err)
	}
	log.Infof(ctx, "[state] Saving full cluster state to Kubernetes")
	stateFile, err := json.Marshal(*fullState)
	if err != nil {
		return err
	}
	timeout := make(chan bool, 1)
	go func() {
		for {
			_, err := k8s.UpdateConfigMap(k8sClient, stateFile, FullStateConfigMapName)
			if err != nil {
				time.Sleep(time.Second * 5)
				continue
			}
			log.Infof(ctx, "[state] Successfully Saved full cluster state to Kubernetes ConfigMap: %s", StateConfigMapName)
			timeout <- true
			break
		}
	}()
	select {
	case <-timeout:
		return nil
	case <-time.After(time.Second * UpdateStateTimeout):
		return fmt.Errorf("[state] Timeout waiting for kubernetes to be ready")
	}
}

func GetStateFromKubernetes(ctx context.Context, kubeCluster *Cluster) (*Cluster, error) {
	log.Infof(ctx, "[state] Fetching cluster state from Kubernetes")
	k8sClient, err := k8s.NewClient(kubeCluster.LocalKubeConfigPath, kubeCluster.K8sWrapTransport)
	if err != nil {
		return nil, fmt.Errorf("Failed to create Kubernetes Client: %v", err)
	}
	var cfgMap *v1.ConfigMap
	var currentCluster Cluster
	timeout := make(chan bool, 1)
	go func() {
		for {
			cfgMap, err = k8s.GetConfigMap(k8sClient, StateConfigMapName)
			if err != nil {
				time.Sleep(time.Second * 5)
				continue
			}
			log.Infof(ctx, "[state] Successfully Fetched cluster state to Kubernetes ConfigMap: %s", StateConfigMapName)
			timeout <- true
			break
		}
	}()
	select {
	case <-timeout:
		clusterData := cfgMap.Data[StateConfigMapName]
		err := yaml.Unmarshal([]byte(clusterData), &currentCluster)
		if err != nil {
			return nil, fmt.Errorf("Failed to unmarshal cluster data")
		}
		return &currentCluster, nil
	case <-time.After(time.Second * GetStateTimeout):
		log.Infof(ctx, "Timed out waiting for kubernetes cluster to get state")
		return nil, fmt.Errorf("Timeout waiting for kubernetes cluster to get state")
	}
}

func GetK8sVersion(localConfigPath string, k8sWrapTransport k8s.WrapTransport) (string, error) {
	logrus.Debugf("[version] Using %s to connect to Kubernetes cluster..", localConfigPath)
	k8sClient, err := k8s.NewClient(localConfigPath, k8sWrapTransport)
	if err != nil {
		return "", fmt.Errorf("Failed to create Kubernetes Client: %v", err)
	}
	discoveryClient := k8sClient.DiscoveryClient
	logrus.Debugf("[version] Getting Kubernetes server version..")
	serverVersion, err := discoveryClient.ServerVersion()
	if err != nil {
		return "", fmt.Errorf("Failed to get Kubernetes server version: %v", err)
	}
	return fmt.Sprintf("%#v", *serverVersion), nil
}

func RebuildState(ctx context.Context, rkeConfig *v3.RancherKubernetesEngineConfig, oldState *FullState, flags ExternalFlags) (*FullState, error) {
	newState := &FullState{
		DesiredState: State{
			RancherKubernetesEngineConfig: rkeConfig.DeepCopy(),
		},
	}

	if flags.CustomCerts {
		certBundle, err := pki.ReadCertsAndKeysFromDir(flags.CertificateDir)
		if err != nil {
			return nil, fmt.Errorf("Failed to read certificates from dir [%s]: %v", flags.CertificateDir, err)
		}
		// make sure all custom certs are included
		if err := pki.ValidateBundleContent(rkeConfig, certBundle, flags.ClusterFilePath, flags.ConfigDir); err != nil {
			return nil, fmt.Errorf("Failed to validates certificates from dir [%s]: %v", flags.CertificateDir, err)
		}
		newState.DesiredState.CertificatesBundle = certBundle
		newState.CurrentState = oldState.CurrentState
		return newState, nil
	}

	// Rebuilding the certificates of the desired state
	if oldState.DesiredState.CertificatesBundle == nil {
		// Get the certificate Bundle
		certBundle, err := pki.GenerateRKECerts(ctx, *rkeConfig, "", "")
		if err != nil {
			return nil, fmt.Errorf("Failed to generate certificate bundle: %v", err)
		}
		newState.DesiredState.CertificatesBundle = certBundle
	} else {
		// Regenerating etcd certificates for any new etcd nodes
		pkiCertBundle := oldState.DesiredState.CertificatesBundle
		if err := pki.GenerateRKEServicesCerts(ctx, pkiCertBundle, *rkeConfig, flags.ClusterFilePath, flags.ConfigDir, false); err != nil {
			return nil, err
		}
		newState.DesiredState.CertificatesBundle = pkiCertBundle
	}
	newState.CurrentState = oldState.CurrentState
	return newState, nil
}

func (s *FullState) WriteStateFile(ctx context.Context, statePath string) error {
	stateFile, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("Failed to Marshal state object: %v", err)
	}
	logrus.Debugf("Writing state file: %s", stateFile)
	if err := ioutil.WriteFile(statePath, stateFile, 0640); err != nil {
		return fmt.Errorf("Failed to write state file: %v", err)
	}
	log.Infof(ctx, "Successfully Deployed state file at [%s]", statePath)
	return nil
}

func GetStateFilePath(configPath, configDir string) string {
	if configPath == "" {
		configPath = pki.ClusterConfig
	}
	baseDir := filepath.Dir(configPath)
	if len(configDir) > 0 {
		baseDir = filepath.Dir(configDir)
	}
	fileName := filepath.Base(configPath)
	baseDir += "/"
	fullPath := fmt.Sprintf("%s%s", baseDir, fileName)
	trimmedName := strings.TrimSuffix(fullPath, filepath.Ext(fullPath))
	return trimmedName + stateFileExt
}

func GetCertificateDirPath(configPath, configDir string) string {
	if configPath == "" {
		configPath = pki.ClusterConfig
	}
	baseDir := filepath.Dir(configPath)
	if len(configDir) > 0 {
		baseDir = filepath.Dir(configDir)
	}
	fileName := filepath.Base(configPath)
	baseDir += "/"
	fullPath := fmt.Sprintf("%s%s", baseDir, fileName)
	trimmedName := strings.TrimSuffix(fullPath, filepath.Ext(fullPath))
	return trimmedName + certDirExt
}

func ReadStateFile(ctx context.Context, statePath string) (*FullState, error) {
	rkeFullState := &FullState{}
	fp, err := filepath.Abs(statePath)
	if err != nil {
		return rkeFullState, fmt.Errorf("failed to lookup current directory name: %v", err)
	}
	file, err := os.Open(fp)
	if err != nil {
		return rkeFullState, fmt.Errorf("Can not find RKE state file: %v", err)
	}
	defer file.Close()
	buf, err := ioutil.ReadAll(file)
	if err != nil {
		return rkeFullState, fmt.Errorf("failed to read state file: %v", err)
	}
	if err := json.Unmarshal(buf, rkeFullState); err != nil {
		return rkeFullState, fmt.Errorf("failed to unmarshal the state file: %v", err)
	}
	rkeFullState.DesiredState.CertificatesBundle = pki.TransformPEMToObject(rkeFullState.DesiredState.CertificatesBundle)
	rkeFullState.CurrentState.CertificatesBundle = pki.TransformPEMToObject(rkeFullState.CurrentState.CertificatesBundle)
	return rkeFullState, nil
}

func removeStateFile(ctx context.Context, statePath string) {
	log.Infof(ctx, "Removing state file: %s", statePath)
	if err := os.Remove(statePath); err != nil {
		logrus.Warningf("Failed to remove state file: %v", err)
		return
	}
	log.Infof(ctx, "State file removed successfully")
}
