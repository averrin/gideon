package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_ttf"
)

const WUNDER = "http://api.wunderground.com/api/%s/conditions/q/%s.json"
const APIKEY = "bddbf374b57d9a80"
const LOCATION = "ru/Vsevolozhsk"

var weather CurrentObservation

func main() {
	app := new(Application)
	os.Exit(app.run())
}

type Application struct {
	Window   *sdl.Window
	Renderer *sdl.Renderer
	Surface  *sdl.Surface
	Scene    *Scene
}

func (app *Application) run() int {
	weather = GetWeather()
	sdl.Init(sdl.INIT_EVERYTHING)
	ttf.Init()

	w := 800
	h := 600
	window, err := sdl.CreateWindow("Glitch", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		w, h, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	app.Window = window
	defer app.Window.Destroy()
	renderer, err := sdl.CreateRenderer(app.Window, -1, sdl.RENDERER_ACCELERATED)
	surface, err := app.Window.GetSurface()
	if err != nil {
		panic(err)
	}
	app.Renderer = renderer
	app.Surface = surface
	renderer.Clear()
	app.Scene = NewScene(app, Geometry{int32(w), int32(h)})
	renderer.Present()
	sdl.Delay(5)
	app.Window.UpdateSurface()

	LoadFonts(18)

	rect := sdl.Rect{10, 10, 100, 20}
	text := fmt.Sprintf("Temp: %v° (%s°)", weather.TempC, weather.FeelslikeC)
	if strconv.Itoa(weather.TempC) == weather.FeelslikeC || weather.FeelslikeC == "" {
		text = fmt.Sprintf("Temp: %v°", weather.TempC)
	}
	temp := NewText(&rect, text, sdl.Color{250, 250, 250, 1})

	rectH := sdl.Rect{10, 30, 100, 20}
	textH := fmt.Sprintf("Hum: %v", weather.RelativeHumidity)
	hum := NewText(&rectH, textH, sdl.Color{250, 250, 250, 1})

	rectW := sdl.Rect{10, 50, 100, 20}
	textW := fmt.Sprintf("Weather: %v", weather.Weather)
	wea := NewText(&rectW, textW, sdl.Color{250, 250, 250, 1})

	rectI := sdl.Rect{10, 70, 100, 100}
	textI := fmt.Sprintf("\uf010")
	icon := NewText(&rectI, textI, sdl.Color{250, 250, 250, 1})

	cwd, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	dir := filepath.Join(cwd, "fonts")
	fontIcons, _ := ttf.OpenFont(path.Join(dir, "weathericons-regular-webfont.ttf"), 60)
	icon.SetFont(fontIcons)

	l, _ := app.Scene.AddLayer("info")
	l.AddItem(&temp)
	l.AddItem(&hum)
	l.AddItem(&wea)
	l.AddItem(&icon)
	go app.Scene.Run()

	running := true
	for running {
		var event sdl.Event
		for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			// fmt.Print(".")
			ret := 1
			switch t := event.(type) {
			case *sdl.QuitEvent:
				ret = 0
			case *sdl.KeyDownEvent:
				// fmt.Printf("[%d ms] Keyboard\ttype:%d\tsym:%s\tmodifiers:%d\tstate:%d\trepeat:%d\n",
				// t.Timestamp, t.Type, sdl.GetScancodeName(t.Keysym.Scancode), t.Keysym.Mod, t.State, t.Repeat)
				key := sdl.GetScancodeName(t.Keysym.Scancode)
				log.Println(key)
				if t.Keysym.Sym == sdl.K_ESCAPE || t.Keysym.Sym == sdl.K_CAPSLOCK {
					ret = 0
				}
			}
			if ret == 0 {
				running = false
			}
		}
	}
	return 0
}

func GetWeather() CurrentObservation {
	url := fmt.Sprintf(WUNDER, APIKEY, LOCATION)
	log.Print(url)
	response, _ := http.Get(url)
	defer response.Body.Close()
	var r WeatherResponse
	body, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(body, &r)
	return r.CurrentObservation
}
