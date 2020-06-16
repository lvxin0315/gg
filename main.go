package main

import (
	"github.com/gin-gonic/gin"
	"github.com/lvxin0315/gg/middlewares"
	"github.com/lvxin0315/gg/routers"
)

func main() {
	engine := gin.Default()
	//加载路由
	routers.InitRouter(engine)
	//中间件-跨域
	engine.Use(middlewares.Cors())

	engine.Run()
}
