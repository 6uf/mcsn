package main

import (
	"context"
	"crypto/tls"
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
		os.Exit(0)
	} else {
		checkifValid()

		writetoFile(jsonValues{Accounts: AccountsVer, Bearers: Confirmed, Config: ConfigsVer, Names: NamesVer, Vps: VpsesVer})

		for _, acc := range Confirmed {
			bearers.Bearers = append(bearers.Bearers, strings.Split(acc, "`")[0])
			bearers.AccountType = append(bearers.AccountType, strings.Split(acc, "`")[2])
		}
	}
}

func grabArray(array []interface{}) ([]string, error) {
	var list []string; for _, names := range array {list = append(list, names.(string))}; if list == nil || len(list) == 0 {return nil, errors.New("empty")}; return list, nil
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

	BearersVer, err = grabArray(config[`Bearers`].([]interface{}))
	if err != nil {
		bearers, err = mcapi2.Auth(AccountsVer)
		if err != nil {
			fmt.Println(err)
		}

		if len(bearers.Bearers) == 0 {
			log.Println("Unable to authenticate your account(s), please Reverify your login details and make sure the accounts are Microsoft OR Giftcard accounts.")
		} else {
			for i := range bearers.Bearers {
				BearersVer = append(BearersVer, bearers.Bearers[i]+"`"+time.Now().Add(time.Duration(time.Second*86400)).Format(time.RFC850)+"`"+bearers.AccountType[i])
			}
		}
	} else {
		if len(BearersVer) < len(AccountsVer) {
			func() {
				BearersVer = []string{}
				bearers, err = mcapi2.Auth(AccountsVer)
				if err != nil {
					fmt.Println(err)
				}

				if len(bearers.Bearers) == 0 {
				} else {
					for i := range bearers.Bearers {
						BearersVer = append(BearersVer, bearers.Bearers[i]+"`"+time.Now().Add(time.Duration(time.Second*86400)).Format(time.RFC850)+"`"+bearers.AccountType[i] + "`" + strings.Split(AccountsVer[i], ":")[0])
					}
				}
			}()
		} else if len(AccountsVer) < len(BearersVer) {
			var confirmedBearers []string
			for _, acc := range AccountsVer {
				for _, num := range BearersVer {
					if acc == strings.Split(num, "`")[3] {
						confirmedBearers = append(confirmedBearers, strings.Split(num, "`")[0]+"`"+time.Now().Add(time.Duration(time.Second*86400)).Format(time.RFC850)+"`"+strings.Split(num, "`")[2] + "`" + acc)
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
	for number, accs := range BearersVer {
		m := strings.Split(accs, "`")
		f, _ := http.NewRequest("GET", "https://api.minecraftservices.com/minecraft/profile/name/boom/available", nil)
		f.Header.Set("Authorization", "Bearer "+m[0])
		j, _ := http.DefaultClient.Do(f)

		if j.StatusCode == 401 {
			fmt.Printf("Account %v turned up invalid. Attempted to Reauth\n", AccountsVer[number])
			bearers, err = mcapi2.Auth([]string{AccountsVer[number]})
			if err == nil {
				Confirmed = append(Confirmed, bearers.Bearers[0] + "`" + time.Now().Add(time.Duration(time.Second*86400)).Format(time.RFC850)+"`" + bearers.AccountType[0] + "`" + strings.Split(AccountsVer[number], ":")[0])
			}

		} else {
			wdad, _ := time.Parse(time.RFC850, m[1])

			if time.Now().After(wdad) {
				reAuth = append(reAuth, AccountsVer[number])
			} else {
				Confirmed = append(Confirmed, m[0] + "`" + m[1] + "`" + m[2] + "`" + strings.Split(AccountsVer[number], "`")[0])
			}
		}
	}

	if len(reAuth) != 0 {
		log.Println("Reauthing some accounts..")
		bearers, _ = mcapi2.Auth(reAuth)

		for i, acc := range bearers.Bearers {
			if checkChange(acc) == true {
				Confirmed = append(Confirmed, acc + "`" + time.Now().Add(time.Duration(time.Second*86400)).Format(time.RFC850) + "`" + bearers.AccountType[i]+ "`" + strings.Split(AccountsVer[i], ":")[0])
			} else {
				log.Printf("Account #%v cannot be name changed.\n", i)
			}
		}
	}
}

func checkChange(bearer string) bool {
	conn, _ := tls.Dial("tcp", "api.minecraftservices.com"+":443", nil)
	defer conn.Close()

	fmt.Fprintln(conn, "GET /minecraft/profile/namechange HTTP/1.1\r\nHost: api.minecraftservices.com\r\nUser-Agent: Dismal/1.0\r\nAuthorization: Bearer "+bearer+"\r\n\r\n")

	conn.Read(authbytes)

	authbytes = []byte(strings.Split(strings.Split(string(authbytes), "\x00")[0], "\r\n\r\n")[1])
	json.Unmarshal(authbytes, &auth)

	switch auth["nameChangeAllowed"].(bool) {
	case true:
		securityResult = true
	case false:
		securityResult = false
	}

	return securityResult
}
