package processes

import (
	"github.com/qlik-trial/go-service-kit/v29/application"
	"github.com/qlik-trial/usage-telemetry-publisher/internal/dependencies"
)

// BuildAppProcesses builds all processes of the usage-telemetry-publisher application
var BuildAppProcesses = func(appCtx *dependencies.ApplicationContext) map[string]application.Runnable {
	processes := map[string]application.Runnable{
		"UsageTelemetryPublisherAPIServer": BuildUsageTelemetryPublisherAPIServer(appCtx),
	}

	return processes
}
