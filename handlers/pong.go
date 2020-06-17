package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lvxin0315/gg/databases"
	"github.com/lvxin0315/gg/etc"
	"github.com/lvxin0315/gg/models"
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
		databases.NewDB().Exec("SELECT * FROM mall_article")
		t--

		//测试模型操作
		articleModel := new(models.MallArticle)
		articleModel.Author = fmt.Sprintf("Author%d", t)
		articleModel.Title = fmt.Sprintf("Title%d", t)
		articleModel.ShareTitle = fmt.Sprintf("ShareTitle%d", t)
		err := databases.NewDB().Model(articleModel).Save(articleModel).Error
		if err != nil {
			logrus.Error(err)
		}

		if t == 0 {
			break
		}
	}

	c.JSON(200, gin.H{
		"message": "pong",
	})
}
