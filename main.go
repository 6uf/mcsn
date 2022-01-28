package main

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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
					go checkAccs()
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
					sendI(fmt.Sprintf("Estimated (Mean) Delay: %v\n", MeanPing()))
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
					delay = float64(AutoOffset())
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
	file, _ := os.Open("accounts.txt")

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		AccountsVer = append(AccountsVer, scanner.Text())
	}

	if len(AccountsVer) == 0 {
		sendE("Unable to continue, you have no accounts added.\n")
		os.Exit(0)
	}

	grabDetails()

	if !acc.ManualBearer {
		if acc.Bearers == nil {
			sendE("No bearers have been found, please check your details.")
			os.Exit(0)
		} else {
			checkifValid()

			for _, acc := range acc.Bearers {
				bearers.Bearers = append(bearers.Bearers, acc.Bearer)
				bearers.AccountType = append(bearers.AccountType, acc.Type)
			}

			if bearers.Bearers == nil {
				sendE("Failed to authorize your bearers, please rerun the sniper.")
				os.Exit(0)
			}
		}
	}
}

func grabDetails() {
	if acc.ManualBearer {
		for _, bearer := range AccountsVer {
			if apiGO.CheckChange(bearer) {
				bearers.Bearers = append(bearers.Bearers, bearer)
				bearers.AccountType = append(bearers.AccountType, isGC(bearer))
			}

			time.Sleep(time.Second)
		}
	}

	if acc.Bearers == nil {
		bearerz, err := apiGO.Auth(AccountsVer)
		if err != nil {
			sendE(err.Error())
			os.Exit(0)
		}

		if len(bearerz.Bearers) == 0 {
			sendE("Unable to authenticate your account(s), please Reverify your login details.\n")
			return
		} else {
			for i := range bearerz.Bearers {
				acc.Bearers = append(acc.Bearers, Bearers{
					Bearer:       bearerz.Bearers[i],
					AuthInterval: int64(time.Hour * 24),
					AuthedAt:     time.Now().Unix(),
					Type:         bearerz.AccountType[i],
					Email:        strings.Split(AccountsVer[i], ":")[0],
					Password:     strings.Split(AccountsVer[i], ":")[1],
				})
			}
			acc.SaveConfig()
			acc.LoadState()
		}
	} else {
		if len(acc.Bearers) < len(AccountsVer) {
			var auth []string
			check := make(map[string]bool)

			for _, acc := range acc.Bearers {
				check[acc.Email+":"+acc.Password] = true
			}

			for _, accs := range AccountsVer {
				if !check[accs] {
					auth = append(auth, accs)
				}
			}

			bearerz, _ := apiGO.Auth(auth)

			if len(bearerz.Bearers) != 0 {
				for i := range bearerz.Bearers {
					acc.Bearers = append(acc.Bearers, Bearers{
						Bearer:       bearerz.Bearers[i],
						AuthInterval: int64(time.Hour * 24),
						AuthedAt:     time.Now().Unix(),
						Type:         bearerz.AccountType[i],
						Email:        strings.Split(AccountsVer[i], ":")[0],
						Password:     strings.Split(AccountsVer[i], ":")[1],
					})
					acc.SaveConfig()
					acc.LoadState()
				}
			}
		} else if len(AccountsVer) < len(acc.Bearers) {
			for _, accs := range AccountsVer {
				for _, num := range acc.Bearers {
					if accs == num.Email+":"+num.Password {
						acc.Bearers = append(acc.Bearers, num)
					}
				}
			}
			acc.SaveConfig()
			acc.LoadState()
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
	for _, accs := range acc.Bearers {
		f, _ := http.NewRequest("GET", "https://api.minecraftservices.com/minecraft/profile/name/boom/available", nil)
		f.Header.Set("Authorization", "Bearer "+accs.Bearer)
		j, _ := http.DefaultClient.Do(f)

		if j.StatusCode == 401 {
			sendI(fmt.Sprintf("Account %v turned up invalid. Attempting to Reauth\n", accs.Email))
			reAuth = append(reAuth, accs.Email+":"+accs.Password)
		}
	}

	if len(reAuth) != 0 {
		sendI(fmt.Sprintf("Reauthing %v accounts..\n", len(reAuth)))
		bearerz, _ := apiGO.Auth(reAuth)

		for i, accs := range bearerz.Bearers {
			if apiGO.CheckChange(accs) {
				acc.Bearers = append(acc.Bearers, Bearers{
					Bearer:       bearerz.Bearers[i],
					AuthInterval: int64(time.Hour * 24),
					AuthedAt:     time.Now().Unix(),
					Type:         bearerz.AccountType[i],
					Email:        strings.Split(reAuth[i], ":")[0],
					Password:     strings.Split(reAuth[i], ":")[1],
				})
			} else {
				sendI(fmt.Sprintf("Account %v cant name change\n", strings.Split(reAuth[i], ":")[0]))
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

	searches := droptimeSiteSearches(name)

	sendI(fmt.Sprintf("Name: %v | Delay: %v | Searches: %v\n", name, delay, searches))

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
			sendInfo(status, dropTime, searches)
			sendS(fmt.Sprintf("Sent @ %v | [%v] @ %v ~ %v", formatTime(sendTime[e]), status, formatTime(recv[e]), strings.Split(emailGot, ":")[0]))
			content += fmt.Sprintf("+ [DISMAL] Sent @ %v | [%v] @ %v ~ %v\n", formatTime(sendTime[e]), status, formatTime(recv[e]), strings.Split(emailGot, ":")[0])
		}
	}

	logSnipe(content, name)
}

// code from Alien https://github.com/wwhtrbbtt/AlienSniper

func ReadFile(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

func (s *Config) ToJson() []byte {
	b, _ := json.MarshalIndent(s, "", "  ")
	return b
}

func (config *Config) SaveConfig() {
	WriteFile("config.json", string(config.ToJson()))
}

func (s *Config) LoadState() {
	data, err := ReadFile("config.json")
	if err != nil {
		log.Println("No state file found, creating new one.")
		s.LoadFromFile()
		s.SaveConfig()
		return
	}

	json.Unmarshal([]byte(data), s)
	s.LoadFromFile()
}

func (c *Config) LoadFromFile() {
	// Load a config file

	jsonFile, err := os.Open("config.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		log.Fatalln("Failed to open config file: ", err)
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &c)
}

func WriteFile(path string, content string) {
	ioutil.WriteFile(path, []byte(content), 0644)
}

func checkAccs() {
	for {
		// check if the last auth was more than a minute ago
		for _, accs := range acc.Bearers {
			if time.Now().Unix() > accs.AuthedAt+accs.AuthInterval {
				sendI(accs.Email + " is due for auth")

				// authenticating account
				bearers, _ := apiGO.Auth([]string{accs.Email + ":" + accs.Password})

				if bearers.Bearers != nil {
					accs.AuthedAt = time.Now().Unix()
					accs.Bearer = bearers.Bearers[0]
					accs.Type = bearers.AccountType[0]
					acc.Bearers = append(acc.Bearers, accs)

					acc.SaveConfig()
					acc.LoadState()

					break // break the loop to update the info.
				}

				// if the account isnt usable, remove it from the list
				var ts Config
				for _, i := range acc.Bearers {
					if i.Email != accs.Email {
						ts.Bearers = append(ts.Bearers, i)
					}
				}

				acc.Bearers = ts.Bearers

				acc.SaveConfig()
				acc.LoadState()
				break // break the loop to update the state.Accounts info.
			}
		}

		time.Sleep(time.Second * 10)
	}
}

func droptimeSiteSearches(username string) string {
	resp, err := http.Get(fmt.Sprintf("https://droptime.site/api/v2/searches/%v", username))

	if err != nil {
		return "0"
	}
	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "0"
	}

	if resp.StatusCode < 300 {
		var res Searches
		err = json.Unmarshal(respBytes, &res)
		if err != nil {
			return "0"
		}

		return res.Searches
	}

	return "0"
}

//

func MeanPing() float64 {
	var values []float64
	for i := 1; i < 11; i++ {
		value := AutoOffset()
		sendI(fmt.Sprintf("%v`st Request(s) gave %v as a estimated delay", i, value))
		values = append(values, value)
	}

	total := 0.0

	for _, v := range values {
		total += v
	}

	return math.Round(total / float64(len(values)))

}
