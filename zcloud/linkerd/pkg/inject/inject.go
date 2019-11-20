package inject

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	k8sResource "k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/helm/pkg/chartutil"
	"sigs.k8s.io/yaml"

	"github.com/zdnscloud/zke/zcloud/linkerd/pkg/charts"
	"github.com/zdnscloud/zke/zcloud/linkerd/pkg/k8s"
	config "github.com/zdnscloud/zke/zcloud/linkerd/pkg/pbconfig"
	"github.com/zdnscloud/zke/zcloud/linkerd/pkg/version"
)

const (
	proxyInitResourceRequestCPU    = "10m"
	proxyInitResourceRequestMemory = "10Mi"
	proxyInitResourceLimitCPU      = "100m"
	proxyInitResourceLimitMemory   = "50Mi"

	traceDefaultSvcAccount = "default"
)

var (
	rTrail = regexp.MustCompile(`\},\s*\]`)

	// ProxyAnnotations is the list of possible annotations that can be applied on a pod or namespace
	ProxyAnnotations = []string{
		k8s.ProxyAdminPortAnnotation,
		k8s.ProxyControlPortAnnotation,
		k8s.ProxyDisableIdentityAnnotation,
		k8s.ProxyDisableTapAnnotation,
		k8s.ProxyEnableDebugAnnotation,
		k8s.ProxyEnableExternalProfilesAnnotation,
		k8s.ProxyImagePullPolicyAnnotation,
		k8s.ProxyInboundPortAnnotation,
		k8s.ProxyInitImageAnnotation,
		k8s.ProxyInitImageVersionAnnotation,
		k8s.ProxyOutboundPortAnnotation,
		k8s.ProxyCPULimitAnnotation,
		k8s.ProxyCPURequestAnnotation,
		k8s.ProxyImageAnnotation,
		k8s.ProxyLogLevelAnnotation,
		k8s.ProxyMemoryLimitAnnotation,
		k8s.ProxyMemoryRequestAnnotation,
		k8s.ProxyUIDAnnotation,
		k8s.ProxyVersionOverrideAnnotation,
		k8s.ProxyIgnoreInboundPortsAnnotation,
		k8s.ProxyIgnoreOutboundPortsAnnotation,
		k8s.ProxyTraceCollectorSvcAddr,
	}
)

// Origin defines where the input YAML comes from. Refer the ResourceConfig's
// 'origin' field
type Origin int

const (
	// OriginCLI is the value of the ResourceConfig's 'origin' field if the input
	// YAML comes from the CLI
	OriginCLI Origin = iota

	// OriginWebhook is the value of the ResourceConfig's 'origin' field if the input
	// YAML comes from the CLI
	OriginWebhook

	// OriginUnknown is the value of the ResourceConfig's 'origin' field if the
	// input YAML comes from an unknown source
	OriginUnknown
)

// OwnerRetrieverFunc is a function that returns a pod's owner reference
// kind and name
type OwnerRetrieverFunc func(*corev1.Pod) (string, string)

// ResourceConfig contains the parsed information for a given workload
type ResourceConfig struct {
	configs        *config.All
	nsAnnotations  map[string]string
	ownerRetriever OwnerRetrieverFunc
	origin         Origin

	workload struct {
		obj      runtime.Object
		metaType metav1.TypeMeta

		// Meta is the workload's metadata. It's exported so that metadata of
		// non-workload resources can be unmarshalled by the YAML parser
		Meta *metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

		ownerRef *metav1.OwnerReference
	}

	pod struct {
		meta        *metav1.ObjectMeta
		labels      map[string]string
		annotations map[string]string
		spec        *corev1.PodSpec
	}
}

type patch struct {
	charts.Values
	PathPrefix            string
	AddRootAnnotations    bool
	Annotations           map[string]string
	AddRootLabels         bool
	AddRootInitContainers bool
	AddRootVolumes        bool
	Labels                map[string]string
	DebugContainer        *charts.DebugContainer
}

// NewResourceConfig creates and initializes a ResourceConfig
func NewResourceConfig(configs *config.All, origin Origin) *ResourceConfig {
	config := &ResourceConfig{
		configs: configs,
		origin:  origin,
	}

	config.pod.meta = &metav1.ObjectMeta{}
	config.pod.labels = map[string]string{k8s.ControllerNSLabel: configs.GetGlobal().GetLinkerdNamespace()}
	config.pod.annotations = map[string]string{}
	return config
}

// WithKind enriches ResourceConfig with the workload kind
func (conf *ResourceConfig) WithKind(kind string) *ResourceConfig {
	conf.workload.metaType = metav1.TypeMeta{Kind: kind}
	return conf
}

// WithNsAnnotations enriches ResourceConfig with the namespace annotations, that can
// be used in shouldInject()
func (conf *ResourceConfig) WithNsAnnotations(m map[string]string) *ResourceConfig {
	conf.nsAnnotations = m
	return conf
}

// WithOwnerRetriever enriches ResourceConfig with a function that allows to retrieve
// the kind and name of the workload's owner reference
func (conf *ResourceConfig) WithOwnerRetriever(f OwnerRetrieverFunc) *ResourceConfig {
	conf.ownerRetriever = f
	return conf
}

// GetOwnerRef returns a reference to the resource's owner resource, if any
func (conf *ResourceConfig) GetOwnerRef() *metav1.OwnerReference {
	return conf.workload.ownerRef
}

// AppendPodAnnotations appends the given annotations to the pod spec in conf
func (conf *ResourceConfig) AppendPodAnnotations(annotations map[string]string) {
	for annotation, value := range annotations {
		conf.pod.annotations[annotation] = value
	}
}

// AppendPodAnnotation appends the given single annotation to the pod spec in conf
func (conf *ResourceConfig) AppendPodAnnotation(k, v string) {
	conf.pod.annotations[k] = v
}

// YamlMarshalObj returns the yaml for the workload in conf
func (conf *ResourceConfig) YamlMarshalObj() ([]byte, error) {
	return yaml.Marshal(conf.workload.obj)
}

// ParseMetaAndYAML extracts the workload metadata and pod specs from the given
// input bytes. The results are stored in the conf's fields.
func (conf *ResourceConfig) ParseMetaAndYAML(bytes []byte) (*Report, error) {
	if err := conf.parse(bytes); err != nil {
		return nil, err
	}

	return newReport(conf), nil
}

// GetPatch returns the JSON patch containing the proxy and init containers specs, if any.
// If injectProxy is false, only the config.linkerd.io annotations are set.
func (conf *ResourceConfig) GetPatch(injectProxy bool) ([]byte, error) {
	clusterDomain := conf.configs.GetGlobal().GetClusterDomain()
	if clusterDomain == "" {
		clusterDomain = "cluster.local"
	}
	values := &patch{
		Values: charts.Values{
			Namespace:     conf.configs.GetGlobal().GetLinkerdNamespace(),
			ClusterDomain: clusterDomain,
		},
		Annotations: map[string]string{},
		Labels:      map[string]string{},
	}
	if strings.ToLower(conf.workload.metaType.Kind) != k8s.Pod {
		values.PathPrefix = "/spec/template"
	}

	if conf.pod.spec != nil {
		conf.injectPodAnnotations(values)
		if injectProxy {
			conf.injectObjectMeta(values)
			conf.injectPodSpec(values)
		}
	}

	rawValues, err := yaml.Marshal(values)
	if err != nil {
		return nil, err
	}

	files := []*chartutil.BufferedFile{
		{Name: chartutil.ChartfileName},
		{Name: "requirements.yaml"},
		{Name: "templates/patch.json"},
	}

	chart := &charts.Chart{
		Name:      "patch",
		Dir:       "patch",
		Namespace: values.Namespace,
		RawValues: rawValues,
		Files:     files,
	}
	buf, err := chart.Render()
	if err != nil {
		return nil, err
	}

	// Get rid of invalid trailing commas
	res := rTrail.ReplaceAll(buf.Bytes(), []byte("}\n]"))

	return res, nil
}

// Note this switch also defines what kinds are injectable
func (conf *ResourceConfig) getFreshWorkloadObj() runtime.Object {
	switch strings.ToLower(conf.workload.metaType.Kind) {
	case k8s.Deployment:
		return &appsv1.Deployment{}
	case k8s.ReplicationController:
		return &corev1.ReplicationController{}
	case k8s.ReplicaSet:
		return &appsv1.ReplicaSet{}
	case k8s.Job:
		return &batchv1.Job{}
	case k8s.DaemonSet:
		return &appsv1.DaemonSet{}
	case k8s.StatefulSet:
		return &appsv1.StatefulSet{}
	case k8s.Pod:
		return &corev1.Pod{}
	}

	return nil
}

// JSONToYAML is a replacement for the same function in sigs.k8s.io/yaml
// that does conserve the field order as portrayed in k8s' api structs
func (conf *ResourceConfig) JSONToYAML(bytes []byte) ([]byte, error) {
	obj := conf.getFreshWorkloadObj()
	if err := json.Unmarshal(bytes, obj); err != nil {
		return nil, err
	}
	return yaml.Marshal(obj)
}

// parse parses the bytes payload, filling the gaps in ResourceConfig
// depending on the workload kind
func (conf *ResourceConfig) parse(bytes []byte) error {
	// The Kubernetes API is versioned and each version has an API modeled
	// with its own distinct Go types. If we tell `yaml.Unmarshal()` which
	// version we support then it will provide a representation of that
	// object using the given type if possible. However, it only allows us
	// to supply one object (of one type), so first we have to determine
	// what kind of object `bytes` represents so we can pass an object of
	// the correct type to `yaml.Unmarshal()`.
	// ---------------------------------------
	// Note: bytes is expected to be YAML and will only modify it when a
	// supported type is found. Otherwise, conf is left unmodified.

	// When injecting the linkerd proxy into a linkerd controller pod. The linkerd proxy's
	// LINKERD2_PROXY_DESTINATION_SVC_ADDR variable must be set to localhost for
	// the following reasons:
	//	1. According to https://github.com/kubernetes/minikube/issues/1568, minikube has an issue
	//     where pods are unable to connect to themselves through their associated service IP.
	//     Setting the LINKERD2_PROXY_DESTINATION_SVC_ADDR to localhost allows the
	//     proxy to bypass kube DNS name resolution as a workaround to this issue.
	//  2. We avoid the TLS overhead in encrypting and decrypting intra-pod traffic i.e. traffic
	//     between containers in the same pod.
	//  3. Using a Service IP instead of localhost would mean intra-pod traffic would be load-balanced
	//     across all controller pod replicas. This is undesirable as we would want all traffic between
	//	   containers to be self contained.
	//  4. We skip recording telemetry for intra-pod traffic within the control plane.

	if err := yaml.Unmarshal(bytes, &conf.workload.metaType); err != nil {
		return err
	}
	obj := conf.getFreshWorkloadObj()

	switch v := obj.(type) {
	case *appsv1.Deployment:
		if err := yaml.Unmarshal(bytes, v); err != nil {
			return err
		}

		conf.workload.obj = v
		conf.workload.Meta = &v.ObjectMeta
		conf.pod.labels[k8s.ProxyDeploymentLabel] = v.Name
		conf.complete(&v.Spec.Template)

	case *corev1.ReplicationController:
		if err := yaml.Unmarshal(bytes, v); err != nil {
			return err
		}

		conf.workload.obj = v
		conf.workload.Meta = &v.ObjectMeta
		conf.pod.labels[k8s.ProxyReplicationControllerLabel] = v.Name
		conf.complete(v.Spec.Template)

	case *appsv1.ReplicaSet:
		if err := yaml.Unmarshal(bytes, v); err != nil {
			return err
		}

		conf.workload.obj = v
		conf.workload.Meta = &v.ObjectMeta
		conf.pod.labels[k8s.ProxyReplicaSetLabel] = v.Name
		conf.complete(&v.Spec.Template)

	case *batchv1.Job:
		if err := yaml.Unmarshal(bytes, v); err != nil {
			return err
		}

		conf.workload.obj = v
		conf.workload.Meta = &v.ObjectMeta
		conf.pod.labels[k8s.ProxyJobLabel] = v.Name
		conf.complete(&v.Spec.Template)

	case *appsv1.DaemonSet:
		if err := yaml.Unmarshal(bytes, v); err != nil {
			return err
		}

		conf.workload.obj = v
		conf.workload.Meta = &v.ObjectMeta
		conf.pod.labels[k8s.ProxyDaemonSetLabel] = v.Name
		conf.complete(&v.Spec.Template)

	case *appsv1.StatefulSet:
		if err := yaml.Unmarshal(bytes, v); err != nil {
			return err
		}

		conf.workload.obj = v
		conf.workload.Meta = &v.ObjectMeta
		conf.pod.labels[k8s.ProxyStatefulSetLabel] = v.Name
		conf.complete(&v.Spec.Template)

	case *corev1.Pod:
		if err := yaml.Unmarshal(bytes, v); err != nil {
			return err
		}

		conf.workload.obj = v
		conf.pod.spec = &v.Spec
		conf.pod.meta = &v.ObjectMeta

		if conf.ownerRetriever != nil {
			kind, name := conf.ownerRetriever(v)
			conf.workload.ownerRef = &metav1.OwnerReference{Kind: kind, Name: name}
			switch kind {
			case k8s.Deployment:
				conf.pod.labels[k8s.ProxyDeploymentLabel] = name
			case k8s.ReplicationController:
				conf.pod.labels[k8s.ProxyReplicationControllerLabel] = name
			case k8s.ReplicaSet:
				conf.pod.labels[k8s.ProxyReplicaSetLabel] = name
			case k8s.Job:
				conf.pod.labels[k8s.ProxyJobLabel] = name
			case k8s.DaemonSet:
				conf.pod.labels[k8s.ProxyDaemonSetLabel] = name
			case k8s.StatefulSet:
				conf.pod.labels[k8s.ProxyStatefulSetLabel] = name
			}
		}

	default:
		// unmarshal the metadata of other resource kinds like namespace, secret,
		// config map etc. to be used in the report struct
		if err := yaml.Unmarshal(bytes, &conf.workload); err != nil {
			return err
		}
	}

	if conf.pod.meta.Annotations == nil {
		conf.pod.meta.Annotations = map[string]string{}
	}

	return nil
}

func (conf *ResourceConfig) complete(template *corev1.PodTemplateSpec) {
	conf.pod.spec = &template.Spec
	conf.pod.meta = &template.ObjectMeta
}

// injectPodSpec adds linkerd sidecars to the provided PodSpec.
func (conf *ResourceConfig) injectPodSpec(values *patch) {
	values.Proxy = &charts.Proxy{
		Component:              conf.pod.labels[k8s.ProxyDeploymentLabel],
		EnableExternalProfiles: conf.enableExternalProfiles(),
		DisableTap:             conf.tapDisabled(),
		Image: &charts.Image{
			Name:       conf.proxyImage(),
			Version:    conf.proxyVersion(),
			PullPolicy: conf.proxyImagePullPolicy(),
		},
		LogLevel: conf.proxyLogLevel(),
		Ports: &charts.Ports{
			Admin:    conf.proxyAdminPort(),
			Control:  conf.proxyControlPort(),
			Inbound:  conf.proxyInboundPort(),
			Outbound: conf.proxyOutboundPort(),
		},
		UID:       conf.proxyUID(),
		Resources: conf.proxyResourceRequirements(),
	}

	if v := conf.pod.meta.Annotations[k8s.ProxyEnableDebugAnnotation]; v != "" {
		debug, err := strconv.ParseBool(v)
		if err != nil {
			log.Warnf("unrecognized value used for the %s annotation: %s", k8s.ProxyEnableDebugAnnotation, v)
			debug = false
		}

		if debug {
			log.Infof("inject debug container")
			values.DebugContainer = &charts.DebugContainer{
				Image: &charts.Image{
					Name:       k8s.DebugSidecarImage,
					Version:    conf.configs.GetGlobal().GetVersion(),
					PullPolicy: conf.proxyImagePullPolicy(),
				},
			}
		}
	}

	saVolumeMount := conf.serviceAccountVolumeMount()

	// use the primary container's capabilities to ensure psp compliance, if
	// enabled
	if conf.pod.spec.Containers != nil && len(conf.pod.spec.Containers) > 0 {
		if sc := conf.pod.spec.Containers[0].SecurityContext; sc != nil && sc.Capabilities != nil {
			values.Proxy.Capabilities = &charts.Capabilities{
				Add:  []string{},
				Drop: []string{},
			}
			for _, add := range sc.Capabilities.Add {
				values.Proxy.Capabilities.Add = append(values.Proxy.Capabilities.Add, string(add))
			}
			for _, drop := range sc.Capabilities.Drop {
				values.Proxy.Capabilities.Drop = append(values.Proxy.Capabilities.Drop, string(drop))
			}
		}
	}

	if saVolumeMount != nil {
		values.Proxy.SAMountPath = &charts.SAMountPath{
			Name:      saVolumeMount.Name,
			MountPath: saVolumeMount.MountPath,
			ReadOnly:  saVolumeMount.ReadOnly,
		}
	}

	if !conf.configs.GetGlobal().GetCniEnabled() {
		conf.injectProxyInit(values)
	}

	idctx := conf.identityContext()
	if idctx == nil {
		values.Proxy.DisableIdentity = true
		return
	}

	values.Identity = &charts.Identity{
		TrustAnchorsPEM: idctx.GetTrustAnchorsPem(),
		TrustDomain:     idctx.GetTrustDomain(),
	}

	values.AddRootVolumes = len(conf.pod.spec.Volumes) == 0

	if trace := conf.trace(); trace != nil {
		log.Infof("tracing enabled: remote service=%s, service account=%s", trace.CollectorSvcAddr, trace.CollectorSvcAccount)
		values.Proxy.Trace = trace
	}
}

func (conf *ResourceConfig) injectProxyInit(values *patch) {
	values.ProxyInit = &charts.ProxyInit{
		Image: &charts.Image{
			Name:       conf.proxyInitImage(),
			PullPolicy: conf.proxyInitImagePullPolicy(),
			Version:    conf.proxyInitVersion(),
		},
		IgnoreInboundPorts:  conf.proxyInboundSkipPorts(),
		IgnoreOutboundPorts: conf.proxyOutboundSkipPorts(),
		Resources: &charts.Resources{
			CPU: charts.Constraints{
				Limit:   proxyInitResourceLimitCPU,
				Request: proxyInitResourceRequestCPU,
			},
			Memory: charts.Constraints{
				Limit:   proxyInitResourceLimitMemory,
				Request: proxyInitResourceRequestMemory,
			},
		},
		Capabilities: values.Proxy.Capabilities,
		SAMountPath:  values.Proxy.SAMountPath,
	}

	values.AddRootInitContainers = len(conf.pod.spec.InitContainers) == 0

}

func (conf *ResourceConfig) serviceAccountVolumeMount() *corev1.VolumeMount {
	// Probably always true, but wanna be super-safe
	if containers := conf.pod.spec.Containers; len(containers) > 0 {
		for _, vm := range containers[0].VolumeMounts {
			if vm.MountPath == k8s.MountPathServiceAccount {
				vm := vm // pin
				return &vm
			}
		}
	}
	return nil
}

func (conf *ResourceConfig) trace() *charts.Trace {
	var (
		svcAddr    = conf.getOverride(k8s.ProxyTraceCollectorSvcAddr)
		svcAccount = conf.getOverride(k8s.ProxyTraceCollectorSvcAccount)
	)

	if svcAddr == "" {
		return nil
	}

	if svcAccount == "" {
		svcAccount = traceDefaultSvcAccount
	}

	hostAndPort := strings.Split(svcAddr, ":")
	hostname := strings.Split(hostAndPort[0], ".")

	var ns string
	if len(hostname) == 1 {
		ns = conf.workload.Meta.Namespace
	} else {
		ns = hostname[1]
	}

	return &charts.Trace{
		CollectorSvcAddr:    svcAddr,
		CollectorSvcAccount: fmt.Sprintf("%s.%s", svcAccount, ns),
	}
}

// Given a ObjectMeta, update ObjectMeta in place with the new labels and
// annotations.
func (conf *ResourceConfig) injectObjectMeta(values *patch) {
	values.Annotations[k8s.ProxyVersionAnnotation] = conf.proxyVersion()

	if conf.identityContext() != nil {
		values.Annotations[k8s.IdentityModeAnnotation] = k8s.IdentityModeDefault
	} else {
		values.Annotations[k8s.IdentityModeAnnotation] = k8s.IdentityModeDisabled
	}

	if len(conf.pod.labels) > 0 {
		values.AddRootLabels = len(conf.pod.meta.Labels) == 0
		for _, k := range sortedKeys(conf.pod.labels) {
			values.Labels[k] = conf.pod.labels[k]
		}
	}
}

func (conf *ResourceConfig) injectPodAnnotations(values *patch) {
	values.AddRootAnnotations = len(conf.pod.meta.Annotations) == 0

	for _, k := range sortedKeys(conf.pod.annotations) {
		values.Annotations[k] = conf.pod.annotations[k]

		// append any additional pod annotations to the pod's meta.
		// for e.g., annotations that were converted from CLI inject options.
		conf.pod.meta.Annotations[k] = conf.pod.annotations[k]
	}
}

func (conf *ResourceConfig) getOverride(annotation string) string {
	if override := conf.pod.meta.Annotations[annotation]; override != "" {
		return override
	}
	return conf.nsAnnotations[annotation]
}

func (conf *ResourceConfig) proxyImage() string {
	if override := conf.getOverride(k8s.ProxyImageAnnotation); override != "" {
		return override
	}
	return conf.configs.GetProxy().GetProxyImage().GetImageName()
}

func (conf *ResourceConfig) proxyImagePullPolicy() string {
	if override := conf.getOverride(k8s.ProxyImagePullPolicyAnnotation); override != "" {
		return override
	}
	return conf.configs.GetProxy().GetProxyImage().GetPullPolicy()
}

func (conf *ResourceConfig) proxyVersion() string {
	if override := conf.getOverride(k8s.ProxyVersionOverrideAnnotation); override != "" {
		return override
	}
	if proxyVersion := conf.configs.GetProxy().GetProxyVersion(); proxyVersion != "" {
		return proxyVersion
	}
	if controlPlaneVersion := conf.configs.GetGlobal().GetVersion(); controlPlaneVersion != "" {
		return controlPlaneVersion
	}
	return version.Version
}

func (conf *ResourceConfig) proxyInitVersion() string {
	if override := conf.getOverride(k8s.ProxyInitImageVersionAnnotation); override != "" {
		return override
	}
	if v := conf.configs.GetProxy().GetProxyInitImageVersion(); v != "" {
		return v
	}
	return version.ProxyInitVersion
}

func (conf *ResourceConfig) proxyControlPort() int32 {
	if override := conf.getOverride(k8s.ProxyControlPortAnnotation); override != "" {
		controlPort, err := strconv.ParseInt(override, 10, 32)
		if err == nil {
			return int32(controlPort)
		}
	}

	return int32(conf.configs.GetProxy().GetControlPort().GetPort())
}

func (conf *ResourceConfig) proxyInboundPort() int32 {
	if override := conf.getOverride(k8s.ProxyInboundPortAnnotation); override != "" {
		inboundPort, err := strconv.ParseInt(override, 10, 32)
		if err == nil {
			return int32(inboundPort)
		}
	}

	return int32(conf.configs.GetProxy().GetInboundPort().GetPort())
}

func (conf *ResourceConfig) proxyAdminPort() int32 {
	if override := conf.getOverride(k8s.ProxyAdminPortAnnotation); override != "" {
		adminPort, err := strconv.ParseInt(override, 10, 32)
		if err == nil {
			return int32(adminPort)
		}
	}
	return int32(conf.configs.GetProxy().GetAdminPort().GetPort())
}

func (conf *ResourceConfig) proxyOutboundPort() int32 {
	if override := conf.getOverride(k8s.ProxyOutboundPortAnnotation); override != "" {
		outboundPort, err := strconv.ParseInt(override, 10, 32)
		if err == nil {
			return int32(outboundPort)
		}
	}

	return int32(conf.configs.GetProxy().GetOutboundPort().GetPort())
}

func (conf *ResourceConfig) proxyLogLevel() string {
	if override := conf.getOverride(k8s.ProxyLogLevelAnnotation); override != "" {
		return override
	}

	return conf.configs.GetProxy().GetLogLevel().GetLevel()
}

func (conf *ResourceConfig) identityContext() *config.IdentityContext {
	if override := conf.getOverride(k8s.ProxyDisableIdentityAnnotation); override != "" {
		value, err := strconv.ParseBool(override)
		if err == nil && value {
			return nil
		}
	}

	return conf.configs.GetGlobal().GetIdentityContext()
}

func (conf *ResourceConfig) tapDisabled() bool {
	if override := conf.getOverride(k8s.ProxyDisableTapAnnotation); override != "" {
		value, err := strconv.ParseBool(override)
		if err == nil && value {
			return true
		}
	}
	return false
}

func (conf *ResourceConfig) proxyResourceRequirements() *charts.Resources {
	var (
		requestCPU    k8sResource.Quantity
		requestMemory k8sResource.Quantity
		limitCPU      k8sResource.Quantity
		limitMemory   k8sResource.Quantity
		err           error
	)
	res := &charts.Resources{}

	if override := conf.getOverride(k8s.ProxyCPURequestAnnotation); override != "" {
		requestCPU, err = k8sResource.ParseQuantity(override)
	} else if defaultRequest := conf.configs.GetProxy().GetResource().GetRequestCpu(); defaultRequest != "" {
		requestCPU, err = k8sResource.ParseQuantity(defaultRequest)
	}
	if err != nil {
		log.Warnf("%s (%s)", err, k8s.ProxyCPURequestAnnotation)
	}
	if !requestCPU.IsZero() {
		res.CPU.Request = requestCPU.String()
	}

	if override := conf.getOverride(k8s.ProxyMemoryRequestAnnotation); override != "" {
		requestMemory, err = k8sResource.ParseQuantity(override)
	} else if defaultRequest := conf.configs.GetProxy().GetResource().GetRequestMemory(); defaultRequest != "" {
		requestMemory, err = k8sResource.ParseQuantity(defaultRequest)
	}
	if err != nil {
		log.Warnf("%s (%s)", err, k8s.ProxyMemoryRequestAnnotation)
	}
	if !requestMemory.IsZero() {
		res.Memory.Request = requestMemory.String()
	}

	if override := conf.getOverride(k8s.ProxyCPULimitAnnotation); override != "" {
		limitCPU, err = k8sResource.ParseQuantity(override)
	} else if defaultLimit := conf.configs.GetProxy().GetResource().GetLimitCpu(); defaultLimit != "" {
		limitCPU, err = k8sResource.ParseQuantity(defaultLimit)
	}
	if err != nil {
		log.Warnf("%s (%s)", err, k8s.ProxyCPULimitAnnotation)
	}
	if !limitCPU.IsZero() {
		res.CPU.Limit = limitCPU.String()
	}

	if override := conf.getOverride(k8s.ProxyMemoryLimitAnnotation); override != "" {
		limitMemory, err = k8sResource.ParseQuantity(override)
	} else if defaultLimit := conf.configs.GetProxy().GetResource().GetLimitMemory(); defaultLimit != "" {
		limitMemory, err = k8sResource.ParseQuantity(defaultLimit)
	}
	if err != nil {
		log.Warnf("%s (%s)", err, k8s.ProxyMemoryLimitAnnotation)
	}
	if !limitMemory.IsZero() {
		res.Memory.Limit = limitMemory.String()
	}

	return res
}

func (conf *ResourceConfig) proxyUID() int64 {
	if overrides := conf.getOverride(k8s.ProxyUIDAnnotation); overrides != "" {
		v, err := strconv.ParseInt(overrides, 10, 64)
		if err == nil {
			return v
		}
	}

	return conf.configs.GetProxy().GetProxyUid()
}

func (conf *ResourceConfig) enableExternalProfiles() bool {
	disableExternalProfiles := conf.configs.GetProxy().GetDisableExternalProfiles()
	if override := conf.getOverride(k8s.ProxyEnableExternalProfilesAnnotation); override != "" {
		value, err := strconv.ParseBool(override)
		if err == nil {
			return value
		}
	}

	return !disableExternalProfiles
}

func (conf *ResourceConfig) proxyInitImage() string {
	if override := conf.getOverride(k8s.ProxyInitImageAnnotation); override != "" {
		return override
	}
	return conf.configs.GetProxy().GetProxyInitImage().GetImageName()
}

func (conf *ResourceConfig) proxyInitImagePullPolicy() string {
	if override := conf.getOverride(k8s.ProxyImagePullPolicyAnnotation); override != "" {
		return override
	}
	return conf.configs.GetProxy().GetProxyInitImage().GetPullPolicy()
}

func (conf *ResourceConfig) proxyInboundSkipPorts() string {
	if override := conf.getOverride(k8s.ProxyIgnoreInboundPortsAnnotation); override != "" {
		return override
	}

	ports := []string{}
	for _, port := range conf.configs.GetProxy().GetIgnoreInboundPorts() {
		portStr := strconv.FormatUint(uint64(port.GetPort()), 10)
		ports = append(ports, portStr)
	}
	return strings.Join(ports, ",")
}

func (conf *ResourceConfig) proxyOutboundSkipPorts() string {
	if override := conf.getOverride(k8s.ProxyIgnoreOutboundPortsAnnotation); override != "" {
		return override
	}

	ports := []string{}
	for _, port := range conf.configs.GetProxy().GetIgnoreOutboundPorts() {
		portStr := strconv.FormatUint(uint64(port.GetPort()), 10)
		ports = append(ports, portStr)
	}
	return strings.Join(ports, ",")
}

// GetOverriddenConfiguration returns a map of the overridden proxy annotations
func (conf *ResourceConfig) GetOverriddenConfiguration() map[string]string {
	proxyOverrideConfig := map[string]string{}
	for _, annotation := range ProxyAnnotations {
		proxyOverrideConfig[annotation] = conf.getOverride(annotation)
	}

	return proxyOverrideConfig
}

func sortedKeys(m map[string]string) []string {
	keys := []string{}
	for k := range m {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	return keys
}
