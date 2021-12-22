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
	"time"

	"github.com/Liza-Developer/mcapi2"
)

type embeds struct {
	Content interface{} `json:"content"`
	Embeds  []embed     `json:"embeds"`
}

type embed struct {
	Description interface{} `json:"description"`
	Color       interface{} `json:"color"`
	Footer      footer      `json:"footer"`
	Time        interface{} `json:"timestamp"`
}

type footer struct {
	Text interface{} `json:"text"`
	Icon interface{} `json:"icon_url"`
}

type skinUrls struct {
	Url     interface{} `json:"url"`
	Varient interface{} `json:"variant"`
}

type Name struct {
	Names string `json:"name"`
}

func init() {
	q, _ := ioutil.ReadFile("accounts.json")

	config = mcapi2.GetConfig(q)

	sendInfo = mcapi2.ServerInfo{
		SkinUrl: config[`Config`].([]interface{})[1].(string),
	}

	content += `
+    __  ______________ _   __
-   /  |/  / ____/ ___// | / /
+  / /|_/ / /    \__ \/  |/ / 
- / /  / / /___ ___/ / /|  /  
+/_/  /_/\____//____/_/ |_/

`
	if config[`Config`].([]interface{})[3].(string) == "Discord ID [REQUIRED]" {
		var ID string
		fmt.Print("Enter your discord ID: \n>> ")
		fmt.Scanln(&ID)
		fmt.Println()

		config[`Config`].([]interface{})[3] = ID

		writetoFile(config)
	}

	flag.Parse()
}

func auto(option string, delay float64) {
	switch delay {
	case 0:
		useAuto = true
	}

	for {
		if option == "list" {
			for _, name := range config["Names"].([]interface{}) {
				names = append(names, name.(string))
			}
		} else {
			names = threeLetters(option)
		}

		sendAuto(option, delay)

		if len(names) == 0 && option == "list" {
			break
		}
	}
}

func jsonValue(f interface{}) []byte {
	g, _ := json.Marshal(f)
	return g
}

func remove(l []string, item string) []string {
	for i, other := range l {
		if other == item {
			l = append(l[:i], l[i+1:]...)
		}
	}
	return l
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
