package main

import (
	"os"

	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"github.com/zdnscloud/zke/cmd"
)

var VERSION = "v1.0.0"
var BUILD string

func main() {
	logrus.SetOutput(colorable.NewColorableStdout())

	if err := mainErr(); err != nil {
		logrus.Fatal(err)
	}
}

func mainErr() error {
	app := cli.NewApp()
	app.Name = "zke"
	app.Version = VERSION
	app.Usage = "ZDNS Kubernetes Engine, an extremely simple, lightning fast Kubernetes installer that works everywhere"
	app.Before = func(ctx *cli.Context) error {
		if ctx.GlobalBool("debug") {
			logrus.SetLevel(logrus.DebugLevel)
		}
		logrus.Debugf("ZKE version %s build at %s", app.Version, BUILD)
		return nil
	}
	app.Author = "ZdnsCloud"
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
