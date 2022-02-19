package src

import (
	"crypto/tls"
	"fmt"
	"image"
	"image/png"
	"math"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/Liza-Developer/apiGO"
	"github.com/disintegration/imaging"
	"github.com/go-resty/resty/v2"
	"github.com/nfnt/resize"
)

func Skinart(name, imageFile string) {
	var Accd string
	var choose string
	var bearerNum int = 0

	if !Acc.ManualBearer {
		SendW("Use config bearer? [yes | no]: ")
		fmt.Scan(&choose)

		if strings.ContainsAny(strings.ToLower(choose), "yes ye y") {
			Acc.LoadState()

			if len(Acc.Bearers) == 0 {
				SendE("Unable to continue, you have no Bearers added.")
			} else {
				var email string
				SendW("Email of the Account you will use: ")
				fmt.Scan(&email)

				fmt.Println()

				for _, details := range Acc.Bearers {
					if strings.EqualFold(strings.ToLower(details.Email), strings.ToLower(email)) {
						if details.Bearer != "" {
							Bearers.Details = append(Bearers.Details, apiGO.Info{
								Bearer:      details.Bearer,
								AccountType: details.Type,
								Email:       details.Email,
							})
							break
						} else {
							SendE("Your bearer is empty.")
							break
						}
					}
				}
			}
		} else {
			SendW("Enter your Account details to continue [email:password]: ")
			fmt.Scan(&Accd)
			fmt.Println()
			Bearers = apiGO.Auth([]string{Accd})
			fmt.Println()
		}
	} else {
		SendW("This will use the first bearer within your Accounts.txt | Press enter to verify: ")
		fmt.Scanf("h")

		AuthAccs()
	}

	SendW("Would you like to use any previously generated skins [yes:no]: ")
	fmt.Scan(&choose)

	if strings.ContainsAny(strings.ToLower(choose), "yes ye y") {
		var folder string
		SendW("Name of the folder [case sensitive]: ")
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
			SendE(err.Error())
		}

		base, err := readImage("images/base.png")
		if err != nil {
			SendE(err.Error())
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
	skin, _ := client.R().SetAuthToken(Bearers.Details[bearerNum].Bearer).SetFormData(map[string]string{
		"variant": "slim",
	}).SetFile(path, path).Post("https://api.minecraftservices.com/minecraft/profile/skins")

	if skin.StatusCode() == 200 {
		SendI("Skin Changed")
	} else {
		SendE("Failed skin change. (sleeping for 30 seconds)")
		time.Sleep(30 * time.Second)
	}

	SendW("Press CTRL+C to Continue : ")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	fmt.Println()
}

func logSnipe(content string, name string) {
	logFile, err := os.OpenFile(fmt.Sprintf("logs/%v.txt", strings.ToLower(name)), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		SendE(err.Error())
	}

	defer logFile.Close()

	logFile.WriteString(content)
}

func MOTD() string {
	rand.Seed(time.Now().UnixNano())
	return list[rand.Intn(len(list))]
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
		SendI(fmt.Sprintf("Request Took: %v", math.Round(value)))
		values = append(values, value)
	}

	return mean(values), time.Since(time1)
}
