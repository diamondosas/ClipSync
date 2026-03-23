package ping

import (
	"log"
	"time"

	"github.com/go-ping/ping"
)

func PingIPS(IPS []string){
	for{
		if IPS == nil{
			time.Sleep(2 * time.Second)
		}else{
			for _, ip := range IPS{
					Pinger, err := ping.NewPinger(ip)
					if err != nil{
						log.Println(err)
					}
					Pinger.Count = 2
					Pinger.Run()
					if err != nil{
						log.Println(err)
					}
			}
		}
	}
}