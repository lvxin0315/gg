package syncer

import (
	"fmt"
	"github.com/lvxin0315/gg/config"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

// rabbitMQ的channel
type rabbitMQChannel struct {
	config  config.ChannelConfig
	amqpUrl string
	connect *amqp.Connection
	ch      *amqp.Channel
}

func (channel *rabbitMQChannel) init(config config.ChannelConfig) error {
	channel.config = config
	channel.amqpUrl = fmt.Sprintf("amqp://%s:%s@%s:%d/",
		config.User,
		config.Password,
		config.Host,
		config.Port)
	conn, err := amqp.Dial(channel.amqpUrl)
	if err != nil {
		logrus.Error("rabbitMQChannel.init - Dial: ", err)
		return err
	}
	ch, err := conn.Channel()
	if err != nil {
		logrus.Error("rabbitMQChannel.init - Channel: ", err)
		return err
	}
	channel.connect = conn
	channel.ch = ch
	return nil
}

func (channel *rabbitMQChannel) sendMessage(subject string, data []byte) error {
	if channel.connect.IsClosed() {
		err := channel.init(channel.config)
		if err != nil {
			return err
		}
	}
	err := channel.ch.Publish(channel.config.Exchange,
		subject,
		true,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        data,
		})
	if err != nil {
		logrus.Error("rabbitMQChannel.sendMessage: ", err)
		return err
	}
	if config.CommonConfig.Debug {
		logrus.Debug("rabbitMQChannel - subject: ", subject, " data: ", string(data))
	}
	return nil
}

func (channel *rabbitMQChannel) healthy() {

}

func (channel *rabbitMQChannel) close() {
	_ = channel.ch.Close()
	_ = channel.connect.Close()
}
