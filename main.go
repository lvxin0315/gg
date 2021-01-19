package main

import (
	"github.com/lvxin0315/gg/config"
	"github.com/lvxin0315/gg/syncer"
	"github.com/sirupsen/logrus"
	"time"
)

func main() {
	// 加载配置
	config.InitConfig()

	logrus.Info(config.SyncerConfig)

	// 启动binlog
	binlogSyncer := new(syncer.BinlogSyncer)
	binlogSyncer.Start()
	go func() {
		time.Sleep(10 * time.Second)
		logrus.Info("sleep")
		binlogSyncer.Close()
	}()
	select {}
}
