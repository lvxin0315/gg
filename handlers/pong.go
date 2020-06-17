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
	//试试配置文件
	//tryConfig()
	//试试连接池
	//tryDBClient()
	//试试模型操作-查询100条
	tryDBSelect()
	//返回值：pong
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

//试试连接池
func tryDBClient() {
	t := 20
	for {
		databases.NewDB().Exec("SELECT * FROM mall_article")
		t--
		//试试模型操作-插入数据
		articleModel := new(models.MallArticle)
		articleModel.Author = fmt.Sprintf("Author%d", t)
		articleModel.Title = fmt.Sprintf("Title%d", t)
		articleModel.ShareTitle = fmt.Sprintf("ShareTitle%d", t)
		err := databases.NewDB().Model(models.MallArticle{}).Save(articleModel).Error
		if err != nil {
			logrus.Error(err)
		}

		if t == 0 {
			break
		}
	}
}

//试试查询
func tryDBSelect() {
	var articleModelList []*models.MallArticle
	err := databases.NewDB().Model(&models.MallArticle{}).Where("id > ?", 10).Limit(100).Scan(&articleModelList).Error
	if err != nil {
		logrus.Error(err)
	}
	logrus.Println(articleModelList[0].Author)
}

//试试配置文件
func tryConfig() {
	logrus.Println("etc.Config.APPName:", etc.Config.APPName)
	logrus.Println("etc.Config.Contacts[0].Name:", etc.Config.Contacts[0].Name)
	logrus.Println("etc.Config.DB.Host:", etc.Config.DB.Host)
	logrus.Println("etc.Config.DB.MaxIdleConns:", etc.Config.DB.MaxIdleConns)
	logrus.Println("etc.Config.DB.MaxOpenConns:", etc.Config.DB.MaxOpenConns)
}
