package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type ClientConfig struct {
	ServerUrl      string `yaml:"url"`
	TryGetTaskTime int    `yaml:"gettasktime"`
	LoopTime       int    `yaml:"looptasktime"`
	KillTaskTime   int    `yaml:"killtasktime"`
}

type ServerConfig struct {
	UploadFilePath string `yaml:"uploadpath"`
}

type Set struct {
	Server ServerConfig `yaml:"server"`
	Client ClientConfig `yaml:"client"`
	init   bool
}

var Setting = Set{}

func NewConfig() *Set {
	if !Setting.init {
		file, err := os.ReadFile("./conf/app.yml")
		if err != nil {
			log.Fatal("fail to read file:", err)
		}

		err = yaml.Unmarshal(file, &Setting)
		if err != nil {
			log.Fatal("fail to yaml unmarshal:", err)
		}
		Setting.init = true
	}
	return &Setting
}
