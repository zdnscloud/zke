package templates

import (
	"bytes"
	"fmt"
	"text/template"
)

var tmpltMap map[string]string = map[string]string{
	"nginx":          NginxIngressTemplate,
	"metrics-server": MetricsServerTemplate,
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
