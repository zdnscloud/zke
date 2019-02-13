package addons

import "github.com/zdnscloud/zke/templates"

func GetCoreDNSManifest(CoreDNSConfig interface{}) (string, error) {

	return templates.CompileTemplateFromMap(templates.CoreDNSTemplate, CoreDNSConfig)
}
