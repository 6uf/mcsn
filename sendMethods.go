package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/Liza-Developer/api"
	"github.com/fatih/color"
)

func sendAuto(option string, delay float64) {
	leng := 0
	delays := config[`Spread`].(float64)

	if useAuto {
		delay = AutoOffset()
	}

	for i, name := range names {

		var sends []time.Time
		var recvs []time.Time
		var statuscodes []string
		var email []string

		if bearers.Bearers == nil || len(bearers.Bearers) == 0 {
			fmt.Println("Attempting to reauth accounts..")
			authAccs()
		}

		dropTime := drops[i]

		fmt.Printf("\n    Name: %v\n   Delay: %v\nDroptime: %v\n\n", name, delay, formatTime(time.Unix(dropTime, 0)))

		api.PreSleep(dropTime)

		payload := bearers.CreatePayloads(name)

		if useAuto {
			api.Sleep(dropTime, AutoOffset())
		} else {
			api.Sleep(dropTime, delay)
		}

		fmt.Println()

		for f, accType := range bearers.AccountType {
			switch accType {
			case "Giftcard":
				leng = 6
			case "Microsoft":
				leng = 2
			}

			for i := 0; i < leng; {
				go func() {
					send, recv, status := payload.SocketSending(int64(f))
					sends = append(sends, send)
					recvs = append(recvs, recv)
					statuscodes = append(statuscodes, status)
				}()
				i++
				time.Sleep(time.Duration(delays) * time.Microsecond)
			}
			email = append(email, strings.Split(AccountsVer[f], ":")[0])
		}

		time.Sleep(500 * time.Millisecond)

		sort.Slice(sends, func(i, j int) bool {
			return sends[i].Before(sends[j])
		})

		for i, accountType := range payload.AccountType {

			switch accountType {
			case "Giftcard":
				leng = 6
			case "Microsoft":
				leng = 2
			}

			for f := 0; f < leng; {
				if statuscodes[f] != "200" {
					content += fmt.Sprintf("- [%v] Sent @ %v | Recv @ %v\n", statuscodes[f], formatTime(sends[f]), formatTime(recvs[f]))
					fmt.Printf("[%v] Sent @ %v | Recv @ %v\n", statuscodes[f], formatTime(sends[f]), formatTime(recvs[f]))
				} else {
					time.Sleep(500 * time.Millisecond)

					bearerGot, emailGot, got := check(statuscodes[f], name, fmt.Sprintf("%v", dropTime))

					content += fmt.Sprintf("+ [%v] Succesfully sniped %v | %v\n", statuscodes[f], name, emailGot)
					color.Green(fmt.Sprintf("[%v] Recv @ %v | %v\n", statuscodes[f], formatTime(recvs[f]), strings.Split(emailGot, ":")[0]))

					if got {
						fmt.Println("[204] Sent Webhook")
					} else {
						fmt.Println("[FAILED] Couldnt send webhook.")
					}

					req, err := sendInfo.ChangeSkin(jsonValue(skinUrls{Url: sendInfo.SkinUrl, Varient: "slim"}), bearerGot)
					if err != nil {
						log.Println(err)
					} else {
						if req.StatusCode == 200 {
							fmt.Println("[200] Changed Skin")
						} else if req.StatusCode == 400 {
							fmt.Println("[400] Used wrong bearer during Skin Change. (acc is giftcard :skull:)")
						} else {
							fmt.Println("[401] Unauthorized")
						}
					}

					bearers.Bearers = remove(bearers.Bearers, bearers.Bearers[i])
					bearers.AccountType = remove(bearers.AccountType, bearers.AccountType[i])
					payload.Payload = remove(payload.Payload, payload.Payload[i])
				}
				f++
			}
		}
	}
}

func singlesniper(name string, delay float64) {
	var sends []time.Time
	var recvs []time.Time
	var statuscodes []string
	var email []string
	var leng int

	delays := config[`Spread`].(float64)
	dropTime = api.DropTime(name)

	if dropTime < 10000 {
		fmt.Print("[ERR] Unix Droptime: ")
		fmt.Scan(&dropTime)
		fmt.Println()
	}

	fmt.Printf("\n    Name: %v\n   Delay: %v\nDroptime: %v\n\n", name, delay, formatTime(time.Unix(dropTime, 0)))

	api.PreSleep(dropTime)

	payload := bearers.CreatePayloads(name)

	api.Sleep(dropTime, delay)

	fmt.Println()

	for f, accType := range bearers.AccountType {
		switch accType {
		case "Giftcard":
			leng = 6
		case "Microsoft":
			leng = 2
		}

		for i := 0; i < leng; {
			go func() {
				send, recv, status := payload.SocketSending(int64(f))
				sends = append(sends, send)
				recvs = append(recvs, recv)
				statuscodes = append(statuscodes, status)
			}()
			i++
			time.Sleep(time.Duration(delays) * time.Microsecond)
		}
		email = append(email, strings.Split(AccountsVer[f], ":")[0])
	}

	time.Sleep(500 * time.Millisecond)

	sort.Slice(sends, func(i, j int) bool {
		return sends[i].Before(sends[j])
	})

	var num int

	for _, accountType := range payload.AccountType {

		switch accountType {
		case "Giftcard":
			leng = 6
		case "Microsoft":
			leng = 2
		}

		for f := 0; f < leng; {
			if statuscodes[f] != "200" {
				content += fmt.Sprintf("- [%v] Sent @ %v | Recv @ %v\n", statuscodes[num], formatTime(sends[num]), formatTime(recvs[num]))
				fmt.Printf("[%v] Sent @ %v | Recv @ %v\n", statuscodes[num], formatTime(sends[num]), formatTime(recvs[num]))
				num++
			} else {
				time.Sleep(500 * time.Millisecond)

				bearerGot, emailGot, got := check(statuscodes[f], name, fmt.Sprintf("%v", dropTime))

				content += fmt.Sprintf("+ [%v] Succesfully sniped %v | %v\n", statuscodes[num], name, strings.Split(emailGot, ":")[0])
				color.Green(fmt.Sprintf("[%v] Recv @ %v | %v\n", statuscodes[f], formatTime(recvs[f]), strings.Split(emailGot, ":")[0]))

				if got {
					fmt.Println("[204] Sent Webhook")
				} else {
					fmt.Println("[FAILED] Couldnt send webhook.")
				}

				req, err := sendInfo.ChangeSkin(jsonValue(skinUrls{Url: sendInfo.SkinUrl, Varient: "slim"}), bearerGot)
				if err != nil {
					log.Println(err)
				} else {
					if req.StatusCode == 200 {
						fmt.Println("[200] Changed Skin")
					} else if req.StatusCode == 400 {
						fmt.Println("[400] Used wrong bearer during Skin Change. (acc is giftcard :skull:)")
					} else {
						fmt.Println("[401] Unauthorized")
					}
				}

				break
			}
			f++
		}
	}

	time.Sleep(time.Second * 3)
}

func searches(name string) float64 {
	var searches float64
	req, _ := http.NewRequest("GET", "https://droptime.site/api/v2/searches/"+name, nil)

	resp, _ := http.DefaultClient.Do(req)

	if resp.StatusCode != 200 {
		searches = 0
	} else {
		body, _ := ioutil.ReadAll(resp.Body)

		var config map[string]interface{}

		json.Unmarshal(body, &config)

		searches = config["searches"].(float64)
	}

	return searches
}
