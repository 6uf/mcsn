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
					Description: "`3nl`, `3c`, or `name-list` OR type a name and itll only queue that",
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
				var namestoAdd []string
				var accounttoUse []string
				var configOption []string
				var bearerstoUse []string
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

				accounts := strings.Split((i.ApplicationCommandData().Options[0].StringValue()), ",")
				func() {

					if config[`Vps`] == nil || len(config[`Vps`].([]interface{})) == 0 {
					} else {
						for _, accs := range config[`Vps`].([]interface{}) {
							vps = append(vps, accs.(string))
						}
					}

					if config[`Names`] == nil || len(config[`Names`].([]interface{})) == 0 {
						for _, names := range accounts {
							namestoAdd = append(namestoAdd, names)
						}
					} else {
						for _, accs := range config[`Names`].([]interface{}) {
							namestoAdd = append(namestoAdd, accs.(string))
						}

						for _, names := range accounts {
							namestoAdd = append(namestoAdd, names)
						}
					}

					if config[`Accounts`] != nil {
						for _, accs := range config[`Accounts`].([]interface{}) {
							accounttoUse = append(accounttoUse, accs.(string))
						}
					}

					if config[`Config`] != nil {
						for _, accs := range config[`Config`].([]interface{}) {
							configOption = append(configOption, accs.(string))
						}
					}

					if config[`Bearers`] != nil {
						for _, accs := range config[`Bearers`].([]interface{}) {
							bearerstoUse = append(bearerstoUse, accs.(string))
						}
					}
				}()

				v, _ := json.MarshalIndent(jsonValues{Accounts: accounttoUse, Bearers: bearerstoUse, Config: configOption, Names: namestoAdd, Vps: vps}, "", "  ")

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
				var namestoAdd []string
				var accounttoUse []string
				var configOption []string
				var bearerstoUse []string
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
					if config[`Accounts`] != nil {
						for _, accs := range config[`Accounts`].([]interface{}) {
							accounttoUse = append(accounttoUse, accs.(string))
						}
					}

					if config[`Config`] != nil {
						for _, accs := range config[`Config`].([]interface{}) {
							configOption = append(configOption, accs.(string))
						}
					}

					if config[`Bearers`] != nil {
						for _, accs := range config[`Bearers`].([]interface{}) {
							bearerstoUse = append(bearerstoUse, accs.(string))
						}
					}

					if config[`Vps`] != nil {
						for _, accs := range config[`Vps`].([]interface{}) {
							vps = append(vps, accs.(string))
						}
					}
				}()

				v, _ := json.MarshalIndent(jsonValues{Accounts: accounttoUse, Bearers: bearerstoUse, Config: configOption, Names: namestoAdd, Vps: vps}, "", "  ")

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
				var namestoAdd []string
				var accounttoUse []string
				var configOption []string
				var bearerstoUse []string
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

				if config[`Accounts`] != nil {
					for _, accs := range config[`Accounts`].([]interface{}) {
						accounttoUse = append(accounttoUse, accs.(string))
					}
				}

				if config[`Config`] != nil {
					for _, accs := range config[`Config`].([]interface{}) {
						configOption = append(configOption, accs.(string))
					}
				}

				if config[`Bearers`] != nil {
					for _, accs := range config[`Bearers`].([]interface{}) {
						bearerstoUse = append(bearerstoUse, accs.(string))
					}
				}

				if config[`Names`] != nil {
					for _, accs := range config[`Names`].([]interface{}) {
						namestoAdd = append(namestoAdd, accs.(string))
					}
				}

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

				file, err := os.Open("test.exe")
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
					if config[`Vps`] != nil {
						for _, accs := range config[`Vps`].([]interface{}) {
							vps = append(vps, accs.(string))
						}

						vps = append(vps, fmt.Sprintf("%v-%v-%v", i.ApplicationCommandData().Options[0].StringValue()+":22", i.ApplicationCommandData().Options[1].StringValue(), i.ApplicationCommandData().Options[2].StringValue()))
					} else {
						vps = append(vps, fmt.Sprintf("%v-%v-%v", i.ApplicationCommandData().Options[0].StringValue()+":22", i.ApplicationCommandData().Options[1].StringValue(), i.ApplicationCommandData().Options[2].StringValue()))
					}
				}

				v, _ := json.MarshalIndent(jsonValues{Accounts: accounttoUse, Bearers: bearerstoUse, Config: configOption, Names: namestoAdd, Vps: vps}, "", "  ")

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
				var namestoAdd []string
				var accounttoUse []string
				var configOption []string
				var bearerstoUse []string
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

				if config[`Vps`] != nil {
					for _, accs := range config[`Vps`].([]interface{}) {
						vps = append(vps, accs.(string))
					}
				}

				if config[`Config`] != nil {
					for _, accs := range config[`Config`].([]interface{}) {
						configOption = append(configOption, accs.(string))
					}
				}

				if config[`Bearers`] != nil {
					for _, accs := range config[`Bearers`].([]interface{}) {
						bearerstoUse = append(bearerstoUse, accs.(string))
					}
				}

				if config[`Names`] != nil {
					for _, accs := range config[`Names`].([]interface{}) {
						namestoAdd = append(namestoAdd, accs.(string))
					}
				}

				if config[`Accounts`] != nil {
					for _, accs := range config[`Accounts`].([]interface{}) {
						accounttoUse = append(accounttoUse, accs.(string))
					}

					accs := strings.Split(i.ApplicationCommandData().Options[0].StringValue(), ",")

					for _, m := range accs {
						accounttoUse = append(accounttoUse, m)
					}

				} else {
					accs := strings.Split(i.ApplicationCommandData().Options[0].StringValue(), ",")
					fmt.Println("test2")
					for _, m := range accs {
						accounttoUse = append(accounttoUse, m)
					}
				}

				authAccs()

				v, _ := json.MarshalIndent(jsonValues{Accounts: accounttoUse, Bearers: bearerstoUse, Config: configOption, Names: namestoAdd, Vps: vps}, "", "  ")

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

				if i.ApplicationCommandData().Options[1].StringValue() == "3nl" {
					choiceofSnipe = "3nl"
				} else if i.ApplicationCommandData().Options[1].StringValue() == "3c" {
					choiceofSnipe = "3c"
				} else if i.ApplicationCommandData().Options[1].StringValue() == "name-list" {
					choiceofSnipe = "name-list"
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

					}

					for i := 0; i < len(config[`Vps`].([]interface{})); i++ {
						var uwu []string
						ListofstringAccType = append(ListofstringAccType, uwu)
					}

					accNum := 0

					if config[`Bearers`] == nil || len(config[`Bearers`].([]interface{})) == 0 {
						var temp []string
						for _, meow := range config[`Accounts`].([]interface{}) {
							temp = append(temp, meow.(string))
						}

						bearers, _ := mcapi2.Auth(temp)

						for _, nums := range bearers.Bearers {
							Listofstring[accNum] = append(Listofstring[accNum], nums)
							accNum++
							if accNum == len(Listofstring) {
								accNum = 0
							}
						}

						for _, nums := range bearers.AccountType {
							ListofstringAccType[accNum] = append(ListofstringAccType[accNum], nums)
							accNum++
							if accNum == len(ListofstringAccType) {
								accNum = 0
							}
						}

						go addToFile(bearers)
					} else {
						for _, nums := range config[`Bearers`].([]interface{}) {
							bearer := strings.Split(nums.(string), "`")
							Listofstring[accNum] = append(Listofstring[accNum], bearer[0])
							ListofstringAccType[accNum] = append(ListofstringAccType[accNum], bearer[2])
							accNum++
							if accNum == len(Listofstring) {
								accNum = 0
							}
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

						fmt.Println(i.ApplicationCommandData().Options[0].IntValue()+int64(meow), List, choiceofSnipe, i.ApplicationCommandData().Options[1].StringValue(), id, Lists, config[`Config`].([]interface{})[1].(string), config[`Config`].([]interface{})[2].(string), config[`Config`].([]interface{})[3].(string))

						go func() {
							if choiceofSnipe == "singlename" {
								err = session.Run(fmt.Sprintf("nohup ./snipee %v %v %v %v %v %v %v %v %v &", i.ApplicationCommandData().Options[0].IntValue()+int64(meow), List, choiceofSnipe, i.ApplicationCommandData().Options[1].StringValue(), id, Lists, config[`Config`].([]interface{})[1].(string), config[`Config`].([]interface{})[2].(string), config[`Config`].([]interface{})[3].(string)))
							} else {
								err = session.Run(fmt.Sprintf("nohup ./snipee %v %v %v %v %v %v %v %v &", i.ApplicationCommandData().Options[0].IntValue()+int64(meow), List, choiceofSnipe, id, Lists, config[`Config`].([]interface{})[1].(string), config[`Config`].([]interface{})[2].(string), config[`Config`].([]interface{})[3].(string)))
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

				if config[`Vpses`] == nil || len(config[`Vps`].([]interface{})) == 0 {
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

					dstFile, err := session.Create("/root/snipee")
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

func addToFile(bearers mcapi2.MCbearers) {
	var namestoAdd []string
	var accounttoUse []string
	var configOption []string
	var vps []string
	var newBearers []string

	if config[`Vps`] != nil {
		for _, accs := range config[`Vps`].([]interface{}) {
			vps = append(vps, accs.(string))
		}
	}

	if config[`Config`] != nil {
		for _, accs := range config[`Config`].([]interface{}) {
			configOption = append(configOption, accs.(string))
		}
	}

	if config[`Names`] != nil {
		for _, accs := range config[`Names`].([]interface{}) {
			namestoAdd = append(namestoAdd, accs.(string))
		}
	}

	for e, bearerz := range bearers.Bearers {
		newBearers = append(newBearers, bearerz+"`"+time.Now().Add(time.Duration(time.Second*86400)).Format(time.RFC850)+"`"+bearers.AccountType[e])
	}

	if config[`Accounts`] != nil {
		for _, accs := range config[`Accounts`].([]interface{}) {
			accounttoUse = append(accounttoUse, accs.(string))
		}
	}

	v, _ := json.MarshalIndent(jsonValues{Accounts: accounttoUse, Bearers: newBearers, Config: configOption, Names: namestoAdd, Vps: vps}, "", "  ")

	ioutil.WriteFile("accounts.json", v, 0)
}
