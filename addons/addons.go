package addons

import (
	"strconv"

	"github.com/zdnscloud/zke/templates"
)

func GetAddonsExecuteJob(addonName, nodeName, image string) (string, error) {
	return getAddonJob(addonName, nodeName, image, false)
}

func getAddonJob(addonName, nodeName, image string, isDelete bool) (string, error) {
	jobConfig := map[string]string{
		"AddonName": addonName,
		"NodeName":  nodeName,
		"Image":     image,
		"DeleteJob": strconv.FormatBool(isDelete),
	}
	return templates.CompileTemplateFromMap(templates.AddonJobTemplate, jobConfig)
}
