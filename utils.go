package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/Liza-Developer/mcapi"
)

type jsonValues struct {
	Accounts []string
	Bearers  []string
	Config   []string
}

var (
	bearers       mcapi.MCbearers
	config        map[string]interface{}
	err           error
	accounts      []string
	removeBearers []string
	accountType   []string
)

func init() {
	q, _ := ioutil.ReadFile("accounts.json")

	config = GetConfig(q)

	if len(config[`Bearers`].([]interface{})) == 0 {
		reauth()
	} else {
		for _, accs := range config[`Bearers`].([]interface{}) {

			m := strings.Split(accs.(string), "`")

			wdad, _ := time.Parse(time.RFC850, m[1])

			if time.Now().After(wdad) == true {
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

		jsonValues := jsonValues{
			Accounts: m,
			Bearers:  removeBearers,
			Config:   configOptions,
		}

		v, _ := json.MarshalIndent(jsonValues, "", "  ")

		ioutil.WriteFile("accounts.json", v, 0)

		if len(bearers.Bearers) == 0 {
			fmt.Println("No valid accounts.. Attempting to reauth..")
			reauth()
		}
	}

}

func GetConfig(owo []byte) map[string]interface{} {
	var config map[string]interface{}
	json.Unmarshal(owo, &config)
	return config
}

func remove(l []string, item string) []string {
	for i, other := range l {
		if other == item {
			l = append(l[:i], l[i+1:]...)
		}
	}
	return l
}


func formatTime(t time.Time) string {
	return t.Format("15:04:05.00000")
}

func reauth() {
	var configOptions []string
	for _, accs := range config[`Accounts`].([]interface{}) {
		accounts = append(accounts, accs.(string))
	}

	for _, accs := range config[`Config`].([]interface{}) {
		configOptions = append(configOptions, accs.(string))
	}

	bearers, err = mcapi.Auth(accounts)
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

	jsonValues := jsonValues{
		Accounts: accounts,
		Bearers:  newBearers,
		Config:   configOptions,
	}

	v, _ := json.MarshalIndent(jsonValues, "", "  ")

	ioutil.WriteFile("accounts.json", v, 0)

	if len(newBearers) == 0 {
		fmt.Println("No valid accounts.")
		os.Exit(0)
	}
}
