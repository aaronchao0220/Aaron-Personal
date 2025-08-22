package e2e

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/qlik-trial/usage-telemetry-publisher/test/util"
	"github.com/stretchr/testify/require"
)

func TestHealth_Telemetry_Publisher(t *testing.T) {
	t.Run("usage-telemetry-publisher service should have a metric endpoint", func(t *testing.T) {
		res, err := http.Get(fmt.Sprintf("%v/metrics", util.GetUsageTelemetryPublisherServiceAddr()))
		if res != nil {
			defer res.Body.Close()
		}
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, res.StatusCode)
		require.True(t, strings.Contains(res.Header.Get("Content-Type"), "text/plain"), "Content-Type for metrics should be text/plain")
	})
}
