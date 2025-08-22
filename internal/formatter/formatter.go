package formatter

import (
	"encoding/json"

	"github.com/qlik-trial/usage-telemetry-publisher/internal/model"
)

// flattenMap recursively flattens nested maps with a given prefix
func flattenMap(data map[string]any, prefix string, result map[string]any) {
	for key, value := range data {
		fullKey := prefix + key

		if nestedMap, ok := value.(map[string]any); ok {
			// Recursively flatten nested maps
			flattenMap(nestedMap, fullKey+".", result)
		} else {
			// Add the leaf value to the result
			result[fullKey] = value
		}
	}
}

// Flatten takes a CloudEvent and returns a flattened representation of its data.
// data should be flattened into dimension single level (dimension.data.<key>)
func Flatten(event *model.ScrubbedEvent) string {
	flattened := make(map[string]any)

	// set standard values
	flattened["idempotencyKey"] = event.Id
	flattened["eventName"] = event.Type
	flattened["timestamp"] = event.Time
	flattened["customerId"] = event.TenantId

	// Recursively flatten the event data
	flattenMap(event.Data, "dimension.data.", flattened)

	b, _ := json.Marshal(flattened)
	return string(b)
}
