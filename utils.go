package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Liza-Developer/apiGO"
	"github.com/disintegration/imaging"
	"github.com/go-resty/resty/v2"
	"github.com/logrusorgru/aurora/v3"
	"github.com/nfnt/resize"
)

type embeds struct {
	Content interface{} `json:"content"`
	Embeds  []embed     `json:"embeds"`
}

type Searches struct {
	Searches string `json:"searches"`
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
	Names string  `json:"name"`
	Drop  float64 `json:"droptime"`
}

type Pixel struct {
	Point image.Point
	Color color.Color
}

type Bearers struct {
	Bearer       string `json:"Bearer"`
	Email        string `json:"Email"`
	Password     string `json:"Password"`
	AuthInterval int64  `json:"AuthInterval"`
	AuthedAt     int64  `json:"AuthedAt"`
	Type         string `json:"Type"`
	NameChange   bool   `json:"NameChange"`
}

type Vps struct {
	IP       string `json:"ip"`
	Port     string `json:"port"`
	Password string `json:"password"`
}

type Config struct {
	Bearers           []Bearers `json:"Bearers"`
	ChangeSkinLink    string    `json:"ChangeSkinLink"`
	ChangeskinOnSnipe bool      `json:"ChangeskinOnSnipe"`
	DiscordBotToken   string    `json:"DiscordBotToken"`
	DiscordID         string    `json:"DiscordID"`
	GcReq             int       `json:"GcReq"`
	MFAReq            int       `json:"MFAReq"`
	ManualBearer      bool      `json:"ManualBearer"`
	SpreadPerReq      int       `json:"SpreadPerReq"`
	Vps               []Vps     `json:"Vps"`
}

var acc Config

var (
	BearersVer  []string
	Confirmed   []string
	VpsesVer    []string
	bearers     apiGO.MCbearers
	AccountsVer []string
	name        string
	config      map[string]interface{}
	emailGot    string

	images    []image.Image
	thirdRow  [][]int = [][]int{{64, 16, 72, 24}, {56, 16, 64, 24}, {48, 16, 56, 24}, {40, 16, 48, 24}, {32, 16, 40, 24}, {24, 16, 32, 24}, {16, 16, 24, 24}, {8, 16, 16, 24}, {0, 16, 8, 24}}
	secondRow [][]int = [][]int{{64, 8, 72, 16}, {56, 8, 64, 16}, {48, 8, 56, 16}, {40, 8, 48, 16}, {32, 8, 40, 16}, {24, 8, 32, 16}, {16, 8, 24, 16}, {8, 8, 16, 16}, {0, 8, 8, 16}}
	firstRow  [][]int = [][]int{{64, 0, 72, 8}, {56, 0, 64, 8}, {48, 0, 56, 8}, {40, 0, 48, 8}, {32, 0, 40, 8}, {24, 0, 32, 8}, {16, 0, 24, 8}, {8, 0, 16, 8}, {0, 0, 8, 8}}
)

func init() {

	acc.LoadState()

	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}

	webhookvar, _ := ioutil.ReadFile("config.json")
	json.Unmarshal(webhookvar, &config)

	_, err := os.Stat("logs")

	if os.IsNotExist(err) {
		err = os.Mkdir("logs", 0755)
		if err != nil {
			fmt.Println("[MCSN] Failed to create Folder.")
		}
	}

	_, err = os.Stat("cropped")

	if os.IsNotExist(err) {
		os.MkdirAll("cropped/logs", 0755)
	}

	fmt.Print(aurora.White(`
    __  _____________  __
   /  |/  / ___/ __/ |/ /
  / /|_/ / /___\ \/    / 
 /_/  /_/\___/___/_/|_/
 `))

	fmt.Print(aurora.Sprintf(aurora.White(`
Ver: %v / %v

`), aurora.Bold(aurora.BrightBlack("4.10")), aurora.Bold(aurora.BrightBlack("Made By Liza"))))
}

func formatTime(t time.Time) string {
	return t.Format("05.00000")
}

func sendE(content string) {
	fmt.Println(aurora.Sprintf(aurora.White("[%v] "+content), aurora.Bold(aurora.Red("ERROR"))))
}

func sendI(content string) {
	fmt.Println(aurora.Sprintf(aurora.White("[%v] "+content), aurora.Yellow("INFO")))
}

func sendS(content string) {
	fmt.Println(aurora.Sprintf(aurora.White("[%v] "+content), aurora.Green("SUCCESS")))
}

func sendW(content string) {
	fmt.Print(aurora.Sprintf(aurora.White("[%v] "+content), aurora.Green("INPUT")))
}

func sendInfo(status string, dropTime int64, searches string) {
	bearerGot, emailGots, _, accs := check(status, name, fmt.Sprintf("%v", dropTime), searches)

	bearers.Bearers = remove(bearers.Bearers, bearerGot)
	bearers.AccountType = remove(bearers.AccountType, accs)

	emailGot = emailGots

	switch {
	case acc.ChangeskinOnSnipe:
		sendInfo := apiGO.ServerInfo{
			SkinUrl: acc.ChangeSkinLink,
		}

		sendInfo.ChangeSkin(jsonValue(skinUrls{Url: sendInfo.SkinUrl, Varient: "slim"}), bearerGot)
	}
}

// - Used to calculate delay, some of it is accurate some isnt! never rely on recommended delay.. simply base ur delay off it. -

func AutoOffset() float64 {
	var pingTimes int64
	conn, _ := tls.Dial("tcp", "api.minecraftservices.com:443", nil)
	defer conn.Close()
	for i := 0; i < 3; i++ {
		recv := make([]byte, 4096)
		time1 := time.Now()
		conn.Write([]byte("PUT /minecraft/profile/name/test HTTP/1.1\r\nHost: api.minecraftservices.com\r\nAuthorization: Bearer TestToken\r\n\r\n"))
		conn.Read(recv)
		pingTimes += time.Since(time1).Milliseconds()
	}
	return (float64(pingTimes) / float64(6000)) * 10000
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

func check(status, name, unixTime, searches string) (string, string, bool, string) {
	var bearerGot string
	var emailGot string
	var send bool
	var accountType string

	if status == `200` {
		for i, bearer := range bearers.Bearers {
			for _, email := range AccountsVer {
				httpReq, err := http.NewRequest("GET", "https://api.minecraftservices.com/minecraft/profile", nil)
				if err != nil {
					continue
				}
				httpReq.Header.Set("Authorization", "Bearer "+bearer)

				uwu, err := http.DefaultClient.Do(httpReq)
				if err != nil {
					continue
				}

				bodyByte, err := ioutil.ReadAll(uwu.Body)
				if err != nil {
					continue
				}

				var info map[string]interface{}
				json.Unmarshal(bodyByte, &info)

				if info[`name`] == nil {
				} else if info[`name`] == name {
					bearerGot = bearer

					if acc.ManualBearer {
						emailGot = email[0:30]
					} else {
						emailGot = email
					}

					accountType = bearers.AccountType[i]

					type data struct {
						Name   string `json:"name"`
						Bearer string `json:"bearer"`
						Unix   string `json:"unix"`
						Config string `json:"config"`
					}

					body, err := json.Marshal(data{Name: name, Bearer: bearerGot, Unix: unixTime, Config: string(jsonValue(embeds{Content: "<@" + acc.DiscordID + ">", Embeds: []embed{{Description: fmt.Sprintf("Succesfully sniped %v with %v searches:skull:", name, searches), Color: 770000, Footer: footer{Text: "MCSN"}, Time: time.Now().Format(time.RFC3339)}}}))})

					if err == nil {
						req, err := http.NewRequest("POST", "https://droptime.site/api/v2/webhook", bytes.NewBuffer(body))
						if err != nil {
							fmt.Println(err)
						}

						req.Header.Set("Content-Type", "application/json")

						resp, err := http.DefaultClient.Do(req)
						if err == nil {
							if resp.StatusCode == 200 {
								send = true
							} else {
								send = false
							}
						}
					}
				}
			}
		}
	}

	return bearerGot, emailGot, send, accountType
}

func threeLetters(option string) ([]string, []int64) {
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

func logSnipe(content string, name string) {
	logFile, err := os.OpenFile(fmt.Sprintf("logs/%v.txt", strings.ToLower(name)), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		sendE(err.Error())
	}

	defer logFile.Close()

	logFile.WriteString(content)
}

func skinart(name string) {
	var accd string
	var choose string
	sendW("Use config bearer? [yes | no]: ")
	fmt.Scan(&choose)

	if strings.ContainsAny(strings.ToLower(choose), "yes ye y") {
		acc.LoadState()

		if len(acc.Bearers) == 0 {
			sendE("Unable to continue, you have no bearers added.")
		} else {
			var email string
			sendW("Email of the account you will use: ")
			fmt.Scan(&email)

			for _, details := range acc.Bearers {
				if strings.EqualFold(details.Email, strings.ToLower(email)) {
					if details.Bearer != "" {
						bearers.Bearers = append(bearers.Bearers, details.Bearer)
						break
					} else {
						sendE("Your bearer is empty.")
						break
					}
				}
			}
		}
	} else {
		sendW("Enter your account details to continue [email:password]: ")
		fmt.Scan(&accd)
		bearers, _ = apiGO.Auth([]string{accd})
	}

	if bearers.Bearers == nil {
		sendE("Unable to continue, no bearers have been found.")
		return
	}

	img, err := readImage("images/image.png")
	if err != nil {
		sendE(err.Error())
	}

	base, err := readImage("images/base.png")
	if err != nil {
		sendE(err.Error())
	}

	if img.Bounds().Size() != image.Pt(72, 24) {
		img = resize.Resize(72, 24, img, resize.Lanczos3)

		writeImage(img, "images/image.png")
	}

	for _, array := range thirdRow {
		images = append(images, cropImage(img, image.Rect(array[0], array[1], array[2], array[3])))
	}

	for _, array := range secondRow {
		images = append(images, cropImage(img, image.Rect(array[0], array[1], array[2], array[3])))
	}

	for _, array := range firstRow {
		images = append(images, cropImage(img, image.Rect(array[0], array[1], array[2], array[3])))
	}

	for i, images := range images {
		imgs := imaging.Paste(base, images, image.Point{
			X: 8,
			Y: 8,
		})

		writeImage(imgs, fmt.Sprintf("cropped/logs/base_%v.png", i))
	}

	for i := 0; i < len(images); {
		changeSkin(i)
		i++
	}
}

func readImage(name string) (image.Image, error) {
	fd, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	img, _, err := image.Decode(fd)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func cropImage(img image.Image, crop image.Rectangle) image.Image {
	type subImager interface {
		SubImage(r image.Rectangle) image.Image
	}

	simg, _ := img.(subImager)

	return simg.SubImage(crop)
}

func writeImage(img image.Image, name string) error {
	fd, err := os.Create(name)
	if err != nil {
		return err
	}
	defer fd.Close()

	return png.Encode(fd, img)
}

func DecodePixelsFromImage(img image.Image, offsetX, offsetY int) []*Pixel {
	pixels := []*Pixel{}
	for y := 0; y <= img.Bounds().Max.Y; y++ {
		for x := 0; x <= img.Bounds().Max.X; x++ {
			p := &Pixel{
				Point: image.Point{x + offsetX, y + offsetY},
				Color: img.At(x, y),
			}
			pixels = append(pixels, p)
		}
	}
	return pixels
}

func changeSkin(num int) {
	client := resty.New()
	skin, _ := client.R().SetAuthToken(bearers.Bearers[0]).SetFormData(map[string]string{
		"variant": "slim",
	}).SetFile(fmt.Sprintf("cropped/logs/base_%v.png", num), fmt.Sprintf("cropped/logs/base_%v.png", num)).Post("https://api.minecraftservices.com/minecraft/profile/skins")

	if skin.StatusCode() == 200 {
		sendI("Skin Changed")
	} else {
		sendE("Failed skin change. (sleeping for 30 seconds)")
		time.Sleep(30 * time.Second)
	}

	sendW("Press CTRL+C to Continue : ")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	fmt.Println()
}

func snipe(name string, delay float64, option string, charType string) {
	var useAuto bool = false
	switch option {
	case "single":
		if name == "" {
			sendE("You have entered a empty name | go run . snipe -u username -d 10 / mcsn.exe snipe -u username -d 10")
			return
		}

		dropTime := apiGO.DropTime(name)
		if dropTime < int64(10000) {
			sendW("-!- Droptime [UNIX] : ")
			fmt.Scan(&dropTime)
			fmt.Println()
		}

		checkVer(name, delay, dropTime)

	case "auto":
		if delay == 0 {
			useAuto = true
		}

		for {

			names, drops := threeLetters(charType)

			for e, name := range names {
				if useAuto {
					delay = float64(AutoOffset())
				}

				if !acc.ManualBearer {
					if len(bearers.Bearers) == 0 {
						sendE("No more usable account(s)")
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
	}

	fmt.Println()

	sendW("Press CTRL+C to Continue : ")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
}

func authAccs() {
	file, _ := os.Open("accounts.txt")

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		AccountsVer = append(AccountsVer, scanner.Text())
	}

	if len(AccountsVer) == 0 {
		sendE("Unable to continue, you have no accounts added.\n")
		os.Exit(0)
	}

	grabDetails()

	if !acc.ManualBearer {
		if acc.Bearers == nil {
			sendE("No bearers have been found, please check your details.")
			os.Exit(0)
		} else {
			checkifValid()

			for _, acc := range acc.Bearers {
				if acc.NameChange {
					bearers.Bearers = append(bearers.Bearers, acc.Bearer)
					bearers.AccountType = append(bearers.AccountType, acc.Type)
				}
			}

			if bearers.Bearers == nil {
				sendE("Failed to authorize your bearers, please rerun the sniper.")
				os.Exit(0)
			}
		}
	}
}

func grabDetails() {
	if acc.ManualBearer {
		for _, bearer := range AccountsVer {
			if apiGO.CheckChange(bearer) {
				bearers.Bearers = append(bearers.Bearers, bearer)
				bearers.AccountType = append(bearers.AccountType, isGC(bearer))
			}

			time.Sleep(time.Second)
		}
	} else {
		if acc.Bearers == nil {
			bearerz, err := apiGO.Auth(AccountsVer)
			if err != nil {
				sendE(err.Error())
				os.Exit(0)
			}

			if len(bearerz.Bearers) == 0 {
				sendE("Unable to authenticate your account(s), please Reverify your login details.\n")
				return
			} else {
				for i := range bearerz.Bearers {
					acc.Bearers = append(acc.Bearers, Bearers{
						Bearer:       bearerz.Bearers[i],
						AuthInterval: 86400,
						AuthedAt:     time.Now().Unix(),
						Type:         bearerz.AccountType[i],
						Email:        strings.Split(AccountsVer[i], ":")[0],
						Password:     strings.Split(AccountsVer[i], ":")[1],
						NameChange:   apiGO.CheckChange(bearerz.Bearers[i]),
					})
				}
				acc.SaveConfig()
				acc.LoadState()
			}
		} else {
			if len(acc.Bearers) < len(AccountsVer) {
				var auth []string
				check := make(map[string]bool)

				for _, acc := range acc.Bearers {
					check[acc.Email+":"+acc.Password] = true
				}

				for _, accs := range AccountsVer {
					if !check[accs] {
						auth = append(auth, accs)
					}
				}

				bearerz, _ := apiGO.Auth(auth)

				if len(bearerz.Bearers) != 0 {
					for i := range bearerz.Bearers {
						acc.Bearers = append(acc.Bearers, Bearers{
							Bearer:       bearerz.Bearers[i],
							AuthInterval: 86400,
							AuthedAt:     time.Now().Unix(),
							Type:         bearerz.AccountType[i],
							Email:        strings.Split(AccountsVer[i], ":")[0],
							Password:     strings.Split(AccountsVer[i], ":")[1],
							NameChange:   apiGO.CheckChange(bearerz.Bearers[i]),
						})
						acc.SaveConfig()
						acc.LoadState()
					}
				}
			} else if len(AccountsVer) < len(acc.Bearers) {
				for _, accs := range AccountsVer {
					for _, num := range acc.Bearers {
						if accs == num.Email+":"+num.Password {
							acc.Bearers = append(acc.Bearers, num)
						}
					}
				}
				acc.SaveConfig()
				acc.LoadState()
			}
		}
	}
}

func isGC(bearer string) string {
	var accountT string
	conn, _ := tls.Dial("tcp", "api.minecraftservices.com"+":443", nil)

	fmt.Fprintln(conn, "GET /minecraft/profile/namechange HTTP/1.1\r\nHost: api.minecraftservices.com\r\nUser-Agent: Dismal/1.0\r\nAuthorization: Bearer "+bearer+"\r\n\r\n")

	e := make([]byte, 12)
	conn.Read(e)

	switch string(e[9:12]) {
	case `404`:
		accountT = "Giftcard"
	default:
		accountT = "Microsoft"
	}

	return accountT
}

func checkifValid() {
	var reAuth []string
	for _, accs := range acc.Bearers {
		f, _ := http.NewRequest("GET", "https://api.minecraftservices.com/minecraft/profile/name/boom/available", nil)
		f.Header.Set("Authorization", "Bearer "+accs.Bearer)
		j, _ := http.DefaultClient.Do(f)

		if j.StatusCode == 401 {
			sendI(fmt.Sprintf("Account %v turned up invalid. Attempting to Reauth\n", accs.Email))
			reAuth = append(reAuth, accs.Email+":"+accs.Password)
		}
	}

	if len(reAuth) != 0 {
		sendI(fmt.Sprintf("Reauthing %v accounts..\n", len(reAuth)))
		bearerz, _ := apiGO.Auth(reAuth)

		for i, accs := range bearerz.Bearers {
			acc.Bearers = append(acc.Bearers, Bearers{
				Bearer:       bearerz.Bearers[i],
				AuthInterval: int64(time.Hour * 24),
				AuthedAt:     time.Now().Unix(),
				Type:         bearerz.AccountType[i],
				Email:        strings.Split(reAuth[i], ":")[0],
				Password:     strings.Split(reAuth[i], ":")[1],
				NameChange:   apiGO.CheckChange(accs),
			})
		}
	}

	acc.SaveConfig()
	acc.LoadState()
}

func checkVer(name string, delay float64, dropTime int64) {
	var content string
	var sendTime []time.Time
	var leng float64
	var recv []time.Time
	var statusCode []string

	searches := droptimeSiteSearches(name)

	sendI(fmt.Sprintf("Name: %v | Delay: %v | Searches: %v\n", name, delay, searches))

	var wg sync.WaitGroup

	apiGO.PreSleep(dropTime)

	payload := bearers.CreatePayloads(name)

	fmt.Println()

	apiGO.Sleep(dropTime, delay)

	fmt.Println()

	for e, account := range payload.AccountType {
		switch account {
		case "Giftcard":
			leng = float64(acc.GcReq)
		case "Microsoft":
			leng = float64(acc.MFAReq)
		}

		for i := 0; float64(i) < leng; i++ {
			wg.Add(1)
			fmt.Fprintln(payload.Conns[e], payload.Payload[e])
			sendTime = append(sendTime, time.Now())
			go func(e int) {
				ea := make([]byte, 1000)
				payload.Conns[e].Read(ea)
				recv = append(recv, time.Now())
				statusCode = append(statusCode, string(ea[9:12]))
				wg.Done()
			}(e)
			time.Sleep(time.Duration(acc.SpreadPerReq) * time.Microsecond)
		}
	}

	wg.Wait()

	sort.Slice(sendTime, func(i, j int) bool {
		return sendTime[i].Before(sendTime[j])
	})

	sort.Slice(recv, func(i, j int) bool {
		return recv[i].Before(recv[j])
	})

	for e, status := range statusCode {
		if status != "200" {
			content += fmt.Sprintf("- [DISMAL] Sent @ %v | [%v] @ %v\n", formatTime(sendTime[e]), status, formatTime(recv[e]))
			sendI(fmt.Sprintf("Sent @ %v | [%v] @ %v", formatTime(sendTime[e]), status, formatTime(recv[e])))
		} else {
			sendInfo(status, dropTime, searches)
			sendS(fmt.Sprintf("Sent @ %v | [%v] @ %v ~ %v", formatTime(sendTime[e]), status, formatTime(recv[e]), strings.Split(emailGot, ":")[0]))
			content += fmt.Sprintf("+ [DISMAL] Sent @ %v | [%v] @ %v ~ %v\n", formatTime(sendTime[e]), status, formatTime(recv[e]), strings.Split(emailGot, ":")[0])
		}
	}

	logSnipe(content, name)
}

// code from Alien https://github.com/wwhtrbbtt/AlienSniper

func ReadFile(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

func (s *Config) ToJson() []byte {
	b, _ := json.MarshalIndent(s, "", "  ")
	return b
}

func (config *Config) SaveConfig() {
	WriteFile("config.json", string(config.ToJson()))
}

func (s *Config) LoadState() {
	data, err := ReadFile("config.json")
	if err != nil {
		log.Println("No state file found, creating new one.")
		s.LoadFromFile()
		s.SaveConfig()
		return
	}

	json.Unmarshal([]byte(data), s)
	s.LoadFromFile()
}

func (c *Config) LoadFromFile() {
	// Load a config file

	jsonFile, err := os.Open("config.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		log.Fatalln("Failed to open config file: ", err)
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &c)
}

func WriteFile(path string, content string) {
	ioutil.WriteFile(path, []byte(content), 0644)
}

func checkAccs() {
	for {
		// check if the last auth was more than a minute ago
		for _, accs := range acc.Bearers {
			if time.Now().Unix() > accs.AuthedAt+accs.AuthInterval {
				sendI(accs.Email + " is due for auth")

				// authenticating account
				bearers, _ := apiGO.Auth([]string{accs.Email + ":" + accs.Password})

				if bearers.Bearers != nil {
					accs.AuthedAt = time.Now().Unix()
					accs.Bearer = bearers.Bearers[0]
					accs.Type = bearers.AccountType[0]
					acc.Bearers = append(acc.Bearers, accs)

					acc.SaveConfig()
					acc.LoadState()

					break // break the loop to update the info.
				}

				// if the account isnt usable, remove it from the list
				var ts Config
				for _, i := range acc.Bearers {
					if i.Email != accs.Email {
						ts.Bearers = append(ts.Bearers, i)
					}
				}

				acc.Bearers = ts.Bearers

				acc.SaveConfig()
				acc.LoadState()
				break // break the loop to update the state.Accounts info.
			}
		}

		time.Sleep(time.Second * 10)
	}
}

func droptimeSiteSearches(username string) string {
	resp, err := http.Get(fmt.Sprintf("https://droptime.site/api/v2/searches/%v", username))

	if err != nil {
		return "0"
	}
	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "0"
	}

	if resp.StatusCode < 300 {
		var res Searches
		err = json.Unmarshal(respBytes, &res)
		if err != nil {
			return "0"
		}

		return res.Searches
	}

	return "0"
}

//

func mean(values []float64) float64 {
	total := 0.0

	for _, v := range values {
		total += float64(v)
	}

	return math.Round(total / float64(len(values)))
}

func MeanPing() (float64, time.Duration) {
	var values []float64
	time1 := time.Now()
	for i := 0; i < 10; i++ {
		value := AutoOffset()
		sendI(fmt.Sprintf("%v. Request(s) gave %v as a estimated delay", i, math.Round(value)))
		values = append(values, value)
	}

	return mean(values), time.Since(time1)
}

func run(acctype string) {
	for {
		time.Sleep(1 * time.Second)
		var delay float64 = AutoOffset()
		sendI(fmt.Sprintf("Testing %v in 3 seconds", delay))
		res := kqzzPing(acctype, 0.15, true, delay)
		if res == true {
			break
		}
	}
}

func kqzzPing(accs string, aim_for float64, log bool, delay float64) interface{} {
	var leng float64

	bearers := apiGO.MCbearers{}
	bearers.Bearers = []string{"testbearer"}
	bearers.AccountType = []string{accs}

	var recv []time.Time
	var wg sync.WaitGroup
	var statuscode []string

	dropTime := time.Now().Add(time.Second * 3).Unix()

	apiGO.PreSleep(dropTime)

	payload := bearers.CreatePayloads("Test")

	fmt.Println()

	apiGO.Sleep(dropTime, delay)

	fmt.Println()

	for e, account := range payload.AccountType {
		switch account {
		case "Giftcard":
			leng = float64(acc.GcReq)
		case "Microsoft":
			leng = float64(acc.MFAReq)
		}

		for i := 0; float64(i) < leng; i++ {
			wg.Add(1)
			fmt.Fprintln(payload.Conns[e], payload.Payload[e])
			go func(e int) {
				ea := make([]byte, 1000)
				payload.Conns[e].Read(ea)
				recv = append(recv, time.Now())
				statuscode = append(statuscode, string(ea[9:12]))
				wg.Done()
			}(e)
			time.Sleep(time.Duration(acc.SpreadPerReq) * time.Microsecond)
		}
	}

	wg.Wait()

	for i, sends := range recv {
		in, _ := strconv.Atoi(fmt.Sprintf("%v", sends.Format(".000")[1:]))

		sendI(fmt.Sprintf("Recv @: %v | [%v]", formatTime(sends), statuscode[i]))

		if InBetween(in, 99, 105) {
			sendS(fmt.Sprintf("%v is a good delay!", delay))
			return true
		}
	}

	return false
}

func InBetween(i, min, max int) bool {
	if (i >= min) && (i <= max) {
		return true
	} else {
		return false
	}
}
