package src

import (
	"crypto/tls"
	"fmt"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/6uf/apiGO"
	"github.com/logrusorgru/aurora/v3"
)

type ReqConfig struct {
	Name     string
	Delay    float64
	Droptime int64

	Proxy bool
}

func (Info *ReqConfig) SnipeReq() {
	var wg sync.WaitGroup
	var content string
	var data SentRequests

	fmt.Print(aurora.Sprintf(aurora.Faint(aurora.White("Name: %v - Delay: %v - Droptime: %v\n")), aurora.Red(Info.Name), aurora.Red(Info.Delay), aurora.Red(time.Unix(Info.Droptime, 0))))

	for time.Now().Before(time.Unix(Info.Droptime, 0).Add(-time.Second * 10)) {
		fmt.Print(aurora.Sprintf(aurora.Faint(aurora.White("Generating Payloads/TLS Connection In: %v      \r")), aurora.Red(time.Until(time.Unix(Info.Droptime, 0).Add(-time.Second*10)).Round(time.Second).Seconds())))
		time.Sleep(time.Second * 1)
	}

	if Info.Proxy {
		Clients := genSockets(Pro, Info.Name)
		time.Sleep(time.Until(time.Unix(Info.Droptime, 0).Add(time.Millisecond * time.Duration(0-Info.Delay)).Add(time.Duration(-float64(time.Since(time.Now()).Nanoseconds())/1000000.0) * time.Millisecond)))
		for _, config := range Clients {
			wg.Add(1)
			go func(config Proxys) {
				var wgs sync.WaitGroup
				for _, Acc := range config.Accounts {
					if Acc.AccountType == "Giftcard" {
						for i := 0; i < Acc.Requests; i++ {
							wgs.Add(1)
							go func(Account apiGO.Info, payloads string) {
								SendTime, recvTime, Status := apiGO.SocketSending(config.Conn, payloads)

								data.Requests = append(data.Requests, Details{
									Bearer:     Account.Bearer,
									SentAt:     SendTime,
									RecvAt:     recvTime,
									StatusCode: Status,
									Success:    Status == "200",
									UnixRecv:   recvTime.Unix(),
									Email:      Account.Email,
									Type:       Account.AccountType,
								})

								wgs.Done()
							}(Acc, fmt.Sprintf("POST /minecraft/profile HTTP/1.1\r\nHost: api.minecraftservices.com\r\nConnection: open\r\nContent-Length:%s\r\nContent-Type: application/json\r\nAccept: application/json\r\nAuthorization: Bearer %s\r\n\r\n"+string([]byte(`{"profileName":"`+Info.Name+`"}`))+"\r\n", strconv.Itoa(len(string([]byte(`{"profileName":"`+Info.Name+`"}`)))), Acc.Bearer))
						}
					} else {
						for i := 0; i < Acc.Requests; i++ {
							wgs.Add(1)
							go func(Account apiGO.Info, payloads string) {
								SendTime, recvTime, Status := apiGO.SocketSending(config.Conn, payloads)

								data.Requests = append(data.Requests, Details{
									Bearer:     Account.Bearer,
									SentAt:     SendTime,
									RecvAt:     recvTime,
									StatusCode: Status,
									Success:    Status == "200",
									UnixRecv:   recvTime.Unix(),
									Email:      Account.Email,
									Type:       Account.AccountType,
								})

								wgs.Done()
							}(Acc, "PUT /minecraft/profile/name/"+Info.Name+" HTTP/1.1\r\nHost: api.minecraftservices.com\r\nConnection: open\r\nUser-Agent: MCSN/1.0\r\nAuthorization: bearer "+Acc.Bearer+"\r\n\r\n")
						}
					}
				}

				wgs.Wait()

				wg.Done()
			}(config)
		}
	} else {
		payload := Bearers.CreatePayloads(Info.Name)
		conn, _ := tls.Dial("tcp", "api.minecraftservices.com:443", nil)
		time.Sleep(time.Until(time.Unix(Info.Droptime, 0).Add(time.Millisecond * time.Duration(0-Info.Delay)).Add(time.Duration(-float64(time.Since(time.Now()).Nanoseconds())/1000000.0) * time.Millisecond)))
		for e, Account := range Bearers.Details {
			for i := 0; float64(i) < float64(Account.Requests); i++ {
				wg.Add(1)
				go func(e int, Account apiGO.Info) {
					SendTime, recvTime, Status := apiGO.SocketSending(conn, payload.Payload[e])

					data.Requests = append(data.Requests, Details{
						Bearer:     Account.Bearer,
						SentAt:     SendTime,
						RecvAt:     recvTime,
						StatusCode: Status,
						Success:    Status == "200",
						UnixRecv:   recvTime.Unix(),
						Email:      Account.Email,
						Type:       Account.AccountType,
					})

					wg.Done()
				}(e, Account)
				time.Sleep(time.Duration(Acc.SpreadPerReq) * time.Microsecond)
			}
		}
	}

	wg.Wait()
	fmt.Println()

	sort.Slice(data.Requests, func(i, j int) bool {
		return data.Requests[i].SentAt.Before(data.Requests[j].SentAt)
	})

	for _, request := range data.Requests {
		if request.Success {
			content += fmt.Sprintf("+ Sent @ %v | [%v] @ %v ~ %v\n", formatTime(request.SentAt), request.StatusCode, formatTime(request.RecvAt), request.Email)
			fmt.Print(aurora.Sprintf(aurora.Faint(aurora.White("%v >> [%v] @ %v O %v\n")), aurora.Green(formatTime(request.SentAt)), aurora.Green(request.StatusCode), aurora.Green(formatTime(request.RecvAt)), aurora.Green(request.Email)))

			fmt.Println()

			if Acc.ChangeskinOnSnipe {
				SendInfo := apiGO.ServerInfo{
					SkinUrl: Acc.ChangeSkinLink,
				}

				resp, _ := SendInfo.ChangeSkin(jsonValue(skinUrls{Url: SendInfo.SkinUrl, Varient: "slim"}), request.Bearer)
				if resp.StatusCode == 200 {
					fmt.Print(aurora.Sprintf(aurora.Faint(aurora.White("[%v] Succesfully Changed your Skin!\n")), aurora.Green(resp.StatusCode)))
				} else {
					fmt.Print(aurora.Sprintf(aurora.Faint(aurora.White("[%v] Couldnt Change your Skin..\n")), aurora.Red("ERROR")))
				}
			}

			removeDetails(request)
			fmt.Print(aurora.Faint(aurora.White("\nIf you enjoy using MCSN feel free to join the discord! https://discord.gg/a8EQ97ZfgK\n")))
			break
		} else {
			content += fmt.Sprintf("- Sent @ %v >> [%v] @ %v ~ %v\n", formatTime(request.SentAt), request.StatusCode, formatTime(request.RecvAt), request.Email)
			if request.Cloudfront {
				fmt.Print(aurora.Sprintf(aurora.Faint(aurora.White("[%v] %v >> [%v] @ %v X %v\n")), aurora.Red("CLOUDFRONT"), aurora.Red(formatTime(request.SentAt)), aurora.Red(request.StatusCode), aurora.Red(formatTime(request.RecvAt)), aurora.Red(request.Email)))
			} else {
				fmt.Print(aurora.Sprintf(aurora.Faint(aurora.White("%v >> [%v] @ %v X %v\n")), aurora.Red(formatTime(request.SentAt)), aurora.Red(request.StatusCode), aurora.Red(formatTime(request.RecvAt)), aurora.Red(request.Email)))
			}
		}
	}

	logSnipe(content, Info.Name)
}
