package main

import (
	"context"
	"encoding/json"
	"errors"
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
							auto("3c", c.Float64("d"))
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
							auto("3l", c.Float64("d"))
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
							auto("3n", c.Float64("d"))
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
						Usage: "Snipe names you have added to your `Names` in the config.",
						Action: func(c *cli.Context) error {
							authAccs()
							auto("list", 0)
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
					fmt.Printf("Estimated Delay: %v\n", math.Round(AutoOffset(true)))
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

	grabDetails()

	if BearersVer == nil {
		fmt.Println("No bearers have been found, please check your details.")
		os.Exit(0)
	} else {
		checkifValid()

		writetoFile(jsonValues{Accounts: AccountsVer, Bearers: Confirmed, Config: ConfigsVer, Names: NamesVer, Vps: VpsesVer, RoleID: config[`RoleID`].(string)})

		for _, acc := range Confirmed {
			bearers.Bearers = append(bearers.Bearers, strings.Split(acc, "`")[0])
			bearers.AccountType = append(bearers.AccountType, strings.Split(acc, "`")[2])
		}
	}

	if bearers.Bearers == nil {
		fmt.Println("Failed to authorize your bearers, please rerun the sniper.")
		os.Exit(0)
	}
}

func grabArray(array []interface{}) ([]string, error) {
	var list []string

	if array != nil {
		for _, names := range array {
			list = append(list, names.(string))
		}
		if len(list) == 0 {
			return nil, errors.New("empty")
		}
	} else {
		return nil, errors.New("empty")
	}

	return list, nil
}

func writetoFile(str interface{}) {
	v, _ := json.MarshalIndent(str, "", "  ")

	ioutil.WriteFile("accounts.json", v, 0)
}

func grabDetails() {
	AccountsVer, err = grabArray(config[`Accounts`].([]interface{}))
	if err != nil {
		log.Println("unable to continue, you have no accounts added.")
		os.Exit(0)
	}
	var empty bool
	if config[`Bearers`] == nil {
		empty = true
	} else {
		BearersVer, _ = grabArray(config[`Bearers`].([]interface{}))
	}

	if empty {
		bearerz, err := mcapi2.Auth(AccountsVer)
		if err != nil {
			fmt.Println(err)
		}

		if len(bearerz.Bearers) == 0 {
			log.Println("Unable to authenticate your account(s), please Reverify your login details and make sure the accounts are Microsoft OR Giftcard accounts.")
		} else {
			for i := range bearerz.Bearers {
				BearersVer = append(BearersVer, bearerz.Bearers[i]+"`"+time.Now().Add(time.Duration(time.Second*86400)).Format(time.RFC850)+"`"+bearerz.AccountType[i]+"`"+AccountsVer[i])
			}
		}
	} else {
		if len(BearersVer) < len(AccountsVer) {
			func() {
				BearersVer = []string{}
				bearerz, err := mcapi2.Auth(AccountsVer)
				if err != nil {
					fmt.Println(err)
				}

				if len(bearerz.Bearers) == 0 {
				} else {
					for i := range bearerz.Bearers {
						BearersVer = append(BearersVer, bearerz.Bearers[i]+"`"+time.Now().Add(time.Duration(time.Second*86400)).Format(time.RFC850)+"`"+bearerz.AccountType[i]+"`"+AccountsVer[i])
					}
				}
			}()
		} else if len(AccountsVer) < len(BearersVer) {
			var confirmedBearers []string
			for _, acc := range AccountsVer {
				for _, num := range BearersVer {
					if acc == strings.Split(num, "`")[3] {
						confirmedBearers = append(confirmedBearers, strings.Split(num, "`")[0]+"`"+time.Now().Add(time.Duration(time.Second*86400)).Format(time.RFC850)+"`"+strings.Split(num, "`")[2]+"`"+acc)
					}
				}
			}

			BearersVer = confirmedBearers
		}
	}

	ConfigsVer, err = grabArray(config[`Config`].([]interface{}))
	if err != nil {
		log.Println("could not find anything from Config in your accounts.json THIS is a critical error, please head to the github and reinstall the accounts.json file")
		os.Exit(0)
	}

	NamesVer, err = grabArray(config[`Names`].([]interface{}))
	if err != nil {
		log.Println("could not find anything from Names in your accounts.json THIS is not a required value.")
	}

	VpsesVer, err = grabArray(config[`Vps`].([]interface{}))
	if err != nil {
		log.Println("could not find anything from Vps in your accounts.json THIS is not a required value.")
		VpsesVer = append(VpsesVer, "placeholder")
	}
}

func checkifValid() {
	var reAuth []string
	for _, accs := range BearersVer {
		m := strings.Split(accs, "`")
		f, _ := http.NewRequest("GET", "https://api.minecraftservices.com/minecraft/profile/name/boom/available", nil)
		f.Header.Set("Authorization", "Bearer "+m[0])
		j, _ := http.DefaultClient.Do(f)

		if j.StatusCode == 401 {
			fmt.Printf("Account %v turned up invalid. Attempting to Reauth\n", m[3])
			reAuth = append(reAuth, m[3])
		} else {
			wdad, _ := time.Parse(time.RFC850, m[1])

			if time.Now().After(wdad) {
				reAuth = append(reAuth, m[3])
			} else {
				Confirmed = append(Confirmed, m[0]+"`"+m[1]+"`"+m[2]+"`"+m[3])
			}
		}
	}

	if len(reAuth) != 0 {
		log.Printf("Reauthing %v accounts..\n", len(reAuth))
		bearerz, _ := mcapi2.Auth(reAuth)

		for i, acc := range bearerz.Bearers {
			Confirmed = append(Confirmed, acc+"`"+time.Now().Add(time.Duration(time.Second*86400)).Format(time.RFC850)+"`"+bearerz.AccountType[i]+"`"+AccountsVer[i])
		}
	}
}

func check(status, name string) string {
	var bearerGot string
	if status == `200` {
		for _, bearer := range bearers.Bearers {

			httpReq, err := http.NewRequest("GET", "https://api.minecraftservices.com/minecraft/profile", nil)
			if err != nil {
				continue
			}
			httpReq.Header.Set("Authorization", "Bearer "+bearer)

			uwu, err := http.DefaultClient.Do(httpReq)
			if err != nil {
				continue
			}

			bodyByte, err := ioutil.ReadAll(uwu.Body)
			if err != nil {
				continue
			}

			var info map[string]interface{}
			json.Unmarshal(bodyByte, &info)

			if info[`name`] == nil {
			} else if info[`name`] == name {
				bearerGot = bearer
			}
		}
	}

	return bearerGot
}
