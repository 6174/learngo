package routers

import (
	"6174/cliapp/modules/settings"
	"6174/cliapp/modules/log"
)

func init() {

}

// global configuration
func GlobalInit() {
	settings.NewContext()
	log.Trace("Custom path: %s", settings.CustomPath)
	log.Trace("Log path: %s", settings.LogRootPath)
}