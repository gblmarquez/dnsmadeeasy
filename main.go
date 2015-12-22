package main
 
import (
    "fmt"
    "time"
    "os"
    "flag"
    "path"
    "path/filepath"
    "encoding/json"
    "net/http"
    "io/ioutil"
    "strconv"

	"bitbucket.org/kardianos/service"
	"bitbucket.org/kardianos/osext"
)

const (
	getIpFmt    	= "http://myip.dnsmadeeasy.com"
	updateIpFmt 	= "http://cp.dnsmadeeasy.com/servlet/updateip?username=%s&password=%s&id=%s&ip=%s"
	fileName 	= "dnsmadeeasy.cfg"
	helpCommands 	= "Dynamic DNS updater for DnsMadeEasy\n\nUsage:\n\n        dnsmadeeasy command [arguments]\n\n\nCommands options: \n\n  install [username] [password] [record_id]\n  remove \n  run [username] [password] [record_id]\n  start \n  stop\n"
)

type Settings struct {
	Interval int
	Username string
	Password string
	Id string
	Ip string
}

var log service.Logger
var settingsFilePath string

func main() {
	var name = "dnsmadeeasy"
	var displayName = "Dynamic DNS Updater"
	var desc = "Dynamic DNS updater for DnsMadeEasy"

	var s, err = service.NewService(name, displayName, desc)
	log = s

	if err != nil {
		fmt.Printf("%s unable to start: %s", displayName, err)
		return
	}

	rootPath, _ := osext.ExecutableFolder()
    settingsFilePath, _ = filepath.Abs(path.Join(rootPath, fileName))

	if len(os.Args) > 1 {
		flag.Parse()
		var err error
		verb := os.Args[1]
		switch verb {
		case "install":
			err = s.Install()
			if err != nil {
				fmt.Printf("Failed to install: %s\n", err)
				return
			}
			fmt.Printf("Service \"%s\" installed.\n", displayName)

			settings := readSettings()

			settings.Username = os.Args[2]
			settings.Password = os.Args[3]
			settings.Id = os.Args[4]

			saveSettings(settings)
			fmt.Printf("Created config file on '%s'", settingsFilePath)

		case "remove":
			err = s.Remove()
			if err != nil {
				fmt.Printf("Failed to remove: %s\n", err)
				return
			}
			fmt.Printf("Service \"%s\" removed.\n", displayName)
		case "run":			

			settings := readSettings()

			settings.Username = os.Args[2]
			settings.Password = os.Args[3]
			settings.Id = os.Args[4]

			saveSettings(settings)
			
			doWork()

		case "start":
			err = s.Start()
			if err != nil {
				fmt.Printf("Failed to start: %s\n", err)
				return
			}
			fmt.Printf("Service \"%s\" started.\n", displayName)
		case "stop":
			err = s.Stop()
			if err != nil {
				fmt.Printf("Failed to stop: %s\n", err)
				return
			}
			fmt.Printf("Service \"%s\" stopped.\n", displayName)
		case "help":
			fmt.Printf(helpCommands)
		}
		return
	} else {
		fmt.Printf(helpCommands)
	}

	err = s.Run(func() error {
		// start
		go doWork()
		return nil
	}, func() error {
		// stop
		stopWork()
		return nil
	})
	if err != nil {
		s.Error(err.Error())
	}
}

var exit = make(chan bool)

func doWork() {
	log.Info("Service is running with config file '%s'", settingsFilePath)

	first := true
	ticker := time.NewTicker(7 * time.Second)

    settings := readSettings()
	tickerUpdate := time.NewTicker(7 * time.Second)

	for {
		select {
		case <-tickerUpdate.C:
			if first == true {
				tickerUpdate = time.NewTicker(time.Duration(settings.Interval) * time.Minute)
				first = false
			}

		    settings := readSettings()
		    id, _ := strconv.Atoi(settings.Id)

		    if id > 0 {
				ip := getExternalIp()

				if ip != settings.Ip {
					// update ip
					settings.Ip = ip

					result := updateIp(settings.Username, settings.Password, settings.Id, settings.Ip)
					if result == "success" {
						log.Info("Updated ID %s with IP %s was success", settings.Id, settings.Ip)
					} else {										
						log.Error("Updated ID %s with IP %s returned %s", settings.Id, settings.Ip, result)
					}
				}
			} else {
				log.Warning("Can't update settings are empty on '%s'", settingsFilePath)
			}

			saveSettings(settings)
		case <-ticker.C:
		case <-exit:
			ticker.Stop()
			tickerUpdate.Stop()
			return
		}
	}
}

func stopWork() {
	log.Info("Service is stopping")
	exit <- true
}

func readSettings() Settings {	
	settings := Settings{}
	settings.Interval = 3
	settings.Username = ""
	settings.Password = ""
	settings.Id = "0"
	settings.Ip = ""

	if _, err := os.Stat(settingsFilePath); err == nil {
		file, err := ioutil.ReadFile(settingsFilePath)
		if err != nil {
			log.Error("readSettings : ioutil.ReadFile: %v", err)
		}    

	    err = json.Unmarshal(file, &settings)
	    if err != nil {
			log.Error("readSettings : json.Unmarshal: %v", err)
		}
	} else {
		log.Warning("Creating default config file on '%s'", settingsFilePath)
	}

	return settings
}

func saveSettings(settings Settings) {
	outfile, err := json.MarshalIndent(settings, "", "  ")
    if err != nil {
		log.Error("saveSettings : json.Marshal: %v", err)
	}

    err = ioutil.WriteFile(settingsFilePath, outfile, 0644)
    if err != nil {
        log.Error("saveSettings : ioutil.WriteFile: %v", err)
    }
}

func updateIp(username string, password string, id string, ip string) string {
	url := fmt.Sprintf(updateIpFmt, username, password, id, ip)

	response, err := http.Get(url)

	if err != nil {
		log.Error("updateIp : get: %v", err)
        return string("")
    } else {
        defer response.Body.Close()
        contents, err := ioutil.ReadAll(response.Body)
        if err != nil {
			log.Error("updateIp : readBody: %v", err)
        }
        return string(contents)
    }
}

func getExternalIp() string {
    response, err := http.Get(getIpFmt)
    if err != nil {    	
		log.Error("getExternalIp : get: %v", err)
        return string("")
    } else {
        defer response.Body.Close()
        contents, err := ioutil.ReadAll(response.Body)
        if err != nil {
			log.Error("getExternalIp : readBody: %v", err)
            return string("")
        }
        return string(contents)
    }
}
