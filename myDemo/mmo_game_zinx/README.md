目录结构  
/apis 存放基本的用户自定义路由业务，一个msgID对应一个业务  
/conf zinx.json 存放zinx配置文件  
/pb
* msg.ptoto 原始protobuf协议文件
* build.sh 编译msg.ptoto的脚本
* msg.pb.go 编译生成的go文件（只读）

/core 存放核心功能
/main.go 服务器主入口
