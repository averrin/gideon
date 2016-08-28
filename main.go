package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/averrin/seker"
	ds "github.com/averrin/shodan/modules/datastream"
	eg "github.com/averrin/shodan/modules/eventghost"
	sh "github.com/averrin/shodan/modules/smarthome"
	wu "github.com/averrin/shodan/modules/weather"
	"github.com/spf13/viper"
	"github.com/veandco/go-sdl2/sdl"
)

const PADDING_TOP = 20
const PADDING_LEFT = 20

var WU wu.WUnderground
var datastream *ds.DataStream
var eventghost *eg.EventGhost
var smarthome sh.SmartHome
var windowed *bool
var FONT_SIZE int32
var icons map[string]string
var VERSION string
var SHODAN_VERSION string
var noti seker.Text

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
	eventghost = eg.Connect(viper.GetStringMapString("eventghost"))
	smarthome = sh.Connect(viper.GetStringMapString("smarthome"))
	app := new(Application)
	os.Exit(app.run())
}

type Application struct {
	seker.Application
}

func (app *Application) run() int {
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
	app.Init("Gideon", uint32(dmode), seker.Geometry{int32(w), int32(h)})
	defer app.Close()
	datastream = ds.Connect(viper.GetStringMapString("datastream"))

	seker.LoadFonts(int(FONT_SIZE))
	app.initWeather()
	app.initClock()
	app.initInterior(PADDING_LEFT, PADDING_TOP+90+(FONT_SIZE+2)*4)

	pingNetwork := app.initPinger("Network", PADDING_LEFT, PADDING_TOP+90+(FONT_SIZE+2)*8)
	pingPC := app.initPinger("PC", PADDING_LEFT, PADDING_TOP+90+(FONT_SIZE+2)*9)
	pingTwin := app.initPinger("Evil twin", PADDING_LEFT, PADDING_TOP+90+(FONT_SIZE+2)*10)

	rectV := sdl.Rect{500, 400, -1, 40}
	textV := fmt.Sprintf("%v (%v)", VERSION, SHODAN_VERSION)
	ver := seker.NewText(&rectV, textV, "#eeeeee")

	l, _ := app.Scene.AddLayer("version")
	ver.SetRules([]seker.HighlightRule{
		{0, -1, "gray", seker.DefaultFont},
	})
	l.AddItem(&ver)

	rectN := sdl.Rect{300, PADDING_TOP + 120, -1, 20}
	noti = seker.NewText(&rectN, "", "#eeeeee")

	l, _ = app.Scene.AddLayer("notification")
	// ver.SetRules([]seker.HighlightRule{
	// 	{0, -1, "gray", seker.DefaultFont},
	// })
	l.AddItem(&noti)

	go TestConnection(pingNetwork, "8.8.8.8")
	go TestConnection(pingPC, "onyx.local")
	go TestConnection(pingTwin, "evil.chip")
	go app.Scene.Run()
	datastream.Heartbeat("gideon")
	shodan := datastream.GetHeartbeat("shodan")
	pingShodan := app.initPinger("Shodan", PADDING_LEFT, PADDING_TOP+90+(FONT_SIZE+2)*11)
	go func() {
		pingShodan.SetRules([]seker.HighlightRule{seker.HighlightRule{0, -1, "red", seker.DefaultFont}})
		for {
			select {
			case ping, ok := <-shodan:
				if ok {
					if ping {
						pingShodan.SetRules([]seker.HighlightRule{seker.HighlightRule{0, -1, "green", seker.DefaultFont}})
					} else {
						pingShodan.SetRules([]seker.HighlightRule{seker.HighlightRule{0, -1, "red", seker.DefaultFont}})
					}
				} else {
					pingShodan.SetRules([]seker.HighlightRule{seker.HighlightRule{0, -1, "red", seker.DefaultFont}})
				}
			default:
			}
		}
		pingShodan.SetRules([]seker.HighlightRule{seker.HighlightRule{0, -1, "red", seker.DefaultFont}})
	}()

	commands := datastream.GetCommands("gideon")
	go func() {
		for {
			select {
			case cmd := <-commands:
				log.Println(cmd.Name)
				if cmd.Name == "kill" {
					datastream.SendStatus(ds.Status{
						"gideon", time.Now(), true, nil,
					})
					log.Println("Killed by " + cmd.Sender)
					os.Exit(1)
				} else if cmd.Name == "pause" || cmd.Name == "play" {
					eventghost.Send("p/p")
					datastream.SendStatus(ds.Status{
						"gideon", time.Now(), true, nil,
					})
				} else if cmd.Name == "noti" {
					noti.SetText(cmd.Args["message"].(string))
					datastream.SendStatus(ds.Status{
						"gideon", time.Now(), true, nil,
					})
				} else if cmd.Name == "update" {
					log.Println("updating...")
					ver.SetText("Updating...")
					err := exec.Command("/home/chip/update.sh").Start()
					if err != nil {
						datastream.SendStatus(ds.Status{
							"gideon", time.Now(), true, nil,
						})
					} else {
						datastream.SendStatus(ds.Status{
							"gideon", time.Now(), false, string(err),
						})
					}
				} else if strings.HasPrefix(cmd.Name, "eg:") {
					log.Println("Send to eg: " + cmd.Name[3:])
					eventghost.Send(cmd.Name[3:])
					datastream.SendStatus(ds.Status{
						"gideon", time.Now(), true, nil,
					})
				} else if strings.HasPrefix(cmd.Name, "sh:") {
					tokens := strings.Split(cmd.Name[3:], ":")
					log.Println("Send to sh: " + cmd.Name[3:])
					code := smarthome.GetCode(tokens[0], tokens[1])
					log.Println(code)
					smarthome.SendCode(code)
					datastream.SendStatus(ds.Status{
						"gideon", time.Now(), true, nil,
					})
				} else {
					datastream.SendStatus(ds.Status{
						"gideon", time.Now(), false, nil,
					})
				}
			default:
			}
		}
	}()

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
	errCount := 0
	go func() {
		d := 0
		tick := 0
		for {
			if errCount > 5 {
				temp.SetText(fmt.Sprintf("Temp: --"))
				hum.SetText(fmt.Sprintf("Humidity: --"))
				break
			}
			time.Sleep(time.Duration(d) * time.Second)
			outRaw, err := exec.Command("python", "sht21.py").Output()
			if err != nil {
				errCount++
				fmt.Print(err)
				continue
			}
			out := strings.TrimRight(string(outRaw), "\n")
			data := strings.Split(out, ";")
			if len(data) != 2 {
				errCount++
				fmt.Print(out + ";")
				continue
			}
			temp.SetText(fmt.Sprintf("Temp: %v°", data[0]))
			hum.SetText(fmt.Sprintf("Humidity: %v%%", data[1]))
			errCount = 0
			if tick > 12*5 {
				datastream.SetRoomTemp(data[0], data[1])
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
