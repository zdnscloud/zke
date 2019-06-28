package cmd

import (
	"fmt"
	"time"

	"github.com/zdnscloud/zke/types"

	"github.com/gin-gonic/gin"
	"github.com/urfave/cli"
	"github.com/zdnscloud/cement/log"
	"github.com/zdnscloud/gorest/adaptor"
	"github.com/zdnscloud/gorest/api"
	resttypes "github.com/zdnscloud/gorest/types"
)

var (
	Version = resttypes.APIVersion{
		Version: "v1",
		Group:   "zkemanager.zcloud.cn",
	}
)

func ServeCommand() cli.Command {
	return cli.Command{
		Name:   "serve",
		Usage:  "start a zke restserver at port 8080",
		Action: clusterServeFromCli,
	}
}

type ZKEManager struct {
	api.DefaultHandler

	Configs map[string]*types.ZKEConfig
}

func NewZKEManager() ZKEManager {
	return ZKEManager{
		Configs: make(map[string]*types.ZKEConfig),
	}
}

func (m *ZKEManager) RegisterSchemas(version *resttypes.APIVersion, schemas *resttypes.Schemas) {
	schemas.MustImportAndCustomize(version, types.ZKEConfig{}, m, func(schema *resttypes.Schema, handler resttypes.Handler) {
		schema.Handler = handler
		schema.CollectionMethods = []string{"GET", "POST", "DELETE"}
		schema.ResourceMethods = []string{"GET", "POST", "DELETE"}
	})
}

func (m *ZKEManager) List(ctx *resttypes.Context) interface{} {
	var configs []*types.ZKEConfig
	for _, c := range m.Configs {
		configs = append(configs, c)
	}
	return configs
}

func (m *ZKEManager) Get(ctx *resttypes.Context) interface{} {
	id := ctx.Object.GetID()
	c, ok := m.Configs[id]
	if ok {
		return c
	}
	return nil
}

func (m *ZKEManager) Create(ctx *resttypes.Context, yamlConf []byte) (interface{}, *resttypes.APIError) {
	c := ctx.Object.(*types.ZKEConfig)
	c.SetID(c.ClusterName)
	c.SetCreationTimestamp(time.Now())
	_, ok := m.Configs[c.ClusterName]
	if ok {
		return c, resttypes.NewAPIError(resttypes.DuplicateResource, "dulicate clusters")
	}
	m.Configs[c.ClusterName] = c

	excuateErr := ClusterUpFromRestClient(c)
	if excuateErr != nil {
		err := &resttypes.APIError{
			Message: fmt.Sprintf("cluster up err: %s", excuateErr),
		}
		return c, err
	}

	return c, nil
}

func (m *ZKEManager) Delete(ctx *resttypes.Context) *resttypes.APIError {
	id := ctx.Object.GetID()
	config, ok := m.Configs[id]
	if !ok {
		return resttypes.NewAPIError(resttypes.NotFound, "cluster doesn't exist")
	}
	err := ClusterRemoveFromRestClient(config)
	if err != nil {
		return &resttypes.APIError{
			Message: fmt.Sprintf("cluster remove err: %s", err),
		}
	}
	delete(m.Configs, id)
	return nil
}

func clusterServeFromCli(ctx *cli.Context) {
	log.InitLogger("debug")

	zkeMgr := NewZKEManager()

	schemas := resttypes.NewSchemas()
	zkeMgr.RegisterSchemas(&Version, schemas)

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	server := api.NewAPIServer()
	if err := server.AddSchemas(schemas); err != nil {
		log.Fatalf("add schemas failed:%s", err.Error())
	}
	server.Use(api.RestHandler)
	adaptor.RegisterHandler(router, server, server.Schemas.UrlMethods())

	addr := "0.0.0.0:8080"
	router.Run(addr)
}
