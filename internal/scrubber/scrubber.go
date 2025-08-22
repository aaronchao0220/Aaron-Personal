package scrubber

import (
	"github.com/qlik-trial/go-service-kit/v29/messaging/events"
	"github.com/qlik-trial/usage-telemetry-publisher/internal/model"
)

// ScrubEvent removes sensitive information from a CloudEvent and returns a ScrubbedEvent.
func ScrubEvent(event events.CloudEvent) model.ScrubbedEvent {
	return model.ScrubbedEvent{
		Id:                 event.Id,
		SpecVersion:        event.SpecVersion,
		TenantId:           event.TenantID,
		Source:             event.Source,
		UserId:             event.UserID,
		SessionId:          event.SessionID,
		Type:               event.Type,
		Time:               event.Time,
		Host:               event.Host,
		OriginIp:           event.OriginIP,
		OwnerId:            event.OwnerID,
		TopLevelResourceId: event.TopLevelResourceID,
		SpaceId:            event.SpaceID,
		ClientId:           event.ClientID,
		Reason:             event.Reason,
		Data:               event.Data.(map[string]any),
	}
}
