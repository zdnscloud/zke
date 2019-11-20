package linkerd

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	yamlDecoder "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/helm/pkg/chartutil"
	"sigs.k8s.io/yaml"

	"github.com/zdnscloud/zke/zcloud/linkerd/pkg/charts"
	"github.com/zdnscloud/zke/zcloud/linkerd/pkg/pbconfig"
)

const (
	configStage          = "config"
	controlPlaneStage    = "control-plane"
	helmDefaultChartName = "linkerd2"
	helmDefaultChartDir  = "linkerd2"
)

var (
	templatesConfigStage = []string{
		"templates/namespace.yaml",
		"templates/identity-rbac.yaml",
		"templates/controller-rbac.yaml",
		"templates/destination-rbac.yaml",
		"templates/heartbeat-rbac.yaml",
		"templates/web-rbac.yaml",
		"templates/serviceprofile-crd.yaml",
		"templates/trafficsplit-crd.yaml",
		"templates/prometheus-rbac.yaml",
		"templates/grafana-rbac.yaml",
		"templates/proxy-injector-rbac.yaml",
		"templates/sp-validator-rbac.yaml",
		"templates/tap-rbac.yaml",
		"templates/psp.yaml",
	}

	templatesControlPlaneStage = []string{
		"templates/_validate.tpl",
		"templates/_affinity.tpl",
		"templates/_config.tpl",
		"templates/_helpers.tpl",
		"templates/_nodeselector.tpl",
		"templates/config.yaml",
		"templates/identity.yaml",
		"templates/controller.yaml",
		"templates/destination.yaml",
		"templates/heartbeat.yaml",
		"templates/web.yaml",
		"templates/prometheus.yaml",
		"templates/grafana.yaml",
		"templates/proxy-injector.yaml",
		"templates/sp-validator.yaml",
		"templates/tap.yaml",
	}
)

func GetDeployYaml(clusterDomain string) (string, error) {
	var buf bytes.Buffer
	if err := installRunE(newInstallOptionsWithDefaults(clusterDomain), &buf); err != nil {
		return err
	}

	return buf.String(), nil
}

func installRunE(options *InstallOptions, out io.Writer) error {
	values, configs, err := options.validateAndBuild()
	if err != nil {
		return err
	}

	return render(out, values, configs)
}

func render(w io.Writer, values *charts.Values, configs *pbconfig.All) error {
	rawValues, err := yaml.Marshal(values)
	if err != nil {
		return err
	}

	files := []*chartutil.BufferedFile{
		{Name: chartutil.ChartfileName},
	}

	if values.Stage == "" || values.Stage == configStage {
		for _, template := range templatesConfigStage {
			files = append(files, &chartutil.BufferedFile{
				Name: template,
			})
		}
	}

	if values.Stage == "" || values.Stage == controlPlaneStage {
		for _, template := range templatesControlPlaneStage {
			files = append(files, &chartutil.BufferedFile{
				Name: template,
			})
		}
	}

	chart := &charts.Chart{
		Name:      helmDefaultChartName,
		Dir:       helmDefaultChartDir,
		Namespace: defaultNamespace,
		RawValues: rawValues,
		Files:     files,
	}
	buf, err := chart.Render()
	if err != nil {
		return err
	}

	return processYAML(&buf, w, resourceTransformerInject{
		injectProxy: true,
		configs:     configs,
	})
}

func processYAML(in io.Reader, out io.Writer, rt resourceTransformerInject) error {
	reader := yamlDecoder.NewYAMLReader(bufio.NewReaderSize(in, 4096))
	for {
		bytes, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		var result []byte
		isList, err := kindIsList(bytes)
		if err != nil {
			return err
		}
		if isList {
			result, err = processList(bytes, rt)
		} else {
			result, err = rt.transform(bytes)
		}
		if err != nil {
			return err
		}

		out.Write(result)
		out.Write([]byte("---\n"))
	}

	return nil
}

func kindIsList(bytes []byte) (bool, error) {
	var meta metav1.TypeMeta
	if err := yaml.Unmarshal(bytes, &meta); err != nil {
		return false, err
	}
	return meta.Kind == "List", nil
}

func processList(bytes []byte, rt resourceTransformerInject) ([]byte, error) {
	var sourceList corev1.List
	if err := yaml.Unmarshal(bytes, &sourceList); err != nil {
		return nil, err
	}

	items := []runtime.RawExtension{}
	for _, item := range sourceList.Items {
		result, err := rt.transform(item.Raw)
		if err != nil {
			return nil, err
		}

		injected, err := yaml.YAMLToJSON(result)
		if err != nil {
			return nil, err
		}

		items = append(items, runtime.RawExtension{Raw: injected})
	}

	sourceList.Items = items
	result, err := yaml.Marshal(sourceList)
	if err != nil {
		return nil, err
	}
	return result, nil
}
