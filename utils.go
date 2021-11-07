package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/Liza-Developer/mcap"
)

func init() {
	q, _ := ioutil.ReadFile("accounts.json")

	config = mcap.GetConfig(q)

	if config[`Bearers`] == nil {
		func() {
			var configOptions []string
			for _, accs := range config[`Accounts`].([]interface{}) {
				accounts = append(accounts, accs.(string))
			}

			for _, accs := range config[`Config`].([]interface{}) {
				configOptions = append(configOptions, accs.(string))
			}

			bearers, err = mcap.Auth(accounts)
			if err != nil {
				fmt.Println(err)
			}

			if bearers.Bearers == nil {
				bearers.Bearers = []string{}
			}

			var newBearers []string

			for e, bearerz := range bearers.Bearers {
				newBearers = append(newBearers, bearerz+"`"+time.Now().Add(time.Duration(time.Second*86400)).Format(time.RFC850)+"`"+bearers.AccountType[e])
			}

			v, _ := json.MarshalIndent(jsonValues{Accounts: accounts, Bearers: newBearers, Config: configOptions}, "", "  ")

			ioutil.WriteFile("accounts.json", v, 0)

			if len(newBearers) == 0 {
				fmt.Println("No valid accounts.")
				os.Exit(0)
			}
		}()
	} else {
		for _, accs := range config[`Bearers`].([]interface{}) {

			m := strings.Split(accs.(string), "`")

			wdad, _ := time.Parse(time.RFC850, m[1])

			if time.Now().After(wdad) {
			} else {
				removeBearers = append(removeBearers, m[0]+"`"+m[1]+"`"+m[2])
				accounts = append(accounts, m[0])
				accountType = append(accountType, m[2])
			}
		}
		var m []string
		var configOptions []string
		for _, accs := range config[`Accounts`].([]interface{}) {
			m = append(m, accs.(string))
		}
		for _, accs := range config[`Config`].([]interface{}) {
			configOptions = append(configOptions, accs.(string))
		}

		if removeBearers == nil {
			removeBearers = []string{}
		} else {
			bearers.Bearers = accounts
			bearers.AccountType = accountType
		}

		v, _ := json.MarshalIndent(jsonValues{Accounts: m, Bearers: removeBearers, Config: configOptions}, "", "  ")

		ioutil.WriteFile("accounts.json", v, 0)

		if len(bearers.Bearers) == 0 {
			fmt.Println("No valid accounts.. Attempting to reauth..")
			func() {
				var configOptions []string
				for _, accs := range config[`Accounts`].([]interface{}) {
					accounts = append(accounts, accs.(string))
				}

				for _, accs := range config[`Config`].([]interface{}) {
					configOptions = append(configOptions, accs.(string))
				}

				bearers, err = mcap.Auth(accounts)
				if err != nil {
					fmt.Println(err)
				}

				if bearers.Bearers == nil {
					bearers.Bearers = []string{}
				}

				var newBearers []string

				for e, bearerz := range bearers.Bearers {
					newBearers = append(newBearers, bearerz+"`"+time.Now().Add(time.Duration(time.Second*86400)).Format(time.RFC850)+"`"+bearers.AccountType[e])
				}

				v, _ := json.MarshalIndent(jsonValues{Accounts: accounts, Bearers: newBearers, Config: configOptions}, "", "  ")

				ioutil.WriteFile("accounts.json", v, 0)

				if len(newBearers) == 0 {
					fmt.Println("No valid accounts.")
					os.Exit(0)
				}
			}()
		}
	}
}

func formatTime(t time.Time) string {
	return t.Format("15:04:05.00000")
}
