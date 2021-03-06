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

func sendToClients(msg Message, userConn *userConnection) {
	result := "{msg&;" +
		msg.User + "&;" +
		time.Unix(msg.Date, 0).Format("2006-01-02 15:04:05") + "&;" + msg.Msg + "}"
	system.Connections.Range(func(key, value interface{}) bool { // 向每个连接到的客户端上写入信息
		conn := value.(*userConnection)
		if userConn.id != key {
			_, _ = conn.uconn.Write([]byte(result))
		}
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
	system.Wg.Done()            // 所有线程完成
}

func catchError(funcName string) {
	err := recover()
	if err != nil {
		fmt.Printf("检测到函数%v崩溃\n", funcName)
		fmt.Println(err)
		quit()
	}
}

func sendHistoryMsg(connection *userConnection) {
	var Messages []Message
	var result = "{msghistory&;"
	_ = DB.Order("id DESC").
		Limit(60).
		Where("date <= ?", time.Now().Unix()).
		Find(&Messages)
	var tmpMsg Message
	var tmpResult string
	for index := range Messages {
		tmpMsg = Messages[len(Messages)-1-index]
		tmpResult = tmpMsg.User + "&;" +
			time.Unix(tmpMsg.Date, 0).Format("2006-01-02 15:04:05") + "&;" +
			tmpMsg.Msg
		if index != len(Messages)-1 {
			tmpResult += "&;"
		}
		result += tmpResult
	}
	result += "}"

	_, _ = connection.uconn.Write([]byte(result))
}

func disconnect(connection *userConnection) {
	promptDisconnect(connection)
	system.Connections.Delete(connection.id) // 从连接池中删除连接
	_ = connection.uconn.Close()
}
