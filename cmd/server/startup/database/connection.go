package database

import (
	"context"
	"database/sql"
	"sync"

	"github.com/dadn-dream-home/x/server/state/topic"
	"github.com/dadn-dream-home/x/server/telemetry"
	"github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

type HookedConnection struct {
	sync.RWMutex
	*sql.Conn

	subscribers map[topic.Topic]map[*subscriber]struct{}
}

type subscriber struct {
	callback func(topic topic.Topic, rowid int64) error
	errCh    chan error
}

func NewHookedConnection(ctx context.Context, db *sql.DB) *HookedConnection {
	log := telemetry.GetLogger(ctx)

	conn := &HookedConnection{
		subscribers: make(map[topic.Topic]map[*subscriber]struct{}),
	}
	var err error

	conn.Conn, err = db.Conn(ctx)
	if err != nil {
		log.Fatal("failed to get connection from database", zap.Error(err))
	}

	conn.Raw(func(driverConn interface{}) error {
		driverConn.(*sqlite3.SQLiteConn).
			RegisterUpdateHook(conn.UpdateHook)
		return nil
	})

	return conn
}

func (conn *HookedConnection) UpdateHook(op int, database string, table string, rowid int64) {
	topic := topic.Topic{Op: op, Table: table}

	conn.RLock()
	defer conn.RUnlock()
	for subscriber := range conn.subscribers[topic] {
		err := subscriber.callback(topic, rowid)
		if err != nil {
			subscriber.errCh <- err
		}
	}
}

func (c *HookedConnection) Subscribe(
	cb func(t topic.Topic, r int64) error,
	ts ...topic.Topic,
) (unsubscribe func(), errCh <-chan error) {
	c.Lock()
	defer c.Unlock()

	s := &subscriber{cb, make(chan error)}

	for _, t := range ts {
		if c.subscribers[t] == nil {
			c.subscribers[t] = make(map[*subscriber]struct{})
		}

		c.subscribers[t][s] = struct{}{}
	}

	return func() {
		c.Lock()
		defer c.Unlock()

		close(s.errCh)
		for _, t := range ts {
			delete(c.subscribers[t], s)
		}
	}, s.errCh
}
