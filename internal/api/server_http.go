package api

import (
	"net/http"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"

	routes "github.com/zackmwangi/shell_golang_svc/internal/api_routes"
	"github.com/zackmwangi/shell_golang_svc/internal/config"
	"go.uber.org/zap"
)

func InitHTTPRoutingEngine(c *config.AppConfig) *gin.Engine {

	globalHTTPRoutingEngine := gin.New()
	var myLogger *zap.Logger

	if c.AppEnv == "prod" {
		gin.SetMode(gin.ReleaseMode)
		myLogger, _ = zap.NewProduction()

	} else {
		gin.SetMode(gin.DebugMode)
		myLogger, _ = zap.NewDevelopment()
	}

	globalHTTPRoutingEngine.Use(ginzap.Ginzap(myLogger, time.RFC3339, true))
	//TODO: Add HTTP middleware for rate limiting/throttling
	//TODO: Add HTTP middleware for TLS enforcement
	//TODO: Add HTTP middleware for request validation
	//TODO: Add HTTP middleware for payload compression
	//TODO: Add HTTP middleware for CORS
	//TODO: Add HTTP middleware for request id/telemetry
	//################
	globalHTTPRoutingEngine.Use(ginzap.RecoveryWithZap(myLogger, true))
	globalHTTPRoutingEngine.Use(gin.Recovery())
	return globalHTTPRoutingEngine
}

func AddHttpEndpoints(httpRoutingEngine *gin.Engine, c *config.AppConfig, s *Servers) *gin.Engine {

	routes.AddHealthEndpoints(httpRoutingEngine, c)
	routes.AddDocsEndpoints(httpRoutingEngine, c)
	//
	//#
	//
	routes.AddGrpcGatewayEndpoints(httpRoutingEngine, s.Grpc.gmux)
	//
	//#
	//Fallback to 404
	httpRoutingEngine.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error_info": "Requested resource was not found"})
	})

	return httpRoutingEngine
}
