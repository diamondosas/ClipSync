package ping

import (
	"time"

	"clipsync/internal/network"
)

// PingIPS sends a UDP ping to each IP and returns only those that respond with a Pong within the timeout.
func PingIPS(ips []string) []string {
	if len(ips) == 0 {
		return nil
	}

	// Drain any leftover responses in PongChan
	for {
		select {
		case <-network.PongChan:
		default:
			goto Send
		}
	}

Send:
	for _, ip := range ips {
		network.SendPing(ip)
	}

	responded := make(map[string]bool)
	timeout := time.After(1 * time.Second)

	for {
		select {
		case ip := <-network.PongChan:
			responded[ip] = true
			if len(responded) == len(ips) {
				goto Filter
			}
		case <-timeout:
			goto Filter
		}
	}

Filter:
	var reachable []string
	for _, ip := range ips {
		if responded[ip] {
			reachable = append(reachable, ip)
		}
	}
	return reachable
}