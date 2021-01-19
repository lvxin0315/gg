package config

type mysqlConfig struct {
	Host     string
	User     string
	Password string
	Port     int
	Flavor   string
}

var MysqlConfig mysqlConfig
