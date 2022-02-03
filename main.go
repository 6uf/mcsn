package main

import (
	"fmt"
	"os"

	"github.com/Liza-Developer/apiGO"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:  "snipe",
				Usage: `Snipes names onto an account.`,
				Action: func(c *cli.Context) error {
					authAccs()
					fmt.Println()
					go checkAccs()
					snipe(c.String("u"), c.Float64("d"), "single", "")
					return nil
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "u",
						Usage: "username to snipe",
						Value: "",
					},
					&cli.Float64Flag{
						Name:        "d",
						DefaultText: "1",
						Usage:       "Snipes a few ms earlier so you can counter ping lag.",
						Value:       0,
					},
				},
			},

			{
				Name:    "botsniper",
				Aliases: []string{"bot", "b", "bs"},
				Usage:   "Runs the discord bot sniper.",
				Action: func(c *cli.Context) error {
					authAccs()
					apiGO.StartDigital()
					go checkAccs()
					go apiGO.TaskThread()
					apiGO.Bot()
					return nil
				},
			},

			{
				Name:    "auto",
				Aliases: []string{"as", "autosniper", "a"},
				Usage:   "Auto sniper automatically snipes 3C, 3L, or 3N for you. -3c -3l -3n are the commands.",
				Subcommands: []*cli.Command{
					{
						Name:  "3c",
						Usage: "Snipe names are are a combination of Numeric and Alphabetic.",
						Action: func(c *cli.Context) error {
							authAccs()
							fmt.Println()
							go checkAccs()
							snipe("", c.Float64("d"), "auto", "3c")
							return nil
						},
						Flags: []cli.Flag{
							&cli.Float64Flag{
								Name:  "d",
								Usage: "Snipes a few ms earlier so you can counter ping lag.",
								Value: 0,
							},
						},
					},
					{
						Name:  "3l",
						Usage: "Snipe only Alphabetic names.",
						Action: func(c *cli.Context) error {
							authAccs()
							fmt.Println()
							go checkAccs()
							snipe("", c.Float64("d"), "auto", "3l")
							return nil
						},
						Flags: []cli.Flag{
							&cli.Float64Flag{
								Name:  "d",
								Usage: "Snipes a few ms earlier so you can counter ping lag.",
								Value: 0,
							},
						},
					},
					{
						Name:  "3n",
						Usage: "Snipe only Numeric names.",
						Action: func(c *cli.Context) error {
							authAccs()
							fmt.Println()
							go checkAccs()
							snipe("", c.Float64("d"), "auto", "3n")
							return nil
						},
						Flags: []cli.Flag{
							&cli.Float64Flag{
								Name:  "d",
								Usage: "Snipes a few ms earlier so you can counter ping lag.",
								Value: 0,
							},
						},
					},

					{
						Name:  "list",
						Usage: "Snipe names are are a combination of Numeric and Alphabetic.",
						Action: func(c *cli.Context) error {
							authAccs()
							fmt.Println()
							go checkAccs()
							snipe("", c.Float64("d"), "auto", "list")
							return nil
						},
						Flags: []cli.Flag{
							&cli.Float64Flag{
								Name:  "d",
								Usage: "Snipes a few ms earlier so you can counter ping lag.",
								Value: 0,
							},
						},
					},
				},
			},

			{
				Name:    "ping",
				Aliases: []string{"p"},
				Usage:   "ping helps give you a rough estimate of your connection to the minecraft API.",
				Action: func(c *cli.Context) error {

					delay, time := MeanPing()
					sendS(fmt.Sprintf("Estimated (Mean) Delay: %v ~ Took: %v\n", delay, time))

					return nil
				},
			},

			{
				Name:  "namemc",
				Usage: `NameMC Skin Art`,
				Action: func(c *cli.Context) error {
					skinart(c.String("n"))
					return nil
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "n",
						Usage: "Name of your Art",
						Value: "",
					},
				},
			},
		},

		HideHelp: false,
		Name:     "MCSN",
		Usage:    "A name sniper dedicated to premium free services",
		Version:  "4.30",
	}

	app.Run(os.Args)
}
