package rpc

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/Jordanzuo/goutil/debugUtil"
	"github.com/Jordanzuo/goutil/logUtil"
	"time"
)

var (
	// 客户端对象
	clientObj *client
)

// 处理客户端逻辑
func handleClient() {
	for {
		content, ok := clientObj.getValidMessage()
		if !ok {
			break
		}

		debugUtil.Printf("content:%s\n", string(content))

		// 处理数据
		// 1、将结果反序列化
		responseMap := make(map[string]interface{})
		err := json.Unmarshal(content, &responseMap)
		if err != nil {
			logUtil.Log(fmt.Sprintf("反序列化%s出错，错误信息为：%s", string(content), err), logUtil.Error, true)
			continue
		}

		// 2、判断Code是否为0
		code_init, ok := responseMap["Code"].(float64)
		if !ok {
			fmt.Println(fmt.Sprintf("Code:%v，不是int类型", responseMap["Code"]))
			continue
		}
		code := int(code_init)
		if code != 0 {
			// 解析Message
			message, ok := responseMap["Message"].(string)
			if !ok {
				fmt.Println(fmt.Sprintf("Message:%v，不是string类型", responseMap["Message"]))
				continue
			}
			fmt.Println("返回结果不正确，错误信息为：", message)
			continue
		}

		// 3、判断CommandType
		commandType_int, ok := responseMap["CommandType"].(float64)
		if !ok {
			fmt.Println(fmt.Sprintf("CommandType:%v，不是int类型", responseMap["CommandType"]))
			continue
		}
		commandType := int(commandType_int)
		switch commandType {
		case 1: //Login
			fmt.Println("登陆成功，可以发送信息")
			loginSucceedCh <- 1
		case 2: //Logout
		case 3: //SendMessage
			// 获取Data
			dataMap, ok := responseMap["Data"].(map[string]interface{})
			if !ok {
				fmt.Println(fmt.Sprintf("Data:%v，不是map[string]interface{}类型", responseMap["Data"]))
				continue
			}

			// 获取Message
			message, ok := dataMap["Message"].(string)
			if !ok {
				fmt.Println(fmt.Sprintf("Message:%v，不是string类型", dataMap["Message"]))
				continue
			}

			// 判断是不是系统消息
			var name string
			var serverName string
			var fromMap map[string]interface{}
			if dataMap["From"] == nil {
				fmt.Printf("[%s]：%s\n", "系统消息", message)
			} else {
				if name, ok = dataMap["From"].(string); ok {

				} else if fromMap, ok = dataMap["From"].(map[string]interface{}); ok {
					// 获取发送者名称：
					if name, ok = fromMap["Name"].(string); !ok {
						fmt.Println(fmt.Sprintf("Name:%v，不是string类型", fromMap["Name"]))
						continue
					}

					if serverName, ok = fromMap["ServerName"].(string); !ok {
						fmt.Println(fmt.Sprintf("ServerName:%v，不是string类型", fromMap["ServerName"]))
						continue
					}
				} else {
					fmt.Println(fmt.Sprintf("From:%v，不是map[string]interface{}类型", dataMap["From"]))
					continue
				}

				fmt.Printf("[%s]%s说：%s\n", serverName, name, message)
			}
		case 4: //UpdatePlayerInfo
		}
	}
}

// 启动客户端
// ch：通道，用于传输连接成功的结果
func start(ch chan int) {
	// 连接指定的端口
	msg := ""
	conn, err := net.DialTimeout("tcp", chatServerAddress, 2*time.Second)
	if err != nil {
		msg = fmt.Sprintf("Dial Error: %s", err)
	} else {
		msg = fmt.Sprintf("Connect to the server. (local address: %s)", conn.LocalAddr())
	}

	logUtil.Log(msg, logUtil.Info, true)
	debugUtil.Println(msg)

	// 发送连接失败的通知
	if err != nil {
		ch <- 0
		return
	}

	// 创建客户端对象
	clientObj = newClient(conn)

	// 发送连接成功的通知
	ch <- 1

	defer func() {
		conn.Close()
		clientObj = nil
	}()

	// 死循环，不断地读取数据，解析数据，发送数据
	for {
		// 先读取数据，每次读取1024个字节
		readBytes := make([]byte, 1024)

		// Read方法会阻塞，所以不用考虑异步的方式
		n, err := conn.Read(readBytes)
		if err != nil {
			var errMsg string

			// 判断是连接关闭错误，还是普通错误
			if err == io.EOF {
				errMsg = fmt.Sprintf("另一端关闭了连接：%s，读取到的字节数为：%d", err, n)
				clientObj.conn.Close()
			} else {
				errMsg = fmt.Sprintf("读取数据错误：%s，读取到的字节数为：%d", err, n)
			}

			fmt.Println(errMsg)
			logUtil.Log(errMsg, logUtil.Error, true)

			//退出
			os.Exit(0)
		}

		// 将读取到的数据追加到已获得的数据的末尾
		clientObj.appendContent(readBytes[:n])

		// 已经包含有效的数据，处理该数据
		handleClient()
	}
}

// 启动客户端（连接ChatServer）
func StartClient(startCh chan int) {
	// 监听连接成功通道
	ch := make(chan int)
	go start(ch)

	//阻塞直到连接成功
	ret := <-ch
	if ret == 0 {
		panic("连接ChatServer失败，请检查配置")
	}

	// 发送login消息
	login()

	//阻塞直到登录成功或超时
	select {
	case <-loginSucceedCh:
		// 发送心跳包
		go heartBeat()
	case <-func() chan bool {
		timeout := make(chan bool, 1)
		go func() {
			time.Sleep(30 * time.Second)
			timeout <- false
		}()
		return timeout
	}():
		debugUtil.Println("Login Timeout")

		// 如果是启动程序调用，则panic，否则不处理
		panic("登录ChatServerCenter超时，请检查配置")
	}

	startCh <- 1
}
