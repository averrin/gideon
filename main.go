package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_ttf"
)

const WUNDER = "http://api.wunderground.com/api/%s/conditions/q/%s.json"
const APIKEY = "bddbf374b57d9a80"
const LOCATION = "ru/Vsevolozhsk"
const FONT_SIZE = 24
const PADDING_TOP = 10

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

var icons map[string]string

func (app *Application) run() int {
	sdl.Init(sdl.INIT_EVERYTHING)
	ttf.Init()

	icons = map[string]string{
		"partlycloudy":   "\uf002",
		"cloudy":         "\uf013",
		"rain":           "\uf019",
		"clear":          "\uf00d",
		"chanceflurries": "\uf07b",
		"chancerain":     "\uf07b",
		"chancesleet":    "\uf07b",
		"chancesnow":     "\uf07b",
		"chancetstorms":  "\uf07b",
		"flurries":       "\uf07b",
		"fog":            "\uf003",
		"hazy":           "\uf0b6",
		"mostlycloudy":   "\uf013",
		"mostlysunny":    "\uf00d",
		"partlysunny":    "\uf07b",
		"sleet":          "\uf07b",
		"snow":           "\uf07b",
		"sunny":          "\uf07b",
		"tstorms":        "\uf01d",
	}

	w := 800
	h := 600
	window, err := sdl.CreateWindow("Gideon", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		w, h, sdl.WINDOW_FULLSCREEN_DESKTOP)
		// w, h, sdl.WINDOW_SHOWN)
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

	LoadFonts(FONT_SIZE)
	app.initWeather()
	app.initClock()
	//pingStatus := app.initPinger()
	//go TestConnection(pingStatus)
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
	blank := CurrentObservation{}
	if weather == blank {
		fmt.Print("wx")
		rectI := sdl.Rect{30, PADDING_TOP, 100, 100}
		icon := NewText(&rectI, "\uf07b", white)
		rectW := sdl.Rect{10, 90 + (FONT_SIZE+2)*1, -1, 20}
		no := NewText(&rectW, "Weather unavailable.", white)
		l, _ := app.Scene.AddLayer("info")
		l.AddItem(&icon)
		l.AddItem(&no)
		return
	}
	//●⬤

	rectI := sdl.Rect{30, PADDING_TOP, 100, 100}
	textI := icons[weather.Icon]
	icon := NewText(&rectI, textI, white)

	cwd, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	dir := filepath.Join(cwd, "fonts")
	fontIcons, _ := ttf.OpenFont(path.Join(dir, "weathericons-regular-webfont.ttf"), 60)
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

	go func() {
		for {
			sdl.Delay(5 * 1000 * 60)
			fmt.Print("wu")
			weather = GetWeather()
			text := fmt.Sprintf("Temp: %v° (%s°)", weather.TempC, weather.FeelslikeC)
			if fmt.Sprintf("%v", weather.TempC) == weather.FeelslikeC || weather.FeelslikeC == "" {
				text = fmt.Sprintf("Temp: %v°", weather.TempC)
			}
			temp.SetText(text)
			textH := fmt.Sprintf("Humidity: %v", weather.RelativeHumidity)
			hum.SetText(textH)
			wea.SetText(weather.Weather)
			icon.SetText(icons[weather.Icon])
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
	icon.SetRules([]HighlightRule{HighlightRule{0, -1, status, font}})

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
			// time.Sleep(500 * time.Millisecond)
                        sdl.Delay(500)
			clock.SetText(time.Now().Format(`15:04`))
		}
	}()
}
