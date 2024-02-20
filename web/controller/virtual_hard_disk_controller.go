package controller

import (
	"astralstack-hyperv/service/virtual_hard_disk_service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func createVirtualHardDiskEndpoint(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "")
}

func listVirtualHardDiskEndpoint(ctx *gin.Context) {
	list, err := virtual_hard_disk_service.List()
	if err != nil {
		log.Panicln(err)
		return
	}
	ctx.JSON(200, list)
}
