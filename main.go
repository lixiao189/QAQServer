package main

import _ "QAQServer/config" // 导入配置

var system System // 全局的系统数据

func main() {
	initDatabase()     // 连接数据库
	startListen()      // 开启连接 connections 是一个指针
	manageConnection() // 管理用户连接
}
