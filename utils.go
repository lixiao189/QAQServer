package main

import (
	"fmt"
	"time"
)

func promptDisconnect(connection *userConnection) {
	fmt.Println(time.Now().String() + " 连接 " + connection.id + " 下线")
}

func promptConnect(connection *userConnection) {
	fmt.Println(time.Now().String() + " 连接 " + connection.id + " 上线")
}

func saveToDB(msg Message) {
	DB.Create(&msg)
}

func sendToClients(msg Message) {
	result := "message&;" + msg.Group +
		"&;" + msg.User + "&;" +
		fmt.Sprint(time.Unix(msg.Date, 0).Format("2006-01-02 15:04:05")) + "&;" + msg.Msg
	system.Connections.Range(func(key, value interface{}) bool { // 向每个连接到的客户端上写入信息
		conn := value.(*userConnection)
		_, _ = conn.uconn.Write([]byte(result))
		return true
	})
}

func quit() { // 关闭所有连接后退出
	system.Cancel()
	system.Connections.Range(func(k, v interface{}) bool { // 关闭用户连接
		_ = v.(*userConnection).uconn.Close()
		return true
	})
	_ = system.Listener.Close() // 关闭系统监听
	fmt.Println("\n程序退出中")
	time.Sleep(time.Second * 3) // 等待所有的连接关闭
	system.Wg.Done()            // 当前线程完成
}

func catchError() {
	err := recover()
	if err != nil {
		fmt.Println("检测到崩溃")
		fmt.Println(err)
		quit()
	}
}

func sendHistoryMsg(connection *userConnection, groupName string) {
	var Messages []Message
	var result = "historyMessage&;" + groupName + "&;"
	_ = DB.Order("id DESC").
		Limit(60).
		Where("`group` = ? AND date <= ?", groupName, time.Now().Unix()).
		Find(&Messages)

	size := 0
	var tmpMsg Message
	var tmpResult string
	for index := range Messages {
		if size > 9280 {
			break
		} else {
			tmpMsg = Messages[len(Messages)-1-index]
			tmpResult = tmpMsg.User + "&;" +
				fmt.Sprint(time.Unix(tmpMsg.Date, 0).Format("2006-01-02 15:04:05")) + "&;" +
				tmpMsg.Msg + "&;"
			result += tmpResult
			size += len([]byte(tmpResult))
		}
	}

	_, _ = connection.uconn.Write([]byte(result))
}

func disconnect(connection *userConnection) {
	promptDisconnect(connection)
	system.Connections.Delete(connection.id) // 从连接池中删除连接
	_ = connection.uconn.Close()
}
