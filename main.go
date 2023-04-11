package main

import (
	"main/node"
	"os"
)

func main() {
	taskMode := os.Getenv("TASKMODE")
	if taskMode == "node" {
		node.NewNodeTaskMgr().Run()
	}
}
