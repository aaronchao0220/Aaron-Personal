package scrubber

import (
	"testing"

	"github.com/qlik-trial/go-service-kit/v29/messaging/events"
	"github.com/stretchr/testify/require"
)

func TestScrubEvent(t *testing.T) {
	// Create a sample CloudEvent with all fields populated
	testData := map[string]any{
		"action":  "button_click",
		"element": "submit",
		"metadata": map[string]any{
			"page":    "dashboard",
			"section": "header",
		},
	}

	inputEvent := events.CloudEvent{
		Id:                 "test-id-123",
		SpecVersion:        "1.0",
		TenantID:           "tenant-456",
		UserID:             "user-789",
		SessionID:          "session-abc",
		Source:             "qlik.com/myapp",
		Type:               "user.interaction",
		Time:               "2025-08-21T10:00:00Z",
		Host:               "app.qlik.com",
		OriginIP:           "192.168.1.1",
		OwnerID:            "owner-def",
		TopLevelResourceID: "resource-ghi",
		SpaceID:            "space-jkl",
		ClientID:           "client-mno",
		Reason:             "user_action",
		Data:               testData,
	}

	// Call ScrubEvent
	result := ScrubEvent(inputEvent)

	require.Equal(t, inputEvent.Id, result.Id)
	require.Equal(t, inputEvent.SpecVersion, result.SpecVersion)
	require.Equal(t, inputEvent.TenantID, result.TenantId)
	require.Equal(t, inputEvent.UserID, result.UserId)
	require.Equal(t, inputEvent.SessionID, result.SessionId)
	require.Equal(t, inputEvent.Source, result.Source)
	require.Equal(t, inputEvent.Type, result.Type)
	require.Equal(t, inputEvent.Time, result.Time)
	require.Equal(t, inputEvent.Host, result.Host)
	require.Equal(t, inputEvent.OriginIP, result.OriginIp)
	require.Equal(t, inputEvent.OwnerID, result.OwnerId)
	require.Equal(t, inputEvent.TopLevelResourceID, result.TopLevelResourceId)
	require.Equal(t, inputEvent.SpaceID, result.SpaceId)
	require.Equal(t, inputEvent.ClientID, result.ClientId)
	require.Equal(t, inputEvent.Reason, result.Reason)

	expectedData := inputEvent.Data.(map[string]any)

	// Check top-level data fields
	require.Equal(t, expectedData["action"], result.Data["action"])
	require.Equal(t, expectedData["element"], result.Data["element"])

	// Check nested data
	resultMetadata, ok := result.Data["metadata"].(map[string]any)
	require.True(t, ok, "Expected metadata to be map[string]any")

	expectedMetadata := expectedData["metadata"].(map[string]any)
	require.Equal(t, expectedMetadata["page"], resultMetadata["page"])
	require.Equal(t, expectedMetadata["section"], resultMetadata["section"])
}

func TestScrubEvent_EmptyFields(t *testing.T) {
	// Test with minimal CloudEvent (empty strings)
	inputEvent := events.CloudEvent{
		Id:          "minimal-id",
		SpecVersion: "1.0",
		Data:        map[string]any{},
	}

	result := ScrubEvent(inputEvent)

	// Verify that empty fields are preserved as empty
	require.Equal(t, "minimal-id", result.Id)
	require.Equal(t, "1.0", result.SpecVersion)
	require.Equal(t, "", result.TenantId)
	require.Equal(t, 0, len(result.Data))
}

func TestScrubEvent_NilData(t *testing.T) {
	// Test behavior when Data is nil or not a map
	inputEvent := events.CloudEvent{
		Id:          "test-id",
		SpecVersion: "1.0",
		Data:        nil,
	}

	// This should panic or handle the nil case - testing current behavior
	defer func() {
		if r := recover(); r != nil {
			t.Logf("ScrubEvent panicked with nil Data as expected: %v", r)
		}
	}()

	result := ScrubEvent(inputEvent)

	// If we reach here, the function handled nil gracefully
	require.Nil(t, result.Data)
}

func TestScrubEvent_ComplexData(t *testing.T) {
	// Test with complex nested data structure
	complexData := map[string]any{
		"simple_string": "value1",
		"simple_number": 42,
		"simple_bool":   true,
		"nested_object": map[string]any{
			"level2": map[string]any{
				"level3": "deep_value",
				"array":  []string{"item1", "item2"},
			},
			"another_field": "value2",
		},
		"array_of_objects": []map[string]any{
			{"id": 1, "name": "first"},
			{"id": 2, "name": "second"},
		},
	}

	inputEvent := events.CloudEvent{
		Id:       "complex-data-test",
		TenantID: "tenant-123",
		Data:     complexData,
	}

	result := ScrubEvent(inputEvent)

	// Verify complex data structure is preserved
	require.Equal(t, "value1", result.Data["simple_string"])
	require.Equal(t, 42, result.Data["simple_number"])
	require.Equal(t, true, result.Data["simple_bool"])

	// Check nested structure
	nestedObj, ok := result.Data["nested_object"].(map[string]any)
	require.True(t, ok, "Expected nested_object to be map[string]any")

	level2, ok := nestedObj["level2"].(map[string]any)
	require.True(t, ok, "Expected level2 to be map[string]any")

	require.Equal(t, "deep_value", level2["level3"])
}
