package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/zackmwangi/shell_golang_svc/internal/config"
)

func AddDocsEndpoints(httpRoutingEngine *gin.Engine, c *config.AppConfig) *gin.Engine {

	vx := "v1" //TODO: Pull from config c
	vxRouterGroup := httpRoutingEngine.Group(vx).Use().(*gin.RouterGroup)
	vxRouterGroup.Static("/docs/openapi", "./docs/dist/")

	return httpRoutingEngine
}
