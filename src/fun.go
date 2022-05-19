package src

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/6uf/apiGO"
	"github.com/disintegration/imaging"
	"github.com/go-resty/resty/v2"
	"github.com/iskaa02/qalam"
	"github.com/iskaa02/qalam/gradient"
	"github.com/nfnt/resize"
)

func Skinart(name, imageFile string) {
	var Accd string
	var choose string
	var bearerNum int = 0

	if !Acc.ManualBearer {
		PrintGrad("Use config bearer? [yes | no]: ")
		fmt.Scan(&choose)
		if strings.ContainsAny(strings.ToLower(choose), "yes ye y") {
			Acc.LoadState()
			if len(Acc.Bearers) == 0 {
				PrintGrad("Unable to continue, you have no Bearers added.\n")
				os.Exit(0)
			} else {
				var email string
				PrintGrad("Email of the Account you will use: ")
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
						}
					}
				}
			}
		} else {
			PrintGrad("Enter your Account details to continue [EMAIL:PASS]: ")
			fmt.Scan(&Accd)
			Bearers = apiGO.Auth([]string{Accd})
		}
	} else {
		PrintGrad("This will use the first bearer within your Accounts.txt | Press enter to verify: ")
		fmt.Scanf("h")
		AuthAccs()
	}
	PrintGrad("Would you like to use any previously generated skins [YES:NO]: ")
	fmt.Scan(&choose)

	if strings.ContainsAny(strings.ToLower(choose), "yes ye y") {
		var folder string
		PrintGrad("Name of the folder: ")
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
			PrintGrad(err.Error())
		}

		base, err := readImage("images/base.png")
		if err != nil {
			PrintGrad(err.Error())
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
	skin, err := client.R().SetAuthToken(Bearers.Details[bearerNum].Bearer).SetFormData(map[string]string{
		"variant": "slim",
	}).SetFile(path, path).Post("https://api.minecraftservices.com/minecraft/profile/skins")
	if err != nil {
		fmt.Println(err)
	} else {
		if skin.StatusCode() == 200 {
			PrintGrad(fmt.Sprintf("[%v] Skin Changed\n", skin.StatusCode()))
		} else {
			PrintGrad(fmt.Sprintf("[%v] Failed skin change. (sleeping for 30 seconds)\n", skin.StatusCode()))
			time.Sleep(30 * time.Second)
		}
	}

	PrintGrad("Press CTRL+C to Continue: ")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	fmt.Println()
}

func Logo(Data string) string {
	g, _ := gradient.NewGradientBuilder().
		HtmlColors(
			"rgb(125,110,221)",
			"rgb(90%,45%,97%)",
			"hsl(229,79%,85%)",
		).
		Mode(gradient.BlendRgb).
		Build()
	return g.Mutline(Data)
}

func PrintGrad(Text string) {
	g, _ := gradient.NewGradientBuilder().
		HtmlColors(
			"rgb(125,110,221)",
			"rgb(90%,45%,97%)",
			"hsl(229,79%,85%)",
		).
		Mode(gradient.BlendRgb).
		Build()
	g.Print(Text)
}

func AddColor(text string, color string) string {
	return qalam.Sprintf("[%v]%v[/%v]", color, text, color)
}
