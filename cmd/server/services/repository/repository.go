package repository

import (
	"context"
	"database/sql"

	"github.com/dadn-dream-home/x/server/state"
)

type baseRepository struct {
	db *sql.DB
}

type repository struct {
	baseRepository
	feedRepository
	feedValueRepository
	notificationRepository
	configRepository
}

func NewRepository(ctx context.Context, state state.State, db *sql.DB) state.Repository {
	r := repository{}
	r.db = db
	r.feedRepository = feedRepository{r.baseRepository}
	r.feedValueRepository = feedValueRepository{r.baseRepository}
	r.notificationRepository = notificationRepository{r.baseRepository}
	r.configRepository = configRepository{r.baseRepository, r.feedRepository}
	return r
}
