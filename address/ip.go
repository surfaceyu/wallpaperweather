package address

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

// IPAddress IPAddress
type IPAddress struct {
	IP       string
	City     string
	Region   string
	Country  string
	Loc      string
	Org      string
	Timezone string
	Readme   string
}

// GetIPAddress GetIPAddress
func GetIPAddress() IPAddress {
	resp, err := http.Get("http://ipinfo.io")
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	addr := new(IPAddress)
	json.Unmarshal(body, &addr)
	return *addr
}
