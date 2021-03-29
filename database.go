package main

import (
	"QAQServer/config"
	_ "QAQServer/config"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

var DB *gorm.DB

func initDatabase() {
	user := config.Config.GetString("database.user")
	pass := config.Config.GetString("database.pass")
	name := config.Config.GetString("database.name")
	host := config.Config.GetString("database.host")
	port := config.Config.GetString("database.port")

	var err error
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", user, pass, host, port, name)
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	// 添加数据表
	if !DB.Migrator().HasTable(&Message{}) { // 自动添加数据表
		_ = DB.Migrator().CreateTable(&Message{})
	}
}
