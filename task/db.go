package task

import (
	"errors"
	"main/model"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type taskDbInfo struct {
	db *gorm.DB
}

func (task *Task) AutoInit() (err error) {
	task.db, err = gorm.Open(sqlite.Open("main.db"), &gorm.Config{})
	if err != nil {
		return errors.New("failed to connect database")
	}

	task.db.AutoMigrate(&model.TaskInfo{})
	task.db.AutoMigrate(&model.FileInfo{})
	task.db.AutoMigrate(&model.NodeInfo{})
	task.db.AutoMigrate(&model.Group{})

	return nil
}
