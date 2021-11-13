package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Liza-Developer/mcapi2"
	"github.com/bwmarrin/discordgo"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name: "snipe",
				Usage: `Usages:8
	-u - This option takes input, it uses the input to snipe the name your going for.
	-d - This option is used for your delay, example 50.. this is needed.
	`,
				Action: func(c *cli.Context) error {
					authAccs()
					singlesniper(c.String("u"), c.Float64("d"))
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

					q, _ := ioutil.ReadFile("accounts.json")

					config = mcapi2.GetConfig(q)

					var err error
					s, err = discordgo.New("Bot " + config[`Config`].([]interface{})[3].(string))
					if err != nil {
						log.Fatalf("Invalid bot parameters: %v", err)
					}

					ctx, cancel = context.WithTimeout(context.Background(), 99999999*time.Minute)
					if err != nil {
						log.Println(err)
					}

					s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
						if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
							h(s, i)
						}
					})

					runBot()
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
							auto("3c")
							return nil
						},
					},
					{
						Name:  "3l",
						Usage: "Snipe only Alphabetic names.",
						Action: func(c *cli.Context) error {
							authAccs()
							auto("3l")
							return nil
						},
					},
					{
						Name:  "3n",
						Usage: "Snipe only Numeric names.",
						Action: func(c *cli.Context) error {
							authAccs()
							auto("3n")
							return nil
						},
					},
					{
						Name:  "list",
						Usage: "Snipe names you have added to your `Names` in the config.",
						Action: func(c *cli.Context) error {
							authAccs()
							auto("list")
							return nil
						},
					},
				},
			},

			{
				Name:    "ping",
				Aliases: []string{"p"},
				Usage:   "ping helps give you a rough estimate of your connection to the minecraft API.",
				Action: func(c *cli.Context) error {
					fmt.Printf("Estimated Delay: %v", math.Round(AutoOffset(true)))
					return nil
				},
			},
		},
		HideHelp: false,
	}

	app.Run(os.Args)
}

func authAccs() {

	q, _ := ioutil.ReadFile("accounts.json")

	config = mcapi2.GetConfig(q)

	if config[`Accounts`] == nil {
		log.Println("unable to continue, you have no accounts added.")
		os.Exit(0)
	}

	if config[`Bearers`] == nil || len(config[`Bearers`].([]interface{})) == 0 {
		func() {
			var configOptions []string
			for _, accs := range config[`Accounts`].([]interface{}) {
				Accounts = append(Accounts, accs.(string))
			}

			for _, accs := range config[`Config`].([]interface{}) {
				configOptions = append(configOptions, accs.(string))
			}

			if config[`Names`] != nil {
				for _, accs := range config[`Names`].([]interface{}) {
					names = append(names, accs.(string))
				}
			}

			if config[`Vps`] != nil {
				for _, accs := range config[`Vps`].([]interface{}) {
					vpses = append(vpses, accs.(string))
				}
			}

			bearers, err = mcapi2.Auth(Accounts)
			if err != nil {
				fmt.Println(err)
			}

			if bearers.Bearers == nil {
				fmt.Println("Was unable to auth accs")
				os.Exit(0)
			} else {
				var newBearers []string

				for e, bearerz := range bearers.Bearers {
					newBearers = append(newBearers, bearerz+"`"+time.Now().Add(time.Duration(time.Second*86400)).Format(time.RFC850)+"`"+bearers.AccountType[e])
				}

				v, _ := json.MarshalIndent(jsonValues{Accounts: Accounts, Bearers: newBearers, Config: configOptions, Names: names, Vps: vpses}, "", "  ")

				ioutil.WriteFile("accounts.json", v, 0)

				if len(newBearers) == 0 {
					fmt.Println("No valid accounts.")
					os.Exit(0)
				}
			}
		}()
	} else {

		if len(config[`Bearers`].([]interface{})) > len(config[`Accounts`].([]interface{})) {
			func() {
				var configOptions []string
				for _, accs := range config[`Accounts`].([]interface{}) {
					Accounts = append(Accounts, accs.(string))
				}

				for _, accs := range config[`Config`].([]interface{}) {
					configOptions = append(configOptions, accs.(string))
				}

				if config[`Names`] != nil {
					for _, accs := range config[`Names`].([]interface{}) {
						names = append(names, accs.(string))
					}
				}

				if config[`Vps`] != nil {
					for _, accs := range config[`Vps`].([]interface{}) {
						vpses = append(vpses, accs.(string))
					}
				}

				bearers, err = mcapi2.Auth(Accounts)
				if err != nil {
					fmt.Println(err)
				}

				if bearers.Bearers == nil {
					fmt.Println("Was unable to auth accs")
					os.Exit(0)
				} else {
					var newBearers []string

					for e, bearerz := range bearers.Bearers {
						newBearers = append(newBearers, bearerz+"`"+time.Now().Add(time.Duration(time.Second*86400)).Format(time.RFC850)+"`"+bearers.AccountType[e])
					}

					v, _ := json.MarshalIndent(jsonValues{Accounts: Accounts, Bearers: newBearers, Config: configOptions, Names: names, Vps: vpses}, "", "  ")

					ioutil.WriteFile("accounts.json", v, 0)

					if len(newBearers) == 0 {
						fmt.Println("No valid accounts.")
						os.Exit(0)
					}
				}
			}()
		} else {
			for number, accs := range config[`Bearers`].([]interface{}) {

				m := strings.Split(accs.(string), "`")

				f, _ := http.NewRequest("GET", "https://api.minecraftservices.com/minecraft/profile/name/boom/available", nil)
				f.Header.Set("Authorization", "Bearer "+m[0])
				j, _ := http.DefaultClient.Do(f)

				if j.StatusCode == 401 {
					fmt.Printf("Account %v turned up invalid. Attempted to Reauth\n", config[`Accounts`].([]interface{})[number])

					invalidAccs = append(invalidAccs, config[`Accounts`].([]interface{})[number].(string))

				} else {
					wdad, _ := time.Parse(time.RFC850, m[1])

					if time.Now().After(wdad) {
					} else {
						removeBearers = append(removeBearers, m[0]+"`"+m[1]+"`"+m[2])
						Accounts = append(Accounts, m[0])
						accountType = append(accountType, m[2])
					}
				}
			}

			if len(invalidAccs) != 0 {
				g, _ := mcapi2.Auth(invalidAccs)
				for i, _ := range invalidAccs {
					Accounts = append(Accounts, g.Bearers...)
					removeBearers = append(removeBearers, g.Bearers[i]+"`"+time.Now().Add(time.Duration(time.Second*86400)).Format(time.RFC850)+"`"+g.AccountType[i])
					accountType = append(accountType, g.AccountType[i])
				}
			}

			var accz []string

			for _, accs := range config[`Accounts`].([]interface{}) {
				accz = append(accz, accs.(string))
			}
			for _, accs := range config[`Config`].([]interface{}) {
				configOptions = append(configOptions, accs.(string))
			}

			if config[`Names`] != nil {
				for _, accs := range config[`Names`].([]interface{}) {
					names = append(names, accs.(string))
				}
			}

			if config[`Vps`] != nil {
				for _, accs := range config[`Vps`].([]interface{}) {
					vpses = append(vpses, accs.(string))
				}
			}

			if removeBearers == nil {
				removeBearers = []string{}
			} else {
				bearers.Bearers = Accounts
				bearers.AccountType = accountType
			}

			v, _ := json.MarshalIndent(jsonValues{Accounts: accz, Bearers: removeBearers, Config: configOptions, Names: names, Vps: vpses}, "", "  ")

			ioutil.WriteFile("accounts.json", v, 0)

			if len(bearers.Bearers) == 0 {
				fmt.Println("No valid accounts.. Attempting to reauth..")
				func() {
					var configOptions []string
					Accounts = []string{}
					for _, accs := range config[`Accounts`].([]interface{}) {
						Accounts = append(Accounts, accs.(string))
					}

					for _, accs := range config[`Config`].([]interface{}) {
						configOptions = append(configOptions, accs.(string))
					}

					if config[`Names`] != nil {
						for _, accs := range config[`Names`].([]interface{}) {
							names = append(names, accs.(string))
						}
					}

					if config[`Vps`] != nil {
						for _, accs := range config[`Vps`].([]interface{}) {
							vpses = append(vpses, accs.(string))
						}
					}

					bearers, err = mcapi2.Auth(Accounts)
					if err != nil {
						fmt.Println(err)
					}

					if bearers.Bearers == nil {
						bearers.Bearers = []string{}
					}

					var newBearers []string

					for e, bearerz := range bearers.Bearers {
						newBearers = append(newBearers, bearerz+"`"+time.Now().Add(time.Duration(time.Second*86400)).Format(time.RFC850)+"`"+bearers.AccountType[e])
					}

					v, _ := json.MarshalIndent(jsonValues{Accounts: Accounts, Bearers: newBearers, Config: configOptions, Names: names, Vps: vpses}, "", "  ")

					ioutil.WriteFile("accounts.json", v, 0)

					if len(newBearers) == 0 {
						fmt.Println("No valid accounts.")
						os.Exit(0)
					}
				}()
			}
		}
	}
}
