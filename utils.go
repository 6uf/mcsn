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
	"math"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"runtime"
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

type Data struct {
	Name   string `json:"name"`
	Bearer string `json:"bearer"`
	Unix   string `json:"unix"`
	Config string `json:"config"`
}

type SentRequests struct {
	Requests []Details
}

type Details struct {
	Bearer     string
	SentAt     string
	RecvAt     string
	StatusCode string
	UnixRecv   int64
	Success    bool
	Email      string
}

type Config apiGO.Config

var acc Config

var (
	BearersVer  []string
	Confirmed   []string
	VpsesVer    []string
	bearers     apiGO.MCbearers
	AccountsVer []string
	config      map[string]interface{}

	images    []image.Image
	thirdRow  [][]int = [][]int{{64, 16, 72, 24}, {56, 16, 64, 24}, {48, 16, 56, 24}, {40, 16, 48, 24}, {32, 16, 40, 24}, {24, 16, 32, 24}, {16, 16, 24, 24}, {8, 16, 16, 24}, {0, 16, 8, 24}}
	secondRow [][]int = [][]int{{64, 8, 72, 16}, {56, 8, 64, 16}, {48, 8, 56, 16}, {40, 8, 48, 16}, {32, 8, 40, 16}, {24, 8, 32, 16}, {16, 8, 24, 16}, {8, 8, 16, 16}, {0, 8, 8, 16}}
	firstRow  [][]int = [][]int{{64, 0, 72, 8}, {56, 0, 64, 8}, {48, 0, 56, 8}, {40, 0, 48, 8}, {32, 0, 40, 8}, {24, 0, 32, 8}, {16, 0, 24, 8}, {8, 0, 16, 8}, {0, 0, 8, 8}}
)

func init() {

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

`), aurora.Bold(aurora.BrightBlack("4.25")), aurora.Bold(aurora.BrightBlack("Made By Liza"))))

	acc.LoadState()

	if acc.DiscordID == "" {
		var ID string
		sendW("Enter a Discord ID: ")
		fmt.Scan(&ID)

		acc.DiscordID = ID

		acc.SaveConfig()
		acc.LoadState()

		fmt.Println()
	}
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

type checkDetails struct {
	Error string `json:"error"`
	Sent  string `json:"sent"`
}

func (account Details) check(name, searches string) {
	var details checkDetails
	body, err := json.Marshal(Data{Name: name, Bearer: account.Bearer, Unix: fmt.Sprintf("%v", account.UnixRecv), Config: string(jsonValue(embeds{Content: "<@" + acc.DiscordID + ">", Embeds: []embed{{Description: fmt.Sprintf("Succesfully sniped %v with %v searches :bow_and_arrow:", name, searches), Color: 770000, Footer: footer{Text: "MCSN"}, Time: time.Now().Format(time.RFC3339)}}}))})
	if err == nil {
		req, _ := http.NewRequest("POST", "http://droptime.site/api/v2/webhook", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := http.DefaultClient.Do(req)
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal(body, &details)

		if details.Error != "" {
			sendE(details.Error)
		} else {
			if details.Sent != "" {
				sendS(details.Sent)
			} else {
				sendE(fmt.Sprintf("Couldnt send the request: %v", resp.StatusCode))
			}
		}
	}

	for i, accs := range acc.Bearers {
		if account.Email == accs.Email {
			acc.Bearers[i].NameChange = false
			acc.SaveConfig()
			acc.LoadState()

			var meow []apiGO.Info
			for _, acc := range bearers.Details {
				if acc.Email != accs.Email {
					meow = append(meow, acc)
				}
			}

			bearers.Details = meow
			break
		}
	}
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
				if strings.EqualFold(strings.ToLower(details.Email), strings.ToLower(email)) {
					if details.Bearer != "" {
						bearers.Details = append(bearers.Details, apiGO.Info{
							Bearer:      details.Bearer,
							AccountType: details.Type,
							Email:       details.Email,
						})
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
		bearers = apiGO.Auth([]string{accd})
	}

	if bearers.Details == nil {
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
	skin, _ := client.R().SetAuthToken(bearers.Details[0].Bearer).SetFormData(map[string]string{
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
			sendW("Droptime [UNIX] : ")
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
					if len(bearers.Details) == 0 {
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
					bearers.Details = append(bearers.Details, apiGO.Info{
						Bearer:      acc.Bearer,
						AccountType: acc.Type,
						Email:       acc.Email,
					})
				}
			}

			if bearers.Details == nil {
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
				bearers.Details = append(bearers.Details, apiGO.Info{
					Bearer:      bearer,
					AccountType: isGC(bearer),
				})
			}

			time.Sleep(time.Second)
		}
	} else {
		if acc.Bearers == nil {
			bearerz := apiGO.Auth(AccountsVer)
			if len(bearerz.Details) == 0 {
				sendE("Unable to authenticate your account(s), please Reverify your login details.\n")
				return
			} else {
				for _, accs := range bearerz.Details {
					acc.Bearers = append(acc.Bearers, apiGO.Bearers{
						Bearer:       accs.Bearer,
						AuthInterval: 86400,
						AuthedAt:     time.Now().Unix(),
						Type:         accs.AccountType,
						Email:        accs.Email,
						Password:     accs.Password,
						NameChange:   apiGO.CheckChange(accs.Bearer),
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

				bearerz := apiGO.Auth(auth)

				if len(bearerz.Details) != 0 {
					for _, accs := range bearerz.Details {
						acc.Bearers = append(acc.Bearers, apiGO.Bearers{
							Bearer:       accs.Bearer,
							AuthInterval: 86400,
							AuthedAt:     time.Now().Unix(),
							Type:         accs.AccountType,
							Email:        accs.Email,
							Password:     accs.Password,
							NameChange:   apiGO.CheckChange(accs.Bearer),
						})
					}

					acc.SaveConfig()
					acc.LoadState()
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
			sendI(fmt.Sprintf("Account %v turned up invalid. Attempting to Reauth", accs.Email))
			reAuth = append(reAuth, accs.Email+":"+accs.Password)
		}
	}

	if len(reAuth) != 0 {
		sendI(fmt.Sprintf("Reauthing %v accounts..", len(reAuth)))
		bearerz := apiGO.Auth(reAuth)

		if len(bearerz.Details) != 0 {
			for point, data := range acc.Bearers {
				for _, accs := range bearerz.Details {
					if data.Email == accs.Email {
						data.Bearer = accs.Bearer
						data.NameChange = apiGO.CheckChange(accs.Bearer)
						data.Type = accs.AccountType
						data.Password = accs.Password
						data.Email = accs.Email
						data.AuthedAt = time.Now().Unix()
						acc.Bearers[point] = data
						acc.SaveConfig()
					}
				}
			}
		}
	}

	acc.LoadState()
}

func checkVer(name string, delay float64, dropTime int64) {
	var content string
	var leng float64
	var data SentRequests

	searches := droptimeSiteSearches(name)

	sendI(fmt.Sprintf("Name: %v | Delay: %v | Searches: %v\n", name, delay, searches))

	var wg sync.WaitGroup

	apiGO.PreSleep(dropTime)

	payload := bearers.CreatePayloads(name)

	fmt.Println()

	apiGO.Sleep(dropTime, delay)

	fmt.Println()

	for e, account := range bearers.Details {
		switch account.AccountType {
		case "Giftcard":
			leng = float64(acc.GcReq)
		case "Microsoft":
			leng = float64(acc.MFAReq)
		}

		for i := 0; float64(i) < leng; i++ {
			wg.Add(1)
			go func(e int, account apiGO.Info) {
				fmt.Fprintln(payload.Conns[e], payload.Payload[e])
				sendTime := time.Now()
				ea := make([]byte, 1000)
				payload.Conns[e].Read(ea)
				recvTime := time.Now()

				data.Requests = append(data.Requests, Details{
					Bearer:     account.Bearer,
					SentAt:     formatTime(sendTime),
					RecvAt:     formatTime(recvTime),
					StatusCode: string(ea[9:12]),
					Success:    strings.Contains(string(ea[9:12]), "200"),
					UnixRecv:   recvTime.Unix(),
					Email:      account.Email,
				})

				wg.Done()
			}(e, account)
			time.Sleep(time.Duration(acc.SpreadPerReq) * time.Microsecond)
		}
	}

	wg.Wait()

	for _, request := range data.Requests {
		if request.Success {
			content += fmt.Sprintf("+ Sent @ %v | [%v] @ %v ~ %v\n", request.SentAt, request.StatusCode, request.RecvAt, request.Email)
			sendS(fmt.Sprintf("Sent @ %v | [%v] @ %v ~ %v\n", request.SentAt, request.StatusCode, request.RecvAt, request.Email))

			if acc.ChangeskinOnSnipe {
				sendInfo := apiGO.ServerInfo{
					SkinUrl: acc.ChangeSkinLink,
				}

				resp, _ := sendInfo.ChangeSkin(jsonValue(skinUrls{Url: sendInfo.SkinUrl, Varient: "slim"}), request.Bearer)
				if resp.StatusCode == 200 {
					sendS("Succesfully Changed your Skin!")
				} else {
					sendE("Couldnt Change your Skin..")
				}
			}

			request.check(name, searches)

			fmt.Println()

			sendI("If you enjoy using MCSN feel free to join the discord! https://discord.gg/a8EQ97ZfgK")
			break
		} else {
			content += fmt.Sprintf("- Sent @ %v | [%v] @ %v ~ %v\n", request.SentAt, request.StatusCode, request.RecvAt, request.Email)
			sendI(fmt.Sprintf("Sent @ %v | [%v] @ %v ~ %v", request.SentAt, request.StatusCode, request.RecvAt, request.Email))
		}
	}

	logSnipe(content, name)
}

// code from Alien https://github.com/wwhtrbbtt/AlienSniper

func (s *Config) ToJson() []byte {
	b, _ := json.MarshalIndent(s, "", "  ")
	return b
}

func (config *Config) SaveConfig() {
	apiGO.WriteFile("config.json", string(config.ToJson()))
}

func (s *Config) LoadState() {
	data, err := apiGO.ReadFile("config.json")
	if err != nil {
		sendI("No config file found, loading one.")
		s.LoadFromFile()
		s.GcReq = 2
		s.MFAReq = 2
		s.SpreadPerReq = 40
		s.ChangeskinOnSnipe = true
		s.ChangeSkinLink = "https://textures.minecraft.net/texture/516accb84322ca168a8cd06b4d8cc28e08b31cb0555eee01b64f9175cefe7b75"
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
		jsonFile, _ = os.Create("config.json")
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &c)
}

func checkAccs() {
	for {
		time.Sleep(time.Second * 10)

		// check if the last auth was more than a minute ago
		for _, accs := range acc.Bearers {
			if time.Now().Unix() > accs.AuthedAt+accs.AuthInterval {
				sendI(accs.Email + " is due for reauth")

				// authenticating account
				bearers := apiGO.Auth([]string{accs.Email + ":" + accs.Password})

				if bearers.Details != nil {
					for point, data := range acc.Bearers {
						for _, accs := range bearers.Details {
							if data.Email == accs.Email {
								data.Bearer = accs.Bearer
								data.NameChange = apiGO.CheckChange(accs.Bearer)
								data.Type = accs.AccountType
								data.Password = accs.Password
								data.Email = accs.Email
								data.AuthedAt = time.Now().Unix()
								acc.Bearers[point] = data
							}
						}
					}

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
				break // break the loop to update the info.
			}
		}
	}
}

func droptimeSiteSearches(username string) string {
	resp, err := http.Get(fmt.Sprintf("http://droptime.site/api/v2/searches/%v", username))

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
	for i := 1; i < 11; i++ {
		value := AutoOffset()
		sendI(fmt.Sprintf("%v. Request(s) gave %v as a estimated delay", i, math.Round(value)))
		values = append(values, value)
	}

	return mean(values), time.Since(time1)
}
