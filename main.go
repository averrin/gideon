package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/averrin/seker"
	sf "github.com/averrin/shodan/modules/sparkfun"
	wu "github.com/averrin/shodan/modules/weather"
	"github.com/spf13/viper"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_ttf"
)

const WUNDER = "http://api.wunderground.com/api/%s/conditions/q/%s.json"
const PADDING_TOP = 20
const PADDING_LEFT = 20

var WU wu.WUnderground
var windowed *bool
var FONT_SIZE int32
var icons map[string]string

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
	FONT_SIZE = int32(viper.GetInt("FONT_SIZE"))

	WU = wu.Connect(viper.GetStringMapString("weather"))
	app := new(Application)
	os.Exit(app.run())
}

type Application struct {
	Window   *sdl.Window
	Renderer *sdl.Renderer
	Surface  *sdl.Surface
	Scene    *seker.Scene
}

func (app *Application) GetSurface() *sdl.Surface {
	return app.Surface
}

func (app *Application) GetWindow() *sdl.Window {
	return app.Window
}

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
	app.Scene = seker.NewScene(app, seker.Geometry{int32(w), int32(h)})
	renderer.Present()
	time.Sleep(5)
	app.Window.UpdateSurface()

	seker.LoadFonts(int(FONT_SIZE))
	app.initWeather()
	app.initClock()
	app.initInterior(PADDING_LEFT, PADDING_TOP+90+(FONT_SIZE+2)*4)

	pingNetwork := app.initPinger("Network", PADDING_LEFT, PADDING_TOP+90+(FONT_SIZE+2)*8)
	pingPC := app.initPinger("PC", PADDING_LEFT, PADDING_TOP+90+(FONT_SIZE+2)*9)
	pingTwin := app.initPinger("Evil twin", PADDING_LEFT, PADDING_TOP+90+(FONT_SIZE+2)*10)

	go TestConnection(pingNetwork, "8.8.8.8")
	go TestConnection(pingPC, "192.168.1.38")
	go TestConnection(pingTwin, "192.168.1.36")
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
	weather := WU.GetWeather()

	rectI := sdl.Rect{30, PADDING_TOP, 100, 100}
	textI := icons[weather.Icon]
	icon := seker.NewText(&rectI, textI, "#eeeeee")

	fontIcons := seker.GetFont("weathericons-regular-webfont.ttf", 60)
	icon.SetFont(fontIcons)
	rect := sdl.Rect{PADDING_LEFT, PADDING_TOP + 90, -1, FONT_SIZE + 2}
	text := fmt.Sprintf("Temp: %v° (%s°)", weather.TempC, weather.FeelslikeC)
	if fmt.Sprintf("%v", weather.TempC) == weather.FeelslikeC || weather.FeelslikeC == "" {
		text = fmt.Sprintf("Temp: %v°", weather.TempC)
	}
	temp := seker.NewText(&rect, text, "#eeeeee")
	temp.SetRules([]seker.HighlightRule{seker.HighlightRule{5, -1, "orange yellow", seker.BoldFont}})

	rectH := sdl.Rect{PADDING_LEFT, PADDING_TOP + 90 + (FONT_SIZE+2)*1, -1, 20}
	textH := fmt.Sprintf("Humidity: %v", weather.RelativeHumidity)
	hum := seker.NewText(&rectH, textH, "#eeeeee")
	hum.SetRules([]seker.HighlightRule{seker.HighlightRule{9, -1, "cornflower", seker.BoldFont}})

	rectW := sdl.Rect{PADDING_LEFT, PADDING_TOP + 90 + (FONT_SIZE+2)*2, -1, 20}
	textW := fmt.Sprintf("%v", weather.Weather)
	wea := seker.NewText(&rectW, textW, "#eeeeee")
	wea.SetRules([]seker.HighlightRule{seker.HighlightRule{0, -1, "#eeeeee", seker.BoldFont}})

	blank := wu.Weather{}
	go func() {
		for {
			// sdl.Delay(5 * 1000 * 60)
			time.Sleep(5 * time.Minute)
			fmt.Print("wu-")
			weather = WU.GetWeather()
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
	l.AddItems([]seker.Drawable{
		&icon,
		&temp,
		&hum,
		&wea,
	})
}

func (app *Application) initInterior(x int32, y int32) {
	rect := sdl.Rect{x, y, -1, FONT_SIZE + 2}
	title := seker.NewText(&rect, "Interior", "#eeeeee")
	title.SetFont(seker.BoldFont)

	rectT := sdl.Rect{x, y + (FONT_SIZE+2)*1, -1, FONT_SIZE + 2}
	textT := fmt.Sprintf("Temp: %v°", 0)
	temp := seker.NewText(&rectT, textT, "#eeeeee")
	temp.SetRules([]seker.HighlightRule{seker.HighlightRule{5, -1, "orange yellow", seker.BoldFont}})

	rectH := sdl.Rect{x, y + (FONT_SIZE+2)*2, -1, 20}
	textH := fmt.Sprintf("Humidity: %v", 0)
	hum := seker.NewText(&rectH, textH, "#eeeeee")
	hum.SetRules([]seker.HighlightRule{seker.HighlightRule{9, -1, "cornflower", seker.BoldFont}})

	l, _ := app.Scene.AddLayer("interior")
	l.AddItem(&title)
	l.AddItem(&temp)
	l.AddItem(&hum)
	go func() {
		sparkfun := sf.Connect(viper.GetStringMap("sparkfun"))
		d := 0
		tick := 0
		for {
			time.Sleep(time.Duration(d) * time.Second)
			outRaw, err := exec.Command("python", "sht21.py").Output()
			if err != nil {
				fmt.Print(err)
				continue
			}
			out := strings.TrimRight(string(outRaw), "\n")
			data := strings.Split(out, ";")
			if len(data) != 2 {
				fmt.Print(out + ";")
				continue
			}
			temp.SetText(fmt.Sprintf("Temp: %v°", data[0]))
			hum.SetText(fmt.Sprintf("Humidity: %v%%", data[1]))
			if tick > 12*5 {
				sparkfun.SendRoomTemp(data[0], data[1])
				tick = 0
			}
			tick++
			d = 5
			// sdl.Delay(500)
		}
	}()
}

func (app *Application) initPinger(title string, x int32, y int32) *seker.Text {
	rect := sdl.Rect{x, y, 100, 20}
	text := "\uf111"
	icon := seker.NewText(&rect, text, "gray")
	label := seker.NewText(&sdl.Rect{x + 30, y, 100, 20}, title, "#eeeeee")

	l, _ := app.Scene.AddLayer("pinger_" + title)
	l.AddItem(&icon)
	l.AddItem(&label)
	return &icon
}

func (app *Application) initClock() {
	rect := sdl.Rect{330, PADDING_TOP + 30, 100, 20}
	clock := seker.NewText(&rect, time.Now().Format(`15:04`), "#eeeeee")
	clock.SetFont(seker.GetFont(seker.BoldFont.Name, 60))
	rect2 := sdl.Rect{350, PADDING_TOP + 90, 100, 20}
	date := seker.NewText(&rect2, time.Now().Format(`2 Jan Mon`), "#eeeeee")
	l, _ := app.Scene.AddLayer("clock")
	l.AddItem(&clock)
	l.AddItem(&date)
	go func() {
		for {
			time.Sleep(500 * time.Millisecond)
			// sdl.Delay(500)
			clock.SetText(time.Now().Format(`15:04`))
			date.SetText(time.Now().Format(`2 Jan Mon`))
		}
	}()
}
