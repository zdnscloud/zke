package handler

import (
	"bytes"
	"net/http"
	"testing"

	ut "github.com/zdnscloud/cement/unittest"
	"github.com/zdnscloud/gorest/types"
)

type Deployment struct {
	types.Resource
	Name       string       `json:"name"singlecloud:"required=true"`
	Replicas   int          `json:"replicas"singlecloud:"required=true"`
	Containers []*Container `json:"containers"singlecloud:"required=true"`
}

type Container struct {
	Name         string           `json:"name"singlecloud:"required=true"`
	Image        string           `json:"image"singlecloud:"required=true"`
	Command      []string         `json:"command,omitempty"`
	Args         []string         `json:"args,omitempty"`
	ConfigName   string           `json:"configName,omitempty"`
	MountPath    string           `json:"mountPath,omitempty"`
	ExposedPorts []DeploymentPort `json:"exposedPorts,omitempty"`
}

type DeploymentPort struct {
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
}

func TestValidate(t *testing.T) {
	schemas := types.NewSchemas()
	schemas.MustImport(&version, Deployment{})
	schema := schemas.Schema(&version, types.GetResourceType(Deployment{}))
	reqbody := bytes.NewBufferString("{\"name\":\"dm1\", \"replicas\": 1, \"containers\": [{\"name\": \"c1\", \"image\": \"testimage\", \"comamnd\": [\"ifconfig\", \"|\", \"grep\"], \"exposedPorts\": [{\"port\": 8080, \"protocol\": \"tcp\"}]}]}")
	req, _ := http.NewRequest("POST", "/apis/testing/v1/deployments", reqbody)
	req.Host = "127.0.0.1:1234"
	ctx := &types.Context{
		Request: req,
		Schemas: schemas,
		Object: &types.Resource{
			Schema: schema,
		},
	}

	var noerr *types.APIError
	_, err := parseCreateBody(ctx)
	ut.Equal(t, err, noerr)
	err = CheckObjectFields(ctx)
	ut.Equal(t, err, noerr)
}
