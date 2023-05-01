package main

import (
	"fmt"
	"log"
	"net"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Backend struct {
	mqtt.Client
}

func newBackend() Backend {
	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://localhost:1883")
	client := mqtt.NewClient(opts)

	if t := client.Connect(); t.Wait() && t.Error() != nil {
		panic(t.Error())
	}

	log.Printf("connected to broker\n")

	return Backend{client}
}

func main() {
	backend := newBackend()
	s := NewBackendGrpcService(backend.Client)

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 8080))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("server listening at %v\n", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
