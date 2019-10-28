package znet

import (
	"fmt"
	"strconv"
	"zinx/utils"
	"zinx/ziface"
)

type MsgHandle struct {
	apis           map[uint32]ziface.IRouter
	workerPoolSize uint32
	taskQueue      []chan ziface.IRequest // taskQueue 创建一个 chan 数组
}

func (m *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	router, ok := m.apis[request.GetMsgId()]
	if !ok {
		fmt.Println("api msgId = ", request.GetMsgId(), " is not FOUND!")
		return
	}
	router.PreHandle(request)
	router.Handle(request)
	router.PostHandle(request)
}

func (m *MsgHandle) AddRouter(msgId uint32, router ziface.IRouter) {

	if _, ok := m.apis[msgId]; ok {
		panic("repeated api , msgId = " + strconv.Itoa(int(msgId)))
	}
	m.apis[msgId] = router
	fmt.Println("Add api msgId = ", msgId)
}

func (m *MsgHandle) startOneWorker(workerId int, taskQueue chan ziface.IRequest) {
	fmt.Println("Worker ID = ", workerId, " is started.")
	for request := range taskQueue {
		m.DoMsgHandler(request)
	}
}

func (m *MsgHandle) StartWorkerPool() {
	for i := 0; i < int(m.workerPoolSize); i++ {
		// 创建一个长度为MaxWorkerTaskLen 的chan ,用来存储request
		m.taskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		go m.startOneWorker(i, m.taskQueue[i])
	}
}

func (m *MsgHandle) SendMsgToTaskQueue(request ziface.IRequest) {
	workerID := request.GetConnection().GetConnID() % m.workerPoolSize
	fmt.Println("Add ConnID=", request.GetConnection().GetConnID(), " request msgID=", request.GetMsgId(), "to workerID=", workerID)
	m.taskQueue[workerID] <- request
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		apis:           make(map[uint32]ziface.IRouter),
		workerPoolSize: utils.GlobalObject.WorkerPoolSize,
		taskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}
