package model

import (
	"database/sql/driver"
	"encoding/json"

	"gorm.io/gorm"
)

type NodeTaskInfo struct {
	TaskInfo
	RealStatus          string            `json:"realStatus"`
	RealStart           int64             `json:"realStart"`
	RealEnd             int64             `json:"realEnd"`
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
	Name  string   `json:"name"`
	Start int64    `json:"start"`
	Nodes NodeList `json:"nodes"`
	TaskParameter
}

type FileInfo struct {
	gorm.Model
	CurrentName  string `json:"currentName"`
	OriginalName string `json:"originalName"`
	Path         string `json:"path"`
	FileSize     int64  `json:"size"`
}

type NodeInfo struct {
	gorm.Model
	NodeId   string `json:"nodeid"` //节点名称
	Name     string `json:"name"`   // 备注名称
	GroupIds []uint `json:"groupids"`
}

type Group struct {
	gorm.Model
	Name string `json:"name" gorm:"uniqueIndex"`
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
