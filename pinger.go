package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"time"

	"github.com/averrin/seker"
	"github.com/tatsushid/go-fastping"
)

func TestConnection(icon *seker.Text, addr string) {
	attempts := 0
	p := fastping.NewPinger()
	ra, err := net.ResolveIPAddr("ip4:icmp", addr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	p.AddIPAddr(ra)
	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		// fmt.Printf("IP Addr: %s receive, RTT: %v\n", addr.String(), rtt)
		icon.SetRules([]seker.HighlightRule{seker.HighlightRule{0, -1, "pine green", seker.DefaultFont}})
		attempts = 0
	}
	p.OnIdle = func() {
		if attempts > 10 {
			icon.SetRules([]seker.HighlightRule{seker.HighlightRule{0, -1, "red", seker.DefaultFont}})
			attempts = 0
		} else if attempts > 5 {
			icon.SetRules([]seker.HighlightRule{seker.HighlightRule{0, -1, "yellow", seker.DefaultFont}})
		}
		attempts++
	}
	p.RunLoop()
	// ticker := time.NewTicker(time.Millisecond * 250)
	select {
	case <-p.Done():
		if err := p.Err(); err != nil {
			log.Printf("Ping failed: %v", err)
			icon.SetRules([]seker.HighlightRule{seker.HighlightRule{0, -1, "red", seker.DefaultFont}})
		}
		// case <-ticker.C:
		// 	break
	}
}
