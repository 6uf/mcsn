package main

import (
	"time"

	"github.com/Liza-Developer/mcap"
)

type jsonValues struct {
	Accounts []string
	Bearers  []string
	Config   []string
}

var (
	bearers       mcap.MCbearers
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
)
