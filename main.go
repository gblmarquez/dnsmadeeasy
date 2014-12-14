package main
 
import (
    "fmt"
    "time"
    "os"
    "os/user"

    "encoding/json"
    "net/http"
    "io/ioutil"

	"bitbucket.org/kardianos/service"
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

var log service.Logger

func main() {
	var name = "dyndns-dnsmadeeasy"
	var displayName = "Dynamic DNS updater for DnsMadeEasy"
	var desc = "Dynamic DNS updater for DnsMadeEasy"

	var s, err = service.NewService(name, displayName, desc)
	log = s

	if err != nil {
		fmt.Printf("%s unable to start: %s", displayName, err)
		return
	}

	if len(os.Args) > 1 {
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
		case "remove":
			err = s.Remove()
			if err != nil {
				fmt.Printf("Failed to remove: %s\n", err)
				return
			}
			fmt.Printf("Service \"%s\" removed.\n", displayName)
		case "run":
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
		}
		return
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

var exit = make(chan struct{})

func doWork() {
	log.Info("DnsMadeEasy updater is running")
	fmt.Printf("DnsMadeEasy updater is running")
	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-ticker.C:

		    settings := readSettings()

		    if settings.Id > 0 {
				settings.Ip = getExternalIp()

				result := updateIp(settings.Username, settings.Password, settings.Id, settings.Ip)
				fmt.Printf("%s\n", result)
			}

			saveSettings(settings)

			// sleep based on interval settings
			time.Sleep(time.Duration(settings.Interval) * time.Minute)
		case <-exit:
			ticker.Stop()
			return
		}
	}
}

func stopWork() {
	log.Info("DnsMadeEasy updater is stopping!")
	exit <- struct{}{}
}

func readSettings() Settings {
	usr, err := user.Current()
    if err != nil {
        log.Error("user.Current: %v", err)
    }
	
	settings := Settings{}
	settings.Interval = 5
	settings.Username = ""
	settings.Password = ""
	settings.Id = 0
	settings.Ip = ""

    filePath := fmt.Sprintf("%s/.dnsmadeeasy", usr.HomeDir)
	if _, err := os.Stat(filePath); err == nil {
		file, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.Error("ioutil.ReadFile: %v", err)
		}    

	    err = json.Unmarshal(file, &settings)
	    if err != nil {
			log.Error("json.Unmarshal: %v", err)
		}
	}

	return settings
}

func saveSettings(settings Settings) {
	outfile, err := json.MarshalIndent(settings, "", "  ")
    if err != nil {
		log.Error("json.Marshal: %v", err)
	}

	usr, err := user.Current()
    if err != nil {
        log.Error("user.Current: %v", err)
    }

    err = ioutil.WriteFile(fmt.Sprintf("%s/.dnsmadeeasy", usr.HomeDir), outfile, 0644)
    if err != nil {
        log.Error("ioutil.WriteFile: %v", err)
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