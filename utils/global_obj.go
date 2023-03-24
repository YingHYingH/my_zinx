package utils

import (
	"encoding/json"
	"io/ioutil"
	"my_zinx/ziface"
)

type GlobalObj struct {
	TcpServer ziface.IServer
	Host      string // 服务器主机监听的IP
	TcpPort   int    // 服务器主机监听的端口号
	Name      string // 服务器名称

	Version        string // Zinx的版本号
	MaxConn        int    // 服务器允许的最大连接数
	MaxPackageSize uint32 //数据包最大值
}

var GlobalObject *GlobalObj

func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("conf/zinx.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, g)
	if err != nil {
		panic(err)
	}
}

func init() {
	GlobalObject = &GlobalObj{
		TcpServer:      nil,
		Host:           "0.0.0.0", // 本地全IP地址
		TcpPort:        8999,
		Name:           "ZinxServerApp",
		Version:        "v0.4",
		MaxConn:        1000,
		MaxPackageSize: 4096,
	}
	GlobalObject.Reload()
}
