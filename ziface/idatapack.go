package ziface

/*
	封包拆包模块，面向TCP连接的数据流，用于处理TCP粘包问题
*/

type IDataPack interface {
	// 获取包头长度的方法
	GetHeadLen() uint32
	// 封包方法
	Pack(msg IMessage) ([]byte, error)
	// 拆包方法
	UnPack([]byte) (IMessage, error)
}
