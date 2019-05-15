package main

import (
	"fmt"
	"time"

	"github.com/zdnscloud/gok8s/client/config"
	"github.com/zdnscloud/gok8s/exec"
)

func main() {
	cmd := exec.Cmd{
		Path: "/bin/sh",
	}

	pod := exec.Pod{
		Namespace:          "default",
		Name:               "kube-cmd-ben",
		Image:              "rancher/rancher-agent:v2.1.6",
		ServiceAccountName: "default",
	}

	cfg, err := config.GetConfigFromFile("/home/vagrant/.kube/config")
	if err != nil {
		panic("get cfg failed:" + err.Error())
	}

	e, err := exec.New(cfg)
	if err != nil {
		panic("create executor failed:" + err.Error())
	}

	if err := e.RunCmd(pod, cmd, nil, 30*time.Second); err != nil {
		fmt.Printf("run cmd failed:%s\n", err.Error())
	}
}
