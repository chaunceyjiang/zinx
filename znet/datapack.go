package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"zinx/utils"
	"zinx/ziface"
)

type DataPack struct {
}

func (d *DataPack) GetHeadLen() uint32 {
	//panic("implement me")
	return 8
}

func (d *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	//将一个msg对象 封装成byte
	// 格式: 数据长度+数据类型+数据
	databuff := bytes.NewBuffer([]byte{})

	if err := binary.Write(databuff, binary.LittleEndian, msg.GetDataLen()); err != nil {
		return nil, err
	}
	if err := binary.Write(databuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}

	if err := binary.Write(databuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}
	return databuff.Bytes(), nil

}
//Unpack 拆包 data 只包含了head 的数据长度,我们先解析head,然后在解析body
func (d *DataPack) Unpack(data []byte) (ziface.IMessage, error) {
	msg := &Message{
		Id:      0,
		DataLen: 0,
		Data:    nil,
	}
	databuff := bytes.NewReader(data)
	if err := binary.Read(databuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}
	if err := binary.Read(databuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}
	//if err := binary.Read(databuff, binary.LittleEndian, &msg.DataLen); err != nil {
	//	return nil, err
	//}
	if utils.GlobalObject.MaxPacketSize > 0 && msg.DataLen > utils.GlobalObject.MaxPacketSize {
		return nil, errors.New("Too large msg data received\n")
	}
	return msg, nil
}

func NewDataPack() *DataPack {
	return &DataPack{}

}
