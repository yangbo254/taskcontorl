package node

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"main/config"
	"main/model"
	"net/http"
	"os"
)

type NodeUtil struct {
	baseUrl string
}

func NewNodeUtilDefault() *NodeUtil {
	conf := config.NewConfig(config.NODEMODE)
	return NewNodeUtil(conf.Client.ServerUrl)
}

func NewNodeUtil(url string) *NodeUtil {
	return &NodeUtil{baseUrl: url}
}

func (node *NodeUtil) GetFile(files []model.TaskVolumeStruct) (map[string]string, error) {
	volumes := make(map[string]string)
	for _, v := range files {
		filename, err := node.downloadFile(v.ServerPath)
		if err != nil {
			return nil, err
		}
		volumes[filename] = v.DestPath
	}
	return volumes, nil
}

func (node *NodeUtil) GetTask() ([]model.TaskInfo, error) {
	conf := config.NewConfig(config.NODEMODE)
	url := node.baseUrl + "/api/node/gettask"
	url += "?nodeid=" + conf.Client.NodeId
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bytes, _ := io.ReadAll(resp.Body)
	var taskinfo []model.TaskInfo
	err = json.Unmarshal(bytes, &taskinfo)
	if err != nil {
		return nil, err
	}
	return taskinfo, nil
}

func (node *NodeUtil) GetKillTask() ([]uint, error) {
	conf := config.NewConfig(config.NODEMODE)
	url := node.baseUrl + "/api/node/getkill"
	url += "?nodeid=" + conf.Client.NodeId
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bytes, _ := io.ReadAll(resp.Body)
	var taskinfo []uint
	err = json.Unmarshal(bytes, &taskinfo)
	if err != nil {
		return nil, err
	}
	return taskinfo, nil
}

// 上报任务状态
func (node *NodeUtil) ReportContainerStats(obj interface{}) error {
	conf := config.NewConfig(config.NODEMODE)
	url := node.baseUrl + "/report/container/stat"
	url += "?nodeid=" + conf.Client.NodeId
	bytesData, _ := json.Marshal(obj)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(bytesData))
	req.Header.Set("X-Report-Ver", "report_container_stat_v1")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fmt.Println("report system info response status:", resp.Status)
	fmt.Println("report system info response headers:", resp.Header)
	body, _ := io.ReadAll(resp.Body)
	fmt.Println("report system info response body:", string(body))
	return nil
}

// 上报任务状态
func (node *NodeUtil) ReportContainerEnd(obj interface{}) error {
	conf := config.NewConfig(config.NODEMODE)
	url := node.baseUrl + "/report/container/end"
	url += "?nodeid=" + conf.Client.NodeId
	bytesData, _ := json.Marshal(obj)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(bytesData))
	req.Header.Set("X-Report-Ver", "report_container_end_v1")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fmt.Println("report container end info response status:", resp.Status)
	fmt.Println("report container end info response headers:", resp.Header)
	body, _ := io.ReadAll(resp.Body)
	fmt.Println("report container end info response body:", string(body))
	return nil
}

// 上报系统信息
func (node *NodeUtil) ReportSystemInfo() error {
	systeminfo := &SystemInfo{}
	systeminfo.GetInfo()
	conf := config.NewConfig(config.NODEMODE)
	url := node.baseUrl + "/report/system"
	url += "?nodeid=" + conf.Client.NodeId
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer([]byte(systeminfo.String())))
	req.Header.Set("X-Report-Ver", "report_sys_v1")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fmt.Println("report system info response status:", resp.Status)
	fmt.Println("report system info response headers:", resp.Header)
	body, _ := io.ReadAll(resp.Body)
	fmt.Println("report system info response body:", string(body))
	return nil
}

func (node *NodeUtil) downloadFile(filename string) (string, error) {
	url := node.baseUrl + "/api/node/downloadfile"
	url += "?fileid=" + filename
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	filerandname := "/tmp/" + generateRandomFilename()
	file, err := os.Create(filerandname)
	if err != nil {
		return "", err
	}
	defer file.Close()
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", err
	}
	return filerandname, nil
}

func generateRandomFilename() string {
	// 生成16字节的随机字节序列
	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}

	// 将字节序列编码为16进制字符串
	randomHex := hex.EncodeToString(randomBytes)

	return randomHex
}
