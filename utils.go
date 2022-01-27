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
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"runtime"
	"strconv"
	"strings"
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

func sendInfo(status string, dropTime int64) {
	bearerGot, emailGots, _, acc := check(status, name, fmt.Sprintf("%v", dropTime))

	bearers.Bearers = remove(bearers.Bearers, bearerGot)
	bearers.AccountType = remove(bearers.AccountType, acc)

	emailGot = emailGots

	switch {
	case config[`ChangeskinOnSnipe`] == true:
		sendInfo := apiGO.ServerInfo{
			SkinUrl: config[`ChangeSkinLink`].(string),
		}

		sendInfo.ChangeSkin(jsonValue(skinUrls{Url: sendInfo.SkinUrl, Varient: "slim"}), bearerGot)
	}
}

// - Used to calculate delay, some of it is accurate some isnt! never rely on recommended delay.. simply base ur delay off it. -

func AutoOffset() float64 {
	var pingTimes []float64
	conn, _ := tls.Dial("tcp", "api.minecraftservices.com:443", nil)

	for i := 0; i < 10; i++ {
		junk := make([]byte, 4069)
		time1 := time.Now()
		conn.Write([]byte("GET /minecraft/profile/name/test HTTP/1.1\r\nHost: api.minecraftservices.com\r\nAuthorization: Bearer TestToken\r\n\r\n"))
		conn.Read(junk)
		time2 := time.Since(time1)
		pingTimes = append(pingTimes, float64(time2.Milliseconds()))
	}

	return float64(apiGO.Sum(pingTimes)/10000) * 5000
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

func check(status, name, unixTime string) (string, string, bool, string) {
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

					if config["ManualBearer"].(bool) {
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

					body, err := json.Marshal(data{Name: name, Bearer: bearerGot, Unix: unixTime, Config: string(jsonValue(embeds{Content: "<@" + config["DiscordID"].(string) + ">", Embeds: []embed{{Description: fmt.Sprintf("Succesfully sniped %v :skull:", name), Color: 770000, Footer: footer{Text: "MCSN"}, Time: time.Now().Format(time.RFC3339)}}}))})

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
	sendI("This uses the first bearer within your config.json, please change the order if this is a issue.")
	sendW("Press enter to continue: ")

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

	if len(bearers.Bearers) == 0 {
		authAccs()
	}

	for i := 0; i < 26; {
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
}
