package config

import (
	"log"
	"os"

	"github.com/gofrs/uuid"
	"gopkg.in/yaml.v2"
)

type RunningMode int

const (
	NODEMODE RunningMode = 0
	SRVMODE  RunningMode = 1
)

type ClientConfig struct {
	ServerUrl      string `yaml:"url"`
	TryGetTaskTime int64  `yaml:"getTaskTime"`
	LoopTime       int64  `yaml:"loopTaskTime"`
	KillTaskTime   int64  `yaml:"killTaskTime"`
	NodeId         string `yaml:"nodeId"`
}

type ServerConfig struct {
	UploadFilePath string `yaml:"uploadPath"`
}

type Set struct {
	Server ServerConfig `yaml:"server"`
	Client ClientConfig `yaml:"client"`
	init   bool         `yaml:"-"`
}

var Setting = Set{}

func NewConfig(mode RunningMode) *Set {
	if !Setting.init {
		filePath := "./conf/app.yml"
		file, err := os.ReadFile(filePath)
		if err != nil {
			log.Fatal("fail to read file:", err)
		}

		err = yaml.Unmarshal(file, &Setting)
		if err != nil {
			log.Fatal("fail to yaml unmarshal:", err)
		}
		if mode == NODEMODE && Setting.Client.NodeId == "" {
			uuidValue, err := uuid.NewV4()
			if err != nil {
				log.Fatalf("failed to generate UUID: %v", err)
				return &Setting
			}
			Setting.Client.NodeId = uuidValue.String()
			byteData, _ := yaml.Marshal(Setting)
			_ = os.WriteFile(filePath, byteData, os.ModePerm)
			return &Setting
		}

		if mode == SRVMODE && Setting.Server.UploadFilePath == "" {
			Setting.Server.UploadFilePath = "./upload/"
			byteData, _ := yaml.Marshal(Setting)
			_ = os.WriteFile(filePath, byteData, os.ModePerm)
			return &Setting
		}
		Setting.init = true
	}
	return &Setting
}
