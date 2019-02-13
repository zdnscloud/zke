package addons

import "github.com/zdnscloud/zke/templates"

func GetNginxIngressManifest(IngressConfig interface{}) (string, error) {

	return templates.CompileTemplateFromMap(templates.NginxIngressTemplate, IngressConfig)
}
