package main

import (
	"bytes"
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
				Usage: `-u - This option takes input, it uses the input to snipe the name your going for.
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
					s, err = discordgo.New("Bot " + config[`DiscordBotToken`].(string))
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
		HideHelp:    true,
		Name:        "MCSN",
		Description: "A name sniper dedicated to premium free services",
		Version:     "3.4",
	}

	app.Run(os.Args)
}

func authAccs() {
	q, _ := ioutil.ReadFile("accounts.json")

	config = mcapi2.GetConfig(q)

	grabDetails()

	if BearersVer == nil {
		fmt.Println("\nNo bearers have been found, please check your details.")
		os.Exit(0)
	} else {
		checkifValid()

		config["Accounts"] = AccountsVer
		config["Bearers"] = Confirmed

		writetoFile(config)

		for _, acc := range Confirmed {
			bearers.Bearers = append(bearers.Bearers, strings.Split(acc, "`")[0])
			bearers.AccountType = append(bearers.AccountType, strings.Split(acc, "`")[2])
		}
	}

	if bearers.Bearers == nil {
		fmt.Println("\nFailed to authorize your bearers, please rerun the sniper.")
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
	var empty bool

	AccountsVer, err = grabArray(config[`Accounts`].([]interface{}))
	if err != nil {
		log.Println("unable to continue, you have no accounts added.")
		os.Exit(0)
	}

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

func check(status, name, unixTime string) (string, string, bool) {
	var bearerGot string
	var emailGot string
	var send bool
	if status == `200` {
		for _, bearer := range bearers.Bearers {
			for _, email := range AccountsVer {
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
					emailGot = email

					type data struct {
						Name   string `json:"name"`
						Bearer string `json:"bearer"`
						Unix   string `json:"unix"`
						Config string `json:"config"`
					}

					body, err := json.Marshal(data{Name: name, Bearer: bearerGot, Unix: unixTime, Config: string(jsonValue(embeds{Content: fmt.Sprintf("<@&%v>", config["RoleID"]), Embeds: []embed{{Description: fmt.Sprintf("Succesfully sniped %v with %v searches :skull:", name, searches(name)), Color: 770000, Footer: footer{Text: "MCSN"}, Time: time.Now().Format(time.RFC3339)}}}))})

					if err == nil {
						req, err := http.NewRequest("POST", "https://droptime.site/api/v2/webhook", bytes.NewBuffer(body))
						if err != nil {
							fmt.Println(err)
						}

						req.Header.Set("Content-Type", "application/json")

						resp, err := http.DefaultClient.Do(req)
						if err == nil {
							if resp.StatusCode == 200 {
								send = true
							} else {
								send = false
							}
						}
					}
				}
			}
		}
	}

	return bearerGot, emailGot, send
}
