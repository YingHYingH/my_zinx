package znet

import (
	"errors"
	"fmt"
	"io"
	"my_zinx/utils"
	"my_zinx/ziface"
	"net"
	"sync"
)

// 连接模块
type Connection struct {
	// 当前Conn属于哪个Server
	TcpServer ziface.IServer

	// 当前连接的socket TCP套接字
	Conn *net.TCPConn

	// 连接的ID
	ConnID uint32

	// 当前的连接状态
	isClosed bool

	// 告知当前连接已经退出的channel
	ExitChan chan bool

	// 无缓冲管道，用于读写Goroutine之间的消息通信
	msgChan chan []byte

	// 该连接处理的方法Router
	MsgHandler ziface.IMsgHandler

	// 连接属性集合
	property map[string]interface{}

	// 保护连接属性的锁
	propertyLock sync.RWMutex
}

// 初始化连接模块
func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, handler ziface.IMsgHandler) *Connection {
	c := &Connection{
		TcpServer:  server,
		Conn:       conn,
		ConnID:     connID,
		isClosed:   false,
		MsgHandler: handler,
		msgChan:    make(chan []byte),
		ExitChan:   make(chan bool, 1),
		property:   map[string]interface{}{},
	}
	// 将conn加入到connManager
	c.TcpServer.GetConnManager().Add(c)
	return c
}

func (c *Connection) StartWriter() {
	fmt.Println("Writer Goroutine is running...")
	defer fmt.Println("connID = ", c.ConnID, " Writer is exit, remote addr is ", c.RemoteAddr().String())

	for {
		select {
		case data := <-c.msgChan:
			// 有数据要写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("send data error: ", err)
				return
			}
		case <-c.ExitChan:
			// 代表reader已经退出，此时writer也要退出
			return
		}
	}
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
		if utils.GlobalObject.WorkerPoolSize > 0 {
			// 开启了工作池机制，将消息发送给worker工作池处理即可
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			go c.MsgHandler.DoMsgHandler(&req)
		}
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
	c.msgChan <- binaryMsg
	return nil
}

func (c *Connection) Start() {
	fmt.Println("Conn Start()...ConnID = ", c.ConnID)
	// 启动当前连接读数据的业务
	go c.StartReader()
	go c.StartWriter()
	// 创建连接之后的hook方法
	c.TcpServer.CallOnConnStart(c)
}

func (c *Connection) Stop() {
	fmt.Println("Conn Stop...ConnID = ", c.ConnID)
	if c.isClosed {
		return
	}
	c.isClosed = true
	// 关闭连接之前的hook方法
	c.TcpServer.CallOnConnStop(c)
	// 关闭socket连接
	c.Conn.Close()
	// 告知writer关闭
	c.ExitChan <- true
	// 将当前连接从connManager剔除
	c.TcpServer.GetConnManager().Remove(c)
	// 回收资源
	close(c.ExitChan)
	close(c.msgChan)
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

func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	c.property[key] = value
}

func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	if value, ok := c.property[key]; ok {
		return value, nil
	}
	return nil, errors.New("no property found")
}

func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	delete(c.property, key)
}
