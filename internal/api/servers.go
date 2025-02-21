package api

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"net/http"

	v1 "github.com/zackmwangi/shell_golang_svc/internal/api_proto/v1"
	"github.com/zackmwangi/shell_golang_svc/internal/config"
	"github.com/zackmwangi/shell_golang_svc/internal/pkg/middleware"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

type (
	Servers struct {
		AppConfig *config.AppConfig
		Grpc      *MyGrpcServer
		Http      *MyHttpServer
		stopFn    sync.Once
	}

	MyHttpServer struct {
		http   *http.Server
		logger *zap.Logger
		hmux   *gin.Engine
		config *config.AppConfig
	}

	MyGrpcServer struct {
		grpc   *grpc.Server
		logger *zap.Logger
		gmux   *runtime.ServeMux
		config *config.AppConfig
	}
)

// ################################################################
// # HTTP
func (s *MyHttpServer) Run(ctx context.Context, httpAddr string, grpcAddress string) error {

	httpListener, err := net.Listen("tcp", httpAddr)

	if err != nil {
		s.logger.Fatal("error on http address : " + httpAddr)
		os.Exit(1)
	}

	hs := &http.Server{
		//TODO: add these to main config
		Handler:        s.hmux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running properly and accessible
	grpcServerEndpoint := flag.String("grpc-server-endpoint", grpcAddress, "gRPC server endpoint")
	mux := runtime.NewServeMux()

	grpcDialOpts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err = v1.RegisterMybackendGrpcSvcHandlerFromEndpoint(ctx, mux, *grpcServerEndpoint, grpcDialOpts)
	if err != nil {
		return err
	}

	s.http = hs

	s.logger.Sugar().Infof("HTTP service listening at %s", httpAddr)
	return s.http.Serve(httpListener)
}

func (s *MyHttpServer) Shutdown(ctx context.Context) {
	s.logger.Sugar().Infof("HTTP service gracefully shutting down ")

	if s.http != nil {
		if err := s.http.Shutdown(ctx); err != nil {
			s.logger.Fatal("graceful shutdown of HTTP service failed ")
		}
	}
}

// ################################################################
// # gRPC
func (s *MyGrpcServer) Run(ctx context.Context, grpcAddress string) error {

	var lc net.ListenConfig

	lis, err := lc.Listen(ctx, "tcp", grpcAddress)

	if err != nil {

		s.logger.Sugar().Fatalf("error on grpc address : %s", err)

	}

	authMiddlewareGrpc := middleware.UnaryAuthInterceptor

	s.grpc = grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: 5 * time.Minute,
		}),
		grpc.UnaryInterceptor(authMiddlewareGrpc),
	)

	reflection.Register(s.grpc)

	//############################################################
	//Register gRPC base services

	MybackendGrpcSvcServerImpl := NewMybackendGrpcSvcServerImpl(s.config)
	v1.RegisterMybackendGrpcSvcServer(s.grpc, MybackendGrpcSvcServerImpl)

	//########################

	//Register grpc-gateway for HTTP-like API
	grpcDialOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		//auth interceptor
	}

	errX := v1.RegisterMybackendGrpcSvcHandlerFromEndpoint(ctx, s.gmux, grpcAddress, grpcDialOpts)

	if errX != nil {
		s.logger.Sugar().Fatalf("failed to serve gRPC gateway : %s", err)
		return errX
	}

	//########################
	s.logger.Sugar().Infof("gRPC service listening at %v", lis.Addr())
	return s.grpc.Serve(lis)

}

func (s *MyGrpcServer) Shutdown(ctx context.Context) {
	s.logger.Sugar().Infof("gRPC service gracefully shutting down ")

	done := make(chan struct{}, 1)

	go func() {
		if s.grpc != nil {
			s.grpc.GracefulStop()
		}
		done <- struct{}{}
	}()

	select {
	case <-done:
	case <-ctx.Done():
		if s.grpc != nil {
			s.grpc.Stop()
		}
		s.logger.Fatal("graceful shutdown of gRPC server failed. ")
	}
}

//################################################################
// both

func (s *Servers) Run(ctx context.Context) (err error) {

	var ec = make(chan error, 2)

	ctx, cancel := context.WithCancel(ctx)

	//grpcGwMux := s.getGrpcGwMux(s.AppConfig)
	grpcGwMux := s.getGrpcGwMux()
	s.Grpc = &MyGrpcServer{
		logger: s.AppConfig.AppLogger,
		gmux:   grpcGwMux,
		config: s.AppConfig,
	}

	httpMux := s.getHttpMux(s.AppConfig)
	s.Http = &MyHttpServer{
		logger: s.AppConfig.AppLogger,
		hmux:   httpMux,
		config: s.AppConfig,
	}

	go func() {
		err := s.Grpc.Run(ctx, net.JoinHostPort(s.AppConfig.AppListenHostname, s.AppConfig.AppListenPortGrpc))
		if err != nil {
			err = fmt.Errorf("gRPC server stopped : %w", err)
		}
		ec <- err
	}()

	go func() {
		err := s.Http.Run(ctx, net.JoinHostPort(s.AppConfig.AppListenHostname, s.AppConfig.AppListenPortHttp), net.JoinHostPort(s.AppConfig.AppListenHostname, s.AppConfig.AppListenPortGrpc))
		if err != nil {
			err = fmt.Errorf("HTTP server stopped : %w", err)
		}
		ec <- err
	}()

	var es []string

	for i := 0; i < cap(ec); i++ {
		if err := <-ec; err != nil {
			es = append(es, err.Error())
			if ctx.Err() == nil {
				s.Shutdown(context.Background())
			}
		}
	}

	if len(es) > 0 {
		err = errors.New(strings.Join(es, ", "))
	}

	cancel()

	return err
}

// func (s *Servers) getGrpcGwMux(c *config.AppConfig) *runtime.ServeMux {
func (s *Servers) getGrpcGwMux() *runtime.ServeMux {
	return runtime.NewServeMux()
}

func (s *Servers) getHttpMux(c *config.AppConfig) *gin.Engine {

	httpRoutingEngine := InitHTTPRoutingEngine(c)
	AddHttpEndpoints(httpRoutingEngine, c, s)
	return httpRoutingEngine
}

func (s *Servers) Shutdown(ctx context.Context) {
	s.stopFn.Do(func() {
		s.Http.Shutdown(ctx)
		s.Grpc.Shutdown(ctx)
	})
}
