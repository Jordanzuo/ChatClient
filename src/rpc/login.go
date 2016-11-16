package rpc

import (
	"github.com/Jordanzuo/ChatServerModel/src/commandType"
)

var (
	// 存储登陆成功信息的通道
	loginSucceedCh = make(chan int)
)

func login() {
	// 定义请求的主体结构
	commandMap := make(map[string]interface{})
	commandMap["Id"] = id
	commandMap["Name"] = name
	commandMap["UnionId"] = unionId
	commandMap["ExtraMsg"] = extraMsg
	commandMap["Sign"] = sign
	commandMap["PartnerId"] = partnerId
	commandMap["ServerId"] = serverId

	requestMap := make(map[string]interface{})
	requestMap["CommandType"] = commandType.Login
	requestMap["Command"] = commandMap

	// 先登陆服务器
	request(requestMap)
}
