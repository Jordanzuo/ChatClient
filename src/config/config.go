package config

import (
	"github.com/Jordanzuo/goutil/configUtil"
	"github.com/Jordanzuo/goutil/debugUtil"
	"github.com/Jordanzuo/goutil/logUtil"
)

var (
	// 应用Key，用于加密
	AppKey string

	// 数据库连接字符串
	ChatServerCenterAPI string

	// 是否是DEBUG模式
	DEBUG bool
)

func init() {
	// 设置日志文件的存储目录
	logUtil.SetLogPath("LOG")

	// 读取配置文件内容
	config, err := configUtil.ReadJsonConfig("config.ini")
	checkError(err)

	// 解析AppKey
	AppKey, err = configUtil.ReadStringJsonValue(config, "AppKey")

	// 解析ChatServerCenterAPI
	ChatServerCenterAPI, err = configUtil.ReadStringJsonValue(config, "ChatServerCenterAPI")
	checkError(err)

	// 解析DEBUG配置
	debug, err := configUtil.ReadBoolJsonValue(config, "DEBUG")
	checkError(err)

	// 为DEBUG模式赋值
	DEBUG = debug

	// 设置debugUtil的状态
	debugUtil.SetDebug(debug)

	debugUtil.Println("ChatServerCenterAPI", ChatServerCenterAPI)
	debugUtil.Println("DEBUG:", debug)
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
