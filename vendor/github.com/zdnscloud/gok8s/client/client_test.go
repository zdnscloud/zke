package client

import (
	"context"
	"fmt"
	"os"
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"

	ut "github.com/zdnscloud/cement/unittest"
	"github.com/zdnscloud/gok8s/testenv"
)

func newDeploy(count int, ns string) *appsv1.Deployment {
	var replicaCount int32 = 2
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{Name: fmt.Sprintf("deployment-name-%v", count), Namespace: ns},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicaCount,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"foo": "bar"},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"foo": "bar"}},
				Spec:       corev1.PodSpec{Containers: []corev1.Container{{Name: "nginx", Image: "nginx"}}},
			},
		},
	}
}

func newPod(count int, ns string) *corev1.Pod {
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: fmt.Sprintf("pod-%v", count), Namespace: ns},
		Spec:       corev1.PodSpec{Containers: []corev1.Container{{Name: "nginx", Image: "nginx"}}},
	}
}

func newNode(count int, ns string) *corev1.Node {
	return &corev1.Node{
		ObjectMeta: metav1.ObjectMeta{Name: fmt.Sprintf("node-name-%v", count)},
		Spec:       corev1.NodeSpec{},
	}
}

func TestPodNode(t *testing.T) {
	env := testenv.NewEnv(os.Getenv("K8S_ASSETS"), nil)
	err := env.Start()
	ut.Assert(t, err == nil, "testenv cluster start failed:%v", err)
	defer func() {
		env.Stop()
	}()

	clientset, err := kubernetes.NewForConfig(env.Config)
	ut.Assert(t, err == nil, "create clientset failed:%v", err)
	c, err := New(env.Config, Options{})
	ut.Assert(t, err == nil, "create client failed:%v", err)

	ns := "default"
	dep := newDeploy(0, ns)
	err = c.Create(context.TODO(), dep)
	ut.Assert(t, err == nil, "create deploy failed:%v", err)
	actual, err := clientset.AppsV1().Deployments(ns).Get(dep.Name, metav1.GetOptions{})
	ut.Assert(t, err == nil, "get deploy failed:%v", err)
	ut.Equal(t, dep, actual)

	var expectedDep appsv1.Deployment
	err = c.Get(context.TODO(), types.NamespacedName{ns, dep.Name}, &expectedDep)
	ut.Assert(t, err == nil, "get deploy failed:%v", err)
	ut.Equal(t, actual, &expectedDep)

	err = c.Delete(context.TODO(), dep)
	ut.Assert(t, err == nil, "delete deploy failed:%v", err)
	_, err = clientset.AppsV1().Deployments(ns).Get(dep.Name, metav1.GetOptions{})
	ut.Assert(t, err != nil, "get deleted deploy should fail")

	podCount := 10
	var pods []*corev1.Pod
	for i := 0; i < podCount; i++ {
		pod := newPod(i, ns)
		err = c.Create(context.TODO(), pod)
		ut.Assert(t, err == nil, "create pod failed:%v", err)
		pods = append(pods, pod)
	}
	podList := &corev1.PodList{}
	err = c.List(context.TODO(), nil, podList)
	ut.Assert(t, err == nil, "list pod failed:%v", err)
	ut.Equal(t, len(podList.Items), podCount)
	for _, pod := range pods {
		err := clientset.CoreV1().Pods(ns).Delete(pod.Name, &metav1.DeleteOptions{})
		ut.Assert(t, err == nil, "delete pod failed:%v", err)
	}
	podList = &corev1.PodList{}
	err = c.List(context.TODO(), nil, podList)
	ut.Equal(t, len(podList.Items), 0)
}
