package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/lvxin0315/gg/databases"
	"github.com/lvxin0315/gg/etc"
	"github.com/sirupsen/logrus"
)

func Pong(c *gin.Context) {
	logrus.Println("我是test4的：", c.Query("a"))
	//TODO 配置文件
	logrus.Println("etc.Config.APPName:", etc.Config.APPName)
	logrus.Println("etc.Config.Contacts[0].Name:", etc.Config.Contacts[0].Name)
	logrus.Println("etc.Config.DB.Host:", etc.Config.DB.Host)
	logrus.Println("etc.Config.DB.MaxIdleConns:", etc.Config.DB.MaxIdleConns)
	logrus.Println("etc.Config.DB.MaxOpenConns:", etc.Config.DB.MaxOpenConns)
	//测试连接池
	t := 2000
	for {
		databases.NewDB().Exec("SELECT * FROM fa_config")
		t--
		if t == 0 {
			break
		}
	}

	c.JSON(200, gin.H{
		"message": "pong",
	})
}
