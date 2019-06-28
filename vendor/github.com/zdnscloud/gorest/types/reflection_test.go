package types

import (
	"testing"

	ut "github.com/zdnscloud/cement/unittest"
)

var version = APIVersion{
	Group:   "testing",
	Version: "v1",
}

type Cluster struct {
	Resource
	Name string
}

type Node struct {
	Resource
	Name string
}

type NameSpace struct {
	Resource
	Name string
}

type Deployment struct {
	Resource
	Name string
}

type DaemonSet struct {
	Resource
	Name string
}

type StatefulSet struct {
	Resource
	Name string
}

type Pod struct {
	Resource
	Name string
}

func TestReflection(t *testing.T) {
	schemas := NewSchemas()
	schemas.MustImportAndCustomize(&version, Cluster{}, nil, func(schema *Schema, handler Handler) {
		schema.CollectionMethods = []string{"GET", "POST"}
		schema.ResourceMethods = []string{"GET", "DELETE", "PUT"}
	})

	schemas.MustImportAndCustomize(&version, Node{}, nil, func(schema *Schema, handler Handler) {
		schema.Parents = []string{GetResourceType(Cluster{})}
		schema.CollectionMethods = []string{"GET", "POST"}
		schema.ResourceMethods = []string{"GET", "DELETE", "PUT"}
	})
	schemas.MustImportAndCustomize(&version, NameSpace{}, nil, func(schema *Schema, handler Handler) {
		schema.Parents = []string{GetResourceType(Cluster{})}
		schema.CollectionMethods = []string{"GET", "POST"}
		schema.ResourceMethods = []string{"GET", "DELETE", "PUT"}
	})
	schemas.MustImportAndCustomize(&version, Pod{}, nil, func(schema *Schema, handler Handler) {
		schema.Parents = []string{GetResourceType(Deployment{}), GetResourceType(DaemonSet{}), GetResourceType(StatefulSet{})}
		schema.CollectionMethods = []string{"GET", "POST"}
		schema.ResourceMethods = []string{"GET", "DELETE", "PUT"}
	})

	schemas.MustImportAndCustomize(&version, Deployment{}, nil, func(schema *Schema, handler Handler) {
		schema.Parents = []string{GetResourceType(NameSpace{})}
		schema.CollectionMethods = []string{"GET", "POST"}
		schema.ResourceMethods = []string{"GET", "DELETE", "PUT"}
	})

	schemas.MustImportAndCustomize(&version, StatefulSet{}, nil, func(schema *Schema, handler Handler) {
		schema.Parents = []string{GetResourceType(NameSpace{})}
		schema.CollectionMethods = []string{"GET", "POST"}
		schema.ResourceMethods = []string{"GET", "DELETE", "PUT"}
	})

	schemas.MustImportAndCustomize(&version, DaemonSet{}, nil, func(schema *Schema, handler Handler) {
		schema.Parents = []string{GetResourceType(NameSpace{})}
		schema.CollectionMethods = []string{"GET", "POST"}
		schema.ResourceMethods = []string{"GET", "DELETE", "PUT"}
	})

	clusterChildren := schemas.GetChildren(GetResourceType(Cluster{}))
	ut.Equal(t, len(clusterChildren), 2)

	schema := schemas.Schema(&version, GetResourceType(Node{}))
	ut.Equal(t, schema.GetType(), GetResourceType(Node{}))
	ut.Equal(t, schema.PluralName, "nodes")
	ut.Equal(t, schema.Version.Group, "testing")
	ut.Equal(t, schema.Version.Version, "v1")
	ut.Equal(t, schema.Parents, []string{GetResourceType(Cluster{})})
	ut.Equal(t, schema.CollectionMethods, []string{"GET", "POST"})
	ut.Equal(t, schema.ResourceMethods, []string{"GET", "DELETE", "PUT"})
	ut.Equal(t, len(schema.ResourceFields), 3)

	expectUrl := []string{
		"/apis/testing/v1/clusters",
		"/apis/testing/v1/clusters/:cluster_id",
		"/apis/testing/v1/clusters/:cluster_id/nodes",
		"/apis/testing/v1/clusters/:cluster_id/nodes/:node_id",
		"/apis/testing/v1/clusters/:cluster_id/namespaces",
		"/apis/testing/v1/clusters/:cluster_id/namespaces/:namespace_id",
		"/apis/testing/v1/clusters/:cluster_id/namespaces/:namespace_id/deployments",
		"/apis/testing/v1/clusters/:cluster_id/namespaces/:namespace_id/deployments/:deployment_id",
		"/apis/testing/v1/clusters/:cluster_id/namespaces/:namespace_id/daemonsets",
		"/apis/testing/v1/clusters/:cluster_id/namespaces/:namespace_id/daemonsets/:daemonset_id",
		"/apis/testing/v1/clusters/:cluster_id/namespaces/:namespace_id/statefulsets",
		"/apis/testing/v1/clusters/:cluster_id/namespaces/:namespace_id/statefulsets/:statefulset_id",
		"/apis/testing/v1/clusters/:cluster_id/namespaces/:namespace_id/deployments/:deployment_id/pods",
		"/apis/testing/v1/clusters/:cluster_id/namespaces/:namespace_id/deployments/:deployment_id/pods/:pod_id",
		"/apis/testing/v1/clusters/:cluster_id/namespaces/:namespace_id/daemonsets/:daemonset_id/pods",
		"/apis/testing/v1/clusters/:cluster_id/namespaces/:namespace_id/daemonsets/:daemonset_id/pods/:pod_id",
		"/apis/testing/v1/clusters/:cluster_id/namespaces/:namespace_id/statefulsets/:statefulset_id/pods",
		"/apis/testing/v1/clusters/:cluster_id/namespaces/:namespace_id/statefulsets/:statefulset_id/pods/:pod_id",
	}
	urlMethods := schemas.UrlMethods()
	ut.Equal(t, len(urlMethods), len(expectUrl))
	for _, url := range expectUrl {
		ut.Equal(t, len(urlMethods[url]) != 0, true)
	}
}
