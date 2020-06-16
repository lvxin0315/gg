package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/lvxin0315/gg/handlers"
)

func InitRouter(engine *gin.Engine) {
	v1 := engine.Group("/v1")
	v1.Any("/ping", handlers.Pong)
}
