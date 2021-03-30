package main

import (
	"fmt"
	"time"
)

func promptDisconnect(connection *userConnection) {
	fmt.Println(time.Now().String() + " 连接 " + connection.id + " 下线")
}

func saveToDB(msg Message) {
	DB.Create(&msg)
}

func sendToClients(msg Message) {
	result := "message&;" + msg.Group +
		"&;" + msg.User + "&;" +
		fmt.Sprint(time.Unix(msg.Date, 0).Format("2006-01-02 15:04:05")) + "&;" + msg.Msg
	system.Connections.Range(func(key, value interface{}) bool { // 向每个连接到的客户端上写入信息
		conn := value.(userConnection)
		_, _ = conn.uconn.Write([]byte(result))
		return true
	})
}

func sendHistoryMsg(connection *userConnection, groupName string) {
	var Messages []Message
	var result = "historyMessage"
	_ = DB.Where("`group` = ? AND date <= ?", groupName, connection.loginTime).Find(&Messages)
	for _, v := range Messages {
		result += v.User + "&;" +
			fmt.Sprint(time.Unix(v.Date, 0).Format("2006-01-02 15:04:05")) + "&;" +
			v.Msg
	}
	_, _ = connection.uconn.Write([]byte(result))
}

func disconnect(connection *userConnection) {
	promptDisconnect(connection)
	system.Connections.Delete(connection.id) // 从连接池中删除连接
	_ = connection.uconn.Close()
}
