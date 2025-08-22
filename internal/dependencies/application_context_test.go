package dependencies

import (
	"context"
	"testing"

	"github.com/qlik-trial/usage-telemetry-publisher/cmd/config"
	"github.com/qlik-trial/usage-telemetry-publisher/internal/messaging"
)

func TestSubscribeToChannelsCalledTwice(t *testing.T) {
	config.Global.SolaceChannels = "channel1,channel2"
	config.Global.IntermediateStorageEnabled = false

	messagingClientMock := &messaging.MockedMessagingClient{}
	appCtx := ApplicationContext{
		MessagingClient: messagingClientMock,
	}
	messagingClientMock.On("SubscribeEvent", "channel1").Return(nil)
	messagingClientMock.On("SubscribeEvent", "channel2").Return(nil)
	appCtx.subscribeToChannels(context.Background())
	messagingClientMock.AssertExpectations(t)

}
func TestSubscribeToChannelsCalledOnce(t *testing.T) {
	config.Global.SolaceChannels = "channel1"
	config.Global.IntermediateStorageEnabled = false

	messagingClientMock := &messaging.MockedMessagingClient{}
	appCtx := ApplicationContext{
		MessagingClient: messagingClientMock,
	}
	messagingClientMock.On("SubscribeEvent", "channel1").Return(nil)
	appCtx.subscribeToChannels(context.Background())
	messagingClientMock.AssertExpectations(t)

}
