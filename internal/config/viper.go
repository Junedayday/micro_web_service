package config

import "github.com/spf13/viper"

// 全局Viper变量
var Viper = viper.New()

func Load(configFilePath string) error {
	Viper.SetConfigName("config")       // config file name without file type
	Viper.SetConfigType("yaml")         // config file type
	Viper.AddConfigPath(configFilePath) // config file path
	return Viper.ReadInConfig()
}
