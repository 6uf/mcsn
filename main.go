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
					},
					&cli.Float64Flag{
						Name:  "d",
						Usage: "Snipes a few ms earlier so you can counter ping lag.",
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
				Aliases: []string{"as", "a"},
				Usage:   "Auto sniper attempts to snipe upcoming 3 character usernames.",
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
				Name:    "proxy",
				Aliases: []string{"p"},
				Usage:   "Proxy snipes names for you",
				Action: func(c *cli.Context) error {
					authAccs()
					fmt.Println()
					go checkAccs()
					proxy(c.String("u"), c.Float64("d"), apiGO.DropTime(c.String("u")))

					return nil
				},

				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "u",
						Usage: "username to snipe",
					},
					&cli.Float64Flag{
						Name:  "d",
						Usage: "Snipes a few ms earlier so you can counter ping lag.",
					},
				},

				Subcommands: []*cli.Command{
					{
						Name:  "3c",
						Usage: "Snipe names are are a combination of Numeric and Alphabetic.",
						Action: func(c *cli.Context) error {
							authAccs()
							fmt.Println()
							go checkAccs()

							var names []string
							var drops []int64
							delay := c.Float64("d")
							names, drops = threeLetters("3c")

							for e, name := range names {
								if delay == 0 {
									delay = float64(AutoOffset())
								}

								if !acc.ManualBearer {
									if len(bearers.Details) == 0 {
										sendE("No more usable account(s)")
										os.Exit(0)
									}
								}

								proxy(name, c.Float64("d"), drops[e])

								checkVer(name, delay, drops[e])

								fmt.Println()
							}
							return nil
						},
						Flags: []cli.Flag{
							&cli.Float64Flag{
								Name:  "d",
								Usage: "Snipes a few ms earlier so you can counter ping lag.",
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

							var names []string
							var drops []int64
							delay := c.Float64("d")
							names, drops = threeLetters("3l")

							for e, name := range names {
								if delay == 0 {
									delay = float64(AutoOffset())
								}

								if !acc.ManualBearer {
									if len(bearers.Details) == 0 {
										sendE("No more usable account(s)")
										os.Exit(0)
									}
								}

								proxy(name, c.Float64("d"), drops[e])

								checkVer(name, delay, drops[e])

								fmt.Println()
							}

							return nil
						},
						Flags: []cli.Flag{
							&cli.Float64Flag{
								Name:  "d",
								Usage: "Snipes a few ms earlier so you can counter ping lag.",
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

							var names []string
							var drops []int64
							delay := c.Float64("d")
							names, drops = threeLetters("3n")

							for e, name := range names {
								if delay == 0 {
									delay = float64(AutoOffset())
								}

								if !acc.ManualBearer {
									if len(bearers.Details) == 0 {
										sendE("No more usable account(s)")
										os.Exit(0)
									}
								}

								proxy(name, c.Float64("d"), drops[e])

								checkVer(name, delay, drops[e])

								fmt.Println()
							}

							return nil
						},
						Flags: []cli.Flag{
							&cli.Float64Flag{
								Name:  "d",
								Usage: "Snipes a few ms earlier so you can counter ping lag.",
							},
						},
					},
				},
			},

			{
				Name:    "turbo",
				Aliases: []string{"t"},
				Usage:   "Turbo a name just in case it drops!",
				Action: func(c *cli.Context) error {
					authAccs()
					fmt.Println()
					go checkAccs()
					snipe(c.String("u"), 0, "turbo", "")
					return nil
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "u",
						Usage: "username to snipe",
					},
				},
			},

			{
				Name:    "namemc",
				Aliases: []string{"n", "nmc", "skinart"},
				Usage:   `NameMC Skin Art`,
				Action: func(c *cli.Context) error {
					skinart(c.String("n"), c.String("i"))
					return nil
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "n",
						Usage: "Name of your Art",
					},
					&cli.StringFlag{
						Name:  "i",
						Usage: "Name of your image.",
					},
				},
			},
		},

		HideHelp: false,
		Name:     "MCSN",
		Usage:    "A name sniper dedicated to premium free services",
		Version:  "4.50",
	}

	app.Run(os.Args)
}
