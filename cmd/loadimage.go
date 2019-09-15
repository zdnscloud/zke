package cmd

import (
	"context"
	"fmt"

	"github.com/zdnscloud/zke/core"
	"github.com/zdnscloud/zke/core/pki"
	"github.com/zdnscloud/zke/pkg/docker"
	"github.com/zdnscloud/zke/pkg/hosts"
	"github.com/zdnscloud/zke/pkg/log"

	"github.com/urfave/cli"
	"github.com/zdnscloud/cement/errgroup"
	cementlog "github.com/zdnscloud/cement/log"
)

func LoadImageCommand() cli.Command {
	loadImageFlags := []cli.Flag{
		cli.StringFlag{
			Name:  "image-file",
			Usage: "Specify the images tar file",
			Value: "zcloud-images.tar",
		},
	}
	return cli.Command{
		Name:   "loadimage",
		Usage:  "load images for the cluster nodes",
		Action: loadImageFromCli,
		Flags:  loadImageFlags,
	}
}

func loadImageFromCli(cliCtx *cli.Context) error {
	imagePath := cliCtx.String("image-file")

	parentCtx := context.Background()
	logger := cementlog.NewLog4jConsoleLogger(log.LogLevel)
	defer logger.Close()
	ctx, err := log.SetLogger(parentCtx, logger)
	if err != nil {
		return err
	}

	clusterFile, err := resolveClusterFile(pki.ClusterConfig)
	if err != nil {
		return fmt.Errorf("Failed to resolve cluster file: %v", err)
	}

	zkeConfig, err := core.ParseConfig(ctx, clusterFile)
	if err != nil {
		return fmt.Errorf("Failed to parse cluster file: %v", err)
	}

	kubeCluster, err := core.InitClusterObject(ctx, zkeConfig)
	if err != nil {
		return err
	}
	if err := kubeCluster.SetupDialers(ctx, hosts.DialersOptions{}); err != nil {
		return err
	}
	err = kubeCluster.TunnelHosts(ctx)
	if err != nil {
		return err
	}
	log.Debugf(ctx, "Starting load images")

	allHosts := hosts.GetUniqueHostList(kubeCluster.ControlPlaneHosts, kubeCluster.EtcdHosts, kubeCluster.WorkerHosts, kubeCluster.EdgeHosts)

	_, err = errgroup.Batch(allHosts, func(h interface{}) (interface{}, error) {
		runHost := h.(*hosts.Host)
		return nil, docker.LoadImage(ctx, runHost.DClient, runHost.Address, imagePath)
	})

	if err != nil {
		log.Fatal(ctx, "load image failed %s", err)
	}

	return nil
}
