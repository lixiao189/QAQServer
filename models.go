package main

import (
	"net"
	"sync"
)

type userConnection struct {
	name      string   // 连接的用户名
	id        string   // 当前连接的标识号
	uconn     net.Conn // 该用户的连接
	loginTime int64    // 该用户上线的时间
}

type Message struct {
	ID   uint
	Msg  string // 消息主体
	User string // 发送者
	Date int64  // 发送时间
}

type System struct { // 存储服务器后端系统数据
	Listener    net.Listener
	Connections sync.Map
}
