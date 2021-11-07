package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/Liza-Developer/mcap"
)

func main() {

	name = os.Args[1]
	delay, _ = strconv.ParseFloat(os.Args[2], 64)
	spread, _ := strconv.Atoi(config[`Config`].([]interface{})[0].(string))

	dropTime = mcap.DropTime(name)

	fmt.Printf(`
    Name: %v
   Delay: %v
Droptime: %v

`, name, delay, dropTime)

	mcap.PreSleep(dropTime)

	y := bearers.CreatePayloads(name)

	mcap.Sleep(dropTime, delay)

	func() {
		for _, accountType := range y.AccountType {
			switch accountType {
			case "Giftcard":
				for i := 0; i < 6; i++ {
					go func() {
						sends, recvs, statuscodes = y.SocketSending(int64(spread))
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
			} else {
				fmt.Printf("[%v] Recv @ %v\n", statuscodes[i], formatTime(recv))
			}
		}
	}()
}
