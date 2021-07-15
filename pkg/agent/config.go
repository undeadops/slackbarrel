package agent

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

// AppConfig
type AppConfig struct {
	Package []Package `yaml:"package" json:"package"`
	File    []File    `yaml:"files" json:"files"`
	Service []Service `yaml:"service" json:"service"`
}

type Package struct {
	Name  string `yaml:"name" json:"name" required`
	State string `yaml:"state" json:"state" required`
}

type Service struct {
	Name   string `yaml:"name" json:"name" required`
	Action string `yaml:"action" json:"action" required`
}

// File
// Mode is string because Zero Padding Possible issues
type File struct {
	Source string `yaml:"src" json:"src" omitempty`
	Path   string `yaml:"path" json:"path" required`
	Owner  string `yaml:"owner" json:"owner" omitempty`
	Group  string `yaml:"group" json:"group" omitempty`
	Mode   string `yaml:"mode" json:"mode" omitempty`
	State  string `yaml:"state" json:"state"`
	ShaSum string `json:"shasum" omitempty`
}

func (a *Agent) ParseAppConfig(c io.Reader) *AppConfig {
	app := &AppConfig{}

	err := json.NewDecoder(c).Decode(app)
	if err != nil {
		a.Logger.Printf("Unable to unmarshall YAML Appconfig: #%v", err)
	}

	return app
}

func SetupClient() *http.Client {
	return &http.Client{
		Timeout: time.Second * 2,
	}
}
