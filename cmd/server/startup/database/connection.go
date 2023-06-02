package database

import (
	"context"
	"database/sql"
	"sync"

	"github.com/dadn-dream-home/x/server/state"
	"github.com/dadn-dream-home/x/server/telemetry"
	"github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

type HookedConnection struct {
	sync.RWMutex
	*sql.Conn

	subscribers map[state.Topic]map[*subscriber]struct{}
}

type subscriber struct {
	callback func(topic state.Topic, rowid int64) any
}

func NewHookedConnection(ctx context.Context, db *sql.DB) *HookedConnection {
	log := telemetry.GetLogger(ctx)

	conn := &HookedConnection{
		subscribers: make(map[state.Topic]map[*subscriber]struct{}),
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
	topic := state.Topic{Op: op, Table: table}

	conn.RLock()
	defer conn.RUnlock()
	for subscriber := range conn.subscribers[topic] {
		subscriber.callback(topic, rowid)
	}
}

func (c *HookedConnection) Subscribe(t state.Topic, cb func(t state.Topic, r int64) any) (unsubscribe func()) {
	c.Lock()
	defer c.Unlock()

	if c.subscribers[t] == nil {
		c.subscribers[t] = make(map[*subscriber]struct{})
	}

	s := &subscriber{cb}

	c.subscribers[t][s] = struct{}{}

	return func() {
		c.Lock()
		defer c.Unlock()

		delete(c.subscribers[t], s)
	}
}
