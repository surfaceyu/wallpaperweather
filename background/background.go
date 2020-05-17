package background

import (
	"bytes"
	"errors"
	"fmt"
	"image/jpeg"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"golang.org/x/image/bmp"
	"golang.org/x/sys/windows/registry"
)

// WallpaperStyle WallpaperStyle
type WallpaperStyle uint

func (wps WallpaperStyle) String() string {
	return wallpaperStyles[wps]
}

const (
	//Fill 填充
	Fill WallpaperStyle = iota
	//Fit 适应
	Fit
	//Stretch 拉伸
	Stretch
	//Tile 平铺
	Tile
	//Center 居中
	Center
	//Cross 跨区
	Cross
)

var wallpaperStyles = map[WallpaperStyle]string{
	0: "填充",
	1: "适应",
	2: "拉伸",
	3: "平铺",
	4: "居中",
	5: "跨区"}

var (
	bgFile       string
	bgStyle      int
	sFile        string
	waitTime     int
	activeScreen bool
	passwd       bool
)

var (
	regist registry.Key
)

//PathExists PathExists
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// GetImgPath GetImgPath
func GetImgPath(imagPath string, name string) string {
	path := "/images/" + name
	bExist, _ := PathExists(path)
	if bExist == false {
		//通过http请求获取图片的流文件
		resp, _ := http.Get(imagPath)
		body, _ := ioutil.ReadAll(resp.Body)
		out, _ := os.Create("." + path)
		io.Copy(out, bytes.NewReader(body))
	}
	// SetDesktopWallpaper(path, Stretch)
	return path
}

func checkVersion() bool {
	version := GetVersion()
	major := version & 0xFF
	if major < 6 {
		return false
	}
	return true
}

//ConvertedWallpaper ConvertedWallpaper
func ConvertedWallpaper(bgfile string) string {
	file, err := os.Open(bgfile)
	checkErr(err)
	defer file.Close()

	img, err := jpeg.Decode(file) //解码
	checkErr(err)

	bmpPath := os.Getenv("USERPROFILE") + `\Local Settings\Application Data\Microsoft\Wallpaper1.bmp`
	bmpfile, err := os.Create(bmpPath)
	checkErr(err)
	defer bmpfile.Close()

	err = bmp.Encode(bmpfile, img)
	checkErr(err)
	return bmpPath
}

//SetDesktopWallpaper SetDesktopWallpaper
func SetDesktopWallpaper(bgFile string, style WallpaperStyle) error {
	// ext := filepath.Ext(bgFile)
	// vista 以下的系统需要转换jpg为bmp（xp、2003）
	// if !checkVersion() && ext != ".bmp" {
	//     setRegistString("ConvertedWallpaper", bgFile)
	bgFile = ConvertedWallpaper("." + bgFile)
	// }
	// 设置桌面背景
	// dir, err := os.Getwd()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// bgFile = dir + bgFile
	setRegistString("WallPaper", bgFile)
	/* 设置壁纸风格和展开方式
	   在Control Panel\Desktop中的两个键值将被设置
	   TileWallpaper
	    0: 图片不被平铺
	    1: 被平铺
	   WallpaperStyle
	    0:  0表示图片居中，1表示平铺
	    2:  拉伸填充整个屏幕
	    6:  拉伸适应屏幕并保持高度比
	    10: 图片被调整大小裁剪适应屏幕保持纵横比
	    22: 跨区
	*/
	var bgTileWallpaper, bgWallpaperStyle string
	bgTileWallpaper = "0"
	switch style {
	case Fill: // (Windows 7 or later)
		bgWallpaperStyle = "10"
	case Fit: // (Windows 7 or later)
		bgWallpaperStyle = "6"
	case Stretch:
		bgWallpaperStyle = "2"
	case Tile:
		bgTileWallpaper = "1"
		bgWallpaperStyle = "0"
	case Center:
		bgWallpaperStyle = "0"
	case Cross: // win10 or later
		bgWallpaperStyle = "22"
	}

	setRegistString("WallpaperStyle", bgWallpaperStyle)
	setRegistString("TileWallpaper", bgTileWallpaper)
	ok := SystemParametersInfo(SPI_SETDESKWALLPAPER, FALSE, nil, SPIF_UPDATEINIFILE|SPIF_SENDWININICHANGE)
	fmt.Println("SystemParametersInfo = ", bgFile, ok)
	if !ok {
		return errors.New("Desktop background Settings fail")
	}
	return nil
}

func setRegistString(name, value string) {
	key, err := registry.OpenKey(registry.CURRENT_USER, "Control Panel\\Desktop\\", registry.ALL_ACCESS)
	checkErr(err)
	err = key.SetStringValue(name, value)
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
