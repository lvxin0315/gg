package main

import (
	"github.com/lvxin0315/gg/config"
	"github.com/lvxin0315/gg/syncer"
	"github.com/sirupsen/logrus"
)

func main() {
	// 加载配置
	config.InitConfig()
	// 日志配置
	if config.CommonConfig.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.WarnLevel)
	}

	// 启动binlog
	binlogSyncer := new(syncer.BinlogSyncer)
	binlogSyncer.Start()
	defer binlogSyncer.Close()
	//go func() {
	//	time.Sleep(10 * time.Second)
	//	logrus.Info("sleep")
	//	binlogSyncer.Close()
	//}()
	select {}
}
