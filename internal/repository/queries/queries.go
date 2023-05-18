package queries

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"sync"
)

type Queries struct {
	pool *pgxpool.Pool
	mu *sync.Mutex
}

func New(pgxPool *pgxpool.Pool) *Queries {
	return &Queries{pool: pgxPool, mu: &sync.Mutex{}}
}