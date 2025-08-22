package processes

import (
	"context"
	"errors"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	gskApplication "github.com/qlik-trial/go-service-kit/v29/application"
	"github.com/qlik-trial/go-service-kit/v29/healthcheck"
	gskMetrics "github.com/qlik-trial/go-service-kit/v29/metrics"
	gskMetricsMux "github.com/qlik-trial/go-service-kit/v29/metrics/mux"
	"github.com/qlik-trial/go-service-kit/v29/operation"
	"github.com/qlik-trial/usage-telemetry-publisher/cmd/config"
	"github.com/qlik-trial/usage-telemetry-publisher/cmd/version"
	"github.com/qlik-trial/usage-telemetry-publisher/internal/dependencies"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
)

type (

	// APIServer is a runnable that serves the API
	APIServer struct {
		address      string
		handler      http.Handler
		readTimeout  time.Duration
		writeTimeout time.Duration
		idleTimeout  time.Duration
	}
)

// BuildUsageTelemetryPublisherAPIServer builds a runnable API server
func BuildUsageTelemetryPublisherAPIServer(appCtx *dependencies.ApplicationContext) gskApplication.Runnable {
	metricsMiddleware := gskMetricsMux.MetricsMiddleware(
		*operation.GetGlobalLogger(),
		gskMetrics.MiddlewareOpts{EnableHistogram: true, EnableSummary: true},
	)

	// Instantiate router
	router := mux.NewRouter()

	// Add metrics handler
	gskMetrics.NewBuildInfo("usage_telemetry_publisher", version.Version, version.Revision)
	router.Handle("/metrics", promhttp.Handler()).Methods("GET").Name("metrics")

	// Add healthcheck handler
	healthHandler := healthcheck.NewHandler(operation.Logger(context.TODO()))
	if config.Global.MessagingEnabled {
		appCtx.MessagingClient.AddReadinessCheck(healthHandler)
	}
	router.Methods(http.MethodGet).Path("/health").Name("health").HandlerFunc(healthHandler.LiveEndpoint)
	router.Methods(http.MethodGet).Path("/ready").Name("ready").HandlerFunc(healthHandler.ReadyEndpoint)

	// Add debug handler
	if config.Global.EnableDebugEndpoints {
		router.PathPrefix("/debug/pprof/").HandlerFunc(pprof.Index)
	}

	// Add subrouter for API
	subrouter := router.PathPrefix("/v1").Subrouter()

	// Tracing middleware needs to be executed before the metrics middleware for exemplars to work
	subrouter.Use(otelmux.Middleware(config.ServiceName))
	subrouter.Use(metricsMiddleware)

	return &APIServer{
		address: config.Global.HTTPAddr,
		handler: handlers.RecoveryHandler()(router),
	}
}

// Start starts a RunnableHttpServer
func (svr *APIServer) Start(ctx context.Context) error {
	label := "api_server/Start"
	server := &http.Server{
		Addr:         svr.address,
		Handler:      svr.handler,
		ReadTimeout:  svr.readTimeout,
		WriteTimeout: svr.writeTimeout,
		IdleTimeout:  svr.idleTimeout,
	}

	go func() {
		<-ctx.Done()
		operation.Logger(ctx).Info("label", label, "message", "shutdown RestServer")
		server.Shutdown(ctx) //revive:disable:unhandled-error
	}()

	err := server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}

	return err
}
