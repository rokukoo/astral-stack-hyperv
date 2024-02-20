package web

import (
	"astralstack-hyperv/web/controller"
	"github.com/gin-gonic/gin"
)

func NewRoute() *gin.Engine {
	app := gin.Default()

	app.GET("/health", func(ctx *gin.Context) {
		ctx.String(200, "Hello!")
	})

	controller.Init(app)

	return app
}
