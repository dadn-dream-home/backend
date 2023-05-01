package main

import (
	"sync"
)

var wg sync.WaitGroup

func main() {
	service := NewBackendGrpcService()
	service.Init()
	service.Serve()
	wg.Wait()
}
