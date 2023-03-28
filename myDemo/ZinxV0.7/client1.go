package main

import (
	"fmt"
	"io"
	"my_zinx/znet"
	"net"
	"time"
)

// 模拟客户端
func main() {
	fmt.Println("client1 start...")

	time.Sleep(1 * time.Second)
	// 连接远程服务器，得到一个conn连接
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client1 start error:", err)
		return
	}

	for {
		dp := znet.NewDataPack()
		binaryMsg, err := dp.Pack(znet.NewMessage(1, []byte("Zinx client1 Test Message")))
		if err != nil {
			fmt.Println("pack error:", err)
			break
		}
		if _, err = conn.Write(binaryMsg); err != nil {
			fmt.Println("write error:", err)
			break
		}

		binaryHead := make([]byte, dp.GetHeadLen())
		if _, err = io.ReadFull(conn, binaryHead); err != nil {
			fmt.Println("read head error:", err)
			break
		}
		msgHead, err := dp.UnPack(binaryHead)
		if err != nil {
			fmt.Println("client head unpack error:", err)
			break
		}

		if msgHead.GetDataLen() > 0 {
			msg := msgHead.(*znet.Message)
			msg.Data = make([]byte, msg.GetDataLen())
			if _, err = io.ReadFull(conn, msg.Data); err != nil {
				fmt.Println("read msg data error:", err)
				break
			}
			fmt.Println("----> Recv MsgID: ", msg.ID, " dataLen: ", msg.DataLen, " data: ", string(msg.Data))
		}
		time.Sleep(1 * time.Second)
	}
}
