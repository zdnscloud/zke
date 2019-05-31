package helper

import (
	"context"
	"strings"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/zdnscloud/gok8s/client"
	"k8s.io/client-go/kubernetes/scheme"
)

const YamlDelimiter = "---\n"

func CreateResourceFromYaml(cli client.Client, yaml string) error {
	return mapOnYamlDocument(yaml, cli.Create)
}

func DeleteResourceFromYaml(cli client.Client, yaml string) error {
	return mapOnYamlDocument(yaml, func(ctx context.Context, obj runtime.Object) error {
		return cli.Delete(ctx, obj)
	})
}

func UpdateResourceFromYaml(cli client.Client, yaml string) error {
	return mapOnYamlDocument(yaml, cli.Update)
}

func mapOnYamlDocument(yaml string, fn func(context.Context, runtime.Object) error) error {
	decode := scheme.Codecs.UniversalDeserializer().Decode
	for _, doc := range strings.Split(yaml, YamlDelimiter) {
		doc = strings.TrimSpace(doc)
		if doc == "" {
			continue
		}

		obj, _, err := decode([]byte(doc), nil, nil)
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
