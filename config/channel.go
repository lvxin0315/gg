package config

const (
	NatsChannel       = "nats"
	NatsStreamChannel = "nats_stream"
	RabbitMQChannel   = "rabbit"
)

var ChannelsConfig channelsConfig

type ChannelConfig struct {
	Type     string
	Host     string
	Port     int
	User     string
	Password string
}

type channelsConfig struct {
	Channels map[string]ChannelConfig
}
