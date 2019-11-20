package linkerd

import (
	"time"

	"github.com/golang/protobuf/ptypes"

	"github.com/zdnscloud/zke/zcloud/linkerd/pkg/charts"
	"github.com/zdnscloud/zke/zcloud/linkerd/pkg/config"
	"github.com/zdnscloud/zke/zcloud/linkerd/pkg/pbconfig"
	"github.com/zdnscloud/zke/zcloud/linkerd/pkg/version"
)

const (
	defaultControllerReplicas         = 1
	defaultControllerLogLevel         = "info"
	defaultControllerUID              = 2103
	defaultIdentityIssuanceLifetime   = 24 * time.Hour
	defaultIdentityClockSkewAllowance = 20 * time.Second
	defaultNamespace                  = "linkerd"
	defaultProxyImage                 = "zdnscloud/linkerd-io-proxy"
	defaultProxyInitImage             = "zdnscloud/linkerd-io-proxy-init"
	defaultImagePullPolicy            = "IfNotPresent"
	defaultProxyUID                   = 2102
	defaultProxyLogLevel              = "warn,linkerd2_proxy=info"
	defaultProxyControlPort           = 4190
	defaultProxyAdminPort             = 4191
	defaultProxyInboundPort           = 4143
	defaultProxyOutboundPort          = 4140
)

type installOptions struct {
	clusterDomain               string
	controlPlaneVersion         string
	controllerReplicas          uint
	controllerLogLevel          string
	highAvailability            bool
	controllerUID               int64
	disableH2Upgrade            bool
	disableHeartbeat            bool
	noInitContainer             bool
	skipChecks                  bool
	omitWebhookSideEffects      bool
	restrictDashboardPrivileges bool
	identityOptions             *installIdentityOptions
	*proxyConfigOptions
	recordedFlags []*pbconfig.Install_Flag
}

type proxyConfigOptions struct {
	proxyVersion           string
	proxyImage             string
	initImage              string
	initImageVersion       string
	imagePullPolicy        string
	ignoreInboundPorts     []uint
	ignoreOutboundPorts    []uint
	proxyUID               int64
	proxyLogLevel          string
	proxyInboundPort       uint
	proxyOutboundPort      uint
	proxyControlPort       uint
	proxyAdminPort         uint
	proxyCPURequest        string
	proxyMemoryRequest     string
	proxyCPULimit          string
	proxyMemoryLimit       string
	enableExternalProfiles bool
	ignoreCluster          bool
	disableIdentity        bool
	disableTap             bool
}

func newInstallOptionsWithDefaults(clusterDomain string) *InstallOptions {
	return &InstallOptions{
		clusterDomain:       clusterDomain,
		controlPlaneVersion: version.Version,
		controllerReplicas:  defaultControllerReplicas,
		controllerLogLevel:  defaultControllerLogLevel,
		controllerUID:       defaultControllerUID,
		disableH2Upgrade:    true,
		disableHeartbeat:    true,
		proxyConfigOptions: &ProxyConfigOptions{
			proxyVersion:      version.Version,
			proxyImage:        defaultProxyImage,
			initImage:         defaultProxyInitImage,
			initImageVersion:  version.ProxyInitVersion,
			imagePullPolicy:   defaultImagePullPolicy,
			proxyUID:          defaultProxyUID,
			proxyLogLevel:     defaultProxyLogLevel,
			proxyControlPort:  defaultProxyControlPort,
			proxyAdminPort:    defaultProxyAdminPort,
			proxyInboundPort:  defaultProxyInboundPort,
			proxyOutboundPort: defaultProxyOutboundPort,
		},
		identityOptions: &InstallIdentityOptions{
			trustDomain:        clusterDomain,
			issuanceLifetime:   defaultIdentityIssuanceLifetime,
			clockSkewAllowance: defaultIdentityClockSkewAllowance,
		},
	}
}

func (options *installOptions) validateAndBuild() (*charts.Values, *pbconfig.All, error) {
	identityValues, err := options.identityOptions.genValues()
	if err != nil {
		return nil, nil, err
	}

	configs := options.configs(toIdentityContext(identityValues))
	values, err := options.buildValuesWithoutIdentity(configs)
	if err != nil {
		return nil, nil, err
	}

	values.Identity = identityValues
	return values, configs, nil
}

func toIdentityContext(idvals *charts.Identity) *pbconfig.IdentityContext {
	if idvals == nil {
		return nil
	}

	il, err := time.ParseDuration(idvals.Issuer.IssuanceLifetime)
	if err != nil {
		il = defaultIdentityIssuanceLifetime
	}

	csa, err := time.ParseDuration(idvals.Issuer.ClockSkewAllowance)
	if err != nil {
		csa = defaultIdentityClockSkewAllowance
	}

	return &pbconfig.IdentityContext{
		TrustDomain:        idvals.TrustDomain,
		TrustAnchorsPem:    idvals.TrustAnchorsPEM,
		IssuanceLifetime:   ptypes.DurationProto(il),
		ClockSkewAllowance: ptypes.DurationProto(csa),
	}
}

func (options *installOptions) buildValuesWithoutIdentity(configs *pbconfig.All) (*charts.Values, error) {
	installValues, err := charts.NewValues(false)
	if err != nil {
		return nil, err
	}

	globalJSON, proxyJSON, installJSON, err := config.ToJSON(configs)
	if err != nil {
		return nil, err
	}

	installValues.ClusterDomain = options.clusterDomain
	installValues.Configs.Global = globalJSON
	installValues.Configs.Proxy = proxyJSON
	installValues.Configs.Install = installJSON
	installValues.UUID = configs.GetInstall().GetUuid()
	return installValues, nil
}

func (options *installOptions) configs(identity *pbconfig.IdentityContext) *pbconfig.All {
	return &pbconfig.All{
		Global:  options.globalConfig(identity),
		Proxy:   options.proxyConfig(),
		Install: options.installConfig(),
	}
}

func (options *installOptions) globalConfig(identity *pbconfig.IdentityContext) *pbconfig.Global {
	return &pbconfig.Global{
		LinkerdNamespace:       defaultNamespace,
		CniEnabled:             options.noInitContainer,
		Version:                options.controlPlaneVersion,
		IdentityContext:        identity,
		OmitWebhookSideEffects: options.omitWebhookSideEffects,
		ClusterDomain:          options.clusterDomain,
	}
}

func (options *installOptions) installConfig() *pbconfig.Install {
	installID := ""
	if id, err := uuid.NewRandom(); err == nil {
		installID = id.String()
	}

	return &pbconfig.Install{
		Uuid:       installID,
		CliVersion: version.Version,
		Flags:      options.recordedFlags,
	}
}

func (options *installOptions) proxyConfig() *pbconfig.Proxy {
	ignoreInboundPorts := []*pbconfig.Port{}
	for _, port := range options.ignoreInboundPorts {
		ignoreInboundPorts = append(ignoreInboundPorts, &pbconfig.Port{Port: uint32(port)})
	}

	ignoreOutboundPorts := []*pbconfig.Port{}
	for _, port := range options.ignoreOutboundPorts {
		ignoreOutboundPorts = append(ignoreOutboundPorts, &pbconfig.Port{Port: uint32(port)})
	}

	return &pbconfig.Proxy{
		ProxyImage: &pbconfig.Image{
			ImageName:  options.proxyImage,
			PullPolicy: options.imagePullPolicy,
		},
		ProxyInitImage: &pbconfig.Image{
			ImageName:  options.initImage,
			PullPolicy: options.imagePullPolicy,
		},
		ControlPort: &pbconfig.Port{
			Port: uint32(options.proxyControlPort),
		},
		IgnoreInboundPorts:  ignoreInboundPorts,
		IgnoreOutboundPorts: ignoreOutboundPorts,
		InboundPort: &pbconfig.Port{
			Port: uint32(options.proxyInboundPort),
		},
		AdminPort: &pbconfig.Port{
			Port: uint32(options.proxyAdminPort),
		},
		OutboundPort: &pbconfig.Port{
			Port: uint32(options.proxyOutboundPort),
		},
		Resource: &pbconfig.ResourceRequirements{
			RequestCpu:    options.proxyCPURequest,
			RequestMemory: options.proxyMemoryRequest,
			LimitCpu:      options.proxyCPULimit,
			LimitMemory:   options.proxyMemoryLimit,
		},
		ProxyUid: options.proxyUID,
		LogLevel: &pbconfig.LogLevel{
			Level: options.proxyLogLevel,
		},
		DisableExternalProfiles: !options.enableExternalProfiles,
		ProxyVersion:            options.proxyVersion,
		ProxyInitImageVersion:   options.initImageVersion,
	}
}
