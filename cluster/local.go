package cluster

import (
	"github.com/zdnscloud/zke/services"
	"github.com/zdnscloud/zke/types"
)

func GetLocalRKEConfig() *types.RancherKubernetesEngineConfig {
	rkeLocalNode := GetLocalRKENodeConfig()
	imageDefaults := types.K8sVersionToRKESystemImages[DefaultK8sVersion]

	rkeServices := types.RKEConfigServices{
		Kubelet: types.KubeletService{
			BaseService: types.BaseService{
				Image:     imageDefaults.Kubernetes,
				ExtraArgs: map[string]string{"fail-swap-on": "false"},
			},
		},
	}
	return &types.RancherKubernetesEngineConfig{
		Nodes:    []types.RKEConfigNode{*rkeLocalNode},
		Services: rkeServices,
	}

}

func GetLocalRKENodeConfig() *types.RKEConfigNode {
	rkeLocalNode := &types.RKEConfigNode{
		Address:          LocalNodeAddress,
		HostnameOverride: LocalNodeHostname,
		User:             LocalNodeUser,
		Role:             []string{services.ControlRole, services.WorkerRole, services.ETCDRole},
	}
	return rkeLocalNode
}
