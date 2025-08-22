package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGlobal(t *testing.T) {
	assert.Equal(t, Global.DependentServicesReadinessCheckEnabled, defaultDependentServicesReadinessCheckEnabled)
	assert.Equal(t, Global.DependenciesHealthcheckTimeoutMilliseconds, defaultDependenciesHealthcheckTimeoutMilliseconds)
	assert.Equal(t, Global.HTTPAddr, defaultHTTPAddr)
	assert.Equal(t, Global.LogLevel, defaultLogLevel)
	assert.Equal(t, Global.Environment, defaultEnvironment)
	assert.Equal(t, Global.TerminationGracePeriodSeconds, defaultTerminationGracePeriodSeconds)
	assert.Equal(t, Global.MessagingEnabled, defaultMessagingEnabled)
	assert.Equal(t, Global.MessagingPublishBufferSize, defaultMessagingPublishBufferSize)
	assert.Equal(t, Global.MessagingConnectionCheckIntervalSeconds, defaultMessagingConnectionCheckIntervalSeconds)
	assert.Equal(t, Global.TracingEnabled, defaultTracingEnabled)
	assert.Equal(t, Global.PanicOnValidationErrors, defaultPanicOnValidationErrors)
	assert.Equal(t, Global.Region, defaultRegion)
	assert.Equal(t, Global.LaunchDarklyStreamURI, defaultLaunchDarklyStreamURI)
	assert.Equal(t, Global.LaunchDarklySdkKey, defaultLaunchDarklySdkKey)
	assert.Equal(t, Global.LaunchDarklyEnabled, defaultLaunchDarklyEnabled)
	assert.Equal(t, Global.AuthEnabled, defaultAuthEnabled)
	assert.Equal(t, Global.KeysUri, defaultKeysUri)
	assert.Equal(t, Global.AuthJwtAud, defaultAuthJwtAud)
	assert.Equal(t, Global.AuthJwtIss, defaultAuthJwtIss)
	assert.Equal(t, Global.SecretKeyFile, defaultSecretKeyFile)
	assert.Equal(t, Global.TokenURI, defaultTokenURI)
	assert.Equal(t, Global.SkipPurgeEvents, defaultSkipPurgeEvents)
}
