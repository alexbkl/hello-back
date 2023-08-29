package rds

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/redis/rueidis"
)

var rdsConn RdsConn

type RdsConn struct {
	Url      string
	Password string

	rds     *redis.Client
	jsonRds rueidis.Client
	ctx     context.Context
}

func (g *RdsConn) Open() {
	rds := redis.NewClient(&redis.Options{
		Addr:     g.Url,
		Password: g.Password, // no password set
		DB:       0,
	})

	if rds != nil {
		log.Infof("redis: %s connected!", g.Url)
	}

	jsonRds, err := rueidis.NewClient(rueidis.ClientOption{InitAddress: []string{g.Url}})
	if err != nil {
		panic(err)
	}

	var ctx = context.Background()

	g.rds = rds
	g.jsonRds = jsonRds
	g.ctx = ctx

	g.Init()
}

// SetRedisProvider sets the Gorm database connection provider.
func SetRedisProvider(conn RdsConn) {
	rdsConn = conn
}
