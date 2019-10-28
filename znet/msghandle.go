package znet

import (
	"fmt"
	"strconv"
	"zinx/ziface"
)

type MsgHandle struct {
	apis map[uint32]ziface.IRouter
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

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{apis: make(map[uint32]ziface.IRouter)}
}