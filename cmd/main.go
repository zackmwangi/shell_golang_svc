package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/alexliesenfeld/health"
	"github.com/zackmwangi/shell_golang_svc/internal/api"
	"github.com/zackmwangi/shell_golang_svc/internal/config"
	"github.com/zackmwangi/shell_golang_svc/internal/metrix"
	applogger "github.com/zackmwangi/shell_golang_svc/internal/pkg/applogger"
	//
)

func main() {

	//CPU Max
	runtime.GOMAXPROCS(runtime.NumCPU())

	//Init context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	//#Init App config
	appConfig := config.InitAppConfig()

	//Init logger
	AppLogger := applogger.InitAppLogger()
	defer AppLogger.Sync()
	appConfig.AppLogger = AppLogger

	//Enable health check
	healthCheckerReady := health.NewChecker(health.WithCacheDuration(1*time.Second),
		health.WithTimeout(10*time.Second))

	appConfig.HealthCheckerReady = healthCheckerReady

	healthCheckerReady.Start()
	healthCheckerLive := health.NewChecker()

	appConfig.HealthCheckerLive = healthCheckerLive

	//###########################################################
	// Trace providers
	tp, errX := metrix.InitTraceProvider(appConfig)
	//
	if errX != nil {
		appConfig.AppLogger.Fatal(errX.Error())
	}

	defer func(ctx context.Context) {

		//defer loggerProvider.Shutdown(ctx)

		ctx, cancel := context.WithTimeout(ctx, time.Second*5)
		defer cancel()
		/*
			if err := lp.Shutdown(ctx); err != nil {
				log.Printf("Error shutting down otlp logger provider: %v", err)
			}
		*/

		if err := tp.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down otlp tracer provider: %v", err)
		}

	}(ctx)
	//###
	appObservabilityTracer := tp.Tracer(appConfig.SvcName)
	appConfig.ObservabilityTracer = appObservabilityTracer
	//###

	//################################################################################################

	s := &api.Servers{
		AppConfig: appConfig,
	}

	//################################################################################################
	//Handle Signals
	ec := make(chan error, 2)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	go func() {
		ec <- s.Run(context.Background())
	}()

	//######
	var err error

	select {

	case err = <-ec:
	case <-ctx.Done():

		haltCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		s.Shutdown(haltCtx)
		stop()

		err = <-ec
	}

	if err != nil {
		AppLogger.Sugar().Infof("service %s has shutdown : %s", appConfig.SvcName, err)
	}

}
