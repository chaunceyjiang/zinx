package znet

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"net"
	"testing"
	"time"
)

func TestNewDataPack(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("server listen err:", err)
		return
	}

	//创建服务器goroutine，负责从客户端goroutine读取粘包的数据，然后进行解析
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("server accept err:", err)
			}

			//处理客户端请求
			go func(conn net.Conn) {
				//创建封包拆包对象dp
				dp := NewDataPack()
				for {
					//1 先读出流中的head部分
					headData := make([]byte, dp.GetHeadLen())
					_, err := io.ReadFull(conn, headData) //ReadFull 会把msg填充满为止
					if err != nil {
						fmt.Println("read head error")
						break
					}
					//将headData字节流 拆包到msg中
					msgHead, err := dp.Unpack(headData)
					if err != nil {
						assert.Fail(t, "server unpack err:", err)
						return
					}

					if msgHead.GetDataLen() > 0 {
						//msg 是有data数据的，需要再次读取data数据
						msg := msgHead.(*Message)
						msg.Data = make([]byte, msg.GetDataLen())

						//根据dataLen从io中读取字节流
						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							assert.Fail(t, "server unpack data err:", err)
						}
						switch msg.GetMsgId() {
						case 1:
							assert.Equal(t, "data pack test", string(msg.GetData()))
						case 2:
							assert.Equal(t, "NewMsgPackage test", string(msg.GetData()))
						}

					}
				}
			}(conn)
		}
	}()

	time.Sleep(1 * time.Second)
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client dial err:", err)
		return
	}

	//创建一个封包对象 dp
	dp := NewDataPack()

	//封装一个msg1包
	msg1 := NewMsgPackage(1, []byte("data pack test"))

	sendData1, err := dp.Pack(msg1)
	if err != nil {
		assert.Fail(t, "client pack msg1 err:", err)
	}

	msg2 := NewMsgPackage(2, []byte("NewMsgPackage test"))
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		assert.Fail(t, "client pack msg2 err:", err)
	}

	//将sendData1，和 sendData2 拼接一起，组成粘包
	sendData1 = append(sendData1, sendData2...)

	//向服务器端写数据
	conn.Write(sendData1)
	time.Sleep(1 * time.Second)
}
