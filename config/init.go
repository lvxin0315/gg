package config

import (
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

var ggViper *viper.Viper

/**
 * @Author lvxin0315@163.com
 * @Description 初始化加载配置文件
 * @Date 9:30 上午 2021/1/18
 * @Param
 * @return
 **/
func InitConfig() {
	ggViper = viper.New()
	ggViper.SetConfigType("yaml")
	ggViper.AddConfigPath("./")
	err := ggViper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	// 解析
	err = ggViper.Unmarshal(&CommonConfig)
	if err != nil {
		panic(err)
	}
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
