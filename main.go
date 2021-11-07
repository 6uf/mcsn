package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/Liza-Developer/mcapi2"
)

func main() {

	if len(os.Args) == 1 {
		fmt.Println("Please use the format of `go run . name delay` so the sniper can run succesfully.")
		os.Exit(0)
	}

	name = os.Args[1]
	delay, _ = strconv.ParseFloat(os.Args[2], 64)

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

	e := 0

	func() {
		for statusNum, accountType := range y.AccountType {
			switch accountType {
			case "Giftcard":
				for i := 0; i < 6; i++ {
					go func() {
						sends, recvs, statuscodes = y.SocketSending(int64(spread))

						if statuscodes[e] == "200" {
							gotNum = statusNum
						}
						e++
					}()
				}
			case "Microsoft":
				for i := 0; i < 2; i++ {
					go func() {
						sends, recvs, statuscodes = y.SocketSending(int64(spread))
					}()
				}
			}
		}

		time.Sleep(6 * time.Second)

		for _, send := range sends {
			fmt.Printf("[%v] Sent @ %v\n", name, formatTime(send))
		}

		fmt.Println()

		for i, recv := range recvs {
			if statuscodes[i] == "200" {
				fmt.Printf("[%v] Recv @ %v | Got %v Succesfully.\n", statuscodes[i], formatTime(recv), name)
				sendInfo.SendWebhook([]byte(`{
						"content": "This is a test.\n\nSent from MCSN",
						"embeds": null
					  }`))
				sendInfo.ChangeSkin([]byte(""), bearers.Bearers[gotNum])

			} else {
				fmt.Printf("[%v] Recv @ %v\n", statuscodes[i], formatTime(recv))
			}
		}

		time.Sleep(time.Second)
	}()
}
