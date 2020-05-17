package weather

import "github.com/surfaceyu/wallpaperweather/iweather"

// //WallPaper WeatherWallPaper
// type WallPaper struct {
// 	weather weatherMoji
// }

// //WeatherCheck WeatherCheck
// func (w *WallPaper) WeatherCheck() {
// 	w.weather.WeatherCheck()
// }

// //Start Start
// func (w *WallPaper) Start() {
// 	w.weather.Run()
// }
const (
	WSPMoji         = "WSPMoji"
	WSPChinaWeather = "WSPChinaWeather"
)

//New new
func New(wsp string) iweather.IWeather {
	switch wsp {
	case WSPMoji:
		return newMoji()
	case WSPChinaWeather:
		return newMoji()
	}
	return newMoji()
}
