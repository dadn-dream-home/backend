package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/dadn-dream-home/x/server/errutils"
	"github.com/dadn-dream-home/x/server/telemetry"

	pb "github.com/dadn-dream-home/x/protobuf"
	"github.com/dadn-dream-home/x/server/state"
)

type configRepository struct {
	baseRepository
	feedRepository
}

var _ state.ConfigRepository = configRepository{}

func (r configRepository) GetFeedConfig(ctx context.Context, feedId string) (config *pb.Config, err error) {
	log := telemetry.GetLogger(ctx)

	config = &pb.Config{}

	config.FeedConfig, err = r.GetFeed(ctx, feedId)
	if err != nil {
		return nil, errutils.Internal(ctx, fmt.Errorf(
			"error getting feed from database: %w", err))
	}

	if config.FeedConfig.Type == pb.FeedType_HUMIDITY || config.FeedConfig.Type == pb.FeedType_TEMPERATURE {
		sensorConfig, err := r.GetSensorConfig(ctx, feedId)
		if err != nil {
			return nil, errutils.Internal(ctx, fmt.Errorf(
				"error getting sensor config from database: %w", err))
		}
		config.TypeConfig = &pb.Config_SensorConfig{SensorConfig: sensorConfig}
	} else {
		sensorConfig, err := r.GetActuatorConfig(ctx, feedId)
		if err != nil {
			return nil, errutils.Internal(ctx, fmt.Errorf(
				"error getting actuator config from database: %w", err))
		}
		config.TypeConfig = &pb.Config_ActuatorConfig{ActuatorConfig: sensorConfig}
	}

	log.Info("got feed config from database successfully")

	return config, nil
}

func (r configRepository) GetSensorConfig(ctx context.Context, feedId string) (*pb.SensorConfig, error) {
	log := telemetry.GetLogger(ctx)

	var config pb.SensorConfig
	var lower pb.Threshold
	var lowerFeedId sql.NullString
	var lowerFeedType sql.NullInt32
	var upper pb.Threshold
	var upperFeedId sql.NullString
	var upperFeedType sql.NullInt32
	if err := r.db.QueryRow(`
		SELECT has_notification,
			   lower_threshold, lower_has_trigger, lower_state,
			   l.id, l.type,
			   upper_threshold, upper_has_trigger, upper_state,
			   u.id, u.type
		FROM   sensor_configs
			LEFT JOIN feeds AS l on l.id = sensor_configs.lower_feed_id
			LEFT JOIN feeds AS u on u.id = sensor_configs.upper_feed_id
		WHERE  feed_id = ?
	`, feedId,
	).Scan(&config.HasNotification,
		&lower.Threshold, &lower.HasTrigger, &lower.State,
		&lowerFeedId, &lowerFeedType,
		&upper.Threshold, &upper.HasTrigger, &upper.State,
		&upperFeedId, &upperFeedType,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errutils.Internal(ctx, fmt.Errorf(
			"error querying sensor config from database: %w", err))
	}

	config.LowerThreshold = &lower
	if lowerFeedId.Valid {
		config.LowerThreshold.Feed = &pb.Feed{
			Id:   lowerFeedId.String,
			Type: pb.FeedType(lowerFeedType.Int32),
		}
	}
	config.UpperThreshold = &upper
	if upperFeedId.Valid {
		config.UpperThreshold.Feed = &pb.Feed{
			Id:   upperFeedId.String,
			Type: pb.FeedType(upperFeedType.Int32),
		}
	}

	log.Info("got sensor config from database successfully")

	return &config, nil
}

func (r configRepository) GetActuatorConfig(ctx context.Context, feedId string) (*pb.ActuatorConfig, error) {
	log := telemetry.GetLogger(ctx)

	var config pb.ActuatorConfig
	if err := r.db.QueryRow(`
		SELECT automatic, turn_on_cron_expr, turn_off_cron_expr
		FROM actuator_configs
		WHERE feed_id = ?
	`, feedId,
	).Scan(&config.Automatic, &config.TurnOnCronExpr, &config.TurnOffCronExpr); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errutils.Internal(ctx, fmt.Errorf(
			"error querying actuator config from database: %w", err))
	}

	log.Info("got actuator config from database successfully")

	return &config, nil
}

func (r configRepository) UpdateFeedConfig(ctx context.Context, config *pb.Config) error {
	log := telemetry.GetLogger(ctx)

	if err := r.updateFeed(ctx, config.FeedConfig); err != nil {
		return errutils.Internal(ctx, fmt.Errorf(
			"error setting feed into database: %w", err))
	}

	if config.TypeConfig != nil {
		switch config.FeedConfig.Type {
		case pb.FeedType_HUMIDITY, pb.FeedType_TEMPERATURE:
			if err := r.UpdateSensorConfig(ctx, config.FeedConfig.Id, config.GetSensorConfig()); err != nil {
				return err
			}
		default:
			if err := r.UpdateActuatorConfig(ctx, config.FeedConfig.Id, config.GetActuatorConfig()); err != nil {
				return err
			}
		}
	}

	log.Info("set feed config into database successfully")

	return nil
}

func (r configRepository) UpdateSensorConfig(ctx context.Context, feedID string, config *pb.SensorConfig) error {
	log := telemetry.GetLogger(ctx)

	lower := &pb.Threshold{}
	var lowerFeedID sql.NullString
	if config.LowerThreshold != nil {
		lower = config.LowerThreshold
		if lower.Feed != nil {
			lowerFeedID.String = config.LowerThreshold.Feed.Id
			lowerFeedID.Valid = true
		}
	}
	upper := &pb.Threshold{}
	var upperFeedID sql.NullString
	if config.UpperThreshold != nil {
		upper = config.UpperThreshold
		if upper.Feed != nil {
			upperFeedID.String = config.UpperThreshold.Feed.Id
			upperFeedID.Valid = true
		}
	}

	if _, err := r.db.Exec(`
		INSERT OR REPLACE INTO sensor_configs (
			feed_id, has_notification,
			lower_threshold, lower_has_trigger, lower_state,
			lower_feed_id,
			upper_threshold, upper_has_trigger, upper_state,
			upper_feed_id
		)
		VALUES (
			?, ?,
			?, ?, ?,
			?,
			?, ?, ?,
			?
		)
	`, feedID, config.HasNotification,
		lower.Threshold, lower.HasTrigger, lower.State,
		lowerFeedID,
		upper.Threshold, upper.HasTrigger, upper.State,
		upperFeedID,
	); err != nil {
		return errutils.Internal(ctx, fmt.Errorf(
			"error inserting sensor config into database: %w", err))
	}

	log.Info("inserted sensor config into database successfully")

	return nil
}

func (r configRepository) UpdateActuatorConfig(ctx context.Context, feedID string, config *pb.ActuatorConfig) error {
	log := telemetry.GetLogger(ctx)

	if _, err := r.db.Exec(`
		INSERT OR REPLACE INTO actuator_configs (feed_id, automatic, turn_on_cron_expr, turn_off_cron_expr)
		VALUES (?, ?, ?, ?)
	`, feedID, config.Automatic, config.TurnOnCronExpr, config.TurnOffCronExpr); err != nil {
		return errutils.Internal(ctx, fmt.Errorf(
			"error inserting actuator config into database: %w", err))
	}

	log.Info("inserted actuator config into database successfully")

	return nil
}

func (r configRepository) ListActuatorConfigs(ctx context.Context) ([]*pb.Config, error) {
	log := telemetry.GetLogger(ctx)

	rows, err := r.db.Query(`
		SELECT id, type, automatic, turn_on_cron_expr, turn_off_cron_expr
		FROM actuator_configs, feeds
		WHERE feed_id = id
	`)
	if err != nil {
		return nil, errutils.Internal(ctx, fmt.Errorf(
			"error querying actuator configs from database: %w", err))
	}
	defer rows.Close()

	var configs []*pb.Config
	for rows.Next() {
		actuatorConfig := &pb.ActuatorConfig{}
		feed := &pb.Feed{}
		if err := rows.Scan(&feed.Id, &feed.Type, &actuatorConfig.Automatic, &actuatorConfig.TurnOnCronExpr, &actuatorConfig.TurnOffCronExpr); err != nil {
			return nil, errutils.Internal(ctx, fmt.Errorf(
				"error scanning actuator config from database: %w", err))
		}
		config := &pb.Config{}
		config.TypeConfig = &pb.Config_ActuatorConfig{ActuatorConfig: actuatorConfig}
		config.FeedConfig = feed
		configs = append(configs, config)
	}

	log.Info("got actuator configs from database successfully")

	return configs, nil
}

func (r configRepository) GetActuatorConfigByRowId(ctx context.Context, rowId int64) (*pb.Config, error) {
	log := telemetry.GetLogger(ctx)

	var feedId string
	var feedType int32
	var automatic bool
	var turnOnCronExpr string
	var turnOffCronExpr string
	if err := r.db.QueryRow(`
		SELECT a.feed_id, f.type, a.automatic, a.turn_on_cron_expr, a.turn_off_cron_expr
		FROM actuator_configs AS a, feeds AS f
		WHERE a.feed_id = f.id AND a.rowid = ?
	`, rowId,
	).Scan(&feedId, &feedType, &automatic, &turnOnCronExpr, &turnOffCronExpr); err != nil {
		return nil, errutils.Internal(ctx, fmt.Errorf(
			"error querying actuator config from database: %w", err))
	}

	config := &pb.Config{
		FeedConfig: &pb.Feed{
			Id:   feedId,
			Type: pb.FeedType(feedType),
		},
		TypeConfig: &pb.Config_ActuatorConfig{
			ActuatorConfig: &pb.ActuatorConfig{
				Automatic:       automatic,
				TurnOnCronExpr:  turnOnCronExpr,
				TurnOffCronExpr: turnOffCronExpr,
			},
		},
	}

	log.Info("got actuator config from database successfully")

	return config, nil
}