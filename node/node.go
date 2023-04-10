package node

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"main/model"
	"net/http"
)

type Node struct {
	baseUrl string
}

func (node *Node) GetTask() error {
	url := node.baseUrl + "/node/gettask"
	url += "?nodeid=" + "111"
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	bytes, err := io.ReadAll(resp.Body)
	var taskinfo []model.TaskInfo
	err = json.Unmarshal(bytes, &taskinfo)
	if err != nil {
		return err
	}
	return nil
}

// 上报任务状态
func (node *Node) ReportTaskStats() error {
	return nil
}

// 上报系统信息
func (node *Node) ReportSystemInfo() error {
	systeminfo := &SystemInfo{}
	systeminfo.GetInfo()
	url := node.baseUrl + "/report/system"
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
