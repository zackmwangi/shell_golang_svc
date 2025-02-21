package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/zackmwangi/shell_golang_svc/internal/pkg/middleware"
)

func AddGrpcGatewayEndpoints(httpRoutingEngine *gin.Engine, s *runtime.ServeMux) *gin.Engine {

	//User
	httpApiVersion := "v1"
	AddGrpcGatewayEndpointsUser(httpApiVersion, httpRoutingEngine, s)

	//tester
	//httpApiVersion := "v1"
	//AddGrpcGatewayEndpointsTester(httpApiVersion, httpRoutingEngine, s)s
	return httpRoutingEngine
}

//group routes further

func AddGrpcGatewayEndpointsUser(httpApiVersion string, httpRoutingEngine *gin.Engine, s *runtime.ServeMux) *gin.Engine {

	authMiddlewareHttp := middleware.AuthMiddlewareHttp()

	httpApiRouterGroup := httpRoutingEngine.Group(httpApiVersion).Use(authMiddlewareHttp).(*gin.RouterGroup)

	httpApiRouterGroup.GET("/user/*{http_to_grpc_gateway}", gin.WrapH(s))
	httpApiRouterGroup.POST("/user/*{http_to_grpc_gateway}", gin.WrapH(s))

	return httpRoutingEngine
}

/*
func AddGrpcGatewayEndpointsTester(httpApiVersion string, httpRoutingEngine *gin.Engine, s *runtime.ServeMux) *gin.Engine {
	httpApiRouterGroup := httpRoutingEngine.Group(httpApiVersion).Use().(*gin.RouterGroup)

	//wildcard
	httpApiRouterGroup.GET("/tester/*{http_to_grpc_gateway}", gin.WrapH(s))
	httpApiRouterGroup.POST("/tester/*{http_to_grpc_gateway}", gin.WrapH(s))

	//specific
	httpApiRouterGroup.GET("/tester", gin.WrapH(s))
	httpApiRouterGroup.POST("/tester}", gin.WrapH(s))

	return httpRoutingEngine
}
*/
