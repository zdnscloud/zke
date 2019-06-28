package handler

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	ut "github.com/zdnscloud/cement/unittest"
	"github.com/zdnscloud/gorest/parse"
	"github.com/zdnscloud/gorest/types"
)

var (
	schemas = types.NewSchemas().MustImportAndCustomize(&version, Foo{}, &Handler{}, func(schema *types.Schema, handler types.Handler) {
		schema.CollectionMethods = []string{"POST", "GET"}
		schema.ResourceMethods = []string{"GET", "POST", "DELETE", "PUT"}
		schema.Handler = handler
		schema.ResourceActions = append(schema.ResourceActions, types.Action{
			Name:  "encode",
			Input: TestInput{},
		}, types.Action{
			Name:  "decode",
			Input: TestInput{},
		})
	})
)

type TestInput struct {
	Data string `json:"data"`
}

type Foo struct {
	types.Resource
	Name string `json:"name"singlecloud:"required=true"`
	Role string `json:"role"singlecloud:"required=true"`
}

type testServer struct {
	ctx *types.Context
}

func (t *testServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var err *types.APIError
	if t.ctx.Action != nil {
		err = ActionHandler(t.ctx)
	} else {
		switch req.Method {
		case http.MethodPost:
			err = CreateHandler(t.ctx)
		case http.MethodPut:
			err = UpdateHandler(t.ctx)
		case http.MethodDelete:
			err = DeleteHandler(t.ctx)
		case http.MethodGet:
			err = ListHandler(t.ctx)
		default:
			panic("unspport method " + req.Method)
		}
	}

	if err != nil {
		WriteResponse(t.ctx, err.Status, err)
	}
}

func TestCreateHandler(t *testing.T) {
	yamlContent := "testContent"
	expectBody := "{\"id\":\"12138\",\"type\":\"foo\",\"links\":{\"collection\":\"http://127.0.0.1:1234/apis/testing/v1/foos\",\"remove\":\"http://127.0.0.1:1234/apis/testing/v1/foos/12138\",\"self\":\"http://127.0.0.1:1234/apis/testing/v1/foos/12138\",\"update\":\"http://127.0.0.1:1234/apis/testing/v1/foos/12138\"},\"creationTimestamp\":null,\"name\":\"bar\",\"role\":\"master\"}"
	req, _ := http.NewRequest("POST", "/apis/testing/v1/foos", bytes.NewBufferString(fmt.Sprintf("{\"name\":\"bar\", \"yaml_\":\"%s\",\"role\":\"master\"}", yamlContent)))
	req.Host = "127.0.0.1:1234"
	w := httptest.NewRecorder()
	ctx, _ := parse.Parse(w, req, schemas)
	server := &testServer{}
	server.ctx = ctx
	server.ServeHTTP(w, req)
	ut.Equal(t, w.Code, 201)
	ut.Equal(t, w.Body.String(), expectBody)
}

func TestDeleteHandler(t *testing.T) {
	req, _ := http.NewRequest("DELETE", "/apis/testing/v1/foos/12138", nil)
	req.Host = "127.0.0.1:1234"
	w := httptest.NewRecorder()
	ctx, _ := parse.Parse(w, req, schemas)
	server := &testServer{}
	server.ctx = ctx
	server.ServeHTTP(w, req)
	ut.Equal(t, w.Code, 204)
}

func TestUpdateHandler(t *testing.T) {
	yamlContent := "testContent"
	expectBody := "{\"id\":\"12138\",\"type\":\"foo\",\"links\":{\"collection\":\"http://127.0.0.1:1234/apis/testing/v1/foos\",\"remove\":\"http://127.0.0.1:1234/apis/testing/v1/foos/12138\",\"self\":\"http://127.0.0.1:1234/apis/testing/v1/foos/12138\",\"update\":\"http://127.0.0.1:1234/apis/testing/v1/foos/12138\"},\"creationTimestamp\":null,\"name\":\"bar\",\"role\":\"worker\"}"
	req, _ := http.NewRequest("PUT", "/apis/testing/v1/foos/12138", bytes.NewBufferString(fmt.Sprintf("{\"name\":\"bar\", \"yaml_\":\"%s\",\"role\": \"worker\"}", yamlContent)))
	req.Host = "127.0.0.1:1234"
	w := httptest.NewRecorder()
	ctx, _ := parse.Parse(w, req, schemas)
	server := &testServer{}
	server.ctx = ctx
	server.ServeHTTP(w, req)
	ut.Equal(t, w.Code, 200)
	ut.Equal(t, w.Body.String(), expectBody)
}

func TestListHandler(t *testing.T) {
	expectCollection := "{\"type\":\"collection\",\"resourceType\":\"foo\",\"links\":{\"self\":\"http://127.0.0.1:1234/apis/testing/v1/foos\"},\"data\":[]}"
	req, _ := http.NewRequest("GET", "/apis/testing/v1/foos", nil)
	req.Host = "127.0.0.1:1234"
	w := httptest.NewRecorder()
	ctx, _ := parse.Parse(w, req, schemas)
	server := &testServer{}
	server.ctx = ctx
	server.ServeHTTP(w, req)
	ut.Equal(t, w.Code, 200)
	ut.Equal(t, w.Body.String(), expectCollection)
}

func TestGetOne(t *testing.T) {
	expectResult := "{\"id\":\"12138\",\"type\":\"foo\",\"links\":{\"collection\":\"http://127.0.0.1:1234/apis/testing/v1/foos\",\"remove\":\"http://127.0.0.1:1234/apis/testing/v1/foos/12138\",\"self\":\"http://127.0.0.1:1234/apis/testing/v1/foos/12138\",\"update\":\"http://127.0.0.1:1234/apis/testing/v1/foos/12138\"},\"creationTimestamp\":null,\"name\":\"bar\",\"role\":\"worker\"}"
	req, _ := http.NewRequest("GET", "/apis/testing/v1/foos/12138", nil)
	req.Host = "127.0.0.1:1234"
	w := httptest.NewRecorder()
	ctx, _ := parse.Parse(w, req, schemas)
	server := &testServer{}
	server.ctx = ctx
	server.ServeHTTP(w, req)
	ut.Equal(t, w.Code, 200)
	ut.Equal(t, w.Body.String(), expectResult)
}

func TestGetNonExists(t *testing.T) {
	expectResult := "{\"code\":\"NotFound\",\"status\":404,\"type\":\"error\",\"message\":\"foo resource with id 23456 doesn't exist\"}"
	req, _ := http.NewRequest("GET", "/apis/testing/v1/foos/23456", nil)
	req.Host = "127.0.0.1:1234"
	w := httptest.NewRecorder()
	ctx, _ := parse.Parse(w, req, schemas)
	server := &testServer{}
	server.ctx = ctx
	server.ServeHTTP(w, req)
	ut.Equal(t, w.Code, 404)
	ut.Equal(t, w.Body.String(), expectResult)
}

func TestActionHandler(t *testing.T) {
	encodeReq, _ := http.NewRequest("POST", "/apis/testing/v1/foos/123?action=encode", bytes.NewBufferString("{\"data\":\"testdata\"}"))
	encodeReq.Host = "127.0.0.1:1234"
	w := httptest.NewRecorder()
	ctx, _ := parse.Parse(w, encodeReq, schemas)
	server := &testServer{}
	server.ctx = ctx
	server.ServeHTTP(w, encodeReq)
	ut.Equal(t, w.Code, 200)
	base64str := base64.StdEncoding.EncodeToString([]byte("testdata"))
	expectResult := "\"" + base64str + "\""
	ut.Equal(t, w.Body.String(), expectResult)

	decodeReq, _ := http.NewRequest("POST", "/apis/testing/v1/foos/123?action=decode", bytes.NewBufferString(fmt.Sprintf("{\"data\":\"%s\"}", base64str)))
	decodeReq.Host = "127.0.0.1:1234"
	w = httptest.NewRecorder()
	ctx, _ = parse.Parse(w, decodeReq, schemas)
	server = &testServer{}
	server.ctx = ctx
	server.ServeHTTP(w, decodeReq)
	ut.Equal(t, w.Code, 200)
	expectResult = "\"testdata\""
	ut.Equal(t, w.Body.String(), expectResult)
}

type Handler struct{}

func (h *Handler) Create(ctx *types.Context, content []byte) (interface{}, *types.APIError) {
	ctx.Object.SetID("12138")
	return ctx.Object, nil
}

func (h *Handler) Delete(ctx *types.Context) *types.APIError {
	return nil
}

func (h *Handler) Update(ctx *types.Context) (interface{}, *types.APIError) {
	ctx.Object.SetID("12138")
	return ctx.Object, nil
}

func (h *Handler) List(ctx *types.Context) interface{} {
	return []types.Object{}
}

func (h *Handler) Get(ctx *types.Context) interface{} {
	if ctx.Object.GetID() == "12138" {
		foo := ctx.Object.(*Foo)
		foo.Name = "bar"
		foo.Role = "worker"
		return foo
	}
	return nil
}

func (h *Handler) Action(ctx *types.Context) (interface{}, *types.APIError) {
	input, ok := ctx.Action.Input.(*TestInput)
	if ok == false {
		return nil, types.NewAPIError(types.InvalidFormat, "action input type invalid")
	}

	var err *types.APIError
	switch ctx.Action.Name {
	case "encode":
		return base64.StdEncoding.EncodeToString([]byte(input.Data)), nil
	case "decode":
		if data, e := base64.StdEncoding.DecodeString(input.Data); e != nil {
			err = types.NewAPIError(types.InvalidFormat, "decode failed: "+e.Error())
		} else {
			return string(data), nil
		}
	default:
		err = types.NewAPIError(types.NotFound, "not found action "+ctx.Action.Name)
	}

	return nil, err
}
