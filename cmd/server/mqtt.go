package main

import (
	"log"
	"sync"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTClient struct {
	sync.Mutex
	oldClient     mqtt.Client
	topicChannels map[string]bool
}

func NewMQTTClient() MQTTClient {
	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://localhost:1883")
	client := mqtt.NewClient(opts)

	if t := client.Connect(); t.Wait() && t.Error() != nil {
		panic(t.Error())
	}

	return MQTTClient{
		oldClient:     client,
		topicChannels: make(map[string]bool),
	}
}

// Establishes a bidirectional channel to the topic
func (c *MQTTClient) Subscribe(topic string) (tx chan<- []byte, rx <-chan []byte) {

	// lock and unlock mutex at the end of the scope
	(func() {
		c.Lock()
		defer c.Unlock()
		if _, ok := c.topicChannels[topic]; ok {
			panic("already subscribed to topic " + topic)
		}
		c.topicChannels[topic] = true
	})()

	tch := make(chan []byte)
	rch := make(chan []byte)

	// QOS = 2 because we want to ensure that the message is delivered exactly once
	// for simplicity, not for reliability
	var qos byte = 2
	if t := c.oldClient.Subscribe(topic, qos, func(_ mqtt.Client, m mqtt.Message) {
		log.Printf("received topic %v: %s\n", topic, string(m.Payload()))
		rch <- m.Payload()
	}); t.Wait() && t.Error() != nil {
		panic(t.Error())
	}

	wg.Add(1)
	go func() {
		c.listen(topic, tch)
		// tx is closed
		delete(c.topicChannels, topic)
		log.Printf("unsubscribed from topic %v\n", topic)
		close(rch)
		wg.Done()
	}()

	log.Printf("subscribed to topic %v\n", topic)

	return tch, rch
}

func (c *MQTTClient) listen(topic string, tx <-chan []byte) {
	for payload := range tx {
		// for simplicity, not for reliability
		var qos byte = 2
		var retained bool = true
		if t := c.oldClient.Publish(topic, qos, retained, payload); t.Wait() && t.Error() != nil {
			panic(t.Error())
		}

		log.Printf("published topic %v: %s\n", topic, string(payload))
	}
}
