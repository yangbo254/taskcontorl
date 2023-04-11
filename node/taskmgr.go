package node

import (
	"log"
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

func (mgr *NodeTaskMgr) LoopTask() {
	// 遍历等待队列，修改status
	mgr.waitlocker.Lock()
	for _, v := range mgr.waitTaskList {
		if v.RealStatus == "new" {
			// 新加入队列的任务
			// 需要执行pull命令/下载mount文件等前置任务
			// 1.pull镜像
			_, err := NewDockerCmd().PullImage(v.TaskImageName)
			if err != nil {
				log.Fatalf("pull docker images failed. %v", err)
				delete(mgr.waitTaskList, v.ID)
			}
			// 2.下载生成mounts
			volumes, err := NewNodeUtilDefault().GetFile(v.TaskVolume)
			if err != nil {
				log.Fatalf("download volume files failed. %v", err)
				delete(mgr.waitTaskList, v.ID)
			}
			v.RealNodeVolumes = volumes
			// 切换状态为等待
			v.RealStatus = "wait"
		} else if v.RealStatus == "wait" {
			if time.Now().After(v.Start) {
				// 切换任务到运行队列
				containerId, err := NewDockerCmd().CreateContainer(v.TaskImageName, v.TaskCmd, v.RealNodeVolumes)
				if err != nil {
					log.Fatalf("run docker failed. %v", err)
					delete(mgr.waitTaskList, v.ID)
				}
				v.RealNodeContainerId = containerId
				v.RealStart = time.Now()
				// swap map
				mgr.runningTaskList[v.ID] = v
				delete(mgr.waitTaskList, v.ID)
			}
		}
	}
	mgr.waitlocker.Unlock()

	containers, err := NewDockerCmd().ListContainer()
	if err != nil {
		return
	}
	containerMap := NewDockerCmd().ConvContainersToMap(containers)

	// 遍历运行队列，等待运行结束
	for _, v := range mgr.runningTaskList {
		if vv, found := containerMap[v.RealNodeContainerId]; found {
			v.RealStatus = vv.State
			if vv.State == "exited" {
				// 任务结束，移入结束队列
				NewDockerCmd().RmContainer(vv.ID)
				v.RealEnd = time.Now()
				mgr.endTaskList[v.ID] = v
				delete(mgr.runningTaskList, v.ID)
			} else {
				continue
			}
		} else {
			// v未找到，移入结束队列
			v.RealStatus = "not found"
			NewDockerCmd().RmContainer(v.RealNodeContainerId)
			v.RealEnd = time.Now()
			mgr.endTaskList[v.ID] = v
			delete(mgr.runningTaskList, v.ID)
		}

	}

	// 遍历结束队列，上报结束状态
	//for _, v := range mgr.endTaskList {
	// 上报v状态，移除自己
	// to do
	//}
}
