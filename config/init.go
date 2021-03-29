/*
读取配置文件信息
*/

package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

var Config = viper.New()

func init() {
	Config.SetConfigName("config")   // 配置文件名称(无扩展名)
	Config.SetConfigType("yaml")     // 如果配置文件的名称中没有扩展名，则需要配置此项
	Config.AddConfigPath("./config") // 查找配置文件所在的路径
	Config.WatchConfig()
	err := Config.ReadInConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
