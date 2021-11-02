package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/Liza-Developer/mcapi"
)

var (
	name     string
	delay    float64
	dropTime int64
	y        mcapi.Payload
)

func main() {

	name = os.Args[1]
	delay, _ = strconv.ParseFloat(os.Args[2], 64)
	spread, _ := strconv.Atoi(config[`Config`].([]interface{})[0].(string))

	dropTime = mcapi.DropTime(name)

	fmt.Printf(`
    Name: %v
   Delay: %v
Droptime: %v

`, name, delay, dropTime)

	mcapi.PreSleep(dropTime)

	y = bearers.CreatePayloads(name)

	mcapi.Sleep(dropTime, delay)

	g, j, i := y.SocketSending(int64(spread))

	for _, send := range g {
		fmt.Printf("[%v] Sent @ %v\n", name, formatTime(send))
	}

	fmt.Println()

	for num, recv := range j {
		if i[num] == "200" {
			fmt.Printf("[%v] Recv @ %v | Got %v Succesfully.\n", i[num], formatTime(recv), name)
		} else {
			fmt.Printf("[%v] Recv @ %v\n", i[num], formatTime(recv))
		}
	}
}
