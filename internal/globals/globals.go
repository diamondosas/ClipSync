package globals

import (
	"sync"
)

var (
	IPSMu    sync.Mutex
	IPS      []string
	Recieved string
	PORT     = 9999
	Username string
)
// ust waiting for some... 

type Device struct{
 	Name string
  	Ip  string
}

var ConnDevices []Device

var ClipHistory []string