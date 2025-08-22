package main

import (
	"context"
	"os"
	"time"

	"github.com/qlik-trial/go-service-kit/v29/application"
	"github.com/qlik-trial/go-service-kit/v29/log"
	"github.com/qlik-trial/go-service-kit/v29/operation"
	"github.com/qlik-trial/go-service-kit/v29/tracing"
	"github.com/qlik-trial/usage-telemetry-publisher/cmd/config"
	"github.com/qlik-trial/usage-telemetry-publisher/cmd/version"
	"github.com/qlik-trial/usage-telemetry-publisher/internal/dependencies"
	"github.com/qlik-trial/usage-telemetry-publisher/internal/processes"
)

func main() {
	ctx := context.Background()
	label := "main"

	logger := log.New(log.Options{
		Level:       config.Global.LogLevel,
		Environment: config.Global.Environment,
		// our own logger (1) and the GSK logger's depth is 6
		// so we get 7 (1 + 6)
		CallerDepth: 7,
	})
	operation.SetGlobalLogger(&logger)
	operation.SetGlobalLogLevel(config.Global.LogLevel)
	operation.Logger(ctx).Info("label", label, "message", "starting usage-telemetry-publisher service", "version", version.Version, "configuration", config.Global)

	shutdownChannel := make(chan struct{})
	appCtx := dependencies.CreateAppContext(ctx, shutdownChannel)

	if config.Global.TracingEnabled {
		tracingStopper, err := tracing.Initialize(ctx, config.ServiceName, tracing.WithLogger(&logger), tracing.WithServiceVersion(version.Version))
		if err != nil {
			operation.Logger(ctx).Error("label", label, "message", "error initializing tracing:", "error", err)
			os.Exit(1)
		}
		defer func() {
			if tracingStopper != nil {
				tracingStopper(ctx)
			}
		}()
	}

	gracePeriod := time.Duration(config.Global.TerminationGracePeriodSeconds) * time.Second
	app, err := application.New(gracePeriod, application.WithCustomShutdown(func() error {
		err := appCtx.Dispose(ctx)
		return err
	}))
	if err != nil {
		operation.Logger(ctx).Error("label", label, "message", "failed to create application, shutting down", "error", err)
		os.Exit(1)
	}

	process := processes.BuildAppProcesses(appCtx)
	for name, p := range process {
		app.Add(name, p)
	}

	err = app.Start()
	if err != nil {
		operation.Logger(ctx).Error("label", label, "message", "application stopped ungracefully", "error", err)
		os.Exit(1)
	}
	os.Exit(0)
}
