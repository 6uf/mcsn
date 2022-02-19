package src

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Liza-Developer/apiGO"
)

func AuthAccs() {
	var AccountsVer []string
	file, _ := os.Open("Accounts.txt")

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		AccountsVer = append(AccountsVer, scanner.Text())
	}

	if len(AccountsVer) == 0 {
		SendE("Unable to continue, you have no Accounts added.\n")
		os.Exit(0)
	}

	grabDetails(AccountsVer)

	if !Acc.ManualBearer {
		if Acc.Bearers == nil {
			SendE("No Bearers have been found, please check your details.")
			os.Exit(0)
		} else {
			checkifValid()

			for _, Accs := range Acc.Bearers {
				if Accs.NameChange {
					if Accs.Type == "Giftcard" {
						Bearers.Details = append(Bearers.Details, apiGO.Info{
							Bearer:      Accs.Bearer,
							AccountType: Accs.Type,
							Email:       Accs.Email,
							Requests:    Acc.GcReq,
						})
					} else {
						Bearers.Details = append(Bearers.Details, apiGO.Info{
							Bearer:      Accs.Bearer,
							AccountType: Accs.Type,
							Email:       Accs.Email,
							Requests:    Acc.MFAReq,
						})
					}
				}
			}

			if Bearers.Details == nil {
				SendE("Failed to authorize your Bearers, please rerun the sniper.")
				os.Exit(0)
			}
		}
	}
}

func grabDetails(AccountsVer []string) {
	if Acc.ManualBearer {
		for _, bearer := range AccountsVer {
			if apiGO.CheckChange(bearer) {
				Bearers.Details = append(Bearers.Details, apiGO.Info{
					Bearer:      bearer,
					AccountType: isGC(bearer),
				})
			}

			time.Sleep(time.Second)
		}
	} else {
		if Acc.Bearers == nil {
			bearerz := apiGO.Auth(AccountsVer)
			if len(bearerz.Details) == 0 {
				SendE("Unable to authenticate your Account(s), please Reverify your login details.\n")
				return
			} else {
				for _, Accs := range bearerz.Details {
					Acc.Bearers = append(Acc.Bearers, apiGO.Bearers{
						Bearer:       Accs.Bearer,
						AuthInterval: 86400,
						AuthedAt:     time.Now().Unix(),
						Type:         Accs.AccountType,
						Email:        Accs.Email,
						Password:     Accs.Password,
						NameChange:   apiGO.CheckChange(Accs.Bearer),
					})
				}
				Acc.SaveConfig()
				Acc.LoadState()
			}
		} else {
			if len(Acc.Bearers) < len(AccountsVer) {
				var auth []string
				check := make(map[string]bool)

				for _, Acc := range Acc.Bearers {
					check[Acc.Email+":"+Acc.Password] = true
				}

				for _, Accs := range AccountsVer {
					if !check[Accs] {
						auth = append(auth, Accs)
					}
				}

				bearerz := apiGO.Auth(auth)

				if len(bearerz.Details) != 0 {
					for _, Accs := range bearerz.Details {
						Acc.Bearers = append(Acc.Bearers, apiGO.Bearers{
							Bearer:       Accs.Bearer,
							AuthInterval: 86400,
							AuthedAt:     time.Now().Unix(),
							Type:         Accs.AccountType,
							Email:        Accs.Email,
							Password:     Accs.Password,
							NameChange:   apiGO.CheckChange(Accs.Bearer),
						})
					}

					Acc.SaveConfig()
					Acc.LoadState()
				}
			} else if len(AccountsVer) < len(Acc.Bearers) {
				for _, Accs := range AccountsVer {
					for _, num := range Acc.Bearers {
						if Accs == num.Email+":"+num.Password {
							Acc.Bearers = append(Acc.Bearers, num)
						}
					}
				}
				Acc.SaveConfig()
				Acc.LoadState()
			}
		}
	}
}

func checkifValid() {
	var reAuth []string
	for _, Accs := range Acc.Bearers {
		f, _ := http.NewRequest("GET", "https://api.minecraftservices.com/minecraft/profile/name/boom/available", nil)
		f.Header.Set("Authorization", "Bearer "+Accs.Bearer)
		j, _ := http.DefaultClient.Do(f)

		if j.StatusCode == 401 {
			SendI(fmt.Sprintf("Account %v turned up invalid. Attempting to Reauth", Accs.Email))
			reAuth = append(reAuth, Accs.Email+":"+Accs.Password)
		}
	}

	if len(reAuth) != 0 {
		SendI(fmt.Sprintf("Reauthing %v Accounts..", len(reAuth)))
		bearerz := apiGO.Auth(reAuth)

		if len(bearerz.Details) != 0 {
			for point, data := range Acc.Bearers {
				for _, Accs := range bearerz.Details {
					if data.Email == Accs.Email {
						data.Bearer = Accs.Bearer
						data.NameChange = apiGO.CheckChange(Accs.Bearer)
						data.Type = Accs.AccountType
						data.Password = Accs.Password
						data.Email = Accs.Email
						data.AuthedAt = time.Now().Unix()
						Acc.Bearers[point] = data
						Acc.SaveConfig()
					}
				}
			}
		}
	}

	Acc.LoadState()
}
