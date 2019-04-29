package templates

import (
	"bytes"
	"fmt"
	"text/template"
)

var tmpltMap map[string]string = map[string]string{
	"coredns":                       CoreDNSTemplate,
	"nginx":                         NginxIngressTemplate,
	"metrics-server":                MetricsServerTemplate,
	"monitoring-prometheus":         PrometheusTemplate,
	"monitoring-alertmanager":       AlertManagerTemplate,
	"monitoring-node-exporter":      NodeExporterTemplate,
	"monitoring-kube-state-metrics": KubeStateMetricsTemplate,
	"monitoring-grafana-conf":       GrafanaConfigMapTemplate,
	"monitoring-grafana":            GrafanaTemplate,
}

func CompileTemplateFromMap(tmplt string, configMap interface{}) (string, error) {
	out := new(bytes.Buffer)
	t := template.Must(template.New("compiled_template").Parse(tmplt))
	if err := t.Execute(out, configMap); err != nil {
		return "", err
	}
	return out.String(), nil
}

func GetManifest(Config interface{}, addonName string, v ...string) (string, error) {
	tmplt, ok := tmpltMap[addonName]
	if ok {
		return CompileTemplateFromMap(tmplt, Config)
	} else {
		return "", fmt.Errorf("[addon] Unknown addon: %s", addonName)
	}
}
