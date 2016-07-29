package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	//"time"
)

func TestConnection(icon *Text) {
	for {
		time.Sleep(10 * time.Second)
		// sdl.Delay(10 * 1000)
		//	icon.SetRules([]HighlightRule{HighlightRule{0, -1, yellow, font}})
		online := Ping()
		status := "red"
		if online {
			status = "green"
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
