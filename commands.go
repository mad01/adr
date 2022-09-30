package main

import (
	"github.com/fatih/color"
	"github.com/urfave/cli"
)

func setCommands(app *cli.App) {
	app.Commands = []cli.Command{
		{
			Name:    "new",
			Aliases: []string{"c"},
			Usage:   "Create a new ADR",
			Flags:   []cli.Flag{},
			Action: func(c *cli.Context) error {
				helper := NewAdrHelper()
				currentConfig := helper.GetConfig()
				currentConfig.CurrentAdr++
				helper.UpdateConfig(currentConfig)
				helper.NewAdr(currentConfig, c.Args())
				return nil
			},
		},

		{
			Name:        "init",
			Aliases:     []string{"i"},
			Usage:       "Initializes the ADR configurations",
			UsageText:   "adr init /home/user/adrs",
			Description: "Initializes the ADR configuration with an optional ADR base directory\n This is a a prerequisite to running any other adr sub-command",
			Action: func(c *cli.Context) error {
				initDir := c.Args().First()
				color.Green("Initializing ADR base at " + initDir)
				helper := NewAdrHelper()
				helper.InitBaseDir(initDir)
				helper.InitConfig()
				helper.InitTemplate()
				return nil
			},
		},
	}
}
