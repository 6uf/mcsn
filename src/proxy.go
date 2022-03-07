package src

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Liza-Developer/apiGO"
	"github.com/logrusorgru/aurora"
	"golang.org/x/net/proxy"
)

// Proxy code

func Proxy(name string, delay float64, dropTime int64) {
	var wg sync.WaitGroup
	var data SentRequests
	var content string

	searches := apiGO.Search(name)
	fmt.Print(aurora.Sprintf(aurora.Faint(aurora.White("Name: %v - Delay: %v - Searches: %v - Proxys: %v\n")), aurora.Red(name), aurora.Red(delay), aurora.Red(searches.Searches), aurora.Red(len(Pro))))

	for time.Now().Before(time.Unix(dropTime, 0).Add(-time.Second * 15)) {
		fmt.Print(aurora.Sprintf(aurora.Faint(aurora.White("Generating Proxy Connections In: %v      \r")), aurora.Red(time.Until(time.Unix(dropTime, 0).Add(-time.Second*15)).Round(time.Second).Seconds())))
		time.Sleep(time.Second * 1)
	}

	clients := genSockets(Pro, name)

	fmt.Println()
	fmt.Print(aurora.Sprintf(aurora.Faint(aurora.White("Sleeping until droptime: %v\n")), aurora.Red(time.Unix(dropTime, 0))))

	time.Sleep(time.Until(time.Unix(dropTime, 0).Add(time.Millisecond * time.Duration(0-delay)).Add(time.Duration(-float64(time.Since(time.Now()).Nanoseconds())/1000000.0) * time.Millisecond)))

	for _, config := range clients {
		wg.Add(1)
		go func(config Proxys) {
			var wgs sync.WaitGroup

			for _, Acc := range config.Accounts {
				if Acc.AccountType == "Giftcard" {
					for i := 0; i < Acc.Requests; i++ {
						wgs.Add(1)
						go func(Account apiGO.Info, payloads string) {
							SendTime, recvTime, Status := apiGO.Payload{}.SocketSending(config.Conn, payloads)

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
						}(Acc, fmt.Sprintf("POST /minecraft/profile HTTP/1.1\r\nHost: api.minecraftservices.com\r\nConnection: open\r\nContent-Length:%s\r\nContent-Type: application/json\r\nAccept: application/json\r\nAuthorization: Bearer %s\r\n\r\n"+string([]byte(`{"profileName":"`+name+`"}`))+"\r\n", strconv.Itoa(len(string([]byte(`{"profileName":"`+name+`"}`)))), Acc.Bearer))
					}
				} else {
					for i := 0; i < Acc.Requests; i++ {
						wgs.Add(1)
						go func(Account apiGO.Info, payloads string) {
							SendTime, recvTime, Status := apiGO.Payload{}.SocketSending(config.Conn, payloads)

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
						}(Acc, "PUT /minecraft/profile/name/"+name+" HTTP/1.1\r\nHost: api.minecraftservices.com\r\nConnection: open\r\nUser-Agent: MCSN/1.0\r\nAuthorization: bearer "+Acc.Bearer+"\r\n\r\n")
					}
				}
			}

			wgs.Wait()

			wg.Done()
		}(config)
	}

	wg.Wait()

	sort.Slice(data.Requests, func(i, j int) bool {
		return data.Requests[i].SentAt.Before(data.Requests[j].SentAt)
	})

	for _, request := range data.Requests {
		if request.Success {
			content += fmt.Sprintf("+ Sent @ %v >> [%v] @ %v ~ %v\n", formatTime(request.SentAt), request.StatusCode, formatTime(request.RecvAt), request.Email)
			fmt.Print(aurora.Sprintf(aurora.Faint(aurora.White("Sent @ %v >> [%v] @ %v ~ %v\n")), aurora.Green(formatTime(request.SentAt)), aurora.Green(request.StatusCode), aurora.Green(formatTime(request.RecvAt)), aurora.Green(request.Email)))

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

			request.check(name, searches.Searches, request.Type)

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

func GenProxys() []string {
	var Proxys []string

	f, _ := os.Open("Proxys.txt")
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		Proxys = append(Proxys, scanner.Text())
	}

	return Proxys
}

func randomInt(Proxys []string) string {
	for {
		rand.Seed(time.Now().UnixNano())
		proxy := Proxys[rand.Intn(len(Proxys))]
		if !used[proxy] {
			used[proxy] = true
			return proxy
		}
	}
}

func Setup(proxy []string) {
	for _, proxy := range proxy {
		used[proxy] = false
	}
}

func genSockets(Pro []string, name string) (pro []Proxys) {
	var Accs [][]apiGO.Info
	var incr int
	var use int
	roots := x509.NewCertPool()
	roots.AppendCertsFromPEM([]byte(rootCert))
	for _, Acc := range Bearers.Details {
		if len(Accs) == 0 {
			Accs = append(Accs, []apiGO.Info{
				Acc,
			})
		} else {
			if incr == 3 {
				incr = 0
				use++
				Accs = append(Accs, []apiGO.Info{})
			}
			Accs[use] = append(Accs[use], Acc)
		}
		incr++
	}
	for _, Accs := range Accs {
		var user, pass, ip, port string
		auth := strings.Split(randomInt(Pro), ":")
		ip, port = auth[0], auth[1]
		if len(auth) > 2 {
			user, pass = auth[2], auth[3]
		}
		req, err := proxy.SOCKS5("tcp", fmt.Sprintf("%v:%v", ip, port), &proxy.Auth{
			User:     user,
			Password: pass,
		}, proxy.Direct)
		if err != nil {
			fmt.Print(aurora.Sprintf(aurora.Faint(aurora.White("Couldnt login: %v - %v\n")), aurora.Red(ip), aurora.Red(err.Error())))
		} else {
			conn, err := req.Dial("tcp", "api.minecraftservices.com:443")
			if err != nil {
				fmt.Print(aurora.Sprintf(aurora.Faint(aurora.White("Couldnt login: %v - %v\n")), aurora.Red(ip), aurora.Red(err.Error())))
			} else {
				fmt.Print(aurora.Sprintf(aurora.Faint(aurora.White("logged into: %v\n")), aurora.Red(ip)))
				pro = append(pro, Proxys{
					Accounts: Accs,
					Conn:     tls.Client(conn, &tls.Config{RootCAs: roots, InsecureSkipVerify: true, ServerName: "api.minecraftservices.com"}),
				})
			}
		}
	}
	return pro
}
