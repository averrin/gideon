package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_ttf"
)

const WUNDER = "http://api.wunderground.com/api/%s/conditions/q/%s.json"
const PADDING_TOP = 10

var weather CurrentObservation
var windowed *bool
var APIKEY string
var LOCATION string
var FONT_SIZE int32

func main() {
	windowed = flag.Bool("windowed", false, "display in window")
	flag.Parse()
	viper.SetConfigType("yaml")
	viper.SetConfigName("config") // name of config file (without extension)
	viper.AddConfigPath(".")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	APIKEY = viper.GetString("APIKEY")
	LOCATION = viper.GetString("LOCATION")
	FONT_SIZE = int32(viper.GetInt("FONT_SIZE"))
	app := new(Application)
	os.Exit(app.run())
}

type Application struct {
	Window   *sdl.Window
	Renderer *sdl.Renderer
	Surface  *sdl.Surface
	Scene    *Scene
}

var icons map[string]string

func (app *Application) run() int {
	sdl.Init(sdl.INIT_EVERYTHING)
	ttf.Init()

	icons = map[string]string{
		"partlycloudy":   "\uf002",
		"cloudy":         "\uf013",
		"rain":           "\uf019",
		"clear":          "\uf00d",
		"chanceflurries": "\uf064",
		"chancerain":     "\uf01c",
		"chancesleet":    "\uf0b2",
		"chancesnow":     "\uf065",
		"chancetstorms":  "\uf010",
		"flurries":       "\uf064",
		"fog":            "\uf007",
		"hazy":           "\uf0b6",
		"mostlycloudy":   "\uf013",
		"mostlysunny":    "\uf00d",
		"partlysunny":    "\uf002",
		"sleet":          "\uf0b5",
		"snow":           "\uf01b",
		"sunny":          "\uf00d",
		"tstorms":        "\uf01d",
		"none":           "\uf07b",
	}

	w := 640
	h := 480
	dmode := sdl.WINDOW_FULLSCREEN_DESKTOP
	if *windowed {
		dmode = sdl.WINDOW_SHOWN
	}
	window, err := sdl.CreateWindow("Gideon", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		w, h, uint32(dmode))
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
	time.Sleep(5)
	app.Window.UpdateSurface()

	LoadFonts(int(FONT_SIZE))
	app.initWeather()
	app.initClock()
	pingStatus := app.initPinger()
	go TestConnection(pingStatus)
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
				// key := sdl.GetScancodeName(t.Keysym.Scancode)
				// log.Println(key)
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

func (app *Application) initWeather() {
	white := sdl.Color{250, 250, 250, 1}
	weather = GetWeather()
	// blank := CurrentObservation{}
	// if weather == blank {
	// 	fmt.Print("wx")
	// 	rectI := sdl.Rect{30, PADDING_TOP, 100, 100}
	// 	icon := NewText(&rectI, "\uf07b", white)
	// 	rectW := sdl.Rect{10, 90 + (FONT_SIZE+2)*1, -1, 20}
	// 	no := NewText(&rectW, "Weather unavailable.", white)
	// 	l, _ := app.Scene.AddLayer("info")
	// 	l.AddItem(&icon)
	// 	l.AddItem(&no)
	// 	return
	// }
	//●⬤

	rectI := sdl.Rect{30, PADDING_TOP, 100, 100}
	textI := icons[weather.Icon]
	icon := NewText(&rectI, textI, white)

	fontIcons := GetFont("weathericons-regular-webfont.ttf", 60)
	icon.SetFont(fontIcons)
	rect := sdl.Rect{10, PADDING_TOP + 90, -1, FONT_SIZE + 2}
	text := fmt.Sprintf("Temp: %v° (%s°)", weather.TempC, weather.FeelslikeC)
	if fmt.Sprintf("%v", weather.TempC) == weather.FeelslikeC || weather.FeelslikeC == "" {
		text = fmt.Sprintf("Temp: %v°", weather.TempC)
	}
	temp := NewText(&rect, text, white)
	temp.SetRules([]HighlightRule{HighlightRule{5, -1, sdl.Color{200, 200, 100, 1}, boldFont}})

	rectH := sdl.Rect{10, PADDING_TOP + 90 + (FONT_SIZE+2)*1, -1, 20}
	textH := fmt.Sprintf("Humidity: %v", weather.RelativeHumidity)
	hum := NewText(&rectH, textH, white)
	hum.SetRules([]HighlightRule{HighlightRule{9, -1, sdl.Color{100, 133, 167, 1}, boldFont}})

	rectW := sdl.Rect{10, PADDING_TOP + 90 + (FONT_SIZE+2)*2, -1, 20}
	textW := fmt.Sprintf("%v", weather.Weather)
	wea := NewText(&rectW, textW, white)
	wea.SetRules([]HighlightRule{HighlightRule{0, -1, white, boldFont}})

	blank := CurrentObservation{}
	go func() {
		for {
			// sdl.Delay(5 * 1000 * 60)
			time.Sleep(5 * time.Minute)
			fmt.Print("wu")
			weather = GetWeather()
			if weather != blank {
				text := fmt.Sprintf("Temp: %v° (%s°)", weather.TempC, weather.FeelslikeC)
				if fmt.Sprintf("%v", weather.TempC) == weather.FeelslikeC || weather.FeelslikeC == "" {
					text = fmt.Sprintf("Temp: %v°", weather.TempC)
				}
				temp.SetText(text)
				textH := fmt.Sprintf("Humidity: %v", weather.RelativeHumidity)
				hum.SetText(textH)
				wea.SetText(weather.Weather)
				icon.SetText(icons[weather.Icon])
			} else {
				icon.SetText(icons["none"])
				wea.SetText("Weather unavailable")
			}
		}
	}()

	l, _ := app.Scene.AddLayer("info")
	l.AddItems([]Drawable{
		&icon,
		&temp,
		&hum,
		&wea,
	})
}

func (app *Application) initPinger() *Text {
	rect := sdl.Rect{10, PADDING_TOP + 90 + (FONT_SIZE+2)*4, 100, 20}
	text := "\uf111"
	red := sdl.Color{246, 61, 28, 1}
	green := sdl.Color{124, 221, 23, 1}
	yellow := sdl.Color{210, 160, 62, 1}
	icon := NewText(&rect, text, yellow)
	label := NewText(&sdl.Rect{40, PADDING_TOP + 90 + (FONT_SIZE+2)*4, 100, 20}, "Network", sdl.Color{250, 250, 250, 1})

	online := Ping()
	status := red
	if online {
		status = green
	}
	icon.SetRules([]HighlightRule{HighlightRule{0, -1, status, defaultFont}})

	l, _ := app.Scene.AddLayer("pinger")
	l.AddItem(&icon)
	l.AddItem(&label)
	return &icon
}

func (app *Application) initClock() {
	rect := sdl.Rect{10, PADDING_TOP + 90 + (FONT_SIZE+2)*5, 100, 20}
	clock := NewText(&rect, "00:00", sdl.Color{250, 250, 250, 1})
	l, _ := app.Scene.AddLayer("clock")
	l.AddItem(&clock)
	go func() {
		for {
			time.Sleep(500 * time.Millisecond)
			// sdl.Delay(500)
			clock.SetText(time.Now().Format(`15:04`))
		}
	}()
}
