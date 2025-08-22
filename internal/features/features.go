package features

import (
	context "context"

	gskFeatures "github.com/qlik-trial/go-service-kit/v29/features"
)

const (
	EventIngestionFlag = "usage-telemetry-event-ingestion"
	//PhasedRolloutFlag  = "TLV_x_USAGE_TELEMETRY_PUBLISHER"
)

type FeaturesClient interface {
	GetBoolGlobalFeature(ctx context.Context, featureFlag string, contextOptions ...gskFeatures.ContextOption) (value bool, err error)
	GetBoolTenantFeature(ctx context.Context, featureFlag, tenantID string, contextOptions ...gskFeatures.ContextOption) (value bool, err error)
	Initialized() bool
}

func NewFeaturesClient(ctx context.Context, ldSdkKey string, opts ...gskFeatures.Option) (c FeaturesClient, err error) {
	featuresClient, err := gskFeatures.NewLaunchDarklyClient(ctx, ldSdkKey, opts...)
	return featuresClient, err
}
