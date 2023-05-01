package main

import (
	"context"
	"fmt"
	"sync"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTClient struct {
	sync.Mutex
	oldClient    mqtt.Client
	feedChannels map[string]bool
}

func NewMQTTClient() MQTTClient {
	hostname := "localhost"
	port := 1883

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", hostname, port))
	client := mqtt.NewClient(opts)

	if t := client.Connect(); t.Wait() && t.Error() != nil {
		panic(t.Error())
	}

	log.WithField("hostname", hostname).WithField("port", port).Info("connected to mqtt broker")

	return MQTTClient{
		oldClient:    client,
		feedChannels: make(map[string]bool),
	}
}

// Establishes a bidirectional channel to the feed
func (c *MQTTClient) Subscribe(ctx context.Context, feed string) (tx chan<- []byte, rx <-chan []byte) {
	log := LoggerFromContext(ctx).WithField("feed", feed)

	// lock and unlock mutex at the end of the scope
	(func() {
		c.Lock()
		defer c.Unlock()
		if _, ok := c.feedChannels[feed]; ok {
			panic("already subscribed to feed " + feed)
		}
		c.feedChannels[feed] = true
	})()

	tch := make(chan []byte)
	rch := make(chan []byte)

	// QOS = 2 because we want to ensure that the message is delivered exactly once
	// for simplicity, not for reliability
	var qos byte = 2
	if t := c.oldClient.Subscribe(feed, qos, func(_ mqtt.Client, m mqtt.Message) {
		log.WithField("payload", string(m.Payload())).Infoln("received payload")
		rch <- m.Payload()
	}); t.Wait() && t.Error() != nil {
		panic(t.Error())
	}

	log.Infoln("subscribed to feed")

	wg.Add(1)
	go func() {
		c.listen(ctx, feed, tch)
		// tx is closed
		delete(c.feedChannels, feed)
		log.Infoln("unsubscribed from feed")
		close(rch)
		wg.Done()
	}()

	return tch, rch
}

func (c *MQTTClient) listen(ctx context.Context, feed string, tx <-chan []byte) {
	log := LoggerFromContext(ctx).WithField("feed", feed)

	for payload := range tx {
		// for simplicity, not for reliability
		var qos byte = 2
		var retained bool = true
		if t := c.oldClient.Publish(feed, qos, retained, payload); t.Wait() && t.Error() != nil {
			panic(t.Error())
		}

		log.WithField("payload", string(payload)).Infoln("published payload to feed")
	}
}
