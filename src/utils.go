package src

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/6uf/apiGO"
	"github.com/logrusorgru/aurora/v3"
)

func formatTime(t time.Time) string {
	return t.Format("05.00000")
}

func removeDetails(Account apiGO.Details) {
	var new []apiGO.Bearers
	for _, Accs := range Acc.Bearers {
		if Account.Email != Accs.Email {
			new = append(new, Accs)
		}
	}

	Acc.Bearers = new

	var meow []apiGO.Info
	for _, Accs := range Acc.Bearers {
		for _, Acc := range Bearers.Details {
			if Acc.Email != Accs.Email {
				meow = append(meow, Acc)
			}
		}
	}

	Bearers.Details = meow

	var Accz []string
	file, _ := os.Open("accounts.txt")
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Split(scanner.Text(), ":")[0] != Account.Email {
			Accz = append(Accz, scanner.Text())
		}
	}

	rewrite("accounts.txt", strings.Join(Accz, "\n"))

	Acc.Logs = append(Acc.Logs, apiGO.Logs{
		Email:   Account.Email,
		Send:    Account.ResponseDetails.SentAt,
		Recv:    Account.ResponseDetails.RecvAt,
		Success: Account.ResponseDetails.StatusCode == "200",
	})

	Acc.SaveConfig()
	Acc.LoadState()
}

func Snipe(name string, delay float64, option string, charType string) {
	switch option {
	case "single":
		dropTime := apiGO.DropTime(name)
		if dropTime < int64(10000) {
			fmt.Print(aurora.Faint(aurora.White("Droptime [UNIX]: ")))
			fmt.Scan(&dropTime)
			fmt.Print("\n")
		}

		Data := apiGO.ReqConfig{
			Name:     name,
			Delay:    delay,
			Droptime: dropTime,
			Proxy:    false,
			Bearers:  Bearers,
			Proxys:   Proxys,
		}

		fmt.Print(aurora.Sprintf(aurora.Faint(aurora.White("Name: %v - Delay: %v - Droptime: %v\n")), aurora.Red(Data.Name), aurora.Red(Data.Delay), aurora.Red(time.Unix(Data.Droptime, 0))))

		ReadReqs(Data.SnipeReq())
	case "auto":
		for {
			var Data []apiGO.Names
			if charType == "list" {
				file, _ := os.Open("names.txt")

				scanner := bufio.NewScanner(file)

				for scanner.Scan() {
					Data = append(Data,
						apiGO.Names{
							Name:     scanner.Text(),
							Droptime: apiGO.DropTime(scanner.Text()),
							Search:   "0",
						},
					)
				}
			} else {
				Data = apiGO.ThreeLetters(charType)
			}

			for _, name := range Data {
				if delay == 0 {
					delay = apiGO.PingMC()
				}

				if !Acc.ManualBearer {
					if len(Bearers.Details) == 0 {
						fmt.Print(aurora.Sprintf(aurora.Faint(aurora.White("[%v] No more usable Account(s)\n")), aurora.Red("ERROR")))
						os.Exit(0)
					}
				}

				Data := apiGO.ReqConfig{
					Name:     name.Name,
					Delay:    delay,
					Droptime: name.Droptime,
					Proxy:    false,
					Proxys:   Proxys,
					Bearers:  Bearers,
				}

				fmt.Print(aurora.Sprintf(aurora.Faint(aurora.White("Name: %v - Delay: %v - Droptime: %v\n")), aurora.Red(Data.Name), aurora.Red(Data.Delay), aurora.Red(time.Unix(Data.Droptime, 0))))

				ReadReqs(Data.SnipeReq())
				fmt.Println()
			}

			if charType == "list" {
				break
			}
		}
	case "turbo":
		for {
			var leng float64
			var data apiGO.SentRequests
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

						data.Requests = append(data.Requests, apiGO.Details{
							ResponseDetails: apiGO.Resp{
								SentAt:     SendTime,
								RecvAt:     recvTime,
								StatusCode: string(ea[9:12]),
							},
							Bearer: Account.Bearer,
							Email:  Account.Email,
							Type:   Account.AccountType,
						})

						wg.Done()
					}(e, Account)
					time.Sleep(time.Duration(Acc.SpreadPerReq) * time.Microsecond)
				}
			}

			wg.Wait()

			for _, status := range data.Requests {
				if status.ResponseDetails.StatusCode == "200" {
					removeDetails(status)

					if Acc.ChangeskinOnSnipe {
						SendInfo := apiGO.ServerInfo{
							SkinUrl: Acc.ChangeSkinLink,
						}

						SendInfo.ChangeSkin(apiGO.JsonValue(SkinUrls{Url: SendInfo.SkinUrl, Varient: "slim"}), status.Bearer)
					}

					fmt.Print(aurora.Sprintf(aurora.Faint(aurora.White("[%v] %v Claimed %v\n")), aurora.Green(status.ResponseDetails.StatusCode), aurora.Green("Succesfully"), aurora.Red(name)))

					break
				} else {
					fmt.Print(aurora.Sprintf(aurora.Faint(aurora.White("[%v] %v to claim %v\n")), aurora.Green(status.ResponseDetails.StatusCode), aurora.Red("Failed"), aurora.Red(name)))
				}
			}

			time.Sleep(time.Duration(delay) * time.Second)
			fmt.Println()
		}
	case "proxy":
		if charType != "" {
			for {
				var Data []apiGO.Names

				if charType == "list" {
					file, _ := os.Open("names.txt")
					scanner := bufio.NewScanner(file)
					for scanner.Scan() {
						Data = append(Data,
							apiGO.Names{
								Name:     scanner.Text(),
								Droptime: apiGO.DropTime(scanner.Text()),
								Search:   "0",
							},
						)
					}
				} else {
					Data = apiGO.ThreeLetters(charType)
				}

				for _, name := range Data {
					if delay == 0 {
						delay = apiGO.PingMC()
					}

					if !Acc.ManualBearer {
						if len(Bearers.Details) == 0 {
							fmt.Print(aurora.Sprintf(aurora.Faint(aurora.White("[%v] No more usable Account(s)\n")), aurora.Red("ERROR")))
							os.Exit(0)
						}
					}

					Data := apiGO.ReqConfig{
						Name:     name.Name,
						Delay:    delay,
						Droptime: name.Droptime,
						Proxy:    false,
						Proxys:   Proxys,
						Bearers:  Bearers,
					}

					fmt.Print(aurora.Sprintf(aurora.Faint(aurora.White("Name: %v - Delay: %v - Droptime: %v\n")), aurora.Red(Data.Name), aurora.Red(Data.Delay), aurora.Red(time.Unix(Data.Droptime, 0))))

					ReadReqs(Data.SnipeReq())
					fmt.Println()
				}

				if charType == "list" {
					break
				}
			}
		}

		Data := apiGO.ReqConfig{
			Name:     name,
			Delay:    delay,
			Droptime: apiGO.DropTime(name),
			Proxy:    true,
			Proxys:   Proxys,
			Bearers:  Bearers,
		}

		fmt.Print(aurora.Sprintf(aurora.Faint(aurora.White("Name: %v - Delay: %v - Droptime: %v\n")), aurora.Red(Data.Name), aurora.Red(Data.Delay), aurora.Red(time.Unix(Data.Droptime, 0))))

		ReadReqs(Data.SnipeReq())
	}
}

func ReadReqs(Data apiGO.SentRequests) {
	for _, request := range Data.Requests {
		switch request.ResponseDetails.StatusCode {
		case "200":
			fmt.Print(aurora.Sprintf(aurora.Faint(aurora.White("%v >> [%v] @ %v O %v\n\n")), aurora.Green(formatTime(request.ResponseDetails.SentAt)), aurora.Green(request.ResponseDetails.StatusCode), aurora.Green(formatTime(request.ResponseDetails.RecvAt)), aurora.Green(request.Email)))
			switch Acc.ChangeskinOnSnipe {
			case true:
				SendInfo := apiGO.ServerInfo{
					SkinUrl: Acc.ChangeSkinLink,
				}
				resp, err := SendInfo.ChangeSkin(apiGO.JsonValue(SkinUrls{Url: SendInfo.SkinUrl, Varient: "slim"}), request.Bearer)
				if err == nil {
					if resp.StatusCode == 200 {
						fmt.Print(aurora.Sprintf(aurora.Faint(aurora.White("[%v] Succesfully Changed your Skin!\n")), aurora.Green(resp.StatusCode)))
					} else {
						fmt.Print(aurora.Sprintf(aurora.Faint(aurora.White("[%v] Couldnt Change your Skin..\n")), aurora.Red("ERROR")))
					}
				}
			}
			removeDetails(request)
		default:
			fmt.Print(aurora.Sprintf(aurora.Faint(aurora.White("%v >> [%v] @ %v X %v\n")), aurora.Red(formatTime(request.ResponseDetails.SentAt)), aurora.Red(request.ResponseDetails.StatusCode), aurora.Red(formatTime(request.ResponseDetails.RecvAt)), aurora.Red(request.Email)))
		}
	}
}
