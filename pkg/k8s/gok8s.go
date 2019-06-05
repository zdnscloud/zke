package k8s

import (
	"fmt"

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
	err = DoDeployFromYaml(cli, yaml)
	if err != nil {
		fmt.Println("===========\n", yaml, "\n=========\n")
	}
	return err

}

func DoDeployFromYaml(cli client.Client, yaml string) error {
	err := helper.CreateResourceFromYaml(cli, yaml)
	return err
}

func GetK8sClientFromConfig(kubeConfigPath string) (client.Client, error) {
	cfg, err := config.GetConfigFromFile(kubeConfigPath)
	if err != nil {
		return nil, err
	}
	cli, err := client.New(cfg, client.Options{})
	return cli, err
}
