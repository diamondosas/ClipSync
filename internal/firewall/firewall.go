package init

import (
	"log"
	"os/exec"
	"runtime"
)

func AllowFirewall() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("powershell", "-ExecutionPolicy", "Bypass", "-File", "allow-firewall-rules.ps1")
		err := cmd.Run()
		if err != nil {
			log.Println("Firewall setup:", err)
		}
	}
}