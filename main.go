package main

import _ "QAQServer/config" // 导入配置

var system System                         // 全局的系统数据
var messageChan = make(chan Message, 128) // 全局的消息通道

func main() {
	initDatabase() // 连接数据库
	startListen()  // 开启连接
	manage()       // 管理用户连接
}
