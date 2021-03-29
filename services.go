/*
所有的后台服务代码
*/
package main

import (
	"QAQServer/config"
	"context"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func start() { // 启动
	system.CTX, system.Cancel = context.WithCancel(context.Background())
	system.MessageChan = make(chan Message, 128)

	system.Wg.Add(1)
	handleStop() // 启动对停止事件的处理

	var err error
	system.Listener, err = net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		fmt.Println("服务器启动")
	}
}

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
		system.Cancel()
		system.Connections.Range(func(k, v interface{}) bool { // 关闭用户连接
			_ = v.(userConnection).uconn.Close()
			return true
		})
		_ = system.Listener.Close() // 关闭系统监听
		fmt.Println("\n程序退出中")
		time.Sleep(time.Second * 3) // 等待所有的连接关闭
		system.Wg.Done()            // 当前线程完成
	}()
}

func manage() { // 管理连接
	go handleMessage() // 启动对消息的监听
	for {
		select {
		case <-system.CTX.Done():
			return
		default:
			conn, err := system.Listener.Accept() // 主循环接收请求
			if err == nil {                       // 当前连接没啥问题就处理这个连接
				userConn := userConnection{
					id:        uuid.NewV1().String(),
					uconn:     conn,
					loginTime: time.Now().Unix(),
				}
				system.Connections.Store(userConn.id, userConn) // 将用户的连接存入连接池
				go handleConnection(&userConn)                  // 开启新的线程管理连接
			}
		}
	}
}

func handleMessage() {
	for {
		select {
		case <-system.CTX.Done():
			return
		default:
			msg := <-system.MessageChan // 获取一个聊天记录
			go saveToDB(msg)
			go sendToClients(msg)
		}
	}
}

func handleConnection(userConn *userConnection) {
	// 接收用户指令
	clientInput := make([]byte, 512)
	var args []string
	for {
		select {
		case <-system.CTX.Done():
			return // 用 return 结束协程
		default: // 如果没有程序停止信号
			// 获取用户输入
			n, err := userConn.uconn.Read(clientInput)
			if err != nil {
				promptDisconnect(userConn)
				_ = userConn.uconn.Close()
				return
			}
			args = strings.Split(string(clientInput[0:n]), "&;")

			// 处理用户输入
			if args[0] == "user" {
				if args[1] == "disconnect" {
					promptDisconnect(userConn)
					_ = userConn.uconn.Close()
					return
				}
				if args[1] == "status" { // TODO: 日后加上心跳包功能

				}
				if args[1] == "connect" {
					userConn.name = args[2] // 设置用户名
				}
			}

			if args[0] == "msg" { // 向处理信息的后台协程传入信息
				if args[1] == "send" {
					system.MessageChan <- Message{
						Msg:   args[3],
						User:  userConn.name,
						Date:  time.Now().Unix(),
						Group: args[2],
					}
				}
				if args[1] == "list" { // 从数据库中获取上线前的所有历史记录

				}
			}

			if args[0] == "group" {
				if args[1] == "list" { // 获取当前的所有小组
					groups := config.Config.GetStringSlice("group")
					result := "group&;"
					for _, v := range groups {
						result += v + "&;"
					}
					_, err = userConn.uconn.Write([]byte(result))
					if err != nil {
						promptDisconnect(userConn)
						_ = userConn.uconn.Close()
						return
					}
				}
			}
		}
	}
}
