package cluster

import (
	"github.com/zdnscloud/zke/services"
	"github.com/zdnscloud/zke/types"
)

func GetLocalRKEConfig() *types.ZcloudKubernetesEngineConfig {
	rkeLocalNode := GetLocalRKENodeConfig()
	imageDefaults := types.K8sVersionToZKESystemImages[DefaultK8sVersion]

	rkeServices := types.ZKEConfigServices{
		Kubelet: types.KubeletService{
			BaseService: types.BaseService{
				Image:     imageDefaults.Kubernetes,
				ExtraArgs: map[string]string{"fail-swap-on": "false"},
			},
		},
	}
	return &types.ZcloudKubernetesEngineConfig{
		Nodes:    []types.ZKEConfigNode{*rkeLocalNode},
		Services: rkeServices,
	}

}

func GetLocalRKENodeConfig() *types.ZKEConfigNode {
	rkeLocalNode := &types.ZKEConfigNode{
		Address:          LocalNodeAddress,
		HostnameOverride: LocalNodeHostname,
		User:             LocalNodeUser,
		Role:             []string{services.ControlRole, services.WorkerRole, services.ETCDRole},
	}
	return rkeLocalNode
}
