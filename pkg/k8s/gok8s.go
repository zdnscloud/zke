package k8s

import (
	"github.com/zdnscloud/zke/pkg/templates"

	"github.com/zdnscloud/gok8s/client"
	"github.com/zdnscloud/gok8s/client/config"
	"github.com/zdnscloud/gok8s/helper"
)

func DoDeployFromTemplate(cli client.Client, template string, templateConfig interface{}) error {
	yaml, err := templates.CompileTemplateFromMap(template, templateConfig)
	if err != nil {
		return err
	}
	return DoDeployFromYaml(cli, yaml)
}

func DoDeployFromYaml(cli client.Client, yaml string) error {
	return helper.CreateResourceFromYaml(cli, yaml)
}

func GetK8sClientFromConfig(kubeConfigPath string) (client.Client, error) {
	cfg, err := config.GetConfigFromFile(kubeConfigPath)
	if err != nil {
		return nil, err
	}
	return client.New(cfg, client.Options{})
}
