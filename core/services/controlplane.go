package services

import (
	"context"

	"github.com/zdnscloud/zke/core/pki"
	"github.com/zdnscloud/zke/pkg/hosts"
	"github.com/zdnscloud/zke/pkg/log"
	"github.com/zdnscloud/zke/pkg/util"
	"github.com/zdnscloud/zke/types"

	zerrgroup "github.com/zdnscloud/cement/errgroup"
	"golang.org/x/sync/errgroup"
)

func RunControlPlaneTest(ctx context.Context, controlHosts []*hosts.Host, prsMap map[string]types.PrivateRegistry, cpNodePlanMap map[string]types.ZKEConfigNodePlan, updateWorkersOnly bool, alpineImage string, certMap map[string]pki.CertificatePKI) error {
	if updateWorkersOnly {
		return nil
	}
	log.Infof(ctx, "[%s] Building up Controller Plane..", ControlRole)
	var errgrp errgroup.Group

	hostsQueue := util.GetObjectQueue(controlHosts)
	for w := 0; w < WorkerThreads; w++ {
		errgrp.Go(func() error {
			var errList []error
			for host := range hostsQueue {
				runHost := host.(*hosts.Host)
				err := doDeployControlHost(ctx, runHost, prsMap, cpNodePlanMap[runHost.Address].Processes, alpineImage, certMap)
				if err != nil {
					errList = append(errList, err)
				}
			}
			return util.ErrList(errList)
		})
	}
	if err := errgrp.Wait(); err != nil {
		return err
	}
	log.Infof(ctx, "[%s] Successfully started Controller Plane..", ControlRole)
	return nil
}

func RunControlPlane(ctx context.Context, controlHosts []*hosts.Host, prsMap map[string]types.PrivateRegistry, cpNodePlanMap map[string]types.ZKEConfigNodePlan, updateWorkersOnly bool, alpineImage string, certMap map[string]pki.CertificatePKI) error {
	if updateWorkersOnly {
		return nil
	}
	log.Infof(ctx, "[%s] Building up Controller Plane..", ControlRole)

	_, err := zerrgroup.Batch(controlHosts, func(h interface{}) (interface{}, error) {
		runHost := h.(*hosts.Host)
		err := doDeployControlHost(ctx, runHost, prsMap, cpNodePlanMap[runHost.Address].Processes, alpineImage, certMap)
		return nil, err
	})
	if err != nil {
		return err
	}
	/*var errgrp errgroup.Group

	hostsQueue := util.GetObjectQueue(controlHosts)
	for w := 0; w < WorkerThreads; w++ {
		errgrp.Go(func() error {
			var errList []error
			for host := range hostsQueue {
				runHost := host.(*hosts.Host)
				err := doDeployControlHost(ctx, runHost, prsMap, cpNodePlanMap[runHost.Address].Processes, alpineImage, certMap)
				if err != nil {
					errList = append(errList, err)
				}
			}
			return util.ErrList(errList)
		})
	}
	if err := errgrp.Wait(); err != nil {
		return err
	}*/
	log.Infof(ctx, "[%s] Successfully started Controller Plane..", ControlRole)
	return nil
}

func RemoveControlPlane(ctx context.Context, controlHosts []*hosts.Host, force bool) error {
	log.Infof(ctx, "[%s] Tearing down the Controller Plane..", ControlRole)
	var errgrp errgroup.Group
	hostsQueue := util.GetObjectQueue(controlHosts)
	for w := 0; w < WorkerThreads; w++ {
		errgrp.Go(func() error {
			var errList []error
			for host := range hostsQueue {
				runHost := host.(*hosts.Host)
				if err := removeKubeAPI(ctx, runHost); err != nil {
					errList = append(errList, err)
				}
				if err := removeKubeController(ctx, runHost); err != nil {
					errList = append(errList, err)
				}
				if err := removeScheduler(ctx, runHost); err != nil {
					errList = append(errList, err)
				}
				// force is true in remove, false in reconcile
				if force {
					if err := removeKubelet(ctx, runHost); err != nil {
						errList = append(errList, err)
					}
					if err := removeKubeproxy(ctx, runHost); err != nil {
						errList = append(errList, err)
					}
					if err := removeSidekick(ctx, runHost); err != nil {
						errList = append(errList, err)
					}
				}
			}
			return util.ErrList(errList)
		})
	}

	if err := errgrp.Wait(); err != nil {
		return err
	}

	log.Infof(ctx, "[%s] Successfully tore down Controller Plane..", ControlRole)
	return nil
}

func RestartControlPlane(ctx context.Context, controlHosts []*hosts.Host) error {
	log.Infof(ctx, "[%s] Restarting the Controller Plane..", ControlRole)
	var errgrp errgroup.Group

	hostsQueue := util.GetObjectQueue(controlHosts)
	for w := 0; w < WorkerThreads; w++ {
		errgrp.Go(func() error {
			var errList []error
			for host := range hostsQueue {
				runHost := host.(*hosts.Host)
				// restart KubeAPI
				if err := RestartKubeAPI(ctx, runHost); err != nil {
					errList = append(errList, err)
				}

				// restart KubeController
				if err := RestartKubeController(ctx, runHost); err != nil {
					errList = append(errList, err)
				}

				// restart scheduler
				err := RestartScheduler(ctx, runHost)
				if err != nil {
					errList = append(errList, err)
				}
			}
			return util.ErrList(errList)
		})
	}
	if err := errgrp.Wait(); err != nil {
		return err
	}
	log.Infof(ctx, "[%s] Successfully restarted Controller Plane..", ControlRole)
	return nil
}

func doDeployControlHost(ctx context.Context, host *hosts.Host, prsMap map[string]types.PrivateRegistry, processMap map[string]types.Process, alpineImage string, certMap map[string]pki.CertificatePKI) error {
	if host.IsWorker {
		if err := removeNginxProxy(ctx, host); err != nil {
			return err
		}
	}
	// run sidekick
	if err := runSidekick(ctx, host, prsMap, processMap[SidekickContainerName]); err != nil {
		return err
	}
	// run kubeapi
	if err := runKubeAPI(ctx, host, prsMap, processMap[KubeAPIContainerName], alpineImage, certMap); err != nil {
		return err
	}
	// run kubecontroller
	if err := runKubeController(ctx, host, prsMap, processMap[KubeControllerContainerName], alpineImage); err != nil {
		return err
	}
	// run scheduler
	return runScheduler(ctx, host, prsMap, processMap[SchedulerContainerName], alpineImage)
}
