package main

import (
	"time"

	"github.com/Liza-Developer/mcapi2"
)

type jsonValues struct {
	Accounts []string
	Bearers  []string
	Config   []string
}

var (
	bearers       mcapi2.MCbearers
	config        map[string]interface{}
	err           error
	accounts      []string
	removeBearers []string
	accountType   []string
	name          string
	delay         float64
	dropTime      int64
	sends         []time.Time
	recvs         []time.Time
	statuscodes   []string
	sendInfo      mcapi2.ServerInfo
	gotNum        int
	m             []string
	configOptions []string
	invalidAccs   []string
)
