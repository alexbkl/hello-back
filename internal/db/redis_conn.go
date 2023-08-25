package db

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var rdsConn RdsConn

type RdsConn struct {
	Url      string
	Password string

	rds *redis.Client
	ctx context.Context
}

func (g *RdsConn) Open() {
	rds := redis.NewClient(&redis.Options{
		Addr:     g.Url,
		Password: g.Password, // no password set
		DB:       0,
	})

	var ctx = context.Background()

	g.rds = rds
	g.ctx = ctx
}

// SetRedisProvider sets the Gorm database connection provider.
func SetRedisProvider(conn RdsConn) {
	rdsConn = conn
}
