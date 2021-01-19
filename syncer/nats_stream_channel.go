package syncer

import (
	"fmt"
	"github.com/lvxin0315/gg/config"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"github.com/sirupsen/logrus"
)

// nats的channel
type natsStreamChannel struct {
	config        config.ChannelConfig
	natsUrl       string
	connect       *nats.Conn
	streamConnect stan.Conn
}

// 初始化nats连接
func (channel *natsStreamChannel) init(config config.ChannelConfig) error {
	// nats://127.0.0.1:4222
	channel.config = config
	channel.natsUrl = fmt.Sprintf("nats://%s:%d", config.Host, config.Port)
	nc, err := nats.Connect(channel.natsUrl)
	if err != nil {
		logrus.Error("natsChannel.init: ", err)
		return err
	}
	sc, err := stan.Connect("test-cluster", "gg", stan.NatsConn(nc))
	if err != nil {
		logrus.Error("natsStreamChannel.init: ", err)
		return err
	}
	channel.connect = nc
	channel.streamConnect = sc
	return nil
}

// 发送消息
func (channel *natsStreamChannel) sendMessage(subject string, data []byte) error {
	if !channel.connect.IsConnected() {
		err := channel.init(channel.config)
		if err != nil {
			return err
		}
	}
	err := channel.streamConnect.Publish(subject, data)
	if config.CommonConfig.Debug {
		logrus.Debug("natsStreamChannel - subject: ", subject, " data: ", string(data))
	}
	return err
}

// 健康
func (channel *natsStreamChannel) healthy() {

}

// 关闭
func (channel *natsStreamChannel) close() {
	channel.connect.Close()
	_ = channel.streamConnect.Close()
}
