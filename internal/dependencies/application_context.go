package dependencies

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	gskFeatures "github.com/qlik-trial/go-service-kit/v29/features"
	gskJWT "github.com/qlik-trial/go-service-kit/v29/jwt"
	"github.com/qlik-trial/go-service-kit/v29/operation"
	"github.com/qlik-trial/usage-telemetry-publisher/cmd/config"
	"github.com/qlik-trial/usage-telemetry-publisher/cmd/version"
	"github.com/qlik-trial/usage-telemetry-publisher/internal/auth"
	"github.com/qlik-trial/usage-telemetry-publisher/internal/events"
	"github.com/qlik-trial/usage-telemetry-publisher/internal/features"
	"github.com/qlik-trial/usage-telemetry-publisher/internal/messaging"
)

type (
	// ApplicationContext is a struct that holds references to resources that are used by multiple runnables
	ApplicationContext struct {
		TokenGenerator  auth.TokenGenerator
		MessagingClient messaging.EventListener
		FeaturesClient  features.FeaturesClient
	}
)

// CreateAppContext creates a new ApplicationContext
var CreateAppContext = func(ctx context.Context, stopChan <-chan struct{}) *ApplicationContext {
	appCtx := ApplicationContext{}
	appCtx.initTokenGenerator(ctx)

	if config.Global.LaunchDarklyEnabled {
		appCtx.initFeaturesClient(ctx)
	}

	if config.Global.MessagingEnabled {
		appCtx.initAndSubscribeMessagingClient(ctx, stopChan)
	}

	return &appCtx
}

func (appCtx *ApplicationContext) initAndSubscribeMessagingClient(ctx context.Context, stopChan <-chan struct{}) {
	label := "application_context/initAndSubscribeMessagingClient"
	hostname, err := os.Hostname()
	if err != nil {
		panic(fmt.Errorf("failed to get hostname: %w", err))
	}

	messagingClient, err := messaging.CreateClient(ctx, appCtx.TokenGenerator, hostname)
	if err != nil {
		operation.Logger(ctx).Error("label", label, "message", "failed to create messaging client", "error", err)
		panic(fmt.Errorf("failed to create messaging client: %w", err))
	}
	err = messagingClient.Connect(stopChan)
	appCtx.MessagingClient = messagingClient
	if err != nil {
		operation.Logger(ctx).Error("label", label, "message", "failed to connect to solace", "error", err)
	}
	appCtx.subscribeToChannels(ctx)

}

func (appCtx *ApplicationContext) subscribeToChannels(ctx context.Context) {
	label := "application_context/subscribeToChannels"
	channels := strings.SplitSeq(config.Global.SolaceChannels, ",")
	for pfx := range channels {
		channel := strings.TrimSpace(pfx)
		if err := appCtx.MessagingClient.SubscribeEvent(channel,
			config.Global.SolaceStreamingQueueGroup,
			events.EventHandler(ctx, appCtx.FeaturesClient)); err != nil {
			operation.Logger(ctx).Error(
				"label", label, "message", "error subscribing to message queue", "error", err, "channel", channel,
			)
		}
		operation.Logger(ctx).Info("label", label, "message", "subscribing to queue", "channel", channel, "group", config.Global.SolaceStreamingQueueGroup)
	}
}

func (appCtx *ApplicationContext) initTokenGenerator(ctx context.Context) {
	tokenURI := gskJWT.WithTokenURL(config.Global.TokenURI)
	label := "application_context/initTokenGenerator"

	signer, err := gskJWT.NewSigner(config.Global.SecretKeyFile, config.ServiceName, tokenURI)
	if err != nil {
		operation.Logger(ctx).Error("label", label, "message", "failed to create auth token handler", "error", err)
		panic(fmt.Errorf("failed to create auth token handler: %w", err))
	}

	err = signer.LoadCredentials()
	if err != nil {
		operation.Logger(ctx).Error("label", label, "message", "initial static load of jwt signing credentials", "error", err)
	}

	loadErrs := signer.LoadCredentialsInterval(ctx, time.Minute)
	go func() {
		select {
		case err := <-loadErrs:
			operation.Logger(ctx).Error("label", label, "message", "error occurred reading jwt signing credentials", "error", err)
		case <-ctx.Done():
			return
		}
	}()

	appCtx.TokenGenerator = signer
}

// Dispose runs the shutdown process for resources owned by ApplicationContext. This includes closing the messaging client.
func (appCtx *ApplicationContext) Dispose(ctx context.Context) (err error) {
	label := "application_context/Dispose"
	operation.Logger(ctx).Info("label", label, "message", "disposing ApplicationContext resources...")
	allErrors := []error{}
	return errors.Join(allErrors...)
}

func (appCtx *ApplicationContext) initFeaturesClient(ctx context.Context) {
	label := "application_context/initFeaturesClient"
	opts := []gskFeatures.Option{
		gskFeatures.WithServiceName("usage-telemetry-publisher"),
		gskFeatures.WithStreamURI(config.Global.LaunchDarklyStreamURI),
		gskFeatures.WithServiceVersion(version.Version),
		gskFeatures.WithRegion(config.Global.Region),
		gskFeatures.WithLDLogger(operation.GetGlobalLogger()),
	}
	featuresClient, err := features.NewFeaturesClient(ctx, config.Global.LaunchDarklySdkKey, opts...)
	if err != nil {
		operation.Logger(ctx).Error("label", label, "message", "failed to create features client", "error", err)
		panic(fmt.Errorf("failed to create features client: %w", err))
	}
	if !featuresClient.Initialized() {
		operation.Logger(ctx).Error("label", label, "message", "features client not initialized")
		panic(fmt.Errorf("features client not initialized"))
	}
	operation.Logger(ctx).Info("label", label, "message", "features client initialized successfully")
	appCtx.FeaturesClient = featuresClient
}
