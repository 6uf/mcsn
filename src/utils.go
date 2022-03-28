package src

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/6uf/apiGO"
	"github.com/logrusorgru/aurora/v3"
	"golang.org/x/net/proxy"
)

func ThreeLetters(option string) ([]string, []int64) {
	var threeL []string
	var names []string
	var droptime []int64
	var drop []int64

	if option == "list" {
		file, _ := os.Open("names.txt")

		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			if scanner.Text() == "" {
				break
			} else {
				threeL = append(threeL, scanner.Text())
				drop = append(drop, apiGO.DropTime(scanner.Text()))
			}
		}
	} else {
		isAlpha := regexp.MustCompile(`^[A-Za-z]+$`).MatchString

		grabName, _ := http.NewRequest("GET", "http://api.coolkidmacho.com/three", nil)
		jsonBody, _ := http.DefaultClient.Do(grabName)
		jsonGather, _ := ioutil.ReadAll(jsonBody.Body)

		var name []Name
		json.Unmarshal(jsonGather, &name)

		for i := range name {
			names = append(names, name[i].Names)
			droptime = append(droptime, int64(name[i].Drop))
		}

		switch option {
		case "3c":
			threeL = names
			drop = droptime
		case "3l":
			for i, username := range names {
				if !isAlpha(username) {
				} else {
					threeL = append(threeL, username)
					drop = append(drop, droptime[i])
				}
			}
		case "3n":
			for i, username := range names {
				if _, err := strconv.Atoi(username); err == nil {
					threeL = append(threeL, username)
					drop = append(drop, droptime[i])
				}
			}
		}
	}

	return threeL, drop
}

func jsonValue(f interface{}) []byte {
	g, _ := json.Marshal(f)
	return g
}

func formatTime(t time.Time) string {
	return t.Format("05.00000")
}

func removeDetails(Account Details) {
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
		Send:    Account.SentAt,
		Recv:    Account.RecvAt,
		Success: Account.Success,
	})

	Acc.SaveConfig()
	Acc.LoadState()
}

func isGC(bearer string) string {
	conn, _ := tls.Dial("tcp", "api.minecraftservices.com"+":443", nil)

	fmt.Fprintln(conn, "GET /minecraft/profile/namechange HTTP/1.1\r\nHost: api.minecraftservices.com\r\nUser-Agent: Dismal/1.0\r\nAuthorization: Bearer "+bearer+"\r\n\r\n")

	e := make([]byte, 12)
	conn.Read(e)

	switch string(e[9:12]) {
	case `404`:
		return "Giftcard"
	default:
		return "Microsoft"
	}
}

func CheckFiles() {
	_, err := os.Stat("logs")

	if os.IsNotExist(err) {
		os.Mkdir("logs", 0755)
	}

	_, err = os.Open("accounts.txt")
	if os.IsNotExist(err) {
		os.Create("accounts.txt")
	}

	_, err = os.Open("proxys.txt")
	if os.IsNotExist(err) {
		os.Create("proxys.txt")
	}

	_, err = os.Open("names.txt")
	if os.IsNotExist(err) {
		os.Create("names.txt")
	}

	_, err = os.Stat("cropped")
	if os.IsNotExist(err) {
		os.MkdirAll("cropped/logs", 0755)
	}
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

	var wg sync.WaitGroup
	for _, Accs := range Accs {
		wg.Add(1)
		go func(Accs []apiGO.Info) {
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

			wg.Done()
		}(Accs)
	}

	wg.Wait()
	return pro
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

		Data := ReqConfig{
			Name:     name,
			Delay:    delay,
			Droptime: dropTime,
			Proxy:    false,
		}

		Data.SnipeReq()
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

				Data := ReqConfig{
					Name:     name,
					Delay:    delay,
					Droptime: drops[e],
					Proxy:    false,
				}

				Data.SnipeReq()
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
					removeDetails(status)

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

			time.Sleep(time.Duration(delay) * time.Second)
			fmt.Println()
		}
	case "proxy":
		if charType != "" {
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

					Data := ReqConfig{
						Name:     name,
						Delay:    delay,
						Droptime: drops[e],
						Proxy:    true,
					}

					Data.SnipeReq()
					fmt.Println()
				}

				if charType == "list" {
					break
				}
			}
		}

		Data := ReqConfig{
			Name:     name,
			Delay:    delay,
			Droptime: apiGO.DropTime(name),
			Proxy:    true,
		}

		Data.SnipeReq()
	}
}
