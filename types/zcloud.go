package types

type StorageConfig struct {
	Lvm  []Deviceconf `yaml:"lvm" json:"lvm,omitempty"`
	Nfs  []Deviceconf `yaml:"nfs" json:"nfs,omitempty"`
	Ceph []Deviceconf `yaml:"ceph" json:"ceph,omitempty"`
}

type Deviceconf struct {
	Host string   `yaml:"host" json:"host,omitempty"`
	Devs []string `yaml:"devs" json:"devs,omitempty"`
}

type MonitorConfig struct {
	// Monitoring server provider
	MetricsProvider string `yaml:"metrics_provider" json:"metrics_provider,omitempty" norman:"default=metrics-server"`
	// Metrics server options
	MetricsOptions                        map[string]string `yaml:"metrics_options" json:"metrics_options,omitempty"`
	PrometheusAlertManagerIngressEndpoint string            `yaml:"prometheus_alertmanager_ingress_endpoint" json:"prometheus_alertmanager_ingress_endpoint"`
	GrafanaIngressEndpoint                string            `yaml:"grafana_ingress_endpoint" json:"grafana_ingress_endpoint"`
	StorageTypeUse                        string            `yaml:"storage_class_use" json:"storage_class_use"`
	PrometheusDiskCapacity                string            `yaml:"prometheus_disk_capacity" json:"prometheus_disk_capacity"`
}

type RegistryConfig struct {
	Isenabled               bool   `yaml:"isenabled" json:"isenabled,omitempty"`
	RegistryIngressURL      string `yaml:"registry_ingress_url" json:"registry_ingress_url"`
	NotaryIngressURL        string `yaml:"notary_ingress_url" json:"notary_ingress_url"`
	RegistryDiskCapacity    string `yaml:"registry_disk_capacity" json:"registry_disk_capacity"`
	DatabaseDiskCapacity    string `yaml:"database_disk_capacity" json:"database_disk_capacity"`
	RedisDiskCapacity       string `yaml:"redis_disk_capacity" json:"redis_disk_capacity"`
	ChartmuseumDiskCapacity string `yaml:"Chartmuseum_disk_capacity" json:"Chartmuseum_disk_capacity"`
	JobserviceDiskCapacity  string `yaml:"jobservice_disk_capacity" json:"jobservice_disk_capacity"`
}
