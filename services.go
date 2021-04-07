/*
所有的后台服务代码
*/
package main

import (
	uuid "github.com/satori/go.uuid"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func handleStop() { // 检测退出信号
	sigs := make(chan os.Signal, 4)
	signal.Notify(
		sigs,
		syscall.SIGHUP,
		syscall.SIGQUIT,
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	go func() {
		<-sigs // 阻塞该代码 直到有终止程序信号被接收
		quit()
	}()
}

func manage() { // 管理连接
	defer catchError()
	go handlePackage() // 处理接收包的服务
	for {
		select {
		case <-system.CTX.Done():
			return
		default:
			conn, err := system.Listener.Accept() // 主循环接收请求
			if err == nil {                       // 当前连接没啥问题就处理这个连接
				go handleConnection(&userConnection{
					name:  "lazy",
					id:    uuid.NewV1().String(),
					uconn: conn,
				})
			}
		}
	}
}

func dropMessage() {
	for {
		select {
		case <-system.CTX.Done():
			return
		default:
			DB.Where("date < ?", time.Now().Unix()-3600*24*3).Delete(Message{})
			time.Sleep(time.Hour) // 每隔一小时清理一次消息记录
		}
	}
}

func handlePackage() {
	defer catchError()
	for {
		select {
		case <-system.CTX.Done():
			return
		default:
			packageData := <-system.PackageChan
			args := strings.Split(packageData, "&;")
			for len(args) < 8 { // 为了防止代码崩溃往后面填充空白参数
				args = append(args, "")
			}

			result, _ := system.Connections.Load(args[0])
			userConn := result.(*userConnection)
			if args[1] == "user" {
				if args[2] == "named" {
					// 设置该连接的用户昵称
					userConn.name = args[3]
				}
			}

			if args[1] == "msg" {
				if args[2] == "list" {
					sendHistoryMsg(userConn)
				}
				if args[2] == "send" {
					msg := Message{
						Msg:  args[3],
						User: userConn.name,
						Date: time.Now().Unix(),
					}
					saveToDB(msg)      // 将消息存到数据库中
					sendToClients(msg) // 将消息发给所有客户端
				}
			}
		}
	}
}

func handleConnection(userConn *userConnection) {
	defer catchError()

	system.Connections.Store(userConn.id, userConn) // 将用户的连接存入连接池
	promptConnect(userConn)                         // 提示上线

	// 处理用户发送的数据
	clientData := make([]byte, 64)
	isStarted := false // 包是否开始
	packageData := ""
	for {
		select {
		case <-system.CTX.Done():
			return
		default:
			// 获取用户输入
			n, err := userConn.uconn.Read(clientData)
			if err != nil {
				// 该协程处理的客户端失去连接
				disconnect(userConn)
				return
			}
			for i := 0; i < n; i++ {
				if clientData[i] == '{' {
					isStarted = true
					packageData += userConn.id + "&;"
					continue
				}
				if clientData[i] == '}' {
					isStarted = false
					system.PackageChan <- packageData
					packageData = "" // 上一个包已经结束 清空包的内容
					continue
				}

				if isStarted {
					packageData += string(clientData[i])
				}
			}
		}
	}
}
