package main_test

import (
	"context"
	"log"
	"net"
	"os"
	"testing"
	"time"

	"github.com/dadn-dream-home/x/server/startup"
	"github.com/dadn-dream-home/x/server/telemetry"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	pb "github.com/dadn-dream-home/x/protobuf"
)

var ctx = context.Background()

var startServerAndClient func(ctx context.Context) (pb.BackendServiceClient, func())

func TestMain(m *testing.M) {
	os.Chdir(os.Getenv("WORKSPACE_DIR"))

	ctx = telemetry.InitLogger(ctx)
	config := startup.OpenConfig(ctx)
	config.DatabaseConfig.ConnectionString = ":memory:"
	mqtt := startup.ConnectMQTT(ctx, config.MQTTConfig)

	startServerAndClient = func(ctx context.Context) (pb.BackendServiceClient, func()) {
		db := startup.OpenDatabase(ctx, config.DatabaseConfig)
		startup.Migrate(ctx, db, config.DatabaseConfig)

		lis := bufconn.Listen(1024 * 1024)
		server := startup.NewServer(ctx, db, mqtt)
		go server.Serve(ctx, lis)

		conn, err := grpc.DialContext(ctx, "",
			grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
				return lis.Dial()
			}), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Panic("error connecting to server: %w", err)
		}

		return pb.NewBackendServiceClient(conn), func() {
			server.Stop()
			db.Close()
		}
	}

	// setup
	m.Run()

	mqtt.Disconnect(0)
}

func TestSmoke(t *testing.T) {
	_, stop := startServerAndClient(ctx)
	defer stop()
}

func TestCreateDuplicateFeeds(t *testing.T) {
	client, stop := startServerAndClient(ctx)
	defer stop()

	feedID := uuid.NewString()

	_, err := client.CreateFeed(ctx, &pb.CreateFeedRequest{
		Feed: &pb.Feed{
			Id:   feedID,
			Type: pb.FeedType_TEMPERATURE,
		},
	})
	if err != nil {
		t.Fatalf("error creating feed: %v", err)
	}

	_, err = client.CreateFeed(ctx, &pb.CreateFeedRequest{
		Feed: &pb.Feed{
			Id:   feedID,
			Type: pb.FeedType_TEMPERATURE,
		},
	})
	if err == nil {
		t.Fatalf("expected error creating feed")
	}
}

func TestCreateAndListFeeds(t *testing.T) {
	client, stop := startServerAndClient(ctx)
	defer stop()

	feedID := uuid.NewString()

	_, err := client.CreateFeed(ctx, &pb.CreateFeedRequest{
		Feed: &pb.Feed{
			Id:   feedID,
			Type: pb.FeedType_TEMPERATURE,
		},
	})
	if err != nil {
		t.Fatalf("error creating feed: %v", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	stream, err := client.StreamFeedsChanges(ctx, &pb.StreamFeedsChangesRequest{})
	if err != nil {
		t.Fatalf("error listing feeds: %v", err)
	}

	feeds := make([]*pb.Feed, 0)

	for {
		res, err := stream.Recv()
		if err != nil || res == nil {
			break
		}

		feeds = append(feeds, res.Change.Addeds...)

		if len(res.Change.RemovedIDs) > 0 {
			t.Fatalf("unexpected removed feeds: %v", res.Change.RemovedIDs)
		}
	}

	if len(feeds) != 1 {
		t.Fatalf("expected 1 feed, got %d", len(feeds))
	}

	if feeds[0].Id != feedID {
		t.Fatalf("expected feed id to be test, got %s", feeds[0].Id)
	}
}

func TestCreateAndDeleteFeeds(t *testing.T) {
	client, stop := startServerAndClient(ctx)
	defer stop()

	feedID := uuid.NewString()

	if _, err := client.CreateFeed(ctx, &pb.CreateFeedRequest{
		Feed: &pb.Feed{
			Id:   feedID,
			Type: pb.FeedType_TEMPERATURE,
		},
	}); err != nil {
		t.Fatalf("error creating feed: %v", err)
	}

	if _, err := client.DeleteFeed(ctx, &pb.DeleteFeedRequest{
		Feed: &pb.Feed{
			Id: feedID,
		},
	}); err != nil {
		t.Fatalf("error deleting feed: %v", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	stream, err := client.StreamFeedsChanges(ctx, &pb.StreamFeedsChangesRequest{})
	if err != nil {
		t.Fatalf("error listing feeds: %v", err)
	}

	feeds := make([]*pb.Feed, 0)

	for {
		res, err := stream.Recv()
		if err != nil || res == nil {
			break
		}

		feeds = append(feeds, res.Change.Addeds...)

		if len(res.Change.RemovedIDs) > 0 {
			t.Fatalf("unexpected removed feeds: %v", res.Change.RemovedIDs)
		}
	}

	if len(feeds) != 0 {
		t.Fatalf("expected 0 feeds, got %d", len(feeds))
	}
}