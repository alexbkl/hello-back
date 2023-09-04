package rds

import (
	"github.com/Hello-Storage/hello-back/internal/event"
	"github.com/redis/go-redis/v9"
)

var log = event.Log

// Rds returns the default *gorm.DB connection.
func Rds() *redis.Client {
	if rdsConn.rds == nil {
		return nil
	}

	return rdsConn.rds
}
