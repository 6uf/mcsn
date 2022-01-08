package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"github.com/Liza-Developer/apiGO"
	"github.com/bwmarrin/discordgo"
	"github.com/logrusorgru/aurora/v3"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name: "snipe",
				Usage: `-u - This option takes input, it uses the input to snipe the name your going for.
	-d - This option is used for your delay, example 50.. this is needed.
	`,
				Action: func(c *cli.Context) error {
					snipe(c.String("u"), c.Float64("d"), "single", "")
					return nil
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "u",
						Usage: "username to snipe",
						Value: "",
					},
					&cli.Float64Flag{
						Name:        "d",
						DefaultText: "1",
						Usage:       "Snipes a few ms earlier so you can counter ping lag.",
						Value:       0,
					},
				},
			},

			{
				Name:    "botsniper",
				Aliases: []string{"bot", "b", "bs"},
				Usage:   "Runs the discord bot sniper.",
				Action: func(c *cli.Context) error {
					var err error
					s, err = discordgo.New("Bot " + config[`DiscordBotToken`].(string))
					if err != nil {
						log.Fatalf("Invalid bot parameters: %v", err)
					}

					s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
						if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
							h(s, i)
						}
					})

					runBot()
					return nil
				},
			},

			{
				Name:    "auto",
				Aliases: []string{"as", "autosniper", "a"},
				Usage:   "Auto sniper automatically snipes 3C, 3L, or 3N for you. -3c -3l -3n are the commands.",
				Subcommands: []*cli.Command{
					{
						Name:  "3c",
						Usage: "Snipe names are are a combination of Numeric and Alphabetic.",
						Action: func(c *cli.Context) error {
							snipe("", c.Float64("d"), "auto", "3c")
							return nil
						},
						Flags: []cli.Flag{
							&cli.Float64Flag{
								Name:  "d",
								Usage: "Snipes a few ms earlier so you can counter ping lag.",
								Value: 0,
							},
						},
					},
					{
						Name:  "3l",
						Usage: "Snipe only Alphabetic names.",
						Action: func(c *cli.Context) error {
							snipe("", c.Float64("d"), "auto", "3l")
							return nil
						},
						Flags: []cli.Flag{
							&cli.Float64Flag{
								Name:  "d",
								Usage: "Snipes a few ms earlier so you can counter ping lag.",
								Value: 0,
							},
						},
					},
					{
						Name:  "3n",
						Usage: "Snipe only Numeric names.",
						Action: func(c *cli.Context) error {
							snipe("", c.Float64("d"), "auto", "3n")
							return nil
						},
						Flags: []cli.Flag{
							&cli.Float64Flag{
								Name:  "d",
								Usage: "Snipes a few ms earlier so you can counter ping lag.",
								Value: 0,
							},
						},
					},
				},
			},

			{
				Name:    "ping",
				Aliases: []string{"p"},
				Usage:   "ping helps give you a rough estimate of your connection to the minecraft API.",
				Action: func(c *cli.Context) error {
					fmt.Println(aurora.Sprintf(aurora.Bold(aurora.White("Estimated Delay: %v\n")), aurora.Bold(aurora.Red(math.Round(AutoOffset())))))
					return nil
				},
			},
		},

		HideHelp:    false,
		Name:        "MCSN",
		Description: "A name sniper dedicated to premium free services",
		Version:     "3.6.5",
	}

	app.Run(os.Args)
}

func snipe(name string, delay float64, option string, charType string) {
	var useAuto bool = false

	switch option {
	case "single":
		dropTime := apiGO.DropTime(name)
		if dropTime < int64(10000) {

			fmt.Println(aurora.Sprintf(aurora.White(aurora.Bold("-!- Droptime %v: [https://www.epochconverter.com]")), aurora.Bold(aurora.BrightBlack("[UNIX]"))))
			fmt.Print(aurora.SlowBlink(aurora.BrightRed(">> ")))
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
					delay = AutoOffset()
				}

				if !config["ManualBearer"].(bool) {
					if len(bearers.Bearers) == 0 {
						bearers, _ = apiGO.Auth()
					}
				}

				checkVer(name, delay, drops[e])
			}
		}
	}

	fmt.Println(aurora.Sprintf(aurora.White(aurora.Bold("\nPress CTRL+C to Exit"))))
	fmt.Print(aurora.Sprintf(aurora.Red(aurora.Bold(">> "))))
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
}

func checkVer(name string, delay float64, dropTime int64) {
	var content string
	var sendTime []time.Time
	var leng float64

	fmt.Println(aurora.Sprintf(aurora.Bold(aurora.White(` Name: %v
Delay: %v
`)), aurora.Bold(aurora.Red(name)), aurora.Bold(aurora.Red(delay))))

	var wg sync.WaitGroup

	apiGO.PreSleep(dropTime)

	payload := bearers.CreatePayloads(name)

	apiGO.Sleep(dropTime, delay)

	fmt.Println()

	for e, account := range payload.AccountType {
		switch account {
		case "Giftcard":
			leng = config[`GcReq`].(float64)
		case "Microsoft":
			leng = config[`MFAReq`].(float64)
		}

		for i := 0; float64(i) < leng; i++ {
			wg.Add(1)
			fmt.Fprintln(payload.Conns[e], payload.Payload[e])
			sendTime = append(sendTime, time.Now())

			go func() {
				ea := make([]byte, 1000)
				payload.Conns[e].Read(ea)
				recv = append(recv, time.Now())
				statusCode = append(statusCode, string(ea[9:12]))

				wg.Done()

				if string(ea[9:12]) == `200` {
					sendInfo(string(ea[9:12]), dropTime)
				}
			}()
			time.Sleep(time.Duration(config["SpreadPerReq"].(float64)) * time.Microsecond)
		}
	}

	wg.Wait()

	for e, status := range statusCode {
		if status != "200" {
			content += fmt.Sprintf("- [DISMAL] Sent @ %v | [%v] @ %v\n", formatTime(sendTime[e]), status, formatTime(recv[e]))
			fmt.Print(aurora.Sprintf(aurora.Bold(aurora.White(("[%v] Sent @ %v | [%v] @ %v\n"))), aurora.Bold(aurora.Red("DISMAL")), aurora.Bold(aurora.Red(formatTime(sendTime[e]))), aurora.Bold(aurora.Red(status)), aurora.Bold(aurora.Red(formatTime(recv[e])))))
		} else {
			fmt.Print(aurora.Sprintf(aurora.Bold(aurora.White(("[%v] Sent @ %v | [%v] @ %v ~ %v\n"))), aurora.Bold(aurora.Green("DISMAL")), aurora.Bold(aurora.Green(formatTime(sendTime[e]))), aurora.Bold(aurora.Green(status)), aurora.Bold(aurora.Green(formatTime(recv[e]))), aurora.Bold(aurora.Green(strings.Split(emailGot, ":")[0]))))
			content += fmt.Sprintf("+ [DISMAL] Sent @ %v | [%v] @ %v ~ %v\n", formatTime(sendTime[e]), status, formatTime(recv[e]), strings.Split(emailGot, ":")[0])
		}
	}

	logSnipe(content, name)
}
