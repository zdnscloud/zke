package monitor

import (
	"context"

	"github.com/zdnscloud/zke/core"
	"github.com/zdnscloud/zke/monitor/components"
	"github.com/zdnscloud/zke/pkg/k8s"
	"github.com/zdnscloud/zke/pkg/log"

	"github.com/zdnscloud/gok8s/client"
)

const DeployNamespace = "zcloud"

type MonitorImage struct {
	PrometheusAlertManager      string `yaml:"prometheus_alert_manager" json:"prometheus_alert_manager"`
	PrometheusConfigMapReloader string `yaml:"prometheus_configmap_reloader" json:"prometheus_configmap_reloader"`
	PrometheusNodeExporter      string `yaml:"prometheus_nodeexporter" json:"prometheus_nodeexporter"`
	PrometheusServer            string `yaml:"prometheus_server" json:"prometheus_server"`
	Grafana                     string `yaml:"grafana" json:"grafana"`
	GrafanaWatcher              string `yaml:"grafana_watcher" json:"grafana_watcher"`
	KubeStateMetrics            string `yaml:"kube_state_metrics" json:"kube_state_metrics"`
	MetricsServer               string `yaml:"metrics_server" json:"metricsServer,omitempty"`
}

var DefaultMonitorImage = MonitorImage{
	PrometheusAlertManager:      "zdnscloud/prometheus-alertmanager:v0.14.0",
	PrometheusConfigMapReloader: "zdnscloud/prometheus-configmap-reload:v0.1",
	PrometheusNodeExporter:      "zdnscloud/prometheus-node-exporter:v0.15.2",
	PrometheusServer:            "zdnscloud/prometheus:v2.2.1",
	Grafana:                     "zdnscloud/grafana:5.0.0",
	GrafanaWatcher:              "zdnscloud/grafana-watcher:v0.0.8",
	KubeStateMetrics:            "zdnscloud/kube-state-metrics:v1.3.1",
	MetricsServer:               "zdnscloud/metrics-server-amd64:v0.3.1",
}

var componentsTemplates = map[string]string{
	"metrics-server": components.MetricsServerTemplate,
	"prometheus":     components.PrometheusTemplate,
	"node-exporter":  components.NodeExporterTemplate,
	"state-metrics":  components.StateMetricsTemplate,
	"alert-manager":  components.AlertManagerTemplate,
	"grafana":        components.GrafanaTemplate,
}

func DeployMonitoring(ctx context.Context, c *core.Cluster) error {
	log.Infof(ctx, "[Monitor] Setting up Monitor Plugin")
	templateConfig, k8sClient, err := prepare(c)
	if err != nil {
		return err
	}
	err = k8s.DoDeployFromYaml(k8sClient, components.GrafanaConfigMapYaml)
	if err != nil {
		return err
	}
	for component, template := range componentsTemplates {
		err := k8s.DoDeployFromTemplate(k8sClient, template, templateConfig)
		if err != nil {
			log.Infof(ctx, "[Monitor] component %s deploy failed", component)
			return err
		}
	}

	log.Infof(ctx, "[Monitor] Successfully deployed Monitor Plugin")
	return nil
}

func prepare(c *core.Cluster) (map[string]interface{}, client.Client, error) {
	templateConfig := map[string]interface{}{
		"PrometheusAlertManagerImage":           DefaultMonitorImage.PrometheusAlertManager,
		"PrometheusConfigMapReloaderImage":      DefaultMonitorImage.PrometheusConfigMapReloader,
		"PrometheusNodeExporterImage":           DefaultMonitorImage.PrometheusNodeExporter,
		"PrometheusServerImage":                 DefaultMonitorImage.PrometheusServer,
		"GrafanaImage":                          DefaultMonitorImage.Grafana,
		"GrafanaWatcherImage":                   DefaultMonitorImage.GrafanaWatcher,
		"KubeStateMetricsImage":                 DefaultMonitorImage.KubeStateMetrics,
		"MetricsServerImage":                    DefaultMonitorImage.MetricsServer,
		"PrometheusAlertManagerIngressEndpoint": c.Monitor.PrometheusAlertManagerIngressEndpoint,
		"GrafanaIngressEndpoint":                c.Monitor.GrafanaIngressEndpoint,
		"RBACConfig":                            c.Authorization.Mode,
		"MetricsServerOptions":                  c.Monitor.MetricsOptions,
		"MetricsServerMajorVersion":             "v0.3",
		"DeployNamespace":                       DeployNamespace,
	}
	k8sClient, err := k8s.GetK8sClientFromConfig("./kube_config_cluster.yml")
	return templateConfig, k8sClient, err
}
