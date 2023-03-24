package ziface

// 把客户端请求的连接信息和请求的数据，包装到一个Request中
type IRequest interface {
	// 得到当前连接
	GetConnection() IConnection
	// 得到请求的消息数据
	GetData() []byte
}
