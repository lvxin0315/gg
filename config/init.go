package config

import (
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

var ggViper *viper.Viper

func InitConfig() {
	ggViper = viper.New()
	ggViper.SetConfigType("yaml")
	ggViper.AddConfigPath("./")
	err := ggViper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	// 解析
	err = ggViper.Unmarshal(&MysqlConfig)
	if err != nil {
		panic(err)
	}
	err = ggViper.Unmarshal(&SyncerConfig)
	if err != nil {
		panic(err)
	}
	// table 定时自动解析
	go updateSyncerConfig()
}
