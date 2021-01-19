package config

import (
	"github.com/sirupsen/logrus"
	"time"
)

type syncerConfig struct {
	Tables   []string
	Raw      bool
	ServerID int
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
		SyncerConfig.Tables = []string{}
		err = ggViper.Unmarshal(&SyncerConfig)
		if err != nil {
			panic(err)
		}
		logrus.Info("tables", SyncerConfig.Tables)
	}

}
