package main

import _ "QAQServer/config" // 导入配置

var system System // 全局的系统数据

func main() {
	start()          // 开始程序
	go manage()      // 管理线程
	system.Wg.Wait() // 等待所有进程的结束
}
