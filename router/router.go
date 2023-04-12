package router

import (
	"fmt"
	"log"
	"main/config"
	"main/cryptoutil"
	"main/model"
	"main/task"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type RouterEngine struct {
}

var route *gin.Engine

func RunRouter() {
	handleRouter := &RouterEngine{}
	router := gin.Default()
	router.MaxMultipartMemory = 8 << 22 // 32 MiB
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	apiBackend := router.Group("/api/node")
	{
		apiBackend.GET("/gettask", handleRouter.handleBackendGetTask)
	}
	apiWeb := router.Group("/api/web")
	{
		apiWeb.GET("/listtasks", handleRouter.handleWebListTask)

		apiWeb.GET("/addtask", handleRouter.handleWebAddTask)
		apiWeb.GET("/deltask", handleRouter.handleWebDelTask)

		apiWeb.POST("/uploadfile", handleRouter.handleWebUpdateFile)
		apiWeb.GET("/listfiles", handleRouter.handleWebListFiles)
	}

	route.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func (router *RouterEngine) returnError(c *gin.Context, code int, msg string, obj interface{}) {
	c.JSON(http.StatusBadRequest, NewMessageUtil().BuildErrorMessage(code, msg, obj))
}

func (router *RouterEngine) returnMessage(c *gin.Context, obj interface{}) {
	c.JSON(http.StatusOK, NewMessageUtil().BuildMessage(obj))
}

func (router *RouterEngine) handleWebUpdateFile(c *gin.Context) {
	// 单文件
	file, _ := c.FormFile("file")
	log.Println(file.Filename)
	fileName := fmt.Sprintf("%v-%s", time.Now().Unix(), cryptoutil.SHA1(file.Filename))
	filePath := config.NewConfig(config.SRVMODE).Server.UploadFilePath + fileName
	// 上传文件至指定的完整文件路径
	err := c.SaveUploadedFile(file, filePath)
	if err != nil {
		router.returnError(c, 1, err.Error(), nil)
		return
	}
	task.NewTask().AddFileInfo(fileName, file.Filename, filePath, file.Size)
	router.returnMessage(c, nil)
}

func (router *RouterEngine) handleWebListFiles(c *gin.Context) {
	files, err := task.NewTask().ListFileInfo()
	if err != nil {
		router.returnError(c, 1, err.Error(), nil)
		return
	}
	router.returnMessage(c, files)
}

func (router *RouterEngine) handleWebAddTask(c *gin.Context) {
}

func (router *RouterEngine) handleWebDelTask(c *gin.Context) {
}

func (router *RouterEngine) handleWebListTask(c *gin.Context) {
	start := c.Query("start")
	end := c.Query("end")
	if start == "" || end == "" {
		router.returnError(c, 1, "query param error", nil)
		return
	}
	startTime, err := strconv.ParseInt(start, 10, 64)
	if err != nil {
		router.returnError(c, 2, err.Error(), nil)
		return
	}
	endTime, err := strconv.ParseInt(end, 10, 64)
	if err != nil {
		router.returnError(c, 3, err.Error(), nil)
		return
	}
	tasks, err := task.NewTask().GetAllTask(startTime, endTime)
	if err != nil {
		router.returnError(c, 4, err.Error(), nil)
		return
	}
	router.returnMessage(c, tasks)
}

func (router *RouterEngine) handleBackendGetTask(c *gin.Context) {
	nodeId := c.Query("nodeid")
	if nodeId == "" {
		router.returnError(c, 1, "query param error", nil)
		return
	}
	tasks, err := task.NewTask().GetPlannedTask(nodeId)
	if err != nil {
		router.returnError(c, 2, err.Error(), nil)
		return
	}

	var result []model.TaskInfo
	now := time.Now().Unix()
	for _, v := range tasks {
		if now+60*1 >= v.Start {
			result = append(result, v)
		}
	}
	router.returnMessage(c, result)
}
