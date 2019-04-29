package templates

import (
	"bytes"
	"fmt"
	"text/template"
)

const (
	Calico = "calico"
)

var tmpltMap map[string]string = map[string]string{
	"flannel":          FlannelTemplate,
	"coredns":          CoreDNSTemplate,
	"nginx":            NginxIngressTemplate,
	"metrics-server":   MetricsServerTemplate,
	"lvm-storageclass": LVMStorageTemplate,
	"nfs-storageclass": NFSStorageTemplate,
	"zcloud-predeploy": ZcloudPreDeployTemplate,
}

var VersionedTemplate = map[string]map[string]string{
	"calico": map[string]string{
		"v1.13.1": CalicoTemplateV113,
		"default": CalicoTemplateV112,
	},
}

func CompileTemplateFromMap(tmplt string, configMap interface{}) (string, error) {
	out := new(bytes.Buffer)
	t := template.Must(template.New("compiled_template").Parse(tmplt))
	if err := t.Execute(out, configMap); err != nil {
		return "", err
	}
	return out.String(), nil
}

func GetVersionedTemplates(templateName string, k8sVersion string) string {
	versionedTemplate := VersionedTemplate[templateName]
	if _, ok := versionedTemplate[k8sVersion]; ok {
		return versionedTemplate[k8sVersion]
	}
	return versionedTemplate["default"]
}

func GetManifest(Config interface{}, addonName string, v ...string) (string, error) {
	if addonName == Calico {
		return CompileTemplateFromMap(GetVersionedTemplates(addonName, v[0]), Config)
	}
	tmplt, ok := tmpltMap[addonName]
	if ok {
		return CompileTemplateFromMap(tmplt, Config)
	} else {
		return "", fmt.Errorf("[addon] Unknown addon: %s", addonName)
	}
}
