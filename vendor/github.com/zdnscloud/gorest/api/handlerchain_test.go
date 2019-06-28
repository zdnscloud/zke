package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	ut "github.com/zdnscloud/cement/unittest"
	"github.com/zdnscloud/gorest/types"
)

var (
	schemas = types.NewSchemas().MustImportAndCustomize(&version, Foo{}, nil, func(schema *types.Schema, handler types.Handler) {
		schema.CollectionMethods = []string{"POST", "GET"}
	})

	version = types.APIVersion{
		Group:   "testing",
		Version: "v1",
	}
)

type Foo struct {
	types.Resource
}

var gnum int

var dumbHandler1 = func(ctx *types.Context) *types.APIError {
	ctx.Set("key", &gnum)
	return nil
}

var dumbHandler2 = func(ctx *types.Context) *types.APIError {
	val_, _ := ctx.Get("key")
	*(val_.(*int)) = 100
	return nil
}

func TestContextPassChain(t *testing.T) {
	req, _ := http.NewRequest("GET", "/apis/testing/v1/foos", nil)
	req.Host = "127.0.0.1:1234"
	w := httptest.NewRecorder()
	s := NewAPIServer()
	s.AddSchemas(schemas)
	s.Use(dumbHandler1)
	s.Use(dumbHandler2)

	ut.Equal(t, gnum, 0)
	s.ServeHTTP(w, req)
	ut.Equal(t, gnum, 100)

	s.Use(RestHandler)
	s.ServeHTTP(w, req)
}
