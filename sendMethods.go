package main

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/Liza-Developer/api"
	"github.com/logrusorgru/aurora/v3"
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
			fmt.Println(aurora.Bold(aurora.White("Attempting to reauth accounts..")))
			authAccs()
		}

		dropTime := drops[i]

		fmt.Print(aurora.Sprintf(aurora.Bold(aurora.White(("\n    Name: %v\n   Delay: %v\nDroptime: %v\n\n"))), aurora.Bold(aurora.Red(name)), aurora.Bold(aurora.Red(delay)), aurora.Bold(aurora.Red(time.Unix(dropTime, 0).Format("15:04:05")))))

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
					fmt.Print(aurora.Sprintf(aurora.Bold(aurora.White(("[%v] Sent @ %v | Rec @ %v\n"))), aurora.Bold(aurora.Red(statuscodes[f])), aurora.Bold(aurora.Red(formatTime(sends[f]))), aurora.Bold(aurora.Red(formatTime(recvs[f])))))
				} else {
					time.Sleep(500 * time.Millisecond)

					bearerGot, emailGot, got := check(statuscodes[f], name, fmt.Sprintf("%v", dropTime))

					content += fmt.Sprintf("+ [%v] Succesfully sniped %v | %v\n", statuscodes[f], name, emailGot)
					fmt.Print(aurora.Sprintf(aurora.Bold(aurora.White(("[%v] Rec @ %v | %v\n"))), aurora.Bold(aurora.Green(statuscodes[f])), aurora.Bold(aurora.Green(formatTime(recvs[f]))), aurora.Bold(aurora.Green(strings.Split(emailGot, ":")[0]))))

					if got {
						fmt.Println(aurora.Sprintf(aurora.Bold(aurora.White("[%v] Sent Webhook")), aurora.Bold(aurora.Green("204"))))
					} else {
						fmt.Println(aurora.Bold(aurora.White("[FAILED] Couldnt send webhook.")))
					}

					req, err := sendInfo.ChangeSkin(jsonValue(skinUrls{Url: sendInfo.SkinUrl, Varient: "slim"}), bearerGot)
					if err != nil {
						log.Println(err)
					} else {
						if req.StatusCode == 200 {
							fmt.Println(aurora.Sprintf(aurora.Bold(aurora.White("[%v] Changed Skin")), aurora.Bold(aurora.Green("200"))))
						} else if req.StatusCode == 400 {
							fmt.Println(aurora.Sprintf(aurora.Bold(aurora.White("[%v] Used wrong bearer during Skin Change. (acc is giftcard :skull:)")), aurora.Bold(aurora.Red("400"))))
						} else {
							fmt.Println(aurora.Sprintf(aurora.Bold(aurora.White("[%v] Unauthorized")), aurora.Bold(aurora.Green("401"))))
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
		fmt.Print(aurora.Bold(aurora.White("[ERR] Unix Droptime: ")))
		fmt.Scan(&dropTime)
	}

	fmt.Print(aurora.Sprintf(aurora.Bold(aurora.White(("\n    Name: %v\n   Delay: %v\nDroptime: %v\n\n"))), aurora.Bold(aurora.Red(name)), aurora.Bold(aurora.Red(delay)), aurora.Bold(aurora.Red(time.Unix(dropTime, 0).Format("15:04:05")))))

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
				content += fmt.Sprintf("- [%v] Sent @ %v | Rec @ %v\n", statuscodes[num], formatTime(sends[num]), formatTime(recvs[num]))
				fmt.Print(aurora.Sprintf(aurora.Bold(aurora.White(("[%v] Sent @ %v | Rec @ %v\n"))), aurora.Bold(aurora.Red(statuscodes[num])), aurora.Bold(aurora.Red(formatTime(sends[num]))), aurora.Bold(aurora.Red(formatTime(recvs[num])))))
				num++
			} else {
				time.Sleep(500 * time.Millisecond)

				bearerGot, emailGot, got := check(statuscodes[f], name, fmt.Sprintf("%v", dropTime))

				content += fmt.Sprintf("+ [%v] Succesfully sniped %v | %v\n", statuscodes[num], name, strings.Split(emailGot, ":")[0])
				fmt.Print(aurora.Sprintf(aurora.Bold(aurora.White(("[%v] Rec @ %v | %v\n"))), aurora.Bold(aurora.Green(statuscodes[f])), aurora.Bold(aurora.Green(formatTime(recvs[f]))), aurora.Bold(aurora.Green(strings.Split(emailGot, ":")[0]))))

				if got {
					fmt.Println(aurora.Sprintf(aurora.Bold(aurora.White("[%v] Sent Webhook")), aurora.Bold(aurora.Green("204"))))
				} else {
					fmt.Println(aurora.Bold(aurora.White("[FAILED] Couldnt send webhook.")))
				}

				req, err := sendInfo.ChangeSkin(jsonValue(skinUrls{Url: sendInfo.SkinUrl, Varient: "slim"}), bearerGot)
				if err != nil {
					log.Println(err)
				} else {
					if req.StatusCode == 200 {
						fmt.Println(aurora.Sprintf(aurora.Bold(aurora.White("[%v] Changed Skin")), aurora.Bold(aurora.Green("200"))))
					} else if req.StatusCode == 400 {
						fmt.Println(aurora.Sprintf(aurora.Bold(aurora.White("[%v] Used wrong bearer during Skin Change. (acc is giftcard :skull:)")), aurora.Bold(aurora.Red("400"))))
					} else {
						fmt.Println(aurora.Sprintf(aurora.Bold(aurora.White("[%v] Unauthorized")), aurora.Bold(aurora.Green("401"))))
					}
				}

				break
			}
			f++
		}
	}

	time.Sleep(time.Second * 3)
}
