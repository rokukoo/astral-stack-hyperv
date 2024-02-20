package controller

import "github.com/gin-gonic/gin"

func Init(app *gin.Engine) {
	initVHDController(app)
}

func initVHDController(app *gin.Engine) {
	virtual_hard_disk := app.Group("/virtual_hard_disk")

	{
		virtual_hard_disk.GET("list", listVirtualHardDiskEndpoint)
		virtual_hard_disk.POST("create", createVirtualHardDiskEndpoint)
	}
}
