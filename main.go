package main

import (
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
		log.Fatal(err)
	}
}

func mainErr() error {
	app := cli.NewApp()
	app.Name = "zke"
	app.Version = VERSION
	app.Usage = "ZDNS Kubernetes Engine, an extremely simple, lightning fast Kubernetes installer that works everywhere"
	app.Before = func(ctx *cli.Context) error {
		if ctx.GlobalBool("debug") {
			log.DefaultLogLevel = cementlog.Debug
		}
		log.InitConsoleLog()
		log.Debugf("ZKE version %s build at %s", app.Version, BUILD)
		return nil
	}
	app.Author = "Zcloud"
	app.Email = "zcloud@zdns.cn"
	app.Commands = []cli.Command{
		cmd.UpCommand(),
		cmd.RemoveCommand(),
		cmd.ConfigCommand(),
	}

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug,d",
			Usage: "Debug logging",
		},
	}

	return app.Run(os.Args)
}
