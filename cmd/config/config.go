package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/qlik-trial/go-service-kit/v29/operation"
	"github.com/qlik-trial/go-service-kit/v29/validator"
	"github.com/spf13/viper"
)

const (
	// ServiceName contains this service's name
	ServiceName                                       = "usage-telemetry-publisher"
	defaultDependentServicesReadinessCheckEnabled     = true
	defaultIntermediateStorageEnabled                 = false
	defaultDependenciesHealthcheckTimeoutMilliseconds = 10000
	defaultHTTPAddr                                   = ":8080"
	defaultLogLevel                                   = "info"
	defaultEnvironment                                = ""
	defaultTerminationGracePeriodSeconds              = 30
	defaultMessagingEnabled                           = false
	defaultMessagingPublishBufferSize                 = 100
	defaultMessagingConnectionCheckIntervalSeconds    = 5
	defaultTracingEnabled                             = false
	defaultPanicOnValidationErrors                    = false
	defaultRegion                                     = "local"
	defaultLaunchDarklyStreamURI                      = "http://ldRelay:8080/relay"
	defaultLaunchDarklySdkKey                         = "key"
	defaultLaunchDarklyEnabled                        = false
	defaultAuthEnabled                                = false
	defaultKeysUri                                    = "http://keys:8080"
	defaultAuthJwtAud                                 = "qlik.api.internal"
	defaultAuthS2SJwtAud                              = "qlik.api.internal/usage-telemetry-publisher"
	defaultAuthJwtIss                                 = "qlik.api.internal"
	defaultSecretKeyFile                              = "/run/secrets/qlik.com/usage-telemetry-publisher/service-key.yaml"
	defaultTokenURI                                   = "http://edge-auth:8080"
	defaultSolaceStreamingQueueGroup                  = "usage-telemetry-publisher"
	defaultSolaceChannels                             = ""
	defaultEventsFilePath                             = "/etc/config/events.yaml"
	defaultSkipPurgeEvents                            = true
	defaultFeatureFlagsEnabled                        = false
)

// Spec defines the schema for configurations
type Spec struct {
	// DependentServicesReadinessCheckEnabled is a flag to enable or disable the readiness of dependent services of usage-telemetry-publisher
	DependentServicesReadinessCheckEnabled     bool   `mapstructure:"dependent_services_readiness_check_enabled"`
	DependenciesHealthcheckTimeoutMilliseconds int    `mapstructure:"dependencies_healthcheck_timeout_milliseconds"`
	HTTPAddr                                   string `mapstructure:"http_addr"`
	LogLevel                                   string `mapstructure:"log_level"`
	Environment                                string `mapstructure:"environment"`
	TracingEnabled                             bool   `mapstructure:"tracing_enabled"`
	TerminationGracePeriodSeconds              int    `mapstructure:"termination_grace_period_seconds"`
	// SecretKeyFile is the file path to the yaml file containing the private key and other auth server settings
	SecretKeyFile string `mapstructure:"secret_key_file" validate:"required"`
	// TokenURI is the uri to use for requesting tokens
	TokenURI string `mapstructure:"token_uri" validate:"required"`

	Region                     string `mapstructure:"region"`
	LaunchDarklyEnabled        bool   `mapstructure:"launchdarkly_enabled"`
	LaunchDarklyStreamURI      string `mapstructure:"launchdarkly_stream_uri"`
	LaunchDarklySdkKey         string `mapstructure:"launchdarkly_sdk_key"`
	LaunchDarklySdkKeyFile     string `mapstructure:"launchdarkly_sdk_key_file"`
	SolaceChannels             string `mapstructure:"solace_channels"`
	MessagingEnabled           bool   `mapstructure:"messaging_enabled"`
	IntermediateStorageEnabled bool   `mapstructure:"intermediate_storage_enabled"`
	// MessagingConnectionCheckIntervalSeconds is the interval of checking the connection to messaging
	MessagingConnectionCheckIntervalSeconds int `mapstructure:"messaging_connection_check_interval_seconds" validate:"gte=0"`
	// MessagingPublishBufferSize is the maximum number of async published msgs to be buffered
	MessagingPublishBufferSize int `mapstructure:"messaging_publish_buffer_size" validate:"gte=0"`

	// PanicOnValidationErrors toggles whether or not to panic if there are validation errors in the config
	PanicOnValidationErrors bool `mapstructure:"panic_on_validation_errors"`
	// EnableDebugEndpoints toggles whether or not to turn on debug endpoints
	EnableDebugEndpoints bool `mapstructure:"enable_debug_endpoints"`
	// SolaceStreamingQueueGroup specifies the streaming group
	SolaceStreamingQueueGroup string `mapstructure:"solace_streaming_queue_group"`
	// AuthEnabled holds a flag indicating whether authentication is enabled.
	AuthEnabled    bool   `mapstructure:"auth_enabled"`
	KeysUri        string `mapstructure:"keys_uri"`
	AuthJwtAud     string `mapstructure:"auth_jwt_Aud"`
	AuthS2SJwtAud  string `mapstructure:"auth_s2s_jwt_Aud"`
	AuthJwtIss     string `mapstructure:"auth_jwt_iss"`
	EventsFilePath string `mapstructure:"events_file_path"`

	// SkipPurgeEvents if true event types that end with '.purged' are not written to mongo
	SkipPurgeEvents bool `mapstructure:"skip_purge_events"`

	// FeatureFlagsEnabled enables or disables the feature flags use
	FeatureFlagsEnabled bool `mapstructure:"feature_flags_enabled"`
}

// Global is a struct variable, holding global configuration values.
var Global = newSpec()

func newSpec() *Spec {
	return &Spec{
		DependentServicesReadinessCheckEnabled:     defaultDependentServicesReadinessCheckEnabled,
		DependenciesHealthcheckTimeoutMilliseconds: defaultDependenciesHealthcheckTimeoutMilliseconds,
		HTTPAddr:                                defaultHTTPAddr,
		LogLevel:                                defaultLogLevel,
		Environment:                             defaultEnvironment,
		TracingEnabled:                          defaultTracingEnabled,
		PanicOnValidationErrors:                 defaultPanicOnValidationErrors,
		TerminationGracePeriodSeconds:           defaultTerminationGracePeriodSeconds,
		Region:                                  defaultRegion,
		LaunchDarklyEnabled:                     defaultLaunchDarklyEnabled,
		LaunchDarklyStreamURI:                   defaultLaunchDarklyStreamURI,
		LaunchDarklySdkKey:                      defaultLaunchDarklySdkKey,
		AuthEnabled:                             defaultAuthEnabled,
		KeysUri:                                 defaultKeysUri,
		AuthJwtAud:                              defaultAuthJwtAud,
		SolaceChannels:                          defaultSolaceChannels,
		SolaceStreamingQueueGroup:               defaultSolaceStreamingQueueGroup,
		AuthS2SJwtAud:                           defaultAuthS2SJwtAud,
		AuthJwtIss:                              defaultAuthJwtIss,
		SecretKeyFile:                           defaultSecretKeyFile,
		TokenURI:                                defaultTokenURI,
		MessagingEnabled:                        defaultMessagingEnabled,
		MessagingPublishBufferSize:              defaultMessagingPublishBufferSize,
		MessagingConnectionCheckIntervalSeconds: defaultMessagingConnectionCheckIntervalSeconds,
		EventsFilePath:                          defaultEventsFilePath,
		SkipPurgeEvents:                         defaultSkipPurgeEvents,
		FeatureFlagsEnabled:                     defaultFeatureFlagsEnabled,
	}
}

func init() {
	loadConfig(Global)

	validator, err := validator.NewValidator(make(map[string]validator.CustomTag))
	if err != nil {
		if Global.PanicOnValidationErrors {
			panic(fmt.Errorf("fatal error creating validator: %s", err))
		}
		fmt.Printf("error creating validator: %s", err) //revive:disable:unhandled-error
	} else {
		if err := validator.ValidateStruct(Global); err != nil {
			if Global.PanicOnValidationErrors {
				panic(fmt.Errorf("fatal error validating config: %s", err))
			}
			fmt.Printf("error validating config: %s", err) //revive:disable:unhandled-error
		}
	}
}

// LoadConfig loads the environment variables into a config Spec struct
func loadConfig(configSpec *Spec) {
	v := viper.New()
	v.AutomaticEnv()

	setDefaults(v, configSpec)
	loadSecrets(v)

	if err := v.Unmarshal(configSpec); err != nil {
		panic(fmt.Errorf("fatal error unmarshalling config %s", err))
	}
}

func loadSecrets(v *viper.Viper) {
	v.SetDefault("launchdarkly_sdk_key", getFromEnvFile("LAUNCHDARKLY_SDK_KEY_FILE", defaultLaunchDarklySdkKey))
}

func getFromEnvFile(envFileVariable, defaultValue string) string {
	label := "config/getFromEnvFile"
	logger := operation.GetGlobalLogger()
	fileName := os.Getenv(envFileVariable)

	if fileName != "" {
		fileContent, err := os.ReadFile(fileName)
		if err != nil {
			logger.Warn("label", label, "message", "could not load secret from file", "envFileVariable", envFileVariable, "fileName", fileName, "error", err)

			return defaultValue
		}

		logger.Info("label", label, "message", "secret loaded from file", "envFileVariable", envFileVariable, "fileName", fileName)

		return strings.TrimSpace(string(fileContent))
	}

	return defaultValue
}

func setDefaults(v *viper.Viper, i any) {
	values := map[string]any{}
	if err := mapstructure.Decode(i, &values); err != nil {
		panic(err)
	}
	for key, defaultValue := range values {
		v.SetDefault(key, defaultValue)
	}
}
