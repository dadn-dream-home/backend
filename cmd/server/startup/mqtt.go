package startup

import (
	"context"

	"github.com/dadn-dream-home/x/server/telemetry"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/zap"
)

func ConnectMQTT(ctx context.Context, config MQTTConfig) mqtt.Client {
	log := telemetry.GetLogger(ctx)

	options := mqtt.NewClientOptions()

	for _, broker := range config.Brokers {
		options.AddBroker(broker)

		log.Info("added mqtt broker", zap.String("broker", broker))
	}

	client := mqtt.NewClient(options)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal("failed to connect to mqtt brokers", zap.Error(token.Error()))
	}

	log.Info("connected to mqtt brokers")

	return client
}
