package config

import "github.com/Hello-Storage/hello-back/internal/rds"

func ConnectRedis() error {
	rdsConn := rds.RdsConn{
		Url:      env.RedisUrl,
		Password: "", // env.RedisPassword
	}

	rdsConn.Open()
	rds.SetRedisProvider(rdsConn)

	return nil
}
