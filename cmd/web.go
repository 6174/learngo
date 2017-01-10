package cmd

import (
	"github.com/urfave/cli"
	"github.com/kataras/iris"
	"fmt"
	"6174/cliapp/modules/routers/website"
	"6174/cliapp/modules/routers"
)

var CmdWeb = cli.Command {
	Name: "web",
	Usage: "Start spm web server",
	Description: "Spm web server will create a script package repositry server and web app",
	Action: runWeb,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name: "port, p",
			Value: "3000",
			Usage: "Web server port number",
		},
		cli.StringFlag{
			Name: "config, c",
			Value: "custom/conf/app.ini",
			Usage: "Custom configuration file path",
		},
	},
}

func runWeb (ctx *cli.Context) error {
	if ctx.IsSet("config") {
		fmt.Println(ctx.String("config"))
	}
	routers.GlobalInit()
	iris.Get("/", website.Home)

	iris.Listen(":" + ctx.String("port"))

	return nil
}