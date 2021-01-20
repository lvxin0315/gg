package config

const (
	NatsChannel       = "nats"
	NatsStreamChannel = "nats_stream"
	RabbitMQChannel   = "rabbitmq"
)

var ChannelsConfig channelsConfig

type ChannelConfig struct {
	Type     string
	Host     string
	Port     int
	User     string
	Password string
	ClientID string // nats_stream
	Exchange string // rabbitmq交换机
}

type channelsConfig struct {
	Channels map[string]ChannelConfig
}
