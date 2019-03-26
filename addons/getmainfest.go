package addons

import "github.com/zdnscloud/zke/templates"

func GetManifest(Config interface{}, Provider string) (string, error) {
       var tmplt string
       switch Provider {
       case "coredns":
	       tmplt = templates.CoreDNSTemplate
       case "nginx":
	       tmplt = templates.NginxIngressTemplate
       case "metrics-server":
	       tmplt = templates.MetricsServerTemplate
       }
       return templates.CompileTemplateFromMap(tmplt, Config)
}
