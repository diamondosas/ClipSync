package network

import (
	"sync"
)

type Device struct {
	Name   string
	Ip     string
	Alive  bool  
}

var (
	IPSMu    sync.Mutex
	IPS     []string

	PORT     ="9999"
	Username string

	ConnDevicesMu sync.Mutex
	ConnDevices   []Device
)


