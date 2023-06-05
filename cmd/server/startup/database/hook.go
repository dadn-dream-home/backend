package database

import (
	"context"
	"database/sql"
	"sync"

	"github.com/dadn-dream-home/x/server/state"
	"github.com/dadn-dream-home/x/server/state/topic"
	"github.com/dadn-dream-home/x/server/telemetry"
	"github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

type hooker struct {
	sync.RWMutex

	subscribers map[topic.Topic]map[*subscriber]struct{}
}

type subscriber struct {
	callback func(topic topic.Topic, rowid int64) error
	errCh    chan error
}

func RegisterHook(ctx context.Context, config Config) state.DatabaseHooker {
	var h hooker
	h.subscribers = make(map[topic.Topic]map[*subscriber]struct{})

	sql.Register(config.Driver,
		&sqlite3.SQLiteDriver{
			ConnectHook: func(conn *sqlite3.SQLiteConn) error {
				conn.RegisterUpdateHook(h.handleUpdate)
				return nil
			},
		})

	return &h
}

func (conn *hooker) handleUpdate(op int, database string, table string, rowid int64) {
	go func() {
		ctx := telemetry.InitLogger(context.Background())
		log := telemetry.GetLogger(ctx)
		log.Debug("database update", zap.Int("op", op), zap.String("table", table), zap.Int64("rowid", rowid))

		topic := topic.Topic{Op: op, Table: table}

		conn.RLock()
		defer conn.RUnlock()
		wg := sync.WaitGroup{}
		defer wg.Wait()
		for subscriber := range conn.subscribers[topic] {
			subscriber := subscriber
			wg.Add(1)
			go func() {
				err := subscriber.callback(topic, rowid)
				if err != nil {
					subscriber.errCh <- err
				}
				wg.Done()
			}()
		}
	}()
}

func (c *hooker) Subscribe(
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
