package startup

import (
	"context"

	"github.com/dadn-dream-home/x/server/telemetry"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func ConnectMQTT(ctx context.Context, config MQTTConfig) mqtt.Client {
	log := telemetry.GetLogger(ctx)

	options := mqtt.NewClientOptions()

	for _, broker := range config.Brokers {
		options.AddBroker(broker)
		
		log.WithField("broker", broker).
			Infof("added mqtt broker")
	}

	client := mqtt.NewClient(options)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.WithError(token.Error()).
			Fatal("failed to connect to mqtt brokers")
	}

	log.Infof("connected to mqtt brokers")

	return client
}
