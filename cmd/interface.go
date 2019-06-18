package cmd

import (
	"context"
	"github.com/zdnscloud/zke/core"
)

type ZkeCommand interface {
	Create(ctx context.Context, c *core.Cluster) error
	Update(ctx context.Context, c *core.Cluster) error
	Delete(ctx context.Context, c *core.Cluster) error
}
