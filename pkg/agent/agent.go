package agent

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

// Agent - Agent config
type Agent struct {
	Reload    bool     `omitempty`
	ServerUrl string   `yaml:"server_url" omitempty`
	Interval  int      `yaml:"interval" omitempty`
	Apps      []string `yaml:"apps" required`
	Logger    *log.Logger
}

func NewAgent(file string, serverurl string, inter int) *Agent {
	a := &Agent{}
	a.Logger = log.New(os.Stdout, "", log.LstdFlags)

	f, err := ioutil.ReadFile(file)
	if err != nil {
		a.Logger.Printf("Unable to read YAML file: #%v", err)
	}

	err = yaml.Unmarshal(f, a)
	if err != nil {
		a.Logger.Printf("Unable to unmarshall YAML file: #%v", err)
	}

	a.Reload = false

	if a.ServerUrl == "" {
		a.ServerUrl = serverurl
	}

	if a.Interval == 0 {
		a.Interval = inter
	}
	return a
}

func (a *Agent) CheckIn(interrupt chan os.Signal) {
	client := &http.Client{
		Timeout: time.Second * 2,
	}
	a.Logger.Println("Running CheckIn..")
	for _, app := range a.Apps {
		// Setup Server Request
		s := a.ServerUrl + "/config/" + app
		a.Logger.Printf("Server URL: %s\n", s)
		req, err := http.NewRequest("GET", s, nil)
		if err != nil {
			a.Logger.Println(err)
		}
		// Set HTTP Client Header
		req.Header.Set("User-Agent", "Barrel-Agent/0.0.1")
		resp, err := client.Do(req)
		if err != nil {
			a.Logger.Println(err)
		}
		if resp.Body != nil {
			// Defer request body close
			defer resp.Body.Close()
		}

		// Process Returned JSON body
		appConfig := a.ParseAppConfig(resp.Body)

		err = a.Packages(appConfig.Package)
		if err != nil {
			a.Logger.Println(err)
		}

		err = a.Files(appConfig.File)
		if err != nil {
			a.Logger.Println(err)
		}

		if a.Reload {
			err = a.ReloadService(appConfig.Service)
			if err != nil {
				a.Logger.Println(err)
			}
		}
	}
}
