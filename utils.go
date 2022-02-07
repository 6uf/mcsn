package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net"
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
	Unix   int64  `json:"unix"`
	Config string `json:"config"`
}

type checkDetails struct {
	Error string `json:"error"`
	Sent  string `json:"sent"`
}

type SentRequests struct {
	Requests []Details
}

type Details struct {
	Bearer     string
	SentAt     time.Time
	RecvAt     time.Time
	StatusCode string
	UnixRecv   int64
	Success    bool
	Email      string
	Type       string
}

type Conns struct {
	Client  *tls.Conn
	Payload string
}

var (
	bearers   apiGO.MCbearers
	list      []string = []string{"Liza#0002 ~ If your seeing this, join up https://discord.gg/a8EQ97ZfgK", "Liza#0002 ~ Nice Ass", "or#0001 ~ i got a dragon cock", "Noobyte#0000 ~ MMMMMMMMMM", "Noobyte#0000 ~ Touhou Epik", "peet v3#4245 ~ Cool Coder Man", "Steven's Weird#9468 ~ fuck plot armor", "Paid Snipers ~ Not worth the bill", "Steven's Weird#9468 ~ Renting a GF, id rather just buy the girl", "Pock#3483 ~ i still miss soothe", "Kqzz#0001 ~ Money Generator", "Liza#0002 ~ Taddy Was The King?", "; everest ?#7184 ~ Shit Coder"}
	proxys    []string
	used      = make(map[string]bool)
	acc       apiGO.Config
	err       error
	images    []image.Image
	thirdRow  [][]int = [][]int{{64, 16, 72, 24}, {56, 16, 64, 24}, {48, 16, 56, 24}, {40, 16, 48, 24}, {32, 16, 40, 24}, {24, 16, 32, 24}, {16, 16, 24, 24}, {8, 16, 16, 24}, {0, 16, 8, 24}}
	secondRow [][]int = [][]int{{64, 8, 72, 16}, {56, 8, 64, 16}, {48, 8, 56, 16}, {40, 8, 48, 16}, {32, 8, 40, 16}, {24, 8, 32, 16}, {16, 8, 24, 16}, {8, 8, 16, 16}, {0, 8, 8, 16}}
	firstRow  [][]int = [][]int{{64, 0, 72, 8}, {56, 0, 64, 8}, {48, 0, 56, 8}, {40, 0, 48, 8}, {32, 0, 40, 8}, {24, 0, 32, 8}, {16, 0, 24, 8}, {8, 0, 16, 8}, {0, 0, 8, 8}}
)

const rootCert = `-----BEGIN CERTIFICATE-----
MIIB+TCCAZ+gAwIBAgIJAL05LKXo6PrrMAoGCCqGSM49BAMCMFkxCzAJBgNVBAYT
AkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEwHwYDVQQKDBhJbnRlcm5ldCBXaWRn
aXRzIFB0eSBMdGQxEjAQBgNVBAMMCWxvY2FsaG9zdDAeFw0xNTEyMDgxNDAxMTNa
Fw0yNTEyMDUxNDAxMTNaMFkxCzAJBgNVBAYTAkFVMRMwEQYDVQQIDApTb21lLVN0
YXRlMSEwHwYDVQQKDBhJbnRlcm5ldCBXaWRnaXRzIFB0eSBMdGQxEjAQBgNVBAMM
CWxvY2FsaG9zdDBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABHGaaHVod0hLOR4d
66xIrtS2TmEmjSFjt+DIEcb6sM9RTKS8TZcdBnEqq8YT7m2sKbV+TEq9Nn7d9pHz
pWG2heWjUDBOMB0GA1UdDgQWBBR0fqrecDJ44D/fiYJiOeBzfoqEijAfBgNVHSME
GDAWgBR0fqrecDJ44D/fiYJiOeBzfoqEijAMBgNVHRMEBTADAQH/MAoGCCqGSM49
BAMCA0gAMEUCIEKzVMF3JqjQjuM2rX7Rx8hancI5KJhwfeKu1xbyR7XaAiEA2UT7
1xOP035EcraRmWPe7tO0LpXgMxlh2VItpc2uc2w=
-----END CERTIFICATE-----
`

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

	acc.LoadState()

	_, err = os.Stat("logs")

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

	proxys = genProxys()
	setup(proxys)

	fmt.Print(aurora.White(`
    __  _____________  __
   /  |/  / ___/ __/ |/ /
  / /|_/ / /___\ \/    / 
 /_/  /_/\___/___/_/|_/
 `))

	fmt.Print(aurora.Sprintf(aurora.White(`
    Ver: %v
   MOTD: %v
Proxies: %v

`), aurora.White(aurora.Sprintf("%v / %v", aurora.Bold(aurora.BrightBlack("4.50b1")), aurora.Bold(aurora.BrightBlack("Made By Liza")))), aurora.Bold(aurora.BrightBlack(MOTD())), aurora.Bold(aurora.BrightBlack(len(proxys)))))

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

func sendT(content string) {
	fmt.Print(aurora.Sprintf(aurora.White("[%v] "+content), aurora.Green("TIMER")))
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

func (account Details) check(name, searches, accType string) {
	var details checkDetails
	body, _ := json.Marshal(Data{Name: name, Bearer: account.Bearer, Unix: account.UnixRecv, Config: string(jsonValue(embeds{Content: "<@" + acc.DiscordID + ">", Embeds: []embed{{Description: fmt.Sprintf("[%v] Succesfully sniped %v with %v searches :bow_and_arrow:", accType, name, searches), Color: 770000, Footer: footer{Text: "MCSN"}, Time: time.Now().Format(time.RFC3339)}}}))})

	req, _ := http.NewRequest("POST", "http://droptime.site/api/v2/webhook", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := http.DefaultClient.Do(req)
	body, _ = ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &details)

	if details.Error != "" {
		sendE(details.Error)
	} else if details.Sent != "" {
		sendS(details.Sent)
	} else {
		sendE(fmt.Sprintf("Couldnt send the request: %v", resp.StatusCode))
	}

	removeDetails(account)
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

func skinart(name, imageFile string) {
	var accd string
	var choose string
	var bearerNum int = 0

	if !acc.ManualBearer {
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

				fmt.Println()

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
			fmt.Println()
			bearers = apiGO.Auth([]string{accd})
			fmt.Println()
		}
	} else {
		sendW("This will use the first bearer within your accounts.txt | Press enter to verify: ")
		fmt.Scanf("h")

		authAccs()
	}

	sendW("Would you like to use any previously generated skins [yes:no]: ")
	fmt.Scan(&choose)

	if strings.ContainsAny(strings.ToLower(choose), "yes ye y") {
		var folder string
		sendW("Name of the folder [case sensitive]: ")
		fmt.Scan(&folder)

		fmt.Println()

		for i := 0; i < 27; {
			changeSkin(bearerNum, fmt.Sprintf("cropped/logs/%v/base_%v.png", folder, i))
			i++
		}
	} else {
		fmt.Println()

		os.MkdirAll("cropped/logs/"+name, 0755)

		img, err := readImage("images/" + imageFile)
		if err != nil {
			sendE(err.Error())
		}

		base, err := readImage("images/base.png")
		if err != nil {
			sendE(err.Error())
		}

		if img.Bounds().Size() != image.Pt(72, 24) {
			img = resize.Resize(72, 24, img, resize.Lanczos3)

			writeImage(img, "images/"+imageFile)
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
			writeImage(imaging.Paste(base, images, image.Point{
				X: 8,
				Y: 8,
			}), fmt.Sprintf("cropped/logs/%v/base_%v.png", name, i))
		}

		for i := 0; i < 27; {
			changeSkin(bearerNum, fmt.Sprintf("cropped/logs/%v/base_%v.png", name, i))
			i++
		}
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

func changeSkin(bearerNum int, path string) {
	client := resty.New()
	skin, _ := client.R().SetAuthToken(bearers.Details[bearerNum].Bearer).SetFormData(map[string]string{
		"variant": "slim",
	}).SetFile(path, path).Post("https://api.minecraftservices.com/minecraft/profile/skins")

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
				names, drops = threeLetters(charType)
			}

			for e, name := range names {
				if delay == 0 {
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
	case "turbo":
		for {
			var leng float64
			var data SentRequests
			var wg sync.WaitGroup

			payload := bearers.CreatePayloads(name)

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
							SentAt:     sendTime,
							RecvAt:     recvTime,
							StatusCode: string(ea[9:12]),
							Success:    strings.Contains(string(ea[9:12]), "200"),
							UnixRecv:   recvTime.Unix(),
							Email:      account.Email,
							Type:       account.AccountType,
						})

						wg.Done()
					}(e, account)
					time.Sleep(time.Duration(acc.SpreadPerReq) * time.Microsecond)
				}
			}

			wg.Wait()

			for _, status := range data.Requests {
				if status.Success {
					status.check(name, "0", status.Type)

					if acc.ChangeskinOnSnipe {
						sendInfo := apiGO.ServerInfo{
							SkinUrl: acc.ChangeSkinLink,
						}

						sendInfo.ChangeSkin(jsonValue(skinUrls{Url: sendInfo.SkinUrl, Varient: "slim"}), status.Bearer)
					}

					sendS("Succesfully Claimed " + name + " " + status.StatusCode)

					break
				} else {
					sendI(fmt.Sprintf("Failed to claim %v | %v", name, status.StatusCode))
				}
			}

			sendI("Sending 2 requests in a minute.")

			time.Sleep(time.Minute)

			fmt.Println()
		}
	}

	fmt.Println()

	sendW("Press CTRL+C to Continue : ")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
}

func authAccs() {
	var AccountsVer []string
	file, _ := os.Open("accounts.txt")

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		AccountsVer = append(AccountsVer, scanner.Text())
	}

	if len(AccountsVer) == 0 {
		sendE("Unable to continue, you have no accounts added.\n")
		os.Exit(0)
	}

	grabDetails(AccountsVer)

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

func grabDetails(AccountsVer []string) {
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

func checkVer(name string, delay float64, dropTime int64) {
	var content string
	var leng float64
	var data SentRequests

	searches, _ := apiGO.Search(name)

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
				ea := make([]byte, 4096)
				payload.Conns[e].Read(ea)
				recvTime := time.Now()

				data.Requests = append(data.Requests, Details{
					Bearer:     account.Bearer,
					SentAt:     sendTime,
					RecvAt:     recvTime,
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

	sort.Slice(data.Requests, func(i, j int) bool {
		return data.Requests[i].SentAt.Before(data.Requests[j].SentAt)
	})

	for _, request := range data.Requests {
		if request.Success {
			content += fmt.Sprintf("+ Sent @ %v | [%v] @ %v ~ %v\n", formatTime(request.SentAt), request.StatusCode, formatTime(request.RecvAt), request.Email)
			sendS(fmt.Sprintf("Sent @ %v | [%v] @ %v ~ %v", formatTime(request.SentAt), request.StatusCode, formatTime(request.RecvAt), request.Email))

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

			request.check(name, searches, request.Type)

			fmt.Println()

			sendI("If you enjoy using MCSN feel free to join the discord! https://discord.gg/a8EQ97ZfgK")
			break
		} else {
			content += fmt.Sprintf("- Sent @ %v | [%v] @ %v ~ %v\n", formatTime(request.SentAt), request.StatusCode, formatTime(request.RecvAt), request.Email)
			sendI(fmt.Sprintf("Sent @ %v | [%v] @ %v ~ %v", formatTime(request.SentAt), request.StatusCode, formatTime(request.RecvAt), request.Email))
		}
	}

	logSnipe(content, name)
}

// code from Alien https://github.com/wwhtrbbtt/AlienSniper

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
				var ts apiGO.Config
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

func removeDetails(account Details) {
	var new []apiGO.Bearers
	for _, accs := range acc.Bearers {
		if account.Email != accs.Email {
			new = append(new, accs)
		}
	}

	acc.Bearers = new

	var meow []apiGO.Info
	for _, accs := range acc.Bearers {
		for _, acc := range bearers.Details {
			if acc.Email != accs.Email {
				meow = append(meow, acc)
			}
		}
	}

	bearers.Details = meow

	var accz []string
	file, _ := os.Open("accounts.txt")
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Split(scanner.Text(), ":")[0] != account.Email {
			accz = append(accz, scanner.Text())
		}
	}

	rewrite("accounts.txt", strings.Join(accz, "\n"))

	acc.Logs = append(acc.Logs, apiGO.Logs{
		Email:   account.Email,
		Send:    account.SentAt,
		Recv:    account.RecvAt,
		Success: account.Success,
	})

	acc.SaveConfig()
	acc.LoadState()
}

func rewrite(path, accounts string) {
	os.Create(path)

	file, _ := os.OpenFile(path, os.O_RDWR, 0644)
	defer file.Close()

	file.WriteAt([]byte(accounts), 0)
}

func MOTD() string {
	rand.Seed(time.Now().UnixNano())
	return list[rand.Intn(len(list))]
}

// Proxy code

func proxy(name string, delay float64, dropTime int64) {
	var leng float64
	var wg sync.WaitGroup
	var data SentRequests
	var content string

	searches, _ := apiGO.Search(name)

	sendI(fmt.Sprintf("Name: %v | Delay: %v | Searches: %v | Proxys: %v\n", name, delay, searches, len(proxys)))

	for time.Now().Before(time.Unix(dropTime, 0).Add(-time.Second * 35)) {
		sendT(fmt.Sprintf("Generating Proxy Connections In: %v      \r", time.Until(time.Unix(dropTime, 0).Add(-time.Second*35)).Round(time.Second).Seconds()))
		time.Sleep(time.Second * 1)
	}

	clients := genSockets(proxys, name)

	fmt.Print("\n\n")

	apiGO.Sleep(dropTime, delay)

	fmt.Println()

	for e, account := range bearers.Details {
		if e == len(clients) {
			break
		} else {
			if account.AccountType == "Giftcard" {
				leng = float64(acc.GcReq)
			} else {
				leng = float64(acc.MFAReq)
			}

			for i := 0; float64(i) < leng; i++ {
				wg.Add(1)
				go func(account apiGO.Info, e int) {
					fmt.Fprintln(clients[e].Client, clients[e].Payload)
					sendTime := time.Now()

					var ea = make([]byte, 4096)

					clients[e].Client.Read(ea)
					recvTime := time.Now()

					data.Requests = append(data.Requests, Details{
						Bearer:     account.Bearer,
						SentAt:     sendTime,
						RecvAt:     recvTime,
						StatusCode: string(ea[9:12]),
						Success:    string(ea[9:12]) == "200",
						UnixRecv:   recvTime.Unix(),
						Email:      account.Email,
						Type:       account.AccountType,
					})

					wg.Done()
				}(account, e)

				time.Sleep(time.Duration(acc.SpreadPerReq) * time.Microsecond)
			}
		}
	}

	wg.Wait()

	sort.Slice(data.Requests, func(i, j int) bool {
		return data.Requests[i].SentAt.Before(data.Requests[j].SentAt)
	})

	for _, request := range data.Requests {
		if request.Success {
			content += fmt.Sprintf("+ Sent @ %v | [%v] @ %v ~ %v\n", formatTime(request.SentAt), request.StatusCode, formatTime(request.RecvAt), request.Email)
			sendS(fmt.Sprintf("Sent @ %v | [%v] @ %v ~ %v", formatTime(request.SentAt), request.StatusCode, formatTime(request.RecvAt), request.Email))

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

			request.check(name, searches, "Proxy")

			fmt.Println()

			sendI("If you enjoy using MCSN feel free to join the discord! https://discord.gg/a8EQ97ZfgK")
			break
		} else {
			content += fmt.Sprintf("- Sent @ %v | [%v] @ %v ~ %v\n", formatTime(request.SentAt), request.StatusCode, formatTime(request.RecvAt), request.Email)
			sendI(fmt.Sprintf("Sent @ %v | [%v] @ %v ~ %v", formatTime(request.SentAt), request.StatusCode, formatTime(request.RecvAt), request.Email))
		}
	}

	logSnipe(content, name)
}

func genProxys() []string {
	var proxys []string

	f, _ := os.Open("proxys.txt")
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		proxys = append(proxys, scanner.Text())
	}

	return proxys
}

func randomInt(proxys []string) string {
	for {
		rand.Seed(time.Now().UnixNano())
		proxy := proxys[rand.Intn(len(proxys))]
		if !used[proxy] {
			used[proxy] = true
			return proxy
		}
	}
}

func setup(proxy []string) {
	for _, proxy := range proxy {
		used[proxy] = false
	}
}

func genSockets(proxy []string, name string) []Conns {
	var DataType []Conns

	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM([]byte(rootCert))
	if !ok {
		log.Fatal("failed to parse root certificate")
	}

	for e, bearer := range bearers.Details {
		if e == len(proxy) {
			break
		} else {
			proxy := randomInt(proxy)
			conn, err := net.Dial("tcp", proxy)
			if err != nil {
				sendE(err.Error())
			} else {

				conn.Write([]byte("CONNECT api.minecraftservices.com:443 HTTP/1.1\r\nHost: api.minecraftservices.com:443\r\nProxy-Connection: keep-alive\r\nUser-Agent: MCSN/1.1\r\n\r\n"))

				var junk = make([]byte, 4096)

				conn.Read(junk)

				config := &tls.Config{RootCAs: roots, InsecureSkipVerify: true, ServerName: strings.Split(proxy, ":")[0]}

				tls := tls.Client(conn, config)

				if bearer.AccountType == "Giftcard" {
					DataType = append(DataType, Conns{
						Client:  tls,
						Payload: fmt.Sprintf("POST /minecraft/profile HTTP/1.1\r\nHost: api.minecraftservices.com\r\nConnection: open\r\nContent-Length:%s\r\nContent-Type: application/json\r\nAccept: application/json\r\nAuthorization: Bearer %s\r\n\r\n"+string([]byte(`{"profileName":"`+name+`"}`))+"\r\n", strconv.Itoa(len(string([]byte(`{"profileName":"`+name+`"}`)))), bearer.Bearer),
					})
				} else {
					DataType = append(DataType, Conns{
						Client:  tls,
						Payload: "PUT /minecraft/profile/name/" + name + " HTTP/1.1\r\nHost: api.minecraftservices.com:443\r\nUser-Agent: MCSN/1.0\r\nAuthorization: bearer " + bearer.Bearer + "\r\n\r\n",
					})
				}
			}
		}
	}

	return DataType
}

//
