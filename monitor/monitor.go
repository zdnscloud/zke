package monitor

import (
	"context"

	"github.com/zdnscloud/zke/core"
	"github.com/zdnscloud/zke/monitor/resources"
	"github.com/zdnscloud/zke/pkg/log"
	"github.com/zdnscloud/zke/pkg/templates"
)

const (
	PrometheusDeployJobName       = "zke-prometheus-deploy-job"
	NodeExporterDeployJobName     = "zke-nodeexporter-deploy-job"
	KubeStateMetricsDeployJobName = "zke-kubestatemetrics-deploy-job"
	AlertManagerDeployJobName     = "zke-alertmanager-deploy-job"
	GrafanaConfigmapDeployJobName = "zke-grafanaconf-deploy-job"
	GrafanaDeployJobName          = "zke-grafana-deploy-job"
	MetricServerDeployJobName     = "zke-metricsserver-deploy-job"
)

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

func DeployMonitoring(ctx context.Context, c *core.Cluster) error {
	log.Infof(ctx, "[Monitor] Setting up MonitorPlugin")
	config := map[string]interface{}{
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
	}
	// deploy metrics server
	if err := doOneDeploy(ctx, c, config, resources.MetricsServerTemplate, MetricServerDeployJobName); err != nil {
		if err, ok := err.(*core.AddonError); ok && err.IsCritical {
			return err
		}
		log.Warnf(ctx, "Failed to deploy addon execute job [MetricServer]: %v", err)
	}
	// deploy prometheus
	if err := doOneDeploy(ctx, c, config, resources.PrometheusTemplate, PrometheusDeployJobName); err != nil {
		return err
	}
	// deploy nodeexporter
	if err := doOneDeploy(ctx, c, config, resources.NodeExporterTemplate, NodeExporterDeployJobName); err != nil {
		return err
	}
	// deploy state metrics
	if err := doOneDeploy(ctx, c, config, resources.StateMetricsTemplate, KubeStateMetricsDeployJobName); err != nil {
		return err
	}
	// deploy alertmanager
	if err := doOneDeploy(ctx, c, config, resources.AlertManagerTemplate, AlertManagerDeployJobName); err != nil {
		return err
	}
	// deploy grafana configmap
	if err := c.DoAddonDeploy(ctx, resources.GrafanaConfigMapYaml, GrafanaConfigmapDeployJobName, true); err != nil {
		return err
	}
	// deploy grafana
	if err := doOneDeploy(ctx, c, config, resources.GrafanaTemplate, GrafanaDeployJobName); err != nil {
		return err
	}
	return nil
}

func doOneDeploy(ctx context.Context, c *core.Cluster, config map[string]interface{}, resourcesTemplate string, deployJobName string) error {
	configYaml, err := templates.CompileTemplateFromMap(resourcesTemplate, config)
	if err != nil {
		return err
	}

	if err := c.DoAddonDeploy(ctx, configYaml, deployJobName, true); err != nil {
		return err
	}
	return nil
}
