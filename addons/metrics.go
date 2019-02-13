package addons

import "github.com/zdnscloud/zke/templates"

func GetMetricsServerManifest(MetricsServerConfig interface{}) (string, error) {

	return templates.CompileTemplateFromMap(templates.MetricsServerTemplate, MetricsServerConfig)
}
