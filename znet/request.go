package znet

import "my_zinx/ziface"

type Request struct {
	// 已经和客户端建立好的连接
	conn ziface.IConnection
	// 客户端请求的数据
	data []byte
}

func (r *Request) GetConnection() ziface.IConnection {
	return r.conn
}

func (r *Request) GetData() []byte {
	return r.data
}
