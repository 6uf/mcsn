package src

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"os"
	"os/signal"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/Liza-Developer/apiGO"
	"github.com/logrusorgru/aurora"
)

func Snipe(name string, delay float64, option string, charType string) {
	switch option {
	case "single":
		if name == "" {
			fmt.Print(aurora.Sprintf(aurora.Faint(aurora.White("[%v] You have entered a empty name | go run . snipe -u username -d 10 / mcsn.exe snipe -u username -d 10\n")), aurora.Red("ERROR")))
			return
		}

		dropTime := apiGO.DropTime(name)
		if dropTime < int64(10000) {
			fmt.Print(aurora.Faint(aurora.White("Droptime [UNIX]: ")))
			fmt.Scan(&dropTime)
			fmt.Println()
		}

		checkVer(name, delay, dropTime)

	case "auto":
		for {

			var names []string
			var drops []int64

			if charType == "list" {
				file, _ := os.Open("names.txt")

				scanner := bufio.NewScanner(file)

				for scanner.Scan() {
					drops = append(drops, apiGO.DropTime(scanner.Text()))
					names = append(names, scanner.Text())

					time.Sleep(1 * time.Second)
				}
			} else {
				names, drops = ThreeLetters(charType)
			}

			for e, name := range names {
				if delay == 0 {
					delay = float64(AutoOffset())
				}

				if !Acc.ManualBearer {
					if len(Bearers.Details) == 0 {
						fmt.Print(aurora.Sprintf(aurora.Faint(aurora.White("[%v] No more usable Account(s)\n")), aurora.Red("ERROR")))
						os.Exit(0)
					}
				}

				checkVer(name, delay, drops[e])

				fmt.Println()
			}

			if charType == "list" {
				break
			}
		}
	case "turbo":
		for {
			var leng float64
			var data SentRequests
			var wg sync.WaitGroup

			payload := Bearers.CreatePayloads(name)

			for e, Account := range Bearers.Details {
				switch Account.AccountType {
				case "Giftcard":
					leng = float64(Acc.GcReq)
				case "Microsoft":
					leng = float64(Acc.MFAReq)
				}

				for i := 0; float64(i) < leng; i++ {
					wg.Add(1)
					go func(e int, Account apiGO.Info) {
						fmt.Fprintln(payload.Conns[e], payload.Payload[e])
						SendTime := time.Now()
						ea := make([]byte, 1000)
						payload.Conns[e].Read(ea)
						recvTime := time.Now()

						data.Requests = append(data.Requests, Details{
							Bearer:     Account.Bearer,
							SentAt:     SendTime,
							RecvAt:     recvTime,
							StatusCode: string(ea[9:12]),
							Success:    string(ea[9:12]) == "200",
							UnixRecv:   recvTime.Unix(),
							Email:      Account.Email,
							Type:       Account.AccountType,
						})

						wg.Done()
					}(e, Account)
					time.Sleep(time.Duration(Acc.SpreadPerReq) * time.Microsecond)
				}
			}

			wg.Wait()

			for _, status := range data.Requests {
				if status.Success {
					status.check(name, "0", status.Type)

					if Acc.ChangeskinOnSnipe {
						SendInfo := apiGO.ServerInfo{
							SkinUrl: Acc.ChangeSkinLink,
						}

						SendInfo.ChangeSkin(jsonValue(skinUrls{Url: SendInfo.SkinUrl, Varient: "slim"}), status.Bearer)
					}

					fmt.Print(aurora.Sprintf(aurora.Faint(aurora.White("[%v] %v Claimed %v\n")), aurora.Green(status.StatusCode), aurora.Green("Succesfully"), aurora.Red(name)))

					break
				} else {
					fmt.Print(aurora.Sprintf(aurora.Faint(aurora.White("[%v] %v to claim %v\n")), aurora.Green(status.StatusCode), aurora.Red("Failed"), aurora.Red(name)))
				}
			}

			time.Sleep(time.Minute)

			fmt.Println()
		}
	}

	fmt.Println()

	fmt.Print((aurora.Faint(aurora.White("CTRL+C To Continue: "))))
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
}

func checkVer(name string, delay float64, dropTime int64) {
	var content string
	var data SentRequests

	searches, _ := apiGO.Search(name)

	fmt.Print(aurora.Sprintf(aurora.Faint(aurora.White("%v: %v - %v: %v - %v: %v\n")), aurora.Red("Name"), name, aurora.Red("Delay"), delay, aurora.Red("Searches"), searches))

	var wg sync.WaitGroup

	for time.Now().Before(time.Unix(dropTime, 0).Add(-time.Second * 5)) {
		fmt.Print(aurora.Sprintf(aurora.Faint(aurora.White("Generating Payloads/TLS Connection In: %v      \r")), aurora.Red(time.Until(time.Unix(dropTime, 0).Add(-time.Second*5)).Round(time.Second).Seconds())))
		time.Sleep(time.Second * 1)
	}

	payload := Bearers.CreatePayloads(name)
	conn, _ := tls.Dial("tcp", "api.minecraftservices.com:443", nil)

	fmt.Println()

	apiGO.Sleep(dropTime, delay)

	fmt.Println()

	for e, Account := range Bearers.Details {
		for i := 0; float64(i) < float64(Account.Requests); i++ {
			wg.Add(1)
			go func(e int, Account apiGO.Info) {
				fmt.Fprintln(conn, payload.Payload[e])
				SendTime := time.Now()

				var ea = make([]byte, 4096)
				conn.Read(ea)

				recvTime := time.Now()

				data.Requests = append(data.Requests, Details{
					Bearer:     Account.Bearer,
					SentAt:     SendTime,
					RecvAt:     recvTime,
					StatusCode: string(ea[9:12]),
					Success:    string(ea[9:12]) == "200",
					UnixRecv:   recvTime.Unix(),
					Email:      Account.Email,
					Type:       Account.AccountType,
					Cloudfront: strings.Contains(string(ea), "Error from cloudfront") && strings.Contains(string(ea), "We can't connect to the server for this app or website at this time. There might be too much traffic or a configuration error. Try again later, or contact the app or website owner."),
				})

				wg.Done()
			}(e, Account)
			time.Sleep(time.Duration(Acc.SpreadPerReq) * time.Microsecond)
		}
	}

	wg.Wait()

	sort.Slice(data.Requests, func(i, j int) bool {
		return data.Requests[i].SentAt.Before(data.Requests[j].SentAt)
	})

	for _, request := range data.Requests {
		if request.Success {
			content += fmt.Sprintf("+ Sent @ %v | [%v] @ %v ~ %v\n", formatTime(request.SentAt), request.StatusCode, formatTime(request.RecvAt), request.Email)
			fmt.Print(aurora.Sprintf(aurora.Faint(aurora.White("Sent @ %v >> [%v] @ %v ~ %v\n")), aurora.Red(formatTime(request.SentAt)), aurora.Green(request.StatusCode), aurora.Red(formatTime(request.RecvAt)), aurora.Red(request.Email)))

			if Acc.ChangeskinOnSnipe {
				SendInfo := apiGO.ServerInfo{
					SkinUrl: Acc.ChangeSkinLink,
				}

				resp, _ := SendInfo.ChangeSkin(jsonValue(skinUrls{Url: SendInfo.SkinUrl, Varient: "slim"}), request.Bearer)
				if resp.StatusCode == 200 {
					fmt.Print(aurora.Faint(aurora.White("Succesfully Changed your Skin!\n")))
				} else {
					fmt.Print(aurora.Sprintf(aurora.Faint(aurora.White("[%v] Couldnt Change your Skin..\n")), aurora.Red("ERROR")))
				}
			}

			request.check(name, searches, request.Type)

			fmt.Println()

			fmt.Print((aurora.Faint(aurora.White("If you enjoy using MCSN feel free to join the discord! https://discord.gg/a8EQ97ZfgK\n"))))
			break
		} else {
			content += fmt.Sprintf("- Sent @ %v >> [%v] @ %v ~ %v\n", formatTime(request.SentAt), request.StatusCode, formatTime(request.RecvAt), request.Email)
			if request.Cloudfront {
				fmt.Print(aurora.Sprintf(aurora.Faint(aurora.White("[%v] Sent @ %v >> [%v] @ %v ~ %v\n")), aurora.Red("CLOUDFRONT"), aurora.Red(formatTime(request.SentAt)), aurora.Red(request.StatusCode), aurora.Red(formatTime(request.RecvAt)), aurora.Red(request.Email)))
			} else {
				fmt.Print(aurora.Sprintf(aurora.Faint(aurora.White("Sent @ %v >> [%v] @ %v ~ %v\n")), aurora.Red(formatTime(request.SentAt)), aurora.Red(request.StatusCode), aurora.Red(formatTime(request.RecvAt)), aurora.Red(request.Email)))
			}
		}
	}

	logSnipe(content, name)
}
