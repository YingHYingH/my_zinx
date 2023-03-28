package ziface

// 将请求的消息封装到一个message中，定义抽象接口

type IMessage interface {
	GetMsgID() uint32
	GetDataLen() uint32
	GetData() []byte

	SetMsgID(uint32)
	SetData([]byte)
	SetDataLen(uint32)
}
