package src

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/6uf/apiGO"
)

func formatTime(t time.Time) string {
	return t.Format("05.00000")
}

func formatTimeStamp(t time.Time) string {
	return t.Format("15:04:05")
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
			PrintGrad("Droptime [UNIX]: ")
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

		PrintGrad(fmt.Sprintf("Name: %v - Delay: %v - Droptime: %v - Searches: %v\n", Data.Name, Data.Delay, formatTimeStamp(time.Unix(Data.Droptime, 0)), apiGO.Search(Data.Name)))
		ReadReqs(Data.SnipeReq(Acc))
	case "auto":
		for {
			var Data []apiGO.Droptime
			if charType == "list" {
				file, _ := os.Open("names.txt")

				scanner := bufio.NewScanner(file)

				for scanner.Scan() {
					Data = append(Data,
						apiGO.Droptime{
							Name:     scanner.Text(),
							Droptime: int(apiGO.DropTime(scanner.Text())),
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
						PrintGrad("[ERROR] No more usable Account(s)\n")
						os.Exit(0)
					}
				}

				Data := apiGO.ReqConfig{
					Name:     name.Name,
					Delay:    delay,
					Droptime: int64(name.Droptime),
					Proxy:    false,
					Proxys:   Proxys,
					Bearers:  Bearers,
				}

				PrintGrad(fmt.Sprintf("Name: %v - Delay: %v - Droptime: %v - Searches: %v\n", Data.Name, Data.Delay, formatTimeStamp(time.Unix(Data.Droptime, 0)), apiGO.Search(Data.Name)))

				ReadReqs(Data.SnipeReq(Acc))
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
						if resp, err := apiGO.ChangeSkin(apiGO.JsonValue(SkinUrls{Url: Acc.ChangeSkinLink, Varient: "slim"}), status.Bearer); err == nil && resp.StatusCode == 200 {
							PrintGrad("Succesfully Changed your Skin!\n")
						}
					}
					PrintGrad(fmt.Sprintf("[%v] Succesfully Claimed %v\n", status.ResponseDetails.StatusCode, name))
					break
				}
			}

			time.Sleep(time.Duration(delay) * time.Second)
			fmt.Println()
		}
	case "proxy":
		if charType != "" {
			for {
				var Data []apiGO.Droptime

				if charType == "list" {
					file, _ := os.Open("names.txt")
					scanner := bufio.NewScanner(file)
					for scanner.Scan() {
						Data = append(Data,
							apiGO.Droptime{
								Name:     scanner.Text(),
								Droptime: int(apiGO.DropTime(scanner.Text())),
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
							PrintGrad("[ERROR] No more usable Account(s)\n")
							os.Exit(0)
						}
					}

					Data := apiGO.ReqConfig{
						Name:     name.Name,
						Delay:    delay,
						Droptime: int64(name.Droptime),
						Proxy:    false,
						Proxys:   Proxys,
						Bearers:  Bearers,
					}

					PrintGrad(fmt.Sprintf("Name: %v - Delay: %v - Droptime: %v - Searches: %v\n", Data.Name, Data.Delay, formatTimeStamp(time.Unix(Data.Droptime, 0)), apiGO.Search(Data.Name)))

					ReadReqs(Data.SnipeReq(Acc))
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

		PrintGrad(fmt.Sprintf("Name: %v - Delay: %v - Droptime: %v - Searches: %v\n", Data.Name, Data.Delay, formatTimeStamp(time.Unix(Data.Droptime, 0)), apiGO.Search(Data.Name)))

		ReadReqs(Data.SnipeReq(Acc))
	}
}

func ReadReqs(Data apiGO.SentRequests) {
	for _, request := range Data.Requests {
		switch request.ResponseDetails.StatusCode {
		case "200":
			if request.Type == "Giftcard" {
				request.Info = apiGO.GetProfileInformation(request.Bearer)
			}
			PrintGrad(fmt.Sprintf("%v >> [%v] @ %v O %v | %v\n", formatTime(request.ResponseDetails.SentAt), request.ResponseDetails.StatusCode, formatTime(request.ResponseDetails.RecvAt), request.Email, request.ClaimNameMC(Acc)))
			switch Acc.ChangeskinOnSnipe {
			case true:
				if resp, err := apiGO.ChangeSkin(apiGO.JsonValue(SkinUrls{Url: Acc.ChangeSkinLink, Varient: "slim"}), request.Bearer); err == nil && resp.StatusCode == 200 {
					PrintGrad("Succesfully Changed your Skin!\n")
				}
			}
			removeDetails(request)
		default:
			PrintGrad(fmt.Sprintf("%v >> [%v] @ %v X %v\n", formatTime(request.ResponseDetails.SentAt), request.ResponseDetails.StatusCode, formatTime(request.ResponseDetails.RecvAt), request.Email))
		}
	}
}

func CheckAccs() {
	for {
		time.Sleep(time.Second * 10)
		for _, Accs := range Acc.Bearers {
			if time.Now().Unix() > Accs.AuthedAt+Accs.AuthInterval {
				bearers := apiGO.Auth([]string{Accs.Email + ":" + Accs.Password})
				for point, data := range Acc.Bearers {
					for _, Accs := range bearers.Details {
						if Accs.Bearer != "" {
							if data.Email == Accs.Email {
								data.Bearer = Accs.Bearer
								data.NameChange = apiGO.CheckChange(Accs.Bearer)
								data.Type = Accs.AccountType
								data.Password = Accs.Password
								data.Email = Accs.Email
								data.AuthedAt = time.Now().Unix()
								data.Info = Accs.Info
								Acc.Bearers[point] = data
								UpdateBearer(Accs)
							}
						}
					}

					Acc.SaveConfig()
					Acc.LoadState()
					break // break the loop to update the info.
				}

				// if the Account isnt usable, remove it from the list
				var ts apiGO.Config
				for _, i := range Acc.Bearers {
					if i.Email != Accs.Email {
						ts.Bearers = append(ts.Bearers, i)
					}
				}

				Acc.Bearers = ts.Bearers

				Acc.SaveConfig()
				Acc.LoadState()
				break // break the loop to update the info.
			}
		}
	}
}

func UpdateBearer(B apiGO.Info) {
	for i, D := range Bearers.Details {
		if D == B {
			Bearers.Details[i] = B
			break
		}
	}
}

func GetNAMEMCKEY() error {
	var choose string
	PrintGrad("Use config bearer? [yes | no]: ")
	fmt.Scan(&choose)
	if strings.ContainsAny(strings.ToLower(choose), "yes ye y") {
		Acc.LoadState()
		if len(Acc.Bearers) == 0 {
			PrintGrad("Unable to continue, you have no Bearers added.\n")
			os.Exit(0)
		} else {
			var email string
			PrintGrad("Email of the Account you will use: ")
			fmt.Scan(&email)

			fmt.Println()
			for _, details := range Acc.Bearers {
				if strings.EqualFold(strings.ToLower(details.Email), strings.ToLower(email)) {
					if details.Bearer != "" {
						Bearers.Details = append(Bearers.Details, apiGO.Info{
							Bearer:      details.Bearer,
							AccountType: details.Type,
							Email:       details.Email,
						})
						break
					}
				}
			}
		}
	} else {
		var Accd string
		PrintGrad("Enter your Account details to continue [EMAIL:PASS]: ")
		fmt.Scan(&Accd)
		Bearers = apiGO.Auth([]string{Accd})
	}

	var data apiGO.Details = apiGO.Details{
		Bearer: Bearers.Details[0].Bearer,
		Info:   Bearers.Details[0].Info,
	}

	PrintGrad(data.ClaimNameMC(Acc))
	return nil
}
