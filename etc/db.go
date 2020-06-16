package etc

//mysql 配置
type mysqlConfig struct {
	Name         string `required:"true" env:"MysqlDatabaseName"`
	User         string `required:"true" env:"MysqlUser"`
	Password     string `required:"true" env:"MysqlPassword"`
	Host         string `required:"true" env:"MysqlHost"`
	Port         uint   `required:"true" env:"MysqlPort"`
	MaxIdleConns int    `yaml:"maxIdleConns" required:"true" env:"MysqlMaxIdleConns"`
	MaxOpenConns int    `yaml:"maxOpenConns" required:"true" env:"MysqlMaxOpenConns"`
}
