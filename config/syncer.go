package config

import (
	"github.com/sirupsen/logrus"
	"time"
)

type SyncerTableConfig struct {
	Name    string
	Channel string
}

type syncerConfig struct {
	Tables                map[string]SyncerTableConfig
	Raw                   bool
	ServerID              int
	UpdateTableColumnTime int
	Subject               string
}

var SyncerConfig syncerConfig

func updateSyncerConfig() {
	for {
		time.Sleep(15 * time.Second)
		err := ggViper.ReadInConfig()
		if err != nil {
			logrus.Error("updateSyncerConfig:", err)
			return
		}
		err = ggViper.Unmarshal(&SyncerConfig)
		if err != nil {
			panic(err)
		}
		logrus.Info("tables", SyncerConfig.Tables)
	}

}
