package main

import (
	"github.com/robfig/cron/v3"
	"github.com/surfaceyu/wallpaperweather/weather"
)

func main() {
	w := weather.New(weather.WSPMoji)
	w.Run()
	c := cron.New(cron.WithSeconds())
	defer c.Stop()
	_, err := c.AddFunc("0 */10 * * * *", w.WeatherCheck)
	if err != nil {
		panic(err)
	}
	c.Start()
	select {}
}
