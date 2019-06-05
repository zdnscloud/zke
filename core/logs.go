package core

import (
	"context"

	"github.com/zdnscloud/zke/pkg/hosts"

	"github.com/zdnscloud/cement/errgroup"
)

func (c *Cluster) CleanDeadLogs(ctx context.Context) error {
	hostList := hosts.GetUniqueHostList(c.EtcdHosts, c.ControlPlaneHosts, c.WorkerHosts, c.StorageHosts, c.EdgeHosts)

	_, err := errgroup.Batch(hostList, func(h interface{}) (interface{}, error) {
		return nil, hosts.DoRunLogCleaner(ctx, h.(*hosts.Host), c.SystemImages.Alpine, c.PrivateRegistriesMap)
	})
	return err
}
