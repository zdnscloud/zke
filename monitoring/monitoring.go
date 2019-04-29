package monitoring

import (
	"context"
	"github.com/zdnscloud/zke/cluster"
	"github.com/zdnscloud/zke/pkg/log"
	"github.com/zdnscloud/zke/templates"
)

const (
	PrometheusResourceName        = "monitoring-prometheus"
	PrometheusDeployJobName       = "zke-prometheus-deploy-job"
	NodeExporterResourceName      = "monitoring-node-exporter"
	NodeExporterDeployJobName     = "zke-nodeexporter-deploy-job"
	KubeStateMetricsResourceName  = "monitoring-kube-state-metrics"
	KubeStateMetricsDeployJobName = "zke-kubestatemetrics-deploy-job"
	AlertManagerResourceName      = "monitoring-alertmanager"
	AlertManagerDeployJobName     = "zke-alertmanager-deploy-job"
	GrafanaConfigmapResourceName  = "monitoring-grafana-conf"
	GrafanaConfigmapDeployJobName = "zke-grafanaconf-deploy-job"
	GrafanaResourceName           = "monitoring-grafana"
	GrafanaDeployJobName          = "zke-grafana-deploy-job"

	PrometheusAlertManagerImage           = "PrometheusAlertManagerImage"
	PrometheusConfigMapReloaderImage      = "PrometheusConfigMapReloaderImage"
	PrometheusAlertManagerIngressEndpoint = "PrometheusAlertManagerIngressEndpoint"
	KubeStateMetricsImage                 = "KubeStateMetricsImage"
	PermetheusNodeExporterImage           = "PermetheusNodeExporterImage"
	PermetheusServerImage                 = "PermetheusServerImage"
	GrafanaImage                          = "GrafanaImage"
	GrafanaWatcherImage                   = "GrafanaWatcherImage"
	GrafanaIngressEndpoint                = "GrafanaIngressEndpoint"
)

func deployMonitoring(ctx context.Context, c *Cluster) error {
	log.Infof(ctx, "[Monitor] Setting up MonitoringPlugin")
	if err := c.doPrometheusDeploy(ctx); err != nil {
		return err
	}
	if err := c.doNodeExporterDeploy(ctx); err != nil {
		return err
	}
	if err := c.doKubeStateMetricsDeploy(ctx); err != nil {
		return err
	}
	if err := c.doAlertManagerDeploy(ctx); err != nil {
		return err
	}
	if err := c.doGrafanaDeploy(ctx); err != nil {
		return err
	}
	return nil
}

func doPrometheusDeploy(ctx context.Context, c *Cluster) error {
	config := map[string]interface{}{
		RBACConfig:                       c.Authorization.Mode,
		PermetheusServerImage:            c.SystemImages.PrometheusServer,
		PrometheusConfigMapReloaderImage: c.SystemImages.PrometheusConfigMapReloader,
	}
	prometheusYaml, err := templates.GetManifest(config, PrometheusResourceName)
	if err != nil {
		return err
	}
	if err := c.doAddonDeploy(ctx, prometheusYaml, PrometheusDeployJobName, true); err != nil {
		return err
	}
	return nil
}

func doNodeExporterDeploy(ctx context.Context, c *Cluster) error {
	config := map[string]interface{}{
		PermetheusNodeExporterImage: c.SystemImages.PrometheusNodeExporter,
	}
	configYaml, err := templates.GetManifest(config, NodeExporterResourceName)
	if err != nil {
		return err
	}
	if err := c.doAddonDeploy(ctx, configYaml, NodeExporterDeployJobName, true); err != nil {
		return err
	}
	return nil
}

func doKubeStateMetricsDeploy(ctx context.Context, c *Cluster) error {
	config := map[string]interface{}{
		RBACConfig:            c.Authorization.Mode,
		KubeStateMetricsImage: c.SystemImages.KubeStateMetrics,
	}
	configYaml, err := templates.GetManifest(config, KubeStateMetricsResourceName)
	if err != nil {
		return err
	}
	if err := c.doAddonDeploy(ctx, configYaml, KubeStateMetricsDeployJobName, true); err != nil {
		return err
	}
	return nil
}

func doAlertManagerDeploy(ctx context.Context, c *Cluster) error {
	config := map[string]interface{}{
		RBACConfig:                            c.Authorization.Mode,
		PrometheusAlertManagerImage:           c.SystemImages.PrometheusAlertManager,
		PrometheusConfigMapReloaderImage:      c.SystemImages.PrometheusConfigMapReloader,
		PrometheusAlertManagerIngressEndpoint: c.Monitoring.PrometheusAlertManagerIngressEndpoint,
	}
	configYaml, err := templates.GetManifest(config, AlertManagerResourceName)
	if err != nil {
		return err
	}
	if err := c.doAddonDeploy(ctx, configYaml, AlertManagerDeployJobName, true); err != nil {
		return err
	}
	return nil
}

func doGrafanaDeploy(ctx context.Context, c *Cluster) error {
	config := map[string]interface{}{
		RBACConfig:             c.Authorization.Mode,
		GrafanaImage:           c.SystemImages.Grafana,
		GrafanaWatcherImage:    c.SystemImages.GrafanaWatcher,
		GrafanaIngressEndpoint: c.Monitoring.GrafanaIngressEndpoint,
	}
	GrafanaConfigmapYaml := templates.GrafanaConfigMapTemplate
	if err := c.doAddonDeploy(ctx, GrafanaConfigmapYaml, GrafanaConfigmapDeployJobName, true); err != nil {
		return err
	}
	configYaml, err := templates.GetManifest(config, GrafanaResourceName)
	if err != nil {
		return err
	}
	if err := c.doAddonDeploy(ctx, configYaml, GrafanaDeployJobName, true); err != nil {
		return err
	}
	return nil
}
