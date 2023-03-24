package main

import (
	"fmt"
	"net"
	"time"
)

// 模拟客户端
func main() {
	fmt.Println("client start...")

	time.Sleep(1 * time.Second)
	// 连接远程服务器，得到一个conn连接
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client start error:", err)
		return
	}

	for {
		// 连接调用write写数据
		_, err := conn.Write([]byte("hello world..."))
		if err != nil {
			fmt.Println("write conn error:", err)
			return
		}
		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read buf err:", err)
			return
		}
		fmt.Printf("server call back:%s, cnt = %d\n", buf, cnt)
		time.Sleep(1 * time.Second)
	}
}
