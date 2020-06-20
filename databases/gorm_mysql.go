package databases

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/lvxin0315/gg/etc"
)

// import _ "github.com/jinzhu/gorm/dialects/postgres"
// import _ "github.com/jinzhu/gorm/dialects/sqlite"
// import _ "github.com/jinzhu/gorm/dialects/mssql"

var gormDB *gorm.DB

//初始化数据库
func InitMysqlDB() {
	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		etc.Config.DB.User,
		etc.Config.DB.Password,
		etc.Config.DB.Host,
		etc.Config.DB.Port,
		etc.Config.DB.Name,
	))
	if err != nil {
		panic(err)
	}
	db.DB().SetMaxIdleConns(etc.Config.DB.MaxIdleConns)
	db.DB().SetMaxOpenConns(etc.Config.DB.MaxOpenConns)
	//plugin
	db.Callback().Create().After("gorm:create").Register("plugin:gg_after_create", ggAfterCreate)
	db.Callback().Query().Before("gorm:query_destination").Register("plugin:gg_before_query_destination", ggBeforeQueryDestination)
	gormDB = db
}

func NewDB() *gorm.DB {
	return gormDB.New()
}
