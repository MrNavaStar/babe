package main

import (
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	app := cli.App{
		Name: "babe",
		Commands: []*cli.Command{
			{
				Name:  "relocate",
				Usage: "babe relocate <file.jar> <from:to>...",
				Args:  true,
				Action: func(c *cli.Context) error {
					return RelocateJar(c.Args().First(), ParseRelocations(c.Args().Slice()[1:]))
				},
			},
			{
				Name:  "minimize",
				Usage: "babe minimize <file.jar>",
				Args:  true,
				Action: func(c *cli.Context) error {
					return MinimizeJar(c.Args().First())
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
