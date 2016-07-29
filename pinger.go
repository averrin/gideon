package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	//"time"
)

func TestConnection(icon *Text) {
	s := 0
	for {
		time.Sleep(time.Duration(s) * time.Second)
		// sdl.Delay(10 * 1000)
		//	icon.SetRules([]HighlightRule{HighlightRule{0, -1, yellow, font}})
		online := Ping()
		status := "red"
		s = 5
		if online {
			status = "asparagus"
			s = 30
			fmt.Print(":")
		}
		icon.SetRules([]HighlightRule{HighlightRule{0, -1, status, defaultFont}})
	}
}

func Ping() bool {
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
