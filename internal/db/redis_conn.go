package db

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var rdsConn RdsConn

type RdsConn struct {
	url      string
	password string

	rds *redis.Client
	ctx context
}

func (g *RdsConn) Open() {
	rds := redis.NewClient(&redis.Options{
		Addr:     g.url,
		Password: g.password, // no password set
		DB:       0,
	})

	var ctx = context.Background()
}
