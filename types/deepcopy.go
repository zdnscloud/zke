package types

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AuthWebhookConfig) DeepCopyInto(out *AuthWebhookConfig) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AuthWebhookConfig.
func (in *AuthWebhookConfig) DeepCopy() *AuthWebhookConfig {
	if in == nil {
		return nil
	}
	out := new(AuthWebhookConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AuthnConfig) DeepCopyInto(out *AuthnConfig) {
	*out = *in
	if in.SANs != nil {
		in, out := &in.SANs, &out.SANs
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Webhook != nil {
		in, out := &in.Webhook, &out.Webhook
		*out = new(AuthWebhookConfig)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AuthnConfig.
func (in *AuthnConfig) DeepCopy() *AuthnConfig {
	if in == nil {
		return nil
	}
	out := new(AuthnConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AuthzConfig) DeepCopyInto(out *AuthzConfig) {
	*out = *in
	if in.Options != nil {
		in, out := &in.Options, &out.Options
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AuthzConfig.
func (in *AuthzConfig) DeepCopy() *AuthzConfig {
	if in == nil {
		return nil
	}
	out := new(AuthzConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BackupConfig) DeepCopyInto(out *BackupConfig) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BackupConfig.
func (in *BackupConfig) DeepCopy() *BackupConfig {
	if in == nil {
		return nil
	}
	out := new(BackupConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BaseService) DeepCopyInto(out *BaseService) {
	*out = *in
	if in.ExtraArgs != nil {
		in, out := &in.ExtraArgs, &out.ExtraArgs
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.ExtraBinds != nil {
		in, out := &in.ExtraBinds, &out.ExtraBinds
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.ExtraEnv != nil {
		in, out := &in.ExtraEnv, &out.ExtraEnv
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BaseService.
func (in *BaseService) DeepCopy() *BaseService {
	if in == nil {
		return nil
	}
	out := new(BaseService)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BastionHost) DeepCopyInto(out *BastionHost) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BastionHost.
func (in *BastionHost) DeepCopy() *BastionHost {
	if in == nil {
		return nil
	}
	out := new(BastionHost)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CalicoNetworkProvider) DeepCopyInto(out *CalicoNetworkProvider) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CalicoNetworkProvider.
func (in *CalicoNetworkProvider) DeepCopy() *CalicoNetworkProvider {
	if in == nil {
		return nil
	}
	out := new(CalicoNetworkProvider)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CloudProvider) DeepCopyInto(out *CloudProvider) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CloudProvider.
func (in *CloudProvider) DeepCopy() *CloudProvider {
	if in == nil {
		return nil
	}
	out := new(CloudProvider)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DNSConfig) DeepCopyInto(out *DNSConfig) {
	*out = *in
	if in.UpstreamNameservers != nil {
		in, out := &in.UpstreamNameservers, &out.UpstreamNameservers
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.ReverseCIDRs != nil {
		in, out := &in.ReverseCIDRs, &out.ReverseCIDRs
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.NodeSelector != nil {
		in, out := &in.NodeSelector, &out.NodeSelector
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DNSConfig.
func (in *DNSConfig) DeepCopy() *DNSConfig {
	if in == nil {
		return nil
	}
	out := new(DNSConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ETCDService) DeepCopyInto(out *ETCDService) {
	*out = *in
	in.BaseService.DeepCopyInto(&out.BaseService)
	if in.ExternalURLs != nil {
		in, out := &in.ExternalURLs, &out.ExternalURLs
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Snapshot != nil {
		in, out := &in.Snapshot, &out.Snapshot
		*out = new(bool)
		**out = **in
	}
	if in.BackupConfig != nil {
		in, out := &in.BackupConfig, &out.BackupConfig
		*out = new(BackupConfig)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ETCDService.
func (in *ETCDService) DeepCopy() *ETCDService {
	if in == nil {
		return nil
	}
	out := new(ETCDService)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FlannelNetworkProvider) DeepCopyInto(out *FlannelNetworkProvider) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FlannelNetworkProvider.
func (in *FlannelNetworkProvider) DeepCopy() *FlannelNetworkProvider {
	if in == nil {
		return nil
	}
	out := new(FlannelNetworkProvider)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HealthCheck) DeepCopyInto(out *HealthCheck) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HealthCheck.
func (in *HealthCheck) DeepCopy() *HealthCheck {
	if in == nil {
		return nil
	}
	out := new(HealthCheck)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IngressConfig) DeepCopyInto(out *IngressConfig) {
	*out = *in
	if in.Options != nil {
		in, out := &in.Options, &out.Options
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.NodeSelector != nil {
		in, out := &in.NodeSelector, &out.NodeSelector
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.ExtraArgs != nil {
		in, out := &in.ExtraArgs, &out.ExtraArgs
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IngressConfig.
func (in *IngressConfig) DeepCopy() *IngressConfig {
	if in == nil {
		return nil
	}
	out := new(IngressConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KubeAPIService) DeepCopyInto(out *KubeAPIService) {
	*out = *in
	in.BaseService.DeepCopyInto(&out.BaseService)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KubeAPIService.
func (in *KubeAPIService) DeepCopy() *KubeAPIService {
	if in == nil {
		return nil
	}
	out := new(KubeAPIService)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KubeControllerService) DeepCopyInto(out *KubeControllerService) {
	*out = *in
	in.BaseService.DeepCopyInto(&out.BaseService)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KubeControllerService.
func (in *KubeControllerService) DeepCopy() *KubeControllerService {
	if in == nil {
		return nil
	}
	out := new(KubeControllerService)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KubeletService) DeepCopyInto(out *KubeletService) {
	*out = *in
	in.BaseService.DeepCopyInto(&out.BaseService)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KubeletService.
func (in *KubeletService) DeepCopy() *KubeletService {
	if in == nil {
		return nil
	}
	out := new(KubeletService)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KubeproxyService) DeepCopyInto(out *KubeproxyService) {
	*out = *in
	in.BaseService.DeepCopyInto(&out.BaseService)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KubeproxyService.
func (in *KubeproxyService) DeepCopy() *KubeproxyService {
	if in == nil {
		return nil
	}
	out := new(KubeproxyService)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KubernetesServicesOptions) DeepCopyInto(out *KubernetesServicesOptions) {
	*out = *in
	if in.KubeAPI != nil {
		in, out := &in.KubeAPI, &out.KubeAPI
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Kubelet != nil {
		in, out := &in.Kubelet, &out.Kubelet
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Kubeproxy != nil {
		in, out := &in.Kubeproxy, &out.Kubeproxy
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.KubeController != nil {
		in, out := &in.KubeController, &out.KubeController
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Scheduler != nil {
		in, out := &in.Scheduler, &out.Scheduler
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KubernetesServicesOptions.
func (in *KubernetesServicesOptions) DeepCopy() *KubernetesServicesOptions {
	if in == nil {
		return nil
	}
	out := new(KubernetesServicesOptions)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NetworkConfig) DeepCopyInto(out *NetworkConfig) {
	*out = *in
	if in.Options != nil {
		in, out := &in.Options, &out.Options
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.CalicoNetworkProvider != nil {
		in, out := &in.CalicoNetworkProvider, &out.CalicoNetworkProvider
		*out = new(CalicoNetworkProvider)
		**out = **in
	}
	if in.FlannelNetworkProvider != nil {
		in, out := &in.FlannelNetworkProvider, &out.FlannelNetworkProvider
		*out = new(FlannelNetworkProvider)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NetworkConfig.
func (in *NetworkConfig) DeepCopy() *NetworkConfig {
	if in == nil {
		return nil
	}
	out := new(NetworkConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PortCheck) DeepCopyInto(out *PortCheck) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PortCheck.
func (in *PortCheck) DeepCopy() *PortCheck {
	if in == nil {
		return nil
	}
	out := new(PortCheck)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PrivateRegistry) DeepCopyInto(out *PrivateRegistry) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PrivateRegistry.
func (in *PrivateRegistry) DeepCopy() *PrivateRegistry {
	if in == nil {
		return nil
	}
	out := new(PrivateRegistry)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Process) DeepCopyInto(out *Process) {
	*out = *in
	if in.Command != nil {
		in, out := &in.Command, &out.Command
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Args != nil {
		in, out := &in.Args, &out.Args
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Env != nil {
		in, out := &in.Env, &out.Env
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.VolumesFrom != nil {
		in, out := &in.VolumesFrom, &out.VolumesFrom
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Binds != nil {
		in, out := &in.Binds, &out.Binds
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	out.HealthCheck = in.HealthCheck
	if in.Labels != nil {
		in, out := &in.Labels, &out.Labels
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Publish != nil {
		in, out := &in.Publish, &out.Publish
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Process.
func (in *Process) DeepCopy() *Process {
	if in == nil {
		return nil
	}
	out := new(Process)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ZKEConfigNode) DeepCopyInto(out *ZKEConfigNode) {
	*out = *in
	if in.Role != nil {
		in, out := &in.Role, &out.Role
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Labels != nil {
		in, out := &in.Labels, &out.Labels
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ZKEConfigNode.
func (in *ZKEConfigNode) DeepCopy() *ZKEConfigNode {
	if in == nil {
		return nil
	}
	out := new(ZKEConfigNode)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ZKEConfigNodePlan) DeepCopyInto(out *ZKEConfigNodePlan) {
	*out = *in
	if in.Processes != nil {
		in, out := &in.Processes, &out.Processes
		*out = make(map[string]Process, len(*in))
		for key, val := range *in {
			(*out)[key] = *val.DeepCopy()
		}
	}
	if in.PortChecks != nil {
		in, out := &in.PortChecks, &out.PortChecks
		*out = make([]PortCheck, len(*in))
		copy(*out, *in)
	}
	if in.Files != nil {
		in, out := &in.Files, &out.Files
		*out = make([]File, len(*in))
		copy(*out, *in)
	}
	if in.Annotations != nil {
		in, out := &in.Annotations, &out.Annotations
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Labels != nil {
		in, out := &in.Labels, &out.Labels
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ZKEConfigNodePlan.
func (in *ZKEConfigNodePlan) DeepCopy() *ZKEConfigNodePlan {
	if in == nil {
		return nil
	}
	out := new(ZKEConfigNodePlan)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ZKEConfigServices) DeepCopyInto(out *ZKEConfigServices) {
	*out = *in
	in.Etcd.DeepCopyInto(&out.Etcd)
	in.KubeAPI.DeepCopyInto(&out.KubeAPI)
	in.KubeController.DeepCopyInto(&out.KubeController)
	in.Scheduler.DeepCopyInto(&out.Scheduler)
	in.Kubelet.DeepCopyInto(&out.Kubelet)
	in.Kubeproxy.DeepCopyInto(&out.Kubeproxy)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ZKEConfigServices.
func (in *ZKEConfigServices) DeepCopy() *ZKEConfigServices {
	if in == nil {
		return nil
	}
	out := new(ZKEConfigServices)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ZKEPlan) DeepCopyInto(out *ZKEPlan) {
	*out = *in
	if in.Nodes != nil {
		in, out := &in.Nodes, &out.Nodes
		*out = make([]ZKEConfigNodePlan, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ZKEPlan.
func (in *ZKEPlan) DeepCopy() *ZKEPlan {
	if in == nil {
		return nil
	}
	out := new(ZKEPlan)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ZKESystemImages) DeepCopyInto(out *ZKESystemImages) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ZKESystemImages.
func (in *ZKESystemImages) DeepCopy() *ZKESystemImages {
	if in == nil {
		return nil
	}
	out := new(ZKESystemImages)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ZcloudKubernetesEngineConfig) DeepCopyInto(out *ZcloudKubernetesEngineConfig) {
	*out = *in
	if in.Nodes != nil {
		in, out := &in.Nodes, &out.Nodes
		*out = make([]ZKEConfigNode, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	in.Services.DeepCopyInto(&out.Services)
	in.Network.DeepCopyInto(&out.Network)
	in.Storage.DeepCopyInto(&out.Storage)
	in.Authentication.DeepCopyInto(&out.Authentication)
	if in.AddonsInclude != nil {
		in, out := &in.AddonsInclude, &out.AddonsInclude
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	out.SystemImages = in.SystemImages
	in.Authorization.DeepCopyInto(&out.Authorization)
	if in.PrivateRegistries != nil {
		in, out := &in.PrivateRegistries, &out.PrivateRegistries
		*out = make([]PrivateRegistry, len(*in))
		copy(*out, *in)
	}
	in.Ingress.DeepCopyInto(&out.Ingress)
	in.CloudProvider.DeepCopyInto(&out.CloudProvider)
	out.BastionHost = in.BastionHost
	in.Monitoring.DeepCopyInto(&out.Monitoring)
	out.Restore = in.Restore
	in.DNS.DeepCopyInto(&out.DNS)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ZcloudKubernetesEngineConfig.
func (in *ZcloudKubernetesEngineConfig) DeepCopy() *ZcloudKubernetesEngineConfig {
	if in == nil {
		return nil
	}
	out := new(ZcloudKubernetesEngineConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RestoreConfig) DeepCopyInto(out *RestoreConfig) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RestoreConfig.
func (in *RestoreConfig) DeepCopy() *RestoreConfig {
	if in == nil {
		return nil
	}
	out := new(RestoreConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SchedulerService) DeepCopyInto(out *SchedulerService) {
	*out = *in
	in.BaseService.DeepCopyInto(&out.BaseService)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SchedulerService.
func (in *SchedulerService) DeepCopy() *SchedulerService {
	if in == nil {
		return nil
	}
	out := new(SchedulerService)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MonitoringConfig) DeepCopyInto(out *MonitoringConfig) {
	*out = *in
	if in.Options != nil {
		in, out := &in.Options, &out.Options
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MonitoringConfig.
func (in *MonitoringConfig) DeepCopy() *MonitoringConfig {
	if in == nil {
		return nil
	}
	out := new(MonitoringConfig)
	in.DeepCopyInto(out)
	return out
}

//
func (in *StorageConfig) DeepCopyInto(out *StorageConfig) {
	*out = *in
	if in.Lvm != nil {
		in, out := &in.Lvm, &out.Lvm
		*out = make([]Lvmconf, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

func (in *StorageConfig) DeepCopy() *StorageConfig {
	if in == nil {
		return nil
	}
	out := new(StorageConfig)
	in.DeepCopyInto(out)
	return out
}

func (in *Lvmconf) DeepCopyInto(out *Lvmconf) {
	*out = *in
	if in.Devs != nil {
		in, out := &in.Devs, &out.Devs
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

func (in *Lvmconf) DeepCopy() *Lvmconf {
	if in == nil {
		return nil
	}
	out := new(Lvmconf)
	in.DeepCopyInto(out)
	return out
}
