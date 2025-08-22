package events

import (
	"context"
	"encoding/json"

	"github.com/qlik-trial/go-service-kit/v29/messaging"
	"github.com/qlik-trial/go-service-kit/v29/operation"
	"github.com/qlik-trial/usage-telemetry-publisher/internal/features"
	"github.com/qlik-trial/usage-telemetry-publisher/internal/model"
)

func EventHandler(ctx context.Context, featuresClient features.FeaturesClient) messaging.MsgHandler {
	label := "event_handler/EventHandler"
	return func(msg *messaging.Message) {
		op, ctx := operation.NewOperation(ctx, "handling_event", operation.RecordMetrics(true))
		var err error
		defer func() {
			op.Finish(err)
		}()

		event, unmarshalErr := parseEvent(ctx, msg)
		if unmarshalErr != nil {
			return
		}

		if !validateEvent(ctx, msg, event) {
			return
		}

		operation.Logger(ctx).Debug("label", label, "event", event.Source, "message", "event handled")
		ackWithLog(ctx, msg, label)
	}
}

func parseEvent(ctx context.Context, msg *messaging.Message) (model.CloudEvent, error) {
	label := "event_handler/parseEvent"
	var event model.CloudEvent
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		operation.Logger(ctx).Info(
			"label", label,
			"message", "failed to unmarshal event",
			"error", err,
			"event", string(msg.Data))
		ackWithLog(ctx, msg, label)
		return event, err
	}
	return event, nil
}

func validateEvent(ctx context.Context, msg *messaging.Message, event model.CloudEvent) bool {
	label := "event_handler/validateEvent"
	if !isValidEvent(event) {
		operation.Logger(ctx).Info(
			"label", label,
			"message", "event is not valid",
			"event", event)
		ackWithLog(ctx, msg, label)
		return false
	}

	return true
}

func ackWithLog(ctx context.Context, msg *messaging.Message, label string) {
	ackErr := msg.Ack()
	if ackErr != nil {
		operation.Logger(ctx).Error(
			"label", label,
			"message", "failed to ack event",
			"error", ackErr)
	}
}

func isValidEvent(event model.CloudEvent) bool {
	if event.EventType == "" || event.Time == "" || event.TenantId == "" {
		return false
	}
	return true
}
