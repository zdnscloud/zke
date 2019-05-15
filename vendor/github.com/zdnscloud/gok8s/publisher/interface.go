package publisher

import (
	"k8s.io/client-go/tools/record"
)

type EventPublisher interface {
	GetEventRecorderFor(name string) record.EventRecorder
}
