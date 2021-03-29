package main

import _ "QAQServer/config" // 导入配置

var system System // 全局的系统数据

func main() {
	initDatabase()   // 连接数据库
	start()          // 开始程序
	go manage()      // 管理用户连接
	system.Wg.Wait() // 等待进程的结束
}
