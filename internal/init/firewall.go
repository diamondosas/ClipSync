package init

import (
	"log"
	"os/exec"
)

func AllowFirewall(){
	cmd := exec.Command("start", )
	err := cmd.Run()
	if err != nil{
		log.Println(err)
	}
}