package main

import (
	"fmt"
	"mcsn/src"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/6uf/apiGO"
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

	header := `
 •   ▄ ·   ▄▄·  ▄▄ ·    ▄ 
·██ ▐███▪▐█ ▌▪▐█ ▀  •█▌▐█
▐█ ▌▐▌▐█·██ ▄▄▄▀▀▀█▄▐█▐▐▌
██ ██▌▐█▌▐███▌▐█▄▪▐███▐█▌
▀▀  █▪▀▀▀·▀▀▀  ▀▀▀▀ ▀▀ █▪
`

	for _, char := range []string{"•", "·", "▪"} {
		header = strings.ReplaceAll(header, char, aurora.Sprintf(aurora.Faint(aurora.Red("%v")), char))
	}
	for _, char := range []string{"█", "▄", "▌", "▀", "▌", "▀"} {
		header = strings.ReplaceAll(header, char, aurora.Sprintf(aurora.Faint(aurora.White("%v")), char))
	}
	for _, char := range []string{"▐"} {
		header = strings.ReplaceAll(header, char, aurora.Sprintf(aurora.Faint(aurora.BrightBlack("%v")), char))
	}

	fmt.Println(header)
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
				Name:    "test",
				Aliases: []string{"t", "ts"},
				Usage:   "Test uses a dummy bearer and email etc to test your speed etc.",
				Action: func(c *cli.Context) error {
					if c.Bool("gc") {
						src.Bearers.Details = append(src.Bearers.Details, apiGO.Info{
							Bearer:      "testbearer",
							AccountType: "Giftcard",
							Email:       "testcommand@mcsn.com",
							Requests:    src.Acc.GcReq,
						})
					} else {
						src.Bearers.Details = append(src.Bearers.Details, apiGO.Info{
							Bearer:      "testbearer",
							AccountType: "Giftcard",
							Email:       "testcommand@mcsn.com",
							Requests:    src.Acc.GcReq,
						})
					}

					src.Snipe("test", 100, "single", "")

					return nil
				},
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "gc",
						Usage: "username to snipe",
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
					src.Snipe(c.String("u"), c.Float64("d"), "proxy", "")
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
							src.Snipe(c.String("u"), c.Float64("d"), "proxy", "3c")
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
							src.Snipe(c.String("u"), c.Float64("d"), "proxy", "3l")
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
							src.Snipe(c.String("u"), c.Float64("d"), "proxy", "3n")
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
						Usage: "Snipe only Numeric names.",
						Action: func(c *cli.Context) error {
							src.AuthAccs()
							fmt.Println()
							go src.CheckAccs()
							src.Snipe(c.String("u"), c.Float64("d"), "proxy", "list")
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
					src.Snipe(c.String("u"), float64(c.Int64("t")), "turbo", "")
					return nil
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "u",
						Usage: "username to snipe",
					},
					&cli.Int64Flag{
						Name:  "t",
						Usage: "the seconds between each request",
						Value: 60,
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
						Usage: "Name of your Art (this will be used for the folder name)",
					},
					&cli.StringFlag{
						Name:  "i",
						Usage: "Name of your image.",
					},
				},
			},
		},
		Name:    "MCSN",
		Usage:   "A name sniper dedicated to premium free services",
		Version: "5.15",
	}

	app.Run(os.Args)
}
