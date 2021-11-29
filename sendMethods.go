package main

import (
	"fmt"
	"time"

	"github.com/Liza-Developer/mcapi2"
)

func sendAuto(option string, delay float64) {
	leng := 0
	for _, name := range names {

		if useAuto {
			delay = AutoOffset(false)
		}

		if bearers.Bearers == nil || len(bearers.Bearers) == 0 {
			fmt.Println("Attempting to reauth accounts..")
			authAccs()
		}

		dropTime := mcapi2.DropTime(name)

		fmt.Printf("    Name: %v\n   Delay: %v\nDroptime: %v\n\n", name, delay, dropTime)

		mcapi2.PreSleep(dropTime)

		payload := bearers.CreatePayloads(name)

		mcapi2.Sleep(dropTime, delay)

		fmt.Println()

		for f, accType := range bearers.AccountType {
			switch accType {
			case "Giftcard":
				leng = 6
			case "Microsoft":
				leng = 2
			}

			for i := 0; i < leng; {
				go func() {
					send, recv, status := payload.SocketSending(int64(f))
					if status == "200" {
						content += fmt.Sprintf("+ [%v] Succesfully sniped %v\n", status, name)
						fmt.Printf("[%v] Succesfully sniped %v\n", status, name)
						sendInfo.ChangeSkin(nil, bearers.Bearers[f])
						bearers.Bearers = remove(bearers.Bearers, bearers.Bearers[f])
						bearers.AccountType = remove(bearers.AccountType, bearers.AccountType[f])
						payload.Payload = remove(payload.Payload, payload.Payload[f])

						sendInfo.SendWebhook(jsonValue(embeds{Content: nil, Embeds: []embed{{Description: fmt.Sprintf("```diff\n%v\n```", content), Color: nil}}}))

						content = `
+    __  ______________ _   __
-   /  |/  / ____/ ___// | / /
+  / /|_/ / /    \__ \/  |/ / 
- / /  / / /___ ___/ / /|  /  
+/_/  /_/\____//____/_/ |_/

`
					} else {
						content += fmt.Sprintf("- [%v] Sent @ %v | Recv @ %v\n", status, formatTime(send), formatTime(recv))
						fmt.Printf("[%v] Sent @ %v | Recv @ %v\n", status, formatTime(send), formatTime(recv))
					}
				}()
				i++
				time.Sleep(40 * time.Microsecond)
			}
		}

		content = `
+    __  ______________ _   __
-   /  |/  / ____/ ___// | / /
+  / /|_/ / /    \__ \/  |/ / 
- / /  / / /___ ___/ / /|  /  
+/_/  /_/\____//____/_/ |_/

`

	}
}

func singlesniper(name string, delay float64) {
	var leng int

	dropTime = mcapi2.DropTime(name)

	fmt.Printf(`    Name: %v
   Delay: %v
Droptime: %v

`, name, delay, formatTime(time.Unix(dropTime, 0)))

	mcapi2.PreSleep(dropTime)

	y := bearers.CreatePayloads(name)

	mcapi2.Sleep(dropTime, delay)

	fmt.Println()

	for length, accountType := range y.AccountType {
		switch accountType {
		case "Giftcard":
			leng = 6
		case "Microsoft":
			leng = 2
		}

		for i := 0; i < leng; {
			go func() {
				send, recv, status := y.SocketSending(int64(length))
				if status == "200" {
					content += fmt.Sprintf("+ [%v] Recv @ %v | Got %v Succesfully.\n", status, formatTime(recv), name)
					fmt.Printf("[%v] Recv @ %v | Got %v Succesfully.\n", status, formatTime(recv), name)
					sendInfo.SendWebhook(jsonValue(embeds{Content: nil, Embeds: []embed{{Description: fmt.Sprintf("```diff\n%v\n```", content), Color: nil}}}))
					sendInfo.ChangeSkin([]byte(""), bearers.Bearers[length])
				} else {
					content += fmt.Sprintf("- [%v] Sent @ %v | Recv @ %v\n", status, formatTime(send), formatTime(recv))
					fmt.Printf("[%v] Sent @ %v | Recv @ %v\n", status, formatTime(send), formatTime(recv))
				}
			}()
			i++
			time.Sleep(40 * time.Microsecond)
		}
	}

	time.Sleep(5 * time.Second)
}
