package messaging

import (
	"github.com/qlik-trial/go-service-kit/v29/healthcheck"
	"github.com/qlik-trial/go-service-kit/v29/messaging"
	"github.com/stretchr/testify/mock"
)

type MockedMessagingClient struct {
	mock.Mock
}

func (client *MockedMessagingClient) SubscribeEvent(subject string, qgroup string, cb messaging.MsgHandler) error {
	returns := client.Called(subject)
	retval, _ := returns[0].(error)
	return retval
}
func (client *MockedMessagingClient) AddReadinessCheck(check healthcheck.Handler) {
	client.Called(check)
}
func (client *MockedMessagingClient) Close() {
	client.Called()
}
func (client *MockedMessagingClient) Connect(channel <-chan struct{}) error {
	returns := client.Called(channel)
	retval, _ := returns[0].(error)
	return retval
}
