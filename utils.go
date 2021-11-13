package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Liza-Developer/mcapi2"
)

type Name struct {
	Names string `json:"name"`
}

func init() {
	q, _ := ioutil.ReadFile("accounts.json")

	config = mcapi2.GetConfig(q)

	flag.Parse()

	sendInfo = mcapi2.ServerInfo{
		Webhook: config[`Config`].([]interface{})[2].(string),
		SkinUrl: config[`Config`].([]interface{})[1].(string),
	}
}

func auto(option string) {

	tick := 0
	leng := 0

	if option == "list" {
		var name []string
		for _, names := range config["Names"].([]interface{}) {
			name = append(name, names.(string))
		}

		for _, names := range name {
			if bearers.Bearers == nil || len(bearers.Bearers) == 0 {
				fmt.Println("Attempting to reauth accounts..")
				authAccs()
			}

			dropTime := mcapi2.DropTime(names)

			fmt.Printf("   Name: %v\n   Delay: %v\nDroptime: %v\n\n", names, AutoOffset(false), dropTime)

			mcapi2.PreSleep(dropTime)

			payload := bearers.CreatePayloads(names)

			mcapi2.Sleep(dropTime, AutoOffset(false))

			fmt.Println()

			for _, accType := range bearers.AccountType {
				switch accType {
				case "Giftcard":
					leng = 6
				case "Microsoft":
					leng = 2
				}

				for i := 0; i < leng; {
					go func() {
						send, recv, status := payload.SocketSending(10)
						if status == "200" {
							fmt.Printf("[%v] Succesfully sniped %v\n", status, names)
							f, _ := json.Marshal(config[`WebhookBody`])
							sendInfo.SendWebhook([]byte(strings.Replace(string(f), "{name}", names, 1)))
							sendInfo.ChangeSkin(nil, bearers.Bearers[tick])
							bearers.Bearers = remove(bearers.Bearers, bearers.Bearers[tick])
							bearers.AccountType = remove(bearers.AccountType, bearers.AccountType[tick])
						} else {
							fmt.Printf("[%v] Sent @ %v | Recv @ %v\n", status, formatTime(send), formatTime(recv))
						}
					}()
					i++
				}
				tick++
			}
		}
	} else {
		names := threeLetters(option)

		for _, name :=  range names {

			if bearers.Bearers == nil || len(bearers.Bearers) == 0 {
				fmt.Println("Attempting to reauth accounts..")
				authAccs()
			}

			dropTime := mcapi2.DropTime(name)

			fmt.Printf("   Name: %v\n   Delay: %v\nDroptime: %v\n\n", name, AutoOffset(false), dropTime)

			mcapi2.PreSleep(dropTime)

			payload := bearers.CreatePayloads(name)

			mcapi2.Sleep(dropTime, AutoOffset(false))

			fmt.Println()

			for _, accType := range bearers.AccountType {
				switch accType {
				case "Giftcard":
					leng = 6
				case "Microsoft":
					leng = 2
				}

				for i := 0; i < leng; {
					go func() {
						send, recv, status := payload.SocketSending(10)
						if status == "200" {
							fmt.Printf("[%v] Succesfully sniped %v\n", status, name)
							f, _ := json.Marshal(config[`WebhookBody`])
							sendInfo.SendWebhook([]byte(strings.Replace(string(f), "{name}", name, 1)))
							sendInfo.ChangeSkin(nil, bearers.Bearers[tick])
							bearers.Bearers = remove(bearers.Bearers, bearers.Bearers[tick])
							bearers.AccountType = remove(bearers.AccountType, bearers.AccountType[tick])
						} else {
							fmt.Printf("[%v] Sent @ %v | Recv @ %v\n", status, formatTime(send), formatTime(recv))
						}
					}()
					i++
				}
				tick++
			}
			if len(names) < 1 {
				names = threeLetters(option)
			}
		}
	}
}

func remove(l []string, item string) []string {
	for i, other := range l {
		if other == item {
			l = append(l[:i], l[i+1:]...)
		}
	}
	return l
}

func singlesniper(name string, delay float64) {

	spread, _ := strconv.Atoi(config[`Config`].([]interface{})[0].(string))

	dropTime = mcapi2.DropTime(name)

	fmt.Printf(`    Name: %v
   Delay: %v
  Spread: %v
Droptime: %v

`, name, delay, spread, formatTime(time.Unix(dropTime, 0)))

	mcapi2.PreSleep(dropTime)

	y := bearers.CreatePayloads(name)

	mcapi2.Sleep(dropTime, delay)

	fmt.Println()

	func() {
		var leng int
		for _, accountType := range y.AccountType {
			switch accountType {
			case "Giftcard":
				leng = 6
			case "Microsoft":
				leng = 2
			}

			for i := 0; i < leng; i++ {
				go func() {
					send, recv, status := y.SocketSending(int64(spread))
					if status == "200" {
						fmt.Printf("[%v] Recv @ %v | Got %v Succesfully.\n", status, formatTime(recv), name)
						f, _ := json.Marshal(config[`WebhookBody`])
						sendInfo.SendWebhook([]byte(strings.Replace(string(f), "{name}", name, 1)))
						sendInfo.ChangeSkin([]byte(""), bearers.Bearers[gotNum])
					} else {
						fmt.Printf("[%v] Sent @ %v | Recv @ %v\n", status, formatTime(send), formatTime(recv))
					}
				}()
			}
		}

		time.Sleep(time.Second)
	}()
}

func formatTime(t time.Time) string {
	return t.Format("15:04:05.00000")
}

func AutoOffset(print bool) float64 {

	payload := []byte("GET /minecraft/profile/name/test HTTP/1.1\r\nHost: api.minecraftservices.com\r\nAuthorization: Bearer TestToken" + "\r\n")
	conn, _ := tls.Dial("tcp", "api.minecraftservices.com"+":443", nil)
	pingTimes := make([]float64, 10)

	for i := 0; i < 10; i++ {
		junk := make([]byte, 1000)
		conn.Write(payload)
		time1 := time.Now()
		conn.Write([]byte("\r\n"))
		conn.Read(junk)
		time2 := time.Since(time1)
		switch print {
		case true:
			fmt.Printf("Took | %v\n", time2)
		}
		pingTimes[i] = float64(time2.Milliseconds())

	}

	// calculates the sum and does the math.. / 10000 to get the decimal version of sum then i * 5100~ (u can also do 5000) but it
	// only times the decimal to get the non deciaml number Example: 57 (the delay recommendations are very similar to python delay scripts ive tested)

	fmt.Println()

	return float64(sum(pingTimes)/10000) * 5000
}

func sum(array []float64) float64 {
	var sum1 float64 = 0
	for i := 0; i < 10; i++ {
		sum1 = sum1 + array[i]
	}
	return sum1
}

func threeLetters(option string) []string {
	var threeL []string
	isAlpha := regexp.MustCompile(`^[A-Za-z]+$`).MatchString

	names := make([]string, 0)
	grabName, _ := http.NewRequest("GET", "http://api.coolkidmacho.com/three", nil)
	jsonBody, _ := http.DefaultClient.Do(grabName)
	jsonGather, _ := ioutil.ReadAll(jsonBody.Body)

	var name []Name
	json.Unmarshal(jsonGather, &name)

	for i := range name {
		names = append(names, name[i].Names)
	}

	switch option {
	case "3c":
		threeL = names
	case "3l":
		for _, username := range names {
			if !isAlpha(username) {
			} else {
				threeL = append(threeL, username)
			}
		}
	case "3n":
		for _, username := range names {
			if _, err := strconv.Atoi(username); err == nil {
				threeL = append(threeL, username)
			}
		}
	}

	return threeL
}
