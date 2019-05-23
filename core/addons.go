package core

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/zdnscloud/zke/pkg/k8s"
	"github.com/zdnscloud/zke/pkg/log"
	"github.com/zdnscloud/zke/pkg/templates"
)

const (
	NginxIngressAddonAppName = "ingress-nginx"
	CoreDNSAddonAppName      = "coredns"
)

type AddonError struct {
	err        string
	IsCritical bool
}

func (e *AddonError) Error() string {
	return e.err
}

func (c *Cluster) DoAddonDeploy(ctx context.Context, addonYaml, resourceName string, IsCritical bool) error {
	addonUpdated, err := c.StoreAddonConfigMap(ctx, addonYaml, resourceName)
	if err != nil {
		return &AddonError{fmt.Sprintf("Failed to save addon ConfigMap: %v", err), IsCritical}
	}
	log.Infof(ctx, "[addons] Executing deploy job %s", resourceName)
	k8sClient, err := k8s.NewClient(c.LocalKubeConfigPath, c.K8sWrapTransport)
	if err != nil {
		return &AddonError{fmt.Sprintf("%v", err), IsCritical}
	}
	node, err := k8s.GetNode(k8sClient, c.ControlPlaneHosts[0].HostnameOverride)
	if err != nil {
		return &AddonError{fmt.Sprintf("Failed to get Node [%s]: %v", c.ControlPlaneHosts[0].HostnameOverride, err), IsCritical}
	}
	addonJob, err := GetAddonsExecuteJob(resourceName, node.Name, c.Services.KubeAPI.Image)
	if err != nil {
		return &AddonError{fmt.Sprintf("Failed to generate addon execute job: %v", err), IsCritical}
	}

	if err = c.ApplySystemAddonExecuteJob(addonJob, addonUpdated); err != nil {
		return &AddonError{fmt.Sprintf("%v", err), IsCritical}
	}
	return nil
}

func (c *Cluster) StoreAddonConfigMap(ctx context.Context, addonYaml string, addonName string) (bool, error) {
	log.Infof(ctx, "[addons] Saving ConfigMap for addon %s to Kubernetes", addonName)
	updated := false
	kubeClient, err := k8s.NewClient(c.LocalKubeConfigPath, c.K8sWrapTransport)
	if err != nil {
		return updated, err
	}
	timeout := make(chan bool, 1)
	go func() {
		for {
			updated, err = k8s.UpdateConfigMap(kubeClient, []byte(addonYaml), addonName)
			if err != nil {
				time.Sleep(time.Second * 5)
				continue
			}
			log.Infof(ctx, "[addons] Successfully saved ConfigMap for addon %s to Kubernetes", addonName)
			timeout <- true
			break
		}
	}()
	select {
	case <-timeout:
		return updated, nil
	case <-time.After(time.Second * UpdateStateTimeout):
		return updated, fmt.Errorf("[addons] Timeout waiting for kubernetes to be ready")
	}
}

func (c *Cluster) ApplySystemAddonExecuteJob(addonJob string, addonUpdated bool) error {
	if err := k8s.ApplyK8sSystemJob(addonJob, c.LocalKubeConfigPath, c.K8sWrapTransport, c.AddonJobTimeout, addonUpdated); err != nil {
		return err
	}
	return nil
}

func GetAddonsExecuteJob(addonName, nodeName, image string) (string, error) {
	getAddonJob := func(addonName, nodeName, image string, isDelete bool) (string, error) {
		jobConfig := map[string]string{
			"AddonName": addonName,
			"NodeName":  nodeName,
			"Image":     image,
			"DeleteJob": strconv.FormatBool(isDelete),
		}
		return templates.CompileTemplateFromMap(addonJobTemplate, jobConfig)
	}
	return getAddonJob(addonName, nodeName, image, false)
}

const addonJobTemplate = `
{{- $addonName := .AddonName }}
{{- $nodeName := .NodeName }}
{{- $image := .Image }}
apiVersion: batch/v1
kind: Job
metadata:
{{- if eq .DeleteJob "true" }}
  name: {{$addonName}}-delete-job
{{- else }}
  name: {{$addonName}}-deploy-job
{{- end }}
  namespace: kube-system
spec:
  backoffLimit: 10
  template:
    metadata:
       name: zke-deploy
    spec:
        tolerations:
        - key: node-role.kubernetes.io/controlplane
          operator: Exists
          effect: NoSchedule
        - key: node-role.kubernetes.io/etcd
          operator: Exists
          effect: NoExecute
        hostNetwork: true
        serviceAccountName: zke-job-deployer
        nodeName: {{$nodeName}}
        containers:
          {{- if eq .DeleteJob "true" }}
          - name: {{$addonName}}-delete-pod
          {{- else }}
          - name: {{$addonName}}-pod
          {{- end }}
            image: {{$image}}
            {{- if eq .DeleteJob "true" }}
            command: ["/bin/sh"]
            args: ["-c" ,"kubectl get --ignore-not-found=true -f /etc/config/{{$addonName}}.yaml -o name | xargs kubectl delete --ignore-not-found=true"]
            {{- else }}
            command: [ "kubectl", "apply", "-f" , "/etc/config/{{$addonName}}.yaml"]
            {{- end }}
            volumeMounts:
            - name: config-volume
              mountPath: /etc/config
        volumes:
          - name: config-volume
            configMap:
              # Provide the name of the ConfigMap containing the files you want
              # to add to the container
              name: {{$addonName}}
              items:
                - key: {{$addonName}}
                  path: {{$addonName}}.yaml
        restartPolicy: Never`