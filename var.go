package main

import (
	"github.com/Liza-Developer/mcapi2"
)

type jsonValues struct {
	Accounts []string
	Bearers  []string
	Config   []string
	Names    []string
	Vps      []string
}

var (
	bearers       mcapi2.MCbearers
	config        map[string]interface{}
	err           error
	removeBearers []string
	accountType   []string
	dropTime      int64
	sendInfo      mcapi2.ServerInfo
	configOptions []string
	invalidAccs   []string
	content string
	useAuto bool
	authbytes []byte
	auth map[string]interface{}
	securityResult bool

	AccountsVer []string
	BearersVer []string
	ConfigsVer []string
	NamesVer []string
	VpsesVer []string
	Confirmed []string
)
