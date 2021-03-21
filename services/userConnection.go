/*
处理每个用户连接的服务
*/

package services

import (
	"bufio"
	uuid "github.com/satori/go.uuid"
	"net"
	"strings"
	"time"
)

type userConn struct { // 定义了一个用户连接
	uuid     string
	conn     net.Conn
	lastTime int64
}

var connections = make([]userConn, 128) // 连接池

func UserConnection(conn net.Conn) {
	// 添加用户连接
	connections = append(
		connections,
		userConn{
			uuid:     uuid.NewV1().String(),
			conn:     conn,
			lastTime: time.Now().Unix(), // 上次的连接时间
		},
	)

	reader := bufio.NewReader(conn)
	//writer := bufio.NewWriter(conn)
	for {
		clientInput := make([]byte, 256)
		n, _ := reader.Read(clientInput) // 获取用户输入
		clientInput = clientInput[:n]
		args := strings.Split(string(clientInput), " ") // 用户输入的命令参数

		if args[0] == "user" {

		}
		if args[0] == "msg" {

		}
		if args[0] == "group" {

		}
	}
}
