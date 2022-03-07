package main

import (
	"fmt"
	"mcsn/src"
	"os"
	"os/exec"
	"runtime"

	"github.com/Liza-Developer/apiGO"
	"github.com/iskaa02/qalam/gradient"
	"github.com/jwalton/go-supportscolor"
	"github.com/logrusorgru/aurora/v3"
	"github.com/urfave/cli/v2"
)

func init() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}

	src.CheckFiles()
	src.Acc.LoadState()

	src.Pro = src.GenProxys()
	src.Setup(src.Pro)

	supportscolor.Stdout()
	g, _ := gradient.
		NewGradient("#FEAC5E", "#C779D0", "#4BC0C8")
	fmt.Println(g.Mutline(`
	• ▌ ▄ ·   ▄▄·  ▄▄ ·    ▄ 
	·██ ▐███▪▐█ ▌▪▐█ ▀  •█▌▐█
	▐█ ▌▐▌▐█·██ ▄▄▄▀▀▀█▄▐█▐▐▌
	██ ██▌▐█▌▐███▌▐█▄▪▐███▐█▌
	▀▀  █▪▀▀▀·▀▀▀  ▀▀▀▀ ▀▀ █▪
	`))

	if src.Acc.DiscordID == "" {
		fmt.Print(aurora.Blink(aurora.Faint(aurora.White("Enter a Discord ID: "))))
		fmt.Scan(&src.Acc.DiscordID)

		src.Acc.SaveConfig()
		src.Acc.LoadState()

		fmt.Println()
	}
}

func main() {

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:  "snipe",
				Usage: `Snipes names onto an src.Account.`,
				Action: func(c *cli.Context) error {
					src.AuthAccs()
					fmt.Println()
					go src.CheckAccs()
					src.Snipe(c.String("u"), c.Float64("d"), "single", "")
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
				Name:    "auto",
				Aliases: []string{"as", "a"},
				Usage:   "Auto sniper attempts to snipe upcoming 3 character usernames.",
				Subcommands: []*cli.Command{
					{
						Name:  "3c",
						Usage: "Snipe names are are a combination of Numeric and Alphabetic.",
						Action: func(c *cli.Context) error {
							src.AuthAccs()
							fmt.Println()
							go src.CheckAccs()
							src.Snipe("", c.Float64("d"), "auto", "3c")
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
							src.AuthAccs()
							fmt.Println()
							go src.CheckAccs()
							src.Snipe("", c.Float64("d"), "auto", "3l")
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
							src.AuthAccs()
							fmt.Println()
							go src.CheckAccs()
							src.Snipe("", c.Float64("d"), "auto", "3n")
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
							src.AuthAccs()
							fmt.Println()
							go src.CheckAccs()
							src.Snipe("", c.Float64("d"), "auto", "list")
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

					delay, time := src.MeanPing()
					fmt.Printf("Estimated (Mean) Delay: %v ~ Took: %v\n", delay, time)

					return nil
				},
			},

			{
				Name:    "proxy",
				Aliases: []string{"p"},
				Usage:   "Proxy snipes names for you",
				Action: func(c *cli.Context) error {
					src.AuthAccs()
					fmt.Println()
					go src.CheckAccs()
					src.Proxy(c.String("u"), c.Float64("d"), apiGO.DropTime(c.String("u")))

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
							src.AuthAccs()
							fmt.Println()
							go src.CheckAccs()

							var names []string
							var drops []int64
							delay := c.Float64("d")
							names, drops = src.ThreeLetters("3c")

							for e, name := range names {
								if delay == 0 {
									delay = float64(src.AutoOffset())
								}

								if !src.Acc.ManualBearer {
									if len(src.Bearers.Details) == 0 {
										fmt.Print(aurora.Faint(aurora.White("No more usable Account(s)\n")))
										os.Exit(0)
									}
								}

								src.Proxy(name, c.Float64("d"), drops[e])

								src.Setup(src.Pro)

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
							src.AuthAccs()
							fmt.Println()
							go src.CheckAccs()

							var names []string
							var drops []int64
							delay := c.Float64("d")
							names, drops = src.ThreeLetters("3l")

							for e, name := range names {
								if delay == 0 {
									delay = float64(src.AutoOffset())
								}

								if !src.Acc.ManualBearer {
									if len(src.Bearers.Details) == 0 {
										fmt.Print(aurora.Faint(aurora.White("No more usable Account(s)\n")))
										os.Exit(0)
									}
								}

								src.Proxy(name, c.Float64("d"), drops[e])

								src.Setup(src.Pro)

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
							src.AuthAccs()
							fmt.Println()
							go src.CheckAccs()

							var names []string
							var drops []int64
							delay := c.Float64("d")
							names, drops = src.ThreeLetters("3n")

							for e, name := range names {
								if delay == 0 {
									delay = float64(src.AutoOffset())
								}

								if !src.Acc.ManualBearer {
									if len(src.Bearers.Details) == 0 {
										fmt.Print(aurora.Faint(aurora.White("No more usable Account(s)\n")))
										os.Exit(0)
									}
								}

								src.Proxy(name, c.Float64("d"), drops[e])

								src.Setup(src.Pro)

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
					src.AuthAccs()
					fmt.Println()
					go src.CheckAccs()
					src.Snipe(c.String("u"), 0, "turbo", "")
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
					src.Skinart(c.String("n"), c.String("i"))
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
		Version:  "4.60",
	}

	app.Run(os.Args)
}
