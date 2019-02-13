package cloudprovider

import (
	"github.com/zdnscloud/zke/cloudprovider/aws"
	"github.com/zdnscloud/zke/cloudprovider/azure"
	"github.com/zdnscloud/zke/cloudprovider/custom"
	"github.com/zdnscloud/zke/cloudprovider/openstack"
	"github.com/zdnscloud/zke/cloudprovider/vsphere"
	"github.com/rancher/types/apis/management.cattle.io/v3"
)

type CloudProvider interface {
	Init(cloudProviderConfig v3.CloudProvider) error
	GenerateCloudConfigFile() (string, error)
	GetName() string
}

func InitCloudProvider(cloudProviderConfig v3.CloudProvider) (CloudProvider, error) {
	var p CloudProvider
	if cloudProviderConfig.Name == aws.AWSCloudProviderName {
		p = aws.GetInstance()
	}
	if cloudProviderConfig.AzureCloudProvider != nil || cloudProviderConfig.Name == azure.AzureCloudProviderName {
		p = azure.GetInstance()
	}
	if cloudProviderConfig.OpenstackCloudProvider != nil || cloudProviderConfig.Name == openstack.OpenstackCloudProviderName {
		p = openstack.GetInstance()
	}
	if cloudProviderConfig.VsphereCloudProvider != nil || cloudProviderConfig.Name == vsphere.VsphereCloudProviderName {
		p = vsphere.GetInstance()
	}
	if cloudProviderConfig.CustomCloudProvider != "" {
		p = custom.GetInstance()
	}

	if p != nil {
		if err := p.Init(cloudProviderConfig); err != nil {
			return nil, err
		}
	}
	return p, nil
}
