package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	rand2 "math/rand"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/Liza-Developer/mcapi2"
	"github.com/bwmarrin/discordgo"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

var (
	GuildID        = flag.String("guild", "", "Test guild ID. If not passed - bot registers commands globally")
	RemoveCommands = flag.Bool("rmcmd", true, "Remove all commands after shutdowning or not")
)

type accounts struct {
	Accounts []string `json:"Accounts"`
	Vps      []string `json:"Vps"`
	Logs     []string `json:"Logs"`
	MeowTime string   `json:"Key"`
}

type answerRes struct {
	ID int `json:"id,omitempty"`
}

type Names struct {
	Names []string `json:"Names"`
}

var (
	s      *discordgo.Session
	ctx    context.Context
	cancel context.CancelFunc
)

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "snipe-name",
			Description: "Enter a name and delay to begin your snipe!",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "delay",
					Description: "Delay to use.",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "queue",
					Description: "`3nl`, `3c`, `list`, `3l` OR type a name and it'll only queue that",
					Required:    true,
				},
			},
		},
		{
			Name:        "add-vps",
			Description: "Add a vps!! (in testing)",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "address",
					Description: "IP of your vps :3",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "username",
					Description: "example `root`",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "password",
					Description: "logs into vps through password.",
					Required:    true,
				},
			},
		},
		{
			Name:        "add-accounts",
			Description: "add your accounts. `email:password,email:password`",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "account-info",
					Description: "email:password,email:password",
					Required:    true,
				},
			},
		},
		{
			Name:        "delete-vps",
			Description: "Remove an account you have loaded in!",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "delete-vps",
					Description: "Delete multiple Vpses `vpsinfo,vpsinfo,vpsinfo` OR `vpsinfo` alone.",
					Required:    true,
				},
			},
		},
		{
			Name:        "vpses-loaded",
			Description: "Check the vpses you have loaded atm!",
		},
		{
			Name:        "add-names",
			Description: "This is for adding names to your list queues",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "name",
					Description: "Format `name,name,name` or `name` alone",
					Required:    true,
				},
			},
		},
		{
			Name:        "delete-names",
			Description: "Remove an account you have loaded in!",
		},
		{
			Name:        "close-sniper",
			Description: "Closes the sniper if u have a name queued.",
		},
		{
			Name:        "update-sniper",
			Description: "Update the Sniper..",
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"add-names": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			go func() {
				var id string

				q, _ := ioutil.ReadFile("accounts.json")

				config = mcapi2.GetConfig(q)

				if i.Member == nil {
					id = i.User.ID
				} else {
					id = i.Member.User.ID
				}

				if id != config[`Config`].([]interface{})[4] {
					embed := &discordgo.MessageEmbed{
						Author:      &discordgo.MessageEmbedAuthor{},
						Color:       000000, // Green
						Description: "```You are not authorized to use this Bot.```",
						Timestamp:   time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
						Title:       "Dismal Logs",
					}

					sendEmbed(embed, id)
					return
				}

				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "...",
					},
				})

				accounts := strings.Split((i.ApplicationCommandData().Options[0].StringValue()), ",")

				authAccs()

				for _, acc := range accounts {
					AccountsVer = append(AccountsVer, acc)
				}

				v, _ := json.MarshalIndent(jsonValues{Accounts: AccountsVer, Bearers: BearersVer, Config: ConfigsVer, Names: NamesVer, Vps: VpsesVer}, "", "  ")

				err := ioutil.WriteFile("accounts.json", v, 0)
				if err == nil {
					embed := &discordgo.MessageEmbed{
						Author:      &discordgo.MessageEmbedAuthor{},
						Color:       000000, // Green
						Description: fmt.Sprintf("```Succesfully added your names %v.```", i.ApplicationCommandData().Options[0].StringValue()),
						Timestamp:   time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
						Title:       "Dismal Logs",
					}

					sendEmbed(embed, id)
					return

				} else {
					embed := &discordgo.MessageEmbed{
						Author:      &discordgo.MessageEmbedAuthor{},
						Color:       000000, // Green
						Description: fmt.Sprintf("```Unsuccesfully added your names.. %v```", err),
						Timestamp:   time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
						Title:       "Dismal Logs",
					}

					sendEmbed(embed, id)
					return
				}
			}()
		},
		"delete-names": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			go func() {
				var id string

				q, _ := ioutil.ReadFile("accounts.json")

				config = mcapi2.GetConfig(q)

				if i.Member == nil {
					id = i.User.ID
				} else {
					id = i.Member.User.ID
				}

				if id != config[`Config`].([]interface{})[4] {
					embed := &discordgo.MessageEmbed{
						Author:      &discordgo.MessageEmbedAuthor{},
						Color:       000000, // Green
						Description: "```You are not authorized to use this Bot.```",
						Timestamp:   time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
						Title:       "Dismal Logs",
					}

					sendEmbed(embed, id)
					return
				}

				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "...",
					},
				})

				authAccs()

				v, _ := json.MarshalIndent(jsonValues{Accounts: AccountsVer, Bearers: BearersVer, Config: ConfigsVer, Names: []string{"placeholder"}, Vps: VpsesVer}, "", "  ")

				err := ioutil.WriteFile("accounts.json", v, 0)
				if err == nil {
					embed := &discordgo.MessageEmbed{
						Author:      &discordgo.MessageEmbedAuthor{},
						Color:       000000, // Green
						Description: "```Succesfully removed all Names```",
						Timestamp:   time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
						Title:       "Dismal Logs",
					}
					sendEmbed(embed, id)
				} else {
					embed := &discordgo.MessageEmbed{
						Author:      &discordgo.MessageEmbedAuthor{},
						Color:       000000, // Green
						Description: "```Unsuccesfully removed all Names```",
						Timestamp:   time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
						Title:       "Dismal Logs",
					}
					sendEmbed(embed, id)
				}
			}()
		},
		"add-vps": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			go func() {
				var id string

				q, _ := ioutil.ReadFile("accounts.json")

				config = mcapi2.GetConfig(q)

				if i.Member == nil {
					id = i.User.ID
				} else {
					id = i.Member.User.ID
				}

				if id != config[`Config`].([]interface{})[4] {
					embed := &discordgo.MessageEmbed{
						Author:      &discordgo.MessageEmbedAuthor{},
						Color:       000000, // Green
						Description: "```You are not authorized to use this Bot.```",
						Timestamp:   time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
						Title:       "Dismal Logs",
					}

					sendEmbed(embed, id)
					return
				}

				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "...",
					},
				})

				authAccs()

				embed := &discordgo.MessageEmbed{
					Author:      &discordgo.MessageEmbedAuthor{},
					Color:       000000, // Green
					Description: "```Adding your vps...```",
					Timestamp:   time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
					Title:       "Dismal Logs",
				}
				sendEmbed(embed, id)

				addr := i.ApplicationCommandData().Options[0].StringValue() + ":22"
				configs := &ssh.ClientConfig{
					HostKeyCallback: ssh.InsecureIgnoreHostKey(),
					User:            i.ApplicationCommandData().Options[1].StringValue(),
					Auth: []ssh.AuthMethod{
						ssh.Password(i.ApplicationCommandData().Options[2].StringValue()),
					},
				}
				conn, err := ssh.Dial("tcp", addr, configs)
				if err != nil {
					embed := &discordgo.MessageEmbed{
						Author:      &discordgo.MessageEmbedAuthor{},
						Color:       000000, // Green
						Description: fmt.Sprintf("```Failed to dial: %v```", err),
						Timestamp:   time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
						Title:       "Dismal Errors",
					}
					sendEmbed(embed, id)
					return
				}

				session, err := sftp.NewClient(conn)
				if err != nil {
					embed := &discordgo.MessageEmbed{
						Author:      &discordgo.MessageEmbedAuthor{},
						Color:       000000, // Green
						Description: fmt.Sprintf("```Failed to create session: %v```", err),
						Timestamp:   time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
						Title:       "Dismal Errors",
					}
					sendEmbed(embed, id)
					return
				}

				defer session.Close()

				sesh, err := conn.NewSession()
				if err != nil {
					embed := &discordgo.MessageEmbed{
						Author:      &discordgo.MessageEmbedAuthor{},
						Color:       000000, // Green
						Description: fmt.Sprintf("```Failed to create session: %v```", err),
						Timestamp:   time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
						Title:       "Dismal Errors",
					}
					sendEmbed(embed, id)
					return
				}
				defer sesh.Close()

				file, err := os.Open("sniper")
				if err != nil {
					embed := &discordgo.MessageEmbed{
						Author:      &discordgo.MessageEmbedAuthor{},
						Color:       000000, // Green
						Description: fmt.Sprintf("```Failed to find Directory: %v```", err),
						Timestamp:   time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
						Title:       "Dismal Errors",
					}
					sendEmbed(embed, id)
					return
				}

				dstFile, err := session.Create("/root/snipe")
				if err != nil {
					embed := &discordgo.MessageEmbed{
						Author:      &discordgo.MessageEmbedAuthor{},
						Color:       000000, // Green
						Description: fmt.Sprintf("```Failed to find Directory: %v```", err),
						Timestamp:   time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
						Title:       "Dismal Errors",
					}
					sendEmbed(embed, id)
					return
				}
				defer dstFile.Close()

				if _, err := dstFile.ReadFrom(file); err != nil {
					embed := &discordgo.MessageEmbed{
						Author:      &discordgo.MessageEmbedAuthor{},
						Color:       000000, // Green
						Description: fmt.Sprintf("```Failed to find Directory: %v```", err),
						Timestamp:   time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
						Title:       "Dismal Errors",
					}
					sendEmbed(embed, id)
					return
				} else {
					if VpsesVer != nil {
						VpsesVer = append(VpsesVer, fmt.Sprintf("%v-%v-%v", i.ApplicationCommandData().Options[0].StringValue()+":22", i.ApplicationCommandData().Options[1].StringValue(), i.ApplicationCommandData().Options[2].StringValue()))
					} else {
						VpsesVer = append(VpsesVer, fmt.Sprintf("%v-%v-%v", i.ApplicationCommandData().Options[0].StringValue()+":22", i.ApplicationCommandData().Options[1].StringValue(), i.ApplicationCommandData().Options[2].StringValue()))
					}
				}

				v, _ := json.MarshalIndent(jsonValues{Accounts: AccountsVer, Bearers: BearersVer, Config: ConfigsVer, Names: NamesVer, Vps: VpsesVer}, "", "  ")

				err = ioutil.WriteFile("accounts.json", v, 0)
				if err == nil {

					var stdoutBuf bytes.Buffer

					sesh.Stdout = &stdoutBuf
					sesh.Run("chmod +x ./snipe\n")

					embed := &discordgo.MessageEmbed{
						Author:      &discordgo.MessageEmbedAuthor{},
						Color:       000000, // Green
						Description: fmt.Sprintf("```Succesfully added your vps and uploaded the sniper. %v.```", i.ApplicationCommandData().Options[0].StringValue()),
						Timestamp:   time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
						Title:       "Dismal Logs",
					}

					sendEmbed(embed, id)
					return

				} else {
					embed := &discordgo.MessageEmbed{
						Author:      &discordgo.MessageEmbedAuthor{},
						Color:       000000, // Green
						Description: fmt.Sprintf("```Unsuccesfully added your names.. %v```", err),
						Timestamp:   time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
						Title:       "Dismal Logs",
					}

					sendEmbed(embed, id)
					return
				}

			}()
		},
		"add-accounts": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			go func() {
				var id string

				q, _ := ioutil.ReadFile("accounts.json")

				config = mcapi2.GetConfig(q)

				if i.Member == nil {
					id = i.User.ID
				} else {
					id = i.Member.User.ID
				}

				if id != config[`Config`].([]interface{})[4] {
					embed := &discordgo.MessageEmbed{
						Author:      &discordgo.MessageEmbedAuthor{},
						Color:       000000, // Green
						Description: "```You are not authorized to use this Bot.```",
						Timestamp:   time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
						Title:       "Dismal Logs",
					}

					sendEmbed(embed, id)
					return
				}

				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "...",
					},
				})

				authAccs()

				v, _ := json.MarshalIndent(jsonValues{Accounts: AccountsVer, Bearers: BearersVer, Config: ConfigsVer, Names: NamesVer, Vps: VpsesVer}, "", "  ")

				err := ioutil.WriteFile("accounts.json", v, 0)
				if err == nil {
					embed := &discordgo.MessageEmbed{
						Author:      &discordgo.MessageEmbedAuthor{},
						Color:       000000, // Green
						Description: fmt.Sprintf("```Succesfully added your accounts %v.```", i.ApplicationCommandData().Options[0].StringValue()),
						Timestamp:   time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
						Title:       "Dismal Logs",
					}

					sendEmbed(embed, id)
					return

				} else {
					embed := &discordgo.MessageEmbed{
						Author:      &discordgo.MessageEmbedAuthor{},
						Color:       000000, // Green
						Description: fmt.Sprintf("```Unsuccesfully added your accounts.. %v```", err),
						Timestamp:   time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
						Title:       "Dismal Logs",
					}

					sendEmbed(embed, id)
					return
				}

			}()
		},
		"vpses-loaded": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			go func() {
				var id string
				var vps []string

				q, _ := ioutil.ReadFile("accounts.json")

				config = mcapi2.GetConfig(q)

				if i.Member == nil {
					id = i.User.ID
				} else {
					id = i.Member.User.ID
				}

				if id != config[`Config`].([]interface{})[4] {
					embed := &discordgo.MessageEmbed{
						Author:      &discordgo.MessageEmbedAuthor{},
						Color:       000000, // Green
						Description: "```You are not authorized to use this Bot.```",
						Timestamp:   time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
						Title:       "Dismal Logs",
					}

					sendEmbed(embed, id)
					return
				}

				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "...",
					},
				})

				func() {
					if config[`Vps`] == nil || len(config[`Vps`].([]interface{})) == 0 {
						embed := &discordgo.MessageEmbed{
							Author:      &discordgo.MessageEmbedAuthor{},
							Color:       000000, // Green
							Description: "```You have no vpses loaded, please add some.```",
							Timestamp:   time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
							Title:       "Dismal Logs",
						}
						sendEmbed(embed, id)
						return
					} else {
						for _, accs := range config[`Vps`].([]interface{}) {
							vps = append(vps, accs.(string))
						}

						embed := &discordgo.MessageEmbed{
							Author:      &discordgo.MessageEmbedAuthor{},
							Color:       000000, // Green
							Description: fmt.Sprintf("```Vpses: %v```", vps),
							Timestamp:   time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
							Title:       "Dismal Logs",
						}
						sendEmbed(embed, id)
					}
				}()
			}()
		},
		"snipe-name": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			go func() {
				var id string
				var choiceofSnipe string

				q, _ := ioutil.ReadFile("accounts.json")

				config = mcapi2.GetConfig(q)

				if i.ApplicationCommandData().Options[1].StringValue() == "3n" {
					choiceofSnipe = "3n"
				} else if i.ApplicationCommandData().Options[1].StringValue() == "3c" {
					choiceofSnipe = "3c"
				} else if i.ApplicationCommandData().Options[1].StringValue() == "list" {
					choiceofSnipe = "list"
				} else if i.ApplicationCommandData().Options[1].StringValue() == "3l" {
					choiceofSnipe = "3l"
				} else {
					choiceofSnipe = "singlename"
				}

				if i.Member == nil {
					id = i.User.ID
				} else {
					id = i.Member.User.ID
				}

				if id != config[`Config`].([]interface{})[4] {
					embed := &discordgo.MessageEmbed{
						Author:      &discordgo.MessageEmbedAuthor{},
						Color:       000000, // Green
						Description: "```You are not authorized to use this Bot.```",
						Timestamp:   time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
						Title:       "Dismal Logs",
					}

					sendEmbed(embed, id)
					return
				}

				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "...",
					},
				})

				if config[`Vps`] == nil || len(config[`Vps`].([]interface{})) == 0 {
					embed := &discordgo.MessageEmbed{
						Author:      &discordgo.MessageEmbedAuthor{},
						Color:       000000, // Green
						Description: "You have no Vps's added!",
						Timestamp:   time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
						Title:       "Dismal Error",
					}
					sendEmbed(embed, id)
					return

				} else if config[`Accounts`] == nil || len(config[`Accounts`].([]interface{})) == 0 {
					embed := &discordgo.MessageEmbed{
						Author:      &discordgo.MessageEmbedAuthor{},
						Color:       000000, // Green
						Description: "You have no accounts added!",
						Timestamp:   time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
						Title:       "Dismal Error",
					}
					sendEmbed(embed, id)
					return
				} else {
					var Listofstring [][]string
					var ListofstringAccType [][]string
					for i := 0; i < len(config[`Vps`].([]interface{})); i++ {
						var uwu []string
						Listofstring = append(Listofstring, uwu)
						ListofstringAccType = append(ListofstringAccType, uwu)

					}

					accNum := 0

					for _, nums := range config[`Bearers`].([]interface{}) {
						bearer := strings.Split(nums.(string), "`")
						Listofstring[accNum] = append(Listofstring[accNum], bearer[0])
						ListofstringAccType[accNum] = append(ListofstringAccType[accNum], bearer[2])
						accNum++
						if accNum == len(Listofstring) {
							accNum = 0
						}
					}

					var num int = 0

					if choiceofSnipe == "singlename" {
						embed := &discordgo.MessageEmbed{
							Author:      &discordgo.MessageEmbedAuthor{},
							Color:       000000, // Green
							Description: fmt.Sprintf("Starting sniper...```\n\nName    - %v\nDelay   - %v\nId      - %v\nOption  - %v```", i.ApplicationCommandData().Options[1].StringValue(), i.ApplicationCommandData().Options[0].IntValue(), id, choiceofSnipe),
							Timestamp:   time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
							Title:       "Dismal Logs",
						}
						sendEmbed(embed, id)
					} else {
						embed := &discordgo.MessageEmbed{
							Author:      &discordgo.MessageEmbedAuthor{},
							Color:       000000, // Green
							Description: fmt.Sprintf("Starting sniper...```\n\nDelay   - %v\nId      - %v\nOption  - %v```", i.ApplicationCommandData().Options[0].IntValue(), id, choiceofSnipe),
							Timestamp:   time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
							Title:       "Dismal Logs",
						}
						sendEmbed(embed, id)
					}

					for meow, ips := range config[`Vps`].([]interface{}) {
						splitIps := strings.Split(ips.(string), "-")

						addr := splitIps[0]
						configs := &ssh.ClientConfig{
							HostKeyCallback: ssh.InsecureIgnoreHostKey(),
							User:            splitIps[1],
							Auth: []ssh.AuthMethod{
								ssh.Password(splitIps[2]),
							},
							//Ciphers: []string{"3des-cbc", "aes256-cbc", "aes192-cbc", "aes128-cbc"},
						}
						conn, err := ssh.Dial("tcp", addr, configs)
						if err != nil {
							embed := &discordgo.MessageEmbed{
								Author:      &discordgo.MessageEmbedAuthor{},
								Color:       000000, // Green
								Description: fmt.Sprintf("```Failed to dial: %v```", err),
								Timestamp:   time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
								Title:       "Dismal Errors",
							}
							sendEmbed(embed, id)
							return
						}

						session, err := conn.NewSession()
						if err != nil {
							log.Println(err)
						}
						var stdoutBuf bytes.Buffer

						session.Stdout = &stdoutBuf

						Listofstr, _ := json.Marshal(Listofstring[meow])
						List := string(Listofstr)
						List = strings.TrimLeft(List, "[")
						List = strings.TrimRight(List, "]")
						List = strings.TrimLeft(List, `"`)
						List = strings.TrimRight(List, `"`)
						List = strings.TrimPrefix(List, `"`)
						List = strings.TrimSuffix(List, `"`)
						List = strings.Trim(List, `"`)

						Listofstrs, _ := json.Marshal(ListofstringAccType[meow])
						Lists := string(Listofstrs)
						Lists = strings.TrimLeft(Lists, "[")
						Lists = strings.TrimRight(Lists, "]")
						Lists = strings.TrimLeft(Lists, `"`)
						Lists = strings.TrimRight(Lists, `"`)
						Lists = strings.TrimPrefix(Lists, `"`)
						Lists = strings.TrimSuffix(Lists, `"`)
						Lists = strings.Trim(Lists, `"`)

						meow := rand2.Intn(3)

						sesh, err := sftp.NewClient(conn)
						file, err := os.Open("accounts.json")
						dstFile, err := sesh.Create("/root/accounts.json")
						if err != nil {
							embed := &discordgo.MessageEmbed{
								Author:      &discordgo.MessageEmbedAuthor{},
								Color:       000000, // Green
								Description: fmt.Sprintf("```Failed to find Directory: %v```", err),
								Timestamp:   time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
								Title:       "Dismal Errors",
							}
							sendEmbed(embed, id)
							return
						}
						defer dstFile.Close()

						if _, err := dstFile.ReadFrom(file); err != nil {
							embed := &discordgo.MessageEmbed{
								Author:      &discordgo.MessageEmbedAuthor{},
								Color:       000000, // Green
								Description: fmt.Sprintf("```Failed to find Directory: %v```", err),
								Timestamp:   time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
								Title:       "Dismal Errors",
							}
							sendEmbed(embed, id)
							return
						}

						go func() {
							if choiceofSnipe == "singlename" {
								err = session.Run(fmt.Sprintf("tmux\n./snipe %v %v %v %v %v &", i.ApplicationCommandData().Options[0].IntValue()+int64(meow), i.ApplicationCommandData().Options[1].StringValue(), List, Lists, choiceofSnipe))
							} else {
								err = session.Run(fmt.Sprintf("tmux\n./snipe %v %v %v %v &", i.ApplicationCommandData().Options[0].IntValue()+int64(meow), List, Lists, choiceofSnipe))
							}
							if err == nil {
								session.Close()
							} else {
								embed := &discordgo.MessageEmbed{
									Author:      &discordgo.MessageEmbedAuthor{},
									Color:       000000, // Green
									Description: fmt.Sprintf("```Error while starting sniper: %v```", err),
									Timestamp:   time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
									Title:       "Dismal Errors",
								}
								sendEmbed(embed, id)
							}
							num++
						}()

						var empty bool
						for _, input := range Listofstring {
							if len(input) == 0 {
								empty = true
							}
						}

						if len(config[`Vps`].([]interface{})) == 1 || empty == true {
							return
						}
					}
				}
			}()
		},
		"close-sniper": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			go func() {
				var id string

				q, _ := ioutil.ReadFile("accounts.json")

				config = mcapi2.GetConfig(q)

				if i.Member == nil {
					id = i.User.ID
				} else {
					id = i.Member.User.ID
				}

				if id != config[`Config`].([]interface{})[4] {
					embed := &discordgo.MessageEmbed{
						Author:      &discordgo.MessageEmbedAuthor{},
						Color:       000000, // Green
						Description: "```You are not authorized to use this Bot.```",
						Timestamp:   time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
						Title:       "Dismal Logs",
					}

					sendEmbed(embed, id)
					return
				}

				if config[`Vps`] == nil || len(config[`Vps`].([]interface{})) == 0 {
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "No vpses loaded",
						},
					})
					return
				} else {
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "...",
						},
					})
				}

				for _, ips := range config[`Vps`].([]interface{}) {
					splitIps := strings.Split(ips.(string), "-")

					addr := splitIps[0]
					config := &ssh.ClientConfig{
						HostKeyCallback: ssh.InsecureIgnoreHostKey(),
						User:            splitIps[1],
						Auth: []ssh.AuthMethod{
							ssh.Password(splitIps[2]),
						},
					}
					conn, err := ssh.Dial("tcp", addr, config)
					if err != nil {
						embed := &discordgo.MessageEmbed{
							Author:      &discordgo.MessageEmbedAuthor{},
							Color:       000000, // Green
							Description: fmt.Sprintf("```Failed to dial: %v```", err),
							Timestamp:   time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
							Title:       "Dismal Errors",
						}
						sendEmbed(embed, id)
						return
					}

					session, err := conn.NewSession()
					defer session.Close()

					var stdoutBuf bytes.Buffer

					session.Stdout = &stdoutBuf
					session.Run(`ps ax`)

					os.Create("logs.txt")

					ioutil.WriteFile("logs.txt", stdoutBuf.Bytes(), 0644)

					file, _ := os.Open("logs.txt")
					reader := bufio.NewReader(file)

					var line string
					var foundLine []string
					for {
						line, err = reader.ReadString('\n')
						if strings.Contains(line, "./snipe") {
							foundLine = append(foundLine, line[0:7])
						}
						if len(line) == 0 {
							break
						}
					}

					var sessions *ssh.Session
					var ammountClosed int

					if len(foundLine) == 0 {
						embed := &discordgo.MessageEmbed{
							Author:      &discordgo.MessageEmbedAuthor{},
							Color:       000000, // Green
							Description: "```No instances of the sniper is currently running```",
							Timestamp:   time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
							Title:       "Dismal Errors",
						}
						sendEmbed(embed, id)
						return
					} else {
						for _, lines := range foundLine {
							sessions, _ = conn.NewSession()
							defer sessions.Close()
							err = sessions.Run(`kill -9 ` + lines)
							if err != nil {
								embed := &discordgo.MessageEmbed{
									Author:      &discordgo.MessageEmbedAuthor{},
									Color:       000000, // Green
									Description: "```Unsuccesfully closed the sniper```",
									Timestamp:   time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
									Title:       "Dismal Errors",
								}
								sendEmbed(embed, id)
							} else {
								ammountClosed++
							}
						}
					}
				}
			}()
		},
		"update-sniper": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			go func() {
				var id string

				q, _ := ioutil.ReadFile("accounts.json")

				config = mcapi2.GetConfig(q)

				if i.Member == nil {
					id = i.User.ID
				} else {
					id = i.Member.User.ID
				}

				if id != config[`Config`].([]interface{})[4] {
					embed := &discordgo.MessageEmbed{
						Author:      &discordgo.MessageEmbedAuthor{},
						Color:       000000, // Green
						Description: "```You are not authorized to use this Bot.```",
						Timestamp:   time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
						Title:       "Dismal Logs",
					}

					sendEmbed(embed, id)
					return
				}

				if config[`Vps`] == nil || len(config[`Vps`].([]interface{})) == 0 {
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "```Please add a vps to your account before using the sniper! the command is /add-vps```",
						},
					})
					return
				} else {
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "...",
						},
					})
				}

				for _, ips := range config[`Vps`].([]interface{}) {
					splitIps := strings.Split(ips.(string), "-")

					addr := splitIps[0]
					config := &ssh.ClientConfig{
						HostKeyCallback: ssh.InsecureIgnoreHostKey(),
						User:            splitIps[1],
						Auth: []ssh.AuthMethod{
							ssh.Password(splitIps[2]),
						},
					}
					conn, err := ssh.Dial("tcp", addr, config)
					if err != nil {
						embed := &discordgo.MessageEmbed{
							Author:      &discordgo.MessageEmbedAuthor{},
							Color:       000000, // Green
							Description: fmt.Sprintf("```Failed to dial: %v```", err),
							Timestamp:   time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
							Title:       "Dismal Errors",
						}
						sendEmbed(embed, id)
						return
					}

					session, err := sftp.NewClient(conn)
					if err != nil {
						embed := &discordgo.MessageEmbed{
							Author:      &discordgo.MessageEmbedAuthor{},
							Color:       000000, // Green
							Description: fmt.Sprintf("```Failed to create session: %v```", err),
							Timestamp:   time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
							Title:       "Dismal Errors",
						}
						sendEmbed(embed, id)
						return
					}
					// Close connection
					defer session.Close()

					var file *os.File

					file, err = os.Open("sniper")
					if err != nil {
						embed := &discordgo.MessageEmbed{
							Author:      &discordgo.MessageEmbedAuthor{},
							Color:       000000, // Green
							Description: fmt.Sprintf("```Failed to find Directory: %v```", err),
							Timestamp:   time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
							Title:       "Dismal Errors",
						}
						sendEmbed(embed, id)
						return
					}

					dstFile, err := session.Create("/root/snipe")
					if err != nil {
						embed := &discordgo.MessageEmbed{
							Author:      &discordgo.MessageEmbedAuthor{},
							Color:       000000, // Green
							Description: fmt.Sprintf("```Failed to find Directory: %v```", err),
							Timestamp:   time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
							Title:       "Dismal Errors",
						}
						sendEmbed(embed, id)
						return
					}
					defer dstFile.Close()

					if _, err := dstFile.ReadFrom(file); err != nil {
						embed := &discordgo.MessageEmbed{
							Author:      &discordgo.MessageEmbedAuthor{},
							Color:       000000, // Green
							Description: fmt.Sprintf("```Failed to find Directory: %v```", err),
							Timestamp:   time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
							Title:       "Dismal Errors",
						}
						sendEmbed(embed, id)
						return
					} else {
						embed := &discordgo.MessageEmbed{
							Author:      &discordgo.MessageEmbedAuthor{},
							Color:       000000, // Green
							Description: "```Succesfully updated the sniper.```",
							Timestamp:   time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
							Title:       "Dismal Logs",
						}
						sendEmbed(embed, id)
					}
				}
			}()
		},
	}
)

func sendEmbed(embed *discordgo.MessageEmbed, id string) {
	go func() {
		s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {})

		s.Open()

		channel, err := s.UserChannelCreate(id)
		if err != nil {
			log.Println("error creating channel:", err)
			return
		}

		_, err = s.ChannelMessageSendEmbed(channel.ID, embed)
		if err != nil {
			log.Println("error sending DM message:", err)
			return
		}
	}()
}

func runBot() {
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Println("Bot is up!")
	})

	err := s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}
	for _, command := range commands {
		s.ApplicationCommandCreate(s.State.User.ID, "", command)
	}

	defer s.Close()

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Gracefully shutdowning")
}

