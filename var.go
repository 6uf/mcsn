package main

import (
	"context"

	"github.com/Liza-Developer/api"
	"github.com/bwmarrin/discordgo"
)

var (
	bearers        api.MCbearers
	config         map[string]interface{}
	err            error
	dropTime       int64
	sendInfo       api.ServerInfo
	content        string
	useAuto        bool
	authbytes      []byte
	auth           map[string]interface{}
	securityResult bool
	names          []string
	s              *discordgo.Session
	ctx            context.Context
	cancel         context.CancelFunc
	drops          []int64

	AccountsVer []string
	BearersVer  []string
	ConfigsVer  []string
	NamesVer    []string
	VpsesVer    []string
	Confirmed   []string
)
