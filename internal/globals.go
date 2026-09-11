package internal

import (
	"sync"
	"time"
)

type Device struct {
	Name   string
	Ip     string
	Alive  bool  
	LastSeen time.Time
}

var (
	IPSMu    sync.Mutex
	IPS     []string

	PORT     string ="9999"
	Username string

	ConnDevicesMu sync.Mutex
	ConnDevices   []Device

	ClipHistoryMu sync.Mutex
	ClipHistory   []string
)


