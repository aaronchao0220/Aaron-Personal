package formatter

import (
	"testing"

	"github.com/qlik-trial/usage-telemetry-publisher/internal/model"
)

func Test_write_events(t *testing.T) {
	tests := []struct {
		events   []*model.ScrubbedEvent
		expected string
	}{
		{events: []*model.ScrubbedEvent{}, expected: ""},
		{
			events: []*model.ScrubbedEvent{
				{
					Id:       "12345",
					Type:     "com.qlik.v1.some_event",
					Time:     "2023-10-01T12:00:00Z",
					TenantId: "tenant_123",
					Data: map[string]any{
						"foo": map[string]any{"bar": map[string]any{"baz": "qux"}},
					},
				},
			},
			expected: "{\"customerId\":\"tenant_123\",\"dimension.data.foo.bar.baz\":\"qux\",\"eventName\":\"com.qlik.v1.some_event\",\"idempotencyKey\":\"12345\",\"timestamp\":\"2023-10-01T12:00:00Z\"}",
		},
		{
			events: []*model.ScrubbedEvent{
				{
					Id:       "12345",
					Type:     "com.qlik.v1.some_event",
					Time:     "2023-10-01T12:00:00Z",
					TenantId: "tenant_123",
					Data: map[string]any{
						"foo": map[string]any{"bar": map[string]any{"baz": "qux"}},
					},
				},
				{
					Id:       "12345",
					Type:     "com.qlik.v1.some_event",
					Time:     "2023-10-01T12:00:00Z",
					TenantId: "tenant_123",
					Data: map[string]any{
						"foo": "aaron"},
				},
			},
			expected: `{"customerId":"tenant_123","dimension.data.foo.bar.baz":"qux","eventName":"com.qlik.v1.some_event","idempotencyKey":"12345","timestamp":"2023-10-01T12:00:00Z"}
{"customerId":"tenant_123","dimension.data.foo":"aaron","eventName":"com.qlik.v1.some_event","idempotencyKey":"12345","timestamp":"2023-10-01T12:00:00Z"}`,
		},
	}

	for _, test := range tests {
		result := Write(test.events)
		if result != test.expected {
			t.Errorf("expected %q, got %q", test.expected, result)
		}
	}
}
