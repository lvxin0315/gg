package etc

import (
	"github.com/jinzhu/configor"
)

var Config = struct {
	APPName  string `default:"app name"`
	Contacts []struct {
		Name  string
		Email string `required:"true"`
	}
	DB *mysqlConfig `default:"db"`
}{}

func init() {
	configor.Load(&Config, "config.yml")
}
