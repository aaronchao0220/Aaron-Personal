package messaging

import (
	"context"
	"time"

	"github.com/qlik-trial/go-service-kit/v29/healthcheck"
	"github.com/qlik-trial/go-service-kit/v29/messaging"
	"github.com/qlik-trial/go-service-kit/v29/operation"
	"github.com/qlik-trial/usage-telemetry-publisher/cmd/config"
	"github.com/qlik-trial/usage-telemetry-publisher/internal/auth"
)

// Client is a messaging client used to interface with a messaging client
type Client struct {
	*messaging.Client
}

const (
	solaceAudience = "qlik.api.internal/messaging"
)

type EventListener interface {
	SubscribeEvent(subject string, qgroup string, cb messaging.MsgHandler) error
	AddReadinessCheck(healthcheck.Handler)
	Close()
	Connect(<-chan struct{}) error
}

// CreateClient creates a messaging Client instance
func CreateClient(ctx context.Context, tokenGenerator auth.TokenGenerator, clientID string) (EventListener, error) {
	// Token handler funcs
	label := "messaging/client/CreateClient"
	var solaceTokenHandler = func() string {
		token, err := tokenGenerator.GetServiceToServiceJWT(context.Background(), solaceAudience, "")
		if err != nil {
			operation.Logger(context.Background()).Error("label", label, "message", "failed to generate service-to-service JWT for messaging", "error", err)
		}
		return token
	}

	options := []messaging.Option{
		messaging.EnableMetrics(true),
	}
	if config.Global.AuthEnabled {
		options = append(options, messaging.WithSolaceTokenHandler(solaceTokenHandler))
		options = append(options, messaging.WithMaxMsgBufferSize(int32(config.Global.MessagingPublishBufferSize)))
	}

	// Set up messaging client
	messagingClient, err := messaging.NewClient(
		clientID,
		config.ServiceName,
		options...,
	)
	if err != nil {
		return nil, err
	}
	return &Client{messagingClient}, err
}

// Connect will attempt to connect until successful or is stop to stop via the provided channel
func (mc *Client) Connect(stopSignal <-chan struct{}) error {
	label := "messaging/client/Connect"
	messagingCtx, messagingCtxCancel := context.WithCancel(context.Background())
	go func() {
		<-stopSignal
		messagingCtxCancel()
	}()

	operation.Logger(context.Background()).Info("label", label, "message", "Connecting to messaging provider...")
	err := mc.ConnectWithCtx(messagingCtx, time.Duration(config.Global.MessagingConnectionCheckIntervalSeconds)*time.Second)

	if err != nil {
		operation.Logger(context.Background()).Error("label", label, "message", "Failed to connect to messaging provider", "error", err)
		return err
	}
	operation.Logger(context.Background()).Info("label", label, "message", "Connected to messaging provider")
	return nil
}

// CloseWithChan instructs the connection to messaging to be closed and will close
// the returned channel when it has completed closing
func (mc *Client) CloseWithChan() <-chan struct{} {
	done := make(chan struct{})
	go func() {
		mc.Close()
		close(done)
	}()
	return done
}

// SubscribeEvent subscribes to the specified STAN queue with the supplied Receiver callback
// A durable queue group allows you to have all members leave but still maintain state.
// When a member re-joins, it starts at the last position in that group.
func (mc *Client) SubscribeEvent(subject, qgroup string, cb messaging.MsgHandler) error {
	subOpts := []messaging.SubscriptionOption{}
	subOpts = append(subOpts, messaging.SetManualAckMode())
	return mc.Subscribe(messaging.NewQueueSubscription(subject+"/>", qgroup, cb, subOpts...)) //nolint:staticcheck
}
