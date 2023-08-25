package config

import "github.com/Hello-Storage/hello-back/internal/db"

func ConnectRedis() error {
	rdsConn := db.RdsConn{
		Url:      env.RedisUrl,
		Password: env.RedisPassword,
	}

	rdsConn.Open()
	db.SetRedisProvider(rdsConn)

	return nil
}
