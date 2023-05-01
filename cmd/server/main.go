package main

import (
	"log"
	"math/rand"
	"sync"

	pb "github.com/dadn-dream-home/x/protobuf"
)

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

	db := NewDatabase()
	defer db.Close()
	
	feed := db.InsertFeed("test feed", pb.FeedType_TEMPERATURE)

	for i := 0; i < 10; i++ {
		db.InsertFeedValue(feed.Id, rand.Float32()*100.0)
	}

	log.Printf("feed values: %v\n", db.SummariseFeedValues(feed.Id))

	wg.Wait()
}
