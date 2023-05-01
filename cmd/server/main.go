package main

import "sync"

var wg sync.WaitGroup

func main() {
	// backend := newBackend()
	// s := NewBackendGrpcService(backend.Client)

	// lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 8080))
	// if err != nil {
	// 	log.Fatalf("failed to listen: %v", err)
	// }

	// log.Printf("server listening at %v\n", lis.Addr())
	// if err := s.Serve(lis); err != nil {
	// 	log.Fatalf("failed to serve: %v", err)
	// }

	mqtt := NewMQTTClient()
	tx, _ := mqtt.Subscribe("test")
	close(tx)

	wg.Wait()
}
