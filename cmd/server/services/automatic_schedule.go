package services

import (
	"context"
	"sync"

	"github.com/dadn-dream-home/x/server/state"
	"github.com/dadn-dream-home/x/server/state/topic"

	pb "github.com/dadn-dream-home/x/protobuf"

	"github.com/robfig/cron/v3"
)

type scheduler struct {
	sync.Mutex
	state.State
	*cron.Cron

	entries map[string][]cron.EntryID
}

func NewAutomaticScheduler(s state.State) state.AutomaticScheduler {
	return &scheduler{
		State:   s,
		Cron:    cron.New(),
		entries: make(map[string][]cron.EntryID),
	}
}

func (s *scheduler) Serve(ctx context.Context) error {
	// initial list
	configs, err := s.Repository().ListActuatorConfigs(ctx)
	if err != nil {
		return err
	}
	for _, config := range configs {
		s.handleActuatorConfigInsert(ctx, config)
	}

	unsubFn, errCh := s.DatabaseHooker().Subscribe(func(t topic.Topic, rowid int64) error {
		config, _ := s.Repository().GetActuatorConfigByRowId(ctx, rowid)
		switch t {
		case topic.Insert("actuator_configs"):
			s.handleActuatorConfigInsert(ctx, config)
		case topic.Update("actuator_configs"):
			s.handleActuatorConfigUpdate(ctx, config)
		}
		return nil
	}, topic.Insert("actuator_configs"), topic.Update("actuator_configs"))
	select {
	case <-ctx.Done():
		unsubFn()
		return ctx.Err()
	case err := <-errCh:
		unsubFn()
		return err
	}
}

func (s *scheduler) handleActuatorConfigInsert(ctx context.Context, config *pb.Config) {
	s.Lock()
	defer s.Unlock()
	s.entries[config.FeedConfig.Id] = make([]cron.EntryID, 2)
	s.entries[config.FeedConfig.Id][0], _ = s.AddFunc(config.GetActuatorConfig().TurnOffCronExpr, func() {
		s.MQTTSubscriber().Publish(ctx, config.FeedConfig.Id, []byte("0"))
	})
	s.entries[config.FeedConfig.Id][1], _ = s.AddFunc(config.GetActuatorConfig().TurnOnCronExpr, func() {
		s.MQTTSubscriber().Publish(ctx, config.FeedConfig.Id, []byte("1"))
	})
}

func (s *scheduler) handleActuatorConfigRemove(ctx context.Context, config *pb.Config) {
	s.Lock()
	defer s.Unlock()
	for _, entryID := range s.entries[config.FeedConfig.Id] {
		s.Remove(entryID)
	}
	delete(s.entries, config.FeedConfig.Id)
}

func (s *scheduler) handleActuatorConfigUpdate(ctx context.Context, config *pb.Config) {
	s.handleActuatorConfigRemove(ctx, config)
	s.handleActuatorConfigInsert(ctx, config)
}
