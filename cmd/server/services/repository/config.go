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
	if err := r.conn.QueryRowContext(ctx, `
		SELECT has_notification, lower_threshold, upper_threshold
		FROM   sensor_configs
		WHERE  feed_id = ?
	`, feedId,
	).Scan(&config.HasNotification, &config.LowerThreshold, &config.UpperThreshold); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errutils.Internal(ctx, fmt.Errorf(
			"error querying sensor config from database: %w", err))
	}

	log.Info("got sensor config from database successfully")

	return &config, nil
}

func (r configRepository) GetActuatorConfig(ctx context.Context, feedId string) (*pb.ActuatorConfig, error) {
	log := telemetry.GetLogger(ctx)

	var config pb.ActuatorConfig
	if err := r.conn.QueryRowContext(ctx, `
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

	if _, err := r.conn.ExecContext(ctx, `
		INSERT INTO sensor_configs (feed_id, has_notification, lower_threshold, upper_threshold)
		VALUES (?, ?, ?, ?)
		ON CONFLICT (feed_id) DO UPDATE SET
			has_notification = EXCLUDED.has_notification,
			lower_threshold = EXCLUDED.lower_threshold,
			upper_threshold = EXCLUDED.upper_threshold
	`, feedID, config.HasNotification, config.LowerThreshold, config.UpperThreshold); err != nil {
		return errutils.Internal(ctx, fmt.Errorf(
			"error inserting sensor config into database: %w", err))
	}

	log.Info("inserted sensor config into database successfully")

	return nil
}

func (r configRepository) UpdateActuatorConfig(ctx context.Context, feedID string, config *pb.ActuatorConfig) error {
	log := telemetry.GetLogger(ctx)

	if _, err := r.conn.ExecContext(ctx, `
		INSERT INTO actuator_configs (feed_id, automatic, turn_on_cron_expr, turn_off_cron_expr)
		VALUES (?, ?, ?, ?)
		ON CONFLICT (feed_id) DO UPDATE SET
			automatic = EXCLUDED.automatic,
			turn_on_cron_expr = EXCLUDED.turn_on_cron_expr,
			turn_off_cron_expr = EXCLUDED.turn_off_cron_expr
	`, feedID, config.Automatic, config.TurnOnCronExpr, config.TurnOffCronExpr); err != nil {
		return errutils.Internal(ctx, fmt.Errorf(
			"error inserting actuator config into database: %w", err))
	}

	log.Info("inserted actuator config into database successfully")

	return nil
}
