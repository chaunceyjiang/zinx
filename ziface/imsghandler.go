package ziface

type ImsgHandle interface {
	DoMsgHandler(request IRequest) // 调用自定义handler
	AddRouter(msgId uint32,router IRouter) // 添加handler
	StartWorkerPool() // 启动worker工作池
	SendMsgToTaskQueue(request IRequest) // 将消息交给taskQueue,由worker处理
}
