package src

import (
	"fmt"
	"time"

	"github.com/Liza-Developer/apiGO"
	"github.com/logrusorgru/aurora"
)

// code from Alien https://github.com/wwhtrbbtt/AlienSniper

func CheckAccs() {
	for {
		time.Sleep(time.Second * 10)
		// check if the last auth was more than a minute ago
		for _, Accs := range Acc.Bearers {
			if time.Now().Unix() > Accs.AuthedAt+Accs.AuthInterval {
				fmt.Print(aurora.Sprintf(aurora.Faint(aurora.White("%v is due for reauth\n")), aurora.Red(Accs.Email)))

				// authenticating Account
				bearers := apiGO.Auth([]string{Accs.Email + ":" + Accs.Password})
				for point, data := range Acc.Bearers {
					for _, Accs := range bearers.Details {
						if Accs.Bearer != "" {
							if data.Email == Accs.Email {
								data.Bearer = Accs.Bearer
								data.NameChange = apiGO.CheckChange(Accs.Bearer)
								data.Type = Accs.AccountType
								data.Password = Accs.Password
								data.Email = Accs.Email
								data.AuthedAt = time.Now().Unix()
								Acc.Bearers[point] = data
							}
						}
					}

					Acc.SaveConfig()
					Acc.LoadState()
					break // break the loop to update the info.
				}

				// if the Account isnt usable, remove it from the list
				var ts apiGO.Config
				for _, i := range Acc.Bearers {
					if i.Email != Accs.Email {
						ts.Bearers = append(ts.Bearers, i)
					}
				}

				Acc.Bearers = ts.Bearers

				Acc.SaveConfig()
				Acc.LoadState()
				break // break the loop to update the info.
			}
		}
	}
}

//
