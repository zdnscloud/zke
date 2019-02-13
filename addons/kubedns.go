package addons

import "github.com/zdnscloud/zke/templates"

func GetKubeDNSManifest(KubeDNSConfig interface{}) (string, error) {

	return templates.CompileTemplateFromMap(templates.KubeDNSTemplate, KubeDNSConfig)
}
