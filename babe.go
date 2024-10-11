package main

import (
	"os"

	"github.com/mrnavastar/babe/babe"
	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.App{
		Name: "babe",
		Commands: []*cli.Command{
			{
				Name:  "relocate",
				Args:  true,
				Action: func(c *cli.Context) error {
					return babe.RelocateJar(c.Args().First(), babe.ParseRelocations(c.Args().Slice()[1:]))
				},
			},
			{
				Name:  "minimize",
				Args:  true,
				Action: func(c *cli.Context) error {
					return babe.MinimizeJar(c.Args().First())
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
