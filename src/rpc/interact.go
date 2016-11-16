package rpc

import (
	"fmt"
	"github.com/Jordanzuo/ChatServerModel/src/channelType"
	"github.com/Jordanzuo/ChatServerModel/src/commandType"
	"github.com/Jordanzuo/goutil/logUtil"
	"strings"
)

func Interaction(ch chan int) {
	// 处理内部未处理的异常，以免导致主线程退出，从而导致系统崩溃
	defer func() {
		if r := recover(); r != nil {
			logUtil.Log(fmt.Sprintf("通过recover捕捉到的未处理异常：%v", r), logUtil.Error, true)
		}

		ch <- 1
	}()

	// 源源不断地从键盘输入数据
	fmt.Println("请输入要发送的信息。如果要退出则输入q")
	for {
		var input string
		n, err := fmt.Scan(&input)
		if err != nil {
			fmt.Printf("input error:%s\n", err)
			break
		}

		if n == 0 {
			fmt.Println("你输入的数据为空，请重新输入。如果要退出则输入q")
			continue
		}

		// 如果选择退出
		if input == "q" {
			clientObj.conn.Close()
			break
		}

		// 发送数据
		request(assembleMessageParam(input))
	}
}

func assembleMessageParam(message string) map[string]interface{} {
	commandMap := make(map[string]interface{})

	// 根据message来判断是世界频道，还是公会频道
	msgList := strings.Split(message, ":")
	if len(msgList) == 1 {
		commandMap["ChannelType"] = channelType.World
		commandMap["Message"] = message
	} else {
		if msgList[0] == "Union" {
			commandMap["ChannelType"] = channelType.Union
			commandMap["Message"] = msgList[1]
		} else if msgList[0] == "Private" {
			commandMap["ChannelType"] = channelType.Private
			commandMap["ToPlayerId"] = msgList[1]
			commandMap["Message"] = msgList[2]
		} else if msgList[0] == "Cross" {
			commandMap["ChannelType"] = channelType.CrossServer
			commandMap["Message"] = message
		}
	}

	requestMap := make(map[string]interface{})
	requestMap["CommandType"] = commandType.SendMessage
	requestMap["Command"] = commandMap

	return requestMap
}
