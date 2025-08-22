package e2e

import (
	"context"
	"testing"

	gskFeatures "github.com/qlik-trial/go-service-kit/v29/features"
	"github.com/qlik-trial/go-service-kit/v29/operation"
	"github.com/qlik-trial/usage-telemetry-publisher/cmd/config"
	"github.com/qlik-trial/usage-telemetry-publisher/internal/features"
	"github.com/stretchr/testify/assert"
)

func TestFeatureClient(t *testing.T) {
	ctx := context.Background()
	opts := []gskFeatures.Option{
		gskFeatures.WithServiceName("usage-telemetry-publisher-component-test"),
		gskFeatures.WithStreamURI(config.Global.LaunchDarklyStreamURI),
		gskFeatures.WithServiceVersion("test-version"),
		gskFeatures.WithRegion("test-region"),
		gskFeatures.WithLDLogger(operation.GetGlobalLogger()),
	}
	client, err := features.NewFeaturesClient(ctx, config.Global.LaunchDarklySdkKey, opts...)
	assert.Nil(t, err)

	assert.True(t, client.Initialized())

	value, err := client.GetBoolGlobalFeature(ctx, features.EventIngestionFlag)
	assert.Nil(t, err)
	assert.False(t, value)

	value, err = client.GetBoolTenantFeature(ctx, features.EventIngestionFlag, "mock-tenant-1")
	assert.Nil(t, err)
	assert.False(t, value)

	value, err = client.GetBoolTenantFeature(ctx, features.EventIngestionFlag, "mock-tenant-2")
	assert.Nil(t, err)
	assert.True(t, value)

}
