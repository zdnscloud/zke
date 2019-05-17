package monitoring

import (
	"context"
	"strings"

	"github.com/zdnscloud/zke/core"
	"github.com/zdnscloud/zke/monitoring/resources"
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
	MonitoringPreDeployJobName    = "zke-monitoring-pre-deploy-job"

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

type metricsServerOptions struct {
	RBACConfig         string
	Options            map[string]string
	MetricsServerImage string
	Version            string
}

func DeployMonitoring(ctx context.Context, c *core.Cluster) error {
	log.Infof(ctx, "[Monitor] Setting up MonitoringPlugin")
	if err := doMonitoringPreDeploy(ctx, c); err != nil {
		return err
	}

	if err := doMetricServerDeploy(ctx, c); err != nil {
		if err, ok := err.(*core.AddonError); ok && err.IsCritical {
			return err
		}
		log.Warnf(ctx, "Failed to deploy addon execute job [MetricServer]: %v", err)
	}

	if err := doPrometheusDeploy(ctx, c); err != nil {
		return err
	}

	if err := doNodeExporterDeploy(ctx, c); err != nil {
		return err
	}

	if err := doStateMetricsDeploy(ctx, c); err != nil {
		return err
	}

	if err := doAlertManagerDeploy(ctx, c); err != nil {
		return err
	}

	if err := doGrafanaDeploy(ctx, c); err != nil {
		return err
	}
	return nil
}

func doPrometheusDeploy(ctx context.Context, c *core.Cluster) error {
	config := map[string]interface{}{
		core.RBACConfig:                  c.Authorization.Mode,
		PermetheusServerImage:            c.SystemImages.PrometheusServer,
		PrometheusConfigMapReloaderImage: c.SystemImages.PrometheusConfigMapReloader,
	}
	prometheusYaml, err := templates.CompileTemplateFromMap(resources.PrometheusTemplate, config)
	if err != nil {
		return err
	}
	if err := c.DoAddonDeploy(ctx, prometheusYaml, PrometheusDeployJobName, true); err != nil {
		return err
	}
	return nil
}

func doNodeExporterDeploy(ctx context.Context, c *core.Cluster) error {
	config := map[string]interface{}{
		PermetheusNodeExporterImage: c.SystemImages.PrometheusNodeExporter,
	}
	configYaml, err := templates.CompileTemplateFromMap(resources.NodeExporterTemplate, config)
	if err != nil {
		return err
	}
	if err := c.DoAddonDeploy(ctx, configYaml, NodeExporterDeployJobName, true); err != nil {
		return err
	}
	return nil
}

func doStateMetricsDeploy(ctx context.Context, c *core.Cluster) error {
	config := map[string]interface{}{
		core.RBACConfig:       c.Authorization.Mode,
		KubeStateMetricsImage: c.SystemImages.KubeStateMetrics,
	}
	configYaml, err := templates.CompileTemplateFromMap(resources.StateMetricsTemplate, config)
	if err != nil {
		return err
	}
	if err := c.DoAddonDeploy(ctx, configYaml, KubeStateMetricsDeployJobName, true); err != nil {
		return err
	}
	return nil
}

func doAlertManagerDeploy(ctx context.Context, c *core.Cluster) error {
	config := map[string]interface{}{
		core.RBACConfig:                       c.Authorization.Mode,
		PrometheusAlertManagerImage:           c.SystemImages.PrometheusAlertManager,
		PrometheusConfigMapReloaderImage:      c.SystemImages.PrometheusConfigMapReloader,
		PrometheusAlertManagerIngressEndpoint: c.Monitoring.PrometheusAlertManagerIngressEndpoint,
	}
	configYaml, err := templates.CompileTemplateFromMap(resources.AlertManagerTemplate, config)
	if err != nil {
		return err
	}
	if err := c.DoAddonDeploy(ctx, configYaml, AlertManagerDeployJobName, true); err != nil {
		return err
	}
	return nil
}

func doGrafanaDeploy(ctx context.Context, c *core.Cluster) error {
	config := map[string]interface{}{
		core.RBACConfig:        c.Authorization.Mode,
		GrafanaImage:           c.SystemImages.Grafana,
		GrafanaWatcherImage:    c.SystemImages.GrafanaWatcher,
		GrafanaIngressEndpoint: c.Monitoring.GrafanaIngressEndpoint,
	}
	if err := c.DoAddonDeploy(ctx, resources.GrafanaConfigMapYaml, GrafanaConfigmapDeployJobName, true); err != nil {
		return err
	}
	configYaml, err := templates.CompileTemplateFromMap(resources.GrafanaTemplate, config)
	if err != nil {
		return err
	}
	if err := c.DoAddonDeploy(ctx, configYaml, GrafanaDeployJobName, true); err != nil {
		return err
	}
	return nil
}

func doMetricServerDeploy(ctx context.Context, c *core.Cluster) error {
	log.Infof(ctx, "[addons] Setting up %s", c.Monitoring.MetricsProvider)
	s := strings.Split(c.SystemImages.MetricsServer, ":")
	versionTag := s[len(s)-1]
	MetricsServerConfig := metricsServerOptions{
		MetricsServerImage: c.SystemImages.MetricsServer,
		RBACConfig:         c.Authorization.Mode,
		Options:            c.Monitoring.MetricsOptions,
		Version:            core.GetTagMajorVersion(versionTag),
	}
	metricsYaml, err := templates.CompileTemplateFromMap(resources.MetricsServerTemplate, MetricsServerConfig)
	if err != nil {
		return err
	}
	if err := c.DoAddonDeploy(ctx, metricsYaml, MetricServerDeployJobName, false); err != nil {
		return err
	}
	log.Infof(ctx, "[addons] %s deployed successfully", c.Monitoring.MetricsProvider)
	return nil
}

func doMonitoringPreDeploy(ctx context.Context, c *core.Cluster) error {
	if err := c.DoAddonDeploy(ctx, resources.PreDeployYaml, MonitoringPreDeployJobName, true); err != nil {
		return err
	}
	return nil
}
