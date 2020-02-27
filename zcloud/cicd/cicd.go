package cicd

import (
	"github.com/zdnscloud/zke/core"
)

const deployNamespace = "zcloud"

func GetDeployConfig(c *core.Cluster) map[string]string {
	return map[string]string{
		"DeployNamespace":          deployNamespace,
		"ControllerImage":          c.Image.CICD.Controller,
		"KubeConfigWriterImage":    c.Image.CICD.KubeConfigWriter,
		"CredsIniterImage":         c.Image.CICD.CredsIniter,
		"GitIniterImage":           c.Image.CICD.GitIniter,
		"EntrypointerImage":        c.Image.CICD.Entrypointer,
		"ImageDigestExporterImage": c.Image.CICD.ImageDigestExporter,
		"PullRequestIniterImage":   c.Image.CICD.PullRequestIniter,
		"GCSFetcherImage":          c.Image.CICD.GCSFetcher,
		"WebhookImage":             c.Image.CICD.Webhook,
		"DashboardImage":           c.Image.CICD.Dashboard,
	}
}
