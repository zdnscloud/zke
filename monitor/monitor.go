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
		"PrometheusAlertManagerImage":           c.SystemImages.PrometheusAlertManager,
		"PrometheusConfigMapReloaderImage":      c.SystemImages.PrometheusConfigMapReloader,
		"PrometheusNodeExporterImage":           c.SystemImages.PrometheusNodeExporter,
		"PrometheusServerImage":                 c.SystemImages.PrometheusServer,
		"GrafanaImage":                          c.SystemImages.Grafana,
		"GrafanaWatcherImage":                   c.SystemImages.GrafanaWatcher,
		"KubeStateMetricsImage":                 c.SystemImages.KubeStateMetrics,
		"MetricsServerImage":                    c.SystemImages.MetricsServer,
		"PrometheusAlertManagerIngressEndpoint": c.Monitor.PrometheusAlertManagerIngressEndpoint,
		"GrafanaIngressEndpoint":                c.Monitor.GrafanaIngressEndpoint,
		"RBACConfig":                            c.Authorization.Mode,
		"MetricsServerOptions":                  c.Monitor.MetricsOptions,
		"MetricsServerMajorVersion":             "v0.3",
		"DeployNamespace":                       DeployNamespace,
		"PrometheusDiskCapacity":                c.Monitor.PrometheusDiskCapacity,
		"StorageTypeUse":                        c.Monitor.StorageTypeUse,
	}
	k8sClient, err := k8s.GetK8sClientFromConfig(c.LocalKubeConfigPath)
	return templateConfig, k8sClient, err
}
