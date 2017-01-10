package cmd

import (
	"github.com/urfave/cli"
	"6174/cliapp/modules/settings"
	"6174/cliapp/modules/log"
	"6174/cliapp/modules/util"
)

var PublishCommand = cli.Command{
	Name: "publish",
	Usage: "Publish package",
	Description: "Publish package to remote repositry",
	Action: publishPackage,
	Flags: []cli.Flag {
		cli.StringFlag {
			Name: "repositry, r",
			Value: "http://spm.idcos.com",
			Usage: "Remote repositry",
		},
	},
}

func publishPackage(ctx *cli.Context) error {
	settings.NewContext()
	log.Info("publish package in the current directory")
	// find the working directory
	pwd := util.PWD()
	log.Info("Current directory is: %s", pwd)

	log.Info("end")
	return nil
}