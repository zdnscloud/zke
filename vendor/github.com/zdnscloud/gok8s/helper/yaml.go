package helper

import (
	"bufio"
	"context"
	//"fmt"
	"io"
	"strings"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/zdnscloud/gok8s/client"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes/scheme"
)

func CreateResourceFromYaml(cli client.Client, yaml string) error {
	return mapOnYamlDocument(yaml, cli.Create)
}

func DeleteResourceFromYaml(cli client.Client, yaml string) error {
	return mapOnYamlDocument(yaml, func(ctx context.Context, obj runtime.Object) error {
		return cli.Delete(ctx, obj, client.PropagationPolicy(metav1.DeletePropagationForeground))
	})
}

func UpdateResourceFromYaml(cli client.Client, yaml string) error {
	return mapOnYamlDocument(yaml, cli.Update)
}

func mapOnYamlDocument(data string, fn func(context.Context, runtime.Object) error) error {
	decode := scheme.Codecs.UniversalDeserializer().Decode
	reader := yaml.NewYAMLReader(bufio.NewReader(strings.NewReader(data)))
	for {
		doc, err := reader.Read()
    //fmt.Printf("---> doc:%s, err:%v\n", string(doc), err)
		if err != nil {
			if err == io.EOF {
				return nil
			} else {
				return err
			}
		}

		obj, _, err := decode(doc, nil, nil)
		if err != nil {
			return err
		}
		if err := fn(context.TODO(), obj); err != nil {
			if apierrors.IsAlreadyExists(err) == false &&
				apierrors.IsNotFound(err) == false {
				return err
			}
		}
	}
	return nil
}
