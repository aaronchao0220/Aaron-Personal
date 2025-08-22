package formatter

import (
	"testing"

	"github.com/qlik-trial/usage-telemetry-publisher/internal/model"
)

func TestShouldFlattenEvent(t *testing.T) {
	tests := []struct {
		event    *model.ScrubbedEvent
		expected string
	}{
		{
			event: &model.ScrubbedEvent{
				Id:       "12345",
				Type:     "com.qlik.v1.some_event",
				Time:     "2023-10-01T12:00:00Z",
				TenantId: "tenant_123",
				Data: map[string]any{
					"user": "john_doe",
					"age":  30,
				},
			},
			expected: `{"customerId":"tenant_123","dimension.data.age":30,"dimension.data.user":"john_doe","eventName":"com.qlik.v1.some_event","idempotencyKey":"12345","timestamp":"2023-10-01T12:00:00Z"}`,
		},
		{
			event: &model.ScrubbedEvent{
				Id:       "12345",
				Type:     "com.qlik.v1.some_event",
				Time:     "2023-10-01T12:00:00Z",
				TenantId: "tenant_123",
				Data: map[string]any{
					"foo": "bar",
				},
			},
			expected: `{"customerId":"tenant_123","dimension.data.foo":"bar","eventName":"com.qlik.v1.some_event","idempotencyKey":"12345","timestamp":"2023-10-01T12:00:00Z"}`,
		},
		{
			event: &model.ScrubbedEvent{
				Id:       "12345",
				Type:     "com.qlik.v1.some_event",
				Time:     "2023-10-01T12:00:00Z",
				TenantId: "tenant_123",
				Data: map[string]any{
					"foo": map[string]any{"bar": "baz"},
				},
			},
			expected: `{"customerId":"tenant_123","dimension.data.foo.bar":"baz","eventName":"com.qlik.v1.some_event","idempotencyKey":"12345","timestamp":"2023-10-01T12:00:00Z"}`,
		},
		{
			event: &model.ScrubbedEvent{
				Id:       "12345",
				Type:     "com.qlik.v1.some_event",
				Time:     "2023-10-01T12:00:00Z",
				TenantId: "tenant_123",
				Data: map[string]any{
					"foo": map[string]any{"bar": map[string]any{"baz": "qux"}},
				},
			},
			expected: `{"customerId":"tenant_123","dimension.data.foo.bar.baz":"qux","eventName":"com.qlik.v1.some_event","idempotencyKey":"12345","timestamp":"2023-10-01T12:00:00Z"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.expected, func(t *testing.T) {
			result := Flatten(test.event)
			if result != test.expected {
				t.Errorf("expected %q, got %q", test.expected, result)
			}
		})
	}
}
