package znet

import (
	"errors"
	"fmt"
	"io"
	"my_zinx/ziface"
	"net"
)

// 连接模块
type Connection struct {
	// 当前连接的socket TCP套接字
	Conn *net.TCPConn

	// 连接的ID
	ConnID uint32

	// 当前的连接状态
	isClosed bool

	// 告知当前连接已经退出的channel
	ExitChan chan bool

	// 该连接处理的方法Router
	MsgHandler ziface.IMsgHandler
}

// 初始化连接模块
func NewConnection(conn *net.TCPConn, connID uint32, handler ziface.IMsgHandler) *Connection {
	c := &Connection{
		Conn:       conn,
		ConnID:     connID,
		isClosed:   false,
		MsgHandler: handler,
		ExitChan:   make(chan bool, 1),
	}
	return c
}

func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running...")
	defer fmt.Println("connID = ", c.ConnID, " Reader is exit, remote addr is ", c.RemoteAddr().String())
	defer c.Stop()

	for {
		// 创建拆包解包的对象
		dp := NewDataPack()
		// 读取客户端的msg head
		headData := make([]byte, dp.GetHeadLen())
		_, err := io.ReadFull(c.GetTCPConnection(), headData)
		if err != nil {
			fmt.Println("read msg head err: ", err)
			break
		}
		// 拆包，得到msgID和msgDataLen，放在msg消息中
		msg, err := dp.UnPack(headData)
		if err != nil {
			fmt.Println("unpack err: ", err)
			break
		}
		// 根据dataLen再次读取Dat|a，放在msg.Data
		var data []byte
		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())
			_, err = io.ReadFull(c.GetTCPConnection(), data)
			if err != nil {
				fmt.Println("read msg data err: ", err)
				break
			}
		}
		msg.SetData(data)
		req := Request{
			conn: c,
			msg:  msg,
		}
		go c.MsgHandler.DoMsgHandler(&req)
	}
}

// 提供一个SendMsg方法，将我们要发送给客户端的数据，先封包，再发送
func (c *Connection) SendMsg(msgID uint32, data []byte) error {
	if c.isClosed {
		return errors.New("connection closed when send msg")
	}
	// 将data进行封包 MsgDataLen/MsgID/Data
	dp := NewDataPack()
	binaryMsg, err := dp.Pack(NewMessage(msgID, data))
	if err != nil {
		fmt.Println("pack error msgID: ", msgID)
		return errors.New("pack msg error")
	}
	if _, err = c.Conn.Write(binaryMsg); err != nil {
		fmt.Println("write error msgID: ", msgID)
		return errors.New("write msg error")
	}
	return nil
}

func (c *Connection) Start() {
	fmt.Println("Conn Start()...ConnID = ", c.ConnID)
	// 启动当前连接读数据的业务
	go c.StartReader()
}

func (c *Connection) Stop() {
	fmt.Println("Conn Stop...ConnID = ", c.ConnID)
	if c.isClosed {
		return
	}
	c.isClosed = true

	// 关闭socket连接
	c.Conn.Close()
	// 回收资源
	close(c.ExitChan)
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}
