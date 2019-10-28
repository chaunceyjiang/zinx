package ziface

type IMessage interface {
	// 获取消息数据段的长度
	GetDataLen() uint32
	// 获取消息ID
	GetMsgId() uint32
	// 获取消息内容
	GetData() []byte


	SetMesId(uint32)

	SetData([]byte)

	SetDataLen(uint32)
}



