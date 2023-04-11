# taskcontorl

## 架构
### web界面(web)
    UI展示
### 控制端后台(task)
#### 对web界面提供：  
- [x] 任务管理、任务新增、任务删除 接口  
- [ ] 节点机器状态、节点任务状态、强制结束节点任务 接口  
- [ ] 客户端上报状态展示 接口  
#### 对节点提供：  
- [ ] 任务下发 强制结束任务 接口  
- [ ] 节点机器状态上报 节点任务状态上报 接口  
#### 对客户端上报程序提供：
- [ ] 客户端上报状态 接口  
### 节点程序(node)
- [x] 提供节点机器状态监控上报  
- [x] 提供任务分发及执行  
- [x] 提供docker image pull 功能  
- [x] 提供docker container list 功能  
- [x] 提供docker container kill 功能  
### 客户端上报程序(container)
    和具体程序相关??
### 数据流
    web <==> task <==> node/container
