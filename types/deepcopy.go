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
func (in *ZKEConfigNetwork) DeepCopyInto(out *ZKEConfigNetwork) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NetworkConfig.
func (in *ZKEConfigNetwork) DeepCopy() *ZKEConfigNetwork {
	if in == nil {
		return nil
	}
	out := new(ZKEConfigNetwork)
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
func (in *ZKENodePlan) DeepCopyInto(out *ZKENodePlan) {
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
func (in *ZKENodePlan) DeepCopy() *ZKENodePlan {
	if in == nil {
		return nil
	}
	out := new(ZKENodePlan)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ZKEConfigCore) DeepCopyInto(out *ZKEConfigCore) {
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
func (in *ZKEConfigCore) DeepCopy() *ZKEConfigCore {
	if in == nil {
		return nil
	}
	out := new(ZKEConfigCore)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ZKEConfigImages) DeepCopyInto(out *ZKEConfigImages) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ZKESystemImages.
func (in *ZKEConfigImages) DeepCopy() *ZKEConfigImages {
	if in == nil {
		return nil
	}
	out := new(ZKEConfigImages)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ZKEConfig) DeepCopyInto(out *ZKEConfig) {
	*out = *in
	if in.Nodes != nil {
		in, out := &in.Nodes, &out.Nodes
		*out = make([]ZKEConfigNode, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	in.Core.DeepCopyInto(&out.Core)
	in.Network.DeepCopyInto(&out.Network)
	in.Authentication.DeepCopyInto(&out.Authentication)
	out.SystemImages = in.SystemImages
	in.Authorization.DeepCopyInto(&out.Authorization)
	if in.PrivateRegistries != nil {
		in, out := &in.PrivateRegistries, &out.PrivateRegistries
		*out = make([]PrivateRegistry, len(*in))
		copy(*out, *in)
	}
	// in.Ingress.DeepCopyInto(&out.Ingress)
	in.Monitor.DeepCopyInto(&out.Monitor)
	// in.DNS.DeepCopyInto(&out.DNS)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ZcloudKubernetesEngineConfig.
func (in *ZKEConfig) DeepCopy() *ZKEConfig {
	if in == nil {
		return nil
	}
	out := new(ZKEConfig)
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
func (in *MonitorConfig) DeepCopyInto(out *MonitorConfig) {
	*out = *in
	if in.MetricsOptions != nil {
		in, out := &in.MetricsOptions, &out.MetricsOptions
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MonitoringConfig.
func (in *MonitorConfig) DeepCopy() *MonitorConfig {
	if in == nil {
		return nil
	}
	out := new(MonitorConfig)
	in.DeepCopyInto(out)
	return out
}
