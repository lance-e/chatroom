package global

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	SensitiveWords  []string
	MessageQueueLen int
)

func initConfig() {
	viper.SetConfigFile("./config/chatroom.yaml")
	viper.SetConfigName("chatroom")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(RootDir + "/config")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	SensitiveWords = viper.GetStringSlice("sensitive")
	MessageQueueLen = viper.GetInt("message-queue")

	//监控配置文件
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		viper.ReadInConfig()
		SensitiveWords = viper.GetStringSlice("sensitive")

	})
}
