package src

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Liza-Developer/apiGO"
	"github.com/logrusorgru/aurora"
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

func (Account Details) check(name, searches, AccType string) {
	var details checkDetails
	body, _ := json.Marshal(Data{Name: name, Bearer: Account.Bearer, Id: Acc.DiscordID, Unix: Account.UnixRecv, Config: string(jsonValue(embeds{Content: "<@" + Acc.DiscordID + ">", Embeds: []embed{{Description: fmt.Sprintf("[%v] Succesfully sniped %v with %v searches :bow_and_arrow:", AccType, name, searches), Color: 770000, Footer: footer{Text: "MCSN"}, Time: time.Now().Format(time.RFC3339)}}}))})

	req, _ := http.NewRequest("POST", "https://droptime.herokuapp.com/webhook", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := http.DefaultClient.Do(req)
	body, _ = ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &details)

	if details.Error != "" {
		fmt.Print(aurora.Sprintf(aurora.Faint(aurora.White("[%v] %v\n")), aurora.Red("ERROR"), details.Error))
	} else if details.Sent != "" {
		fmt.Print(aurora.Sprintf(aurora.Faint(aurora.White("[%v] %v\n")), aurora.Green("200"), details.Sent))
	}

	removeDetails(Account)
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
