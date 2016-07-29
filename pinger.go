package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/tatsushid/go-fastping"
	//"time"
)

func TestConnection(icon *Text, addr string) {
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
		fmt.Print("+")
		icon.SetRules([]HighlightRule{HighlightRule{0, -1, "asparagus", defaultFont}})
		attempts = 0
	}
	p.OnIdle = func() {
		if attempts > 5 {
			icon.SetRules([]HighlightRule{HighlightRule{0, -1, "red", defaultFont}})
			attempts = 0
		} else {
			icon.SetRules([]HighlightRule{HighlightRule{0, -1, "yellow", defaultFont}})
			attempts++
		}
	}
	p.RunLoop()
	ticker := time.NewTicker(time.Millisecond * 250)
	select {
	case <-p.Done():
		if err := p.Err(); err != nil {
			log.Printf("Ping failed: %v", err)
			icon.SetRules([]HighlightRule{HighlightRule{0, -1, "red", defaultFont}})
		}
	case <-ticker.C:
		break
	}
}

func PingGet() bool {
	url := "http://google.com"
	fmt.Print(".")
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != 200 {
		fmt.Print(err)
		body, _ := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		print(body)
		return false
	}
	return true
}

func Ping(addr string) bool {
	p := fastping.NewPinger()
	ra, err := net.ResolveIPAddr("ip4:icmp", addr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	p.AddIPAddr(ra)
	at := 0
	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		at = 0
		fmt.Printf("IP Addr: %s receive, RTT: %v\n", addr.String(), rtt)
		p.Stop()
	}
	p.OnIdle = func() {
		at++
		fmt.Print(".")
		if at > 5 {
			p.Stop()
		}
	}
	p.RunLoop()
	select {
	case <-p.Done():
		fmt.Print("=")
		return true
	default:
		fmt.Print(at)
		if at <= 5 {
			return true
		}
		return false
	}
}
