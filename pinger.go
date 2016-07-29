package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	//"time"

	"github.com/veandco/go-sdl2/sdl"
)

func TestConnection(icon *Text) {
	red := sdl.Color{246, 61, 28, 1}
	green := sdl.Color{124, 221, 23, 1}
	//yellow := sdl.Color{210, 160, 62, 1}
	for {
		time.Sleep(10 * time.Second)
		// sdl.Delay(10 * 1000)
		//	icon.SetRules([]HighlightRule{HighlightRule{0, -1, yellow, font}})
		online := Ping()
		status := red
		if online {
			status = green
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
