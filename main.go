package main

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/Liza-Developer/apiGO"
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
					sendI(fmt.Sprintf("Estimated Delay: %v\n", math.Round(AutoOffset())))
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

		HideHelp:    false,
		Name:        "MCSN",
		Description: "A name sniper dedicated to premium free services",
		Version:     "3.4",
	}

	app.Run(os.Args)
}

func snipe(name string, delay float64, option string, charType string) {
	var useAuto bool = false

	switch option {
	case "single":
		dropTime := apiGO.DropTime(name)
		if dropTime < int64(10000) {
			sendW("-!- Droptime [UNIX] : ")
			fmt.Scan(&dropTime)
			fmt.Println()
		}

		checkVer(name, delay, dropTime)

	case "auto":
		if delay == 0 {
			useAuto = true
		}

		for {

			names, drops := threeLetters(charType)

			for e, name := range names {
				if useAuto {
					delay = AutoOffset()
				}

				if !config["ManualBearer"].(bool) {
					if len(bearers.Bearers) == 0 {
						authAccs()
						if bearers.Bearers == nil {
							sendE("No more usable account(s)")
							os.Exit(0)
						}
					}
				}

				checkVer(name, delay, drops[e])

				fmt.Println()
			}

			if charType == "list" {
				break
			}
		}
	}

	fmt.Println()

	sendW("Press CTRL+C to Continue : ")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
}

func authAccs() {
	q, _ := ioutil.ReadFile("config.json")

	config = apiGO.GetConfig(q)

	grabDetails()

	if !config["ManualBearer"].(bool) {
		if BearersVer == nil {
			sendE("No bearers have been found, please check your details.")
			os.Exit(0)
		} else {

			checkifValid()

			fmt.Println()

			config["Bearers"] = Confirmed

			writetoFile(config)

			for _, acc := range Confirmed {
				bearers.Bearers = append(bearers.Bearers, strings.Split(acc, "`")[0])
				bearers.AccountType = append(bearers.AccountType, strings.Split(acc, "`")[2])
			}

			if bearers.Bearers == nil {
				sendE("Failed to authorize your bearers, please rerun the sniper.")
				os.Exit(0)
			}
		}
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

	ioutil.WriteFile("config.json", v, 0)
}

func grabDetails() {
	var empty bool = false
	file, _ := os.Open("accounts.txt")

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		AccountsVer = append(AccountsVer, scanner.Text())
	}

	if len(AccountsVer) == 0 {
		sendE("Unable to continue, you have no accounts added.\n")
		return
	}

	if config["ManualBearer"].(bool) {
		for _, bearer := range AccountsVer {
			if apiGO.CheckChange(bearer) {
				bearers.Bearers = append(bearers.Bearers, bearer)
				bearers.AccountType = append(bearers.AccountType, isGC(bearer))
			}

			time.Sleep(time.Second)
		}
	} else {
		if config[`Bearers`] == nil {
			bearerz, err := apiGO.Auth(AccountsVer)
			if err != nil {
				sendE(fmt.Sprintf("%v", err))
				os.Exit(0)
			}

			if len(bearerz.Bearers) == 0 {
				sendE("Unable to authenticate your account(s), please Reverify your login details.\n")
				return
			} else {
				for i := range bearerz.Bearers {
					BearersVer = append(BearersVer, bearerz.Bearers[i]+"`"+time.Now().Add(time.Duration(time.Second*86400)).Format(time.RFC850)+"`"+bearerz.AccountType[i]+"`"+AccountsVer[i])
				}
			}
		} else {
			if config[`Bearers`] == nil {
				empty = true
			} else {
				BearersVer, _ = grabArray(config[`Bearers`].([]interface{}))
			}

			if empty {
				bearerz, err := apiGO.Auth(AccountsVer)
				if err != nil {
					fmt.Println(err)
				}

				if len(bearerz.Bearers) == 0 {
					sendE("Unable to authenticate your account(s), please Reverify your login details.")
				} else {
					for i := range bearerz.Bearers {
						BearersVer = append(BearersVer, bearerz.Bearers[i]+"`"+time.Now().Add(time.Duration(time.Second*86400)).Format(time.RFC850)+"`"+bearerz.AccountType[i]+"`"+AccountsVer[i])
					}
				}
			} else {
				if len(BearersVer) < len(AccountsVer) {
					check := make(map[string]bool)
					var acc []string

					for _, i := range BearersVer {
						check[strings.Split(i, "`")[3]] = true
					}

					for _, accs := range AccountsVer {
						if !check[accs] {
							acc = append(acc, accs)
						}
					}

					bearerz, _ := apiGO.Auth(acc)

					if len(bearerz.Bearers) != 0 {
						for i := range bearerz.Bearers {
							BearersVer = append(BearersVer, bearerz.Bearers[i]+"`"+time.Now().Add(time.Duration(time.Second*86400)).Format(time.RFC850)+"`"+bearerz.AccountType[i]+"`"+AccountsVer[i])
						}
					}
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
	}
}

func isGC(bearer string) string {
	var accountT string
	conn, _ := tls.Dial("tcp", "api.minecraftservices.com"+":443", nil)

	fmt.Fprintln(conn, "GET /minecraft/profile/namechange HTTP/1.1\r\nHost: api.minecraftservices.com\r\nUser-Agent: Dismal/1.0\r\nAuthorization: Bearer "+bearer+"\r\n\r\n")

	e := make([]byte, 12)
	conn.Read(e)

	switch string(e[9:12]) {
	case `404`:
		accountT = "Giftcard"
	default:
		accountT = "Microsoft"
	}

	return accountT
}

func checkifValid() {
	var reAuth []string
	for _, accs := range BearersVer {
		m := strings.Split(accs, "`")
		f, _ := http.NewRequest("GET", "https://api.minecraftservices.com/minecraft/profile/name/boom/available", nil)
		f.Header.Set("Authorization", "Bearer "+m[0])
		j, _ := http.DefaultClient.Do(f)

		if j.StatusCode == 401 {
			sendI(fmt.Sprintf("Account %v turned up invalid. Attempting to Reauth\n", strings.Split(m[3], ":")[0]))
			reAuth = append(reAuth, m[3])
		} else {
			wdad, _ := time.Parse(time.RFC850, m[1])

			if time.Now().After(wdad) {
				reAuth = append(reAuth, m[3])
			} else {
				if apiGO.CheckChange(m[0]) {
					Confirmed = append(Confirmed, m[0]+"`"+m[1]+"`"+m[2]+"`"+m[3])
				} else {
					sendI(fmt.Sprintf("Account %v cant name change\n", strings.Split(m[3], ":")[0]))
				}
			}
		}
	}

	if len(reAuth) != 0 {
		sendI(fmt.Sprintf("Reauthing %v accounts..\n", len(reAuth)))
		bearerz, _ := apiGO.Auth(reAuth)

		for i, acc := range bearers.Bearers {
			if apiGO.CheckChange(acc) {
				Confirmed = append(Confirmed, acc+"`"+time.Now().Add(time.Duration(time.Second*86400)).Format(time.RFC850)+"`"+bearerz.AccountType[i]+"`"+AccountsVer[i])
			} else {
				sendI(fmt.Sprintf("Account %v cant name change\n", len(acc[0:5])))
			}
		}
	}
}

func checkVer(name string, delay float64, dropTime int64) {
	var content string
	var sendTime []time.Time
	var leng float64
	var recv []time.Time
	var statusCode []string

	sendI(fmt.Sprintf("Name: %v | Delay: %v\n", name, delay))

	var wg sync.WaitGroup

	apiGO.PreSleep(dropTime)

	payload := bearers.CreatePayloads(name)

	fmt.Println()

	apiGO.Sleep(dropTime, delay)

	fmt.Println()

	for e, account := range payload.AccountType {
		switch account {
		case "Giftcard":
			leng = config[`GcReq`].(float64)
		case "Microsoft":
			leng = config[`MFAReq`].(float64)
		}

		for i := 0; float64(i) < leng; i++ {
			wg.Add(1)
			fmt.Fprintln(payload.Conns[e], payload.Payload[e])
			sendTime = append(sendTime, time.Now())
			go func(e int) {
				ea := make([]byte, 1000)
				payload.Conns[e].Read(ea)
				recv = append(recv, time.Now())
				statusCode = append(statusCode, string(ea[9:12]))

				if string(ea[9:12]) == `200` {
					sendInfo(string(ea[9:12]), dropTime)
				}

				wg.Done()
			}(e)
			time.Sleep(time.Duration(config["SpreadPerReq"].(float64)) * time.Microsecond)
		}
	}

	wg.Wait()

	sort.Slice(sendTime, func(i, j int) bool {
		return sendTime[i].Before(sendTime[j])
	})

	sort.Slice(recv, func(i, j int) bool {
		return recv[i].Before(recv[j])
	})

	for e, status := range statusCode {
		if status != "200" {
			content += fmt.Sprintf("- [DISMAL] Sent @ %v | [%v] @ %v\n", formatTime(sendTime[e]), status, formatTime(recv[e]))
			sendI(fmt.Sprintf("Sent @ %v | [%v] @ %v", formatTime(sendTime[e]), status, formatTime(recv[e])))
		} else {
			sendS(fmt.Sprintf("Sent @ %v | [%v] @ %v ~ %v", formatTime(sendTime[e]), status, formatTime(recv[e]), strings.Split(emailGot, ":")[0]))
			content += fmt.Sprintf("+ [DISMAL] Sent @ %v | [%v] @ %v ~ %v\n", formatTime(sendTime[e]), status, formatTime(recv[e]), strings.Split(emailGot, ":")[0])
		}
	}

	logSnipe(content, name)
}
