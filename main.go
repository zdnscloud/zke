package main

import (
	"fmt"
	"os"

	"github.com/zdnscloud/zke/cmd"
	"github.com/zdnscloud/zke/pkg/log"

	"github.com/urfave/cli"
	cementlog "github.com/zdnscloud/cement/log"
)

var VERSION = "v1.0.0"
var BUILD string

func main() {
	if err := mainErr(); err != nil {
		fmt.Printf("zke err %s", err.Error())
	}
}

func mainErr() error {
	app := cli.NewApp()
	app.Name = "zke"
	app.Version = VERSION
	app.Usage = "ZDNS Kubernetes Engine, an extremely simple, lightning fast Kubernetes installer that works everywhere"
	app.Before = func(ctx *cli.Context) error {
		if ctx.GlobalBool("debug") {
			log.LogLevel = cementlog.Debug
		}
		// log.Debugf("ZKE version %s build at %s", app.Version, BUILD)
		return nil
	}
	app.Author = "Zcloud"
	app.Email = "zcloud@zdns.cn"
	app.Commands = []cli.Command{
		cmd.UpCommand(),
		cmd.RemoveCommand(),
		cmd.ConfigCommand(),
		cmd.LoadImageCommand(),
	}

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug,d",
			Usage: "Debug logging",
		},
	}

	return app.Run(os.Args)
}
