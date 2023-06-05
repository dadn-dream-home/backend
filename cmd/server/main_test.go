package main_test

import (
	"context"
	"log"
	"net"
	"os"
	"testing"
	"time"

	"github.com/dadn-dream-home/x/server/startup"
	"github.com/dadn-dream-home/x/server/startup/database"
	"github.com/dadn-dream-home/x/server/telemetry"
	mqtt "github.com/eclipse/paho.mqtt.golang"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	pb "github.com/dadn-dream-home/x/protobuf"
)

var ctx = context.Background()

var startServerAndClient func(ctx context.Context) (pb.BackendServiceClient, mqtt.Client, func())

func TestMain(m *testing.M) {
	os.Chdir(os.Getenv("WORKSPACE_DIR"))

	ctx = telemetry.InitLogger(ctx)
	config := startup.OpenConfig(ctx)
	config.DatabaseConfig.ConnectionString = "file::memory:?cache=shared&_mutex=full"
	config.DatabaseConfig.Driver = uuid.NewString()
	hooker := database.RegisterHook(ctx, config.DatabaseConfig)
	mqttClient := startup.ConnectMQTT(ctx, config.MQTTConfig)

	startServerAndClient = func(ctx context.Context) (pb.BackendServiceClient, mqtt.Client, func()) {
		db := database.OpenDatabase(ctx, config.DatabaseConfig)
		database.Migrate(ctx, db, config.DatabaseConfig)

		lis := bufconn.Listen(1024 * 1024)
		server := startup.NewServer(ctx, db, mqttClient, hooker)
		go server.Serve(ctx, lis)

		conn, err := grpc.DialContext(ctx, "",
			grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
				return lis.Dial()
			}), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Panic("error connecting to server: %w", err)
		}

		return pb.NewBackendServiceClient(conn), mqttClient, func() {
			server.Stop()
			db.Close()
		}
	}

	// setup
	m.Run()

	mqttClient.Disconnect(0)
}

func TestSmoke(t *testing.T) {
	_, _, stop := startServerAndClient(ctx)
	defer stop()
}

func TestCreateDuplicateFeeds(t *testing.T) {
	client, _, stop := startServerAndClient(ctx)
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
	client, _, stop := startServerAndClient(ctx)
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
	client, _, stop := startServerAndClient(ctx)
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

func TestGetFeedConfig(t *testing.T) {
	client, _, stop := startServerAndClient(ctx)
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

	if _, err := client.GetFeedConfig(ctx, &pb.GetFeedConfigRequest{Feed: &pb.Feed{Id: feedID}}); err != nil {
		t.Fatalf("error getting feed config: %v", err)
	}
}

func TestGetFeedConfigNotFound(t *testing.T) {
	client, _, stop := startServerAndClient(ctx)
	defer stop()

	if _, err := client.GetFeedConfig(ctx, &pb.GetFeedConfigRequest{Feed: &pb.Feed{Id: "notfound"}}); err == nil {
		t.Fatalf("expected error getting feed config")
	}
}

func TestUpdateFeedConfig(t *testing.T) {
	client, _, stop := startServerAndClient(ctx)
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

	if _, err := client.UpdateFeedConfig(ctx, &pb.UpdateFeedConfigRequest{
		Config: &pb.Config{
			FeedConfig: &pb.Feed{
				Id:   feedID,
				Type: pb.FeedType_TEMPERATURE,
			},
			TypeConfig: &pb.Config_SensorConfig{
				SensorConfig: &pb.SensorConfig{
					HasNotification: true,
					LowerThreshold: &pb.Threshold{
						Threshold: 10,
					},
					UpperThreshold: &pb.Threshold{
						Threshold: 40,
					},
				},
			},
		},
	}); err != nil {
		t.Fatalf("error updating feed config: %v", err)
	}
}

func TestLowThresholdNotification(t *testing.T) {
	client, mqtt, stop := startServerAndClient(ctx)
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

	if _, err := client.UpdateFeedConfig(ctx, &pb.UpdateFeedConfigRequest{
		Config: &pb.Config{
			FeedConfig: &pb.Feed{
				Id:   feedID,
				Type: pb.FeedType_TEMPERATURE,
			},
			TypeConfig: &pb.Config_SensorConfig{
				SensorConfig: &pb.SensorConfig{
					HasNotification: true,
					LowerThreshold: &pb.Threshold{
						Threshold: 2,
					},
					UpperThreshold: &pb.Threshold{
						Threshold: 98,
					},
				},
			},
		},
	}); err != nil {
		t.Fatalf("error updating feed config: %v", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	stream, err := client.StreamNotifications(ctx, &pb.StreamNotificationsRequest{
		Feed: &pb.Feed{
			Id: feedID,
		},
	})
	if err != nil {
		t.Fatalf("error streaming notifications: %v", err)
	}

	if token := mqtt.Publish(feedID, 0, false, "1"); token.Wait() && token.Error() != nil {
		t.Fatalf("error publishing to mqtt: %v", token.Error())
	}
	if token := mqtt.Publish(feedID, 0, false, "30"); token.Wait() && token.Error() != nil {
		t.Fatalf("error publishing to mqtt: %v", token.Error())
	}
	if token := mqtt.Publish(feedID, 0, false, "99"); token.Wait() && token.Error() != nil {
		t.Fatalf("error publishing to mqtt: %v", token.Error())
	}

	var notifications []*pb.Notification

	for {
		res, err := stream.Recv()
		if err != nil || res == nil {
			break
		}

		notifications = append(notifications, res.Notification)
	}

	if len(notifications) != 2 {
		t.Fatalf("expected 2 notification, got %d", len(notifications))
	}

	if notifications[0].Feed.Id != feedID {
		t.Fatalf("expected feed id %s, got %s", feedID, notifications[0].Feed.Id)
	}
}

func TestThresholdDisableNotification(t *testing.T) {
	client, mqtt, stop := startServerAndClient(ctx)
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

	if _, err := client.UpdateFeedConfig(ctx, &pb.UpdateFeedConfigRequest{
		Config: &pb.Config{
			FeedConfig: &pb.Feed{
				Id:   feedID,
				Type: pb.FeedType_TEMPERATURE,
			},
			TypeConfig: &pb.Config_SensorConfig{
				SensorConfig: &pb.SensorConfig{
					HasNotification: false,
					LowerThreshold: &pb.Threshold{
						Threshold: 2,
					},
					UpperThreshold: &pb.Threshold{
						Threshold: 98,
					},
				},
			},
		},
	}); err != nil {
		t.Fatalf("error updating feed config: %v", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	stream, err := client.StreamNotifications(ctx, &pb.StreamNotificationsRequest{
		Feed: &pb.Feed{
			Id: feedID,
		},
	})
	if err != nil {
		t.Fatalf("error streaming notifications: %v", err)
	}

	if token := mqtt.Publish(feedID, 0, false, "1"); token.Wait() && token.Error() != nil {
		t.Fatalf("error publishing to mqtt: %v", token.Error())
	}
	if token := mqtt.Publish(feedID, 0, false, "30"); token.Wait() && token.Error() != nil {
		t.Fatalf("error publishing to mqtt: %v", token.Error())
	}
	if token := mqtt.Publish(feedID, 0, false, "99"); token.Wait() && token.Error() != nil {
		t.Fatalf("error publishing to mqtt: %v", token.Error())
	}

	var notifications []*pb.Notification

	for {
		res, err := stream.Recv()
		if err != nil || res == nil {
			break
		}

		notifications = append(notifications, res.Notification)
	}

	if len(notifications) != 0 {
		t.Fatalf("expected 0 notification, got %d", len(notifications))
	}
}

// tests for feed values

func TestFeedValue(t *testing.T) {
	client, mqtt, stop := startServerAndClient(ctx)
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

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	stream, err := client.StreamSensorValues(ctx, &pb.StreamSensorValuesRequest{
		Feed: &pb.Feed{
			Id: feedID,
		},
	})
	if err != nil {
		t.Fatalf("error streaming feed values: %v", err)
	}

	if token := mqtt.Publish(feedID, 0, false, "1"); token.Wait() && token.Error() != nil {
		t.Fatalf("error publishing to mqtt: %v", token.Error())
	}
	if token := mqtt.Publish(feedID, 0, false, "30"); token.Wait() && token.Error() != nil {
		t.Fatalf("error publishing to mqtt: %v", token.Error())
	}
	if token := mqtt.Publish(feedID, 0, false, "99"); token.Wait() && token.Error() != nil {
		t.Fatalf("error publishing to mqtt: %v", token.Error())
	}

	var feedValues []float32
	for {
		res, err := stream.Recv()
		if err != nil || res == nil {
			break
		}

		feedValues = append(feedValues, res.Value)
	}

	if len(feedValues) != 3 {
		t.Fatalf("expected 3 feed values, got %d", len(feedValues))
	}
}

func TestActivityLog(t *testing.T) {
	client, mqtt, stop := startServerAndClient(ctx)
	defer stop()

	feedID := uuid.NewString()

	if _, err := client.CreateFeed(ctx, &pb.CreateFeedRequest{
		Feed: &pb.Feed{
			Id:   feedID,
			Type: pb.FeedType_LIGHT,
		},
	}); err != nil {
		t.Fatalf("error creating feed: %v", err)
	}

	if token := mqtt.Publish(feedID, 0, false, "1"); token.Wait() && token.Error() != nil {
		t.Fatalf("error publishing to mqtt: %v", token.Error())
	}
	if token := mqtt.Publish(feedID, 0, false, "0"); token.Wait() && token.Error() != nil {
		t.Fatalf("error publishing to mqtt: %v", token.Error())
	}

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	stream, err := client.StreamActivities(ctx, &pb.StreamActivitiesRequest{})
	if err != nil {
		t.Fatalf("error streaming activity log: %v", err)
	}

	if token := mqtt.Publish(feedID, 0, false, "1"); token.Wait() && token.Error() != nil {
		t.Fatalf("error publishing to mqtt: %v", token.Error())
	}

	var activities []*pb.Activity
	for {
		res, err := stream.Recv()
		if err != nil || res == nil {
			break
		}

		activities = append(activities, res.Activity)
	}

	t.Logf("activities: %v", activities)

	if len(activities) != 3 {
		t.Fatalf("expected 3 activity log, got %d", len(activities))
	}
}

func TestActivityLogWithTriggerConfig(t *testing.T) {
	client, mqtt, stop := startServerAndClient(ctx)
	defer stop()

	feedA := &pb.Feed{
		Id:   uuid.NewString(),
		Type: pb.FeedType_TEMPERATURE,
	}

	if _, err := client.CreateFeed(ctx, &pb.CreateFeedRequest{
		Feed: feedA,
	}); err != nil {
		t.Fatalf("error creating feed: %v", err)
	}

	feedB := &pb.Feed{
		Id:   uuid.NewString(),
		Type: pb.FeedType_LIGHT,
	}

	if _, err := client.CreateFeed(ctx, &pb.CreateFeedRequest{
		Feed: feedB,
	}); err != nil {
		t.Fatalf("error creating feed: %v", err)
	}

	if _, err := client.UpdateFeedConfig(ctx, &pb.UpdateFeedConfigRequest{
		Config: &pb.Config{
			FeedConfig: feedA,
			TypeConfig: &pb.Config_SensorConfig{
				SensorConfig: &pb.SensorConfig{
					HasNotification: true,
					LowerThreshold: &pb.Threshold{
						Threshold:  10,
						HasTrigger: true,
						Feed:       feedB,
						State:      true,
					},
					UpperThreshold: &pb.Threshold{
						Threshold:  90,
						HasTrigger: true,
						Feed:       feedB,
						State:      false,
					},
				},
			},
		},
	}); err != nil {
		t.Fatalf("error updating feed config: %v", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	stream, err := client.StreamActivities(ctx, &pb.StreamActivitiesRequest{})
	if err != nil {
		t.Fatalf("error streaming activity log: %v", err)
	}

	if token := mqtt.Publish(feedA.Id, 0, false, "0"); token.Wait() && token.Error() != nil {
		t.Fatalf("error publishing to mqtt: %v", token.Error())
	}

	var activities []*pb.Activity
	for {
		res, err := stream.Recv()
		if err != nil || res == nil {
			break
		}

		activities = append(activities, res.Activity)
	}

	t.Logf("activities: %v", activities)

	if len(activities) != 1 {
		t.Fatalf("expected 1 activity log, got %d", len(activities))
	}
}
