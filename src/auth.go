package src

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Liza-Developer/apiGO"
	"github.com/logrusorgru/aurora"
)

func AuthAccs() {
	var AccountsVer []string
	file, _ := os.Open("accounts.txt")

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		AccountsVer = append(AccountsVer, scanner.Text())
	}

	if len(AccountsVer) == 0 {
		fmt.Print(aurora.Sprintf(aurora.Faint(aurora.White("[%v] Unable to continue, you have no Accounts added.\n")), aurora.Red("ERROR")))
		os.Exit(0)
	}

	AccountsVer = CheckDupes(AccountsVer)
	AccountsVer = grabDetails(AccountsVer)

	if !Acc.ManualBearer {
		if len(Acc.Bearers) == 0 {
			fmt.Print(aurora.Sprintf(aurora.Faint(aurora.White("[%v] No Bearers have been found, please check your details.\n")), aurora.Red("ERROR")))
			rewrite("accounts.txt", strings.Join(AccountsVer, "\n"))

			os.Exit(0)
		} else {
			AccountsVer = checkifValid(AccountsVer)

			rewrite("accounts.txt", strings.Join(AccountsVer, "\n"))

			if len(AccountsVer) != 0 {
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
			} else {
				fmt.Print(aurora.Sprintf(aurora.Faint(aurora.White("[%v] Unable to find any usable Accounts.\n")), aurora.Red("ERROR")))
				os.Exit(0)
			}
		}
	}
}

func grabDetails(AccountsVer []string) []string {
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
	} else if Acc.Bearers == nil {
		fmt.Print(aurora.Sprintf(aurora.Faint(aurora.White("Attempting to authenticate %v account(s)\n\n")), aurora.Red(len(AccountsVer))))
		for _, Accs := range apiGO.Auth(AccountsVer).Details {
			if Accs.Error != "" {
				AccountsVer = remove(AccountsVer, Accs.Email+":"+Accs.Password)
				fmt.Print(aurora.Sprintf(aurora.Faint(aurora.White("[%v] Account %v came up Invalid: %v\n")), aurora.Red("ERROR"), aurora.Red(Accs.Email), aurora.Red(Accs.Error)))
			} else {
				if Accs.Bearer != "" {
					if apiGO.CheckChange(Accs.Bearer) {
						fmt.Print(aurora.Sprintf(aurora.Faint(aurora.White("Succesfully Reauthed %v\n")), aurora.Green(Accs.Email)))
						Acc.Bearers = append(Acc.Bearers, apiGO.Bearers{
							Bearer:       Accs.Bearer,
							AuthInterval: 86400,
							AuthedAt:     time.Now().Unix(),
							Type:         Accs.AccountType,
							Email:        Accs.Email,
							Password:     Accs.Password,
							NameChange:   true,
						})
					} else {
						AccountsVer = remove(AccountsVer, Accs.Email+":"+Accs.Password)
						fmt.Print(aurora.Sprintf(aurora.Faint(aurora.White("[%v] Account %v Cannot Name Change.\n")), aurora.Red("ERROR"), aurora.Red(Accs.Email)))
					}
				} else {
					fmt.Print(aurora.Sprintf(aurora.Faint(aurora.White("[%v] Account %v bearer is nil.\n")), aurora.Red("ERROR"), aurora.Red(Accs.Email)))
				}
			}
		}
	} else if len(Acc.Bearers) < len(AccountsVer) {
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

		fmt.Print(aurora.Sprintf(aurora.Faint(aurora.White("Attempting to authenticate %v account(s)\n\n")), aurora.Red(len(AccountsVer))))
		for _, Accs := range apiGO.Auth(auth).Details {
			if Accs.Error != "" {
				AccountsVer = remove(AccountsVer, Accs.Email+":"+Accs.Password)
				fmt.Print(aurora.Sprintf(aurora.Faint(aurora.White("[%v] Account %v came up Invalid: %v\n")), aurora.Red("ERROR"), aurora.Red(Accs.Email), aurora.Red(Accs.Error)))
			} else {
				if Accs.Bearer != "" {
					if apiGO.CheckChange(Accs.Bearer) {
						fmt.Print(aurora.Sprintf(aurora.Faint(aurora.White("Succesfully Reauthed %v\n")), aurora.Green(Accs.Email)))
						Acc.Bearers = append(Acc.Bearers, apiGO.Bearers{
							Bearer:       Accs.Bearer,
							AuthInterval: 86400,
							AuthedAt:     time.Now().Unix(),
							Type:         Accs.AccountType,
							Email:        Accs.Email,
							Password:     Accs.Password,
							NameChange:   true,
						})
					} else {
						AccountsVer = remove(AccountsVer, Accs.Email+":"+Accs.Password)
						fmt.Print(aurora.Sprintf(aurora.Faint(aurora.White("[%v] Account %v Cannot Name Change.\n")), aurora.Red("ERROR"), aurora.Red(Accs.Email)))
					}
				} else {
					fmt.Print(aurora.Sprintf(aurora.Faint(aurora.White("[%v] Account %v bearer is nil.\n")), aurora.Red("ERROR"), aurora.Red(Accs.Email)))
				}
			}
		}
	} else if len(AccountsVer) < len(Acc.Bearers) {
		for _, Accs := range AccountsVer {
			for _, num := range Acc.Bearers {
				if Accs == num.Email+":"+num.Password {
					Acc.Bearers = append(Acc.Bearers, num)
				}
			}
		}
	}

	Acc.SaveConfig()
	Acc.LoadState()

	return AccountsVer
}

func checkifValid(AccountsVer []string) []string {
	var reAuth []string
	for _, Accs := range Acc.Bearers {
		f, _ := http.NewRequest("GET", "https://api.minecraftservices.com/minecraft/profile/name/boom/available", nil)
		f.Header.Set("Authorization", "Bearer "+Accs.Bearer)
		j, _ := http.DefaultClient.Do(f)

		if j.StatusCode == 401 {
			fmt.Print(aurora.Sprintf(aurora.Faint(aurora.White("Account %v turned up invalid. Attempting to Reauth\n")), aurora.Red(Accs.Email)))
			reAuth = append(reAuth, Accs.Email+":"+Accs.Password)
		}
	}

	if len(reAuth) != 0 {
		fmt.Print(aurora.Sprintf(aurora.Faint(aurora.White("\nReauthing %v Accounts..\n\n")), aurora.Red(len(reAuth))))
		bearerz := apiGO.Auth(reAuth)

		for point, data := range Acc.Bearers {
			for _, Accs := range bearerz.Details {
				if Accs.Error != "" {
					AccountsVer = remove(AccountsVer, Accs.Email+":"+Accs.Password)
					fmt.Print(aurora.Sprintf(aurora.Faint(aurora.White("[%v] Account %v came up Invalid: %v\n")), aurora.Red("ERROR"), aurora.Red(Accs.Email), aurora.Red(Accs.Error)))
				} else if Accs.Bearer != "" {
					if data.Email == Accs.Email {
						if Accs.Bearer != "" {
							if apiGO.CheckChange(Accs.Bearer) {
								fmt.Print(aurora.Sprintf(aurora.Faint(aurora.White("Succesfully Reauthed %v\n")), aurora.Green(Accs.Email)))
								data.Bearer = Accs.Bearer
								data.NameChange = true
								data.Type = Accs.AccountType
								data.Password = Accs.Password
								data.Email = Accs.Email
								data.AuthedAt = time.Now().Unix()
								Acc.Bearers[point] = data
							} else {
								AccountsVer = remove(AccountsVer, Accs.Email+":"+Accs.Password)
								fmt.Print(aurora.Sprintf(aurora.Faint(aurora.White("[%v] Account %v Cannot Name Change.\n")), aurora.Red("ERROR"), aurora.Red(Accs.Email)))
							}
						} else {
							fmt.Print(aurora.Sprintf(aurora.Faint(aurora.White("[%v] Account %v bearer is nil.\n")), aurora.Red("ERROR"), aurora.Red(Accs.Email)))
						}
					}
				}
			}
		}
	}

	Acc.SaveConfig()
	Acc.LoadState()

	return AccountsVer
}

func remove(l []string, item string) []string {
	for i, other := range l {
		if other == item {
			l = append(l[:i], l[i+1:]...)
		}
	}
	return l
}

func rewrite(path, accounts string) {
	os.Create(path)

	file, _ := os.OpenFile(path, os.O_RDWR, 0644)
	defer file.Close()

	file.WriteAt([]byte(accounts), 0)
}

// _diamondburned_#4507 thanks to them for the epic example below.

func CheckDupes(strs []string) []string {
	dedup := strs[:0] // re-use the backing array
	track := make(map[string]bool, len(strs))

	for _, str := range strs {
		if track[str] {
			continue
		}
		dedup = append(dedup, str)
		track[str] = true
	}

	return dedup
}
