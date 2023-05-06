package services

import (
	"context"
	"fmt"

	"github.com/dadn-dream-home/x/server/telemetry"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func NewMQTTClient(ctx context.Context) (mqtt.Client, error) {
	log := telemetry.GetLogger(ctx)

	options := mqtt.NewClientOptions()
	broker := "tcp://localhost:1883"
	options.AddBroker(broker)
	client := mqtt.NewClient(options)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, fmt.Errorf("failed to connect to mqtt broker at %s: %w", broker, token.Error())
	}

	log.WithField("broker", broker).Infof("connected to mqtt broker\n")

	return client, nil
}
