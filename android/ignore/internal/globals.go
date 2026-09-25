package internal

import (
	"sync"
	"time"
)

type Device struct {
	Name     string
	Ip       string
	Alive    bool
	LastSeen time.Time
}

var (
	PORT string = "9999"

	Hostname string

	ConnDevices   []Device
	ConnDevicesMu sync.Mutex

	ClipHistory   []string
	ClipHistoryMu sync.Mutex

	LastRecvClip []byte
	LastRecvClipMu sync.Mutex

	SecretKey  []byte = []byte("clipboardsyncapp")
)	
