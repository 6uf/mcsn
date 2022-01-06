package main

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Liza-Developer/api"
	"github.com/logrusorgru/aurora/v3"
)

func SortKey(item string) int {
	n, _ := strconv.Atoi(strings.Split(item, ":")[1])

	return n
}

func randomized() string {
	var returnS string

	if left {
		returnS = `/`
		left = false
	} else if !left {
		returnS = `\`
		left = true
	}

	return returnS
}

func sortInfo(statuscodes []string) []string {
	sort.Slice(statuscodes, func(i, j int) bool {
		return SortKey(statuscodes[i]) < SortKey(statuscodes[j])
	})

	list := make([][]string, len(statuscodes))

	for _, send := range statuscodes {
		item, _ := strconv.Atoi(strings.Split(send, ":")[1])
		list[item] = append(list[item], send)
	}

	for _, lists := range list {
		sort.Slice(lists, func(i, j int) bool {
			n1, _ := strconv.Atoi(strings.Split(lists[i], ":")[2])
			n2, _ := strconv.Atoi(strings.Split(lists[j], ":")[2])

			return n1 < n2
		})
	}

	var sendBack []string

	for _, items := range list {
		sendBack = append(sendBack, items...)
	}

	return sendBack
}

func sendAuto(option string, delay float64) {
	leng := 0
	delays := config[`Spread`].(float64)

	if useAuto {
		delay = AutoOffset()
	}

	for i, name := range names {
		var statuscodes []string

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

		var wg sync.WaitGroup

		for f, accType := range bearers.AccountType {
			switch accType {
			case "Giftcard":
				leng = 6
			case "Microsoft":
				leng = 2
			}

			for i := 0; i < leng; {
				wg.Add(1)
				go func(f, i int) {
					send, recv, status := payload.SocketSending(int64(f))
					statuscodes = append(statuscodes, status+fmt.Sprintf(":%v:%v:%v:%v", f, i, fmt.Sprintf("%v", formatTime(send)), fmt.Sprintf("%v", formatTime(recv))))
					wg.Done()
				}(f, i)
				i++
				time.Sleep(time.Duration(delays) * time.Microsecond)
			}
		}

		wg.Wait()

		list := sortInfo(statuscodes)

		for _, item := range list {
			items := strings.Split(item, ":")
			send := items[3]
			recv := items[4]
			status := items[0]

			if status != "200" {
				fmt.Print(aurora.Sprintf(aurora.Bold(aurora.White(("[%v] Sent @ %v | [%v] @ %v %v\n"))), aurora.Bold(aurora.Red("MCSN")), aurora.Bold(aurora.Red(send)), aurora.Bold(aurora.Red(status)), aurora.Bold(aurora.Red(recv)), randomized()))
			} else {
				time.Sleep(500 * time.Millisecond)

				bearerGot, emailGot, got := check(status, name, fmt.Sprintf("%v", dropTime))
				fmt.Print(aurora.Sprintf(aurora.Bold(aurora.White(("[%v] Sent @ %v | [%v] @ %v ~ %v %v\n"))), aurora.Bold(aurora.Green("MCSN")), aurora.Bold(aurora.Green(send)), aurora.Bold(aurora.Green(status)), aurora.Bold(aurora.Green(recv)), aurora.Bold(aurora.Green(strings.Split(emailGot, ":")[0])), randomized()))

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
			}
		}
	}
}

func singlesniper(name string, delay float64) {
	var statuscodes []string
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

	var wg sync.WaitGroup

	for f, accType := range bearers.AccountType {
		switch accType {
		case "Giftcard":
			leng = 6
		case "Microsoft":
			leng = 2
		}

		for i := 0; i < leng; {
			wg.Add(1)
			go func(f, i int) {
				send, recv, status := payload.SocketSending(int64(f))
				statuscodes = append(statuscodes, status+fmt.Sprintf(":%v:%v:%v:%v", f, i, fmt.Sprintf("%v", formatTime(send)), fmt.Sprintf("%v", formatTime(recv))))
				wg.Done()
			}(f, i)
			i++
			time.Sleep(time.Duration(delays) * time.Microsecond)
		}
	}

	wg.Wait()

	list := sortInfo(statuscodes)

	for _, item := range list {
		items := strings.Split(item, ":")
		send := items[3]
		recv := items[4]
		status := items[0]

		if status != "200" {
			fmt.Print(aurora.Sprintf(aurora.Bold(aurora.White(("[%v] Sent @ %v | [%v] @ %v %v\n"))), aurora.Bold(aurora.Red("MCSN")), aurora.Bold(aurora.Red(send)), aurora.Bold(aurora.Red(status)), aurora.Bold(aurora.Red(recv)), randomized()))
		} else {
			time.Sleep(500 * time.Millisecond)

			bearerGot, emailGot, got := check(status, name, fmt.Sprintf("%v", dropTime))
			fmt.Print(aurora.Sprintf(aurora.Bold(aurora.White(("[%v] Sent @ %v | [%v] @ %v ~ %v %v\n"))), aurora.Bold(aurora.Green("MCSN")), aurora.Bold(aurora.Green(send)), aurora.Bold(aurora.Green(status)), aurora.Bold(aurora.Green(recv)), aurora.Bold(aurora.Green(strings.Split(emailGot, ":")[0])), randomized()))

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
		}
	}

	time.Sleep(time.Second * 3)
}
