package configserver

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

// Config
type Config struct {
	Apps AppName `yaml:"apps" json:"apps"`
}

// AppName
type AppName map[string]App

// App
type App struct {
	Package []Package `yaml:"package" json:"package"`
	File    []File    `yaml:"files" json:"files"`
	Service []Service `yaml:"service" json:"service"`
}

// Package
type Package struct {
	Name  string `yaml:"name" json:"name" required`
	State string `yaml:"state" json:"state" required`
}

// Service
type Service struct {
	Name   string `yaml:"name" json:"name" required`
	Action string `yaml:"action" json:"action" required`
}

// Files
type File struct {
	Source string `yaml:"src" json:"src" omitempty`
	Path   string `yaml:"path" json:"path" required`
	Owner  string `yaml:"owner" json:"owner" omitempty`
	Group  string `yaml:"group" json:"group" omitempty`
	Mode   string `yaml:"mode" json:"mode" omitempty`
	Sate   string `yaml:"state" json:"state" required`
	ShaSum string `json:"shasum" omitempty`
}

// SetupConfigServer
func SetupConfigServer(file string, dataDir string) *Config {
	c := &Config{}

	f, err := ioutil.ReadFile(file)
	if err != nil {
		log.Printf("Unable to read YAML file: #%v", err)
	}

	err = yaml.Unmarshal(f, c)
	if err != nil {
		log.Printf("Unable to unmarshall YAML file: #%v", err)
	}
	for x, app := range c.Apps {
		for i, file := range app.File {
			// for each file listed.
			f := dataDir + "/" + file.Source
			fsum, err := ShaSumFile(f)
			if err != nil {
				log.Printf("Unable to shaSum file: %v", err)
			}
			c.Apps[x].File[i].ShaSum = fsum
		}
	}
	return c
}

func ShaSumFile(file string) (string, error) {
	// Open File, Determin Sha256 Hash
	// Used to know if file is different on remote end
	f, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
