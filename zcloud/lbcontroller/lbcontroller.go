package lbcontroller

import (
	"context"

	"github.com/zdnscloud/zke/core"

	"github.com/zdnscloud/gok8s/client"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
)

const (
	deployReplicas       = 1
	deployName           = "elb-controller"
	deployNamespace      = "zcloud"
	deployServiceAccount = "zcloud-cluster-admin"
)

func CreateOrUpdate(cli client.Client, cluster *core.Cluster) error {
	deploy := appsv1.Deployment{}
	err := cli.Get(context.TODO(), k8stypes.NamespacedName{deployNamespace, deployName}, &deploy)
	if apierrors.IsNotFound(err) {
		return cli.Create(context.TODO(), genDeploy(cluster))
	}
	return cli.Update(context.TODO(), genDeploy(cluster))
}

func DeleteIfExist(cli client.Client, cluster *core.Cluster) error {
	deploy := appsv1.Deployment{}
	err := cli.Get(context.TODO(), k8stypes.NamespacedName{deployNamespace, deployName}, &deploy)
	if apierrors.IsNotFound(err) {
		return nil
	}

	if err != nil {
		return err
	}
	return cli.Delete(context.TODO(), &deploy)
}

func genDeploy(c *core.Cluster) *appsv1.Deployment {
	replicas := int32(deployReplicas)
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deployName,
			Namespace: deployNamespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": deployName,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": deployName,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						corev1.Container{
							Name:    deployName,
							Image:   c.Image.ZcloudLBController,
							Command: []string{"/elbc"},
							Args:    []string{"-masterserver", c.LoadBalance.MasterServer, "-backupserver", c.LoadBalance.BackupServer, "-user", c.LoadBalance.User, "-password", c.LoadBalance.Password, "-cluster", c.ClusterName},
						},
					},
					RestartPolicy:      corev1.RestartPolicyAlways,
					DNSPolicy:          corev1.DNSClusterFirst,
					ServiceAccountName: deployServiceAccount,
				},
			},
		},
	}
}
