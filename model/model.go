package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type NodeTaskInfo struct {
	TaskInfo
	RealStatus          string            `json:"realStatus"`
	RealStart           time.Time         `json:"realStart"`
	RealEnd             time.Time         `json:"realEnd"`
	RealNodeVolumes     map[string]string `json:"realVolumes"`
	RealNodeContainerId string            `json:"realContainerId"`
}

type TaskVolumeStruct struct {
	ServerPath string `json:"serverpath"`
	DestPath   string `json:"destpath"`
}

type TaskParameter struct {
	TaskImageName string             `json:"taskImageName"`
	TaskHostMode  string             `json:"taskHostMode"`
	TaskEnv       string             `json:"taskEnv"`
	TaskCmd       []string           `json:"taskCmd"`
	TaskVolume    []TaskVolumeStruct `json:"taskVolume"`
}

type TaskInfo struct {
	gorm.Model
	Name  string    `json:"name"`
	Start time.Time `json:"start"`
	Nodes NodeList  `json:"nodes"`
	TaskParameter
}

type NodeList []string

// gorm 自定义结构需要实现 Value Scan 两个方法
// Value 实现方法
func (p NodeList) Value() (driver.Value, error) {
	return json.Marshal(p)
}

// Scan 实现方法
func (p *NodeList) Scan(data interface{}) error {
	return json.Unmarshal(data.([]byte), &p)
}
