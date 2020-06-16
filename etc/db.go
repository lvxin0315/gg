package etc

//mysql 配置
type MysqlConfig struct {
	Name     string `required:"true" default:"test" env:"MysqlDatabaseName"`
	User     string `required:"true" default:"root" env:"MysqlUser"`
	Password string `required:"true" default:"root" env:"MysqlPassword"`
	Host     string `required:"true" default:"127.0.0.1" env:"MysqlHost"`
	Port     uint   `required:"true" default:"3306" env:"MysqlPort"`
}
