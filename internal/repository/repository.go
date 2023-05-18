package repository

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
	"context"
	"postgreintern/internal/model"
	"postgreintern/internal/repository/queries"
)

type repo struct {
	logger *logrus.Logger
	pool   *pgxpool.Pool
	*queries.Queries
}

// инициализация подключения к бд
func dbSetup(logger *logrus.Logger) *pgxpool.Pool {
	dsn := "postgres://admin:123@localhost:5442/events?sslmode=disable"

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		logger.Fatalln(err)
	}

	pool, err := pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		logger.Fatalln(err)
	}

	return pool
}

func NewRepo(logger *logrus.Logger) Repository {
	pool := dbSetup(logger)

	return &repo{
		logger: logger,
		pool:   dbSetup(logger),
		Queries: queries.New(pool),
	}
}

type Repository interface {
	// SaveEvent сохраняет событие в бд и возвращает id события
	SaveEvent(ctx context.Context, event model.FileEvent) (int, error)

	// SaveCommandExecution сохраняет выполнение команды в бд с id событий как foreign key
	SaveCommandExecution(ctx context.Context, event model.CmdEvent, id int) error
}