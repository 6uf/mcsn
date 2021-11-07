package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Liza-Developer/mcapi2"
)

func init() {
	q, _ := ioutil.ReadFile("accounts.json")

	config = mcapi2.GetConfig(q)

	sendInfo = mcapi2.ServerInfo{
		Webhook: config[`Config`].([]interface{})[2].(string),
		SkinUrl: config[`Config`].([]interface{})[1].(string),
	}

	if config[`Bearers`] == nil || len(config[`Bearers`].([]interface{})) == 0 {
		func() {
			var configOptions []string
			for _, accs := range config[`Accounts`].([]interface{}) {
				accounts = append(accounts, accs.(string))
			}

			for _, accs := range config[`Config`].([]interface{}) {
				configOptions = append(configOptions, accs.(string))
			}

			bearers, err = mcapi2.Auth(accounts)
			if err != nil {
				fmt.Println(err)
			}

			if bearers.Bearers == nil {
				fmt.Println("Was unable to auth accs")
				os.Exit(0)
			} else {
				var newBearers []string

				fmt.Println(bearers.AccountType)

				for e, bearerz := range bearers.Bearers {
					newBearers = append(newBearers, bearerz+"`"+time.Now().Add(time.Duration(time.Second*86400)).Format(time.RFC850)+"`"+bearers.AccountType[e])
				}

				v, _ := json.MarshalIndent(jsonValues{Accounts: accounts, Bearers: newBearers, Config: configOptions}, "", "  ")

				ioutil.WriteFile("accounts.json", v, 0)

				if len(newBearers) == 0 {
					fmt.Println("No valid accounts.")
					os.Exit(0)
				}
			}
		}()
	} else {

		if len(config[`Bearers`].([]interface{})) > len(config[`Accounts`].([]interface{})) {
			func() {
				var configOptions []string
				for _, accs := range config[`Accounts`].([]interface{}) {
					accounts = append(accounts, accs.(string))
				}

				for _, accs := range config[`Config`].([]interface{}) {
					configOptions = append(configOptions, accs.(string))
				}

				bearers, err = mcapi2.Auth(accounts)
				if err != nil {
					fmt.Println(err)
				}

				if bearers.Bearers == nil {
					fmt.Println("Was unable to auth accs")
					os.Exit(0)
				} else {
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
				}
			}()
		} else {

			for number, accs := range config[`Bearers`].([]interface{}) {

				m := strings.Split(accs.(string), "`")

				f, _ := http.NewRequest("GET", "https://api.minecraftservices.com/minecraft/profile/name/boom/available", nil)
				f.Header.Set("Authorization", "Bearer "+m[0])
				j, _ := http.DefaultClient.Do(f)

				if j.StatusCode == 401 {
					fmt.Printf("Account %v turned up invalid. Attempted to Reauth\n", config[`Accounts`].([]interface{})[number])

					invalidAccs = append(invalidAccs, config[`Accounts`].([]interface{})[number].(string))

				} else {
					wdad, _ := time.Parse(time.RFC850, m[1])

					if time.Now().After(wdad) {
					} else {
						removeBearers = append(removeBearers, m[0]+"`"+m[1]+"`"+m[2])
						accounts = append(accounts, m[0])
						accountType = append(accountType, m[2])
					}
				}
			}

			if len(invalidAccs) != 0 {
				g, _ := mcapi2.Auth(invalidAccs)
				for i, _ := range invalidAccs {
					accounts = append(accounts, g.Bearers...)
					removeBearers = append(removeBearers, g.Bearers[i]+"`"+time.Now().Add(time.Duration(time.Second*86400)).Format(time.RFC850)+"`"+g.AccountType[i])
					accountType = append(accountType, g.AccountType[i])
				}
			}

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

					bearers, err = mcapi2.Auth(accounts)
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
}

func formatTime(t time.Time) string {
	return t.Format("15:04:05.00000")
}
