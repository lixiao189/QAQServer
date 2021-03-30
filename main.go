package main

import (
	_ "QAQServer/config"
	"context"
	"fmt"
	"net"
	"os"
) // 导入配置

var system System // 全局的系统数据

func main() {
	// 初始化程序
	initDatabase() // 初始化数据库
	system.CTX, system.Cancel = context.WithCancel(context.Background())
	system.MessageChan = make(chan Message, 128)
	system.Wg.Add(1)
	var err error
	system.Listener, err = net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		fmt.Println("服务器启动")
	}

	handleStop()     // 启动对停止事件的处理
	go manage()      // 开启管理服务协程
	system.Wg.Wait() // 等待所有进程的结束
}
