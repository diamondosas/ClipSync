package ping

import (
	"log"
	"time"

	"github.com/go-ping/ping"
)

func PingIPS(IPS []string){
	if IPS != nil{
		for _, ip := range IPS{
				pinger, err := ping.NewPinger(ip)
				if err != nil{
					log.Println(err)
				}
				pinger.Count = 2
				err = pinger.Run()
				if err != nil{
					log.Println(err)
				}
				stats := pinger.Statistics()
				if stats.PacketsRecv == 0{
					log.Print(ip)
					IPS = append(IPS[:s], IPS[s+1:])
					//Remove the particaular element of the thing
				}
		}
	}else{
		time.Sleep(2 * time.Second)
		return
	}
}