package routes

import (
	"github.com/alexliesenfeld/health"
	"github.com/gin-gonic/gin"
	"github.com/zackmwangi/shell_golang_svc/internal/config"
)

func AddHealthEndpoints(httpRoutingEngine *gin.Engine, c *config.AppConfig) *gin.Engine {

	httpRoutingEngine.GET(c.HealthEndpointPrefix+c.HealthEndpointLive, gin.WrapH(health.NewHandler(c.HealthCheckerLive)))
	httpRoutingEngine.GET(c.HealthEndpointPrefix+c.HealthEndpointReady, gin.WrapH(health.NewHandler(c.HealthCheckerReady)))

	return httpRoutingEngine
}
