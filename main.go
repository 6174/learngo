package main

import (
	"os"
	"6174/cliapp/cmd"
	"fmt"
	"github.com/urfave/cli"
)

var Version = "0.0.1+dev";

func main() {
	app := cli.NewApp()
	app.Name = "spm"
	app.Usage = "A shell script package manage tool"
	app.Version = Version
	app.Flags = []cli.Flag {
		cli.StringFlag{
			Name: "lang",
			Value: "English",
			Usage: "Language for greeting",
		},
	}

	app.Commands = []cli.Command {
		cmd.CmdWeb,
		cmd.PublishCommand,
		cli.Command{
			Name: "get",
			Usage: "Install package from url",
			Action: func(c *cli.Context) error {
				fmt.Println("get", c.Args().Get(0), c.String("registry"), c.String("lang"))
				return nil
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name: "registry",
					Usage: "Registry of packages",
				},
			},
		},
	}

	app.Action = func (c *cli.Context) error {
		//cli.DefaultAppComplete(c)
		//cli.HandleExitCoder(errors.New("not an exit coder, though"))
		cli.ShowAppHelp(c)
		//cli.ShowCommandCompletions(c, "nope")
		//cli.ShowCommandHelp(c, "get")
		//cli.ShowCompletions(c)
		//cli.ShowSubcommandHelp(c)
		//cli.ShowVersion(c)
		return nil
	}

	app.Run(os.Args)
}
