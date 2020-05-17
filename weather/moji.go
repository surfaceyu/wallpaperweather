package weather

import (
	"fmt"
	"log"
	"regexp"

	"github.com/gocolly/colly"
	"github.com/surfaceyu/wallpaperweather/address"
	"github.com/surfaceyu/wallpaperweather/background"
	"github.com/surfaceyu/wallpaperweather/iweather"
)

type weatherMoji struct {
	IPAddress address.IPAddress
	ImgURL    string
}

func (w *weatherMoji) GetImgURL() {

	url := fmt.Sprintf("https://tianqi.moji.com/weather/china/%s/%s", w.IPAddress.City, w.IPAddress.Region)
	c := colly.NewCollector()
	c.OnHTML("div[data-url]", func(e *colly.HTMLElement) {
		w.ImgURL = e.Attr("data-url")
	})
	c.Visit(url)
}

func (w *weatherMoji) WeatherCheck() {
	w.GetImgURL()
	r, err := regexp.Compile("[^/]+")
	if err != nil {
		log.Fatalln(err)
	}
	imgName := r.FindAllString(w.ImgURL, -1)
	path := background.GetImgPath(w.ImgURL, imgName[len(imgName)-1])
	fmt.Println("path = ", path)
	background.SetDesktopWallpaper(path, background.Stretch)
}

//Run Run
func (w *weatherMoji) Run() {
	w.WeatherCheck()
}

//NewMoji NewMoji
func newMoji() iweather.IWeather {
	return &weatherMoji{
		IPAddress: address.GetIPAddress(),
		ImgURL:    "",
	}
}
