package main
 
import (
	"github.com/gblmar/sdkajsdk"
    "fmt"
    "log"
    "encoding/json"
    "net/http"
    "io/ioutil"
    "time"
)

const (
	getIpFmt	= "http://www.dnsmadeeasy.com/myip.jsp"
	updateIpFmt = "http://www.dnsmadeeasy.com/servlet/updateip?username=%s&password=%s&id=%d&ip=%s"
)

type Settings struct {
	Interval int
	Username string
	Password string
	Id int
	Ip string
}

func main() {

	// infity loop
    for {
	    settings := readSettings()

		settings.Ip = getExternalIp()

		result := updateIp(settings.Username, settings.Password, settings.Id, settings.Ip)
		fmt.Printf("%s\n", result)

		saveSettings(settings)

		// sleep based on interval settings
		time.Sleep(time.Duration(settings.Interval) * time.Second)
	}

}

func readSettings() Settings {
	file, err := ioutil.ReadFile("settings.json")
	if err != nil {
		log.Fatal("ioutil.ReadFile", err)
	}    
    var settings Settings
    err = json.Unmarshal(file, &settings)
    if err != nil {
		log.Fatal("json.Unmarshal", err)
	}
	return settings
}

func saveSettings(settings Settings) {
	outfile, err := json.MarshalIndent(settings, "", "  ")
    if err != nil {
		log.Fatal("json.Marshal", err)
	}
    err = ioutil.WriteFile("settings.json", outfile, 0644)
    if err != nil {
        fmt.Printf("ioutil.WriteFile: %+v", err)
    }
}

func updateIp(username string, password string, id int, ip string) string {
	url := fmt.Sprintf(updateIpFmt, username, password, id, ip)

	response, err := http.Get(url)

	if err != nil {
        return string("")
    } else {
        defer response.Body.Close()
        contents, err := ioutil.ReadAll(response.Body)
        if err != nil {
            return string("")
        }
        return string(contents)
    }
}

func getExternalIp() string {
    response, err := http.Get(getIpFmt)
    if err != nil {
        return string("")
    } else {
        defer response.Body.Close()
        contents, err := ioutil.ReadAll(response.Body)
        if err != nil {
            return string("")
        }
        return string(contents)
    }
}