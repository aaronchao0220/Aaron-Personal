package e2e

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/qlik-trial/usage-telemetry-publisher/test/util"
	"github.com/stretchr/testify/require"
)

func TestLiveEndpoint(t *testing.T) {
	t.Run("service service should be live", func(t *testing.T) {
		res, err := http.Get(fmt.Sprintf("%v/health", util.GetUsageTelemetryPublisherServiceAddr()))
		if res != nil {
			defer res.Body.Close()
		}
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, res.StatusCode)
	})
}

func TestReadyEndpoint(t *testing.T) {
	t.Run("service service should be ready", func(t *testing.T) {
		res, err := http.Get(fmt.Sprintf("%v/ready", util.GetUsageTelemetryPublisherServiceAddr()))
		if res != nil {
			defer res.Body.Close()
		}
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, res.StatusCode)
	})
}
