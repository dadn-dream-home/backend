package state

import (
	"context"
	"database/sql"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type State interface {
	MQTT() mqtt.Client
	DB() *sql.DB
	PubSub() PubSub
}

type PubSub interface {
	Subscribe(ctx context.Context, feed string, id string, ch chan<- []byte) (done <-chan struct{}, err error)
	Unsubscribe(ctx context.Context, feed string, id string) error
	AddFeed(ctx context.Context, feed string) error
	RemoveFeed(ctx context.Context, feed string) error
}
