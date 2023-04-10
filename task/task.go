package task

import (
	"main/model"
	"time"
)

type Task struct {
	taskDbInfo
}

var gTask *Task = nil

func NewTask() *Task {
	if gTask == nil {
		gTask = &Task{}
		gTask.AutoInit()
	}
	return gTask
}

func (task *Task) AddTask(taskName string, start time.Time, nodes []string, taskParam model.TaskParameter) (err error) {
	taskInfo := model.TaskInfo{
		Name:          taskName,
		Start:         start,
		Nodes:         nodes,
		TaskParameter: taskParam,
	}
	result := task.db.Create(&taskInfo)
	return result.Error
}

func (task *Task) DelTask(taskId uint) {
	task.db.Delete(&model.TaskInfo{}, taskId)
}

func (task *Task) GetAllTask(StartBegin, StartEnd time.Time) (tasks []model.TaskInfo, err error) {
	result := task.db.Where("start >= ? and start <= ?", StartBegin, StartEnd).Find(&tasks)
	return tasks, result.Error
}

func (task *Task) GetPlannedTask(NodeID string) (tasks []model.TaskInfo, err error) {
	var allTasks []model.TaskInfo
	result := task.db.Where("Start >= ?", time.Now()).Order("id").Find(&allTasks)
	if result.Error != nil {
		return nil, result.Error
	}
	for _, v := range allTasks {
		for _, e := range v.Nodes {
			if e == NodeID {
				tasks = append(tasks, v)
				break
			}
		}
	}
	return tasks, nil
}
