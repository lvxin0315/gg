package syncer

import (
	"fmt"
	"github.com/lvxin0315/gg/config"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

// nats的channel
type natsChannel struct {
	config  config.ChannelConfig
	natsUrl string
	connect *nats.Conn
}

// 初始化nats连接
func (channel *natsChannel) init(config config.ChannelConfig) error {
	// nats://127.0.0.1:4222
	channel.config = config
	channel.natsUrl = fmt.Sprintf("nats://%s:%d", config.Host, config.Port)
	nc, err := nats.Connect(channel.natsUrl)
	if err != nil {
		logrus.Error("natsChannel.init: ", err)
		return err
	}
	channel.connect = nc
	return nil
}

// 发送消息
func (channel *natsChannel) sendMessage(subject string, data []byte) error {
	if !channel.connect.IsConnected() {
		err := channel.init(channel.config)
		if err != nil {
			return err
		}
	}
	err := channel.connect.Publish(subject, data)
	if config.CommonConfig.AppDebug {
		logrus.Debug("subject: ", subject, " data: ", string(data))
	}
	return err
}

// 健康
func (channel *natsChannel) healthy() {

}

// 关闭
func (channel *natsChannel) close() {

}
