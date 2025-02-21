package config

import (
	"log"

	"github.com/alexliesenfeld/health"
	"github.com/joho/godotenv"
	"github.com/zackmwangi/shell_golang_svc/internal/pkg/envo"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type (
	AppConfig struct {
		AppLogger          *zap.Logger
		HealthCheckerLive  health.Checker
		HealthCheckerReady health.Checker

		//################################

		SvcName    string
		SvcVersion string
		AppEnv     string

		AppListenHostname string
		AppListenPortHttp string
		AppListenPortGrpc string

		HealthEndpointPrefix string
		HealthEndpointLive   string
		HealthEndpointReady  string

		//###########
		//
		MetricsEndpointPrefix string
		MetricsEndpoint       string
		//
		OtelExporterOtlpURL string
		ObservabilityTracer trace.Tracer
	}
)

func InitAppConfig() *AppConfig {

	envConfigFile := ".env"
	err := godotenv.Load(envConfigFile)

	if err != nil {
		log.Printf("could not load .env file %s", err)
	} else {
		log.Printf("loaded .env file at %s", envConfigFile)
	}

	//#HARD DEFAULTS
	//########################

	svcNameDefault := "mybackend-svc-api"
	svcVersionDefault := "v1.0.0"

	appEnvDefault := "dev"

	//appEnvDefault := "prod"

	appListenHostnameDefault := "0.0.0.0"
	//appListenHostnameDefault := "172.20.0.28"
	appListenPortHttpDefault := "8081"
	appListenPortGrpcDefault := "8082"

	healthEndpointPrefixDefault := "/"
	healthEndpointLiveDefault := "health"
	healthEndpointReadyDefault := "ready"
	//###################################

	metricsEndpointPrefixDefault := "/"
	metricsEndpointDefault := "metrics"

	//otelExporterOtlpURLDefault := "http://172.21.0.101:14268/api/traces"
	otelExporterOtlpURLDefault := "http://localhost:14268/api/traces"

	//################################################################################################

	//# Swap hard defaults with .env/Configmap supplied values
	svcName := envo.EnvString("SVC_NAME", svcNameDefault)
	svcVersion := envo.EnvString("SVC_VERSION", svcVersionDefault)
	appEnv := envo.EnvString("APP_ENV", appEnvDefault)
	//ginMode := envo.EnvString("GIN_MODE", ginModeDefault)

	appListenHostname := envo.EnvString("APP_LISTEN_HOSTNAME", appListenHostnameDefault)
	appListenPortHttp := envo.EnvString("APP_LISTEN_PORT_HTTP", appListenPortHttpDefault)
	appListenPortGrpc := envo.EnvString("APP_LISTEN_PORT_GRPC", appListenPortGrpcDefault)

	healthEndpointPrefix := envo.EnvString("HEALTH_ENDPOINT_PREFIX", healthEndpointPrefixDefault)
	healthEndpointLive := envo.EnvString("HEALTH_ENDPOINT_LIVE", healthEndpointLiveDefault)
	healthEndpointReady := envo.EnvString("HEALTH_ENDPOINT_READY", healthEndpointReadyDefault)

	//################################################################################################

	metricsEndpointPrefix := envo.EnvString("METRICS_ENDPOINT_PREFIX", metricsEndpointPrefixDefault)
	metricsEndpoint := envo.EnvString("METRICS_ENDPOINT", metricsEndpointDefault)

	otelExporterOtlpURL := envo.EnvString("OTEL_EXPORTER_OTLP_ENDPOINT", otelExporterOtlpURLDefault)

	//################################################################################################
	return &AppConfig{

		SvcName:    svcName,
		SvcVersion: svcVersion,
		AppEnv:     appEnv,

		AppListenHostname: appListenHostname,
		AppListenPortHttp: appListenPortHttp,
		AppListenPortGrpc: appListenPortGrpc,

		HealthEndpointPrefix: healthEndpointPrefix,
		HealthEndpointLive:   healthEndpointLive,
		HealthEndpointReady:  healthEndpointReady,

		MetricsEndpointPrefix: metricsEndpointPrefix,
		MetricsEndpoint:       metricsEndpoint,
		//
		OtelExporterOtlpURL: otelExporterOtlpURL,
	}

}
