package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lvxin0315/gg/databases"
	"github.com/lvxin0315/gg/etc"
	"github.com/lvxin0315/gg/middlewares"
	"github.com/lvxin0315/gg/routers"
)

func main() {
	engine := gin.Default()
	//加载路由
	routers.InitRouter(engine)
	//中间件-跨域
	engine.Use(middlewares.Cors())
	//加载db
	databases.InitMysqlDB()
	databases.InitMemDB()
	engine.Run(fmt.Sprintf(":%s", etc.Config.Port))
}
