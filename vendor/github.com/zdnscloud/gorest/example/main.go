package main

import (
	"encoding/base64"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zdnscloud/cement/uuid"
	"github.com/zdnscloud/gorest/adaptor"
	"github.com/zdnscloud/gorest/api"
	"github.com/zdnscloud/gorest/types"
)

var (
	version = types.APIVersion{
		Group:   "zdns.cloud.example",
		Version: "example/v1",
	}
)

type Cluster struct {
	types.Resource `json:",inline"`
	Name           string `json:"name,omitempty"`
}

type Node struct {
	types.Resource `json:",inline"`
	Name           string `json:"name,omitempty"`
}

type Handler struct {
	objects map[string]types.Object
}

func newHandler() *Handler {
	return &Handler{
		objects: make(map[string]types.Object),
	}
}

func (h *Handler) Create(ctx *types.Context, content []byte) (interface{}, *types.APIError) {
	id, _ := uuid.Gen()
	switch ctx.Object.GetType() {
	case "cluster":
		cluster := ctx.Object.(*Cluster)
		for _, object := range h.objects {
			if object.GetType() == "cluster" && object.(*Cluster).Name == cluster.Name {
				return nil, types.NewAPIError(types.DuplicateResource, "cluster "+cluster.Name+" already exists")
			}
		}

		cluster.SetID(id)
		cluster.SetCreationTimestamp(time.Now())
		h.objects[id] = cluster
		return cluster, nil
	case "node":
		if parent := ctx.Object.GetParent(); parent != nil {
			if h.hasID(parent.GetID()) == false {
				return nil, types.NewAPIError(types.NotFound, "cluster "+parent.GetID()+" is non-exists")
			}
		}

		node := ctx.Object.(*Node)
		for _, object := range h.objects {
			if object.GetType() == "node" && object.(*Node).Name == node.Name {
				return nil, types.NewAPIError(types.DuplicateResource, "node "+node.Name+" already exists")
			}
		}

		node.SetID(id)
		node.SetCreationTimestamp(time.Now())
		h.objects[id] = node
		return node, nil
	default:
		return nil, types.NewAPIError(types.NotFound, "no found resource type "+ctx.Object.GetType())
	}
}

func (h *Handler) hasObject(obj types.Object) *types.APIError {
	if parent := obj.GetParent(); parent != nil {
		if h.hasID(parent.GetID()) == false {
			return types.NewAPIError(types.NotFound, "cluster "+parent.GetID()+" is non-exists")
		}
	}

	if h.hasID(obj.GetID()) == false {
		return types.NewAPIError(types.NotFound, "no found resource "+obj.GetType()+" with id "+obj.GetID())
	}

	return nil
}

func (h *Handler) hasID(id string) bool {
	_, ok := h.objects[id]
	return ok
}

func (h *Handler) hasChild(id string) bool {
	for _, obj := range h.objects {
		if parent := obj.GetParent(); parent != nil && parent.GetID() == id {
			return true
		}
	}

	return false
}

func (h *Handler) Delete(ctx *types.Context) *types.APIError {
	if err := h.hasObject(ctx.Object); err != nil {
		return err
	}

	if h.hasChild(ctx.Object.GetID()) {
		return types.NewAPIError(types.DeleteParent, "resource has child resource")
	}

	delete(h.objects, ctx.Object.GetID())
	return nil
}

func (h *Handler) Update(ctx *types.Context) (interface{}, *types.APIError) {
	if err := h.hasObject(ctx.Object); err != nil {
		return nil, err
	}

	h.objects[ctx.Object.GetID()] = ctx.Object
	return ctx.Object, nil
}

func (h *Handler) List(ctx *types.Context) interface{} {
	var result []types.Object
	for _, object := range h.objects {
		if object.GetType() == ctx.Object.GetType() {
			result = append(result, object)
		}
	}
	return result
}

func (h *Handler) Get(ctx *types.Context) interface{} {
	if parent := ctx.Object.GetParent(); parent != nil && h.hasID(parent.GetID()) == false {
		return nil
	}

	return h.objects[ctx.Object.GetID()]
}

func (h *Handler) Action(ctx *types.Context) (interface{}, *types.APIError) {
	err := h.hasObject(ctx.Object)
	if err != nil {
		return nil, err
	}

	input, ok := ctx.Action.Input.(*Input)
	if ok == false {
		return nil, types.NewAPIError(types.InvalidFormat, "action input type invalid")
	}

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

type Input struct {
	Data string `json:"data,omitempty"`
}

func main() {
	router := gin.Default()
	apiServer := getApiServer()
	adaptor.RegisterHandler(router, apiServer, apiServer.Schemas.UrlMethods())
	router.Run("0.0.0.0:1234")
}

func getApiServer() *api.Server {
	server := api.NewAPIServer()
	schemas := types.NewSchemas()
	handler := newHandler()
	schemas.MustImportAndCustomize(&version, Cluster{}, handler, func(schema *types.Schema, handler types.Handler) {
		schema.Handler = handler
		schema.CollectionMethods = []string{"GET", "POST"}
		schema.ResourceMethods = []string{"GET", "PUT", "DELETE", "POST"}
		schema.ResourceActions = append(schema.ResourceActions, types.Action{
			Name:  "encode",
			Input: Input{},
		}, types.Action{
			Name:  "decode",
			Input: Input{},
		})
	})

	schemas.MustImportAndCustomize(&version, Node{}, handler, func(schema *types.Schema, handler types.Handler) {
		schema.Parents = []string{types.GetResourceType(Cluster{})}
		schema.Handler = handler
		schema.CollectionMethods = []string{"GET", "POST"}
		schema.ResourceMethods = []string{"GET", "PUT", "DELETE", "POST"}
	})

	if err := server.AddSchemas(schemas); err != nil {
		panic(err.Error())
	}

	server.Use(api.RestHandler)
	return server
}
