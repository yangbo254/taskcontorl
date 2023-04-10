package node

import (
	"main/model"
	"sync"
	"time"
)

type NodeTaskMgr struct {
	waitlocker      *sync.Mutex
	waitTaskList    map[uint]model.NodeTaskInfo
	runningTaskList map[uint]model.NodeTaskInfo
	endTaskList     map[uint]model.NodeTaskInfo
}

type taskRunningInfo struct {
	ContainerId string
}

func (mgr *NodeTaskMgr) LoopTask() {
	// 遍历等待队列，修改status
	mgr.waitlocker.Lock()
	for _, v := range mgr.waitTaskList {
		if v.RealStatus == "new" {
			// 新加入队列的任务
			// 需要执行pull命令/下载mount文件等前置任务
			// to do

			// 切换状态为等待
			v.RealStatus = "wait"
		} else if v.RealStatus == "wait" {
			if time.Now().After(v.Start) {
				// 切换任务到运行队列
				// to do

				// 从队列中移除
				delete(mgr.waitTaskList, v.ID)

			}
		}
	}
	mgr.waitlocker.Unlock()

	// 遍历运行队列，等待运行结束
	for _, v := range mgr.runningTaskList {
		// v未找到，移入结束队列
		// to do
	}

	// 遍历结束队列，上报结束状态
	for _, v := range mgr.endTaskList {
		// 上报v状态，移除自己
		// to do
	}
}
