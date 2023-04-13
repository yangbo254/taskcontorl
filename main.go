package main

import (
	"main/node"
	"main/router"
	"os"
)

func main() {
	taskMode := os.Getenv("TASKMODE")
	if taskMode == "node" {
		node.NewNodeTaskMgr().Run()
	} else {
		router.RunRouter()
	}
}
