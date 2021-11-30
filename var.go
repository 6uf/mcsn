package main

import (
	"context"

	"github.com/Liza-Developer/mcapi2"
	"github.com/bwmarrin/discordgo"
)

type jsonValues struct {
	Accounts []string
	Bearers  []string
	Config   []string
	Names    []string
	Vps      []string
}

var (
	bearers        mcapi2.MCbearers
	config         map[string]interface{}
	err            error
	dropTime       int64
	sendInfo       mcapi2.ServerInfo
	content        string
	useAuto        bool
	authbytes      []byte
	auth           map[string]interface{}
	securityResult bool
	names          []string
	s              *discordgo.Session
	ctx            context.Context
	cancel         context.CancelFunc

	AccountsVer []string
	BearersVer  []string
	ConfigsVer  []string
	NamesVer    []string
	VpsesVer    []string
	Confirmed   []string
)
