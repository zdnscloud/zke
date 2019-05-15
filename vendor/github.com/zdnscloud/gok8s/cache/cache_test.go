package cache

import (
	"context"
	"os"
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	kcache "k8s.io/client-go/tools/cache"

	ut "github.com/zdnscloud/cement/unittest"
	"github.com/zdnscloud/gok8s/client"
	"github.com/zdnscloud/gok8s/testenv"
)

func newPod(name, ns string, labels map[string]string, restartPolicy corev1.RestartPolicy) *corev1.Pod {
	three := int64(3)
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			Containers:            []corev1.Container{{Name: "nginx", Image: "nginx"}},
			RestartPolicy:         restartPolicy,
			ActiveDeadlineSeconds: &three,
		},
	}
}

func TestPodCache(t *testing.T) {
	env := testenv.NewEnv(os.Getenv("K8S_ASSETS"), nil)
	err := env.Start()
	ut.Assert(t, err == nil, "testenv cluster start failed:%v", err)
	defer func() {
		env.Stop()
	}()

	cli, err := client.New(env.Config, client.Options{})
	ut.Assert(t, err == nil, "create client failed:%v", err)

	testNamespaceOne := "test-namespace-1"
	testNamespaceTwo := "test-namespace-2"
	err = cli.Create(context.TODO(), newPod("test-pod-1", testNamespaceOne, map[string]string{"test-label": "test-pod-1"}, corev1.RestartPolicyNever))
	ut.Assert(t, err == nil, "create pod failed:%v", err)
	err = cli.Create(context.TODO(), newPod("test-pod-2", testNamespaceTwo, map[string]string{"test-label": "test-pod-2"}, corev1.RestartPolicyAlways))
	ut.Assert(t, err == nil, "create pod failed:%v", err)
	err = cli.Create(context.TODO(), newPod("test-pod-3", testNamespaceTwo, map[string]string{"test-label": "test-pod-3"}, corev1.RestartPolicyOnFailure))
	ut.Assert(t, err == nil, "create pod failed:%v", err)

	stop := make(chan struct{})
	defer close(stop)
	c, err := New(env.Config, Options{})
	ut.Assert(t, err == nil, "create cache failed:%v", err)
	go c.Start(stop)
	ut.Assert(t, c.WaitForCacheSync(stop), "wait for sync should ok")

	//read test
	svcs := &corev1.ServiceList{}
	err = c.List(context.TODO(), nil, svcs)
	ut.Assert(t, err == nil, "list services failed:%v", err)
	hasKubeService := false
	for _, svc := range svcs.Items {
		if svc.Namespace == "default" && svc.Name == "kubernetes" {
			hasKubeService = true
			break
		}
	}
	ut.Assert(t, hasKubeService, "no kubeservice found")

	svc := &corev1.Service{}
	svcKey := client.ObjectKey{Namespace: "default", Name: "kubernetes"}
	err = c.Get(context.TODO(), svcKey, svc)
	ut.Assert(t, err == nil, "list services failed:%v", err)
	ut.Equal(t, svc.Name, "kubernetes")
	ut.Equal(t, svc.Namespace, "default")

	pods := corev1.PodList{}
	listOpt := &client.ListOptions{}
	listOpt.InNamespace(testNamespaceTwo)
	listOpt.MatchingLabels(map[string]string{"test-label": "test-pod-2"})
	err = c.List(context.TODO(), listOpt, &pods)
	ut.Assert(t, err == nil, "list pod failed:%v", err)
	ut.Equal(t, len(pods.Items), 1)
	ut.Equal(t, pods.Items[0].Labels["test-label"], "test-pod-2")

	pod := &corev1.Pod{}
	informer, err := c.GetInformer(pod)
	ut.Assert(t, err == nil, "get informer for pod failed:%v", err)
	ut.Assert(t, informer.HasSynced(), "pod informer should synced")
	out := make(chan interface{}, 10)
	addFunc := func(obj interface{}) {
		out <- obj
	}
	informer.AddEventHandler(kcache.ResourceEventHandlerFuncs{AddFunc: addFunc})
	//sync the state, which will return 3 pod
	for i := 0; i < 3; i++ {
		<-out
	}
	err = cli.Create(context.TODO(), newPod("test-pod-4", testNamespaceOne, map[string]string{"test-label": "test-pod-4"}, corev1.RestartPolicyOnFailure))
	newCreatePod := <-out
	ut.Equal(t, newCreatePod.(*corev1.Pod).Labels["test-label"], "test-pod-4")

	pods = corev1.PodList{}
	listOpt = &client.ListOptions{}
	c.List(context.TODO(), listOpt, &pods)
	ut.Equal(t, len(pods.Items), 4)

	err = cli.Delete(context.TODO(), newPod("test-pod-1", testNamespaceOne, nil, corev1.RestartPolicyNever))
	ut.Assert(t, err == nil, "delete pod failed:%v", err)
	err = cli.Delete(context.TODO(), newPod("test-pod-2", testNamespaceTwo, nil, corev1.RestartPolicyAlways))
	ut.Assert(t, err == nil, "delete pod failed:%v", err)
	err = cli.Delete(context.TODO(), newPod("test-pod-3", testNamespaceTwo, nil, corev1.RestartPolicyOnFailure))
	ut.Assert(t, err == nil, "delete pod failed:%v", err)
	err = cli.Delete(context.TODO(), newPod("test-pod-4", testNamespaceOne, nil, corev1.RestartPolicyOnFailure))
	ut.Assert(t, err == nil, "delete pod failed:%v", err)

	<-time.After(time.Second)
	pods = corev1.PodList{}
	c.List(context.TODO(), listOpt, &pods)
	ut.Equal(t, len(pods.Items), 0)
}

func TestPodCacheIndex(t *testing.T) {
	env := testenv.NewEnv(os.Getenv("K8S_ASSETS"), nil)
	err := env.Start()
	ut.Assert(t, err == nil, "testenv cluster start failed:%v", err)
	defer func() {
		env.Stop()
	}()

	cli, err := client.New(env.Config, client.Options{})
	ut.Assert(t, err == nil, "create client failed:%v", err)

	testNamespaceOne := "test-namespace-1"
	testNamespaceTwo := "test-namespace-2"
	err = cli.Create(context.TODO(), newPod("test-pod-1", testNamespaceOne, map[string]string{"test-label": "test-pod-1"}, corev1.RestartPolicyNever))
	ut.Assert(t, err == nil, "create pod failed:%v", err)
	err = cli.Create(context.TODO(), newPod("test-pod-2", testNamespaceTwo, map[string]string{"test-label": "test-pod-2"}, corev1.RestartPolicyAlways))
	ut.Assert(t, err == nil, "create pod failed:%v", err)
	err = cli.Create(context.TODO(), newPod("test-pod-3", testNamespaceTwo, map[string]string{"test-label": "test-pod-3"}, corev1.RestartPolicyOnFailure))
	ut.Assert(t, err == nil, "create pod failed:%v", err)

	stop := make(chan struct{})
	defer close(stop)
	c, err := New(env.Config, Options{})
	ut.Assert(t, err == nil, "create cache failed:%v", err)
	indexFunc := func(obj runtime.Object) []string {
		return []string{string(obj.(*corev1.Pod).Spec.RestartPolicy)}
	}
	c.IndexField(&corev1.Pod{}, "spec.restartPolicy", indexFunc)
	ut.Assert(t, err == nil, "index pod failed:%v", err)

	go c.Start(stop)
	ut.Assert(t, c.WaitForCacheSync(stop), "wait for sync should ok")

	pods := corev1.PodList{}
	listOpt := &client.ListOptions{}
	listOpt.MatchingField("spec.restartPolicy", "OnFailure")
	c.List(context.TODO(), listOpt, &pods)
	ut.Equal(t, len(pods.Items), 1)
	ut.Equal(t, pods.Items[0].Name, "test-pod-3")
}
