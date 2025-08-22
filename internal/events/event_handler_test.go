package events

import (
	"testing"

	"github.com/qlik-trial/usage-telemetry-publisher/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestIsValidEvent(t *testing.T) {
	tests := []struct {
		eventType string
		tenantID  string
		time      string
		version   string
		expected  bool
	}{
		{"", "tenant_id", "time", "1.0", false},
		{"com.qlik.v1.audit.purged", "tenant_id", "", "1.0", false},
		{"com.qlik.v1.ai-platform.detect-hallucinations-scan-job.triggered", "", "time", "1.0", false},
		{"com.qlik.v1.ai-platform.detect-hallucinations-scan-job.triggered", "tenant_id", "time", "1.0", true},
	}

	for _, test := range tests {
		t.Run(test.eventType, func(t *testing.T) {
			event := model.CloudEvent{
				EventType:   test.eventType,
				TenantId:    test.tenantID,
				Time:        test.time,
				SpecVersion: test.version,
			}
			result := isValidEvent(event)
			assert.Equal(t, test.expected, result)
		})
	}
}
