package global

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	SensitiveWords []string
)

func initConfig() {
	viper.SetConfigFile("./config/chatroom.yaml")
	viper.SetConfigName("chatroom")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	SensitiveWords = viper.GetStringSlice("sensitive")
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		viper.ReadInConfig()
		SensitiveWords = viper.GetStringSlice("sensitive")

	})
}
