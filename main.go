package main

import (
	"encoding/json"
	"fmt"
	"github.com/Jordanzuo/ChatClient/src/config"
	"github.com/Jordanzuo/ChatClient/src/rpc"
	"github.com/Jordanzuo/ChatServerModel/src/centerResponseObject"
	"github.com/Jordanzuo/goutil/debugUtil"
	"github.com/Jordanzuo/goutil/logUtil"
	"github.com/Jordanzuo/goutil/securityUtil"
	"github.com/Jordanzuo/goutil/webUtil"
	"os"
	"strconv"
	"sync"
)

var (
	wg sync.WaitGroup
)

func init() {
	wg.Add(1)
}

func main() {
	// 获得输入的信息，并计算sign
	if len(os.Args[1:]) != 6 {
		fmt.Println("输入错误，请重新输入，格式为：id name partnerId serverId unionId extraMsg")
		os.Exit(1)
	}

	id, name, partnerId_str, serverId_str, unionId, extraMsg := os.Args[1], os.Args[2], os.Args[3], os.Args[4], os.Args[5], os.Args[6]
	extraMsg = `{"lv":1,"vip":10}`
	sign := securityUtil.Md5String(fmt.Sprintf("%s-%s-%s", id, name, config.AppKey), false)

	partnerId, err := strconv.Atoi(partnerId_str)
	if err != nil {
		fmt.Printf("PartnerId:%v的格式错误\n", partnerId_str)
	}
	serverId, err := strconv.Atoi(serverId_str)
	if err != nil {
		fmt.Printf("ServerId:%v的格式错误\n", serverId_str)
	}

	// 获取聊天服务器的地址
	chatServerUrl := getChatServerAddress()
	debugUtil.Println("ChatServerUrl:", chatServerUrl)
	if chatServerUrl == "" {
		fmt.Println("没有获取到ChatServerUrl，请检查服务器配置")
		return
	}

	// 设置属性
	rpc.SetConfig(chatServerUrl, id, name, unionId, extraMsg, sign, partnerId, serverId)

	// 启动客户端
	startCh := make(chan int)
	go rpc.StartClient(startCh)
	<-startCh

	// 与用户交互
	ch := make(chan int)
	go rpc.Interaction(ch)

	<-ch
}

func getChatServerAddress() string {
	// 定义请求参数
	postDict := make(map[string]string)

	// 连接服务器，以获取数据
	returnBytes, err := webUtil.PostWebData(config.ChatServerCenterAPI, postDict, nil)
	if err != nil {
		logUtil.Log(fmt.Sprintf("获取聊天服务器信息出错，错误信息为：%s", err), logUtil.Error, true)
		panic(err)
	}

	// 解析返回值
	responseObj := new(centerResponseObject.ResponseObject)
	if err = json.Unmarshal(returnBytes, &responseObj); err != nil {
		logUtil.Log(fmt.Sprintf("获取服务器组列表出错，反序列化返回值出错，错误信息为：%s", err), logUtil.Error, true)
		panic(err)
	}

	// 判断返回状态是否为成功
	if responseObj.Code != centerResponseObject.Con_Success {
		msg := fmt.Sprintf("获取服务器组列表出错，返回状态：%d，信息为：%s", responseObj.Code, responseObj.Message)
		logUtil.Log(msg, logUtil.Error, true)
		panic(msg)
	}

	if tmpDataMap, ok := responseObj.Data.(map[string]interface{}); !ok {
		msg := "获取聊天服务器信息出错，返回的数据不是map[string]interface{}类型"
		logUtil.Log(msg, logUtil.Error, true)
		panic(msg)
	} else {
		if chatServerUrl, ok := tmpDataMap["ChatServerUrl"].(string); ok {
			return chatServerUrl
		} else {
			panic("没有找到可用的聊天服务器")
		}
	}
}
