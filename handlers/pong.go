package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/lvxin0315/gg/etc"
	"github.com/sirupsen/logrus"
)

func Pong(c *gin.Context) {
	logrus.Println("我是test4的：", c.Query("a"))
	//TODO 配置文件
	logrus.Println("etc.Config.APPName:", etc.Config.APPName)
	logrus.Println("etc.Config.APPName:", etc.Config.Contacts[0].Name)
	logrus.Println("etc.Config.DB.Host:", etc.Config.DB.Host)
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
