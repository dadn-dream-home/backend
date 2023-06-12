package services

import (
	"context"
	"sync"

	"github.com/dadn-dream-home/x/server/state"
	"github.com/dadn-dream-home/x/server/state/topic"
	"github.com/dadn-dream-home/x/server/telemetry"
	"go.uber.org/zap"

	pb "github.com/dadn-dream-home/x/protobuf"

	"github.com/go-logr/zapr"
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
		State: s,
		Cron: cron.New(cron.WithLogger(zapr.NewLogger(
			telemetry.GetLogger(telemetry.InitLogger(context.Background()))))),
		entries: make(map[string][]cron.EntryID),
	}
}

func (s *scheduler) Serve(ctx context.Context) error {
	s.Cron.Start()

	// initial list
	configs, err := s.Repository().ListActuatorConfigs(ctx)
	if err != nil {
		return err
	}
	for _, config := range configs {
		s.handleActuatorConfigInsert(ctx, config)
	}

	s.Cron.AddFunc("* * * * *", func() {
		log := telemetry.GetLogger(ctx)
		log.Info("cron job running")
	})

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
	log := telemetry.GetLogger(ctx)
	if !config.GetActuatorConfig().Automatic {
		s.entries[config.FeedConfig.Id] = make([]cron.EntryID, 0)
		return
	}
	var err error
	s.entries[config.FeedConfig.Id] = make([]cron.EntryID, 2)
	s.entries[config.FeedConfig.Id][0], err = s.AddFunc(config.GetActuatorConfig().TurnOffCronExpr, func() {
		s.MQTTSubscriber().Publish(ctx, config.FeedConfig.Id, []byte("0"))
	})
	if err != nil {
		log.Fatal("failed to add turn-off cron job", zap.Error(err))
	}
	s.entries[config.FeedConfig.Id][1], err = s.AddFunc(config.GetActuatorConfig().TurnOnCronExpr, func() {
		s.MQTTSubscriber().Publish(ctx, config.FeedConfig.Id, []byte("1"))
	})
	if err != nil {
		log.Fatal("failed to add turn-on cron job", zap.Error(err))
	}
	log.Info("added cron job", zap.String("turn-off", config.GetActuatorConfig().TurnOffCronExpr))
	log.Info("added cron job", zap.String("turn-on", config.GetActuatorConfig().TurnOnCronExpr))
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
