package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

func TestDataPack(t *testing.T) {
	/*
		模拟的服务器
	*/
	// 创建socketTCP
	listener, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("server listen err: ", err)
		return
	}
	// 从客户端读取数据，拆包处理
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("server accept err: ", err)
				continue
			}
			go func(conn net.Conn) {
				// 处理客户端请求
				dp := NewDataPack()
				for {
					headData := make([]byte, dp.GetHeadLen())
					_, err = io.ReadFull(conn, headData)
					if err != nil {
						fmt.Println("read head err: ", err)
						break
					}
					msgHead, err := dp.UnPack(headData)
					if err != nil {
						fmt.Println("head unPack err: ", err)
						return
					}
					if msgHead.GetDataLen() > 0 {
						msg := msgHead.(*Message)
						msg.Data = make([]byte, msg.GetDataLen())

						_, err = io.ReadFull(conn, msg.Data)
						if err != nil {
							fmt.Println("sever unPack data err: ", err)
							return
						}
						fmt.Println("----> Recv MsgID: ", msg.ID, " dataLen: ", msg.DataLen, " data: ", string(msg.Data))
					}

				}
			}(conn)
		}
	}()
	/*
		模拟客户端
	*/

	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client dial err: ", err)
		return
	}
	dp := NewDataPack()
	// 封装两个msg一起发送，模拟粘包
	msg1 := &Message{
		ID:      1,
		DataLen: 4,
		Data:    []byte{'a', 'b', 'b', 'c'},
	}
	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("client pack msg1 err: ", err)
		return
	}
	msg2 := &Message{
		ID:      2,
		DataLen: 7,
		Data:    []byte{'c', 'a', 'b', 'c', 'a', 'b', 'c'},
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("client pack msg2 err: ", err)
		return
	}
	sendData1 = append(sendData1, sendData2...)
	conn.Write(sendData1)

	select {}
}
