package util

import (
	"log"

	"github.com/qlik-trial/go-service-kit/v29/messaging"
)

type MsgConfig struct {
	ClientID    string
	ServiceName string
}

func MessagingConnect(cfg *MsgConfig) (*messaging.Client, error) {
	msg, err := messaging.NewClient(cfg.ClientID, cfg.ServiceName)
	if err != nil {
		return nil, err
	}
	err = msg.Connect()
	if err != nil {
		log.Fatalf("could not connect to messaging. err: %v", err)
	}
	return msg, nil
}
