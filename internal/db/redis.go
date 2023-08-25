package db

import "github.com/redis/go-redis/v9"

// Rds returns the default *gorm.DB connection.
func Rds() *redis.Client {
	if rdsConn.rds == nil {
		return nil
	}

	return rdsConn.rds
}
