package linkerd

import (
	jsonpatch "github.com/evanphx/json-patch"
	"sigs.k8s.io/yaml"

	"github.com/zdnscloud/zke/zcloud/linkerd/pkg/inject"
	"github.com/zdnscloud/zke/zcloud/linkerd/pkg/k8s"
	"github.com/zdnscloud/zke/zcloud/linkerd/pkg/pbconfig"
)

type resourceTransformerInject struct {
	injectProxy         bool
	configs             *pbconfig.All
	overrideAnnotations map[string]string
	enableDebugSidecar  bool
}

func (rt resourceTransformerInject) transform(bytes []byte) ([]byte, error) {
	conf := inject.NewResourceConfig(rt.configs, inject.OriginCLI)

	if rt.enableDebugSidecar {
		conf.AppendPodAnnotation(k8s.ProxyEnableDebugAnnotation, "true")
	}

	report, err := conf.ParseMetaAndYAML(bytes)
	if err != nil {
		return nil, err
	}

	if b, _ := report.Injectable(); !b {
		return bytes, nil
	}

	if rt.injectProxy {
		conf.AppendPodAnnotation(k8s.CreatedByAnnotation, k8s.CreatedByAnnotationValue())
	} else {
		conf.AppendPodAnnotation(k8s.ProxyInjectAnnotation, k8s.ProxyInjectEnabled)
	}

	if len(rt.overrideAnnotations) > 0 {
		conf.AppendPodAnnotations(rt.overrideAnnotations)
	}

	patchJSON, err := conf.GetPatch(rt.injectProxy)
	if err != nil {
		return nil, err
	}
	if len(patchJSON) == 0 {
		return bytes, nil
	}
	patch, err := jsonpatch.DecodePatch(patchJSON)
	if err != nil {
		return nil, err
	}
	origJSON, err := yaml.YAMLToJSON(bytes)
	if err != nil {
		return nil, err
	}
	injectedJSON, err := patch.Apply(origJSON)
	if err != nil {
		return nil, err
	}
	injectedYAML, err := conf.JSONToYAML(injectedJSON)
	if err != nil {
		return nil, err
	}
	return injectedYAML, nil
}
