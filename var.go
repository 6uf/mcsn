package main

import (
	"flag"

	"github.com/Liza-Developer/api"
	"github.com/bwmarrin/discordgo"
)

var (
	bearers        api.MCbearers
	config         map[string]interface{}
	dropTime       int64
	sendInfo       api.ServerInfo
	content        string
	useAuto        bool
	names          []string
	s              *discordgo.Session
	drops          []int64
	GuildID             = flag.String("guild", "", "Test guild ID. If not passed - bot registers commands globally")
	RemoveCommands      = flag.Bool("rmcmd", true, "Remove all commands after shutdowning or not")
	left           bool = false

	AccountsVer []string
	BearersVer  []string
	ConfigsVer  []string
	NamesVer    []string
	VpsesVer    []string
	Confirmed   []string
)

type Names struct {
	Names []string `json:"Names"`
}

type embeds struct {
	Content interface{} `json:"content"`
	Embeds  []embed     `json:"embeds"`
}

type embed struct {
	Description interface{} `json:"description"`
	Color       interface{} `json:"color"`
	Footer      footer      `json:"footer"`
	Time        interface{} `json:"timestamp"`
}

type footer struct {
	Text interface{} `json:"text"`
	Icon interface{} `json:"icon_url"`
}

type skinUrls struct {
	Url     interface{} `json:"url"`
	Varient interface{} `json:"variant"`
}

type Name struct {
	Names string `json:"name"`
	Drop  int64  `json:"droptime"`
}
