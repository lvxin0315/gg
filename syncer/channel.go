package syncer

import (
	"github.com/lvxin0315/gg/config"
	"github.com/sirupsen/logrus"
	"sync"
)

type channel interface {
	// 初始化
	init(config config.ChannelConfig) error
	// 发送消息
	sendMessage(subject string, data []byte) error
	// 健康检查（无阻塞）
	healthy()
	// 关闭
	close()
}

// 初始化消息channel
func (syncer *channelSyncer) initChannels() {
	for key, channelConfig := range config.ChannelsConfig.Channels {
		var c channel
		switch channelConfig.Type {
		case config.NatsChannel:
			c = new(natsChannel)
		case config.RabbitMQChannel:
			// TODO
			continue
			//c = new(natsChannel)
		default:
			logrus.Warn("暂时不支持的通讯方式：", channelConfig.Type)
			continue
		}
		// 初始化
		err := c.init(channelConfig)
		if err != nil {
			logrus.Error("init error: ", err)
			continue
		}
		// 健康
		go c.healthy()
		// 记录
		syncer.channelList.Store(key, c)
	}
}

// 关闭channel
func (syncer *channelSyncer) closeChannels() {
	syncer.channelList.Range(func(key, value interface{}) bool {
		value.(channel).close()
		return true
	})
}

// 消息处理
type channelSyncer struct {
	// 总的通讯集合
	channelList sync.Map
}
