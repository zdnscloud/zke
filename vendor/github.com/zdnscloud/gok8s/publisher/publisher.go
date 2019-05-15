package publisher

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"
)

type publisher struct {
	scheme           *runtime.Scheme
	eventBroadcaster record.EventBroadcaster
}

func New(config *rest.Config, scheme *runtime.Scheme) (EventPublisher, error) {
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to init clientSet: %v", err)
	}

	p := &publisher{scheme: scheme}
	p.eventBroadcaster = record.NewBroadcaster()
	p.eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: clientSet.CoreV1().Events("")})
	return p, nil
}

func (p *publisher) GetEventRecorderFor(name string) record.EventRecorder {
	return p.eventBroadcaster.NewRecorder(p.scheme, corev1.EventSource{Component: name})
}
