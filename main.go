package main

import (
	_ "QAQServer/config"
	_ "QAQServer/database"
	"QAQServer/services"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		fmt.Println("服务器启动")
	}

	// 检测系统事件
	sigs := make(chan os.Signal, 4)
	signal.Notify(
		sigs,
		syscall.SIGHUP,
		syscall.SIGQUIT,
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	go func() {
		<-sigs               // 阻塞该代码 直到有终止程序信号被接收
		_ = listener.Close() // 关闭监听
		fmt.Println("程序退出中...")
		time.Sleep(time.Second)
		os.Exit(0)
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
		} else {
			go services.UserConnection(conn) // 如果连接成功开启线程接管连接
		}
	}
}
